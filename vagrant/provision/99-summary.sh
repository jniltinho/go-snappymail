#!/usr/bin/env bash
set -euo pipefail
source /vagrant/provision/common.sh

SNAPPY_ADMIN_PW="$(docker exec gosm-snappymail cat /var/lib/snappymail/_data_/_default_/admin_password.txt 2>/dev/null || echo 'ver: docker exec gosm-snappymail cat /var/lib/snappymail/_data_/_default_/admin_password.txt')"

cat <<EOF

================================================================================
  go-snappymail LAB (híbrido) — ambiente pronto
================================================================================

  VM IP ..........: ${VM_IP}
  Hostname .......: ${MAIL_FQDN}

  Adicione no /etc/hosts do seu host:
    ${VM_IP}  ${MAIL_FQDN}

--------------------------------------------------------------------------------
  Nativo na VM (systemd)                 | Docker (${COMPOSE_DIR})
--------------------------------------------------------------------------------
  go-snappymail ..: http://${VM_IP}:8082 | MariaDB ......: 127.0.0.1:3306
  go-cubemail ....: http://${VM_IP}:8080 | SnappyMail PHP: http://${VM_IP}:8888
  go-postfixadmin : http://${VM_IP}:8081 |
  Postfix ........: :25                  |
  Dovecot ........: :143 / :993 (SSL)    |

--------------------------------------------------------------------------------
  Credenciais de teste
--------------------------------------------------------------------------------
  Mailbox ........: ${MAIL_USER} / ${MAIL_PASS}
  PostfixAdmin ...: ${ADMIN_EMAIL} / ${ADMIN_PASS}
  SnappyMail admin: admin / ${SNAPPY_ADMIN_PW}
  (painel admin SnappyMail: http://${VM_IP}:8888/?admin)

--------------------------------------------------------------------------------
  Validação rápida
--------------------------------------------------------------------------------
  # Apps Go (systemd)
  systemctl status go-snappymail go-cubemail postfixadmin --no-pager
  journalctl -u go-snappymail -f

  # Containers
  docker compose -f ${COMPOSE_DIR}/docker-compose.yml ps

  # IMAP login
  doveadm auth test '${MAIL_USER}' '${MAIL_PASS}'

  # Enviar email de teste
  echo "teste" | mail -s "Lab test" ${MAIL_USER}

  # API go-snappymail
  curl http://${VM_IP}:8082/api/v1/version

  SSH ............: vagrant ssh  (ou root/vagrant123)

================================================================================

EOF
