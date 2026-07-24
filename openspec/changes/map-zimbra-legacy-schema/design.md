# Design: legacy Zimbra data-model mapping

## Scope

Documentation only. The output is `docs/zimbra-legacy/`. No code, no schema, no
runtime change. Everything documented was extracted from a live reference
install; it is a factual schema map, not Zimbra source.

## Extraction method (reproducible)

| Data | Source of truth | How extracted |
| --- | --- | --- |
| MySQL databases/tables | `information_schema.tables` | `mysql -N -e "show databases/tables"` |
| MySQL columns/types | `information_schema.columns` | one `SELECT … ORDER BY ordinal_position` |
| LDAP object classes | live directory | `zmprov ga|gd|gc <entity> objectClass` |
| LDAP attribute names | live directory | `zmprov ga|gd|gc` → attribute-name extraction |
| Counts | live | `zmprov gad|gaa|gac` |

Re-run any of these against a Zimbra box to refresh the docs.

## Decisions

- **D1 — read-only, always.** Any future code derived from these models treats
  Zimbra stores as read-only. The panel never AutoMigrates or writes to a Zimbra
  DB/LDAP. Rationale: Zimbra owns that data; corruption risk is unacceptable.
- **D2 — no coupling now.** The GORM structs live in the docs, not in
  `internal/`. They compile into nothing until a future change explicitly adds an
  isolated `internal/zimbra/` adapter behind a feature flag.
- **D3 — summarize, don't dump.** LDAP exposes ~450 account / 495 COS
  attributes. We document the DIT, object classes, and the operationally
  important attributes grouped by function; the exhaustive list stays queryable
  live rather than copied wholesale.
- **D4 — sharding is explicit.** Mailbox tables have no fixed table name; models
  target `mboxgroup<N>` resolved from `zimbra.mailbox.group_id`. Documented as
  `db.Table("mboxgroup%d.mail_item")`, not a `TableName()`.
- **D5 — reuse is concept-first.** We do not reuse Zimbra data as storage. The
  deliverable ranks *concepts* to borrow (COS, DLs, status enum, delegated-admin
  grants, usage reporting), each a candidate for its own later proposal.

## Document set

```
docs/zimbra-legacy/
  README.md                      index + how-Zimbra-stores-data
  01-mysql-databases.md          databases, tables, ER diagram (mermaid)
  02-mysql-tables-reference.md   column reference for key tables
  03-gorm-models.md              GORM structs + a usage example
  04-ldap-structure.md           DIT, object classes, key attributes, mapping
  05-architecture-reuse.md       architecture diagram + ranked reuse analysis
```

## Out of scope

- Implementing a COS/DL/status feature (each its own future change).
- Building the read-only Zimbra adapter.
- Any migration tooling.
