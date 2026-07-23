#!/usr/bin/env bash
# Shared environment for go-snappymail lab VM
set -euo pipefail

_VAGRANT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
if [[ -f "${_VAGRANT_DIR}/.env" ]]; then
  set -a
  # shellcheck disable=SC1091
  source "${_VAGRANT_DIR}/.env"
  set +a
fi

export MAIL_DOMAIN="${MAIL_DOMAIN:-test.local}"
export MAIL_HOST="${MAIL_HOST:-mail.test.local}"
export MAIL_FQDN="${MAIL_FQDN:-mail.test.local}"
export VM_IP="${VM_IP:-192.168.56.20}"

export DB_NAME="${DB_NAME:-postfix}"
export DB_USER="${DB_USER:-postfix}"
export DB_PASS="${DB_PASS:-postfixPassword}"
export MARIADB_ROOT_PASS="${MARIADB_ROOT_PASS:-rootpassword}"

export MAIL_USER="${MAIL_USER:-user@test.local}"
export MAIL_PASS="${MAIL_PASS:-Password1@}"
export ADMIN_EMAIL="${ADMIN_EMAIL:-admin@test.local}"
export ADMIN_PASS="${ADMIN_PASS:-Password1@}"

export VMAIL_UID=1001
export VMAIL_GID=1001

export GO_POSTFIXADMIN_VERSION="1.0.86"
export GO_CUBEMAIL_VERSION="0.0.25"
export SNAPPYMAIL_VERSION="2.38.2"
export SNAPPY_ADMIN_PASS="Admin1@lab"

export SSL_CERT="/etc/ssl/certs/mail.test.local.crt"
export SSL_KEY="/etc/ssl/private/mail.test.local.key"

export DEBIAN_FRONTEND=noninteractive

log() { echo "==> $*"; }

render_template() {
  local src="$1" dest="$2"
  sed \
    -e "s|@@MAIL_DOMAIN@@|${MAIL_DOMAIN}|g" \
    -e "s|@@MAIL_HOST@@|${MAIL_HOST}|g" \
    -e "s|@@MAIL_FQDN@@|${MAIL_FQDN}|g" \
    -e "s|@@DB_PASS@@|${DB_PASS}|g" \
    -e "s|@@DB_USER@@|${DB_USER}|g" \
    -e "s|@@DB_NAME@@|${DB_NAME}|g" \
    -e "s|@@VMAIL_UID@@|${VMAIL_UID}|g" \
    -e "s|@@VMAIL_GID@@|${VMAIL_GID}|g" \
    -e "s|@@SSL_CERT@@|${SSL_CERT}|g" \
    -e "s|@@SSL_KEY@@|${SSL_KEY}|g" \
    "$src" > "$dest"
}

ensure_self_signed_cert() {
  if [[ ! -f "$SSL_CERT" ]]; then
    log "Gerando certificado autoassinado para ${MAIL_FQDN}"
    openssl req -x509 -nodes -days 3650 -newkey rsa:2048 \
      -keyout "$SSL_KEY" \
      -out "$SSL_CERT" \
      -subj "/CN=${MAIL_FQDN}/O=go-snappymail-lab/C=BR"
    chmod 600 "$SSL_KEY"
  fi
}
