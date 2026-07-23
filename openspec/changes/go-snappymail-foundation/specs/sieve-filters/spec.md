## ADDED Requirements

### Requirement: ManageSieve connection

The system SHALL connect to the user's ManageSieve server (when configured) to list and manage Sieve filter scripts.

#### Scenario: List sieve scripts

- **WHEN** the user opens the filters settings page
- **THEN** the system returns all Sieve script names and their active status from the ManageSieve server

### Requirement: Sieve script editor

The frontend SHALL provide a visual and raw-text editor for creating and editing Sieve filter rules.

#### Scenario: Create filter rule

- **WHEN** the user creates a rule "move emails from spam@example.com to Junk"
- **THEN** the system generates valid Sieve syntax and uploads the script to ManageSieve

#### Scenario: Activate script

- **WHEN** the user activates a script
- **THEN** the system sets it as the active script on the ManageSieve server

#### Scenario: Delete script

- **WHEN** the user deletes a script
- **THEN** the script is removed from the ManageSieve server

### Requirement: Sieve availability detection

The system SHALL detect whether ManageSieve is available for the user's account based on domain/server configuration.

#### Scenario: Sieve not configured

- **WHEN** no Sieve server is configured for the domain
- **THEN** the filters section is hidden or shows "not available"
