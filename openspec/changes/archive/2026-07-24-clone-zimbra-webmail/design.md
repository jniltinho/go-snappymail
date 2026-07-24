# Design: clone-zimbra-webmail

## Context

`add-zimbra-skin` reached palette/chrome parity with measured harmony values. The remaining gap (user-reported + observed): typography renders larger than the reference, ▾ controls are static, some hover/press effects and secondary controls are missing. Both UIs now render the same mailbox (`nilton@linuxpro.com.br` via the VM's IMAP), enabling honest A/B audits.

## Goals / Non-Goals

**Goals:** indistinguishable typography, working menus, matching effects; QA-agent audit with zero P1 as the gate.
**Non-Goals:** see proposal (no Briefcase/Calendar/threads/tags; squared corners kept; text branding).

## Decisions

1. **Type scale via Tailwind v4 theme tokens, scoped**: Tailwind v4 utilities resolve sizes from `--text-*` variables, so `[data-skin='zimbra'] { --text-sm: 12px; --text-sm--line-height: 1.35; --text-base: 13px; --text-xs: 11px; }` retunes the whole chrome without touching components. Alternative (per-component overrides) rejected: dozens of selectors, drift-prone.
2. **One 30-line `DropdownMenu.vue`** (slot trigger + items array or slot body; click-outside + Esc close; no dependency, no portal). Used by: user menu, Actions, New message ▾. Alternative (headless-ui lib) rejected: new dependency for three menus.
3. **Actions menu contents** = Mark read/unread, Flag/Unflag, Spam — mirrors Zimbra's Actions dropdown for messages; toolbar keeps the primary buttons (Reply/Reply to All/Forward/Archive/Delete/Spam) exactly as the reference shows them.
4. **QA agent as the loop driver**: after each iteration, run `qa-frontend-cloner` with both URLs + credentials; fix P1s first, then P2s; stop when 0 P1 and P2s are either fixed or waived.

## Risks / Trade-offs

- [Tailwind token override may not cover hardcoded px] → audit also measures elements styled by custom classes (`.tbtn`, `.sort-header` already 12px).
- [Zimbra 10 reference differs slightly from 8.8.15 harmony] → the VM is the contract; criarenet captures stay as secondary reference.

## Migration Plan

Frontend-only; ship per-iteration commits. Rollback = revert.

## Open Questions

None.
