FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates openssl wget \
    postfix postfix-mysql \
    dovecot-core dovecot-imapd dovecot-pop3d dovecot-lmtpd dovecot-mysql \
    mariadb-client supervisor \
    && rm -rf /var/lib/apt/lists/*

RUN groupadd -g 1001 vmail || true \
    && useradd -g vmail -u 1001 vmail -d /var/vmail -m || true \
    && mkdir -p /var/vmail /etc/postfix/sql /var/spool/postfix/private \
    && chown -R vmail:vmail /var/vmail

COPY templates/ /templates/
COPY entrypoint.sh /entrypoint.sh
COPY supervisord.conf /etc/supervisor/conf.d/mail.conf
RUN chmod +x /entrypoint.sh

EXPOSE 25 143 993

VOLUME ["/var/vmail", "/etc/ssl/mail"]

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf"]
