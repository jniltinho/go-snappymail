#!/usr/bin/env bash
set -euo pipefail
source /vagrant/provision/common.sh

log "Atualizando sistema e instalando pacotes base"
export DEBIAN_FRONTEND=noninteractive

debconf-set-selections <<< "postfix postfix/mailname string ${MAIL_FQDN}"
debconf-set-selections <<< "postfix postfix/main_mailer_type string 'Internet Site'"

apt-get update -qq
apt-get upgrade -y -qq
apt-get install -y -qq \
  curl wget ca-certificates gnupg openssl \
  postfix postfix-mysql \
  dovecot-core dovecot-imapd dovecot-pop3d dovecot-lmtpd dovecot-mysql \
  mariadb-server \
  nginx \
  php-fpm php-cli php-mbstring php-json php-xml php-zip php-intl php-curl php-gd php-sqlite3 php-mysql \
  sqlite3 \
  mailutils swaks

# Hostname
hostnamectl set-hostname "${MAIL_FQDN}"
grep -q "${MAIL_FQDN}" /etc/hosts || echo "${VM_IP} ${MAIL_FQDN} ${MAIL_DOMAIN}" >> /etc/hosts

ensure_self_signed_cert

# vmail user
if ! id vmail &>/dev/null; then
  groupadd -g "${VMAIL_GID}" vmail
  useradd -g vmail -u "${VMAIL_UID}" vmail -d /var/vmail -m
fi
chown -R vmail:vmail /var/vmail

# Root SSH for convenience in lab
echo "root:vagrant123" | chpasswd
sed -i 's/^#\?PermitRootLogin.*/PermitRootLogin yes/' /etc/ssh/sshd_config
sed -i 's/^#\?PasswordAuthentication.*/PasswordAuthentication yes/' /etc/ssh/sshd_config
systemctl restart ssh 2>/dev/null || systemctl restart sshd 2>/dev/null || true

log "Base concluída"
