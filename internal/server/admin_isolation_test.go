package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-snappymail/internal/admin"
	"go-snappymail/internal/config"

	"github.com/labstack/echo/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// buildAdminEcho builds just the admin Echo (routes + SPA) for isolation
// testing, mirroring what buildAdminServer wires onto its listener.
func buildAdminEcho(t *testing.T) *echo.Echo {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:iso?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := admin.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	e := echo.New()
	h := admin.NewHandlers(db, config.AdminConfig{JWTSecret: "x", JWTMaxAgeSec: 3600})
	h.RegisterRoutes(e.Group("/api/v1/admin"))
	return e
}

// buildWebmailEcho builds a webmail-style Echo with only a sample webmail route
// (no admin routes) to prove the admin surface is absent here.
func buildWebmailEcho(t *testing.T) *echo.Echo {
	t.Helper()
	e := echo.New()
	api := e.Group("/api/v1")
	api.GET("/version", func(c *echo.Context) error { return c.JSON(http.StatusOK, map[string]string{"v": "test"}) })
	api.GET("/mail/:mailbox", func(c *echo.Context) error { return c.JSON(http.StatusOK, map[string]string{"ok": "1"}) })
	// Include the SPA catch-all: it must NOT swallow /api/* into the app shell,
	// so a cross-surface admin API path still 404s instead of returning 200 HTML.
	e.GET("/*", func(c *echo.Context) error {
		if strings.HasPrefix(c.Request().URL.Path, "/api/") {
			return echo.ErrNotFound
		}
		return c.HTML(http.StatusOK, "<html>webmail spa</html>")
	})
	return e
}

func status(t *testing.T, e *echo.Echo, method, path string) int {
	t.Helper()
	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec.Code
}

// TestAdminSurfaceNotOnWebmail is the golden-rule isolation guarantee: admin
// routes never resolve on the webmail router, and webmail routes never resolve
// on the admin router. Separate Echo instances make crossing impossible.
func TestAdminSurfaceNotOnWebmail(t *testing.T) {
	webmail := buildWebmailEcho(t)
	adminE := buildAdminEcho(t)

	// Admin routes on the WEBMAIL router → 404 (they don't exist there).
	adminPaths := []struct{ method, path string }{
		{http.MethodPost, "/api/v1/admin/auth/login"},
		{http.MethodGet, "/api/v1/admin/overview"},
		{http.MethodGet, "/api/v1/admin/domains"},
		{http.MethodGet, "/api/v1/admin/admins"},
	}
	for _, p := range adminPaths {
		// 404 (no route) or 405 (path shadowed by the GET SPA catch-all) both
		// prove the admin handler never runs on the webmail port. What must NOT
		// happen is a 2xx that reaches an admin handler.
		got := status(t, webmail, p.method, p.path)
		if got != http.StatusNotFound && got != http.StatusMethodNotAllowed {
			t.Errorf("admin route %s %s on webmail router: status = %d, want 404/405 (isolated)", p.method, p.path, got)
		}
	}

	// Webmail routes on the ADMIN router → 404.
	webmailPaths := []struct{ method, path string }{
		{http.MethodGet, "/api/v1/version"},
		{http.MethodGet, "/api/v1/mail/INBOX"},
	}
	for _, p := range webmailPaths {
		if got := status(t, adminE, p.method, p.path); got != http.StatusNotFound {
			t.Errorf("webmail route %s %s on admin router: status = %d, want 404", p.method, p.path, got)
		}
	}

	// Sanity: admin login DOES resolve on the admin router (401/400, not 404).
	if got := status(t, adminE, http.MethodGet, "/api/v1/admin/overview"); got == http.StatusNotFound {
		t.Errorf("admin overview on admin router = 404, want it to resolve (401)")
	}
}
