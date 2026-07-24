# Proposal: admin-panel-zimbra

## Why

O ecossistema "Zimbra em Go" recria o Zimbra em Go + Vue 3, leve e rápido, num **único binário para instalação simples**. Este projeto (**go-snappymail**) já é esse binário: serve o webmail (clone do Zimbra Web Classic, IMAP/SMTP passthrough) com Go + Echo v5 + Vue 3 + Tailwind + GORM embutidos. Falta o **painel administrativo** — um clone do **ZimbraAdmin** (console legado em `:7071`) para gerenciar domínios, contas, aliases, listas e admins.

**Decisão do dono:** tudo é feito **aqui, no go-snappymail** (o binário único que já roda o webmail). O **go-postfixadmin não será alterado** — serve apenas de **referência/exemplo** do backend (schema Postfix/Dovecot, handlers, RBAC), já que muita coisa lá funciona e mostra o caminho. Portamos/adaptamos os padrões para cá, sem importar código.

## What Changes

- **Painel admin no próprio binário go-snappymail**, servido em `:7071` (http/https), ao lado do webmail (`:8082`). Um `serve`, dois serviços, ligados por config — sem processos/instalações separadas.
- **Frontend admin criado DO ZERO** — Vue 3 + TypeScript + Vite + **TailwindCSS** + tema/layout do ZimbraAdmin, em pasta própria **`frontend-admin/`**, separada do webmail (`frontend/`). Nada reaproveitado de UI existente; o tema ZimbraAdmin é recriado do zero (mesmo stack Tailwind do webmail para consistência).
- **Login idêntico ao ZimbraAdmin** (card com topo azul "Admin Console", labels à esquerda, botão Login no canto). **Toda a UI em inglês.**
- **Backend admin novo neste projeto** (GORM → MariaDB/PostgreSQL), tendo o go-postfixadmin como referência:
  - Modelos GORM do schema Postfix/Dovecot (domains, mailboxes, aliases, admins) — **banco separado** do banco de sessão do webmail.
  - Handlers REST sob `/api/v1/admin/*` (envelope JSON consistente) + `GET /api/v1/admin/overview` (contadores reais para a Home).
  - Auth própria do admin (JWT + RBAC granular: superadmin/domain_admin + permissões) — independente da sessão-cookie do webmail.
- **Backend melhorado com as skills `golang-*`** (security, testing, database, lint, error-handling, concurrency) onde precisar — **estender o que já existe no go-snappymail**, sem tocar no go-postfixadmin.
- **Testes de backend** (backend-first, `-race`, table-driven) e **lint** (golangci-lint + eslint/vue-tsc) como itens obrigatórios.
- **Documentação**: arquitetura do binário único (webmail + admin), organização de pastas, mapeamento telas↔endpoints; screenshots sempre em `docs/prints/`.

## Capabilities

### New Capabilities

- `admin-panel-ui`: o layout ZimbraAdmin Classic em Vue 3 + TailwindCSS (login, top bar, árvore de navegação, Home overview, list views, toolbars), servido pelo go-snappymail em `:7071`, com paridade visual contra o console real (`192.168.56.30:7071`) e UI em inglês.
- `admin-backend-api`: o backend admin no go-snappymail (modelos GORM Postfix/Dovecot → MariaDB/PostgreSQL, handlers `/api/v1/admin/*`, overview, JWT/RBAC), inspirado no go-postfixadmin (referência).

### Modified Capabilities

<!-- nenhuma — não há specs principais extraídas ainda -->

## Non-goals

- **Não alterar o go-postfixadmin** — é só referência; nada é importado nem modificado lá.
- Não recriar o toolkit DWT/AJAX do Zimbra — só o **look** em Vue + Tailwind.
- Sem logo/marca registrada do Zimbra (branding em texto).
- Fora desta fase: Monitor avançado, Zimlets, Class of Service completo, migração, backup/restore, certificados, Resources/Distribution Lists reais (stub por ora).
- Não reaproveitar frontends existentes para o admin — criado do zero.

## Impact

- **go-snappymail** (este projeto):
  - **`frontend-admin/` (NOVO)** — painel ZimbraAdmin, Vue 3 + TailwindCSS, do zero; build → `web/admin-dist`.
  - `main.go` — `//go:embed` passa a incluir `web/admin-dist`; `Makefile` — alvo `frontend-admin` + build do painel.
  - `internal/model/` + `internal/database/` — novos modelos GORM (Domain/Mailbox/Alias/Admin) e conexão ao banco Postfix/Dovecot (MariaDB/Postgres), separada da sessão.
  - `internal/handler/` (ou novo `internal/admin/`) — handlers `/api/v1/admin/*` + overview; auth JWT/RBAC.
  - `internal/server/` — segundo listener `:7071` (multi-listener, middleware/SPA/shutdown por serviço); rotas do painel.
  - `internal/config/` — blocos `[admin]` (porta 7071, tls) e `[database]` do mail; `[webmail]` mantém o atual.
  - `docs/` — arquitetura do binário único, dev-environment, mapeamento telas↔endpoints.
  - `openspec/` — esta change.
- **go-postfixadmin**: **inalterado** (referência de leitura apenas).
- Banco: GORM → MariaDB (lab Docker `:3306`) ou PostgreSQL, schema Postfix/Dovecot; separado do banco de sessão (SQLite/GORM) do webmail.
