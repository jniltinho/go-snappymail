# Skins (layout themes)

Skins control colors, typography, and future layout variants (SnappyMail, Gmail, Outlook).

## Server config (`config.toml`)

```toml
[ui]
skin = "snappymail"   # snappymail | gmail | outlook
```

Legacy alias: `theme = "snappymail-default"` maps to `snappymail`.

The SPA loads the active skin from `GET /api/v1/ui/config` on startup.

## Frontend structure

```
frontend/src/skins/
├── types.ts          # SkinId, UIConfigResponse
├── registry.ts       # metadata + normalizeSkinId()
├── apply.ts          # sets data-skin on <html>
├── bootstrap.ts      # fetch server config before mount
├── snappymail.css    # active skin (full)
├── gmail.css         # placeholder tokens
└── outlook.css       # placeholder tokens
```

## Adding a new skin

1. Add id in `internal/ui/skins.go` (`AvailableSkins`, `NormalizeSkin`)
2. Create `frontend/src/skins/<id>.css` with `[data-skin='<id>']` CSS variables
3. Register in `frontend/src/skins/registry.ts` (`SKIN_REGISTRY`, set `ready: true` when layout is done)
4. Import CSS in `frontend/src/style.css`
5. Optionally add skin-specific Vue layout components under `frontend/src/skins/<id>/`

Dark mode is client-side (`localStorage`) and combines with `[data-skin].dark` CSS blocks.
