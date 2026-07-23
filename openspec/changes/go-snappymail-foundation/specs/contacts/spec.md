## ADDED Requirements

### Requirement: Contact CRUD

The system SHALL store contacts in the database and expose CRUD via REST API.

#### Scenario: List contacts

- **WHEN** the user calls `GET /api/v1/contacts`
- **THEN** the system returns all contacts with name, email, phone, and notes

#### Scenario: Create contact

- **WHEN** the user sends `POST /api/v1/contacts` with name and email
- **THEN** the contact is persisted and returned with an ID

#### Scenario: Update contact

- **WHEN** the user sends `PUT /api/v1/contacts/:id`
- **THEN** the contact fields are updated

#### Scenario: Delete contact

- **WHEN** the user sends `DELETE /api/v1/contacts/:id`
- **THEN** the contact is removed

### Requirement: Contact import and export

The system SHALL support importing and exporting contacts in CSV and vCard formats.

#### Scenario: Import CSV

- **WHEN** the user uploads a CSV file via `POST /api/v1/contacts/import`
- **THEN** the system creates contacts from each row and returns an import summary

#### Scenario: Export contacts

- **WHEN** the user calls `GET /api/v1/contacts/export?format=vcard`
- **THEN** the system downloads a vCard file with all contacts

### Requirement: Autocomplete in composer

The frontend SHALL suggest contacts from the address book when the user types in To/Cc/Bcc fields.

#### Scenario: Autocomplete suggestion

- **WHEN** the user types "joh" in the To field
- **THEN** matching contacts appear as selectable suggestions
