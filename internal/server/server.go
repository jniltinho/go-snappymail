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

	"github.com/jniltinho/go-snappymail/internal/config"
	"github.com/jniltinho/go-snappymail/internal/handler"
	appMiddleware "github.com/jniltinho/go-snappymail/internal/server/middleware"
	"github.com/jniltinho/go-snappymail/internal/session"
	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"
)

// AppVersion is set via ldflags at build time.
var AppVersion = "dev"

// Start wires middleware, routes, and blocks until shutdown.
func Start(cfg *config.Config, db *gorm.DB, embeddedFiles embed.FS) error {
	session.InitDB(db)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	e := echo.New()

	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.RequestID())
	e.Use(appMiddleware.SecurityHeaders())
	e.Use(appMiddleware.CSRF())
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

	go func() {
		<-ctx.Done()
		slog.Info("shutting down server")
		shutCtx, shutCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutCancel()
		_ = srv.Shutdown(shutCtx)
	}()

	slog.Info("go-snappymail listening", "addr", addr, "version", AppVersion)
	var listenErr error
	if cfg.Server.TLSCert != "" && cfg.Server.TLSKey != "" {
		listenErr = srv.ListenAndServeTLS(cfg.Server.TLSCert, cfg.Server.TLSKey)
	} else {
		listenErr = srv.ListenAndServe()
	}
	if listenErr == http.ErrServerClosed {
		return nil
	}
	return listenErr
}
