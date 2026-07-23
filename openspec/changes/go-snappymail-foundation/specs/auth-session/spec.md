## ADDED Requirements

### Requirement: IMAP-based login

The system SHALL authenticate users by validating IMAP LOGIN credentials against the configured mail server.

#### Scenario: Successful login

- **WHEN** the user submits valid email and password to `POST /api/v1/auth/login`
- **THEN** the system creates a server-side session and returns a `gsn_session` HTTP-only cookie

#### Scenario: Failed login

- **WHEN** the user submits invalid credentials
- **THEN** the system returns HTTP 401 and does NOT create a session

### Requirement: Login rate limiting

The system SHALL rate-limit login attempts to prevent brute-force attacks.

#### Scenario: Rate limit exceeded

- **WHEN** more than 10 failed login attempts occur from the same IP within 5 minutes
- **THEN** the system returns HTTP 429 until the window expires

### Requirement: Session management

The system SHALL store sessions in the database with configurable max age and periodic cleanup of expired sessions.

#### Scenario: Authenticated request

- **WHEN** a request includes a valid `gsn_session` cookie
- **THEN** the auth middleware attaches the user context and allows access to protected endpoints

#### Scenario: Session expiry

- **WHEN** a session exceeds its max age
- **THEN** the system rejects the request with HTTP 401

### Requirement: Logout

The system SHALL invalidate the session on logout.

#### Scenario: User logout

- **WHEN** the user calls `POST /api/v1/auth/logout`
- **THEN** the session is deleted from the database and the cookie is cleared

### Requirement: Current user endpoint

The system SHALL expose `GET /api/v1/auth/me` returning the authenticated user's email and display name.

#### Scenario: Get current user

- **WHEN** an authenticated user calls `/auth/me`
- **THEN** the system returns the user's email address and session metadata
