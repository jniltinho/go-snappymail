package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestAdminConfigLoadDefaults(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	// Only enable admin; leave everything else unset to exercise defaults.
	viper.Set("admin.enabled", true)

	cfg := Load()

	if !cfg.Admin.Enabled {
		t.Fatal("Admin.Enabled = false, want true")
	}
	if cfg.Admin.Port != 7071 {
		t.Errorf("Admin.Port = %d, want default 7071", cfg.Admin.Port)
	}
	if cfg.Admin.Skin != "serenity" {
		t.Errorf("Admin.Skin = %q, want default %q", cfg.Admin.Skin, "serenity")
	}
	if cfg.Admin.JWTMaxAgeSec != 3600 {
		t.Errorf("Admin.JWTMaxAgeSec = %d, want default 3600", cfg.Admin.JWTMaxAgeSec)
	}
}

func TestAdminConfigLoadExplicit(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	viper.Set("admin.enabled", true)
	viper.Set("admin.host", "127.0.0.1")
	viper.Set("admin.port", 9090)
	viper.Set("admin.tls", true)
	viper.Set("admin.skin", "carbon")
	viper.Set("admin.jwt_secret", "s3cr3t")
	viper.Set("admin.jwt_max_age_sec", 1800)
	viper.Set("admin.database.driver", "postgres")
	viper.Set("admin.database.dsn", "postgres://u:p@localhost/mail")

	cfg := Load()

	if cfg.Admin.Host != "127.0.0.1" || cfg.Admin.Port != 9090 {
		t.Errorf("host:port = %s:%d, want 127.0.0.1:9090", cfg.Admin.Host, cfg.Admin.Port)
	}
	if !cfg.Admin.TLS {
		t.Error("Admin.TLS = false, want true")
	}
	if cfg.Admin.Skin != "carbon" {
		t.Errorf("Admin.Skin = %q, want carbon", cfg.Admin.Skin)
	}
	if cfg.Admin.JWTSecret != "s3cr3t" || cfg.Admin.JWTMaxAgeSec != 1800 {
		t.Errorf("jwt = %q/%d, want s3cr3t/1800", cfg.Admin.JWTSecret, cfg.Admin.JWTMaxAgeSec)
	}
	if cfg.Admin.Database.Driver != "postgres" || cfg.Admin.Database.DSN == "" {
		t.Errorf("admin db = %q/%q, want postgres/<dsn>", cfg.Admin.Database.Driver, cfg.Admin.Database.DSN)
	}
}

func TestAdminConfigValidate(t *testing.T) {
	valid := AdminConfig{
		Enabled:   true,
		Port:      7071,
		JWTSecret: "secret",
		Database:  DatabaseConfig{Driver: "mysql", DSN: "u:p@/mail"},
	}

	tests := []struct {
		name    string
		mutate  func(*AdminConfig)
		wantErr bool
	}{
		{"disabled skips all checks", func(a *AdminConfig) { *a = AdminConfig{Enabled: false} }, false},
		{"valid", func(a *AdminConfig) {}, false},
		{"missing jwt secret", func(a *AdminConfig) { a.JWTSecret = "" }, true},
		{"missing db driver", func(a *AdminConfig) { a.Database.Driver = "" }, true},
		{"missing db dsn", func(a *AdminConfig) { a.Database.DSN = "" }, true},
		{"port out of range low", func(a *AdminConfig) { a.Port = 0 }, true},
		{"port out of range high", func(a *AdminConfig) { a.Port = 70000 }, true},
		{"tls without cert", func(a *AdminConfig) { a.TLS = true }, true},
		{"tls with cert/key", func(a *AdminConfig) { a.TLS = true; a.TLSCert = "c"; a.TLSKey = "k" }, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := valid
			tt.mutate(&c)
			err := c.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAdminAddr(t *testing.T) {
	tests := []struct {
		name string
		cfg  AdminConfig
		want string
	}{
		{"default port", AdminConfig{Host: ""}, ":7071"},
		{"local bind", AdminConfig{Host: "127.0.0.1", Port: 7071}, "127.0.0.1:7071"},
		{"custom port", AdminConfig{Host: "0.0.0.0", Port: 9090}, "0.0.0.0:9090"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.AdminAddr(); got != tt.want {
				t.Errorf("AdminAddr() = %q, want %q", got, tt.want)
			}
		})
	}
}
