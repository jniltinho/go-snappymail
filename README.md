# go-snappymail

Self-hosted webmail in Go with SnappyMail-inspired UX. Single binary, IMAP/SMTP passthrough, embedded SPA.

## Stack

| Component | Description |
|-----------|-------------|
| **go-snappymail** | This project — Echo v5, Cobra, embedded UI |
| **Reference** | [go-cubemail](https://github.com/jniltinho/go-cubemail) architecture |
| **UI target** | SnappyMail 3-column layout (Vue 3 in P1+) |

## Quick start

```bash
make build-prod          # → dist/go-snappymail (UPX compressed)
./dist/go-snappymail init
./dist/go-snappymail migrate
./dist/go-snappymail serve
```

Default HTTP port: **8082** (see `web/files/config.default.toml`).

## Docker lab

```bash
cd docker
cp .env.example .env
docker compose up -d --build
bash scripts/bootstrap.sh
```

See [docker/LAB_ACCOUNTS.md](docker/LAB_ACCOUNTS.md) for test mailboxes (4 domains, 15 accounts).

| Service | Port |
|---------|------|
| go-snappymail | 8082 |
| go-cubemail | 8080 |
| PostfixAdmin | 8081 |
| SnappyMail (PHP) | 8888 |

## Vagrant

Bare-metal or Docker mode on `192.168.56.20` — see [vagrant/README.md](vagrant/README.md).

## Spec

OpenSpec change: `openspec/changes/go-snappymail-foundation/`

## License

MIT
