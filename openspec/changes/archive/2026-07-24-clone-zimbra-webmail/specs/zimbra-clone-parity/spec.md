# Spec: zimbra-clone-parity

Reference: Zimbra FOSS 10.1.17 Classic (advanced) at https://192.168.56.30, account `nilton@linuxpro.com.br`. Clone: go-snappymail zimbra skin, same mailbox via the VM's IMAP. Auditor: agent `qa-frontend-cloner`.

## ADDED Requirements

### Requirement: Typography parity
Under the zimbra skin, UI text SHALL resolve to Zimbra Classic metrics: 12px base for chrome text (rows, buttons, labels, tree, tabs), matching weight (400 normal / 700 unread-bold) and the harmony font stack.

#### Scenario: Type-scale audit
- **WHEN** the QA agent measures computed font-size/family/weight on toolbar buttons, list rows, tree items, and tabs in both UIs
- **THEN** family (first resolved), size, and weight are equal

### Requirement: Functional dropdown menus
Every control rendered with a ▾ SHALL open a real menu: user menu (Dark mode, Logout), toolbar Actions menu (Mark read/unread, Flag, Spam), New message split-button menu (New message). Menus close on click-outside and Esc.

#### Scenario: Menu probe
- **WHEN** the QA agent clicks each ▾ control
- **THEN** a menu opens with the listed items and each item performs its action

### Requirement: Interaction effects parity
Hover and active states SHALL match the reference: toolbar buttons (#CCE5F3 hover / #99CAE7 active), list rows and tree items (#CCE5F3 hover), tabs, and sash hover, measured via computed style after hover.

#### Scenario: Hover audit
- **WHEN** the QA agent hovers each control class in both UIs and re-reads background-color
- **THEN** the values are equal (or the finding is listed as intentional divergence)

### Requirement: QA gate
The change SHALL be accepted only when a full `qa-frontend-cloner` audit reports **zero P1 findings**, with remaining P2/P3 explicitly listed and either fixed or waived in the proposal Non-goals.

#### Scenario: Final audit
- **WHEN** the QA agent runs the full per-region audit (top bar, tabs, toolbar, tree, list, reading pane, composer, login, dark)
- **THEN** the report verdict shows 0 P1
