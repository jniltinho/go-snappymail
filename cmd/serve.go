package cmd

// serveCmd starts the HTTP server. It opens the configured database, initialises
// the handler and session layers, and hands off to server.Start.
// Port and debug-mode flags can override the config file values.

import (
	"fmt"
	"log/slog"
	"os"

	"go-snappymail/internal/config"
	"go-snappymail/internal/database"
	"go-snappymail/internal/server"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web server",
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})))

		cfg := config.Load()

		db, err := database.Open(cfg)
		if err != nil {
			return fmt.Errorf("open database: %w", err)
		}

		server.AppVersion = Version
		return server.Start(cfg, db, globalFS)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().String("port", "", "HTTP port (overrides config)")
	serveCmd.Flags().Bool("debug", false, "Echo debug mode")
	_ = viper.BindPFlag("server.port", serveCmd.Flags().Lookup("port"))
	_ = viper.BindPFlag("server.debug", serveCmd.Flags().Lookup("debug"))
}
