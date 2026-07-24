# Spec: admin-backend-api

Contrato REST do painel admin **no go-snappymail** (novos endpoints sob `/api/v1/admin/*`, GORM → MariaDB/PostgreSQL, auth JWT/RBAC), tendo o **go-postfixadmin como referência** (schema/handlers adaptados, não importados).

## ADDED Requirements

### Requirement: Config multi-serviço num binário
O binário SHALL ler blocos `[admin]`, `[webmail]` e `[database]` do `config.toml` e iniciar por um único `serve` todos os serviços com `enabled=true`, cada um em sua porta.

#### Scenario: Só admin
- **WHEN** `[admin] enabled=true` e `[webmail] enabled=false`
- **THEN** o processo serve apenas `:7071` e não abre a porta do webmail

#### Scenario: Ambos
- **WHEN** ambos `enabled=true`
- **THEN** um processo serve `:7071` (admin) e `:8082` (webmail) sem interferência

### Requirement: Isolamento total das superfícies admin × webmail
Os endpoints de admin (`/api/v1/admin/*`) e a SPA admin (`web/admin-dist`) SHALL ser registrados **apenas** no listener admin (`:7071`), em uma **instância Echo separada**, e **nunca** serem alcançáveis pela porta/listener do webmail (`:8082`) — nem por rota, nem por proxy, nem por SPA fallback. O inverso também: rotas do webmail não existem no router do admin.

#### Scenario: Admin não vaza pela porta do webmail
- **WHEN** uma requisição chega em `http://host:8082/api/v1/admin/overview` (porta do webmail)
- **THEN** o router do webmail responde **404** (a rota admin não existe ali); o handler admin nunca é invocado

#### Scenario: Webmail não vaza pela porta do admin
- **WHEN** uma requisição chega em `http://host:7071/api/v1/mail/INBOX` (porta do admin)
- **THEN** o router do admin responde **404**

#### Scenario: Bind restrito opcional
- **WHEN** `[admin] host = "127.0.0.1"`
- **THEN** o painel admin só aceita conexões locais em `:7071`, sem exposição externa

### Requirement: Overview endpoint
O backend SHALL expor `GET /api/v1/admin/overview` (prefixo `/api/v1/admin`, envelope JSON consistente do go-snappymail) retornando os contadores **que existem no schema** — accounts (mailboxes), domains, aliases, admins — para a Home, exigindo permissão de admin. Campos do console legado sem fonte no schema (COS, servers, active sessions, mail queue) SHALL ser omitidos ou retornados como `null`/`"n/a"`, nunca inventados.

#### Scenario: Agregados reais
- **WHEN** um admin autenticado chama `/api/v1/admin/overview`
- **THEN** recebe counts de accounts/domains/aliases/admins derivados via GORM; campos sem fonte vêm nulos/"n/a"

### Requirement: CRUD via endpoints v1 existentes
As telas SHALL operar sobre novos endpoints sob `/api/v1/admin/*` (domains, mailboxes, aliases, admins) — criados no go-snappymail, adaptados do go-postfixadmin — sob JWT/RBAC.

#### Scenario: Criar domínio
- **WHEN** o admin cria um domínio no painel via `/api/v1/domains`
- **THEN** o backend persiste via GORM e a lista reflete o novo domínio no envelope padrão

### Requirement: Persistência GORM MariaDB/PostgreSQL
O admin SHALL persistir no banco Postfix/Dovecot em MariaDB ou PostgreSQL via GORM (separado do banco de sessão do webmail), selecionável por `[database] driver`. Unit tests usam sqlmock ou container.

#### Scenario: Troca de driver
- **WHEN** `[database] driver` muda entre `mysql` e `postgres` com DSN válido
- **THEN** o painel funciona sem alteração de código (só config/migração)

### Requirement: Cobertura total de testes do backend admin
**Toda** rota `/api/v1/admin/*` SHALL ter testes automatizados cobrindo sucesso, validação, erro e permissão — nenhum endpoint sem teste. `go test -race ./...` SHALL passar; a suíte cobre CRUD de domains/mailboxes/aliases/admins, overview, auth/RBAC, isolamento de superfície e persistência (MariaDB e PostgreSQL). Regra do projeto: backend-first, cada rota entra com seu teste.

#### Scenario: Rota sem teste bloqueia
- **WHEN** um endpoint admin é adicionado/alterado sem teste correspondente
- **THEN** o CI falha (cobertura/lint) e o merge é bloqueado

#### Scenario: Suíte verde
- **WHEN** `go test -race ./...` roda no CI
- **THEN** todas as rotas admin passam (sucesso, validação, 401/403, not-found, erro de banco)

### Requirement: Auth admin (JWT/RBAC granular real)
O painel SHALL usar o RBAC existente com papéis reais (`superadmin`/`domain_admin`/…) e permissões granulares (`domains:read`, `mailboxes:write`, …). Cada rota/nó exige a permissão correspondente; sem ela → 403. Um `domain_admin` acessa Manage dos seus domínios mas não Configure (Servers/Global Settings).

#### Scenario: Acesso negado por permissão
- **WHEN** um `domain_admin` chama uma rota que exige `servers:read` (Global Settings)
- **THEN** o backend responde 403 e o nó correspondente fica oculto na árvore
