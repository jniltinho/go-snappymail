#!/usr/bin/env bash
set -euo pipefail
source /vagrant/provision/common.sh

log "Instalando Go-PostfixAdmin e configurando Postfix/Dovecot"

# --- Go-PostfixAdmin ---
DEB="/tmp/go-postfixadmin_${GO_POSTFIXADMIN_VERSION}_amd64.deb"
if [[ ! -f "$DEB" ]]; then
  wget -q -O "$DEB" \
    "https://github.com/jniltinho/go-postfixadmin/releases/download/v${GO_POSTFIXADMIN_VERSION}/go-postfixadmin_${GO_POSTFIXADMIN_VERSION}_amd64.deb"
fi
dpkg -i "$DEB" || apt-get install -f -y -qq

SESSION_SECRET=$(openssl rand -hex 32)

cat > /opt/go-postfixadmin/config.toml <<EOF
[database]
host   = "localhost"
port   = "3306"
user   = "${DB_USER}"
pass   = "${DB_PASS}"
name   = "${DB_NAME}"
driver = "mysql"
debug  = false

[server]
port            = 8081
cleanup_maildir = false
ssl_enable      = false
session_secret  = "${SESSION_SECRET}"
swagger_enable  = true

jwt_access_ttl  = "15m"
jwt_refresh_ttl = "168h"

[quota]
enabled      = false
domain_quota = true
multiplier   = 1048576

[vacation]
enabled = true

[alias]
edit_alias   = true
alias_domain = true

[transport]
host          = "127.0.0.1:12221"
cache         = "10m"
hostname      = "${MAIL_FQDN}"
localdelivery = "smtp:${MAIL_FQDN}"
delivery      = "lmtp:unix:private/dovecot-lmtp"

[features]
fetchmail = false

[rbac]
enabled = false

[smtp]
server = "localhost"
port   = 25
type   = "plain"
EOF

cd /opt/go-postfixadmin
PA_BIN="./postfixadmin"
[[ -x /opt/go-postfixadmin/postfixadmin ]] && PA_BIN="/opt/go-postfixadmin/postfixadmin"

$PA_BIN migrate

# Bootstrap domain, admin, mailbox, transport
$PA_BIN admin --add-superadmin "${ADMIN_EMAIL}:${ADMIN_PASS}" || true
$PA_BIN domain --add "${MAIL_DOMAIN}" \
  --description "Lab domain" \
  --max-aliases 100 \
  --max-mailboxes 50 || true
$PA_BIN mailbox --add "${MAIL_USER}:${MAIL_PASS}" || true
$PA_BIN transport --add "local:lmtp:unix:private/dovecot-lmtp" || true

systemctl enable postfixadmin
systemctl restart postfixadmin

# --- Postfix SQL maps ---
mkdir -p /etc/postfix/sql
for f in mysql_virtual_domains_maps.cf mysql_virtual_mailbox_maps.cf \
         mysql_virtual_alias_maps.cf mysql_virtual_alias_domain_maps.cf \
         mysql_virtual_alias_domain_catchall_maps.cf mysql_virtual_alias_domain_mailbox_maps.cf; do
  render_template "/vagrant/templates/postfix/sql/${f}" "/etc/postfix/sql/${f}"
done
chmod 640 /etc/postfix/sql/*.cf
chown root:postfix /etc/postfix/sql/*.cf

render_template /vagrant/templates/postfix/main.cf /etc/postfix/main.cf

# --- Dovecot ---
render_template /vagrant/templates/dovecot/dovecot.conf /etc/dovecot/dovecot.conf
render_template /vagrant/templates/dovecot/dovecot-sql.conf /etc/dovecot/dovecot-sql.conf
chmod 640 /etc/dovecot/dovecot-sql.conf
chown root:dovecot /etc/dovecot/dovecot-sql.conf

# Disable split config includes that override our single-file config
for f in /etc/dovecot/conf.d/*.conf; do
  [[ -f "$f" ]] && mv "$f" "${f}.disabled" 2>/dev/null || true
done

systemctl enable postfix dovecot
systemctl restart postfix dovecot

# Send test email to seed INBOX
sleep 3
echo "Test message from go-snappymail lab" | mail -s "Welcome to ${MAIL_DOMAIN}" "${MAIL_USER}" || true

log "Mailserver configurado (Postfix + Dovecot + Go-PostfixAdmin)"
