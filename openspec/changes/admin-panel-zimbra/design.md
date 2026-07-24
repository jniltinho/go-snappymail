# Design: admin-panel-zimbra

## Context

Este projeto (**go-snappymail**) é o **binário único** do ecossistema "Zimbra em Go": já serve o webmail (Zimbra Web Classic clone) com Go + Echo v5 + Vue 3 + Tailwind + GORM embutidos (`main.go` faz `//go:embed all:web/dist all:web/files`; rotas em `/api/v1`; sessão via cookie IMAP). Adicionamos aqui o **painel admin** (clone do ZimbraAdmin, `:7071`).

**go-postfixadmin = referência apenas.** Não é alterado nem importado; serve de exemplo do backend (schema Postfix/Dovecot, handlers `*V1`, GORM, RBAC granular). Portamos os padrões para o go-snappymail.

## Ordem de trabalho (regra do projeto)

**Sempre começar pelo backend** (API Go + testes) antes da UI — regra 1 do AGENTS.md. Sequência: modelos GORM + `/api/v1/admin/*` + auth/RBAC + testes → só então o `frontend-admin/`. Antes da UI, capturar **todos os prints do ZimbraAdmin** e os **tokens de tema** (já extraídos em [`zimbra-admin-theme.md`](zimbra-admin-theme.md), do skin `serenity` real).

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

**Antes de codar a UI:** capturar **todos os prints do ZimbraAdmin** (`192.168.56.30:7071`) em `docs/prints/zimbra-admin/` — login, Home, cada nó da árvore, telas de Domains/Accounts/Aliases/Admins com seus modais, toolbars e toasts — e extrair deles os tokens (paleta harmony, tipografia Helvetica/Arial, cantos 3px). O tema Tailwind do admin é derivado desses prints para ficar **idêntico** ao legado. Tirar novos prints sempre que precisar comparar (convenção do projeto: prints só em `docs/prints/`). **Toda a UI em inglês** (o console legado está em pt-BR por locale da VM; o clone usa inglês).

**Organização Vue (nos dois frontends):** **um modal por arquivo** em `components/modals/` (`DomainModal.vue`, `MailboxModal.vue`, …) — nunca vários modais no mesmo arquivo. Estrutura consistente: `layouts/` (shell), `views/` (uma tela por arquivo), `components/` (reuso), `components/modals/` (modais), `stores/`, `api/`, `theme/`. O webmail (`frontend/`) segue o mesmo padrão (organizar modais existentes, ex. composer, em `components/modals/`).

### D3b — Multi-skin nos DOIS frontends (deixar pronto para trocar de skin)
Ambos os frontends são **multi-skin**, cada um com seu próprio catálogo (não compartilham skins entre si):
- **Webmail** (`frontend/`): já é multi-skin — `[ui] skin` = `zimbra` (default) | `snappymail` | `gmail` | `outlook` | `carbonio`, com o catálogo Go (`internal/ui/skins.go`) + manifest TS + CSS por skin e `make validate-skins`. **Manter e deixar extensível** para novas skins.
- **Admin** (`frontend-admin/`): **espelhar a mesma arquitetura** — catálogo de skins do admin, tokens por skin em `[data-skin='<id>']`/`@theme`, flag `ready`, config `[admin] skin`. Skins iniciais extraídas do ZimbraAdmin real: **`serenity`** (default) e **`carbon`** (2º), com espaço para `vami2`/outras. Tokens em [`zimbra-admin-theme.md`](zimbra-admin-theme.md).
  - Mesma disciplina do webmail: um `validate-skins` do admin (catálogo ↔ manifest ↔ CSS) e dark mode client-side por skin quando fizer sentido.
- **Cada frontend tem seu conjunto de skins** — o webmail não usa skins do admin e vice-versa; cada catálogo vive na sua pasta.

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

### D6 — Multi-listener num processo, com REGRA DE OURO de isolamento
`serve` sobe dois listeners Echo **independentes**: `:8082` (webmail, atual) e `:7071` (admin), cada um com seu próprio `*echo.Echo`, seu middleware, seu `fs.Sub` de SPA e **graceful shutdown** coordenado (errgroup + signal). Bloco `[admin] enabled=false` → listener não abre.

**Regra de ouro (isolamento de superfície — inviolável):** os endpoints de admin **nunca** podem ser alcançados pela porta/listener do webmail, e vice-versa.
- **Instâncias Echo separadas**, não um Echo com dois grupos: as rotas `/api/v1/admin/*` e a SPA `web/admin-dist` são registradas **apenas** no Echo do `:7071`; as rotas do webmail e `web/dist` **apenas** no Echo do `:8082`. Não existe um único router que conheça as duas árvores.
- `registerAdminRoutes` é chamado **só** para o Echo admin; `registerAPIRoutes` (webmail) **só** para o Echo webmail. Um não referencia o outro.
- **Bind opcional em interface distinta:** o admin pode escutar em `127.0.0.1:7071` (só local) enquanto o webmail escuta em `0.0.0.0:8082`, reduzindo exposição — configurável por `[admin] host`.
- **Sem proxy/rewrite** que reencaminhe `/admin` a partir da porta do webmail. Bater em `http://host:8082/api/v1/admin/...` retorna 404 (a rota não existe naquele router), nunca alcança o handler admin.
- Teste obrigatório: requisição a rota de admin **na porta do webmail → 404**; e rota do webmail **na porta do admin → 404** (ver tasks §5).

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
