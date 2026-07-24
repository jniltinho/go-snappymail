# Proposal: add-zimbra-skin

## Why

go-snappymail ships four skins (snappymail, gmail, outlook, carbonio) but none reproduces the classic Zimbra 8 web client — the UI the CRIARE fleet users already know from webmail.criarenet.com (Zimbra 8.8.15, Harmony skin). Adding a faithful `zimbra` skin and making it the default lets those users migrate to go-snappymail with zero visual retraining, and a review pass over the core mail flows ensures the webmail is actually usable day-to-day under that skin.

## What Changes

- New `zimbra` skin (token-driven CSS, `[data-skin='zimbra']` + `.dark` variant) reproducing Zimbra 8 Classic, from live capture of webmail.criarenet.com (reference screenshots in `docs/prints/zimbra/`, gitignored):
  - Zimbra blue top bar (`#007cc3`) with text-only branding area (no Zimbra logo — trademark) and user menu on the right.
  - Tab-strip look under the top bar: light blue strip (`rgba(0,135,195,0.1)` over white), active tab white, inactive tabs white text on blue.
  - Squared toolbar buttons on light chrome; primary action button in Zimbra blue (`#0087c3`).
  - Selection color `#99cae7` for folder tree and message rows; list header `#f2f2f2`; unread rows bold `#333`.
  - Yellow toast/notice surfaces (`#ffffcc`) and yellow focus fill on compose/search inputs.
  - Login page: Zimbra-blue card (`#007cc3`) centered on white page (subtle white→`#ededed` gradient), white labels/inputs, squared Login button.
  - Squared everywhere — respects the global no-border-radius rule (the real Zimbra uses 3–5px radii; we deliberately keep 0).
- Registered in all sync points: `internal/ui/skins.go` (catalog), `frontend/src/skins/manifest.ts`, `frontend/src/skins/index.css`; `make validate-skins` passes. Alias: `classic` (no collision with existing aliases).
- `ready: true` from day one (no preview banner).
- **BREAKING (default change)**: `zimbra` becomes the default skin. All places that hardcode the current default change together:
  - `defaultSkinID` in `internal/ui/skins.go` (+ fallback cases in `internal/ui/skins_test.go`)
  - `DEFAULT_SKIN` in `frontend/src/skins/manifest.ts`
  - `data-skin="snappymail"` pre-bootstrap attribute in `frontend/index.html` (prevents wrong-skin flash before `/ui/config` resolves)
  - `[ui] skin` value and comment line in `web/files/config.default.toml`
  - Docs that state the default: `AGENTS.md` (skins summary), `docs/skins.md` (catalog table, default-fallback text, sample validate output, sample JSON)
  - Note: deployments with an explicit `[ui] skin` are unaffected; the Docker lab sets no skin and will intentionally flip to zimbra.
- `LoginView.vue` brand area: small markup tweak so the skin can render the Zimbra-style brand block as text (product name styled white on blue). No logo asset.
- Core mail flows under the new skin — split honestly between *verify* and *add*:
  - **Verify (exists today)**: search, message list (flat — no threading), read pane (sanitized HTML + inline CIDs), flag toggle, delete, folder counts.
  - **Add (small UI gaps found in review)**: Archive toolbar action (uses existing `POST /mail/:mailbox/:uid/move`), mark read/unread action (uses existing `flag=seen` API), draft/send verified at API level (`/compose/draft`, `/compose/send`) — the composer UI itself is tracked in the `go-snappymail-foundation` change, not here.

## Capabilities

### New Capabilities

- `zimbra-skin`: the Zimbra 8 Classic look — tokens, login page, chrome details, dark variant, registration in the sync points, and default-skin behavior.
- `core-mail-flows-review`: acceptance checks for the mail flows under the zimbra skin — verification of existing flows plus the two small UI additions (archive action, mark read/unread).

### Modified Capabilities

<!-- none — no main specs exist yet; foundation change predates spec extraction -->

## Non-goals

- No per-skin Vue layout (`--skin-layout` stays a hint; no AppShell fork) — the skin is CSS tokens + scoped chrome rules, plus the one LoginView brand-area markup tweak.
- No Zimbra logo or trademarked assets — text branding only.
- No "Stay signed in" checkbox — it needs backend remember-me session semantics; deferred.
- No pixel-perfect clone of Zimbra's DWT widget toolkit, tab-per-message navigation, Briefcase, Tasks, or Zimlets.
- No message threading, no labels beyond IMAP flags, no Sieve.
- No composer UI work here (drafts/send checked at API level; composer belongs to `go-snappymail-foundation`).
- No i18n framework — English only, as the SPA already is.

## Impact

- `internal/ui/skins.go` + `internal/ui/skins_test.go` — new catalog entry, `defaultSkinID` change, updated fallback tests.
- `frontend/src/skins/{manifest.ts,index.css,zimbra.css}` — new skin files; `DEFAULT_SKIN` change affects pre-login bootstrap.
- `frontend/index.html` — pre-bootstrap `data-skin` attribute.
- `frontend/src/components/LoginView.vue` — brand-area tweak (text-only).
- `frontend/src/components/AppToolbar.vue` + `frontend/src/stores/mail.ts` — Archive action and mark read/unread action (existing backend APIs; no Go changes).
- `web/files/config.default.toml` — default `[ui] skin = "zimbra"` + comment.
- `AGENTS.md`, `docs/skins.md` — default-skin references and catalog table.
- `GET /api/v1/ui/config` — reports the new default; API shape unchanged.
- go-cubemail reference: skin registration mirrors the same single-catalog pattern already used here; no cross-repo dependency.
