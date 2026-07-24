package admin

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestHashAndCheckPassword(t *testing.T) {
	hash, err := HashPassword("s3nha-forte")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if !CheckPassword(hash, "s3nha-forte") {
		t.Error("CheckPassword failed for correct password")
	}
	if CheckPassword(hash, "errada") {
		t.Error("CheckPassword accepted a wrong password")
	}
	if _, err := HashPassword(""); err == nil {
		t.Error("HashPassword(\"\") = nil error, want error")
	}
}

func TestCheckPasswordDovecotPrefix(t *testing.T) {
	hash, err := HashMailboxPassword("dovecot-pass")
	if err != nil {
		t.Fatalf("HashMailboxPassword: %v", err)
	}
	if !strings.HasPrefix(hash, "{BLF-CRYPT}") {
		t.Errorf("mailbox hash missing {BLF-CRYPT} prefix: %q", hash)
	}
	if !CheckPassword(hash, "dovecot-pass") {
		t.Error("CheckPassword failed for {BLF-CRYPT} hash")
	}
}

func TestClaimsCan(t *testing.T) {
	super := Claims{Superadmin: true}
	da := Claims{Domains: []string{"example.com"}}

	tests := []struct {
		name   string
		claims Claims
		perm   Permission
		domain string
		want   bool
	}{
		{"super does anything", super, PermAdminsWrite, "", true},
		{"super any domain", super, PermMailboxesWrite, "other.com", true},
		{"da reads own domain", da, PermDomainsRead, "example.com", true},
		{"da not other domain read", da, PermDomainsRead, "other.com", false},
		{"da writes own mailbox", da, PermMailboxesWrite, "example.com", true},
		{"da not other mailbox", da, PermMailboxesWrite, "other.com", false},
		{"da writes own alias", da, PermAliasesWrite, "example.com", true},
		{"da cannot create domain", da, PermDomainsWrite, "example.com", false},
		{"da cannot manage admins", da, PermAdminsRead, "", false},
		{"da global domains read", da, PermDomainsRead, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.claims.Can(tt.perm, tt.domain); got != tt.want {
				t.Errorf("Can(%s, %q) = %v, want %v", tt.perm, tt.domain, got, tt.want)
			}
		})
	}
}

func TestClaimsRole(t *testing.T) {
	if (Claims{Superadmin: true}).Role() != RoleSuperadmin {
		t.Error("superadmin claims should have RoleSuperadmin")
	}
	if (Claims{}).Role() != RoleDomainAdmin {
		t.Error("non-super claims should have RoleDomainAdmin")
	}
}

func TestAuthenticate(t *testing.T) {
	db := newTestDB(t)

	hash, _ := HashPassword("adminpass")
	if err := db.Create(&Admin{Username: "boss@x", Password: hash, Superadmin: true, Active: true}).Error; err != nil {
		t.Fatalf("seed superadmin: %v", err)
	}
	daHash, _ := HashPassword("dapass")
	db.Create(&Admin{Username: "da@x", Password: daHash, Active: true})
	db.Create(&DomainAdmin{Username: "da@x", Domain: "example.com", Active: true})
	inactiveHash, _ := HashPassword("x")
	// GORM's default:true tag turns a zero-value (false) Active into true on
	// create, so deactivate explicitly via Update.
	db.Create(&Admin{Username: "off@x", Password: inactiveHash})
	db.Model(&Admin{}).Where("username = ?", "off@x").Update("active", false)

	t.Run("superadmin ok", func(t *testing.T) {
		c, err := Authenticate(db, "boss@x", "adminpass")
		if err != nil || !c.Superadmin {
			t.Fatalf("got %+v err %v, want superadmin", c, err)
		}
	})
	t.Run("domain_admin scope loaded", func(t *testing.T) {
		c, err := Authenticate(db, "da@x", "dapass")
		if err != nil {
			t.Fatalf("err %v", err)
		}
		if c.Superadmin || len(c.Domains) != 1 || c.Domains[0] != "example.com" {
			t.Errorf("claims = %+v, want domain_admin of example.com", c)
		}
	})
	t.Run("wrong password", func(t *testing.T) {
		if _, err := Authenticate(db, "boss@x", "nope"); err != ErrInvalidCredentials {
			t.Errorf("err = %v, want ErrInvalidCredentials", err)
		}
	})
	t.Run("unknown user", func(t *testing.T) {
		if _, err := Authenticate(db, "ghost@x", "whatever"); err != ErrInvalidCredentials {
			t.Errorf("err = %v, want ErrInvalidCredentials", err)
		}
	})
	t.Run("inactive account", func(t *testing.T) {
		if _, err := Authenticate(db, "off@x", "x"); err != ErrInactiveAccount {
			t.Errorf("err = %v, want ErrInactiveAccount", err)
		}
	})
}

func TestJWTRoundTrip(t *testing.T) {
	// Use the real clock: ParseToken validates expiry against time.Now().
	now := time.Now()
	claims := Claims{Username: "boss@x", Superadmin: true}

	tok, err := IssueToken(claims, "secret", time.Hour, now)
	if err != nil {
		t.Fatalf("IssueToken: %v", err)
	}
	got, err := ParseToken(tok, "secret")
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}
	if got.Username != "boss@x" || !got.Superadmin {
		t.Errorf("parsed = %+v, want boss@x/superadmin", got)
	}
}

func TestJWTRejectsWrongSecret(t *testing.T) {
	tok, _ := IssueToken(Claims{Username: "u"}, "right", time.Hour, time.Now())
	if _, err := ParseToken(tok, "wrong"); err == nil {
		t.Error("ParseToken with wrong secret = nil error, want error")
	}
}

func TestJWTRejectsExpired(t *testing.T) {
	past := time.Now().Add(-2 * time.Hour)
	tok, _ := IssueToken(Claims{Username: "u"}, "secret", time.Hour, past) // expired 1h ago
	if _, err := ParseToken(tok, "secret"); err == nil {
		t.Error("ParseToken of expired token = nil error, want error")
	}
}

func TestJWTRejectsWrongSigningMethod(t *testing.T) {
	// Forge a token with 'none' alg; ParseToken must reject non-HMAC.
	tok := jwt.NewWithClaims(jwt.SigningMethodNone, Claims{Username: "attacker"})
	s, _ := tok.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if _, err := ParseToken(s, "secret"); err == nil {
		t.Error("ParseToken accepted a 'none'-signed token, want rejection")
	}
}

func TestIssueTokenEmptySecret(t *testing.T) {
	if _, err := IssueToken(Claims{Username: "u"}, "", time.Hour, time.Now()); err == nil {
		t.Error("IssueToken with empty secret = nil error, want error")
	}
}

func TestParseTokenEmptySecret(t *testing.T) {
	tok, _ := IssueToken(Claims{Username: "u"}, "secret", time.Hour, time.Now())
	if _, err := ParseToken(tok, ""); err == nil {
		t.Error("ParseToken with empty secret = nil error, want error")
	}
}

func TestJWTRejectsAlgConfusion(t *testing.T) {
	// A token signed HS384/HS512 with the same secret must be rejected — only
	// HS256 is accepted.
	for _, m := range []*jwt.SigningMethodHMAC{jwt.SigningMethodHS384, jwt.SigningMethodHS512} {
		tok := jwt.NewWithClaims(m, Claims{Username: "attacker"})
		s, err := tok.SignedString([]byte("secret"))
		if err != nil {
			t.Fatalf("sign %v: %v", m.Alg(), err)
		}
		if _, err := ParseToken(s, "secret"); err == nil {
			t.Errorf("ParseToken accepted %s token, want rejection", m.Alg())
		}
	}
}
