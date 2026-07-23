## ADDED Requirements

### Requirement: Mail search

The system SHALL search messages across folders using IMAP SEARCH via `GET /api/v1/search`.

#### Scenario: Search by subject

- **WHEN** the user searches with `q=invoice&field=subject`
- **THEN** the system returns matching messages from all searchable folders with uid, folder, subject, from, and date

#### Scenario: Search by sender

- **WHEN** the user searches with `q=john@example.com&field=from`
- **THEN** the system returns messages from that sender

#### Scenario: Empty search results

- **WHEN** no messages match the query
- **THEN** the system returns an empty array with HTTP 200

### Requirement: Search scope

The system SHALL allow restricting search to the current folder or all folders.

#### Scenario: Folder-scoped search

- **WHEN** the user searches with `folder=INBOX`
- **THEN** results are limited to the INBOX folder
