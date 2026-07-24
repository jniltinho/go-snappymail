# Zimbra MySQL — Column Reference (key tables)

Column-level reference for the tables most relevant to a provisioning/reporting
panel. Types are as reported by `information_schema` on the reference install.

## `zimbra.mailbox` — account → mailbox index

The bridge between the LDAP directory and the sharded mailbox data.

| Column | Type | Notes |
| --- | --- | --- |
| `id` | int unsigned | **PK** — the `mailbox_id` used inside `mboxgroup<N>` |
| `group_id` | int unsigned | Shard number → `mboxgroup<group_id>` |
| `account_id` | varchar(127) | **UNIQUE** — the LDAP `zimbraId` of the account |
| `index_volume_id` | tinyint unsigned | Search-index volume |
| `item_id_checkpoint` | int unsigned | Next item id |
| `contact_count` | int unsigned | Cached contact count |
| `size_checkpoint` | bigint unsigned | **Mailbox size in bytes** (quota usage) |
| `change_checkpoint` | int unsigned | Change sequence |
| `new_messages` | int unsigned | Unread/new counter |
| `last_soap_access` | int unsigned | Last activity (epoch secs) |
| `last_backup_at` / `last_purge_at` | int unsigned | Housekeeping timestamps |
| `version` | varchar(16) | Mailbox schema version |
| `comment` | varchar(255) | Free text |

> Reporting value: `size_checkpoint` (quota used) + `new_messages` +
> `last_soap_access` (last activity) are exactly the columns PostfixAdmin lacks.

## `zimbra.mailbox_metadata` — per-mailbox blobs

| Column | Type | Notes |
| --- | --- | --- |
| `mailbox_id` | int unsigned | **PK** part |
| `section` | varchar(64) | **PK** part — e.g. widget/state section name |
| `metadata` | mediumtext | Serialized blob |

## `zimbra.volume` — blob store volumes

| Column | Type | Notes |
| --- | --- | --- |
| `id` | tinyint unsigned | **PK** |
| `type` | tinyint | 1=primary msg, 2=secondary, 10=index |
| `name` | varchar(255) | UNIQUE |
| `path` | text | UNIQUE — on-disk root |
| `compress_blobs` | tinyint(1) | Whether blobs are gzip'd |
| `compression_threshold` | bigint | Min size to compress |
| `store_type` | tinyint(1) | 1=internal (disk), 2=external (S3/…) |
| `store_manager_class` | varchar(255) | Pluggable store impl |

## `zimbra.config` — server config KV

| Column | Type | Notes |
| --- | --- | --- |
| `name` | varchar(255) | **PK** |
| `value` | text | |
| `description` | text | |
| `modified` | timestamp | |

## `mboxgroup<N>.mail_item` — the core item table

| Column | Type | Notes |
| --- | --- | --- |
| `mailbox_id` | int unsigned | **PK** part |
| `id` | int unsigned | **PK** part — item id within the mailbox |
| `type` | tinyint | Discriminator (see 01-mysql-databases.md) |
| `parent_id` | int unsigned | Conversation / parent item |
| `folder_id` | int unsigned | Containing folder (self-ref to a `type=1` item) |
| `prev_folders` | text | Move history |
| `index_id` | int unsigned | Search-index id |
| `imap_id` | int unsigned | IMAP UID |
| `date` | int unsigned | Received/created (epoch secs) |
| `size` | bigint unsigned | Item size in bytes |
| `locator` | varchar(1024) | Blob locator (volume + path) |
| `blob_digest` | varchar(44) | Content hash (dedupe) |
| `unread` | int unsigned | Unread flag/count |
| `flags` | int | Bitmask (flagged, replied, draft…) |
| `tags` | bigint | Tag bitmask (legacy) |
| `tag_names` | text | Tag names (current) |
| `sender` | varchar(128) | From (messages) |
| `recipients` | varchar(128) | To |
| `subject` | text | Subject / item name |
| `name` | varchar(255) | Folder/contact/file name |
| `metadata` | mediumtext | Serialized item metadata |
| `mod_metadata` / `mod_content` | int unsigned | Change stamps |
| `change_date` | int unsigned | Last change (epoch secs) |
| `uuid` | varchar(127) | Stable UUID |

## `mboxgroup<N>.tag` and `tagged_item`

`tag`: `mailbox_id`,`id` (PK), `name` varchar(128), `color` bigint,
`item_count`, `unread`, `listed` tinyint(1), `sequence`, `policy` varchar(1024).

`tagged_item`: join of `mail_item` ↔ `tag` (by `mailbox_id` + item/tag ids).

## `mboxgroup<N>.imap_folder` — external IMAP sync

| Column | Type |
| --- | --- |
| `mailbox_id` | int unsigned (PK) |
| `item_id` | int unsigned (PK) |
| `data_source_id` | char(36) |
| `local_path` / `remote_path` | varchar(1000) |
| `uid_validity` | int unsigned |

## `mboxgroup<N>.appointment` — calendar time index

| Column | Type |
| --- | --- |
| `mailbox_id` | int unsigned (PK) |
| `uid` | varchar(255) (PK) |
| `item_id` | int unsigned → `mail_item.id` |
| `start_time` | datetime |
| `end_time` | datetime |

## `mboxgroup<N>.revision` — item versions

`mailbox_id`,`item_id`,`version` (PK), `date`, `size`, `locator`,
`blob_digest`, `name`, `metadata`, `mod_metadata`, `change_date`, `mod_content`.
