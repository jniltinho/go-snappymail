# Webmail skins

Skins define the **look and feel** of the SPA: colors, typography, login page, and (in the future) layout variants (SnappyMail 3-column, Gmail, Outlook, custom).

The active skin is set **server-side** in `config.toml` — all users see the same skin unless you add per-user overrides later.

## Quick switch (operator)

```toml
[ui]
skin = "snappymail"   # snappymail | gmail | outlook | your-id
```

```bash
# or via environment
GOSM_UI_SKIN=gmail ./dist/go-snappymail serve
```

Restart the server and reload the browser. The SPA reads `GET /api/v1/ui/config` on startup.

| Skin | Status | Notes |
|------|--------|-------|
| `snappymail` | Ready | Default — SnappyMail-inspired blue |
| `gmail` | Tokens only | Red accent palette; full layout TBD |
| `outlook` | Tokens only | Microsoft blue palette; full layout TBD |

Legacy alias: `ui.theme = "snappymail-default"` maps to `snappymail`.

---

## Create a new skin (5 steps)

### 1. Scaffold CSS from template

```bash
make new-skin ID=acme
# creates frontend/src/skins/acme.css from _template.css
```

Or copy manually:

```bash
cp frontend/src/skins/_template.css frontend/src/skins/acme.css
# replace MY_SKIN_ID with acme in the new file
```

### 2. Register on the server (Go)

Edit [internal/ui/skins.go](../internal/ui/skins.go):

```go
const SkinAcme = "acme"

var available = []string{SkinSnappyMail, SkinGmail, SkinOutlook, SkinAcme}

// In NormalizeSkin switch:
case "acme":
    return SkinAcme
```

Run `go test ./internal/ui/...`.

### 3. Register on the frontend (TypeScript)

| File | Change |
|------|--------|
| [frontend/src/skins/types.ts](../frontend/src/skins/types.ts) | Add `'acme'` to `SkinId` union |
| [frontend/src/skins/registry.ts](../frontend/src/skins/registry.ts) | Add entry to `SKIN_REGISTRY` and `normalizeSkinId()` |
| [frontend/src/style.css](../frontend/src/style.css) | Add `@import "./skins/acme.css";` |

Set `ready: true` in the registry when the skin has a complete layout (not just colors).

### 4. Build and test

```bash
# config.toml
[ui]
skin = "acme"

make frontend-dev   # terminal 1 — Vite :5173
make run            # terminal 2 — Go :8082
```

Toggle dark mode in the toolbar to verify `.dark` CSS variables.

### 5. Ship

```bash
make build-prod
make test
```

---

## CSS variables reference

Each skin file targets `[data-skin='your-id']`. The SPA sets `data-skin` on `<html>` at boot.

### Required (colors + login)

| Variable | Used for |
|----------|----------|
| `--color-accent` | Primary brand, buttons |
| `--color-accent-2` | Hover / secondary accent |
| `--color-accent-bar` | Headers, login bar |
| `--color-line` | Borders |
| `--color-panel` | Main panels |
| `--color-panel-2` | Sidebar / hover backgrounds |
| `--color-app-bg` | Page background |
| `--color-ink` | Primary text |
| `--color-ink-sub` | Secondary text |
| `--color-ink-mute` | Muted text |
| `--color-row-selected` | Selected message row |
| `--skin-login-bg` | Login page background |
| `--skin-login-card` | Login form card |
| `--skin-login-input-bg` | Input background |
| `--skin-login-input-border` | Input border |
| `--font-sans` | Font stack |

### Optional

| Variable | Purpose |
|----------|---------|
| `--skin-layout` | Hint for future layout switch (`three-column`, `gmail`, …) |

### Dark mode

Duplicate variables under `[data-skin='your-id'].dark { ... }`. Dark mode is toggled client-side (`localStorage` key `gsn_dark`); server config does not control it today.

---

## Architecture

```
config.toml  ui.skin
      │
      ▼
GET /api/v1/ui/config  ──►  bootstrap.ts  ──►  applySkin()  ──►  <html data-skin="…">
      │                              │
      │                              └──►  settings store (Pinia)
      │
      └── available_skins[]  (from internal/ui/skins.go)
```

Shared UI components (`FolderSidebar`, `MessageList`, …) use Tailwind classes mapped to `--color-*` tokens — **no hardcoded hex** in Vue files when possible.

---

## Advanced: custom layout per skin

When a skin needs a different structure (not just colors):

1. Create `frontend/src/skins/<id>/AppShell.vue`
2. In [App.vue](../frontend/src/App.vue), switch layout by `settings.skin` (or dynamic `defineAsyncComponent`)
3. Set `ready: true` only when the shell is complete

Start with CSS-only skins; add Vue layouts only when the reference product differs structurally (e.g. Gmail conversation view).

---

## API

See [api.md](api.md#get-apiv1uiconfig) — `GET /api/v1/ui/config` returns:

```json
{
  "skin": "snappymail",
  "available_skins": ["snappymail", "gmail", "outlook"],
  "rows_per_page": 50,
  "datetime_format": "02/01/2006 15:04",
  "compose_html": true
}
```

---

## Checklist (copy when adding a skin)

- [ ] `frontend/src/skins/<id>.css` (from `_template.css`)
- [ ] `@import` in `frontend/src/style.css`
- [ ] `internal/ui/skins.go` — constant, `available`, `NormalizeSkin`
- [ ] `frontend/src/skins/types.ts` — `SkinId`
- [ ] `frontend/src/skins/registry.ts` — `SKIN_REGISTRY`, `normalizeSkinId`
- [ ] `config.toml` → `ui.skin = "<id>"` tested
- [ ] Light + dark mode checked
- [ ] `go test ./...` and `cd frontend && npm run build`

More detail in [frontend/src/skins/README.md](../frontend/src/skins/README.md) (developer quick reference).
