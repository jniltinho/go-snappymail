package admin

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

// envelope is the consistent JSON response shape for all /api/v1/admin/* routes.
type envelope struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

// ok writes a 200 JSON response wrapping data.
func ok(c *echo.Context, data any) error {
	return c.JSON(http.StatusOK, envelope{Data: data})
}

// created writes a 201 JSON response wrapping data.
func created(c *echo.Context, data any) error {
	return c.JSON(http.StatusCreated, envelope{Data: data})
}

// fail writes an error JSON response with the given status code.
func fail(c *echo.Context, code int, msg string) error {
	return c.JSON(code, envelope{Error: msg})
}
