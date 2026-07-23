#!/usr/bin/env bash
# Docker + compose plugin, e stack em containers: MariaDB + SnappyMail PHP.
# MariaDB publica 3306 no loopback da VM — Postfix/Dovecot nativos usam 127.0.0.1.
set -euo pipefail
source /vagrant/provision/common.sh

log "Instalando Docker + compose plugin"

if ! command -v docker &>/dev/null; then
  install -m 0755 -d /etc/apt/keyrings
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor --batch --yes -o /etc/apt/keyrings/docker.gpg
  chmod a+r /etc/apt/keyrings/docker.gpg

  echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] \
    https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" \
    > /etc/apt/sources.list.d/docker.list

  apt-get update -qq
  apt-get install -y -qq docker-ce docker-ce-cli containerd.io docker-compose-plugin
fi

usermod -aG docker vagrant || true
systemctl enable --now docker

log "Escrevendo ${COMPOSE_DIR}/docker-compose.yml (mariadb + snappymail)"
mkdir -p "$COMPOSE_DIR"

cat > "${COMPOSE_DIR}/docker-compose.yml" <<EOF
# Infra em containers do lab híbrido go-snappymail.
# Postfix/Dovecot e binários Go rodam NATIVOS na VM (systemd).
services:
  mariadb:
    image: mariadb:10.11
    container_name: gosm-mariadb
    restart: unless-stopped
    environment:
      MARIADB_ROOT_PASSWORD: "${MARIADB_ROOT_PASS}"
      MARIADB_DATABASE: "${DB_NAME}"
      MARIADB_USER: "${DB_USER}"
      MARIADB_PASSWORD: "${DB_PASS}"
    ports:
      - "127.0.0.1:3306:3306"
    volumes:
      - mariadb_data:/var/lib/mysql
    healthcheck:
      test: ["CMD", "mariadb-admin", "ping", "-h", "localhost", "-u", "root", "-p${MARIADB_ROOT_PASS}"]
      interval: 5s
      timeout: 5s
      retries: 10

  snappymail:
    image: ${SNAPPYMAIL_IMAGE}
    container_name: gosm-snappymail
    restart: unless-stopped
    ports:
      - "8888:8888"
    volumes:
      - snappymail_data:/var/lib/snappymail
    extra_hosts:
      - "${MAIL_FQDN}:host-gateway"

volumes:
  mariadb_data:
  snappymail_data:
EOF

log "Subindo containers (docker compose up -d)"
cd "$COMPOSE_DIR"
docker compose up -d

wait_for_mariadb

log "Docker pronto: mariadb (127.0.0.1:3306) + snappymail (:8888)"
