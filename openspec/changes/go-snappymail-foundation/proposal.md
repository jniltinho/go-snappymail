## Why

SnappyMail is a mature, lightweight webmail client, but it depends on PHP and a traditional LAMP-style deployment. A Go reimplementation delivers the same SnappyMail UX as a single self-hosted binary — following the proven architecture of [go-cubemail](https://github.com/jniltinho/go-cubemail) — with better performance, simpler operations, and no PHP runtime. The reference SnappyMail source is available locally at `base/snappymail/` for feature and UI parity analysis.

## What Changes

- Bootstrap a new Go project `go-snappymail` with Cobra CLI (`init`, `migrate`, `serve`, `version`), Viper config (`GOSM_` env prefix), and embedded SPA — **default port 8082** for side-by-side lab comparison with go-cubemail (8080) and SnappyMail PHP (8888).
- Implement Echo v5 HTTP server with versioned REST API (`/api/v1`), CSRF protection, session cookies, and structured logging.
- Connect directly to user IMAP/SMTP servers (no email migration); messages stay on the mail server.
- Build a Vue 3 + Tailwind CSS v4 frontend replicating SnappyMail's 3-column layout: folder sidebar, message list, reading pane — plus composer modal, settings, contacts, and dark mode.
- Port core SnappyMail mail features: folder management, message list/read/flag/move/delete, compose/send/draft, attachments, search, safe HTML rendering with remote image blocking.
- Add SnappyMail-differentiating features in later phases: PGP (OpenPGP.js on frontend, key handling on backend), Sieve filter editor, admin panel for domain configuration.
- Store only sessions, settings, contacts, and PGP key metadata in SQLite/MySQL/PostgreSQL — never duplicate mail bodies.

## Capabilities

### New Capabilities

- `project-scaffold`: Go module layout, CLI commands, config loading, database migrations, embedded SPA serving, Makefile build pipeline.
- `auth-session`: IMAP credential validation at login, secure session cookies, logout, `/auth/me`, rate limiting on login.
- `imap-mailbox`: Folder tree CRUD, message listing with pagination/sort/search, read message with MIME parsing, flag/move/delete/empty-trash, attachment download.
- `compose-send`: Rich-text compose (TipTap), send via SMTP, save draft to IMAP, attachment upload, reply/forward/reply-all.
- `search`: Full-text search across folders via IMAP SEARCH.
- `contacts`: Address book CRUD, import/export (CSV/vCard), autocomplete in composer.
- `settings`: User preferences (layout, reading pane, dark mode, signatures), multiple identities.
- `ui-snappymail`: Vue 3 SPA with SnappyMail-like visual design — 3-column responsive layout, keyboard shortcuts (j/k/r/c/e), SSE new-mail notifications, theme switching.
- `pgp-crypto`: PGP key management, encrypt/sign on send, decrypt/verify on read (phase 2).
- `sieve-filters`: ManageSieve protocol support and visual filter editor (phase 2).
- `admin-panel`: Admin authentication, domain/IMAP-SMTP preset configuration, white-list (phase 2).

### Modified Capabilities

_(none — greenfield project, no existing specs in `openspec/specs/`)_

## Non-goals

- **ActiveSync (EAS)**, **CalDAV**, and **CardDAV servers** — deferred; go-cubemail already covers groupware sync. go-snappymail focuses on webmail parity with SnappyMail.
- **POP3 support** — intentionally excluded (SnappyMail removed it).
- **Multi-account per user** — deferred; v1 supports one IMAP account per session.
- **PHP plugin compatibility** — no `.phar` plugin runtime; extensibility via Go plugins or hooks is future work.
- **1:1 API compatibility** with SnappyMail's JSON `?/Ajax/` endpoints — we use REST JSON API like go-cubemail, not RainLoop action names.

## Impact

- **New repository structure**: `cmd/`, `internal/` (config, handler, imap, smtp, session, repository, server), `frontend/` (Vue 3), `web/dist/` (build output), `docs/` (Swagger).
- **Dependencies**: Echo v5, GORM, emersion/go-imap/v2, wneessen/go-mail, microcosm-cc/bluemonday (HTML sanitization), Cobra, Viper, Swaggo — aligned with go-cubemail.
- **Frontend dependencies**: Vue 3, Vite, Pinia, Tailwind CSS v4, TipTap, Lucide icons, OpenPGP.js (phase 2).
- **Reference code**: `base/snappymail/` (read-only reference for UI/UX and feature mapping); go-cubemail (architecture and handler patterns to reuse/adapt).
- **License**: MIT (distinct from SnappyMail's AGPL — clean-room reimplementation, no PHP code copy).
