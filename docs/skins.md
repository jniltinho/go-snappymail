# Skins — guia de implementação

Skins definem **aparência** do webmail: cores, tipografia, tela de login e (no futuro) layout por produto (SnappyMail, Gmail, Outlook, marca customizada).

Este documento explica **como implementar uma skin do zero**, como manter sincronia entre backend e frontend, e as regras para contribuir com CSS/Vue sem quebrar outras skins.

---

## Índice

1. [Conceitos](#conceitos)
2. [Trocar skin (operador)](#trocar-skin-operador)
3. [Tutorial: skin `acme` do zero](#tutorial-skin-acme-do-zero)
4. [Registro manual (passo a passo)](#registro-manual-passo-a-passo)
5. [Arquivos do sistema](#arquivos-do-sistema)
6. [Variáveis CSS — referência completa](#variáveis-css--referência-completa)
7. [Como o Vue consome tokens](#como-o-vue-consome-tokens)
8. [Modo escuro](#modo-escuro)
9. [Flag `ready` e banner de preview](#flag-ready-e-banner-de-preview)
10. [Aliases em `config.toml`](#aliases-em-configtoml)
11. [Validação (`make validate-skins`)](#validação-make-validate-skins)
12. [Erros comuns e soluções](#erros-comuns-e-soluções)
13. [Regras para componentes Vue](#regras-para-componentes-vue)
14. [Layout customizado (futuro)](#layout-customizado-futuro)
15. [API](#api)
16. [Checklist final](#checklist-final)

---

## Conceitos

| Termo | Significado |
|-------|-------------|
| **Skin id** | Identificador lowercase (`snappymail`, `gmail`, `acme-corp`). Usado em `config.toml`, API e `[data-skin='…']` no CSS. |
| **Catálogo Go** | Lista em `internal/ui/skins.go` — fonte da verdade no servidor. |
| **Manifest TS** | Lista em `frontend/src/skins/manifest.ts` — espelho do catálogo no cliente. |
| **Tokens CSS** | Variáveis `--color-*` e `--skin-login-*` definidas por skin. |
| **`ready`** | `true` = skin completa (layout + cores). `false` = só tokens/preview — SPA mostra banner. |
| **Modo escuro** | Toggle **client-side** (`localStorage` `gsn_dark`). Não vem do servidor. |

Fluxo ao abrir o app:

```
config.toml [ui] skin
       ↓
GET /api/v1/ui/config  →  skin id normalizado
       ↓
bootstrap.ts  →  applySkin(id)  →  <html data-skin="acme">
       ↓
CSS da skin ativa sobrescreve --color-* e --skin-login-*
       ↓
Componentes Vue usam classes Tailwind (bg-panel, text-ink, …)
```

A skin ativa é **server-side**: todos os usuários veem a mesma, salvo overrides futuros por usuário.

---

## Trocar skin (operador)

```toml
# config.toml ou web/files/config.default.toml
[ui]
skin = "snappymail"   # snappymail | gmail | outlook | seu-id
```

```bash
# override por ambiente
GOSM_UI_SKIN=gmail ./dist/go-snappymail serve
```

Reinicie o servidor e recarregue o browser.

| Skin | `ready` | Observação |
|------|---------|------------|
| `snappymail` | ✅ | Padrão — azul inspirado no SnappyMail |
| `gmail` | ❌ | Paleta vermelha; layout completo TBD |
| `outlook` | ❌ | Paleta Microsoft; layout completo TBD |

Alias legado: `[ui] theme = "snappymail-default"` → `snappymail`.

Mais opções de config: [configuration.md](configuration.md#ui--skins).

---

## Tutorial: skin `acme` do zero

Exemplo completo — marca fictícia **Acme Corp** com id `acme`.

### Passo 1 — Scaffold + registro automático

```bash
make new-skin ID=acme REGISTER=1
```

Isso cria e registra:

| Ação | Arquivo |
|------|---------|
| CSS da skin | `frontend/src/skins/acme.css` |
| Catálogo servidor | `internal/ui/skins.go` |
| Manifest cliente | `frontend/src/skins/manifest.ts` |
| Import CSS | `frontend/src/skins/index.css` |
| Validação | roda `make validate-skins` |

**Regras do id:** letras minúsculas, dígitos e hífens; deve começar com letra (`acme`, `acme-corp`). Não use `_` nem maiúsculas.

### Passo 2 — Editar cores

Abra `frontend/src/skins/acme.css`. Estrutura mínima:

```css
[data-skin='acme'] {
  --color-accent: #e65100;
  --color-accent-2: #ff9800;
  --color-accent-bar: #bf360c;
  /* … demais tokens — copie de _template.css … */

  --skin-login-bg: #37474f;
  --skin-login-card: #e65100;
  /* … tokens de login … */
}

[data-skin='acme'].dark {
  --color-panel: #263238;
  /* … overrides escuros … */
}
```

**Dica:** copie `snappymail.css` ou `gmail.css` como ponto de partida e ajuste hex values.

### Passo 3 — Aliases (opcional)

Se quiser que `acme-brand` no `config.toml` resolva para `acme`:

**Go** (`internal/ui/skins.go`):

```go
Aliases: []string{"acme-brand", "acme-mail"},
```

**TS** (`frontend/src/skins/manifest.ts`):

```ts
aliases: ['acme-brand', 'acme-mail'],
```

Aliases devem ser **idênticos** nos dois arquivos.

### Passo 4 — Ativar e testar

```toml
[ui]
skin = "acme"
```

```bash
# terminal 1
make frontend-dev

# terminal 2
make run
```

Abra `http://localhost:8082` (ou `:5173` com proxy Vite).

**Verificar:**

- [ ] Tela de login com cores Acme (fundo, card, botão)
- [ ] Inbox após login (sidebar, lista, toolbar)
- [ ] Toggle dark mode na toolbar — login + inbox
- [ ] `make validate-skins` passa

### Passo 5 — Marcar como pronta (quando aplicável)

Enquanto a skin só altera cores (layout padrão 3 colunas), pode marcar `ready: true`:

```go
// internal/ui/skins.go
Ready: true,
```

```ts
// frontend/src/skins/manifest.ts
ready: true,
```

Skins Gmail/Outlook ficam `ready: false` até terem layout próprio.

### Passo 6 — Build de produção

```bash
make validate-skins
make test
make build-prod
```

---

## Registro manual (passo a passo)

Use quando `REGISTER=1` não for opção ou precisar entender cada arquivo.

### 1. Criar CSS

```bash
make new-skin ID=acme
# ou: cp frontend/src/skins/_template.css frontend/src/skins/acme.css
#     e substituir __SKIN_ID__ por acme
```

### 2. Go — `internal/ui/skins.go`

Insira **antes** de `} // catalog-end`:

```go
	{
		ID:      "acme",
		Label:   "Acme Corp",
		Ready:   false,
		Aliases: []string{},
	},
```

Campos:

| Campo | Descrição |
|-------|-----------|
| `ID` | Id canônico (igual ao nome do `.css`) |
| `Label` | Nome exibido na UI / API |
| `Ready` | Skin completa para produção |
| `Aliases` | Valores alternativos aceitos em `config.toml` |

### 3. TypeScript — `frontend/src/skins/manifest.ts`

Insira **antes** de `] as const // manifest-end`:

```ts
  {
    id: 'acme',
    label: 'Acme Corp',
    ready: false,
    aliases: [],
  },
```

### 4. Import CSS — `frontend/src/skins/index.css`

Adicione **antes** do comentário `imports-begin`:

```css
@import "./acme.css";
```

### 5. Validar

```bash
make validate-skins
go test ./internal/ui/...
cd frontend && npm run build
```

**Não é necessário** editar `style.css` (importa `./skins/index.css`), `registry.ts` (re-exporta manifest), nem `types.ts` separado.

---

## Arquivos do sistema

```
config.toml                    [ui] skin
internal/ui/skins.go           catálogo Go + NormalizeSkin()
internal/handler/ui.go         GET /api/v1/ui/config
frontend/src/skins/
  _template.css                scaffold para novas skins
  acme.css                     tokens da skin acme
  index.css                    @import de todas as *.css
  manifest.ts                  ids, labels, aliases, ready
  registry.ts                  re-export (use manifest.ts)
  bootstrap.ts                 fetch /ui/config antes do mount
  apply.ts                     document.documentElement data-skin
frontend/src/style.css         @import skins/index.css + classes login
frontend/src/stores/settings.ts  Pinia — skin, dark mode
frontend/src/components/LoginView.vue
scripts/new-skin.sh            scaffold + REGISTER=1
scripts/validate-skins.sh        sync Go ↔ TS ↔ CSS
```

Marcadores para auto-registro (`make new-skin REGISTER=1`):

| Marcador | Arquivo |
|----------|---------|
| `catalog-begin` … `catalog-end` | `internal/ui/skins.go` |
| `manifest-begin` … `manifest-end` | `frontend/src/skins/manifest.ts` |
| `imports-begin` | `frontend/src/skins/index.css` |

---

## Variáveis CSS — referência completa

Cada skin define tokens em `[data-skin='seu-id']`. O SPA seta `data-skin` em `<html>` na inicialização.

### Inbox / componentes compartilhados

| Variável | Uso |
|----------|-----|
| `--color-accent` | Destaque principal, badges |
| `--color-accent-2` | Hover, contadores unread |
| `--color-accent-bar` | Bordas de destaque, header login |
| `--color-line` | Bordas (`border-line`) |
| `--color-panel` | Painéis principais (`bg-panel`) |
| `--color-panel-2` | Hover, fundos secundários |
| `--color-app-bg` | Fundo da página (`bg-app-bg`) |
| `--color-ink` | Texto principal (`text-ink`) |
| `--color-ink-sub` | Texto secundário |
| `--color-ink-mute` | Texto muted, labels |
| `--color-row-selected` | Linha selecionada na lista |
| `--font-sans` | Font stack global |

### Login (`LoginView.vue`)

Classes em `style.css` leem **somente** estes tokens — nunca condicionais por skin no Vue.

| Variável | Elemento |
|----------|----------|
| `--skin-login-bg` | `.login-page` |
| `--skin-login-card` | `.login-card` |
| `--skin-login-header-bg` | `.login-header` |
| `--skin-login-header-border` | `.login-header` border |
| `--skin-login-text` | Texto do formulário |
| `--skin-login-input-bg` | `.login-input` |
| `--skin-login-input-border` | `.login-input` border |
| `--skin-login-input-text` | Texto digitado |
| `--skin-login-btn-bg` | `.login-btn` |
| `--skin-login-btn-text` | Texto do botão |
| `--skin-login-error-bg` | `.login-error` |
| `--skin-login-error-text` | Texto de erro |
| `--skin-login-error-border` | Borda de erro |

### Opcional

| Variável | Uso futuro |
|----------|------------|
| `--skin-layout` | Hint para AppShell (`three-column`, `gmail`, …) |

### Exemplo comparativo

**SnappyMail** — login escuro, card azul:

```css
--skin-login-bg: #48525c;
--skin-login-card: #1b3a6b;
--skin-login-text: #ffffff;
```

**Gmail** — login claro, card branco:

```css
--skin-login-bg: #f6f8fc;
--skin-login-card: #ffffff;
--skin-login-text: var(--color-ink);
--skin-login-btn-bg: var(--color-accent);
--skin-login-btn-text: #ffffff;
```

Consulte `frontend/src/skins/snappymail.css`, `gmail.css` e `outlook.css` como referência.

---

## Como o Vue consome tokens

### Tailwind v4 + `@theme`

`frontend/src/style.css` declara defaults em `@theme { --color-accent: … }`. A skin ativa **sobrescreve** essas variáveis via `[data-skin='…']`.

Classes usadas nos componentes:

| Classe Tailwind | Variável |
|-----------------|----------|
| `bg-panel` | `--color-panel` |
| `bg-app-bg` | `--color-app-bg` |
| `border-line` | `--color-line` |
| `text-ink` | `--color-ink` |
| `text-ink-sub` | `--color-ink-sub` |
| `text-ink-mute` | `--color-ink-mute` |
| `text-accent-2` | `--color-accent-2` |

Classes globais em `style.css` (não Tailwind):

| Classe | Uso |
|--------|-----|
| `.side-item` / `.side-item.active` | Sidebar de pastas |
| `.msg-row` / `.msg-row.selected` | Lista de mensagens |
| `.tbtn` | Botões da toolbar |
| `.login-*` | Tela de login |

### Fluxo no código

```ts
// main.ts
const { config } = await bootstrapUI()
useSettingsStore(pinia).initFromServer(config)

// apply.ts
document.documentElement.setAttribute('data-skin', skinId)
```

---

## Modo escuro

- Toggle na toolbar → `settings.darkMode` → classe `.dark` em `<html>`.
- Persistência: `localStorage` chave `gsn_dark` (`1` = escuro).
- **Cada skin** deve ter bloco `[data-skin='id'].dark { … }`.

Exemplo mínimo (dark):

```css
[data-skin='acme'].dark {
  --color-panel: #263238;
  --color-app-bg: #1c2529;
  --color-ink: #eceff1;
  --color-line: #37474f;
  --skin-login-bg: #1c2529;
  --skin-login-card: #263238;
}
```

Tailwind usa `@custom-variant dark (&:where(.dark, .dark *));` — classes `dark:` funcionam dentro de qualquer skin.

---

## Flag `ready` e banner de preview

Quando `ready: false` (Go + manifest), `App.vue` exibe:

> Skin preview: Gmail — full layout coming soon

Isso evita confundir tokens-only com produto final. Atualize **nos dois lugares**:

```go
Ready: true,  // skins.go
```

```ts
ready: true,  // manifest.ts
```

---

## Aliases em `config.toml`

Permite valores legados ou nomes comerciais:

```toml
[ui]
skin = "google"    # → gmail
skin = "microsoft" # → outlook
skin = "default"   # → snappymail
```

Implementação:

- **Servidor:** `ui.NormalizeSkin()` em `skins.go`
- **Cliente:** `normalizeSkinId()` em `manifest.ts`

Valores desconhecidos caem no default `snappymail`.

---

## Validação (`make validate-skins`)

```bash
make validate-skins
```

Verifica:

1. IDs iguais em Go catalog, TS manifest e `index.css` imports
2. Arquivo `frontend/src/skins/<id>.css` existe para cada id
3. Selector `[data-skin='<id>']` presente
4. Bloco `[data-skin='<id>'].dark` presente

Saída esperada:

```
Go catalog:  gmail outlook snappymail
TS manifest: gmail outlook snappymail
CSS imports: gmail outlook snappymail
OK — skin catalog, manifest, and CSS are in sync.
```

Rode sempre após adicionar, renomear ou remover uma skin.

---

## Erros comuns e soluções

| Problema | Causa | Solução |
|----------|-------|---------|
| Skin não muda no browser | Servidor não reiniciado | Reinicie `make run` após editar `config.toml` |
| Cores default (azul SnappyMail) | Id desconhecido ou typo | Confira spelling; aliases em `skins.go` |
| Login ok, inbox errado | Tokens inbox faltando | Preencha `--color-panel`, `--color-ink`, etc. |
| Dark mode quebrado | Sem bloco `.dark` | Adicione `[data-skin='id'].dark` |
| `validate-skins` falha | Go/TS/CSS dessincronizados | Compare ids; use `REGISTER=1` ou registro manual |
| Build TS falha após novo id | Manifest sem entrada | Adicione em `manifest.ts` entre marcadores |
| Banner "Skin preview" | `ready: false` | Normal para gmail/outlook; set `true` quando pronto |
| CSS não carrega | Falta `@import` | Adicione em `frontend/src/skins/index.css` |

---

## Regras para componentes Vue

### Faça

- Use classes Tailwind mapeadas a tokens (`bg-panel`, `text-ink-mute`).
- Use classes globais `.login-*` na tela de login.
- Novas cores → nova variável CSS na skin + `@theme` default em `style.css`.

### Não faça

- `v-if="settings.skin === 'gmail'"` para estilizar — use tokens.
- Hex hardcoded em `.vue` (`#1b3a6b`) — quebra outras skins.
- Editar `style.css` para uma skin específica — use `frontend/src/skins/<id>.css`.
- Esquecer bloco `.dark` ao adicionar variável de cor.

---

## Layout customizado (futuro)

Hoje: **desktop 3 colunas fixas** em `App.vue` — mobile/responsive **adiado** para quando skins tiverem layout próprio.

Quando um produto exigir estrutura diferente (ex.: Gmail):

1. Crie `frontend/src/skins/<id>/AppShell.vue`
2. Em `App.vue`, troque layout via `settings.skin` ou `defineAsyncComponent`
3. Defina `--skin-layout: gmail` no CSS
4. Marque `ready: true` só com shell completo (+ mobile, se aplicável)

Até lá: skins são **CSS-only** sobre o layout SnappyMail padrão.

---

## API

`GET /api/v1/ui/config` (público, sem auth):

```json
{
  "skin": "snappymail",
  "skins": [
    {"id": "snappymail", "label": "SnappyMail", "ready": true},
    {"id": "gmail", "label": "Gmail", "ready": false},
    {"id": "outlook", "label": "Outlook", "ready": false}
  ],
  "available_skins": ["snappymail", "gmail", "outlook"],
  "rows_per_page": 50,
  "datetime_format": "02/01/2006 15:04",
  "compose_html": true
}
```

Detalhes: [api.md](api.md#get-apiv1uiconfig).

---

## Checklist final

Copie ao implementar uma skin:

```
Skin id: ___________

[ ] make new-skin ID=___ REGISTER=1  (ou registro manual completo)
[ ] frontend/src/skins/_____.css — tokens light
[ ] frontend/src/skins/_____.css — bloco .dark
[ ] Aliases iguais em skins.go e manifest.ts (se houver)
[ ] make validate-skins — OK
[ ] config.toml skin = "___" testado
[ ] Login — light + dark
[ ] Inbox — light + dark
[ ] ready: true/false correto em Go + manifest
[ ] go test ./...
[ ] cd frontend && npm run build
```

Referência rápida no código: [frontend/src/skins/README.md](../frontend/src/skins/README.md).
