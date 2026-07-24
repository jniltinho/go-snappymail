# Tasks: add-zimbra-skin

## 1. Skin registration

- [x] 1.1 Scaffold `frontend/src/skins/zimbra.css` (`make new-skin ID=zimbra REGISTER=1` or manual) and fill light + `.dark` token sets from the measured palette in design.md
- [x] 1.2 Register in `internal/ui/skins.go` catalog (`ready: true`, alias `classic`) and in `frontend/src/skins/manifest.ts`; import in `frontend/src/skins/index.css`
- [x] 1.3 `make validate-skins` passes

## 2. Default skin flip

- [x] 2.1 `defaultSkinID = "zimbra"` in `internal/ui/skins.go`; update fallback cases in `internal/ui/skins_test.go`; `go test ./internal/ui/`
- [x] 2.2 `DEFAULT_SKIN = 'zimbra'` in `manifest.ts`; `data-skin="zimbra"` in `frontend/index.html`
- [x] 2.3 `[ui] skin = "zimbra"` + updated comment in `web/files/config.default.toml`
- [x] 2.4 Update default-skin references in `AGENTS.md` and `docs/skins.md` (catalog table, fallback text, sample outputs)

## 3. Zimbra chrome

- [x] 3.1 Scoped chrome rules in `zimbra.css`: blue top bar, tab-strip look, white toolbar with 1px `#bfbfbf` squared buttons, primary action `#0087c3`, selection `#99cae7`, list header `#f2f2f2`, unread bold `#333`, focus fill `#ffffcc`
- [x] 3.2 LoginView brand-area tweak: text-only brand block stylable by skin tokens (white on `#007cc3` card, white→`#ededed` page gradient)

## 4. Flow additions (existing APIs, no Go changes)

- [x] 4.1 `mail.ts` store: `moveMessage` (uses `POST /mail/:mailbox/:uid/move`, resolves/creates `Archive` folder) and `setSeen` (uses `flag=seen`)
- [x] 4.2 AppToolbar: Archive button + mark read/unread button, wired to store actions, English labels/tooltips

## 5. Flow verification (under zimbra skin, docker lab)

- [x] 5.1 Verify search, list (unread bold + counts), read pane (sanitized HTML + inline CIDs), flag toggle, delete
- [x] 5.2 Verify drafts + send at API level (`/compose/draft`, `/compose/send` → Drafts/Sent folders)
- [x] 5.3 Verify archive + mark unread end-to-end in the UI

## 6. Visual parity loop (acceptance bar)

- [x] 6.1 `make frontend && make build`; run app against docker lab; capture login + inbox + dark screenshots
- [x] 6.2 Side-by-side compare with `docs/prints/zimbra/00-login.png`, `01-advanced-inbox.png`, `02-advanced-compose.png`; adjust tokens/chrome; repeat until indistinguishable (squared corners + text brand excepted)
- [x] 6.3 Final: `make test`, `make validate-skins`, screenshots saved to `docs/prints/` per convention
