// Package cmd defines Cobra CLI commands for go-snappymail.
package cmd

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	Version   = "dev"
	BuildDate = "unknown"
	GitCommit = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "go-snappymail",
	Short: "Go SnappyMail — SnappyMail UX webmail in Go",
	Long: `go-snappymail is a self-hosted webmail with SnappyMail-like UI,
backed by IMAP/SMTP and a Vue 3 frontend embedded in a single binary.

Commands:
  init      Create a default configuration file
  serve     Start the web server
  migrate   Run database migrations
  version   Show version information`,
}

var globalFS embed.FS

// Execute runs the root command with embedded static assets.
func Execute(fs embed.FS) {
	globalFS = fs
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "path to config.toml")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/go-snappymail/")
	}

	viper.SetEnvPrefix("GOSM")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "error reading config: %v\n", err)
		}
	}
}
