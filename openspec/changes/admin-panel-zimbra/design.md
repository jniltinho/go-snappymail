# Design: admin-panel-zimbra

## Context

Este projeto (**go-snappymail**) é o **binário único** do ecossistema "Zimbra em Go": já serve o webmail (Zimbra Web Classic clone) com Go + Echo v5 + Vue 3 + Tailwind + GORM embutidos (`main.go` faz `//go:embed all:web/dist all:web/files`; rotas em `/api/v1`; sessão via cookie IMAP). Adicionamos aqui o **painel admin** (clone do ZimbraAdmin, `:7071`).

**go-postfixadmin = referência apenas.** Não é alterado nem importado; serve de exemplo do backend (schema Postfix/Dovecot, handlers `*V1`, GORM, RBAC granular). Portamos os padrões para o go-snappymail.

## Goals / Non-Goals

**Goals**
- Painel ZimbraAdmin em Vue 3 + TailwindCSS, `:7071`, no mesmo binário do webmail.
- Backend admin novo aqui (GORM → MariaDB/PostgreSQL), inspirado no go-postfixadmin.
- **Dois frontends embutidos** (`go:embed`) num único binário; instalação simples.

**Non-Goals** — ver proposal (não tocar go-postfixadmin; sem DWT/logo/monitor/zimlets nesta fase; não reusar frontend existente).

## Decisions

### D1 — Host = go-snappymail (sem problema de import entre módulos)
Como tudo é feito **neste** projeto, não há import de pacotes `internal/` entre módulos (o problema do plano anterior desaparece). O webmail já está aqui; só **adicionamos** o serviço admin ao mesmo binário. go-postfixadmin é lido como referência e reimplementado/adaptado — nada de `go.mod require` dele.

### D2 — Dois frontends, ambos `go:embed` (binário único)
Os **dois** frontends são embutidos no binário:

```
main.go
//go:embed all:web/dist        → SPA do webmail   (frontend/)
//go:embed all:web/admin-dist  → SPA do admin     (frontend-admin/)   ← NOVO
//go:embed all:web/files
```

- Webmail: `fs.Sub(embed, "web/dist")` servido em `/` na porta `:8082` (como hoje).
- Admin: `fs.Sub(embed, "web/admin-dist")` servido em base `/admin/` (ou `/zimbraAdmin/`) na porta `:7071`.
- Nenhum asset em runtime fora do binário; `go build` embute os dois builds. `make build` roda os dois frontends antes.

### D3 — Frontend admin DO ZERO (Vue 3 + TailwindCSS), pasta `frontend-admin/`
Criado do zero, sem reaproveitar o webmail nem qualquer UI existente. **Cada frontend tem seu próprio tema e layout, totalmente separados:**

| | Webmail (`frontend/`) | Admin (`frontend-admin/`) |
|---|---|---|
| Layout | Zimbra Web Classic (3 colunas: pastas/lista/leitura) | ZimbraAdmin Classic (top bar + árvore + content pane) |
| Tema | skins do webmail (`src/skins/`, incl. `zimbra`) | tema admin próprio (`src/theme/`, ZimbraAdmin) |
| Tailwind config | do webmail | **próprio** (`frontend-admin/tailwind.config.ts`) |
| Tokens | `@theme` do webmail | `@theme` do admin (mesma família harmony, mas arquivo separado) |
| Build → embed | `web/dist` | `web/admin-dist` |

Nada de tema/layout/config compartilhado entre os dois: cada pasta tem seu `tailwind.config`, seu `@theme`, seu shell/layout e seus componentes. O stack é o mesmo (Vue 3 + Tailwind) só por familiaridade, **não** por compartilhamento — um não importa nada do outro. Se um dia surgir token realmente comum, extrai-se um pacote; por ora, isolamento total.

```
frontend/                 # webmail (existente) → web/dist            (embed 1)
frontend-admin/           # NOVO painel ZimbraAdmin → web/admin-dist  (embed 2)
  src/
    App.vue  main.ts
    layouts/    AdminShell.vue (top bar + nav tree + content pane)
    views/      Login, Home, Domains, Accounts, Aliases, Admins
    components/ NavTree, Toolbar, ListView, Toast
    theme/      tailwind + tokens ZimbraAdmin (harmony, 3px, tipografia)
    stores/  router/  api/
  index.html  package.json  tsconfig.json
  tailwind.config.ts  postcss.config.js
  vite.config.ts        # base: '/admin/'  ·  outDir: '../web/admin-dist'
```
Tailwind utility-first para estrutura + camada de tokens ZimbraAdmin (`@theme`/CSS vars) para paleta harmony, cantos 3px, tipografia Helvetica/Arial (reaproveitar valores medidos no webmail).

**Organização Vue (nos dois frontends):** **um modal por arquivo** em `components/modals/` (`DomainModal.vue`, `MailboxModal.vue`, …) — nunca vários modais no mesmo arquivo. Estrutura consistente: `layouts/` (shell), `views/` (uma tela por arquivo), `components/` (reuso), `components/modals/` (modais), `stores/`, `api/`, `theme/`. O webmail (`frontend/`) segue o mesmo padrão (organizar modais existentes, ex. composer, em `components/modals/`).

### D4 — Backend admin novo aqui (telas ↔ endpoints), go-postfixadmin como referência
Novos modelos GORM do schema Postfix/Dovecot + handlers no go-snappymail.

| Tela | Endpoint (novo, go-snappymail) |
|---|---|
| Home overview | `GET /api/v1/admin/overview` (counts reais do banco) |
| Domains | `/api/v1/admin/domains*` |
| Accounts (mailboxes) | `/api/v1/admin/mailboxes*` |
| Aliases | `/api/v1/admin/aliases*` |
| Admins | `/api/v1/admin/admins*` |
| Distribution Lists / Resources / COS | stub nesta fase (schema/RBAC não cobre) |

**Overview realista** — só o que o schema tem: accounts, domains, aliases, admins. Version/servers/queue/sessions = `—`/"n/a", nunca inventados.

Persistência: **banco separado** do mail (MariaDB/PostgreSQL via GORM), distinto do banco de sessão do webmail (SQLite/GORM). Config `[database]` do mail aparte. Unit tests com sqlmock; matrix MariaDB+Postgres em container.

### D5 — Auth admin (JWT + RBAC), independente do webmail
O webmail usa cookie de sessão IMAP; o admin usa **JWT próprio + RBAC granular** (superadmin/domain_admin + permissões `domains:read`, `mailboxes:write`, …) contra a tabela de admins do banco. Cada nó/rota exige a permissão; sem ela → 403; nós fora da permissão ficam ocultos. **Isolamento de cookies/CSRF:** como admin e webmail podem coexistir no mesmo host, usar nomes/Path/SameSite/CSRF distintos por serviço (ex.: `gsn_session` webmail × `gsn_admin_jwt` admin).

### D6 — Multi-listener num processo
`serve` sobe dois listeners Echo: `:8082` (webmail, atual) e `:7071` (admin), cada um com seu middleware, seu `fs.Sub` de SPA e **graceful shutdown** coordenado (errgroup + signal). Bloco `[admin] enabled=false` → listener não abre.

### D7 — Separação de arquivos: rotas e templates (webmail × admin)
Rotas e templates/SPA de cada serviço ficam em **arquivos próprios**, nunca misturados:

```
internal/server/
  routes.go            # rotas do WEBMAIL (/api/v1/*, SPA web/dist) — existente
  routes_admin.go      # rotas do ADMIN  (/api/v1/admin/*, SPA web/admin-dist) — NOVO
  render.go            # render/SPA-fallback do WEBMAIL (serve web/dist) — existente
  render_admin.go      # render/SPA-fallback do ADMIN  (serve web/admin-dist) — NOVO
  server.go            # sobe os dois listeners (webmail :8082, admin :7071)
internal/handler/      # handlers do webmail (existente)
internal/admin/        # handlers + modelos + auth do ADMIN (NOVO, isolado)
  routes.go  handlers.go  models.go  auth.go  overview.go

web/
  dist/                # template/SPA do WEBMAIL   (go:embed)
  admin-dist/          # template/SPA do ADMIN     (go:embed)   ← separado
```

- **Rotas** em arquivos distintos: `registerAPIRoutes` (webmail, `routes.go`) e `registerAdminRoutes` (admin, `routes_admin.go`); cada um monta seu `Group` e seu middleware. Nada de rota de admin no arquivo do webmail e vice-versa.
- **Render** em arquivos distintos: o SPA-fallback/serving de cada serviço tem seu próprio arquivo (`render.go` × `render_admin.go`), servindo **só** o seu `fs.Sub` (`web/dist` × `web/admin-dist`). Sem fallback cruzado (uma rota do admin nunca cai na SPA do webmail).
- **Templates/SPA**: `web/dist` (webmail) e `web/admin-dist` (admin) são embeds separados.
- **Config**: blocos `[webmail]` e `[admin]` separados; nada compartilhado além de `[database]` do mail (só o admin usa).

## Melhorias de backend (skills `golang-*`)

Estender o backend do go-snappymail com as skills onde precisar: `golang-security`/`golang-safety` (rotas, cookies/CSRF, headers), `golang-database` (agregados, transações, índices GORM), `golang-concurrency` (multi-listener/errgroup/shutdown), `golang-error-handling`/`golang-observability` (`%w`, slog), `golang-testing` (table-driven, `-race`, sqlmock), `golang-lint` (golangci-lint), `golang-modernize`. Cada melhoria com teste que a cobre.

## Infra a criar

- `.golangci.yml` + rodar `golangci-lint` no CI.
- `frontend-admin/` com ESLint + Prettier + `vue-tsc`.
- CI test+lint (backend e frontend) a cada push.
- Ambiente de banco: MariaDB (Docker `:3306`) e PostgreSQL para a matrix.

## Risks / Trade-offs

- [Portar o schema/handlers do go-postfixadmin] → adaptar por tela, com teste; go-postfixadmin fica de referência viva (roda no lab).
- [Paridade visual do ZimbraAdmin] → fatiar por tela (Login → Home → Domains → …), validar cada uma com `qa-frontend-cloner` contra `:7071`, prints em `docs/prints/`.
- [Dois serviços no mesmo host] → cookies/CSRF isolados por serviço (D5); portas distintas.
- [Binário maior com 2 SPAs] → aceitável; instalação continua "um binário + um config".

## Migration Plan

Fatiado: (1) config `[admin]`/`[database]` + multi-listener servindo admin vazio; (2) scaffold `frontend-admin/` (Vue+Tailwind) + embed 2; (3) login + top bar + árvore + Home; (4) modelos GORM + `/api/v1/admin/*` + overview; (5) telas Domains/Accounts/Aliases/Admins; (6) validação QA por tela. Rollback = `[admin] enabled=false`.

## Open Questions

- Base path do admin: `/zimbraAdmin/` (espelha o legado) vs `/admin/`.
- Schema Postfix/Dovecot: reusar o do go-postfixadmin como referência 1:1 ou enxugar ao necessário do painel.
