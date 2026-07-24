package admin

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// claimsKey is the context key under which validated admin claims are stored.
const claimsKey = "admin_claims"

// JWTMiddleware validates the Bearer token on every admin API request AND
// re-validates the admin against the database: the account must still exist and
// be active, and (for domain_admins) the domain scope is refreshed from the DB
// so deactivations and scope removals take effect immediately — not only at
// token expiry. Missing/invalid/expired/revoked tokens yield 401.
func JWTMiddleware(db *gorm.DB, secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			token, hasBearer := strings.CutPrefix(auth, "Bearer ")
			if !hasBearer || token == "" {
				return fail(c, http.StatusUnauthorized, "missing bearer token")
			}
			tokenClaims, err := ParseToken(token, secret)
			if err != nil {
				return fail(c, http.StatusUnauthorized, "invalid or expired token")
			}

			// Re-validate against the DB (state may have changed since login).
			var a Admin
			if err := db.First(&a, "username = ?", tokenClaims.Username).Error; err != nil {
				return fail(c, http.StatusUnauthorized, "account no longer valid")
			}
			if !a.Active {
				return fail(c, http.StatusUnauthorized, "account is inactive")
			}

			fresh := Claims{Username: a.Username, Superadmin: a.Superadmin}
			if !a.Superadmin {
				var das []DomainAdmin
				if err := db.Where("username = ? AND active = ?", a.Username, true).Find(&das).Error; err != nil {
					return fail(c, http.StatusInternalServerError, "scope load failed")
				}
				for _, da := range das {
					fresh.Domains = append(fresh.Domains, da.Domain)
				}
			}
			c.Set(claimsKey, fresh)
			return next(c)
		}
	}
}

// claimsFrom returns the validated claims placed by JWTMiddleware. The bool is
// false when the middleware did not run (defensive; routes always chain it).
func claimsFrom(c *echo.Context) (Claims, bool) {
	v := c.Get(claimsKey)
	claims, ok := v.(Claims)
	return claims, ok
}

// requirePermission returns 403 unless the request's claims allow perm for the
// domain. Pass "" for global (non-domain-scoped) resources.
func requirePermission(c *echo.Context, perm Permission, domain string) (Claims, bool) {
	claims, ok := claimsFrom(c)
	if !ok {
		_ = fail(c, http.StatusUnauthorized, "not authenticated")
		return Claims{}, false
	}
	if !claims.Can(perm, domain) {
		_ = fail(c, http.StatusForbidden, "forbidden")
		return Claims{}, false
	}
	return claims, true
}
