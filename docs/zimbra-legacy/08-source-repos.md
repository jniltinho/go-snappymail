# Zimbra FOSS — Source Repositories Map

Where the actual source lives, so a future deep dive reads the code from the
upstream open-source repos instead of copying anything into this project. The
reference install is the **FOSS 10.1.x** build (packaged via
`maldua/zimbra-foss`, which assembles the upstream `Zimbra/*` repos).

> We reference these repos; we do not vendor or copy their code. Zimbra server
> code is copyrighted. Read upstream, document behavior here.

## Build / distribution

| Repo | What it is |
| --- | --- |
| `github.com/maldua/zimbra-foss` | Community FOSS build/installer (the releases the lab VM uses) |
| `github.com/Zimbra/zm-build` | Official build orchestration |
| `github.com/Zimbra/zm-zcs` | Top-level assembly (ZCS package) |
| `github.com/Zimbra/zm-core-utils` | `zmprov`, `zmlocalconfig`, `zmcontrol`, admin CLIs |

## Front-end (what our clone cares about)

| Repo | Maps to install path | Contents |
| --- | --- | --- |
| `github.com/Zimbra/zm-ajax` | `webapps/*/js/ajax`, `css/dwt.css` | **AjaxTk / DWT** widget toolkit (`Dwt*` classes), shared by webmail + admin |
| `github.com/Zimbra/zm-web-client` | `webapps/zimbra` | Webmail client (`ZmMail`, JSPs, **the 21 skins**, templates, `zm.css`, `login.css`) |
| `github.com/Zimbra/zm-admin-ajax` | `webapps/zimbraAdmin` | **Admin console** (`ZaApp*`), admin skins (`serenity`, `carbon`, `vami2`), admin CSS |
| `github.com/Zimbra/zm-timezones` | tz data | Calendar timezones |

The **skins** we mirror (e.g. `serenity/skin.properties`) live under
`zm-web-client` / `zm-admin-ajax` in each repo's `WebRoot/skins/<name>/`. That
is the authoritative source for the theme tokens documented in
`openspec/changes/admin-panel-zimbra/zimbra-admin-theme.md`.

## Backend (not a reuse candidate — reference only)

| Repo | Contents |
| --- | --- |
| `github.com/Zimbra/zm-mailbox` | The mail **store**, **SOAP/REST API**, provisioning (LDAP), Lucene index, redolog. This is where `mail_item`, `mailbox`, COS, etc. semantics are defined |
| `github.com/Zimbra/zm-common` | Shared Java utilities |
| `github.com/Zimbra/zm-charset`, `zm-freebusy-provider`, `zm-ews-*` | Peripheral services |

The MySQL schema documented in [01-mysql-databases.md](01-mysql-databases.md)
and the LDAP schema in [04-ldap-structure.md](04-ldap-structure.md) are defined
(and versioned) inside `zm-mailbox`:
- LDAP attribute definitions: `zm-mailbox/store/src/…/zimbra.ldiff` +
  `attrs/zimbra-attrs.xml` (the generator for all `zimbra*` attributes).
- MySQL DDL: `zm-mailbox/store/src/db/…` (`db.sql`, `create_database.sql`).

Read those for authoritative types/constraints; our docs capture the live,
as-deployed shape.

## Documentation / schema references (public)

| Source | Use |
| --- | --- |
| `wiki.zimbra.com` — "Account mailbox database structure" | `mail_item.type` values, folder model |
| `wiki.zimbra.com` — "Directory Schema (ZCS)" | LDAP object classes / attribute reference |
| `zm-mailbox` `.../attrs/zimbra-attrs.xml` | Machine-readable attribute definitions (types, defaults, flags) |

## How to fetch the exact version source (future task)

```bash
# match the installed build version (10.1.17)
git clone --branch <10.1.x-tag> https://github.com/Zimbra/zm-admin-ajax
git clone --branch <10.1.x-tag> https://github.com/Zimbra/zm-web-client
# or start from the community build that produced the VM:
git clone https://github.com/maldua/zimbra-foss
```

Then read `WebRoot/skins/serenity/` for the theme, `WebRoot/js/` for the UI, and
`zm-mailbox/store/src/db` for the schema — **for study, not for copying**.
