package admin

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Auth errors surfaced to callers (handlers map these to 401/403).
var (
	ErrInvalidCredentials = errors.New("admin: invalid credentials")
	ErrInactiveAccount    = errors.New("admin: account is inactive")
	ErrForbidden          = errors.New("admin: forbidden")
)

// dummyHash is a real bcrypt hash at DefaultCost, used for unknown-user logins
// so bcrypt actually runs the KDF and login timing does not leak account
// existence. Computed once at init.
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("timing-equalizer"), bcrypt.DefaultCost)

// Role is a coarse admin role. Fine-grained access is decided by Permission
// checks that also consider domain scope.
type Role string

const (
	RoleSuperadmin  Role = "superadmin"
	RoleDomainAdmin Role = "domain_admin"
)

// Permission is a resource:action capability checked per request.
type Permission string

const (
	PermDomainsRead    Permission = "domains:read"
	PermDomainsWrite   Permission = "domains:write"
	PermMailboxesRead  Permission = "mailboxes:read"
	PermMailboxesWrite Permission = "mailboxes:write"
	PermAliasesRead    Permission = "aliases:read"
	PermAliasesWrite   Permission = "aliases:write"
	PermAdminsRead     Permission = "admins:read"
	PermAdminsWrite    Permission = "admins:write"
)

// Claims is the admin JWT payload.
type Claims struct {
	Username   string   `json:"username"`
	Superadmin bool     `json:"superadmin"`
	Domains    []string `json:"domains"` // managed domains (domain_admin scope)
	jwt.RegisteredClaims
}

// Role returns the role implied by the claims.
func (c Claims) Role() Role {
	if c.Superadmin {
		return RoleSuperadmin
	}
	return RoleDomainAdmin
}

// managesDomain reports whether a domain_admin manages the given domain.
func (c Claims) managesDomain(domain string) bool {
	for _, d := range c.Domains {
		if strings.EqualFold(d, domain) {
			return true
		}
	}
	return false
}

// Can reports whether the claims allow the permission. The domain argument
// scopes domain_admin checks; pass "" for global resources.
//
// superadmin can do everything. domain_admin can read/write mailboxes and
// aliases within its own domains and read its own domains, but cannot create or
// delete domains, and cannot touch admins.
func (c Claims) Can(perm Permission, domain string) bool {
	if c.Superadmin {
		return true
	}
	switch perm {
	case PermDomainsRead:
		return domain == "" || c.managesDomain(domain)
	case PermMailboxesRead, PermMailboxesWrite, PermAliasesRead, PermAliasesWrite:
		return domain != "" && c.managesDomain(domain)
	default:
		// domains:write, admins:read, admins:write → superadmin only.
		return false
	}
}

// HashPassword hashes an admin-panel login password with bcrypt.
func HashPassword(plain string) (string, error) {
	if plain == "" {
		return "", errors.New("admin: empty password")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("admin: hash password: %w", err)
	}
	return string(h), nil
}

// CheckPassword verifies a plaintext password against a stored hash. It accepts
// a bare bcrypt hash or a Dovecot-style {BLF-CRYPT} prefixed bcrypt hash.
func CheckPassword(hash, plain string) bool {
	h := strings.TrimPrefix(hash, "{BLF-CRYPT}")
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(plain)) == nil
}

// HashMailboxPassword produces a Dovecot-compatible {BLF-CRYPT} hash for
// storing mailbox passwords the mail server can verify.
func HashMailboxPassword(plain string) (string, error) {
	h, err := HashPassword(plain)
	if err != nil {
		return "", err
	}
	return "{BLF-CRYPT}" + h, nil
}

// Authenticate verifies an admin's credentials against the admin table and
// returns validated Claims (without a signed token). Timing is kept roughly
// constant by always running a bcrypt comparison.
func Authenticate(db *gorm.DB, username, password string) (Claims, error) {
	var a Admin
	err := db.First(&a, "username = ?", username).Error

	// Always run bcrypt (real hash for unknown users) to equalize timing.
	stored := a.Password
	if err != nil {
		stored = string(dummyHash)
	}
	ok := CheckPassword(stored, password)
	if err != nil || !ok {
		return Claims{}, ErrInvalidCredentials
	}
	if !a.Active {
		return Claims{}, ErrInactiveAccount
	}

	claims := Claims{Username: a.Username, Superadmin: a.Superadmin}
	if !a.Superadmin {
		// Fail closed: if the scope query errors we must not authenticate with
		// an empty (or partial) domain set.
		var das []DomainAdmin
		if err := db.Where("username = ? AND active = ?", a.Username, true).Find(&das).Error; err != nil {
			return Claims{}, fmt.Errorf("admin: load domain scope: %w", err)
		}
		for _, da := range das {
			claims.Domains = append(claims.Domains, da.Domain)
		}
	}
	return claims, nil
}

// IssueToken signs a JWT for the claims with the given secret and lifetime.
func IssueToken(claims Claims, secret string, maxAge time.Duration, now time.Time) (string, error) {
	if secret == "" {
		return "", errors.New("admin: empty jwt secret")
	}
	claims.RegisteredClaims = jwt.RegisteredClaims{
		Subject:   claims.Username,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(maxAge)),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("admin: sign token: %w", err)
	}
	return s, nil
}

// ParseToken validates a signed JWT and returns its claims. It rejects wrong
// signing methods and expired tokens.
func ParseToken(tokenStr, secret string) (Claims, error) {
	if secret == "" {
		return Claims{}, errors.New("admin: empty jwt secret")
	}
	var claims Claims
	_, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (any, error) {
		// Pin to HS256 exactly — reject alg confusion within the HMAC family
		// (HS384/HS512) and any non-HMAC method.
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("admin: unexpected signing method %v", t.Header["alg"])
		}
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return Claims{}, fmt.Errorf("admin: parse token: %w", err)
	}
	return claims, nil
}
