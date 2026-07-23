#!/usr/bin/env bash
set -euo pipefail
source /vagrant/provision/common.sh

log "Configurando MariaDB"

systemctl enable mariadb
systemctl start mariadb

mariadb <<SQL
CREATE DATABASE IF NOT EXISTS ${DB_NAME};
CREATE USER IF NOT EXISTS '${DB_USER}'@'localhost' IDENTIFIED BY '${DB_PASS}';
GRANT ALL ON ${DB_NAME}.* TO '${DB_USER}'@'localhost';
FLUSH PRIVILEGES;
SQL

log "MariaDB pronta (db=${DB_NAME}, user=${DB_USER})"
