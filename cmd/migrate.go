package cmd

import (
	"fmt"

	"go-snappymail/internal/config"
	"go-snappymail/internal/database"
	"go-snappymail/internal/model"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Load()

		db, err := database.Open(cfg)
		if err != nil {
			return fmt.Errorf("connect database: %w", err)
		}

		fmt.Println("Running migrations...")
		if err := db.AutoMigrate(&model.Session{}); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
		fmt.Println("Migrations completed.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
