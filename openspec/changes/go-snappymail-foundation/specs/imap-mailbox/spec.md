## ADDED Requirements

### Requirement: Folder tree listing

The system SHALL return the full IMAP folder hierarchy with unread counts via `GET /api/v1/folders`.

#### Scenario: List folders

- **WHEN** an authenticated user requests the folder list
- **THEN** the system returns a nested JSON tree with folder names, types (inbox/sent/drafts/trash/spam/custom), and unread message counts

### Requirement: Folder management

The system SHALL support creating, renaming, and deleting IMAP folders.

#### Scenario: Create subfolder

- **WHEN** the user sends `POST /api/v1/folders` with a parent folder and name
- **THEN** the system creates the folder on the IMAP server

#### Scenario: Delete folder

- **WHEN** the user sends `POST /api/v1/folders/delete` for a non-system folder
- **THEN** the system removes the folder from the IMAP server

### Requirement: Message listing

The system SHALL list messages in a folder with pagination, sorting, and optional search filter via `GET /api/v1/mail/:mailbox`.

#### Scenario: Paginated message list

- **WHEN** the user requests a folder with `limit=50&offset=0&sort=date_desc`
- **THEN** the system returns up to 50 messages with uid, subject, from, date, flags (seen/flagged/answered), and has-attachments indicator

### Requirement: Read message

The system SHALL fetch and parse a full message via `GET /api/v1/mail/:mailbox/:uid`.

#### Scenario: Read message with attachments

- **WHEN** the user opens a message
- **THEN** the system returns headers, HTML and plain-text body, inline images, and a list of attachments with download URLs

### Requirement: Message flag operations

The system SHALL support toggling IMAP flags (seen, flagged, answered) via `POST /api/v1/mail/:mailbox/:uid/flag`.

#### Scenario: Mark as read

- **WHEN** the user marks a message as read
- **THEN** the system sets the `\Seen` flag on the IMAP server

### Requirement: Move and delete messages

The system SHALL move messages between folders and delete messages.

#### Scenario: Move message

- **WHEN** the user sends `POST /api/v1/mail/:mailbox/:uid/move` with a target folder
- **THEN** the message is moved on the IMAP server

#### Scenario: Delete message

- **WHEN** the user sends `DELETE /api/v1/mail/:mailbox/:uid`
- **THEN** the message is marked `\Deleted` and expunged or moved to trash per server config

#### Scenario: Empty trash

- **WHEN** the user sends `DELETE /api/v1/mail/:mailbox` for the trash folder
- **THEN** all messages in trash are permanently deleted

### Requirement: Attachment download

The system SHALL serve attachment content via `GET /api/v1/mail/:mailbox/:uid/attachment/:part`.

#### Scenario: Download attachment

- **WHEN** the user requests an attachment by part ID
- **THEN** the system streams the binary content with correct Content-Type and Content-Disposition headers

### Requirement: Safe HTML rendering

The system SHALL sanitize HTML email bodies server-side using a whitelist-based HTML sanitizer (e.g., bluemonday) before returning them to the frontend.

#### Scenario: Sanitized HTML body

- **WHEN** a message contains HTML with script tags or event handlers
- **THEN** the API returns sanitized HTML with dangerous elements stripped

### Requirement: Remote image blocking

The system SHALL block remote images in HTML messages by default and allow the user to opt in to display them.

#### Scenario: Blocked remote images

- **WHEN** the user opens a message with external `<img src="https://...">` tags
- **THEN** the reading pane replaces remote image URLs with placeholders until the user clicks "Show images"
