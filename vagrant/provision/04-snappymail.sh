#!/usr/bin/env bash
set -euo pipefail
source /vagrant/provision/common.sh

log "Instalando SnappyMail ${SNAPPYMAIL_VERSION}"

TAR="/tmp/snappymail-${SNAPPYMAIL_VERSION}.tar.gz"
if [[ ! -f "$TAR" ]]; then
  wget -q -O "$TAR" \
    "https://github.com/the-djmaze/snappymail/releases/download/v${SNAPPYMAIL_VERSION}/snappymail-${SNAPPYMAIL_VERSION}.tar.gz"
fi

rm -rf /var/www/snappymail
mkdir -p /var/www/snappymail /var/lib/snappymail
tar -xzf "$TAR" -C /var/www/snappymail

# Data directory outside web root
if [[ -d /var/www/snappymail/data ]] && [[ ! -L /var/www/snappymail/data ]]; then
  cp -a /var/www/snappymail/data/. /var/lib/snappymail/ 2>/dev/null || true
  rm -rf /var/www/snappymail/data
fi
ln -sfn /var/lib/snappymail /var/www/snappymail/data

chown -R www-data:www-data /var/www/snappymail /var/lib/snappymail
find /var/www/snappymail -type d -exec chmod 755 {} \;
find /var/www/snappymail -type f -exec chmod 644 {} \;
chmod 640 /var/www/snappymail/include.php 2>/dev/null || true

render_template /vagrant/templates/nginx/snappymail.conf /etc/nginx/sites-available/snappymail.conf
ln -sf /etc/nginx/sites-available/snappymail.conf /etc/nginx/sites-enabled/snappymail.conf
rm -f /etc/nginx/sites-enabled/default

systemctl enable nginx php8.3-fpm
systemctl restart php8.3-fpm nginx

# Initialize SnappyMail (creates data dirs)
curl -s http://127.0.0.1:8888/ >/dev/null || true

# Set known admin password for lab
ADMIN_HASH=$(php -r "echo password_hash('${SNAPPY_ADMIN_PASS}', PASSWORD_DEFAULT);")
APP_INI="/var/lib/snappymail/_data_/_default_/configs/application.ini"
if [[ -f "$APP_INI" ]]; then
  sed -i "s|^admin_password = .*|admin_password = \"${ADMIN_HASH}\"|" "$APP_INI"
  echo "${SNAPPY_ADMIN_PASS}" > /root/snappymail-admin-password.txt
  chmod 600 /root/snappymail-admin-password.txt
  chown www-data:www-data "$APP_INI"
fi

# Pre-configure IMAP domain for test.local
DOMAIN_JSON="/var/lib/snappymail/_data_/_default_/domains/${MAIL_DOMAIN}.json"
if [[ ! -f "$DOMAIN_JSON" ]]; then
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
  chown www-data:www-data "$DOMAIN_JSON"
fi

log "SnappyMail disponível em http://${VM_IP}:8888 (admin: admin / ${SNAPPY_ADMIN_PASS})"
