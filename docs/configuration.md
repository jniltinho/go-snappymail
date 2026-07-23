# Configuration

Load order: **`config.toml`** → **`GOSM_*` environment variables** (Viper).

```bash
./dist/go-snappymail init
# or
cp web/files/config.default.toml config.toml
```

## Env mapping

`server.secret_key` → `GOSM_SERVER_SECRET_KEY`

```bash
GOSM_SERVER_PORT=9000 ./dist/go-snappymail serve
```

## Key settings

### Server (default port **8082**)

| Key | Notes |
|-----|-------|
| `secret_key` | 32+ chars — encrypts IMAP passwords in session |
| `base_url` | Public URL for cookies/links |
| `tls_cert` / `tls_key` | Optional HTTPS |

### IMAP

| Key | Notes |
|-----|-------|
| `host` | IMAP server hostname |
| `tls_server_name` | SNI when host is Docker service name |
| `insecure_skip_verify` | Lab only — never in production |

### Database

| driver | dsn example |
|--------|-------------|
| `sqlite` | `./data/app.db` (default) |
| `mariadb` | `user:pass@tcp(host:3306)/db?charset=utf8mb4&parseTime=True` |

### Session cookie

`gsn_session` — HttpOnly, configurable `max_age`, set `secure=true` with HTTPS.

## Docker secrets

```bash
cp docker/.env.example docker/.env
```

| Variable | Purpose |
|----------|---------|
| `SERVER_SECRET_KEY` | App encryption key |
| `SESSION_SECRET` | PostfixAdmin JWT |
| `MAIL_PASS` | Lab mailbox password |

Full template: [web/files/config.default.toml](../web/files/config.default.toml)
