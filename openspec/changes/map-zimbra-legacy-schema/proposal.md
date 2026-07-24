# Proposal: Map the legacy Zimbra data model (MySQL + LDAP)

## Why

The go-snappymail admin panel clones ZimbraAdmin visually and manages a
PostfixAdmin SQL provisioning schema. Zimbra itself stores provisioning in
**OpenLDAP** and mailbox metadata in **sharded MySQL**. Before we decide whether
to adopt any of Zimbra's richer concepts (Class of Service, distribution lists,
per-account features, usage reporting) — or to support reading a Zimbra box for
migration/coexistence — we need an accurate, documented map of that data model.

This is a **documentation / research** change. It produces no application code
and touches no runtime behavior. It exists so a future implementation decision
is grounded in the real schema, not guesswork.

## What changes

- Add `docs/zimbra-legacy/` — a factual map of the legacy Zimbra data model,
  extracted live from the reference install (Zimbra 10.1.17 FOSS,
  `192.168.56.30`):
  - MySQL databases, every table, and an ER diagram
  - Column-level reference for the key tables
  - Ready-to-use GORM structs for the reusable/reportable tables
  - OpenLDAP directory: DIT, object classes, key attributes
  - Architecture diagram + a ranked reuse analysis
- No changes to `internal/`, `frontend/`, `frontend-admin/`, or the DB schema.

## What does NOT change

- No new dependencies, no migrations, no endpoints, no models compiled into the
  binary. The GORM structs are documentation samples for a *possible future*
  isolated read-only adapter, not code that ships now.
- The panel remains authoritative on PostfixAdmin SQL; nothing here couples it
  to Zimbra LDAP or the Zimbra mailstore.

## Impact

- Affected: `docs/zimbra-legacy/**` (new), this OpenSpec change.
- Risk: none (documentation only).
- Enables later, separately-proposed changes: a COS-like defaults bundle,
  first-class distribution lists, an account-status enum, and an optional
  read-only Zimbra reporting/migration adapter.
