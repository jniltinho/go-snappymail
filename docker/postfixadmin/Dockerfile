ARG GO_POSTFIXADMIN_VERSION=1.0.86
FROM alpine:3.21

ARG GO_POSTFIXADMIN_VERSION
ENV TZ=America/Sao_Paulo

RUN apk add --no-cache ca-certificates tzdata wget netcat-openbsd

WORKDIR /app

RUN wget -q -O /tmp/pfa.tar.gz \
    "https://github.com/jniltinho/go-postfixadmin/releases/download/v${GO_POSTFIXADMIN_VERSION}/postfixadmin_${GO_POSTFIXADMIN_VERSION}_linux_amd64.tar.gz" \
    && tar -xzf /tmp/pfa.tar.gz -C /app \
    && chmod +x /app/postfixadmin \
    && rm /tmp/pfa.tar.gz

COPY config/postfixadmin.toml /app/config.toml.in
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

EXPOSE 8080 12221

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["./postfixadmin", "server"]
