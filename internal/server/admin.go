package server

import (
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"go-snappymail/internal/admin"
	"go-snappymail/internal/config"
	appMiddleware "go-snappymail/internal/server/middleware"

	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"
)

// buildAdminServer wires the ISOLATED admin listener. It uses its OWN Echo
// instance — the admin API (/api/v1/admin/*) and the admin SPA (web/admin-dist)
// exist only here, never on the webmail router. This is the surface-isolation
// guarantee: a request for an admin route on the webmail port cannot resolve.
func buildAdminServer(cfg *config.Config, adminDB *gorm.DB, embeddedFiles embed.FS) (*http.Server, error) {
	adminFS, err := fs.Sub(embeddedFiles, "web/admin-dist")
	if err != nil {
		return nil, fmt.Errorf("open embedded admin dist: %w", err)
	}

	e := echo.New()
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.RequestID())
	// Defensive HTTP headers (X-Frame-Options, CSP, frame-ancestors none) on the
	// admin listener too — the webmail port sets these and the panel must not be
	// left more exposed to clickjacking/XSS.
	e.Use(appMiddleware.SecurityHeaders())
	// Admin auth is a stateless JWT Bearer flow (no session cookie), so cookie
	// CSRF is not applicable; the JWT middleware guards every protected route.
	e.Use(echoMiddleware.RequestLoggerWithConfig(echoMiddleware.RequestLoggerConfig{
		LogMethod: true, LogURI: true, LogStatus: true, LogLatency: true, LogRemoteIP: true, LogRequestID: true,
		LogValuesFunc: func(c *echo.Context, v echoMiddleware.RequestLoggerValues) error {
			slog.Info("admin request",
				"method", v.Method, "uri", v.URI, "status", v.Status,
				"latency_ms", v.Latency.Milliseconds(), "remote_ip", v.RemoteIP, "request_id", v.RequestID)
			return nil
		},
	}))

	h := admin.NewHandlers(adminDB, cfg.Admin)
	h.RegisterRoutes(e.Group("/api/v1/admin"))
	registerAdminSPA(e, adminFS)

	srv := &http.Server{
		Addr:              cfg.Admin.AdminAddr(),
		Handler:           e,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      0, // JSON responses are small; 0 avoids cutting slow clients on a shared listener
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}
	return srv, nil
}

// adminBase is the SPA base path (matches vite base in frontend-admin). Assets
// are requested as /admin/assets/*, so the base is stripped before the embedded
// file lookup (which is rooted at web/admin-dist).
const adminBase = "/admin/"

// registerAdminSPA serves the admin SPA (web/admin-dist) with client-side
// routing fallback, scoped to THIS listener only. It never falls back to the
// webmail SPA (web/dist).
func registerAdminSPA(e *echo.Echo, adminFS fs.FS) {
	serveIndex := func(c *echo.Context) error {
		data, err := fs.ReadFile(adminFS, "index.html")
		if err != nil {
			return echo.ErrNotFound
		}
		return c.HTML(http.StatusOK, string(data))
	}

	// Bare host or "/admin" → the app lives under /admin/.
	e.GET("/", func(c *echo.Context) error {
		return c.Redirect(http.StatusFound, adminBase)
	})
	e.GET("/admin", func(c *echo.Context) error {
		return c.Redirect(http.StatusFound, adminBase)
	})

	e.GET("/*", func(c *echo.Context) error {
		urlPath := c.Request().URL.Path
		// API paths never fall through to the admin SPA (clean 404).
		if strings.HasPrefix(urlPath, "/api/") {
			return echo.ErrNotFound
		}
		// Map the request to a path relative to the admin dist root by stripping
		// the SPA base prefix (/admin/assets/x.js → assets/x.js).
		rel := strings.TrimPrefix(strings.TrimPrefix(urlPath, adminBase), "/")
		if ext := strings.ToLower(filepath.Ext(rel)); ext != "" && rel != "" {
			data, err := fs.ReadFile(adminFS, rel)
			if err != nil {
				return echo.ErrNotFound
			}
			ct := mime.TypeByExtension(ext)
			if ct == "" {
				ct = "application/octet-stream"
			}
			return c.Blob(http.StatusOK, ct, data)
		}
		// Any non-asset path (/admin/, /admin/login, …) → the SPA shell.
		return serveIndex(c)
	})
}
