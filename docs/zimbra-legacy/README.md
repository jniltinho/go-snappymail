# Zimbra Legacy — Data Model Mapping (reference)

> **Status:** documentation / research only. No application code depends on this.
> Purpose: map the legacy Zimbra data model (MySQL + LDAP) so we can decide what,
> if anything, to reuse or emulate in the go-snappymail admin panel.

The legacy reference is Zimbra 10.1.17 FOSS (VM `192.168.56.30`). All structure
here was extracted live from that install (`information_schema`, `zmprov`,
`ldapsearch`) — it is a **factual schema map**, not Zimbra source.

## Contents

| Doc | What it covers |
| --- | --- |
| [01-mysql-databases.md](01-mysql-databases.md) | The MySQL databases, every table, and the entity-relationship diagram |
| [02-mysql-tables-reference.md](02-mysql-tables-reference.md) | Column-level reference for the key tables |
| [03-gorm-models.md](03-gorm-models.md) | Ready-to-use GORM structs for the reusable tables |
| [04-ldap-structure.md](04-ldap-structure.md) | The OpenLDAP directory: DIT, object classes, key attributes |
| [05-architecture-reuse.md](05-architecture-reuse.md) | Zimbra architecture diagram + what we can realistically reuse |

## How Zimbra stores data (one paragraph)

Zimbra splits its state across two stores:

- **OpenLDAP** — the *directory*: accounts, domains, aliases, distribution
  lists, classes of service (COS), servers, global config. This is the
  authoritative source for "who exists and what they're allowed to do". Our
  admin panel manages exactly this class of object (today via a PostfixAdmin
  SQL schema instead of LDAP).
- **MySQL/MariaDB** — the *mailbox metadata*: every message, folder, tag,
  contact, appointment, etc., stored as rows. Message **blobs** live on disk
  (the "store"), referenced from MySQL by `locator`/`blob_digest`. MySQL is
  **sharded**: a central `zimbra` DB maps each account to a `mboxgroup<N>` DB
  that holds that mailbox's items.

## Snapshot of the reference install

- Domains: **3** · Accounts: **43** · Classes of Service: **2**
- MySQL databases: `zimbra`, `mboxgroup1..6`, `chat`
- LDAP account objectClasses: `inetOrgPerson`, `zimbraAccount`, `amavisAccount`
- ~450 `zimbra*` attributes available on an account; **495** on a COS
