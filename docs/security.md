# Security & git hygiene

## Never commit

| Category | Examples |
|----------|----------|
| Reference tree | `base/snappymail/` (~100 MB PHP) |
| Binaries | `dist/go-snappymail`, UPX output |
| Secrets | `docker/.env`, `vagrant/.env`, local `config.toml` |
| Keys | `id_rsa`, `*.pem`, private key blocks |
| Tokens | `ghp_*`, `gho_*`, `AKIA*` |
| Tooling | `.agents/`, `.cursor/`, `vagrant/.vagrant/` |
| Runtime | `*.db`, `data/`, `coverage.out` |

## Pre-push check

```bash
make check-git
```

Script: [scripts/check-git-clean.sh](../scripts/check-git-clean.sh)

## Where secrets belong

| Secret | Location |
|--------|----------|
| Mail/DB passwords | `docker/.env` or `vagrant/.env` |
| `server.secret_key` | Local `config.toml` or `GOSM_SERVER_SECRET_KEY` |
| PostfixAdmin JWT | `SESSION_SECRET` in `docker/.env` |

Copy templates:

```bash
cp docker/.env.example docker/.env
cp vagrant/.env.example vagrant/.env
```

## Application security (P0)

| Feature | Implementation |
|---------|----------------|
| Session passwords | AES-GCM + `secret_key` |
| CSRF | Double-submit cookie + `X-CSRF-Token` header |
| Login brute force | Rate limit 10/min per IP |
| Security headers | HSTS, XSS, frame deny |
| IMAP credentials | Never logged; validated at login only |

## Lab vs production

Lab docs mention `Password1@` as a **local default** in `.env.example` only.

Production:

- [ ] New random `secret_key`
- [ ] HTTPS + `session.secure = true`
- [ ] `imap.insecure_skip_verify = false`
- [ ] Strong DB credentials in env, not in repo
- [ ] Reverse proxy (Caddy/nginx) with TLS

## Session storage

GORM stores **encrypted** IMAP password — not plaintext. Mail bodies are **never** stored in the app DB.
