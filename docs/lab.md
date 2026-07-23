# Lab environment

Local mail stack to compare **go-snappymail** against SnappyMail (PHP) and go-cubemail.

## Two ways to run

| Mode | Where | Best for |
|------|-------|----------|
| **Docker on host** | Your machine | Fast iteration, `make build` + compose |
| **Vagrant VM** | `192.168.56.20` | Side-by-side validation, agent-browser screenshots |

Both use the same ports and test accounts.

## Hosts file

Add one line to `/etc/hosts`:

```
192.168.56.20  mail.test.local
```

For Docker on localhost only:

```
127.0.0.1  mail.test.local
```

## Service map

| Service | Port | URL (VM) |
|---------|------|----------|
| **go-snappymail** | 8082 | http://192.168.56.20:8082 |
| go-cubemail | 8080 | http://192.168.56.20:8080 |
| PostfixAdmin | 8081 | http://192.168.56.20:8081 |
| SnappyMail (PHP) | 8888 | http://192.168.56.20:8888 |
| SMTP | 25 | mail.test.local:25 |
| IMAP | 143 / 993 | mail.test.local |

## Quick start (Docker)

```bash
cd docker
cp .env.example .env          # edit secrets if needed
docker compose up -d --build
bash scripts/bootstrap.sh     # seed domains + mailboxes
```

From repo root (app outside Docker):

```bash
make build-prod
./dist/go-snappymail init && ./dist/go-snappymail migrate
# point config.toml imap.host to mailserver or 127.0.0.1:993
./dist/go-snappymail serve
```

When go-snappymail runs **inside** compose, IMAP host is `mailserver` (Docker network).

## Test accounts

**15 mailboxes** across **4 domains** — full list in [docker/LAB_ACCOUNTS.md](../docker/LAB_ACCOUNTS.md).

Default password: value of `MAIL_PASS` in `docker/.env` (`.env.example` uses `Password1@` for lab only).

Quick login:

| Email | Password |
|-------|----------|
| `user@test.local` | `Password1@` |
| `ceo@acme.local` | `Password1@` |
| PostfixAdmin admin | `admin@test.local` / `Password1@` |

## Vagrant VM

```bash
cd vagrant
cp .env.example .env
vagrant up
```

**Docker mode** (recommended for full stack including go-snappymail :8082):

```bash
vagrant provision --provision-with docker
```

Do **not** run bare-metal and Docker modes at the same time — same ports.

Details: [vagrant/README.md](../vagrant/README.md)

## Verify IMAP

```bash
cd docker
docker compose exec mailserver doveadm auth test 'user@test.local' 'Password1@'
```

## Side-by-side UX comparison

1. Log into SnappyMail (:8888) and go-snappymail (:8082) with the same mailbox
2. Compare folder list, message read, compose (when P1+ is ready)
3. Optional: save screenshots locally under `docs/prints/` (gitignored) with [agent-browser](https://github.com/nicholasgriffintn/agent-browser)

## Architecture

```
Browser
   ├── :8082  go-snappymail (Go)
   ├── :8080  go-cubemail (Go reference)
   ├── :8888  SnappyMail (PHP reference)
   └── :8081  PostfixAdmin
              │
              ▼
         mailserver (Postfix + Dovecot)
              │
              ▼
           MariaDB
```

## Reset lab

```bash
cd docker
docker compose down -v    # destroys volumes
docker compose up -d --build
bash scripts/bootstrap.sh
```

## Related docs

- [docker/README.md](../docker/README.md) — compose commands
- [configuration.md](configuration.md) — `config.toml` for lab IMAP
- [security.md](security.md) — never commit `.env`
