# Skins — developer quick reference

**Full guide:** [docs/skins.md](../../../docs/skins.md)

## New skin in one command

```bash
make new-skin ID=acme
```

Then complete the printed checklist (Go registry, TypeScript, `@import`, `config.toml`).

## Files in this folder

| File | Role |
|------|------|
| `_template.css` | Copy-paste starter (all CSS variables) |
| `*.css` | One file per skin — `[data-skin='id']` tokens |
| `registry.ts` | Labels + `ready` flag for SPA |
| `bootstrap.ts` | Loads skin from `GET /api/v1/ui/config` |
| `apply.ts` | Sets `data-skin` on `<html>` |

## Rules

- Use **CSS variables only** for colors in shared Vue components.
- Add **`.dark`** block for each skin.
- Set `ready: false` until layout components exist (shows preview banner).
