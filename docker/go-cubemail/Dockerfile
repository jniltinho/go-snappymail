ARG GO_CUBEMAIL_VERSION=0.0.25
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates wget \
    && rm -rf /var/lib/apt/lists/*

ARG GO_CUBEMAIL_VERSION
WORKDIR /app

RUN wget -q -O /tmp/go-cubemail.tar.gz \
    "https://github.com/jniltinho/go-cubemail/releases/download/v${GO_CUBEMAIL_VERSION}/go-cubemail_${GO_CUBEMAIL_VERSION}_linux_amd64.tar.gz" \
    && tar -xzf /tmp/go-cubemail.tar.gz -C /app \
    && chmod +x /app/go-cubemail \
    && rm /tmp/go-cubemail.tar.gz \
    && mkdir -p /app/data /app/tmp/uploads

COPY config.toml /app/config.toml
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]
