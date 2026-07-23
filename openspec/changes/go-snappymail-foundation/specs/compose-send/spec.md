## ADDED Requirements

### Requirement: Send email via SMTP

The system SHALL send email through the configured SMTP server via `POST /api/v1/compose/send`.

#### Scenario: Send new message

- **WHEN** the user submits to, cc, bcc, subject, HTML body, and optional attachments
- **THEN** the system sends the message via SMTP and returns success with a message ID

#### Scenario: Reply to message

- **WHEN** the user sends with `inReplyTo` and `references` headers set
- **THEN** the system sends a properly threaded reply

### Requirement: Save draft

The system SHALL save drafts to the IMAP Drafts folder via `POST /api/v1/compose/draft`.

#### Scenario: Save draft

- **WHEN** the user saves a compose form as draft
- **THEN** the system stores the message in the Drafts folder on the IMAP server

### Requirement: Attachment upload

The system SHALL accept file uploads for compose via `POST /api/v1/compose/upload`.

#### Scenario: Upload attachment

- **WHEN** the user uploads a file during compose
- **THEN** the system returns a temporary attachment ID usable in the send request

### Requirement: Rich text editor

The frontend SHALL provide a rich-text compose editor supporting bold, italic, links, lists, and inline images.

#### Scenario: Compose with formatting

- **WHEN** the user formats text in the composer
- **THEN** the sent email contains valid HTML with the applied formatting

### Requirement: Default sender identity

The system SHALL use the authenticated user's email as the default From address when composing. Full identity selection with signatures is deferred to the settings capability (P3).

#### Scenario: Send with default identity

- **WHEN** the user sends a message without selecting an identity
- **THEN** the From header matches the authenticated user's email address
