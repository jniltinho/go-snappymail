## 1. Project Scaffold (P0)

- [x] 1.1 Initialize Go module `go-snappymail` with go.mod (Go 1.26+, Echo v5, Cobra, Viper, GORM, go-imap/v2)
- [x] 1.2 Create `main.go` with `go:embed web/dist` and delegate to `cmd` package
- [x] 1.3 Implement Cobra root command with subcommands: `init`, `migrate`, `serve`, `version`
- [x] 1.4 Implement `internal/config` with TOML loading and `GOSM_` env overrides
- [x] 1.5 Create `web/files/config.default.toml` (port 8082 default)
- [x] 1.6 Implement `cmd/init.go` to generate default config file
- [x] 1.7 Set up GORM with SQLite/MySQL/PostgreSQL drivers and Session model migrate
- [x] 1.8 Implement `cmd/migrate.go` to run database migrations
- [x] 1.9 Scaffold `internal/server/server.go` with Echo v5, middleware stack (recover, request-id, security headers, CSRF, slog logger)
- [x] 1.10 Implement static file serving for embedded SPA with SPA fallback routing
- [x] 1.11 Create Makefile with targets: `build`, `build-prod`, `release`, `run`, `clean` (binaries in `dist/`, UPX)
- [ ] 1.12 Add `.gitignore`, `LICENSE` (MIT), and initial `README.md`
- [x] 1.13 Docker service `go-snappymail` on port 8082 in `docker/docker-compose.yml`

## 2. Frontend Scaffold (P0)

- [x] 2.5 Build placeholder login page (`web/dist/index.html`) with SnappyMail-like styling
- [ ] 2.1 Initialize Vue 3 + TypeScript + Vite project in `frontend/`
- [ ] 2.2 Configure Tailwind CSS v4 with dark mode support
- [ ] 2.3 Set up Pinia stores skeleton: `auth`, `mail`, `settings`
- [ ] 2.4 Create API client utility with CSRF token handling and cookie credentials
- [ ] 2.5 Build placeholder `LoginView.vue` and `App.vue` shell
- [ ] 2.6 Configure Vite build output to `web/dist/`
- [ ] 2.7 Verify `make all` produces working binary serving the SPA

## 3. Auth & Session (P0)

- [x] 3.1 Implement `internal/session` package with in-memory session store + GORM Session model
- [x] 3.2 Implement `internal/imap/client.go` wrapper for LOGIN validation
- [x] 3.3 Implement `POST /api/v1/auth/login` with IMAP credential validation
- [x] 3.4 Implement auth middleware checking `gsn_session` cookie
- [x] 3.5 Implement `POST /api/v1/auth/logout` and `GET /api/v1/auth/me`
- [x] 3.6 Add login rate limiting middleware (10 attempts / min per IP)
- [x] 3.7 Implement session cleanup background goroutine
- [x] 3.8 Wire login page to auth API with CSRF token handling

## 4. IMAP Mailbox — Core Mail (P1)

- [ ] 4.1 Implement `GET /api/v1/folders` returning nested folder tree with unread counts
- [ ] 4.2 Implement folder CRUD: create, rename, delete endpoints
- [ ] 4.3 Implement `GET /api/v1/mail/:mailbox` with pagination, sort, and search params
- [ ] 4.4 Implement `GET /api/v1/mail/:mailbox/:uid` with MIME parsing (HTML, plain, attachments)
- [ ] 4.5 Implement `POST /api/v1/mail/:mailbox/:uid/flag` for seen/flagged/answered
- [ ] 4.6 Implement `POST /api/v1/mail/:mailbox/:uid/move` and `DELETE` for single message
- [ ] 4.7 Implement `DELETE /api/v1/mail/:mailbox` for empty trash
- [ ] 4.8 Implement attachment download endpoint with streaming
- [ ] 4.9 Implement `GET /api/v1/folders/:name/count` for unread badge updates

## 5. Compose & Send (P1)

- [ ] 5.1 Implement `internal/smtp/sender.go` using wneessen/go-mail
- [ ] 5.2 Implement `POST /api/v1/compose/send` with HTML body, attachments, threading headers
- [ ] 5.3 Implement `POST /api/v1/compose/draft` saving to IMAP Drafts folder
- [ ] 5.4 Implement `POST /api/v1/compose/upload` for temporary attachment storage
- [ ] 5.5 Build `ComposerModal.vue` with TipTap rich-text editor
- [ ] 5.6 Add reply, reply-all, and forward actions in mail store

## 6. Search (P1)

- [ ] 6.1 Implement `GET /api/v1/search` using IMAP SEARCH with field filters (subject, from, body)
- [ ] 6.2 Add search bar in `AppToolbar.vue` with results display in message list

## 7. UI — SnappyMail Layout (P2)

- [ ] 7.1 Build `FolderSidebar.vue` with folder tree, unread badges, and system folder icons
- [ ] 7.2 Build `MessageList.vue` with sender, subject, date, flag indicators, and selection
- [ ] 7.3 Build `ReadingPane.vue` with HTML rendering, attachment list, and action toolbar
- [ ] 7.4 Build `AppToolbar.vue` with search, refresh, compose, and settings buttons
- [ ] 7.5 Implement responsive layout: 3-column desktop, single-column mobile with navigation
- [ ] 7.6 Apply SnappyMail-inspired color palette (Default + NightShine themes) via Tailwind
- [ ] 7.7 Implement dark/light mode toggle with CSS variables and Tailwind `dark:` variants
- [ ] 7.8 Implement keyboard shortcuts (j/k/r/c/e/#) in App.vue
- [ ] 7.9 Implement SSE endpoint `/api/v1/events` with server-side new-mail detection (IMAP IDLE or polling) and frontend listener
- [ ] 7.10 Implement HTML sanitization with bluemonday and remote image blocking in ReadingPane
- [ ] 7.11 Add toast notification component for new mail alerts

## 8. Settings & Identities (P3)

- [ ] 8.1 Implement `GET/PUT /api/v1/settings` for user preferences
- [ ] 8.2 Implement identities CRUD endpoints (`/api/v1/identities`)
- [ ] 8.3 Build `SettingsPane.vue` with layout, reading pane, messages-per-page, and theme options
- [ ] 8.4 Wire identity selection in composer

## 9. Contacts (P3)

- [ ] 9.1 Implement contact model and GORM repository
- [ ] 9.2 Implement contacts CRUD REST endpoints
- [ ] 9.3 Implement CSV and vCard import/export endpoints
- [ ] 9.4 Build `ContactsPane.vue` with list, create/edit modal, and import/export buttons
- [ ] 9.5 Add contact autocomplete in composer To/Cc/Bcc fields

## 10. PGP Crypto (P4)

- [ ] 10.1 Implement PGP key model and encrypted storage in database
- [ ] 10.2 Implement key import/export REST endpoints
- [ ] 10.3 Integrate OpenPGP.js in frontend for encrypt/sign on send
- [ ] 10.4 Integrate OpenPGP.js for decrypt/verify on read with passphrase prompt
- [ ] 10.5 Add PGP settings section in SettingsPane

## 11. Sieve Filters (P4)

- [ ] 11.1 Implement ManageSieve client in `internal/sieve/`
- [ ] 11.2 Implement REST endpoints for list/create/update/delete/activate scripts
- [ ] 11.3 Build Sieve filter editor UI (visual + raw text modes)
- [ ] 11.4 Detect Sieve availability per domain and conditionally show filters section

## 12. Admin Panel (P4)

- [ ] 12.1 Implement admin password hash verification and admin session
- [ ] 12.2 Implement domain CRUD endpoints with IMAP/SMTP/Sieve presets
- [ ] 12.3 Implement domain whitelist enforcement at login
- [ ] 12.4 Build admin panel UI at `/admin` with domain management forms

## 13. Documentation & Quality

- [ ] 13.1 Add Swagger annotations to all API handlers and generate docs
- [ ] 13.2 Write `DOCUMENTS/docs/DEVELOPMENT.md` with local dev setup
- [ ] 13.3 Write `DOCUMENTS/setup/README.md` with production deployment guide (systemd + MariaDB)
- [x] 13.4 Add table-driven tests for auth, session crypto, middleware, routes, and database
- [ ] 13.5 Add GitHub Actions CI workflow (go test, golangci-lint, frontend build)
