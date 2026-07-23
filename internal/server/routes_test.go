package server

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"go-snappymail/internal/config"
	"go-snappymail/internal/handler"
	"github.com/labstack/echo/v5"
)

func TestVersionRoute(t *testing.T) {
	e := echo.New()
	h := handler.New(&config.Config{})
	pass := func(next echo.HandlerFunc) echo.HandlerFunc { return next }
	registerAPIRoutes(e.Group("/api/v1"), h, pass, pass)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/version", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "go-snappymail") {
		t.Fatalf("body = %s", rec.Body.String())
	}
}

func TestStaticIndexFallback(t *testing.T) {
	distFS := fstest.MapFS{
		"index.html": &fstest.MapFile{
			Data: []byte("<html><body>go-snappymail</body></html>"),
		},
	}

	e := echo.New()
	h := handler.New(&config.Config{})
	registerRoutes(e, &config.Config{}, h, fs.FS(distFS))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "go-snappymail") {
		t.Fatalf("body = %s", rec.Body.String())
	}
}
