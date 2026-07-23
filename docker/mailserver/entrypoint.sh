#!/bin/bash
set -euo pipefail

render() {
  local src="$1" dest="$2"
  sed \
    -e "s|@@MAIL_DOMAIN@@|${MAIL_DOMAIN}|g" \
    -e "s|@@MAIL_FQDN@@|${MAIL_FQDN}|g" \
    -e "s|@@DB_HOST@@|${DB_HOST}|g" \
    -e "s|@@DB_PASS@@|${DB_PASS}|g" \
    -e "s|@@DB_USER@@|${DB_USER}|g" \
    -e "s|@@DB_NAME@@|${DB_NAME}|g" \
    -e "s|@@VMAIL_UID@@|${VMAIL_UID}|g" \
    -e "s|@@VMAIL_GID@@|${VMAIL_GID}|g" \
    -e "s|@@SSL_CERT@@|${SSL_CERT}|g" \
    -e "s|@@SSL_KEY@@|${SSL_KEY}|g" \
    "$src" > "$dest"
}

DB_HOST="${DB_HOST:-mariadb}"
DB_PORT="${DB_PORT:-3306}"
MAIL_DOMAIN="${MAIL_DOMAIN:-test.local}"
MAIL_FQDN="${MAIL_FQDN:-mail.test.local}"
DB_USER="${DB_USER:-postfix}"
DB_PASS="${DB_PASS:-postfixPassword}"
DB_NAME="${DB_NAME:-postfix}"
VMAIL_UID="${VMAIL_UID:-1001}"
VMAIL_GID="${VMAIL_GID:-1001}"

SSL_DIR="/etc/ssl/mail"
SSL_CERT="${SSL_DIR}/mail.crt"
SSL_KEY="${SSL_DIR}/mail.key"
mkdir -p "$SSL_DIR"
if [[ ! -f "$SSL_CERT" ]]; then
  openssl req -x509 -nodes -days 3650 -newkey rsa:2048 \
    -keyout "$SSL_KEY" -out "$SSL_CERT" \
    -subj "/CN=${MAIL_FQDN}/O=go-snappymail-docker/C=BR" \
    -addext "subjectAltName=DNS:${MAIL_FQDN},DNS:mailserver,DNS:localhost"
fi

echo "Waiting for MariaDB at ${DB_HOST}:${DB_PORT}..."
until mariadb-admin ping -h"$DB_HOST" -u"$DB_USER" -p"$DB_PASS" --silent 2>/dev/null; do
  sleep 2
done

for f in mysql_virtual_domains_maps.cf mysql_virtual_mailbox_maps.cf \
         mysql_virtual_alias_maps.cf mysql_virtual_alias_domain_maps.cf \
         mysql_virtual_alias_domain_catchall_maps.cf mysql_virtual_alias_domain_mailbox_maps.cf; do
  render "/templates/postfix/sql/${f}" "/etc/postfix/sql/${f}"
done
chmod 640 /etc/postfix/sql/*.cf
chown root:postfix /etc/postfix/sql/*.cf

render /templates/postfix/main.cf /etc/postfix/main.cf
render /templates/dovecot/dovecot.conf /etc/dovecot/dovecot.conf
render /templates/dovecot/dovecot-sql.conf /etc/dovecot/dovecot-sql.conf
chmod 640 /etc/dovecot/dovecot-sql.conf
chown root:dovecot /etc/dovecot/dovecot-sql.conf

for f in /etc/dovecot/conf.d/*.conf; do
  [[ -f "$f" ]] && mv "$f" "${f}.disabled" 2>/dev/null || true
done

chown -R vmail:vmail /var/vmail

exec "$@"
