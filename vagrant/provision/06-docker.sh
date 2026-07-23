#!/usr/bin/env bash
# Install Docker and start the containerized lab stack.
# Stops bare-metal services first to free ports 25/143/993/8080/8081/8888.
set -euo pipefail
source /vagrant/provision/common.sh

log "Instalando Docker"

if ! command -v docker &>/dev/null; then
  install -m 0755 -d /etc/apt/keyrings
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor --batch --yes -o /etc/apt/keyrings/docker.gpg 2>/dev/null || \
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
  chmod a+r /etc/apt/keyrings/docker.gpg

  echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] \
    https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" \
    > /etc/apt/sources.list.d/docker.list

  apt-get update -qq
  apt-get install -y -qq docker-ce docker-ce-cli containerd.io docker-compose-plugin
fi

usermod -aG docker vagrant || true
systemctl enable docker 2>/dev/null || true
systemctl start docker

log "Parando serviços bare-metal (liberando portas para containers)"
for svc in go-cubemail postfixadmin nginx php8.3-fpm postfix dovecot; do
  systemctl stop "$svc" 2>/dev/null || true
  systemctl disable "$svc" 2>/dev/null || true
done

LAB_DIR="/opt/go-snappymail-docker"
rm -rf "$LAB_DIR"
mkdir -p "$LAB_DIR"
cp -a /vagrant/go-snappymail/docker/. "$LAB_DIR/"
chmod +x "$LAB_DIR/scripts/bootstrap.sh"

log "Construindo e iniciando stack Docker..."
cd "$LAB_DIR"
docker compose build --quiet
docker compose up -d

log "Aguardando serviços..."
sleep 15
bash scripts/bootstrap.sh || true

log "Stack Docker pronta em ${LAB_DIR}"
docker compose ps
