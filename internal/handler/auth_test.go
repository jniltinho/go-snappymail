package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jniltinho/go-snappymail/internal/config"
	"github.com/jniltinho/go-snappymail/internal/session"
	"github.com/labstack/echo/v5"
)

func testConfig() *config.Config {
	return &config.Config{
		Server: config.ServerConfig{
			SecretKey: "test-secret-key-32-chars!!",
			Debug:     false,
		},
		IMAP: config.IMAPConfig{
			Host:       "127.0.0.1",
			Port:       1993,
			TLS:        false,
			TimeoutSec: 1,
		},
		Session: config.SessionConfig{
			Name:     "gsn_session",
			MaxAge:   3600,
			HTTPOnly: true,
		},
		UI: config.UIConfig{
			DatetimeFormat: "2006-01-02 15:04",
		},
	}
}

func TestAuthDoLoginValidation(t *testing.T) {
	h := &AuthHandler{cfg: testConfig()}
	e := echo.New()

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantErr    string
	}{
		{
			name:       "missing username",
			body:       "password=secret",
			wantStatus: http.StatusBadRequest,
			wantErr:    "Username and password are required",
		},
		{
			name:       "missing password",
			body:       "username=user@example.com",
			wantStatus: http.StatusBadRequest,
			wantErr:    "Username and password are required",
		},
		{
			name:       "username too long",
			body:       "username=" + strings.Repeat("a", 300) + "@x.com&password=x",
			wantStatus: http.StatusBadRequest,
			wantErr:    "Input too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if err := h.DoLogin(c); err != nil {
				t.Fatalf("DoLogin() error = %v", err)
			}
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
			}
			var resp map[string]string
			if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(resp["error"], tt.wantErr) {
				t.Fatalf("error = %q, want substring %q", resp["error"], tt.wantErr)
			}
		})
	}
}

func TestAuthDoLogout(t *testing.T) {
	session.Set("logout-test", &session.IMAPSession{Username: "user@example.com"})

	h := &AuthHandler{cfg: testConfig()}
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req.AddCookie(&http.Cookie{Name: "gsn_session", Value: "logout-test"})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.DoLogout(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if session.Get("logout-test") != nil {
		t.Fatal("session should be deleted")
	}
}

func TestAuthMe(t *testing.T) {
	s := &session.IMAPSession{Username: "alice@example.com"}
	session.Set("me-test", s)

	h := &AuthHandler{cfg: testConfig()}
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("imap_session", s)

	if err := h.Me(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp["username"] != "alice@example.com" {
		t.Fatalf("username = %q", resp["username"])
	}
	if resp["datetime_format"] == "" {
		t.Fatal("expected datetime_format")
	}

	session.Delete("me-test")
}
