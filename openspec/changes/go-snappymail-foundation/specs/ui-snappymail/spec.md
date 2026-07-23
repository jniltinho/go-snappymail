## ADDED Requirements

### Requirement: SnappyMail 3-column layout

The frontend SHALL render a responsive 3-column layout matching SnappyMail: folder sidebar (left), message list (center), reading pane (right).

#### Scenario: Desktop layout

- **WHEN** the viewport width is ≥ 1024px
- **THEN** all three columns are visible simultaneously

#### Scenario: Mobile layout

- **WHEN** the viewport width is < 768px
- **THEN** the UI shows one column at a time with navigation between folder list, message list, and reading pane

### Requirement: Login screen

The frontend SHALL display a login form with email, password, and optional server override fields styled consistently with SnappyMail's login page.

#### Scenario: Login form display

- **WHEN** an unauthenticated user visits the app
- **THEN** the login view is shown instead of the mail interface

### Requirement: Keyboard shortcuts

The frontend SHALL support SnappyMail-compatible keyboard shortcuts for mail navigation and actions.

#### Scenario: Navigate messages

- **WHEN** the user presses `j` or `k` with no input focused
- **THEN** the selection moves to the next or previous message

#### Scenario: Compose shortcut

- **WHEN** the user presses `c`
- **THEN** the compose modal opens

#### Scenario: Reply shortcut

- **WHEN** the user presses `r` with a message selected
- **THEN** a reply compose is opened for that message

### Requirement: Dark mode

The frontend SHALL support light and dark themes with a toggle in the toolbar or settings.

#### Scenario: Toggle dark mode

- **WHEN** the user toggles dark mode
- **THEN** the entire UI switches color scheme without page reload

### Requirement: Real-time new mail notification

The frontend SHALL connect to SSE endpoint `/api/v1/events` and display a notification when new mail arrives.

#### Scenario: New mail notification

- **WHEN** a new message arrives in INBOX
- **THEN** the message list refreshes and a toast notification is shown

### Requirement: Visual design reference

The frontend SHALL use Tailwind CSS v4 with color palette and spacing inspired by SnappyMail's Default and NightShine themes.

#### Scenario: Theme colors

- **WHEN** the app renders in default light theme
- **THEN** the primary accent, sidebar background, and message list styling visually resemble SnappyMail Default theme
