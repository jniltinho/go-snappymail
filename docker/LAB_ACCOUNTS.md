# Lab Test Accounts

All mailboxes use password **`Password1@`** unless noted. IMAP host: `mailserver` (Docker) or `mail.test.local` (host/VM). Port **993** (TLS).

## Domains

| Domain | Description |
|--------|-------------|
| `test.local` | Primary lab domain |
| `acme.local` | Acme Corp sandbox |
| `demo.local` | Demo tenant |
| `org.local` | Organization tenant |

## Mailboxes

### test.local

| Email | Role |
|-------|------|
| `user@test.local` | Primary test user |
| `alice@test.local` | Secondary user |
| `bob@test.local` | Secondary user |
| `admin@test.local` | Admin mailbox |

### acme.local

| Email | Role |
|-------|------|
| `ceo@acme.local` | Executive |
| `sales@acme.local` | Sales team |
| `support@acme.local` | Support team |
| `billing@acme.local` | Billing |

### demo.local

| Email | Role |
|-------|------|
| `demo@demo.local` | Demo user |
| `guest@demo.local` | Guest account |
| `trainer@demo.local` | Training account |

### org.local

| Email | Role |
|-------|------|
| `team@org.local` | General team |
| `hr@org.local` | HR |
| `dev@org.local` | Developers |
| `ops@org.local` | Operations |

## Webmail URLs (VM `192.168.56.20`)

| Service | URL |
|---------|-----|
| go-snappymail | http://192.168.56.20:8082 |
| go-cubemail | http://192.168.56.20:8080 |
| SnappyMail (PHP) | http://192.168.56.20:8888 |
| PostfixAdmin | http://192.168.56.20:8081 |

## Seed command

```bash
cd docker
docker compose up -d
bash scripts/bootstrap.sh   # includes seed-lab.sh
```

To re-run accounts only:

```bash
bash scripts/seed-lab.sh
```
