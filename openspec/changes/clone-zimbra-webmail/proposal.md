# Proposal: clone-zimbra-webmail

## Why

The client asked to port the Zimbra (Java) web client layout to the Golang webmail. The `add-zimbra-skin` change delivered the palette and core chrome, but a side-by-side against the live reference (Zimbra FOSS 10.1.17 Classic at https://192.168.56.30, account `nilton@linuxpro.com.br`) shows the clone is not yet indistinguishable: fonts don't match (14px Tailwind sizes vs Zimbra's 12px), dropdown menus (▾) are static, several hover/press effects and secondary controls are missing. This change closes the gap with a **designer + QA workflow**: the designer track extracts measured truth from the running reference; the QA track (agent `qa-frontend-cloner`) gates every iteration with a severity-ranked parity audit.

## What Changes

- **QA infrastructure**: `.claude/agents/qa-frontend-cloner.md` — a QA agent expert in frontend cloning; it drives `agent-browser` over both UIs (same viewport, same mailbox — the test instance points at the VM's IMAP so both render `nilton@linuxpro.com.br`), captures computed styles/snapshots/screenshots, probes interactions, and reports P1/P2/P3 findings. Acceptance gate: **zero P1**.
- **Typography parity (P1)**: Zimbra Classic renders at 12px; the clone's Tailwind `text-sm`/`text-base` resolve to 14/16px. Override the type scale tokens under `[data-skin='zimbra']` so utilities resolve to Zimbra sizes.
- **Working menus (P1)**: real dropdowns for controls that show ▾ — user menu (top right), Actions menu in the toolbar, New message split-button menu. Simple Vue dropdown (click-toggle, click-outside close), styled per skin.
- **Effects (P2)**: hover/active states audited per control against the reference (buttons, rows, tree items, tabs, sash) and fixed where missing.
- **Secondary controls parity (P2)**: toolbar Read More/View group (right side), folder-section gear affordances, list checkbox column — implemented or explicitly waived in Non-goals after client review.
- **QA loop**: run `qa-frontend-cloner` after each iteration; findings feed tasks until the audit returns zero P1 and the client signs off on remaining P3s.

## Capabilities

### New Capabilities

- `zimbra-clone-parity`: measurable parity requirements between the Golang webmail (zimbra skin) and the live Zimbra Classic reference — typography, menus, effects, controls — with the QA agent audit as the acceptance mechanism.

### Modified Capabilities

- `zimbra-skin`: type-scale override and new chrome (dropdown menus, toolbar right group) extend the skin's requirements from `add-zimbra-skin`.

## Non-goals

- Features the Go backend doesn't have yet: Briefcase, Tasks, Calendar, Preferences content (tabs stay visual-only), conversation threading, tags/saved searches, mini-calendar, Zimlets.
- Zimbra logo/trademark assets (text branding stays), rounded corners (project squared rule).
- Print flow behind the print dropdown (button may render, action deferred).
- Rich composer (full Compose tab with formatting toolbar, attach row, spell-check, Options) — the clone ships a basic modal composer; the rich composer (TipTap) is tracked in `go-snappymail-foundation`. Recorded as the single remaining P2 in the final audit.

## Impact

- `frontend/src/skins/zimbra.css` + `style.css` — type tokens, menu styles, missing states.
- `frontend/src/components/` — small `DropdownMenu.vue`, AppToolbar/user-menu wiring.
- `.claude/agents/qa-frontend-cloner.md` — QA agent (already added).
- Reference environment: `vagrant/zimbra` VM (docs/zimbra-vagrant-foss.md); test instance config points IMAP/SMTP at the VM.
- No Go/API changes expected.
