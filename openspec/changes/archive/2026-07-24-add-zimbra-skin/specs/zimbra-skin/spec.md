# Spec: zimbra-skin

## ADDED Requirements

### Requirement: Zimbra skin is registered and validated
The system SHALL register a `zimbra` skin in the three sync points — Go catalog (`internal/ui/skins.go`), TS manifest (`frontend/src/skins/manifest.ts`), and CSS imports (`frontend/src/skins/index.css`) — with `ready: true` and alias `classic`.

#### Scenario: validate-skins passes
- **WHEN** `make validate-skins` runs after registration
- **THEN** it exits 0 with `zimbra` present in Go, TS, and CSS listings

#### Scenario: Skin exposed via API
- **WHEN** a client calls `GET /api/v1/ui/config`
- **THEN** `zimbra` appears in `available_skins` and in `skins[]` with `ready: true`

### Requirement: Zimbra is the default skin
The system SHALL use `zimbra` as the default skin: `defaultSkinID` in `internal/ui/skins.go`, `DEFAULT_SKIN` in `manifest.ts`, and `[ui] skin = "zimbra"` in `web/files/config.default.toml`.

#### Scenario: Empty or unknown config value
- **WHEN** `config.toml` has no `[ui] skin` value, or an unknown value
- **THEN** `NormalizeSkin` resolves to `zimbra` and the SPA renders with `data-skin="zimbra"`

#### Scenario: Explicit skin preserved
- **WHEN** `config.toml` sets `[ui] skin = "carbonio"`
- **THEN** the served skin remains `carbonio` (no forced migration)

### Requirement: Zimbra 8 Classic visual identity
The skin SHALL reproduce the Zimbra 8.8.15 Classic look via CSS tokens under `[data-skin='zimbra']`: blue top bar with white branding area, light chrome with gray-gradient squared toolbar buttons, yellow focus fill on text inputs, and Arial/Helvetica type. All UI text SHALL be English. No `border-radius` SHALL be introduced (global squared rule).

#### Scenario: Main app chrome
- **WHEN** the mail view renders with the zimbra skin
- **THEN** the top bar is Zimbra blue, the toolbar buttons are squared with light gray fill and 1px borders, and the selected folder row uses the Zimbra soft-blue selection color

#### Scenario: Focused input highlight
- **WHEN** a text input (search, compose To/Subject) receives focus
- **THEN** its background is the Zimbra yellow focus fill

### Requirement: Zimbra login page
The login page under the zimbra skin SHALL show a solid Zimbra-blue card (`#007cc3`) centered on a white page (subtle white→`#ededed` gradient), with white labels, white input fields, and a squared Login button — matching the reference capture of webmail.criarenet.com. The brand area SHALL be text-only (no Zimbra logo). No "Stay signed in" checkbox (deferred — needs backend remember-me).

#### Scenario: Login render
- **WHEN** an unauthenticated user opens `/`
- **THEN** the page background is white, the card is Zimbra blue with a text-only brand area, and all texts are English

#### Scenario: No wrong-skin flash
- **WHEN** the SPA loads before `/api/v1/ui/config` resolves
- **THEN** the pre-bootstrap `data-skin` in `frontend/index.html` is `zimbra`, so no other skin flashes first

### Requirement: Visual parity with reference
The implemented skin SHALL be iterated (screenshot loop against `docs/prints/zimbra/` live captures of webmail.criarenet.com) until login page and mail view are visually equivalent to the reference — same palette, same chrome hierarchy, same density — allowing only for the squared-corners rule and text-only branding. This is the acceptance bar for the change.

#### Scenario: Side-by-side check
- **WHEN** login and inbox screenshots of go-snappymail (zimbra skin) are compared side-by-side with the reference captures
- **THEN** top bar, tab strip, toolbar, folder tree, list rows, selection colors, and login card are indistinguishable in palette and layout hierarchy

### Requirement: Dark variant
The skin SHALL define a `.dark` variant with legible token overrides for all surfaces, text, and login tokens.

#### Scenario: Dark toggle
- **WHEN** the user enables dark mode (`gsn_dark`)
- **THEN** panels, text, lines, and selection colors switch to the dark token set with no unreadable contrast pairs
