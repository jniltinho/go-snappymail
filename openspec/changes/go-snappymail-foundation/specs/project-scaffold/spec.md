## ADDED Requirements

### Requirement: Single binary with embedded SPA

The system SHALL compile to a single executable that embeds the Vue 3 frontend build output via `go:embed web/dist` and serves it without requiring Node.js at runtime.

#### Scenario: Serve embedded frontend

- **WHEN** the user runs `./go-snappymail serve`
- **THEN** the binary serves the SPA at `/` and static assets from the embedded filesystem

### Requirement: CLI commands

The system SHALL provide Cobra CLI commands: `init`, `migrate`, `serve`, and `version`.

#### Scenario: Initialize configuration

- **WHEN** the user runs `./go-snappymail init`
- **THEN** the system generates a default `config.toml` with server, database, and IMAP/SMTP placeholders

#### Scenario: Run database migrations

- **WHEN** the user runs `./go-snappymail migrate`
- **THEN** the system creates or updates all required database tables

### Requirement: Configuration via TOML and environment

The system SHALL load configuration from `config.toml` with environment variable overrides prefixed `GOSM_`.

#### Scenario: Override port via environment

- **WHEN** `GOSM_SERVER_PORT=9090` is set
- **THEN** the server listens on port 9090 regardless of the TOML value

### Requirement: Build pipeline

The project SHALL include a Makefile with `make all` that builds the frontend (Vite) and compiles the Go binary with embedded assets.

#### Scenario: Full build from source

- **WHEN** the developer runs `make all`
- **THEN** `web/dist/` is populated and `bin/go-snappymail` is produced
