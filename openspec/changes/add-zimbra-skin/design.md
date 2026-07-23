# Design: add-zimbra-skin

## Context

go-snappymail skins are pure CSS token sets scoped to `[data-skin='<id>']` (see `carbonio.css` as the closest precedent — also a live-capture clone). The SPA chrome (LoginView, AppToolbar, FolderSidebar, MessageList, ReadingPane) consumes `--color-*` / `--skin-login-*` tokens via `style.css` and the Tailwind `@theme`. The reference is Zimbra 8.8.15 Advanced (Ajax) client, Harmony skin, captured live from webmail.criarenet.com on 2026-07-23 (screenshots in `docs/prints/zimbra/`, gitignored; computed styles extracted via browser eval).

The quality bar set for this change: **iterate until go-snappymail under the zimbra skin is visually indistinguishable from the reference captures** (palette, chrome hierarchy, density), allowing only the squared-corners rule and text-only branding.

## Goals / Non-Goals

**Goals:**
- Faithful Zimbra 8 Harmony palette and chrome as a token-driven skin; `zimbra` becomes the default everywhere the default is encoded.
- Close the two small flow gaps found in review (Archive action, mark read/unread) using existing backend APIs.
- Screenshot-loop acceptance against the reference captures.

**Non-Goals:**
- Zimbra logo/assets, "Stay signed in" remember-me, threading, composer UI, per-skin Vue layouts, i18n (see proposal Non-goals).

## Design plan (frontend-design)

The brief pins the direction completely — this is a faithful reproduction, so every value below is measured from the live capture, not invented. The one aesthetic risk is deliberate fidelity: keeping Zimbra's 2011-era flat-gray toolbar chrome instead of "modernizing" it.

**Palette (measured, named):**

| Token | Value | Source (measured element) |
|---|---|---|
| `zimbra-blue` | `#007cc3` | top bar `skin_spacing_top_row`, login card `.contentBox`, login button |
| `zimbra-blue-action` | `#0087c3` | "New message" toolbar button |
| `zimbra-selection` | `#99cae7` | selected tree item + selected list row + login input border |
| `zimbra-strip` | `rgba(0,135,195,0.1)` ≈ `#e5f3f9` | app tab strip row |
| `zimbra-header` | `#f2f2f2` | list column header |
| `zimbra-btn-border` | `#bfbfbf` | toolbar button 1px border, subject input border |
| `zimbra-ink` | `#333333` | row text, unread bold text |
| `zimbra-notice` | `#ffffcc` | toast bg, status info — reused as input focus fill |
| `zimbra-page-fade` | `#ededed` | login page gradient end, dialogs `#d3d3d3` |
| `zimbra-footer-ink` | `#656565` | login footer text |

**Type:** Segoe UI 12px is Zimbra's UI face (measured `font-family: "Segoe UI"; font-size: 12px`), weight 400, unread rows 700. Roles: UI text = `'Segoe UI', 'Helvetica Neue', Helvetica, Arial, sans-serif`; no display face — the chrome IS the identity. Density is Zimbra's: compact 12px rows.

**Layout concept (existing 3-column shell, re-skinned):**

```
┌────────────────────────────────────────────────────────────┐
│ ███ zimbra-blue top bar: brand-text · search · user ▾      │
│ ▒▒▒ strip: [Mail]▉white-active-tab   inactive: white text  │
│ [New message ▾]  [Reply][Reply all][Forward][Archive][Del] │ ← white toolbar, 1px #bfbfbf btns
├──────────┬───────────────────────┬─────────────────────────┤
│ folders  │ list: hdr #f2f2f2     │ reading pane (white)    │
│ selected │ row: ●bold sender/date│ subject #333, hdr line  │
│ #99cae7  │ 2nd line: subj+snippet│ sanitized HTML body     │
└──────────┴───────────────────────┴─────────────────────────┘
```

**Signature element:** the Zimbra-blue login card — solid `#007cc3` block on a white→`#ededed` gradient page, white text-only brand block ("go-snappymail" set in white, small caps sub-line), fields right-aligned to labels. It is the one moment users instantly recognize as "Zimbra".

**Critique pass (against generic defaults):** none of the three AI-default looks apply — palette, type, and chrome are measured from the subject. The template answer would be a soft-shadow rounded card and a 14px type scale; we keep 12px, hard edges, flat fills, exactly like the reference.

## Decisions

1. **Token-only skin + minimal scoped chrome rules** (like `carbonio.css`), not a per-skin layout. Rationale: the existing shell already matches Zimbra's 3-column hierarchy; tokens + a few `[data-skin='zimbra']` scoped rules (tab-strip look on the toolbar area, uppercase-free buttons, selection colors) reach parity. Alternative (AppShell fork per `--skin-layout`) rejected: large surface, not needed for parity.
2. **Default change is a constant flip in 4 code points + docs**, no migration logic. `NormalizeSkin` already handles unknown→default. Alternative (config migration) rejected: explicit `[ui] skin` values keep working.
3. **Archive + read/unread as store actions + toolbar buttons**, no Go changes. Archive resolves the folder named `Archive` (create via existing `POST /folders` if missing). Alternative (new backend endpoint) rejected: `move` + `flag=seen` APIs already exist.
4. **Focus-fill yellow `#ffffcc`** applied via scoped rule on text inputs — Zimbra's signature focus behavior from the HTML client, also used by toasts. Kept subtle (background only, border stays).
5. **Squared corners kept** (project rule) even though the real Zimbra uses 3–5px radii on inputs/buttons. Documented divergence; acceptance allows it.
6. **Text-only brand** (no logo asset): brand block renders the product name styled like Zimbra's banner area. Trademark-safe.

## Risks / Trade-offs

- [Visual parity is subjective] → acceptance is the side-by-side screenshot loop on the three reference captures; iterate until indistinguishable per spec scenario.
- [Default flip surprises lab/docker users] → intentional; release notes line in commit message.
- [Archive folder may not exist on some IMAP servers] → resolve-or-create by name on first use; failure surfaces the API error toast.
- [Segoe UI absent on Linux clients] → font stack falls back to Helvetica/Arial; density unchanged.

## Migration Plan

Ship in one PR: skin files + defaults + toolbar actions + docs. Rollback = revert commit (no data/schema changes). Deployments pinning `[ui] skin` see no change.

## Open Questions

None — gap review resolved scope ambiguities (composer stays in foundation change; checkbox deferred).
