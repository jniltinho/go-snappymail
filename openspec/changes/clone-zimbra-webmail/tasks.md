# Tasks: clone-zimbra-webmail

## 1. QA infrastructure

- [x] 1.1 Agent `.claude/agents/qa-frontend-cloner.md` (agent-browser driver, P1/P2/P3 report)
- [x] 1.2 A/B environment: test instance IMAP/SMTP → VM Zimbra (same mailbox both UIs)
- [x] 1.3 Baseline audit run (6 P1 / 14 P2 / 5 P3); findings triaged

## 2. Typography (P1)

- [x] 2.1 Scoped Tailwind v4 type tokens under `[data-skin='zimbra']` (`--text-sm: 12px`, `--text-base: 13px`, `--text-xs: 11px`)
- [x] 2.2 QA re-measure: font family/size/weight equal on toolbar, rows, tree, tabs

## 3. Menus (P1)

- [x] 3.1 `DropdownMenu.vue` (click toggle, click-outside + Esc close, slot items)
- [x] 3.2 User menu (top right): username ▾ → Dark mode toggle, Logout
- [x] 3.3 Toolbar Actions ▾ → Mark read/unread, Flag/Unflag, Spam
- [x] 3.4 New message ▾ split-button menu
- [x] 3.5 Zimbra menu chrome (white, 1px #bfbfbf, #CCE5F3 hover, squared)

## 4. Effects and controls (P2)

- [x] 4.1 Audit hover/active on every control vs reference; fix missing states
- [x] 4.2 Toolbar right group (Read More / View ▾) — render per reference
- [x] 4.3 Any remaining P2s from audit — fix or waive in proposal Non-goals

## 5. QA loop (acceptance)

- [x] 5.1 Re-run `qa-frontend-cloner` after each iteration
- [x] 5.2 Final audit (iteration 3): **0 P1, 1 P2 (composer — waived to foundation), 4 P3 (fixed post-audit)** — acceptance PASS; all 13 iteration-2 findings verified FIXED
