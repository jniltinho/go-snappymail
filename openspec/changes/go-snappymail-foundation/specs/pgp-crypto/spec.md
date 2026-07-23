## ADDED Requirements

### Requirement: PGP key management

The system SHALL store public keys in the database and keep private keys exclusively on the client (IndexedDB via OpenPGP.js), protected by the user's passphrase — matching SnappyMail's client-side key model.

#### Scenario: Import public key

- **WHEN** the user uploads an ASCII-armored public key
- **THEN** the system parses and stores the public key associated with the user's account

#### Scenario: Import private key

- **WHEN** the user uploads a private key protected by passphrase
- **THEN** the private key is stored in browser IndexedDB via OpenPGP.js and is never sent to the server

### Requirement: Encrypt and sign on send

The frontend SHALL support encrypting and/or signing outgoing messages using OpenPGP.js before sending via the compose API.

#### Scenario: Send encrypted message

- **WHEN** the user enables encryption and selects recipient keys
- **THEN** the message body is encrypted with OpenPGP before SMTP submission

#### Scenario: Send signed message

- **WHEN** the user enables signing with their private key
- **THEN** the message includes a valid OpenPGP signature

### Requirement: Decrypt and verify on read

The frontend SHALL decrypt incoming PGP messages and verify signatures using OpenPGP.js.

#### Scenario: Read encrypted message

- **WHEN** the user opens a PGP-encrypted message and provides their passphrase
- **THEN** the reading pane displays the decrypted plaintext

#### Scenario: Verify signature

- **WHEN** a signed message is opened
- **THEN** the UI shows verification status (valid/invalid/unknown key)
