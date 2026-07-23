# Architecture

## Overview

go-snappymail is a **single Go binary** that serves:

1. A versioned **REST API** (`/api/v1/*`)
2. An **embedded SPA** (`go:embed web/dist`)
3. **Session state** in SQLite/MySQL/PostgreSQL (GORM)

Mail bodies and folders live on the **IMAP server** — nothing is migrated into the app database.

```
Browser (SPA)  ──HTTPS──►  go-snappymail (Echo + Cobra + GORM)
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
         SQLite/MySQL    IMAP (Dovecot)   SMTP (Postfix)
         sessions only    mail storage     send mail
```

## Project layout

```
go-snappymail/
├── main.go                 # go:embed web/dist + web/files
├── cmd/                    # Cobra: init, migrate, serve, version
├── internal/
│   ├── config/             # Viper TOML + GOSM_* env
│   ├── handler/            # REST handlers
│   ├── imap/               # IMAP client wrapper
│   ├── session/            # In-memory + GORM session store
│   ├── model/              # GORM models
│   ├── database/           # DB open + slog logger
│   └── server/             # Echo, routes, middleware
├── web/dist/               # Embedded SPA (placeholder → Vue)
├── docker/                 # Lab stack
├── vagrant/                # Validation VM
├── docs/                   # Documentation
└── dist/                   # Binary output (gitignored)
```

## Middleware stack

Recover → Request ID → Security headers → CSRF → slog request logger.

Login uses **rate limiting** (10 requests/minute per IP).

## Session model

- Cookie `gsn_session`
- IMAP password: **AES-GCM** with `server.secret_key`
- GORM table `sessions` + in-memory cache

## Delivery phases

| Phase | Backend | Frontend |
|-------|---------|----------|
| P0 ✅ | Auth API, Echo, GORM, tests | `index.html` placeholder |
| P1 | Folders, messages, compose | Deferred |
| P2 | SSE, sanitization | Vue 3 SnappyMail layout |
| P3+ | Contacts, PGP, Sieve | Matching UI |

**Backend first, frontend second** for every feature slice.

## References

- [go-cubemail](https://github.com/jniltinho/go-cubemail) — architecture
- SnappyMail (`base/snappymail/` local) — UI parity
- OpenSpec `go-snappymail-foundation` — full specs
