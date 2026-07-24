# Zimbra MySQL — GORM Models (reference structs)

Ready-to-use GORM structs for the legacy tables we might read from. These are
**reference only** — copy into `internal/zimbra/` (a future, isolated read-only
package) if/when we build the reporting features in
[05-architecture-reuse.md](05-architecture-reuse.md).

Design notes:
- All timestamps in `mail_item` etc. are **Unix epoch seconds stored as ints**,
  not SQL datetimes — model them as `int64` and convert in code.
- The mailbox tables are **sharded**; a model has no fixed table name. Use
  GORM's `db.Table("mboxgroup1.mail_item")` (or `.Scopes`) to target the shard
  resolved from `zimbra.mailbox.group_id`.
- Treat every legacy table as **read-only**. Never let the admin panel write to
  a live Zimbra store.

## Central `zimbra` database

```go
package zimbra

// Mailbox is zimbra.mailbox — the account -> shard index and usage counters.
type Mailbox struct {
    ID              uint32 `gorm:"column:id;primaryKey"`
    GroupID         uint32 `gorm:"column:group_id"`          // -> mboxgroup<GroupID>
    AccountID       string `gorm:"column:account_id"`        // = LDAP zimbraId
    IndexVolumeID   uint8  `gorm:"column:index_volume_id"`
    ItemIDCheckpoint uint32 `gorm:"column:item_id_checkpoint"`
    ContactCount    *uint32 `gorm:"column:contact_count"`
    SizeCheckpoint  uint64 `gorm:"column:size_checkpoint"`   // mailbox bytes used
    ChangeCheckpoint uint32 `gorm:"column:change_checkpoint"`
    NewMessages     uint32 `gorm:"column:new_messages"`
    LastSoapAccess  uint32 `gorm:"column:last_soap_access"`  // epoch secs
    LastBackupAt    *uint32 `gorm:"column:last_backup_at"`
    LastPurgeAt     uint32 `gorm:"column:last_purge_at"`
    Version         string `gorm:"column:version"`
    Comment         string `gorm:"column:comment"`
}

func (Mailbox) TableName() string { return "zimbra.mailbox" }

// MailboxMetadata is zimbra.mailbox_metadata (per-mailbox KV blobs).
type MailboxMetadata struct {
    MailboxID uint32 `gorm:"column:mailbox_id;primaryKey"`
    Section   string `gorm:"column:section;primaryKey"`
    Metadata  string `gorm:"column:metadata"`
}

func (MailboxMetadata) TableName() string { return "zimbra.mailbox_metadata" }

// Volume is zimbra.volume (blob store volumes).
type Volume struct {
    ID                   uint8  `gorm:"column:id;primaryKey"`
    Type                 int8   `gorm:"column:type"` // 1=primary 2=secondary 10=index
    Name                 string `gorm:"column:name"`
    Path                 string `gorm:"column:path"`
    CompressBlobs        bool   `gorm:"column:compress_blobs"`
    CompressionThreshold int64  `gorm:"column:compression_threshold"`
    StoreType            int8   `gorm:"column:store_type"` // 1=internal 2=external
    StoreManagerClass    string `gorm:"column:store_manager_class"`
}

func (Volume) TableName() string { return "zimbra.volume" }

// Config is zimbra.config (server key/value).
type Config struct {
    Name        string    `gorm:"column:name;primaryKey"`
    Value       string    `gorm:"column:value"`
    Description string    `gorm:"column:description"`
    Modified    time.Time `gorm:"column:modified"`
}

func (Config) TableName() string { return "zimbra.config" }
```

## Sharded `mboxgroup<N>` database

Because the table lives in a shard, don't rely on `TableName()`; pass the shard
explicitly.

```go
package zimbra

// MailItem is mboxgroup<N>.mail_item — the core item row (message, folder,
// contact, appointment, …), discriminated by Type.
type MailItem struct {
    MailboxID   uint32  `gorm:"column:mailbox_id;primaryKey"`
    ID          uint32  `gorm:"column:id;primaryKey"`
    Type        int8    `gorm:"column:type"` // 1=folder 5=message 6=contact 11=appt ...
    ParentID    *uint32 `gorm:"column:parent_id"`
    FolderID    *uint32 `gorm:"column:folder_id"`
    IndexID     *uint32 `gorm:"column:index_id"`
    ImapID      *uint32 `gorm:"column:imap_id"`
    Date        uint32  `gorm:"column:date"` // epoch secs
    Size        uint64  `gorm:"column:size"`
    Locator     string  `gorm:"column:locator"`
    BlobDigest  string  `gorm:"column:blob_digest"`
    Unread      *uint32 `gorm:"column:unread"`
    Flags       int32   `gorm:"column:flags"`
    Tags        int64   `gorm:"column:tags"`
    TagNames    string  `gorm:"column:tag_names"`
    Sender      string  `gorm:"column:sender"`
    Recipients  string  `gorm:"column:recipients"`
    Subject     string  `gorm:"column:subject"`
    Name        string  `gorm:"column:name"`
    Metadata    string  `gorm:"column:metadata"`
    ChangeDate  *uint32 `gorm:"column:change_date"`
    UUID        string  `gorm:"column:uuid"`
}

// Tag is mboxgroup<N>.tag.
type Tag struct {
    MailboxID uint32 `gorm:"column:mailbox_id;primaryKey"`
    ID        int32  `gorm:"column:id;primaryKey"`
    Name      string `gorm:"column:name"`
    Color     *int64 `gorm:"column:color"`
    ItemCount int32  `gorm:"column:item_count"`
    Unread    int32  `gorm:"column:unread"`
    Listed    bool   `gorm:"column:listed"`
    Sequence  uint32 `gorm:"column:sequence"`
}

// Appointment is mboxgroup<N>.appointment (calendar time index).
type Appointment struct {
    MailboxID uint32    `gorm:"column:mailbox_id;primaryKey"`
    UID       string    `gorm:"column:uid;primaryKey"`
    ItemID    uint32    `gorm:"column:item_id"`
    StartTime time.Time `gorm:"column:start_time"`
    EndTime   *time.Time `gorm:"column:end_time"`
}
```

## Example: reading real mailbox usage for one account

```go
// 1) resolve the account's shard + counters from the central DB
var mb zimbra.Mailbox
central.Where("account_id = ?", zimbraID).First(&mb)

// 2) count messages in that mailbox from its shard (type=5 == message)
var msgs int64
central.Table(fmt.Sprintf("mboxgroup%d.mail_item", mb.GroupID)).
    Where("mailbox_id = ? AND type = 5", mb.ID).
    Count(&msgs)

// mb.SizeCheckpoint = bytes used, mb.NewMessages = unread, mb.LastSoapAccess = last activity
```

> `zimbraID` comes from LDAP (`zimbraId` on the account) — see
> [04-ldap-structure.md](04-ldap-structure.md). This is the join between the two
> stores.
