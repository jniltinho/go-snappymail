package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"go-snappymail/internal/config"
	"go-snappymail/internal/ui"
)

// UIHandler exposes server-driven UI defaults (skin, etc.).
type UIHandler struct {
	cfg *config.Config
}

// Config returns the active skin and list of known skin ids for the SPA.
func (h *UIHandler) Config(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"skin":             h.cfg.UI.ResolvedSkin(),
		"available_skins":  ui.AvailableSkins(),
		"rows_per_page":    h.cfg.UI.RowsPerPage,
		"datetime_format":  h.cfg.UI.DatetimeFormat,
		"compose_html":     h.cfg.UI.ComposeHTML,
	})
}
