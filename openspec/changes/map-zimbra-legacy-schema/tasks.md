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

## 4b. Front-end webapp structure (expanded scope)
- [x] 4b.1 `06-webapp-structure.md` — webmail + admin folders, CSS/JS/HTML, skins, JS bundles
- [x] 4b.2 `07-framework-and-jars.md` — AjaxTk/DWT framework, backend Java stack, JAR/WAR manifest + how to open archives
- [x] 4b.3 `08-source-repos.md` — upstream FOSS repos map (links only, no copied source)
- [x] 4b.4 Copyright guardrail: structure/schema/behavior only, no source dumps

## 5. Validate
- [x] 5.1 Review the docs with codex (accuracy of schema/relationships)
- [~] 5.2 Review with agent — CLI returned no usable output non-interactively (timed out)
- [~] 5.3 Review with kilo — CLI returned no usable output non-interactively (timed out)
- [x] 5.4 Apply codex corrections (mail_item.type enum, GORM imports + nullable *string, COS wording, AjaxTk naming)

## 6. Follow-ups (each a separate future change — not in scope here)
- [ ] 6.1 Proposal: Class-of-Service-style defaults bundle
- [ ] 6.2 Proposal: first-class distribution lists
- [ ] 6.3 Proposal: account-status enum (active/locked/closed/…)
- [ ] 6.4 Proposal: optional read-only Zimbra reporting/migration adapter
