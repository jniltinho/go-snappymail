package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
)

func TestUIConfigHandler(t *testing.T) {
	cfg := testConfig()
	cfg.UI.Skin = "gmail"

	h := &UIHandler{cfg: cfg}
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/ui/config", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.Config(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"skin":"gmail"`) {
		t.Fatalf("body = %s", rec.Body.String())
	}
}
