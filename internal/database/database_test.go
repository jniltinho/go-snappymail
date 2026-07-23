package database

import (
	"testing"

	"github.com/jniltinho/go-snappymail/internal/config"
	"github.com/jniltinho/go-snappymail/internal/model"
)

func TestOpenSQLite(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Driver: "sqlite",
			DSN:    "file::memory:?cache=shared",
			Debug:  false,
		},
	}

	db, err := Open(cfg)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	if err := db.AutoMigrate(&model.Session{}); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}
}
