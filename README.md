# go-snappymail

Self-hosted webmail in Go with **SnappyMail-inspired UX**. Single binary, IMAP/SMTP passthrough, embedded SPA — no PHP runtime.

| | |
|---|---|
| **Module** | `go-snappymail` |
| **Default port** | `8082` |
| **Stack** | Go 1.26 · Echo v5 · GORM · Cobra · Viper |
| **Reference** | [go-cubemail](https://github.com/jniltinho/go-cubemail) (architecture) · SnappyMail (UI/UX) |

## Quick start

```bash
git clone git@github.com:jniltinho/go-snappymail.git
cd go-snappymail

make build-prod          # → dist/go-snappymail (UPX)
./dist/go-snappymail init
./dist/go-snappymail migrate
./dist/go-snappymail serve
```

Open http://localhost:8082

## Documentation

| Guide | Description |
|-------|-------------|
| [Architecture](docs/architecture.md) | Components, project layout, delivery phases |
| [Development](docs/development.md) | Build, test, backend-first workflow |
| [Configuration](docs/configuration.md) | `config.toml` and `GOSM_*` environment variables |
| [Skins / themes](docs/skins.md) | **Guia de implementação** — criar skins custom (CSS, registro, validação) |
| [API](docs/api.md) | REST endpoints |
| [Security](docs/security.md) | Secrets, `.env`, git hygiene |
| [Lab environment](docs/lab.md) | Docker + Vagrant comparison stack |
| [Docker lab](docker/README.md) | Container setup and commands |
| [Vagrant VM](vagrant/README.md) | Ubuntu 24.04 validation VM |
| [Test accounts](docker/LAB_ACCOUNTS.md) | 4 domains, 15 mailboxes |
| [OpenSpec](openspec/changes/go-snappymail-foundation/) | Proposal, design, tasks |

## Lab URLs (comparison)

Add to `/etc/hosts`: `192.168.56.20 mail.test.local`

| Service | Port | Purpose |
|---------|------|---------|
| **go-snappymail** | 8082 | This project |
| go-cubemail | 8080 | Go reference webmail |
| PostfixAdmin | 8081 | Mailbox admin |
| SnappyMail (PHP) | 8888 | UX/behavior reference |

```bash
cd docker && cp .env.example .env
docker compose up -d --build && bash scripts/bootstrap.sh
```

## Make targets

```bash
make help          # all targets
make test          # unit tests + race + coverage
make check-git     # block secrets, base/, binaries before push
```

## Principles

1. **Backend first, frontend second** — complete Go API + tests before Vue UI (P1+).
2. **No mail duplication** — messages stay on the IMAP server; GORM stores sessions/settings only.
3. **Clean git** — never commit `base/`, `dist/`, `.env`, or private keys. Run `make check-git`.

## License

MIT — see [LICENSE](LICENSE)
