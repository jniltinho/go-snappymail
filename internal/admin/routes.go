package admin

import (
	"time"

	appMiddleware "go-snappymail/internal/server/middleware"

	"github.com/labstack/echo/v5"
)

// RegisterRoutes mounts the admin API under the given group (expected to be
// "/api/v1/admin" on the ADMIN Echo instance only). The login route is public;
// every other route is behind JWTMiddleware. This function must never be called
// on the webmail Echo — that is what keeps the admin surface isolated.
func (h *Handlers) RegisterRoutes(g *echo.Group) {
	// Public: obtain a token. Rate-limited to blunt online brute-force of admin
	// passwords (the webmail login has the same guard).
	g.POST("/auth/login", h.Login, appMiddleware.NewRateLimit(10, time.Minute))

	// Authenticated: everything else.
	auth := g.Group("", JWTMiddleware(h.db, h.cfg.JWTSecret))
	auth.GET("/me", h.Me)
	auth.GET("/overview", h.Overview)

	auth.GET("/domains", h.ListDomains)
	auth.POST("/domains", h.CreateDomain)
	auth.GET("/domains/:domain", h.GetDomain)
	auth.PUT("/domains/:domain", h.UpdateDomain)
	auth.DELETE("/domains/:domain", h.DeleteDomain)

	auth.GET("/mailboxes", h.ListMailboxes)
	auth.POST("/mailboxes", h.CreateMailbox)
	auth.GET("/mailboxes/:username", h.GetMailbox)
	auth.PUT("/mailboxes/:username", h.UpdateMailbox)
	auth.DELETE("/mailboxes/:username", h.DeleteMailbox)

	auth.GET("/aliases", h.ListAliases)
	auth.POST("/aliases", h.CreateAlias)
	auth.GET("/aliases/:address", h.GetAlias)
	auth.PUT("/aliases/:address", h.UpdateAlias)
	auth.DELETE("/aliases/:address", h.DeleteAlias)

	auth.GET("/admins", h.ListAdmins)
	auth.POST("/admins", h.CreateAdmin)
	auth.GET("/admins/:username", h.GetAdmin)
	auth.PUT("/admins/:username", h.UpdateAdmin)
	auth.DELETE("/admins/:username", h.DeleteAdmin)
}
