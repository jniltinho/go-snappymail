package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
)

func TestCSRF(t *testing.T) {
	e := echo.New()
	e.Use(CSRF(false))
	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	e.POST("/submit", func(c *echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	t.Run("GET issues csrf cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d", rec.Code)
		}
		cookies := rec.Result().Cookies()
		var token string
		for _, c := range cookies {
			if c.Name == csrfCookieName {
				token = c.Value
			}
		}
		if token == "" {
			t.Fatal("expected csrf_token cookie")
		}
	})

	t.Run("POST without token is forbidden", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/submit", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("POST with matching token succeeds", func(t *testing.T) {
		reqGet := httptest.NewRequest(http.MethodGet, "/", nil)
		recGet := httptest.NewRecorder()
		e.ServeHTTP(recGet, reqGet)

		var token string
		for _, c := range recGet.Result().Cookies() {
			if c.Name == csrfCookieName {
				token = c.Value
			}
		}

		req := httptest.NewRequest(http.MethodPost, "/submit", nil)
		req.AddCookie(&http.Cookie{Name: csrfCookieName, Value: token})
		req.Header.Set(csrfHeaderName, token)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
		}
	})
}

func TestRequireAuth(t *testing.T) {
	cookieName := "gsn_session"

	e := echo.New()
	e.GET("/me", func(c *echo.Context) error {
		return c.String(http.StatusOK, "ok")
	}, RequireAuth(cookieName))

	t.Run("missing cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d", rec.Code)
		}
		if !strings.Contains(rec.Body.String(), "Not authenticated") {
			t.Fatalf("body = %s", rec.Body.String())
		}
	})
}

func TestNewRateLimit(t *testing.T) {
	e := echo.New()
	e.Use(NewRateLimit(2, time.Minute))
	e.GET("/limited", func(c *echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	ip := "203.0.113.10"
	doReq := func() *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "/limited", nil)
		req.RemoteAddr = ip + ":1234"
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		return rec
	}

	if rec := doReq(); rec.Code != http.StatusOK {
		t.Fatalf("first request status = %d", rec.Code)
	}
	if rec := doReq(); rec.Code != http.StatusOK {
		t.Fatalf("second request status = %d", rec.Code)
	}
	if rec := doReq(); rec.Code != http.StatusTooManyRequests {
		t.Fatalf("third request status = %d, want 429", rec.Code)
	}
}
