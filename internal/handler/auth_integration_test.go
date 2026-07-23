//go:build integration

package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	appMiddleware "github.com/jniltinho/go-snappymail/internal/server/middleware"
	"github.com/labstack/echo/v5"
)

func TestAuthDoLoginIntegration(t *testing.T) {
	host := os.Getenv("IMAP_TEST_HOST")
	user := os.Getenv("IMAP_TEST_USER")
	pass := os.Getenv("IMAP_TEST_PASS")
	if host == "" || user == "" || pass == "" {
		t.Skip("set IMAP_TEST_HOST, IMAP_TEST_USER, IMAP_TEST_PASS for integration test")
	}

	cfg := testConfig()
	cfg.IMAP.Host = host
	cfg.IMAP.Port = 993
	cfg.IMAP.TLS = true
	cfg.IMAP.InsecureSkipVerify = os.Getenv("IMAP_TEST_INSECURE") == "1"
	if tlsName := os.Getenv("IMAP_TEST_TLS_SERVER_NAME"); tlsName != "" {
		cfg.IMAP.TLSServerName = tlsName
	}

	h := &AuthHandler{cfg: cfg}
	e := echo.New()
	e.Use(appMiddleware.CSRF())
	e.GET("/", func(c *echo.Context) error { return c.NoContent(http.StatusOK) })
	e.POST("/login", h.DoLogin)

	// Obtain CSRF cookie
	reqGet := httptest.NewRequest(http.MethodGet, "/", nil)
	recGet := httptest.NewRecorder()
	e.ServeHTTP(recGet, reqGet)
	var csrfToken string
	for _, c := range recGet.Result().Cookies() {
		if c.Name == "csrf_token" {
			csrfToken = c.Value
		}
	}

	body := "username=" + user + "&password=" + pass
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(&http.Cookie{Name: "csrf_token", Value: csrfToken})
	req.Header.Set("X-CSRF-Token", csrfToken)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
}
