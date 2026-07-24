// Package admin implements the ZimbraAdmin-style admin panel backend: GORM
// models over the Postfix/Dovecot mail database, JWT/RBAC auth, and the
// /api/v1/admin/* handlers. It is fully isolated from the webmail: its own
// database connection, its own routes, and its own listener.
//
// The schema mirrors PostfixAdmin (domains, mailboxes, aliases, admins,
// domain_admins), adapted here — go-postfixadmin is the reference, not a
// dependency.
package admin

import "time"

// Domain is a virtual mail domain.
type Domain struct {
	Domain      string    `gorm:"column:domain;type:varchar(255);primaryKey" json:"domain"`
	Description string    `gorm:"column:description;type:varchar(255);not null;default:''" json:"description"`
	Aliases     int       `gorm:"column:aliases;not null;default:0" json:"aliases"`
	Mailboxes   int       `gorm:"column:mailboxes;not null;default:0" json:"mailboxes"`
	MaxQuota    int64     `gorm:"column:maxquota;not null;default:0" json:"maxquota"`
	Transport   string    `gorm:"column:transport;type:varchar(255);not null;default:'virtual'" json:"transport"`
	BackupMX    bool      `gorm:"column:backupmx;not null;default:false" json:"backupmx"`
	Active      bool      `gorm:"column:active;not null;default:true" json:"active"`
	Created     time.Time `gorm:"column:created;autoCreateTime;not null;default:'2000-01-01 00:00:00'" json:"created"`
	Modified    time.Time `gorm:"column:modified;autoUpdateTime;not null;default:'2000-01-01 00:00:00'" json:"modified"`
}

// TableName maps Domain to the PostfixAdmin "domain" table.
func (Domain) TableName() string { return "domain" }

// Mailbox is a virtual mailbox (account) belonging to a Domain.
type Mailbox struct {
	Username  string    `gorm:"column:username;type:varchar(255);primaryKey" json:"username"`
	Password  string    `gorm:"column:password;type:varchar(255);not null" json:"-"`
	Name      string    `gorm:"column:name;type:varchar(255);not null;default:''" json:"name"`
	Maildir   string    `gorm:"column:maildir;type:varchar(255);not null;default:''" json:"maildir"`
	Quota     int64     `gorm:"column:quota;not null;default:0" json:"quota"`
	LocalPart string    `gorm:"column:local_part;type:varchar(255);not null;default:''" json:"local_part"`
	Domain    string    `gorm:"column:domain;type:varchar(255);not null;index:idx_mailbox_domain,priority:1" json:"domain"`
	Active    bool      `gorm:"column:active;not null;default:true" json:"active"`
	Created   time.Time `gorm:"column:created;autoCreateTime;not null;default:'2000-01-01 00:00:00'" json:"created"`
	Modified  time.Time `gorm:"column:modified;autoUpdateTime;not null;default:'2000-01-01 00:00:00'" json:"modified"`
}

// TableName maps Mailbox to the PostfixAdmin "mailbox" table.
func (Mailbox) TableName() string { return "mailbox" }

// Alias forwards an address to one or more destinations (Goto, comma-separated).
type Alias struct {
	Address  string    `gorm:"column:address;type:varchar(255);primaryKey" json:"address"`
	Goto     string    `gorm:"column:goto;type:text;not null" json:"goto"`
	Domain   string    `gorm:"column:domain;type:varchar(255);not null;index:idx_alias_domain,priority:1" json:"domain"`
	Active   bool      `gorm:"column:active;not null;default:true" json:"active"`
	Created  time.Time `gorm:"column:created;autoCreateTime;not null;default:'2000-01-01 00:00:00'" json:"created"`
	Modified time.Time `gorm:"column:modified;autoUpdateTime;not null;default:'2000-01-01 00:00:00'" json:"modified"`
}

// TableName maps Alias to the PostfixAdmin "alias" table.
func (Alias) TableName() string { return "alias" }

// Admin is a panel administrator. Superadmin grants full access; otherwise
// scope is limited to the domains linked via DomainAdmin (domain_admin role).
type Admin struct {
	Username   string    `gorm:"column:username;type:varchar(255);primaryKey" json:"username"`
	Password   string    `gorm:"column:password;type:varchar(255);not null" json:"-"`
	Superadmin bool      `gorm:"column:superadmin;not null;default:false" json:"superadmin"`
	Active     bool      `gorm:"column:active;not null;default:true" json:"active"`
	Created    time.Time `gorm:"column:created;autoCreateTime;not null;default:'2000-01-01 00:00:00'" json:"created"`
	Modified   time.Time `gorm:"column:modified;autoUpdateTime;not null;default:'2000-01-01 00:00:00'" json:"modified"`
}

// TableName maps Admin to the PostfixAdmin "admin" table.
func (Admin) TableName() string { return "admin" }

// DomainAdmin links a non-superadmin Admin to a Domain it may manage. It backs
// the domain_admin RBAC scope.
type DomainAdmin struct {
	Username string    `gorm:"column:username;type:varchar(255);primaryKey" json:"username"`
	Domain   string    `gorm:"column:domain;type:varchar(255);primaryKey;index:idx_domain_admins_domain" json:"domain"`
	Created  time.Time `gorm:"column:created;autoCreateTime;not null;default:'2000-01-01 00:00:00'" json:"created"`
	Active   bool      `gorm:"column:active;not null;default:true" json:"active"`
}

// TableName maps DomainAdmin to the PostfixAdmin "domain_admins" table.
func (DomainAdmin) TableName() string { return "domain_admins" }

// AllModels returns every admin model, in dependency order, for AutoMigrate.
func AllModels() []any {
	return []any{
		&Domain{},
		&Mailbox{},
		&Alias{},
		&Admin{},
		&DomainAdmin{},
	}
}
