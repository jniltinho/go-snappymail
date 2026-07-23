#!/usr/bin/env bash
# Binários Go NATIVOS na VM, cada um com unit systemd:
#   go-snappymail :8082 (buildado do repo sincronizado)
#   go-cubemail   :8080 (release binária)
#   go-postfixadmin :8081 (instalado em 03-mailserver.sh)
set -euo pipefail
source /vagrant/provision/common.sh

# --- Go toolchain ---
if [[ ! -x /usr/local/go/bin/go ]] || ! /usr/local/go/bin/go version | grep -q "go${GO_VERSION}"; then
  log "Instalando Go ${GO_VERSION}"
  wget -q -O /tmp/go.tar.gz "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
  rm -rf /usr/local/go
  tar -C /usr/local -xzf /tmp/go.tar.gz
fi
export PATH="/usr/local/go/bin:$PATH"
grep -q '/usr/local/go/bin' /etc/profile.d/golang.sh 2>/dev/null || \
  echo 'export PATH=/usr/local/go/bin:$PATH' > /etc/profile.d/golang.sh

# --- go-snappymail (build do repo) ---
log "Buildando go-snappymail a partir de ${REPO_DIR}"

if [[ ! -f "${REPO_DIR}/web/dist/index.html" ]]; then
  echo "ERRO: ${REPO_DIR}/web/dist/index.html não existe." >&2
  echo "Rode 'make frontend' no host (fora da VM) e provisione de novo." >&2
  exit 1
fi

mkdir -p /opt/go-snappymail/data /opt/go-snappymail/tmp/uploads /etc/go-snappymail
cd "$REPO_DIR"
CGO_ENABLED=1 go build -ldflags="-s -w" -o /opt/go-snappymail/go-snappymail .

if [[ ! -f /etc/go-snappymail/config.toml ]]; then
  SECRET_KEY=$(openssl rand -hex 16)
  cat > /etc/go-snappymail/config.toml <<EOF
[server]
host           = "0.0.0.0"
port           = 8082
debug          = false
secret_key     = "${SECRET_KEY}"
base_url       = "http://${VM_IP}:8082"
tls_cert       = ""
tls_key        = ""
swagger_enable = false

[imap]
host            = "127.0.0.1"
port            = 993
tls             = true
timeout_sec     = 30
show_host_input = false
tls_server_name = "${MAIL_FQDN}"
insecure_skip_verify = true   # lab only — Dovecot usa cert autoassinado

[smtp]
host        = "127.0.0.1"
port        = 25
starttls    = false
timeout_sec = 30

[database]
driver = "sqlite"
dsn    = "/opt/go-snappymail/data/app.db"
debug  = false

[session]
name      = "gsn_session"
max_age   = 86400
secure    = false
http_only = true

[ui]
theme           = "snappymail-default"
rows_per_page   = 50
timezone        = "America/Sao_Paulo"
date_format     = "02/01/2006"
datetime_format = "02/01/2006 15:04"
compose_html    = true

[upload]
max_size_mb = 25
temp_dir    = "/opt/go-snappymail/tmp/uploads"
EOF
fi

/opt/go-snappymail/go-snappymail migrate --config /etc/go-snappymail/config.toml

cat > /etc/systemd/system/go-snappymail.service <<EOF
[Unit]
Description=go-snappymail webmail (este repo)
After=network.target docker.service dovecot.service postfix.service

[Service]
Type=simple
WorkingDirectory=/opt/go-snappymail
ExecStart=/opt/go-snappymail/go-snappymail serve --config /etc/go-snappymail/config.toml
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# --- go-cubemail (referência Golang, release binária) ---
log "Instalando go-cubemail v${GO_CUBEMAIL_VERSION}"

INSTALL_DIR="/opt/go-cubemail"
mkdir -p "$INSTALL_DIR/data" "$INSTALL_DIR/tmp/uploads"

TAR="/tmp/go-cubemail_${GO_CUBEMAIL_VERSION}_linux_amd64.tar.gz"
if [[ ! -f "$TAR" ]]; then
  wget -q -O "$TAR" \
    "https://github.com/jniltinho/go-cubemail/releases/download/v${GO_CUBEMAIL_VERSION}/go-cubemail_${GO_CUBEMAIL_VERSION}_linux_amd64.tar.gz"
fi

tar -xzf "$TAR" -C "$INSTALL_DIR"
chmod +x "$INSTALL_DIR/go-cubemail"

if [[ ! -f "$INSTALL_DIR/config.toml" ]]; then
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
fi

cd "$INSTALL_DIR"
./go-cubemail migrate

cat > /etc/systemd/system/go-cubemail.service <<EOF
[Unit]
Description=Go CubeMail (referência Golang do lab go-snappymail)
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
systemctl enable go-snappymail go-cubemail
systemctl restart go-snappymail go-cubemail

log "Apps Go rodando: go-snappymail :8082, go-cubemail :8080"
