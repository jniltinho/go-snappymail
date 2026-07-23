# Documentation index

Guides for **go-snappymail** — SnappyMail UX webmail in Go.

## Getting started

| Document | Audience |
|----------|----------|
| [../README.md](../README.md) | Everyone — quick start |
| [lab.md](lab.md) | Local mail lab (Docker / Vagrant) |
| [development.md](development.md) | Contributors — build, test, workflow |
| [configuration.md](configuration.md) | Operators — `config.toml`, env vars |
| [skins.md](skins.md) | UI skins — create SnappyMail / Gmail / Outlook / custom |

## Design & API

| Document | Content |
|----------|---------|
| [architecture.md](architecture.md) | System design, folders, phases |
| [api.md](api.md) | REST API reference (current + planned) |
| [../openspec/changes/go-snappymail-foundation/](../openspec/changes/go-snappymail-foundation/) | Full OpenSpec proposal |

## Operations & security

| Document | Content |
|----------|---------|
| [security.md](security.md) | Secrets, git rules, production checklist |
| [../docker/README.md](../docker/README.md) | Docker Compose lab |
| [../vagrant/README.md](../vagrant/README.md) | Vagrant VM |
| [../docker/LAB_ACCOUNTS.md](../docker/LAB_ACCOUNTS.md) | Test mailboxes |

## Local-only (not in git)

| Path | Purpose |
|------|---------|
| `base/snappymail/` | SnappyMail PHP reference (~100 MB) — clone locally for UI parity |
| `dist/go-snappymail` | Compiled binary (UPX) |
| `docs/prints/` | Local screenshots from lab validation (agent-browser) |
| `docker/.env`, `vagrant/.env` | Lab secrets |
