package admin

import (
	"fmt"
	"time"

	"go-snappymail/internal/config"
	"go-snappymail/internal/database"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Open connects to the mail (Postfix/Dovecot) database described by the admin
// configuration. It is a separate connection from the webmail session database.
//
// Supported drivers: "mysql"/"mariadb", "postgres"/"postgresql", and "sqlite"
// (dev/tests). Returns a ready *gorm.DB or a wrapped error.
func Open(cfg config.AdminConfig) (*gorm.DB, error) {
	logLevel := logger.Error
	if cfg.Database.Debug {
		logLevel = logger.Info
	}
	gormConfig := &gorm.Config{
		Logger: database.NewGormSlogLogger(logLevel, time.Second),
	}

	var db *gorm.DB
	var err error
	switch cfg.Database.Driver {
	case "mariadb", "mysql":
		db, err = gorm.Open(mysql.Open(cfg.Database.DSN), gormConfig)
	case "postgres", "postgresql":
		db, err = gorm.Open(postgres.Open(cfg.Database.DSN), gormConfig)
	case "sqlite", "":
		db, err = gorm.Open(sqlite.Open(cfg.Database.DSN), gormConfig)
	default:
		return nil, fmt.Errorf("admin: unsupported database driver %q", cfg.Database.Driver)
	}
	if err != nil {
		return nil, fmt.Errorf("admin: connect mail database (driver=%s): %w", cfg.Database.Driver, err)
	}
	return db, nil
}

// Migrate creates or updates every admin table (the whole mail schema:
// domain, mailbox, alias, admin, domain_admins). Safe to call repeatedly.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(AllModels()...); err != nil {
		return fmt.Errorf("admin: migrate mail schema: %w", err)
	}
	return nil
}
