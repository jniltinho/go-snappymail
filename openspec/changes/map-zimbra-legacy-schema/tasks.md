# Tasks: map the legacy Zimbra data model

## 1. Extract (from the reference install)
- [x] 1.1 List MySQL databases and tables (`information_schema`)
- [x] 1.2 Dump column/type info for `zimbra` + a representative `mboxgroup`
- [x] 1.3 Capture LDAP object classes for account / domain / COS
- [x] 1.4 Capture LDAP attribute names + entity counts

## 2. Document MySQL
- [x] 2.1 `01-mysql-databases.md` — databases, tables, `mail_item.type` map, ER diagram
- [x] 2.2 `02-mysql-tables-reference.md` — column reference for key tables
- [x] 2.3 `03-gorm-models.md` — GORM structs + read-only usage example

## 3. Document LDAP
- [x] 3.1 `04-ldap-structure.md` — DIT diagram, object classes, key attributes
- [x] 3.2 Mapping table: Zimbra LDAP concept → our PostfixAdmin schema

## 4. Architecture + reuse
- [x] 4.1 `05-architecture-reuse.md` — data-store architecture diagram
- [x] 4.2 Ranked "what we can borrow" analysis + explicit non-goals
- [x] 4.3 `README.md` index

## 5. Validate
- [ ] 5.1 Review the docs with codex (accuracy of schema/relationships)
- [ ] 5.2 Review with agent (completeness, reuse recommendations)
- [ ] 5.3 Review with kilo (architecture/diagram correctness)
- [ ] 5.4 Apply any corrections from the three reviews

## 6. Follow-ups (each a separate future change — not in scope here)
- [ ] 6.1 Proposal: Class-of-Service-style defaults bundle
- [ ] 6.2 Proposal: first-class distribution lists
- [ ] 6.3 Proposal: account-status enum (active/locked/closed/…)
- [ ] 6.4 Proposal: optional read-only Zimbra reporting/migration adapter
