package cmd

import (
	"fmt"

	"go-snappymail/internal/admin"
	"go-snappymail/internal/config"

	"github.com/spf13/cobra"
)

// migrateAdminCmd creates the whole admin (Postfix/Dovecot) mail schema in the
// database configured under [admin.database]. Separate from `migrate`, which
// handles the webmail session database.
var migrateAdminCmd = &cobra.Command{
	Use:   "migrate-admin",
	Short: "Create/update the admin (mail) database schema",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Load()

		if cfg.Admin.Database.Driver == "" || cfg.Admin.Database.DSN == "" {
			return fmt.Errorf("admin.database.driver and admin.database.dsn must be set")
		}

		db, err := admin.Open(cfg.Admin)
		if err != nil {
			return err
		}

		fmt.Println("Creating admin mail schema (domain, mailbox, alias, admin, domain_admins)...")
		if err := admin.Migrate(db); err != nil {
			return err
		}
		fmt.Println("Admin schema ready.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrateAdminCmd)
}
