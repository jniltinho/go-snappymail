# Resumo para validação — go-snappymail

Documento de handoff para revisão do estado atual do projeto (backend, frontend, skins, docs, build).

| Campo | Valor |
|-------|-------|
| **Repositório** | https://github.com/jniltinho/go-snappymail |
| **Branch** | `main` |
| **Última revisão** | jul/2026 — hardening backend + skin Outlook ready |
| **Módulo Go** | `go-snappymail` |
| **Porta padrão** | `8082` |
| **Princípio** | Backend first, frontend second |

---

## 1. Objetivo do projeto

Webmail self-hosted em Go com UX inspirada no SnappyMail:

- Binário único com SPA Vue 3 **embedded** (`go:embed web/dist`)
- IMAP/SMTP passthrough (sem armazenar mail localmente)
- Sem runtime PHP

---

## 2. O que já foi entregue

### P0 — Foundation ✅

| Item | Status |
|------|--------|
| Cobra CLI (`init`, `migrate`, `serve`, `version`) | ✅ |
| Config TOML + env `GOSM_*` | ✅ |
| GORM + SQLite (MariaDB/Postgres preparados) | ✅ |
| Echo v5 + middleware (recover, CSRF, rate limit login, slog) | ✅ |
| Sessão cookie `gsn_session` | ✅ |
| Auth IMAP (`POST /auth/login`) | ✅ |
| Embed SPA + fallback `index.html` | ✅ |
| Makefile (`build`, `build-prod`, `test`, `run`) | ✅ |
| Docker lab porta 8082 | ✅ |
| Testes unitários (handler, session, server, ui) | ✅ |
| Documentação base (`docs/`) | ✅ |

### P1 — Mail API ✅

| Área | Endpoints | Status |
|------|-----------|--------|
| Pastas | `GET/POST /folders`, rename, delete, count | ✅ |
| Mensagens | list, read, flag, move, delete, attachments, raw | ✅ |
| Compose | send, draft, upload | ✅ |
| Search | `GET /search` (IMAP SEARCH) | ✅ |
| Quota | `GET /auth/quota` | ✅ |

### P2 — Frontend (parcial) 🟡

| Item | Status | Notas |
|------|--------|-------|
| Vue 3 + Vite + TypeScript + Pinia | ✅ | `frontend/` |
| Tailwind CSS v4 + dark mode toggle | ✅ | Classe `.dark` em `<html>` |
| Build → `web/dist/` embedded | ✅ | `make build` |
| LoginView + auth wired | ✅ | CSRF, cookies |
| Layout 3 colunas desktop | ✅ | `App.vue` — grid fixo |
| FolderSidebar, MessageList, ReadingPane, AppToolbar | 🟡 | Scaffold funcional, features limitadas |
| Atalhos j/k | 🟡 | Parcial em `App.vue` |
| ComposerModal (TipTap) | ❌ | Não implementado |
| Search bar na toolbar | ❌ | API existe, UI não |
| SSE / notificações | ❌ | |
| Sanitização HTML (bluemonday) | ✅ | Server-side em `message_read.go` (policy package-level) |
| Layout mobile / responsive | ❌ | **Adiado** — virá com layouts por skin |

### Hardening backend ✅ (review com skills golang)

Revisão completa de `internal/`, `cmd/` e `main.go` com 18 correções aplicadas (`go vet` limpo, `make test -race` passando):

| Categoria | Correções |
|-----------|-----------|
| Segurança | Gate do `imap_host` por `show_host_input` (anti-SSRF); validação de `secret_key` ≥ 32 bytes no startup; `Content-Disposition` via `mime.FormatMediaType` (anti header-injection); cookie CSRF com flag `Secure`; erros IMAP/SMTP internos não vazam mais para o cliente (logados via `slog`) |
| Concorrência | Data race em `session.Get` (LastUsed sob write lock); rate limiter movido para dentro do closure (sem estado global compartilhado) |
| Robustez | `AutoMigrate` de sessões no `serve` (persistência em instalações novas); `GetQuota` retorna erros reais (só quota-unsupported → nil); erros de anexo/append não são mais engolidos; `RunE` retorna erro em vez de `os.Exit` |
| Limpeza | Código morto removido (`session.All`, `imap.MessageCount`, `FetchAllUIDs`, `LoginRateLimit`); política bluemonday hoisted para package-level; helper único de random ID; comentários `GORC_*`→`GOSM_*`; `slices.Contains` em `NormalizeSkin` |

Pendências conhecidas do review (deliberadas): endpoint `/compose/upload` grava arquivos que nada consome (remover ou integrar ao Send); structs de config especulativos (`PushConfig`, `ActiveSyncConfig`) sem consumidor.

### Skins ✅ (v3 — Outlook ready)

| Item | Status |
|------|--------|
| `config.toml` → `[ui] skin` | ✅ |
| Catálogo Go (`internal/ui/skins.go`) | ✅ |
| Manifest TS (`frontend/src/skins/manifest.ts`) | ✅ |
| CSS por skin (`[data-skin='id']`) | ✅ snappymail, gmail, outlook |
| Tokens login (`--skin-login-*`) | ✅ |
| `GET /api/v1/ui/config` + array `skins[]` | ✅ |
| `make new-skin ID=x REGISTER=1` | ✅ |
| `make validate-skins` | ✅ |
| Guia completo `docs/skins.md` | ✅ |

Skins built-in:

| Id | `ready` | Descrição |
|----|---------|-----------|
| `snappymail` | `true` | Default, azul SnappyMail |
| `gmail` | `false` | Só tokens/cores |
| `outlook` | `true` | Fluent 2 minimalista (light/dark), accent bar na seleção |
| `carbonio` | `true` | Zextras Carbonio (light/dark), login navy + card branco; ref. `docs/prints/carbonio/` |

Regra de UI global: **layout quadrado, sem cantos arredondados** (removido `rounded-*`; skins não devem adicionar `border-radius`).

Prints atualizados em `docs/prints/` (07–12): login e inbox light/dark nas skins snappymail e outlook.

---

## 3. Arquitetura resumida

```
config.toml
    ↓
Go binary (Echo :8082)
    ├── /api/v1/*     REST (auth, mail, compose, search, ui/config)
    └── /*            SPA embedded (web/dist via go:embed)

frontend/ (dev)  →  npm run build  →  web/dist/  →  go build embeds
```

**Embed** (`main.go`):

```go
//go:embed all:web/dist
//go:embed all:web/files
var embeddedFiles embed.FS
```

Em produção **não precisa Node** — só o binário compilado após `make build`.

---

## 4. API implementada (prefixo `/api/v1`)

### Público (sem auth)

| Método | Path | Descrição |
|--------|------|-----------|
| GET | `/version` | Versão do app |
| GET | `/ui/config` | Skin, catálogo, rows_per_page, etc. |

### Auth

| Método | Path | Descrição |
|--------|------|-----------|
| POST | `/auth/login` | Login IMAP + sessão |
| POST | `/auth/logout` | Logout |
| GET | `/auth/me` | Usuário atual + skin |
| GET | `/auth/quota` | Quota IMAP |

### Mail (auth required)

| Método | Path |
|--------|------|
| GET/POST | `/folders`, `/folders/rename`, `/folders/delete` |
| GET | `/folders/:name/count` |
| GET | `/mail/:mailbox` |
| GET | `/mail/:mailbox/:uid` |
| GET | `/mail/:mailbox/:uid/attachment/:part` |
| GET | `/mail/:mailbox/:uid/download`, `/raw` |
| POST | `/mail/:mailbox/:uid/flag`, `/move` |
| DELETE | `/mail/:mailbox/:uid`, `/mail/:mailbox` |
| POST | `/compose/send`, `/compose/draft`, `/compose/upload` |
| GET | `/search` |

Detalhes: [docs/api.md](api.md)

---

## 5. Frontend — estrutura

```
frontend/src/
├── api/client.ts          # axios + CSRF
├── components/
│   ├── LoginView.vue
│   ├── FolderSidebar.vue
│   ├── MessageList.vue
│   ├── ReadingPane.vue
│   └── AppToolbar.vue
├── stores/
│   ├── auth.ts
│   ├── mail.ts
│   └── settings.ts
├── skins/
│   ├── manifest.ts        # catálogo TS (sync com Go)
│   ├── index.css          # imports das skins
│   ├── snappymail.css | gmail.css | outlook.css
│   ├── bootstrap.ts       # fetch /ui/config no boot
│   └── apply.ts           # data-skin no <html>
├── App.vue                # layout 3 colunas
├── main.ts
└── style.css              # Tailwind + classes login
```

---

## 6. Documentação

| Documento | Conteúdo |
|-----------|----------|
| [README.md](../README.md) | Quick start |
| [docs/README.md](README.md) | Índice |
| [docs/architecture.md](architecture.md) | Design, fases P0–P3 |
| [docs/development.md](development.md) | Build, test, dev workflow |
| [docs/configuration.md](configuration.md) | config.toml, env vars |
| [docs/api.md](api.md) | REST reference |
| [docs/skins.md](skins.md) | **Guia completo de implementação de skins** |
| [docs/lab.md](lab.md) | Docker/Vagrant lab |
| [docs/security.md](security.md) | Secrets, git hygiene |
| [openspec/changes/go-snappymail-foundation/](../openspec/changes/go-snappymail-foundation/) | Proposta OpenSpec + tasks |

---

## 6.1 Validações executadas no lab (jul/2026)

| Validação | Resultado |
|-----------|-----------|
| Envio local via webmail (`POST /compose/send`, SMTP `localhost:25` → mailserver → INBOX) | ✅ `{"status":"sent"}` — prints 13–14 |
| Login SnappyMail PHP (`:8888`) com conta nova `user01@linuxpro.com.br` | ✅ — print 19 |
| Login Go CubeMail (`:8080`) com `user02@criarenet.com` | ✅ — print 20 |
| 4 domínios novos via painel Go-PostfixAdmin (linuxpro.com.br, criare-net.com.br, uol.com.br, criarenet.com) | ✅ — print 18 |
| 20 contas por domínio novo (80 via CLI `postfixadmin mailbox --add` + 4 `user15@` via painel web) | ✅ — print 21 |
| 5 aliases por domínio novo (postmaster/abuse/admin→ti, info→contato, atendimento→suporte) | ✅ — `docker/lab/aliases.txt` |

Correções de lab feitas no caminho:

- **postfix ↔ MariaDB**: templates `docker/mailserver/templates/postfix/sql/*.cf` fixavam `hosts = localhost`; agora usam `@@DB_HOST@@` (bug de "Temporary lookup failure" no RCPT).
- **Cert do mailserver**: gerado sem SANs — Go rejeita; `entrypoint.sh` agora emite `subjectAltName` (mail.test.local, mailserver, localhost).
- **go-cubemail**: monta `mail_ssl` e confia no cert do lab no entrypoint (era "certificate signed by unknown authority").
- **SnappyMail domain template**: faltava bloco `Sieve` (erro fatal no 2.38.2); adicionado em `docker/snappymail/domains/_default.json`.
- **SMTP do go-snappymail**: novo `[smtp] insecure_skip_verify` (lab/dev) espelhando o do IMAP.
- **seed-lab.sh**: agora semeia aliases de `docker/lab/aliases.txt` via SQL.

## 7. Lab / testes

**Usuário de teste (Docker lab):**

- Email: `user@test.local`
- Senha: `Password1@`

**Comandos de validação:**

```bash
# Sync skins Go ↔ TS ↔ CSS
make validate-skins

# Testes Go (race + coverage)
make test

# Build frontend + binário embedded
make frontend
make build

# Testes de integração IMAP (lab Docker rodando)
make test-integration

# Dev (2 terminais)
make run              # :8082
make frontend-dev     # Vite :5173
```

**Checklist manual sugerido:**

- [ ] `make validate-skins` — OK
- [ ] `make test` — passa
- [ ] `cd frontend && npm run build` — passa
- [ ] `make build && ./dist/go-snappymail serve` — SPA carrega em `/`
- [ ] Login com lab account — redirect para inbox
- [ ] Pastas e mensagens carregam
- [ ] Toggle dark mode funciona
- [ ] Trocar `[ui] skin` no config — login e inbox mudam cores
- [ ] `GET /api/v1/ui/config` retorna `skins[]` e `skin` correto
- [ ] Binário roda **sem** Node instalado (só o `.dist/go-snappymail`)

---

## 8. Decisões explícitas (não são bugs)

| Decisão | Motivo |
|---------|--------|
| Mobile/responsive adiado | Virá com layouts por skin (`AppShell`), não breakpoint genérico |
| Skins gmail/outlook `ready: false` | Só paleta de cores; layout completo TBD |
| Skin server-side | Todos os users veem a mesma skin (`config.toml`) |
| Dark mode client-side | `localStorage` `gsn_dark`, não vem do servidor |
| `web/dist/` no `.gitignore` | Gerado no build; embed no binário |
| Backend before frontend | API P1 completa antes de polir UI P2 |

---

## 9. Pendências principais (próximas fases)

| Prioridade | Item |
|------------|------|
| P2 | ComposerModal + TipTap |
| P2 | Search bar na toolbar |
| P2 | SSE `/api/v1/events` |
| P2 | Polir componentes 7.x (flags, attachments UI, compose actions) |
| P3 | Settings, identities, contacts |
| Futuro | Layout mobile + AppShell por skin |

Ver tasks completas: [openspec/changes/go-snappymail-foundation/tasks.md](../openspec/changes/go-snappymail-foundation/tasks.md)

---

## 10. Commits recentes (referência)

```
c8d8589 refactor: refine skin system with manifest sync and implementation docs
2b1f020 docs: add skin creation guide, template and make new-skin scaffold
54a37ee feat: add configurable UI skins via config.toml and frontend registry
c65f74f docs: update tasks and dev guide for Vue frontend workflow
2437ed4 feat: add Vue 3 frontend with SnappyMail-style inbox shell
5d57ba1 feat: add compose, draft, upload and search API (P1)
21a80d8 feat: add P1 mail API (folders, messages, quota)
f34478e docs: add project guides and exclude local prints from git
```

---

## 11. Pontos de atenção para o revisor

1. **Sincronia skins** — Go catalog, TS manifest e `index.css` devem bater; rodar `make validate-skins`.
2. **Embed** — `go build` sem `make frontend` antes pode embutir SPA desatualizada ou vazia.
3. **OpenSpec tasks 7.1–7.4** — marcadas `[ ]` no tasks.md, mas componentes Vue existem em versão scaffold; validar profundidade vs checklist.
4. **Integração** — `make test-integration` requer lab Docker (`IMAP_TEST_*` env vars).
5. **Segurança** — sanitização bluemonday é server-side; remote image blocking (parte do item 7.10) ainda pendente no frontend.
6. **AGENTS.md** — guia de orientação para IAs na raiz do repo; manter em sync com docs.

---

*Gerado para handoff de validação — jul/2026.*
