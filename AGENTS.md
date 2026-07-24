# AGENTS.md — orientação para agentes de IA

Guia rápido para qualquer IA (Claude, Copilot, Cursor, etc.) trabalhar neste repositório.

## O que é este projeto

**go-snappymail**: webmail self-hosted em Go, UX inspirada no SnappyMail.
Binário único (Go embute a SPA Vue 3 via `go:embed web/dist`), IMAP/SMTP passthrough (não armazena mail), sem PHP. Porta padrão `8082`.

## Regras do projeto

1. **Backend first, frontend second** — API Go completa (com testes) antes de UI.
2. **Sync de skins é obrigatório** — catálogo Go (`internal/ui/skins.go`), manifest TS (`frontend/src/skins/manifest.ts`) e imports CSS (`frontend/src/skins/index.css`) devem bater. Sempre rode `make validate-skins` após mexer em skins.
3. **`web/dist/` é gerado** (gitignored) — nunca edite; `go build` sem `make frontend` antes embute SPA desatualizada.
4. **Nunca commite segredos** — veja `docs/security.md`. `docs/prints/` também é gitignored (screenshots locais).
7. **Screenshots sempre em `docs/prints/`** — ao validar UI/paridade, tire prints e salve **somente** em `docs/prints/` (subpasta por contexto, ex. `docs/prints/zimbra/`, `docs/prints/zimbra/qaN/`). Nunca em `/tmp` ou na raiz. Tire prints sempre que precisar comparar com a referência (webmail `192.168.56.30`, admin `192.168.56.30:7071`); o `agent-browser` grava com caminho absoluto para essa pasta. É gitignored — servem de evidência local do QA.
5. **Tasks oficiais** em `openspec/changes/go-snappymail-foundation/tasks.md` — marque `[x]` ao concluir; use as skills OpenSpec (`opsx:*`) para propostas/arquivo.
6. **Go idiomático** — erros com `fmt.Errorf("...: %w", err)`, `log/slog` para logs, testes table-driven com `-race`. Skills `samber/cc-skills-golang@*` disponíveis via CLAUDE/Skill.

## Estrutura

```
main.go              # go:embed web/dist + web/files; delega para cmd/
cmd/                 # Cobra: init, migrate, serve, version
internal/
├── config/          # TOML (config.toml) + env GOSM_*
├── database/        # GORM (SQLite default; MariaDB/Postgres prontos)
├── handler/         # HTTP handlers REST /api/v1/*
├── imap/            # cliente go-imap/v2 (auth, folders, messages, search)
├── smtp/            # envio
├── model/           # modelos GORM
├── server/          # Echo v5, middleware (CSRF, rate-limit, slog), rotas, SPA fallback
├── session/         # sessão cookie gsn_session (memória + GORM)
└── ui/              # catálogo de skins (autoridade server-side)
frontend/            # Vue 3 + TS + Vite + Pinia + Tailwind v4 (build → web/dist/)
├── src/components/  # LoginView, FolderSidebar, MessageList, ReadingPane, AppToolbar
├── src/stores/      # auth, mail, settings
└── src/skins/       # manifest.ts, *.css por skin ([data-skin='id']), bootstrap/apply
web/files/           # config.default.toml (embedded)
docker/              # lab: mailserver, mariadb, postfixadmin, go-snappymail
docs/                # architecture, api, configuration, development, skins, lab, security
openspec/            # proposta + tasks (fonte de verdade do roadmap)
scripts/             # new-skin.sh, validate-skins.sh
```

## Comandos essenciais

```bash
make test              # go test -race + coverage
make validate-skins    # sync Go ↔ TS ↔ CSS das skins
make frontend          # npm build → web/dist/
make build             # frontend + binário embedded → dist/go-snappymail
make run               # serve :8082 (dev)
make frontend-dev      # Vite :5173 (proxy para :8082)
make test-integration  # requer lab Docker (IMAP_TEST_* env)
make new-skin ID=x REGISTER=1   # scaffold de skin nova
```

## Skins (resumo — guia completo em docs/skins.md)

- Server-side: `config.toml` → `[ui] skin` (zimbra | snappymail | gmail | outlook | carbonio). Exposto em `GET /api/v1/ui/config`.
- 100% token-driven: cada skin define `--color-*` e `--skin-login-*` em `[data-skin='<id>']` (+ variante `.dark`). Consumidos por `frontend/src/style.css` e pelo `@theme` do Tailwind.
- `ready: false` exibe banner de preview na UI. `zimbra` (default, Zimbra 8 Classic), `snappymail`, `outlook` (Fluent minimalista) e `carbonio` (Zextras) estão `ready`.
- Regra global de UI: layout quadrado — sem `border-radius`/`rounded-*`.
- Dark mode é client-side (`localStorage` `gsn_dark`, classe `.dark` no `<html>`).

## Lab / credenciais de teste

Docker lab (`docker/docker-compose.yml`): mailserver IMAP em `localhost:143/993`, app em `:8082`.
Usuário de teste: `user@test.local` / `Password1@`.

## API

Prefixo `/api/v1`. Público: `/version`, `/ui/config`. Auth: `/auth/login|logout|me|quota`. Mail: `/folders*`, `/mail/:mailbox[...]`, `/compose/*`, `/search`. Referência completa: `docs/api.md`.

## Documentação

| Doc | Conteúdo |
|-----|----------|
| `docs/architecture.md` | design e fases P0–P3 |
| `docs/development.md` | setup dev, workflows |
| `docs/configuration.md` | config.toml e env `GOSM_*` |
| `docs/api.md` | referência REST |
| `docs/skins.md` | criar/implementar skins |
| `docs/lab.md` | lab Docker/Vagrant |
| `docs/deployment-dev.md` | VM Vagrant híbrida: Go via systemd, MariaDB/SnappyMail via docker compose, Postfix/Dovecot nativos |
| `docs/dev-environment.md` | subir todo o ambiente de dev (webmail + admin + Zimbra por snapshot) |
| `docs/zimbra-vagrant-foss.md` | VM Zimbra FOSS 10.1.17 (referência da skin zimbra): instalação, domínios, contas |
| `docs/security.md` | higiene de segredos |
| `docs/validation-summary.md` | estado atual do projeto (handoff) |

## Subir o ambiente de desenvolvimento

Guia completo em [`docs/dev-environment.md`](docs/dev-environment.md). Resumo:

### 1. Servidor de e-mail (Zimbra FOSS — já instalado, via snapshot)

A VM `vagrant/zimbra` roda **Zimbra FOSS 10.1.17.p2** e serve de referência do webmail e de servidor IMAP/SMTP real. **Não reinstale do zero** (leva 20-40 min) — restaure o snapshot `zimbra-installed`:

```bash
cd vagrant/zimbra
vagrant snapshot restore zimbra-installed   # sobe a VM já instalada em ~1-2 min
# 1ª vez (sem a VM criada): vagrant up  → depois  vagrant snapshot save zimbra-installed
```

- IP `192.168.56.30` · FQDN `mail.zimbra.test` (adicione no `/etc/hosts` do host)
- Webmail `https://192.168.56.30/` · Admin console `https://mail.zimbra.test:7071/`
- Contas: `admin@zimbra.test`, `nilton@linuxpro.com.br` (+53 outras), senha `Password1@`
- Domínios: `zimbra.test`, `linuxpro.com.br`, `criarenet.com`
- Caixa `nilton@linuxpro.com.br` tem ~500 mensagens de teste (com anexos) para dev
- Detalhes/administração (`zmprov`): [`docs/zimbra-vagrant-foss.md`](docs/zimbra-vagrant-foss.md)

### 2. go-snappymail (webmail Go)

```bash
make run                       # serve :8082 (dev, config.toml local)
# ou apontando para a VM Zimbra (mesma caixa do webmail de referência):
#   [imap] host=192.168.56.30 port=993 tls=true insecure_skip_verify=true tls_server_name="mail.zimbra.test"
#   [smtp] host=192.168.56.30 port=587 starttls=true insecure_skip_verify=true
make frontend-dev              # Vite :5173 (HMR, proxy → :8082)
```

Login de teste: `nilton@linuxpro.com.br` / `Password1@`.

### 3. Lab Docker alternativo (mailserver local)

`docker/docker-compose.yml`: mailserver IMAP `:143/:993`, app `:8082`, usuário `user@test.local` / `Password1@`. Ver [`docs/lab.md`](docs/lab.md).

### Portas do ecossistema "Zimbra em Go"

| Serviço | Porta | Repo |
|---|---|---|
| go-snappymail (webmail) | 8082 | este |
| **painel-admin (ZimbraAdmin-like)** | **7071 (http/https)** | go-postfixadmin (backend) |
| Zimbra webmail (referência) | 443 (VM) | vagrant/zimbra |
| Zimbra admin console (referência) | 7071 (VM) | vagrant/zimbra |

## Arquitetura do ecossistema (recriar o Zimbra em Go + Vue)

O projeto é um **guarda-chuva** de binários Go independentes, cada um com seu Vue 3 embutido, todos usando o mesmo servidor Postfix/Dovecot + banco (MariaDB/PostgreSQL via GORM). Cada peça tem seu próprio repo/escopo e **não interfere** nos outros:

```
Zimbra-em-Go
├── go-snappymail      webmail (este repo)          — porta 8082 — clone do Zimbra Web
├── go-postfixadmin    backend admin + painel :7071 — clone do ZimbraAdmin (fase atual)
│   ├── internal/         API Go (GORM → MariaDB/PostgreSQL), JWT/RBAC
│   ├── frontend/         Vue 3 UI atual (neo-brutalism) → web/dist       (NÃO tocar)
│   └── frontend-admin/   Vue 3 painel ZimbraAdmin → web/admin-dist       (NOVO, separado)
├── mailserver         Postfix + Dovecot            — VM/Docker
└── vagrant/zimbra     Zimbra FOSS instalado        — referência visual/funcional
```

Padrão de cada peça: **single binary** (Go `embed` da SPA Vue), config TOML + env, GORM, testes `-race`, skins/layout clonados da UI legada correspondente do Zimbra. Roadmap detalhado da fase admin: proposta OpenSpec `admin-panel-zimbra` em `go-postfixadmin/openspec/`.

## Estado atual (jul/2026)

P0 (foundation) e P1 (mail API) completos. P2 (frontend) do webmail: layout Zimbra Classic clonado e validado (skin `zimbra` default, abas Contacts/Calendar/Tasks/Preferences, composer, DnD, toast). **Próxima fase**: painel admin `:7071` (go-postfixadmin + frontend ZimbraAdmin). Detalhe: `docs/validation-summary.md`.
