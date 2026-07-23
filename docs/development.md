# Development guide

## Prerequisites

- Go **1.26+**, **CGO** enabled
- Optional: UPX, Docker Compose v2

## Backend-first workflow

1. Handler in `internal/handler/`
2. Route in `internal/server/routes.go`
3. Table-driven tests
4. `make test`
5. Vue UI only after API is stable (P2+)

## Build & run

```bash
make build-prod          # builds Vue → web/dist, then Go binary
./dist/go-snappymail init
./dist/go-snappymail migrate
./dist/go-snappymail serve
```

Frontend dev (API proxy to :8082):

```bash
make run                 # terminal 1 — backend
make frontend-dev        # terminal 2 — Vite on :5173
```

## CLI

| Command | Description |
|---------|-------------|
| `init` | Generate config from template |
| `migrate` | AutoMigrate sessions table |
| `serve` | HTTP server |
| `version` | Build info |

## Testing

```bash
make test
make test-integration   # needs Docker lab + -tags=integration
make check-git          # before every push
```

## Conventions

- Module: `go-snappymail`
- Env prefix: `GOSM_`
- Errors: `%w` wrapping
- Logging: `log/slog`
- Docs/code: English

See [configuration.md](configuration.md) and [security.md](security.md).
