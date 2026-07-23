#!/bin/bash
# Seed multiple domains and mailboxes for local lab testing.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LAB_DIR="${SCRIPT_DIR}/../lab"
DOMAINS_FILE="${LAB_DIR}/domains.txt"
MAILBOXES_FILE="${LAB_DIR}/mailboxes.txt"
SNAPPY_TEMPLATE="${SCRIPT_DIR}/../snappymail/domains/_default.json"

pfa() {
  docker compose exec -T postfixadmin ./postfixadmin "$@" < /dev/null
}

echo "==> Creating lab domains..."
while IFS=$'\t' read -r domain description _; do
  [[ -z "${domain}" || "${domain}" =~ ^# ]] && continue
  pfa domain --add "${domain}" \
    --description "${description:-Lab domain}" \
    --max-aliases 200 \
    --max-mailboxes 100 2>/dev/null || true
  echo "    domain: ${domain}"
done < "${DOMAINS_FILE}"

echo "==> Creating lab mailboxes..."
while IFS= read -r line; do
  line="${line%%#*}"
  line="$(echo "${line}" | xargs)"
  [[ -z "${line}" ]] && continue
  pfa mailbox --add "${line}" 2>/dev/null || true
  echo "    mailbox: ${line%%:*}"
done < "${MAILBOXES_FILE}"

echo "==> Ensuring local transport..."
pfa transport --add "local:lmtp:unix:private/dovecot-lmtp" 2>/dev/null || true

echo "==> Seeding welcome messages..."
while IFS= read -r line; do
  line="${line%%#*}"
  line="$(echo "${line}" | xargs)"
  [[ -z "${line}" ]] && continue
  email="${line%%:*}"
  docker compose exec -T mailserver bash -c \
    "echo 'Welcome to the go-snappymail lab' | mail -s 'Lab mailbox ready' '${email}'" \
    < /dev/null 2>/dev/null || true
done < "${MAILBOXES_FILE}"

echo "==> Configuring SnappyMail domains..."
for i in $(seq 1 30); do
  if docker compose ps snappymail --format '{{.Status}}' 2>/dev/null | grep -qi up; then
    break
  fi
  sleep 2
done

docker compose exec -T snappymail mkdir -p /var/lib/snappymail/_data_/_default_/domains 2>/dev/null || true
while IFS=$'\t' read -r domain _; do
  [[ -z "${domain}" || "${domain}" =~ ^# ]] && continue
  docker cp "${SNAPPY_TEMPLATE}" "gosm-snappymail:/var/lib/snappymail/_data_/_default_/domains/${domain}.json" 2>/dev/null || true
  echo "    snappymail: ${domain}.json"
done < "${DOMAINS_FILE}"

docker compose restart snappymail 2>/dev/null || true

count="$(grep -vE '^\s*(#|$)' "${MAILBOXES_FILE}" | wc -l | tr -d ' ')"
echo "==> Lab seed complete (${count} mailboxes)"
