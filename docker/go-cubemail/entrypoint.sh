#!/bin/sh
set -e
# Trust the lab mailserver self-signed cert (mounted from mail_ssl volume)
if [ -f /etc/ssl/mail/mail.crt ] && [ -d /usr/local/share/ca-certificates ]; then
  cp /etc/ssl/mail/mail.crt /usr/local/share/ca-certificates/lab-mailserver.crt
  update-ca-certificates >/dev/null 2>&1 || cat /etc/ssl/mail/mail.crt >> /etc/ssl/certs/ca-certificates.crt
fi
./go-cubemail migrate
exec ./go-cubemail serve
