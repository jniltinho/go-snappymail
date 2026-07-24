package server

import (
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
	"go-snappymail/internal/config"
	"go-snappymail/internal/handler"
	appMiddleware "go-snappymail/internal/server/middleware"
)

func registerAPIRoutes(g *echo.Group, h *handler.Handlers, authMiddleware, authRateLimit echo.MiddlewareFunc) {
	g.GET("/version", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"version": AppVersion, "app": "go-snappymail"})
	})

	g.GET("/ui/config", h.UI.Config)

	auth := g.Group("/auth")
	auth.POST("/login", h.Auth.DoLogin, authRateLimit)
	auth.POST("/logout", h.Auth.DoLogout)
	auth.GET("/me", h.Auth.Me, authMiddleware)
	auth.GET("/quota", h.Auth.Quota, authMiddleware)

	api := g.Group("", authMiddleware)

	api.GET("/folders", h.Mailbox.FoldersJSON)
	api.POST("/folders", h.Mailbox.CreateSubfolder)
	api.POST("/folders/rename", h.Mailbox.RenameFolder)
	api.POST("/folders/delete", h.Mailbox.DeleteFolder)
	api.GET("/folders/:name/count", h.Mailbox.UnreadCountJSON)

	api.GET("/mail/:mailbox", h.Mailbox.List)
	api.GET("/mail/:mailbox/:uid", h.Message.Read)
	api.GET("/mail/:mailbox/:uid/download", h.Message.Download)
	api.GET("/mail/:mailbox/:uid/raw", h.Message.Raw)
	api.POST("/mail/:mailbox/:uid/flag", h.Message.Flag)
	api.POST("/mail/:mailbox/:uid/move", h.Message.Move)
	api.DELETE("/mail/:mailbox/:uid", h.Message.Delete)
	api.DELETE("/mail/:mailbox", h.Message.EmptyTrash)
	api.GET("/mail/:mailbox/:uid/attachment/:part", h.Message.Attachment)

	api.POST("/compose/send", h.Compose.Send)
	api.POST("/compose/draft", h.Compose.SaveDraft)
	api.POST("/compose/upload", h.Compose.UploadAttachment)

	api.GET("/search", h.Search.Results)
}

func registerRoutes(e *echo.Echo, cfg *config.Config, h *handler.Handlers, distFS fs.FS) {
	authRateLimit := appMiddleware.NewRateLimit(10, time.Minute)
	authMiddleware := appMiddleware.RequireAuth(cfg.Session.Name)

	registerAPIRoutes(e.Group("/api/v1"), h, authMiddleware, authRateLimit)

	e.GET("/*", func(c *echo.Context) error {
		urlPath := c.Request().URL.Path
		// Never serve the SPA for API paths — an unknown API route is a clean
		// 404, not the app shell. Keeps admin API paths from resolving to the
		// webmail SPA on this listener.
		if strings.HasPrefix(urlPath, "/api/") {
			return echo.ErrNotFound
		}
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
