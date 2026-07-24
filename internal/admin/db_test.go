package admin

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"go-snappymail/internal/config"

	"gorm.io/gorm"
)

// newTestDB opens an isolated in-memory SQLite database with the full admin
// schema migrated. Each call gets a fresh database.
func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	// Unique in-memory database per test → real isolation (no shared cache leak).
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := Open(config.AdminConfig{
		Database: config.DatabaseConfig{Driver: "sqlite", DSN: dsn},
	})
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if err := Migrate(db); err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}
	t.Cleanup(func() {
		// Drop every table so the shared in-memory DB is clean for the next test.
		for _, m := range AllModels() {
			_ = db.Migrator().DropTable(m)
		}
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
	})
	return db
}

func TestMigrateCreatesAllTables(t *testing.T) {
	db := newTestDB(t)

	wantTables := []string{"domain", "mailbox", "alias", "admin", "domain_admins"}
	for _, tbl := range wantTables {
		if !db.Migrator().HasTable(tbl) {
			t.Errorf("table %q not created by Migrate", tbl)
		}
	}
}

func TestModelsTableNames(t *testing.T) {
	tests := []struct {
		name  string
		model interface{ TableName() string }
		want  string
	}{
		{"domain", Domain{}, "domain"},
		{"mailbox", Mailbox{}, "mailbox"},
		{"alias", Alias{}, "alias"},
		{"admin", Admin{}, "admin"},
		{"domain_admins", DomainAdmin{}, "domain_admins"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.model.TableName(); got != tt.want {
				t.Errorf("TableName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestModelsCRUDRoundTrip(t *testing.T) {
	db := newTestDB(t)

	// Create a domain, a mailbox and an alias under it.
	dom := &Domain{Domain: "example.com", Description: "Example", Active: true}
	if err := db.Create(dom).Error; err != nil {
		t.Fatalf("create domain: %v", err)
	}
	mb := &Mailbox{Username: "user@example.com", Password: "hash", Domain: "example.com", LocalPart: "user", Active: true}
	if err := db.Create(mb).Error; err != nil {
		t.Fatalf("create mailbox: %v", err)
	}
	al := &Alias{Address: "info@example.com", Goto: "user@example.com", Domain: "example.com", Active: true}
	if err := db.Create(al).Error; err != nil {
		t.Fatalf("create alias: %v", err)
	}

	// Read back.
	var got Domain
	if err := db.First(&got, "domain = ?", "example.com").Error; err != nil {
		t.Fatalf("read domain: %v", err)
	}
	if got.Description != "Example" {
		t.Errorf("domain description = %q, want Example", got.Description)
	}

	// Count relationships.
	var mbCount, alCount int64
	db.Model(&Mailbox{}).Where("domain = ?", "example.com").Count(&mbCount)
	db.Model(&Alias{}).Where("domain = ?", "example.com").Count(&alCount)
	if mbCount != 1 || alCount != 1 {
		t.Errorf("counts mb=%d al=%d, want 1/1", mbCount, alCount)
	}

	// Update + delete.
	if err := db.Model(&got).Update("active", false).Error; err != nil {
		t.Fatalf("update domain: %v", err)
	}
	if err := db.Delete(&Alias{}, "address = ?", "info@example.com").Error; err != nil {
		t.Fatalf("delete alias: %v", err)
	}
	db.Model(&Alias{}).Count(&alCount)
	if alCount != 0 {
		t.Errorf("alias count after delete = %d, want 0", alCount)
	}
}

func TestTimestampsAutoManaged(t *testing.T) {
	db := newTestDB(t)

	dom := &Domain{Domain: "ts.example", Active: true}
	if err := db.Create(dom).Error; err != nil {
		t.Fatalf("create: %v", err)
	}
	if dom.Created.IsZero() || dom.Modified.IsZero() {
		t.Fatalf("autoCreateTime/autoUpdateTime not set: created=%v modified=%v", dom.Created, dom.Modified)
	}
	firstModified := dom.Modified

	time.Sleep(1100 * time.Millisecond) // sqlite autoUpdateTime has 1s resolution
	if err := db.Model(dom).Update("description", "changed").Error; err != nil {
		t.Fatalf("update: %v", err)
	}
	var reloaded Domain
	db.First(&reloaded, "domain = ?", "ts.example")
	if !reloaded.Modified.After(firstModified) {
		t.Errorf("Modified not bumped on update: was %v, now %v", firstModified, reloaded.Modified)
	}
}

func TestPasswordNotSerialized(t *testing.T) {
	mb := Mailbox{Username: "u@x", Password: "SECRET-HASH"}
	b, err := json.Marshal(mb)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(b), "SECRET-HASH") || strings.Contains(string(b), "password") {
		t.Errorf("password leaked in JSON: %s", b)
	}
}

func TestOpenUnsupportedDriver(t *testing.T) {
	_, err := Open(config.AdminConfig{Database: config.DatabaseConfig{Driver: "oracle", DSN: "x"}})
	if err == nil {
		t.Fatal("Open() with unsupported driver = nil error, want error")
	}
}
