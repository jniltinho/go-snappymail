// Package server initialises the Echo HTTP server for go-snappymail.
package server

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-snappymail/internal/admin"
	"go-snappymail/internal/config"
	"go-snappymail/internal/handler"
	"go-snappymail/internal/model"
	appMiddleware "go-snappymail/internal/server/middleware"
	"go-snappymail/internal/session"

	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

// AppVersion is set via ldflags at build time.
var AppVersion = "dev"

// Start wires middleware, routes, and blocks until shutdown.
func Start(cfg *config.Config, db *gorm.DB, embeddedFiles embed.FS) error {
	if len(cfg.Server.SecretKey) < 32 {
		return fmt.Errorf("server.secret_key deve ter pelo menos 32 bytes (atual: %d); gere uma chave com: openssl rand -hex 32", len(cfg.Server.SecretKey))
	}

	if err := db.AutoMigrate(&model.Session{}); err != nil {
		return fmt.Errorf("auto-migrate sessions: %w", err)
	}
	session.InitDB(db)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	e := echo.New()

	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.RequestID())
	e.Use(appMiddleware.SecurityHeaders())
	e.Use(appMiddleware.CSRF(cfg.Session.Secure))
	e.Use(echoMiddleware.RequestLoggerWithConfig(echoMiddleware.RequestLoggerConfig{
		LogMethod:    true,
		LogURI:       true,
		LogStatus:    true,
		LogLatency:   true,
		LogRemoteIP:  true,
		LogRequestID: true,
		LogValuesFunc: func(c *echo.Context, v echoMiddleware.RequestLoggerValues) error {
			slog.Info("request",
				"method", v.Method,
				"uri", v.URI,
				"status", v.Status,
				"latency_ms", v.Latency.Milliseconds(),
				"remote_ip", v.RemoteIP,
				"request_id", v.RequestID,
			)
			return nil
		},
	}))

	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				maxIdle := 30 * time.Minute
				if cfg.Session.MaxAge > 0 {
					maxIdle = time.Duration(cfg.Session.MaxAge) * time.Second
				}
				session.Cleanup(maxIdle)
			}
		}
	}()

	distFS, err := fs.Sub(embeddedFiles, "web/dist")
	if err != nil {
		return fmt.Errorf("open embedded dist: %w", err)
	}

	h := handler.New(cfg)
	registerRoutes(e, cfg, h, distFS)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           e,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      0,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Optionally build the ISOLATED admin listener (separate Echo, own DB).
	var adminSrv *http.Server
	if cfg.Admin.Enabled {
		if err := cfg.Admin.Validate(); err != nil {
			return fmt.Errorf("admin config: %w", err)
		}
		adminDB, err := admin.Open(cfg.Admin)
		if err != nil {
			return err
		}
		// Schema changes are an explicit, operator-run step (`migrate-admin`), not
		// a startup side effect: the admin models map onto the shared Postfix/
		// Dovecot database, and silent AutoMigrate there could alter production
		// tables. serve only opens the DB and serves.
		adminSrv, err = buildAdminServer(cfg, adminDB, embeddedFiles)
		if err != nil {
			return err
		}
	}

	// Coordinated graceful shutdown of every listener on signal.
	go func() {
		<-ctx.Done()
		slog.Info("shutting down servers")
		shutCtx, shutCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutCancel()
		_ = srv.Shutdown(shutCtx)
		if adminSrv != nil {
			_ = adminSrv.Shutdown(shutCtx)
		}
	}()

	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		slog.Info("go-snappymail webmail listening", "addr", addr, "version", AppVersion)
		return serveHTTP(srv, cfg.Server.TLSCert, cfg.Server.TLSKey)
	})
	if adminSrv != nil {
		g.Go(func() error {
			slog.Info("go-snappymail admin listening", "addr", adminSrv.Addr, "skin", cfg.Admin.Skin)
			return serveHTTP(adminSrv, cfg.Admin.TLSCert, cfg.Admin.TLSKey)
		})
	}
	if err := g.Wait(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// serveHTTP starts a server with optional TLS and treats ErrServerClosed as a
// clean shutdown.
func serveHTTP(srv *http.Server, certFile, keyFile string) error {
	var err error
	if certFile != "" && keyFile != "" {
		err = srv.ListenAndServeTLS(certFile, keyFile)
	} else {
		err = srv.ListenAndServe()
	}
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}
