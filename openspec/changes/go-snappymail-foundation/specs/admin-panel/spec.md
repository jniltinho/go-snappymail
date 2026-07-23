## ADDED Requirements

### Requirement: Admin authentication

The system SHALL provide a separate admin login protected by a password hash stored in configuration (not IMAP credentials).

#### Scenario: Admin login

- **WHEN** the admin submits the correct password to `POST /api/v1/admin/login`
- **THEN** an admin session cookie is created with elevated privileges

#### Scenario: Admin login failure

- **WHEN** the admin submits an incorrect password
- **THEN** the system returns HTTP 401 and logs the attempt

### Requirement: Domain configuration

The system SHALL allow the admin to configure IMAP/SMTP/Sieve server presets per email domain.

#### Scenario: Add domain

- **WHEN** the admin creates a domain with IMAP host, port, and security settings
- **THEN** users with matching email domains auto-fill server settings at login

#### Scenario: Domain whitelist

- **WHEN** the admin sets a domain whitelist
- **THEN** only email addresses from listed domains can log in

### Requirement: Admin panel UI

The frontend SHALL provide an admin panel accessible at `/admin` with domain management and system settings.

#### Scenario: Admin panel access

- **WHEN** an authenticated admin visits `/admin`
- **THEN** the admin dashboard with domain list and configuration forms is displayed

#### Scenario: Non-admin blocked

- **WHEN** a regular user attempts to access `/admin`
- **THEN** the system redirects to the login page or returns HTTP 403
