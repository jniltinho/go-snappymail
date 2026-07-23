## ADDED Requirements

### Requirement: User preferences

The system SHALL persist user preferences in the database and expose them via `GET/PUT /api/v1/settings`.

#### Scenario: Load settings

- **WHEN** an authenticated user requests settings
- **THEN** the system returns layout preferences, reading pane mode, messages per page, and locale

#### Scenario: Save settings

- **WHEN** the user updates settings via PUT
- **THEN** preferences are persisted and applied on next page load

### Requirement: Multiple identities

The system SHALL support multiple sending identities with per-identity signatures.

#### Scenario: List identities

- **WHEN** the user calls `GET /api/v1/identities`
- **THEN** the system returns all identities with email, display name, signature, and default flag

#### Scenario: Create identity

- **WHEN** the user sends `POST /api/v1/identities` with email and display name
- **THEN** a new identity is created

#### Scenario: Set default identity

- **WHEN** the user sends `POST /api/v1/identities/:id/default`
- **THEN** that identity becomes the default for new compose messages

### Requirement: Theme preference

The system SHALL persist the user's theme choice (light/dark/system) in settings.

#### Scenario: Dark mode preference

- **WHEN** the user selects dark mode in settings
- **THEN** the preference is saved and the UI renders in dark mode on reload
