package server

import (
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/jniltinho/go-snappymail/internal/config"
	"github.com/jniltinho/go-snappymail/internal/handler"
	appMiddleware "github.com/jniltinho/go-snappymail/internal/server/middleware"
	"github.com/labstack/echo/v5"
)

func registerAPIRoutes(g *echo.Group, h *handler.Handlers, authMiddleware, authRateLimit echo.MiddlewareFunc) {
	g.GET("/version", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"version": AppVersion, "app": "go-snappymail"})
	})

	auth := g.Group("/auth")
	auth.POST("/login", h.Auth.DoLogin, authRateLimit)
	auth.POST("/logout", h.Auth.DoLogout)
	auth.GET("/me", h.Auth.Me, authMiddleware)
}

func registerRoutes(e *echo.Echo, cfg *config.Config, h *handler.Handlers, distFS fs.FS) {
	authRateLimit := appMiddleware.NewRateLimit(10, time.Minute)
	authMiddleware := appMiddleware.RequireAuth(cfg.Session.Name)

	registerAPIRoutes(e.Group("/api/v1"), h, authMiddleware, authRateLimit)

	e.GET("/*", func(c *echo.Context) error {
		urlPath := c.Request().URL.Path
		ext := strings.ToLower(filepath.Ext(urlPath))
		if ext != "" {
			fsPath := strings.TrimPrefix(urlPath, "/")
			data, err := fs.ReadFile(distFS, fsPath)
			if err != nil {
				return echo.ErrNotFound
			}
			ct := mime.TypeByExtension(ext)
			if ct == "" {
				ct = "application/octet-stream"
			}
			return c.Blob(http.StatusOK, ct, data)
		}
		data, err := fs.ReadFile(distFS, "index.html")
		if err != nil {
			return echo.ErrNotFound
		}
		return c.HTML(http.StatusOK, string(data))
	})
}
