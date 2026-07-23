#!/usr/bin/env bash
set -euo pipefail
source /vagrant/provision/common.sh

SNAPPY_ADMIN_PW="${SNAPPY_ADMIN_PASS}"

cat <<EOF

================================================================================
  go-snappymail LAB — ambiente pronto
================================================================================

  VM IP ..........: ${VM_IP}
  Hostname .......: ${MAIL_FQDN}

  Adicione no /etc/hosts do seu host:
    ${VM_IP}  ${MAIL_FQDN}

--------------------------------------------------------------------------------
  Serviços Web
--------------------------------------------------------------------------------
  Go-PostfixAdmin : http://${VM_IP}:8081   (admin mail server)
  SnappyMail PHP  : http://${VM_IP}:8888   (webmail referência)
  Go webmail (Go) : http://${VM_IP}:8080   (go-cubemail — referência Golang)

--------------------------------------------------------------------------------
  Credenciais de teste
--------------------------------------------------------------------------------
  Mailbox ........: ${MAIL_USER} / ${MAIL_PASS}
  PostfixAdmin ...: ${ADMIN_EMAIL} / ${ADMIN_PASS}

  SnappyMail admin: admin / ${SNAPPY_ADMIN_PW}
  (painel admin SnappyMail: http://${VM_IP}:8888/?admin)

--------------------------------------------------------------------------------
  IMAP / SMTP (para clientes de email)
--------------------------------------------------------------------------------
  IMAP ...........: ${MAIL_FQDN}:993 (SSL/TLS)
  IMAP (plain) ...: ${MAIL_FQDN}:143
  SMTP ...........: ${MAIL_FQDN}:25
  SMTP submission : ${MAIL_FQDN}:587 (STARTTLS — se habilitado)

--------------------------------------------------------------------------------
  Validação rápida
--------------------------------------------------------------------------------
  # Testar IMAP login
  doveadm auth test '${MAIL_USER}' '${MAIL_PASS}'

  # Enviar email de teste
  echo "teste" | mail -s "Lab test" ${MAIL_USER}

  # Status dos serviços
  systemctl status postfix dovecot postfixadmin go-cubemail nginx --no-pager

  SSH ............: vagrant ssh  (ou root/vagrant123)

================================================================================
  Nota: go-snappymail ainda não tem código — go-cubemail serve como referência
  Golang na porta 8080. Quando go-snappymail P0 estiver pronto, substitua o
  serviço go-cubemail pelo binário go-snappymail.
================================================================================

EOF
