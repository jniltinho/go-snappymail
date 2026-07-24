# Spec: core-mail-flows-review

Acceptance checks for the core mail flows rendered under the zimbra skin. Requirements are split between **verification** of flows that exist today (a failing check is a defect to fix) and two **small UI additions** identified in review (archive action, mark read/unread). The composer UI is out of scope (tracked in `go-snappymail-foundation`); drafts and send are checked at API level.

## ADDED Requirements

### Requirement: Search emails (verify)
Search SHALL query via `GET /api/v1/search` and render results in the message list with the zimbra skin styles.

#### Scenario: Search by text
- **WHEN** the user types a term in the search field and submits
- **THEN** matching messages render in the list and clearing the search restores the full mailbox listing

### Requirement: List messages (verify)
`GET /api/v1/mail/:mailbox` SHALL render pages of messages (flat list — threading is a non-goal) with sender, subject, date, unread state, and attachment indicator; selection SHALL use the zimbra selection color `#99cae7`.

#### Scenario: Unread visual state
- **WHEN** a mailbox with unread messages is listed
- **THEN** unread rows render bold (`font-weight: 700`) and the folder unread count matches

### Requirement: Read message content (verify)
Opening a message SHALL render sanitized HTML (bluemonday, `data:` CIDs inline), sender/recipient info, date, and attachment list.

#### Scenario: HTML message
- **WHEN** the user opens an HTML message with inline images
- **THEN** content renders sanitized with inline CIDs and no layout bleed outside the reading pane

### Requirement: Drafts via API (verify)
`POST /api/v1/compose/draft` SHALL persist a draft to the Drafts folder.

#### Scenario: Save draft via API
- **WHEN** a logged-in client posts to/subject/body to `/api/v1/compose/draft`
- **THEN** the message appears in the Drafts folder listing

### Requirement: Send via API (verify)
`POST /api/v1/compose/send` SHALL deliver via SMTP, append to Sent, and return errors in English.

#### Scenario: Successful send via API
- **WHEN** a logged-in client posts a valid message to `/api/v1/compose/send`
- **THEN** the message is delivered and appears in the Sent folder

### Requirement: Flag messages (verify)
`POST /api/v1/mail/:mailbox/:uid/flag` SHALL toggle the \Flagged IMAP flag and the list SHALL reflect the state immediately.

#### Scenario: Toggle flag
- **WHEN** the user flags a message
- **THEN** the flag icon fills in the list without a full reload

### Requirement: Mark read and unread (add)
The toolbar SHALL offer mark read/unread for the open or selected message, wired to the existing `flag=seen` API, reflected in row style and folder unread count.

#### Scenario: Mark unread
- **WHEN** the user marks a read message as unread
- **THEN** the row returns to the bold unread style and the folder count increments

### Requirement: Archive action (add)
The toolbar SHALL offer an Archive action that moves the message to the Archive folder via the existing `POST /api/v1/mail/:mailbox/:uid/move`, creating/resolving the Archive folder by name.

#### Scenario: Archive from toolbar
- **WHEN** the user clicks Archive on an open message
- **THEN** the message leaves the current list and appears in the Archive folder
