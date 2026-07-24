package admin

import (
	"errors"
	"net/http"
	"time"

	"go-snappymail/internal/config"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// Handlers holds the dependencies for the admin API: the mail database and the
// admin configuration (JWT secret / lifetime).
type Handlers struct {
	db  *gorm.DB
	cfg config.AdminConfig
}

// NewHandlers builds the admin API handlers.
func NewHandlers(db *gorm.DB, cfg config.AdminConfig) *Handlers {
	return &Handlers{db: db, cfg: cfg}
}

// loginRequest is the JSON body for POST /api/v1/admin/auth/login.
type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// loginResponse returns the signed token and the authenticated identity.
type loginResponse struct {
	Token      string   `json:"token"`
	Username   string   `json:"username"`
	Superadmin bool     `json:"superadmin"`
	Domains    []string `json:"domains,omitempty"`
}

// Login authenticates an admin and returns a JWT. Public route (no middleware).
func (h *Handlers) Login(c *echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return fail(c, http.StatusBadRequest, "invalid request body")
	}
	if req.Username == "" || req.Password == "" {
		return fail(c, http.StatusBadRequest, "username and password are required")
	}

	claims, err := Authenticate(h.db, req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			return fail(c, http.StatusUnauthorized, "invalid credentials")
		case errors.Is(err, ErrInactiveAccount):
			return fail(c, http.StatusForbidden, "account is inactive")
		default:
			return fail(c, http.StatusInternalServerError, "authentication error")
		}
	}

	maxAge := time.Duration(h.cfg.JWTMaxAgeSec) * time.Second
	token, err := IssueToken(claims, h.cfg.JWTSecret, maxAge, time.Now())
	if err != nil {
		return fail(c, http.StatusInternalServerError, "could not issue token")
	}
	return ok(c, loginResponse{
		Token:      token,
		Username:   claims.Username,
		Superadmin: claims.Superadmin,
		Domains:    claims.Domains,
	})
}

// Me returns the current admin identity from the token.
func (h *Handlers) Me(c *echo.Context) error {
	claims, ok2 := claimsFrom(c)
	if !ok2 {
		return fail(c, http.StatusUnauthorized, "not authenticated")
	}
	return ok(c, map[string]any{
		"username":   claims.Username,
		"superadmin": claims.Superadmin,
		"domains":    claims.Domains,
		"role":       string(claims.Role()),
	})
}

// overviewResponse is the Home dashboard payload. Counts are real; fields the
// mail schema cannot source are nil (rendered as "n/a"), never fabricated.
type overviewResponse struct {
	Accounts int64   `json:"accounts"`
	Domains  int64   `json:"domains"`
	Aliases  int64   `json:"aliases"`
	Admins   int64   `json:"admins"`
	Version  *string `json:"version"`         // n/a (not in mail schema)
	Servers  *int64  `json:"servers"`         // n/a
	Queue    *int64  `json:"queue"`           // n/a
	Sessions *int64  `json:"active_sessions"` // n/a
}

// Overview aggregates real counts for the Home dashboard. domain_admins see
// only their scoped domains/accounts/aliases; superadmins see everything.
func (h *Handlers) Overview(c *echo.Context) error {
	claims, allowed := requirePermission(c, PermDomainsRead, "")
	if !allowed {
		return nil
	}

	var res overviewResponse
	domainQ := h.db.Model(&Domain{})
	mbQ := h.db.Model(&Mailbox{})
	aliasQ := h.db.Model(&Alias{})

	if !claims.Superadmin {
		scope := claims.Domains
		if len(scope) == 0 {
			return ok(c, res) // no domains → all zeros
		}
		domainQ = domainQ.Where("domain IN ?", scope)
		mbQ = mbQ.Where("domain IN ?", scope)
		aliasQ = aliasQ.Where("domain IN ?", scope)
	}

	if err := domainQ.Count(&res.Domains).Error; err != nil {
		return fail(c, http.StatusInternalServerError, "overview query failed")
	}
	if err := mbQ.Count(&res.Accounts).Error; err != nil {
		return fail(c, http.StatusInternalServerError, "overview query failed")
	}
	if err := aliasQ.Count(&res.Aliases).Error; err != nil {
		return fail(c, http.StatusInternalServerError, "overview query failed")
	}
	// Admins count is a global (superadmin-only) figure.
	if claims.Superadmin {
		if err := h.db.Model(&Admin{}).Count(&res.Admins).Error; err != nil {
			return fail(c, http.StatusInternalServerError, "overview query failed")
		}
	}
	return ok(c, res)
}
