#!/usr/bin/env bash
set -euo pipefail
source /vagrant/provision/common.sh

log "Instalando Go webmail (go-cubemail v${GO_CUBEMAIL_VERSION} — referência até go-snappymail existir)"

INSTALL_DIR="/opt/go-cubemail"
mkdir -p "$INSTALL_DIR/data" "$INSTALL_DIR/tmp/uploads"

TAR="/tmp/go-cubemail_${GO_CUBEMAIL_VERSION}_linux_amd64.tar.gz"
if [[ ! -f "$TAR" ]]; then
  wget -q -O "$TAR" \
    "https://github.com/jniltinho/go-cubemail/releases/download/v${GO_CUBEMAIL_VERSION}/go-cubemail_${GO_CUBEMAIL_VERSION}_linux_amd64.tar.gz"
fi

tar -xzf "$TAR" -C "$INSTALL_DIR"
chmod +x "$INSTALL_DIR/go-cubemail"

SECRET_KEY=$(openssl rand -hex 16)

cat > "$INSTALL_DIR/config.toml" <<EOF
[server]
host           = "0.0.0.0"
port           = 8080
debug          = false
secret_key     = "${SECRET_KEY}"
base_url       = "http://${VM_IP}:8080"
tls_cert       = ""
tls_key        = ""
swagger_enable = true

[imap]
host            = "${MAIL_FQDN}"
port            = 993
tls             = true
timeout_sec     = 30
show_host_input = false

[smtp]
host        = "${MAIL_FQDN}"
port        = 25
starttls    = false
timeout_sec = 30

[database]
driver = "sqlite"
dsn    = "${INSTALL_DIR}/data/app.db"
debug  = false

[session]
name      = "gorc_session"
max_age   = 86400
secure    = false
http_only = true

[ui]
theme           = "larry"
rows_per_page   = 50
timezone        = "America/Sao_Paulo"
date_format     = "02/01/2006"
datetime_format = "02/01/2006 15:04"
compose_html    = true

[upload]
max_size_mb = 25
temp_dir    = "${INSTALL_DIR}/tmp/uploads"

[activesync]
enabled = false

[push]
vapid_public_key  = ""
vapid_private_key = ""
vapid_contact     = "mailto:${ADMIN_EMAIL}"
EOF

cd "$INSTALL_DIR"
./go-cubemail migrate

cat > /etc/systemd/system/go-cubemail.service <<EOF
[Unit]
Description=Go CubeMail (Golang webmail reference for go-snappymail lab)
After=network.target dovecot.service postfix.service

[Service]
Type=simple
WorkingDirectory=${INSTALL_DIR}
ExecStart=${INSTALL_DIR}/go-cubemail serve
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable go-cubemail
systemctl restart go-cubemail

log "Go webmail disponível em http://${VM_IP}:8080"
