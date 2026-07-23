# Skins — referência rápida (dev)

**Guia completo de implementação:** [docs/skins.md](../../../docs/skins.md)

## Nova skin (recomendado)

```bash
make new-skin ID=acme REGISTER=1   # scaffold + Go + manifest + CSS import
make validate-skins
# editar frontend/src/skins/acme.css
# config.toml → [ui] skin = "acme"
make frontend-dev && make run
```

## Arquivos

| Arquivo | O que fazer |
|---------|-------------|
| `<id>.css` | Tokens `[data-skin='id']` + `.dark` |
| `manifest.ts` | Entrada com `id`, `label`, `ready`, `aliases` |
| `internal/ui/skins.go` | Mesmo id/label/ready/aliases (servidor) |
| `index.css` | `@import "./<id>.css";` |
| `_template.css` | Base para novas skins — não editar |

## Variáveis obrigatórias

**Inbox:** `--color-accent`, `--color-accent-2`, `--color-accent-bar`, `--color-line`, `--color-panel`, `--color-panel-2`, `--color-app-bg`, `--color-ink`, `--color-ink-sub`, `--color-ink-mute`, `--color-row-selected`, `--font-sans`

**Login:** `--skin-login-bg`, `--skin-login-card`, `--skin-login-header-bg`, `--skin-login-header-border`, `--skin-login-text`, `--skin-login-input-bg`, `--skin-login-input-border`, `--skin-login-input-text`, `--skin-login-btn-bg`, `--skin-login-btn-text`, `--skin-login-error-bg`, `--skin-login-error-text`, `--skin-login-error-border`

Lista completa + exemplos Gmail/SnappyMail: [docs/skins.md](../../../docs/skins.md#variáveis-css--referência-completa)

## Regras

1. **Sem hex em `.vue`** — só classes Tailwind / tokens CSS
2. **Sem `if (skin === 'gmail')`** — use `--skin-login-*` e `--color-*`
3. **Sempre bloco `.dark`** por skin
4. **`ready: false`** até layout completo (banner em App.vue)
5. **`make validate-skins`** após qualquer mudança no catálogo

## Exemplos no repo

| Skin | Arquivo | Estilo login |
|------|---------|--------------|
| SnappyMail | `snappymail.css` | Escuro, card azul |
| Gmail | `gmail.css` | Claro, botão vermelho |
| Outlook | `outlook.css` | Fundo azul MS, card branco |
