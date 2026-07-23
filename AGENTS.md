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
| `docs/security.md` | higiene de segredos |
| `docs/validation-summary.md` | estado atual do projeto (handoff) |

## Estado atual (jul/2026)

P0 (foundation) e P1 (mail API) completos. P2 (frontend) parcial: layout 3 colunas, login, pastas/mensagens, dark mode, skins ok; **pendentes**: ComposerModal (TipTap), sanitização HTML (bluemonday), SSE, settings/contacts (P3). Detalhe: `docs/validation-summary.md` §9.
