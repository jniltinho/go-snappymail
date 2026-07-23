#!/usr/bin/env bash
# Seed do lab: domínios, mailboxes e aliases de docker/lab/*.txt.
# go-postfixadmin roda nativo; aliases via mysql client contra o container (127.0.0.1).
set -euo pipefail
source /vagrant/provision/common.sh

LAB_DIR="${REPO_DIR}/docker/lab"
DOMAINS_FILE="${LAB_DIR}/domains.txt"
MAILBOXES_FILE="${LAB_DIR}/mailboxes.txt"
ALIASES_FILE="${LAB_DIR}/aliases.txt"

cd /opt/go-postfixadmin
pfa() { /opt/go-postfixadmin/postfixadmin "$@"; }

log "Criando domínios do lab..."
while IFS=$'\t' read -r domain description _; do
  [[ -z "${domain}" || "${domain}" =~ ^# ]] && continue
  pfa domain --add "${domain}" \
    --description "${description:-Lab domain}" \
    --max-aliases 200 \
    --max-mailboxes 100 2>/dev/null || true
  echo "    domain: ${domain}"
done < "${DOMAINS_FILE}"

log "Criando mailboxes do lab..."
while IFS= read -r line; do
  line="${line%%#*}"
  line="$(echo "${line}" | xargs)"
  [[ -z "${line}" ]] && continue
  pfa mailbox --add "${line}:${MAIL_PASS}" 2>/dev/null || true
  echo "    mailbox: ${line}"
done < "${MAILBOXES_FILE}"

log "Criando aliases do lab..."
if [[ -f "${ALIASES_FILE}" ]]; then
  while IFS=$'\t' read -r addr goto _; do
    [[ -z "${addr}" || "${addr}" =~ ^# ]] && continue
    dom="${addr#*@}"
    mysql -h"${DB_HOST}" -u"${DB_USER}" -p"${DB_PASS}" "${DB_NAME}" -e \
      "INSERT IGNORE INTO alias (address, goto, domain, created, modified, active) VALUES ('${addr}', '${goto}', '${dom}', NOW(), NOW(), 1);" \
      2>/dev/null || true
    echo "    alias: ${addr} -> ${goto}"
  done < "${ALIASES_FILE}"
fi

log "Garantindo transport local..."
pfa transport --add "local:lmtp:unix:private/dovecot-lmtp" 2>/dev/null || true

log "Enviando mensagens de boas-vindas..."
while IFS= read -r line; do
  line="${line%%#*}"
  line="$(echo "${line}" | xargs)"
  [[ -z "${line}" ]] && continue
  echo 'Welcome to the go-snappymail lab' | mail -s 'Lab mailbox ready' "${line}" || true
done < "${MAILBOXES_FILE}"

log "Configurando domínios no SnappyMail (container)..."
DOMAIN_JSON="$(mktemp)"
cat > "$DOMAIN_JSON" <<JSON
{
    "IMAP": {
        "host": "${MAIL_FQDN}",
        "port": 993,
        "secure": 1,
        "shortLogin": 0,
        "ssl": {
            "verify_peer": false,
            "verify_peer_name": false,
            "allow_self_signed": true
        }
    },
    "SMTP": {
        "host": "${MAIL_FQDN}",
        "port": 25,
        "secure": 0,
        "shortLogin": 0,
        "useAuth": true
    },
    "whiteList": ""
}
JSON

docker exec gosm-snappymail mkdir -p /var/lib/snappymail/_data_/_default_/domains 2>/dev/null || true
while IFS=$'\t' read -r domain _; do
  [[ -z "${domain}" || "${domain}" =~ ^# ]] && continue
  docker cp "$DOMAIN_JSON" "gosm-snappymail:/var/lib/snappymail/_data_/_default_/domains/${domain}.json" 2>/dev/null || true
  echo "    snappymail: ${domain}.json"
done < "${DOMAINS_FILE}"
rm -f "$DOMAIN_JSON"

cd "$COMPOSE_DIR" && docker compose restart snappymail >/dev/null 2>&1 || true

count="$(grep -vE '^\s*(#|$)' "${MAILBOXES_FILE}" | wc -l | tr -d ' ')"
log "Seed completo (${count} mailboxes, senha padrão: ${MAIL_PASS})"
