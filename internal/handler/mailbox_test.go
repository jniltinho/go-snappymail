package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-snappymail/internal/session"
	"github.com/labstack/echo/v5"
)

func TestMailboxCreateSubfolderValidation(t *testing.T) {
	h := &MailboxHandler{cfg: testConfig()}
	e := echo.New()
	s := &session.IMAPSession{Username: "user@test.local"}
	session.Set("mb-test", s)
	defer session.Delete("mb-test")

	req := httptest.NewRequest(http.MethodPost, "/folders", strings.NewReader("parent=INBOX"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("imap_session", s)

	if err := h.CreateSubfolder(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
}

func TestMessageFlagValidation(t *testing.T) {
	h := &MessageHandler{cfg: testConfig()}
	e := echo.New()
	s := &session.IMAPSession{Username: "user@test.local"}

	e.POST("/mail/:mailbox/:uid/flag", func(c *echo.Context) error {
		c.Set("imap_session", s)
		return h.Flag(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/mail/INBOX/abc/flag", strings.NewReader("flag=seen&value=1"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
}
