# Tasks: admin-panel-zimbra

## 1. Binário único / config multi-serviço (no go-snappymail)

- [x] 1.1 Blocos `[webmail]`, `[admin]`, `[database]` (mail) no `config.toml` (host, porta, tls, enabled, driver/dsn) + env `GOSM_*`; `[admin] host` permite bind restrito (ex.: `127.0.0.1`)
- [ ] 1.2 `serve` sobe dois listeners Echo: webmail `:8082` (atual) + admin `:7071` — cada `enabled=true`
- [ ] 1.3 Multi-listener com graceful shutdown coordenado (errgroup + signal); `[admin] enabled=false` não abre o listener
- [ ] 1.4 Verificar isolamento: admin `:7071` e webmail `:8082` sem colisão; cookies/CSRF por serviço

## 1b. Separação de arquivos (backend): rotas, render, templates

- [ ] 1b.1 `internal/server/routes_admin.go` (admin) separado de `routes.go` (webmail); `registerAdminRoutes` próprio
- [ ] 1b.2 `internal/server/render_admin.go` servindo **só** `web/admin-dist` (sem fallback cruzado com `web/dist`)
- [ ] 1b.3 `internal/admin/` isolado (handlers/models/auth/overview) — nada de handler admin no pacote do webmail
- [ ] 1b.4 `main.go` `//go:embed` inclui `web/admin-dist` além de `web/dist`; ambos embutidos

## 2. Backend admin (GORM → MariaDB/PostgreSQL) — go-postfixadmin como referência

- [x] 2.1 Modelos GORM do schema Postfix/Dovecot (Domain/Mailbox/Alias/Admin), banco separado do de sessão
- [ ] 2.2 `GET /api/v1/admin/overview` (counts reais accounts/domains/aliases/admins; version/servers/queue/sessions = n/a), permissão admin
- [ ] 2.3 Handlers `/api/v1/admin/domains|mailboxes|aliases|admins` (CRUD) — portados/adaptados do go-postfixadmin
- [x] 2.4 Auth JWT + RBAC granular (superadmin/domain_admin + permissões); 403 sem permissão
- [ ] 2.5 Migração/seed de dev (MariaDB Docker `:3306` ou PostgreSQL; sqlmock para unit)

## 3. Frontend admin DO ZERO (Vue 3 + TailwindCSS) — pasta `frontend-admin/`

- [ ] 3.-1 **Capturar TODOS os prints do ZimbraAdmin** (`192.168.56.30:7071`, login `admin@zimbra.test`/`Password1@`) ANTES de codar a UI, salvando em `docs/prints/zimbra-admin/`: login, Home, cada nó da árvore, Domains (list+modal New/Edit), Accounts, Aliases, Admins, toolbars, toasts, dark se houver. Extrair tokens (paleta/tipografia/cantos) desses prints — o tema tem que ficar **igual**. Tirar novos prints sempre que precisar comparar.
- [ ] 3.0 Scaffold `frontend-admin/` novo (Vue 3 + TS + Vite + **TailwindCSS** + PostCSS): `base:'/admin/'`, `outDir:'../web/admin-dist'`, package.json/tsconfig/tailwind.config próprios; adicionar `web/admin-dist` ao `//go:embed` e `make frontend-admin` (não usar nada do `frontend/` antigo). **Toda a UI em inglês.**
- [ ] 3.1 Tema ZimbraAdmin em Tailwind: tokens harmony (paleta, cantos 3px, tipografia Helvetica/Arial) via `@theme`/CSS vars
- [ ] 3.2 Top bar (marca textual, busca, `admin@… ▾`, refresh)
- [ ] 3.3 Árvore de navegação (Home/Monitor/Manage/Configure/Tools/Search)
- [ ] 3.4 Home (Overview + cards de setup) ligada ao `/api/v1/admin/overview`
- [ ] 3.5 List view + toolbar (colunas, seleção com outline, paginação)
- [ ] 3.6 **Multi-skin do admin** (espelhar o webmail): catálogo de skins do admin (`serenity` default + `carbon`), tokens por `[data-skin='<id>']`/`@theme`, `[admin] skin` na config, `validate-skins` do admin; cada frontend com seu próprio conjunto de skins

## 3b. Organização do Vue — um modal por arquivo (nos DOIS frontends)

- [ ] 3b.1 `frontend-admin/src/components/modals/` — **um arquivo por modal** (`DomainModal.vue`, `MailboxModal.vue`, `AliasModal.vue`, `AdminModal.vue`, `ConfirmModal.vue`), nunca vários modais num só arquivo
- [ ] 3b.2 `frontend/src/components/modals/` (webmail) — mover/organizar os modais existentes (ex.: composer) para um arquivo por modal, mesmo padrão
- [ ] 3b.3 Estrutura consistente nos dois: `layouts/` (shell), `views/` (telas), `components/` (reuso), `components/modals/` (modais), `stores/`, `api/`, `theme/`
- [ ] 3b.4 **Tema/layout separados por frontend**: cada um com seu `tailwind.config`, seu `@theme`/tokens, seu shell/layout — nada compartilhado; um não importa do outro
- [ ] 3b.5 **Componente base `Modal.vue`** (overlay, foco, Esc/click-outside, header/body/footer) + store de UI para modal aberto; admin é **modal-heavy** (toda ação em diálogo, não página)
- [ ] 3b.6 **Cliente API JSON tipado** (`api/`): request/response JSON com envelope consistente; modais enviam e populam JSON (sem form-encoded)

## 4. Telas de gerenciamento (uma view + um modal por recurso)

- [ ] 4.1 Domains — `views/Domains.vue` (list) + `components/modals/DomainModal.vue` (New/Edit) + Delete
- [ ] 4.2 Accounts/Mailboxes — `views/Accounts.vue` + `components/modals/MailboxModal.vue`
- [ ] 4.3 Aliases — `views/Aliases.vue` + `components/modals/AliasModal.vue` (Distribution Lists: stub)
- [ ] 4.4 Admins — `views/Admins.vue` + `components/modals/AdminModal.vue` (papéis RBAC)

## 5. Testes de backend — COBRIR TODAS as rotas admin (obrigatório, backend-first)

> Meta: **nenhum endpoint admin sem teste**. Cada rota `/api/v1/admin/*` tem testes de
> sucesso, validação, erro e permissão. `go test -race ./...` verde; cobertura reportada.

- [ ] 5.1 `GET /api/v1/admin/overview` (agregados corretos; permissão exigida; 403 sem permissão) — table-driven, `-race`
- [ ] 5.2 **Domains** — GET list, GET by id, POST create, PUT update, DELETE: sucesso + validação + duplicado + not-found + erro de banco
- [ ] 5.3 **Mailboxes/Accounts** — list/get/create/update/delete + validação de e-mail/senha/quota + domínio inexistente
- [ ] 5.4 **Aliases** — list/get/create/update/delete + destino inválido + loop/duplicado
- [ ] 5.5 **Admins** — list/get/create/update/delete + atribuição de papel RBAC
- [ ] 5.6 **Auth/RBAC** — JWT válido/expirado/ausente → 401; sem permissão → 403; superadmin × domain_admin (escopo por domínio)
- [ ] 5.7 **Isolamento de superfície** — rota `/api/v1/admin/*` na porta do webmail (:8082) → 404; rota do webmail na porta do admin (:7071) → 404
- [ ] 5.8 **Config multi-serviço** — blocos enabled/disabled → portas/bind certos; `[admin] host` restrito; parsing TOML/env
- [ ] 5.9 **Persistência GORM** — MariaDB **e** PostgreSQL (matrix; sqlmock para unit); migrações aplicam limpo; transações/rollback
- [ ] 5.10 Cobertura mínima acordada por pacote admin e CI verde bloqueando merge
- [ ] 5.4 Testes de config multi-serviço (blocos enabled/disabled → portas certas; parsing TOML/env)
- [ ] 5.5 Testes de persistência GORM em MariaDB **e** PostgreSQL (matrix; sqlmock para unit — ConnectDB não tem SQLite), migrações aplicam limpo
- [ ] 5.6 Cobertura mínima acordada e `go test -race ./...` verde no CI

## 5b. Melhorias de backend (skills golang-*, estender não reescrever)

- [ ] 5b.1 `golang-security`/`golang-safety`: auditar rotas do painel, cookies/CSRF por serviço (D5), headers
- [ ] 5b.2 `golang-database`: query de agregados do overview, transações, NULLs, índices GORM
- [ ] 5b.3 `golang-concurrency`: multi-listener com errgroup + graceful shutdown coordenado
- [ ] 5b.4 `golang-error-handling`/`golang-observability`: `%w`, slog estruturado, logs/metrics das rotas admin
- [ ] 5b.5 `golang-modernize`/`golang-structs-interfaces`: refino idiomático onde tocar; cada melhoria com teste

## 6. Lint / qualidade (obrigatório)

- [ ] 6.1 `golangci-lint run` limpo no backend (govet, staticcheck, errcheck, gosec, revive…); config `.golangci.yml`
- [ ] 6.2 `go vet ./...` e `gofmt`/`goimports` sem diffs
- [ ] 6.3 Lint do frontend admin: `eslint` + `vue-tsc --noEmit` (typecheck) sem erros; `prettier` consistente
- [ ] 6.4 CI roda test + lint (backend e frontend) a cada push; falha bloqueia merge

## 7. Validação e docs

- [ ] 7.1 Auditar cada tela com `qa-frontend-cloner` contra `:7071` até 0 P1 (prints em `docs/prints/`)
- [ ] 7.2 Doc de arquitetura do binário único + mapeamento telas↔endpoints
- [ ] 7.3 Atualizar guia de dev (subir binário único: webmail + admin + Zimbra por snapshot)
