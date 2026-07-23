FROM golang:1.26-bookworm AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /go-snappymail .

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /go-snappymail /app/go-snappymail
COPY docker/go-snappymail/config.toml /app/config.toml
COPY docker/go-snappymail/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh \
    && mkdir -p /app/data /app/tmp/uploads

EXPOSE 8082

ENTRYPOINT ["/app/entrypoint.sh"]
