package admin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-snappymail/internal/config"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

const testSecret = "test-secret"

// testEnv wires an in-memory admin API for handler tests.
type testEnv struct {
	e  *echo.Echo
	db *gorm.DB
	h  *Handlers
}

func newEnv(t *testing.T) *testEnv {
	t.Helper()
	db := newTestDB(t)
	cfg := config.AdminConfig{JWTSecret: testSecret, JWTMaxAgeSec: 3600}
	h := NewHandlers(db, cfg)
	e := echo.New()
	h.RegisterRoutes(e.Group("/api/v1/admin"))
	return &testEnv{e: e, db: db, h: h}
}

// tokenFor authenticates the given claims into a signed token (bypassing the
// login handler for setup convenience).
func tokenFor(t *testing.T, claims Claims) string {
	t.Helper()
	tok, err := IssueToken(claims, testSecret, time.Hour, time.Now())
	if err != nil {
		t.Fatalf("IssueToken: %v", err)
	}
	return tok
}

func (env *testEnv) do(t *testing.T, method, path, token string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	env.e.ServeHTTP(rec, req)
	return rec
}

func decodeData(t *testing.T, rec *httptest.ResponseRecorder, into any) {
	t.Helper()
	var env envelope
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatalf("decode envelope: %v (body=%s)", err, rec.Body.String())
	}
	if into != nil && env.Data != nil {
		b, _ := json.Marshal(env.Data)
		if err := json.Unmarshal(b, into); err != nil {
			t.Fatalf("decode data: %v", err)
		}
	}
}

// superToken seeds an active superadmin and returns its token. The DB-backed
// middleware requires the account to exist and be active.
func superToken(t *testing.T, env *testEnv) string {
	t.Helper()
	hash, _ := HashPassword("password1")
	env.db.Where(Admin{Username: "boss@x"}).
		Assign(Admin{Password: hash, Superadmin: true, Active: true}).
		FirstOrCreate(&Admin{})
	return tokenFor(t, Claims{Username: "boss@x", Superadmin: true})
}

// daToken seeds an active domain_admin scoped to the given domains and returns
// its token. Scope is loaded from the DB by the middleware, so the domain_admins
// rows must exist.
func daToken(t *testing.T, env *testEnv, username string, domains ...string) string {
	t.Helper()
	hash, _ := HashPassword("password1")
	env.db.Where(Admin{Username: username}).
		Assign(Admin{Password: hash, Superadmin: false, Active: true}).
		FirstOrCreate(&Admin{})
	for _, d := range domains {
		env.db.Where(DomainAdmin{Username: username, Domain: d}).
			Assign(DomainAdmin{Active: true}).
			FirstOrCreate(&DomainAdmin{})
	}
	return tokenFor(t, Claims{Username: username, Domains: domains})
}

// ── Auth / middleware ────────────────────────────────────────────

func TestLoginFlow(t *testing.T) {
	env := newEnv(t)
	hash, _ := HashPassword("password1")
	env.db.Create(&Admin{Username: "boss@x", Password: hash, Superadmin: true, Active: true})

	t.Run("success returns token", func(t *testing.T) {
		rec := env.do(t, http.MethodPost, "/api/v1/admin/auth/login", "",
			loginRequest{Username: "boss@x", Password: "password1"})
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200 (body=%s)", rec.Code, rec.Body.String())
		}
		var resp loginResponse
		decodeData(t, rec, &resp)
		if resp.Token == "" || !resp.Superadmin {
			t.Errorf("resp = %+v, want token + superadmin", resp)
		}
	})
	t.Run("wrong password 401", func(t *testing.T) {
		rec := env.do(t, http.MethodPost, "/api/v1/admin/auth/login", "",
			loginRequest{Username: "boss@x", Password: "nope"})
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want 401", rec.Code)
		}
	})
	t.Run("missing fields 400", func(t *testing.T) {
		rec := env.do(t, http.MethodPost, "/api/v1/admin/auth/login", "", loginRequest{})
		if rec.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", rec.Code)
		}
	})
}

func TestAuthRequired(t *testing.T) {
	env := newEnv(t)
	tests := []struct{ method, path string }{
		{http.MethodGet, "/api/v1/admin/overview"},
		{http.MethodGet, "/api/v1/admin/domains"},
		{http.MethodGet, "/api/v1/admin/admins"},
	}
	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			rec := env.do(t, tt.method, tt.path, "", nil)
			if rec.Code != http.StatusUnauthorized {
				t.Errorf("no token: status = %d, want 401", rec.Code)
			}
			rec = env.do(t, tt.method, tt.path, "garbage.token.here", nil)
			if rec.Code != http.StatusUnauthorized {
				t.Errorf("bad token: status = %d, want 401", rec.Code)
			}
		})
	}
}

// ── Domains CRUD ─────────────────────────────────────────────────

func TestDomainsCRUD(t *testing.T) {
	env := newEnv(t)
	tok := superToken(t, env)

	// Create
	rec := env.do(t, http.MethodPost, "/api/v1/admin/domains", tok,
		domainWriteRequest{Domain: "example.com", Description: "Example"})
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: status = %d, want 201 (body=%s)", rec.Code, rec.Body.String())
	}

	// Duplicate → 409
	rec = env.do(t, http.MethodPost, "/api/v1/admin/domains", tok,
		domainWriteRequest{Domain: "example.com"})
	if rec.Code != http.StatusConflict {
		t.Errorf("duplicate: status = %d, want 409", rec.Code)
	}

	// Invalid name → 400
	rec = env.do(t, http.MethodPost, "/api/v1/admin/domains", tok,
		domainWriteRequest{Domain: "not a domain"})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("invalid: status = %d, want 400", rec.Code)
	}

	// List
	rec = env.do(t, http.MethodGet, "/api/v1/admin/domains", tok, nil)
	var list []Domain
	decodeData(t, rec, &list)
	if len(list) != 1 || list[0].Domain != "example.com" {
		t.Errorf("list = %+v, want [example.com]", list)
	}

	// Get
	rec = env.do(t, http.MethodGet, "/api/v1/admin/domains/example.com", tok, nil)
	if rec.Code != http.StatusOK {
		t.Errorf("get: status = %d, want 200", rec.Code)
	}
	// Get not found
	rec = env.do(t, http.MethodGet, "/api/v1/admin/domains/ghost.com", tok, nil)
	if rec.Code != http.StatusNotFound {
		t.Errorf("get missing: status = %d, want 404", rec.Code)
	}

	// Update
	active := false
	rec = env.do(t, http.MethodPut, "/api/v1/admin/domains/example.com", tok,
		domainWriteRequest{Description: "Updated", Active: &active})
	var updated Domain
	decodeData(t, rec, &updated)
	if updated.Description != "Updated" || updated.Active {
		t.Errorf("update = %+v, want Updated/inactive", updated)
	}

	// Delete
	rec = env.do(t, http.MethodDelete, "/api/v1/admin/domains/example.com", tok, nil)
	if rec.Code != http.StatusOK {
		t.Errorf("delete: status = %d, want 200", rec.Code)
	}
	rec = env.do(t, http.MethodDelete, "/api/v1/admin/domains/example.com", tok, nil)
	if rec.Code != http.StatusNotFound {
		t.Errorf("delete again: status = %d, want 404", rec.Code)
	}
}

// ── Mailboxes CRUD + scoping ─────────────────────────────────────

func TestMailboxesCRUDAndScope(t *testing.T) {
	env := newEnv(t)
	super := superToken(t, env)
	env.db.Create(&Domain{Domain: "example.com", Active: true})
	env.db.Create(&Domain{Domain: "other.com", Active: true})

	// Create mailbox (superadmin)
	rec := env.do(t, http.MethodPost, "/api/v1/admin/mailboxes", super,
		mailboxCreateRequest{Username: "user@example.com", Password: "password1", Name: "User"})
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: status = %d, want 201 (body=%s)", rec.Code, rec.Body.String())
	}
	var mb Mailbox
	decodeData(t, rec, &mb)
	if mb.Password != "" {
		t.Error("mailbox password leaked in response")
	}

	// Orphan (domain missing) → 400
	rec = env.do(t, http.MethodPost, "/api/v1/admin/mailboxes", super,
		mailboxCreateRequest{Username: "x@ghost.com", Password: "password1"})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("orphan: status = %d, want 400", rec.Code)
	}

	// Weak password → 400
	rec = env.do(t, http.MethodPost, "/api/v1/admin/mailboxes", super,
		mailboxCreateRequest{Username: "weak@example.com", Password: "short"})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("weak pass: status = %d, want 400", rec.Code)
	}

	// domain_admin of example.com can manage its mailbox…
	da := daToken(t, env, "da@x", "example.com")
	rec = env.do(t, http.MethodGet, "/api/v1/admin/mailboxes/user@example.com", da, nil)
	if rec.Code != http.StatusOK {
		t.Errorf("da own: status = %d, want 200", rec.Code)
	}
	// …but not a mailbox in other.com
	env.db.Create(&Mailbox{Username: "u2@other.com", Domain: "other.com", LocalPart: "u2", Active: true})
	rec = env.do(t, http.MethodGet, "/api/v1/admin/mailboxes/u2@other.com", da, nil)
	if rec.Code != http.StatusForbidden {
		t.Errorf("da other: status = %d, want 403", rec.Code)
	}
	// list is scoped to its domain only
	rec = env.do(t, http.MethodGet, "/api/v1/admin/mailboxes", da, nil)
	var mbs []Mailbox
	decodeData(t, rec, &mbs)
	if len(mbs) != 1 || mbs[0].Domain != "example.com" {
		t.Errorf("da list = %+v, want only example.com", mbs)
	}

	// Update password + delete (superadmin)
	rec = env.do(t, http.MethodPut, "/api/v1/admin/mailboxes/user@example.com", super,
		mailboxUpdateRequest{Password: "newpassword1", Name: "Renamed"})
	if rec.Code != http.StatusOK {
		t.Errorf("update: status = %d, want 200", rec.Code)
	}
	rec = env.do(t, http.MethodDelete, "/api/v1/admin/mailboxes/user@example.com", super, nil)
	if rec.Code != http.StatusOK {
		t.Errorf("delete: status = %d, want 200", rec.Code)
	}
}

// ── Aliases CRUD ─────────────────────────────────────────────────

func TestAliasesCRUD(t *testing.T) {
	env := newEnv(t)
	tok := superToken(t, env)
	env.db.Create(&Domain{Domain: "example.com", Active: true})

	rec := env.do(t, http.MethodPost, "/api/v1/admin/aliases", tok,
		aliasWriteRequest{Address: "info@example.com", Goto: "user@example.com, boss@example.com"})
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: status = %d, want 201 (body=%s)", rec.Code, rec.Body.String())
	}
	// invalid destination → 400
	rec = env.do(t, http.MethodPost, "/api/v1/admin/aliases", tok,
		aliasWriteRequest{Address: "bad@example.com", Goto: "not-an-email"})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("bad goto: status = %d, want 400", rec.Code)
	}
	// domain missing → 400
	rec = env.do(t, http.MethodPost, "/api/v1/admin/aliases", tok,
		aliasWriteRequest{Address: "x@ghost.com", Goto: "y@ghost.com"})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("orphan: status = %d, want 400", rec.Code)
	}
	// update + delete
	rec = env.do(t, http.MethodPut, "/api/v1/admin/aliases/info@example.com", tok,
		aliasWriteRequest{Goto: "only@example.com"})
	if rec.Code != http.StatusOK {
		t.Errorf("update: status = %d, want 200", rec.Code)
	}
	rec = env.do(t, http.MethodDelete, "/api/v1/admin/aliases/info@example.com", tok, nil)
	if rec.Code != http.StatusOK {
		t.Errorf("delete: status = %d, want 200", rec.Code)
	}
}

// ── Admins CRUD + superadmin-only ────────────────────────────────

func TestAdminsCRUDSuperadminOnly(t *testing.T) {
	env := newEnv(t)
	super := superToken(t, env)
	env.db.Create(&Domain{Domain: "example.com", Active: true})
	da := daToken(t, env, "da@x", "example.com")

	// domain_admin cannot list/create admins → 403
	if rec := env.do(t, http.MethodGet, "/api/v1/admin/admins", da, nil); rec.Code != http.StatusForbidden {
		t.Errorf("da list admins: status = %d, want 403", rec.Code)
	}
	if rec := env.do(t, http.MethodPost, "/api/v1/admin/admins", da,
		adminCreateRequest{Username: "n@x", Password: "password1"}); rec.Code != http.StatusForbidden {
		t.Errorf("da create admin: status = %d, want 403", rec.Code)
	}

	// superadmin creates a domain_admin with domain links
	rec := env.do(t, http.MethodPost, "/api/v1/admin/admins", super,
		adminCreateRequest{Username: "da2@x", Password: "password1", Domains: []string{"example.com"}})
	if rec.Code != http.StatusCreated {
		t.Fatalf("create admin: status = %d, want 201 (body=%s)", rec.Code, rec.Body.String())
	}
	var view adminView
	decodeData(t, rec, &view)
	if view.Superadmin || len(view.Domains) != 1 {
		t.Errorf("view = %+v, want domain_admin of one domain", view)
	}

	// weak password → 400
	if rec := env.do(t, http.MethodPost, "/api/v1/admin/admins", super,
		adminCreateRequest{Username: "w@x", Password: "short"}); rec.Code != http.StatusBadRequest {
		t.Errorf("weak: status = %d, want 400", rec.Code)
	}

	// update role + delete
	sa := true
	rec = env.do(t, http.MethodPut, "/api/v1/admin/admins/da2@x", super,
		adminUpdateRequest{Superadmin: &sa})
	decodeData(t, rec, &view)
	if !view.Superadmin {
		t.Errorf("update role: view = %+v, want superadmin", view)
	}
	if rec := env.do(t, http.MethodDelete, "/api/v1/admin/admins/da2@x", super, nil); rec.Code != http.StatusOK {
		t.Errorf("delete: status = %d, want 200", rec.Code)
	}
}

// ── Revocation & scoping edge cases (codex coverage gaps) ────────

func TestStaleTokenAfterDeactivation(t *testing.T) {
	env := newEnv(t)
	tok := daToken(t, env, "da@x", "example.com")
	env.db.Create(&Domain{Domain: "example.com", Active: true})

	// Works while active.
	if rec := env.do(t, http.MethodGet, "/api/v1/admin/overview", tok, nil); rec.Code != http.StatusOK {
		t.Fatalf("active: status = %d, want 200", rec.Code)
	}
	// Deactivate the admin → the same (unexpired) token must stop working now.
	env.db.Model(&Admin{}).Where("username = ?", "da@x").Update("active", false)
	if rec := env.do(t, http.MethodGet, "/api/v1/admin/overview", tok, nil); rec.Code != http.StatusUnauthorized {
		t.Errorf("deactivated: status = %d, want 401", rec.Code)
	}
}

func TestStaleTokenAfterScopeRemoval(t *testing.T) {
	env := newEnv(t)
	tok := daToken(t, env, "da@x", "example.com")
	env.db.Create(&Domain{Domain: "example.com", Active: true})
	env.db.Create(&Mailbox{Username: "u@example.com", Domain: "example.com", LocalPart: "u", Active: true})

	if rec := env.do(t, http.MethodGet, "/api/v1/admin/mailboxes/u@example.com", tok, nil); rec.Code != http.StatusOK {
		t.Fatalf("with scope: status = %d, want 200", rec.Code)
	}
	// Remove the domain scope → immediate effect (scope refreshed from DB).
	env.db.Delete(&DomainAdmin{}, "username = ? AND domain = ?", "da@x", "example.com")
	if rec := env.do(t, http.MethodGet, "/api/v1/admin/mailboxes/u@example.com", tok, nil); rec.Code != http.StatusForbidden {
		t.Errorf("scope removed: status = %d, want 403", rec.Code)
	}
}

func TestAliasForeignDomainForbidden(t *testing.T) {
	env := newEnv(t)
	env.db.Create(&Domain{Domain: "example.com", Active: true})
	env.db.Create(&Domain{Domain: "other.com", Active: true})
	env.db.Create(&Alias{Address: "i@other.com", Goto: "x@other.com", Domain: "other.com", Active: true})
	da := daToken(t, env, "da@x", "example.com")

	for _, m := range []string{http.MethodGet, http.MethodPut, http.MethodDelete} {
		var body any
		if m == http.MethodPut {
			body = aliasWriteRequest{Goto: "z@other.com"}
		}
		rec := env.do(t, m, "/api/v1/admin/aliases/i@other.com", da, body)
		if rec.Code != http.StatusForbidden {
			t.Errorf("%s foreign alias: status = %d, want 403", m, rec.Code)
		}
	}
}

func TestListDomainFilterScoped(t *testing.T) {
	env := newEnv(t)
	env.db.Create(&Domain{Domain: "example.com", Active: true})
	env.db.Create(&Domain{Domain: "other.com", Active: true})
	env.db.Create(&Mailbox{Username: "a@example.com", Domain: "example.com", LocalPart: "a", Active: true})
	env.db.Create(&Mailbox{Username: "b@other.com", Domain: "other.com", LocalPart: "b", Active: true})
	da := daToken(t, env, "da@x", "example.com")

	// A domain_admin cannot use ?domain= to peek at a foreign domain.
	if rec := env.do(t, http.MethodGet, "/api/v1/admin/mailboxes?domain=other.com", da, nil); rec.Code != http.StatusForbidden {
		t.Errorf("da ?domain=other.com: status = %d, want 403", rec.Code)
	}
	// Superadmin can filter by domain.
	super := superToken(t, env)
	rec := env.do(t, http.MethodGet, "/api/v1/admin/mailboxes?domain=other.com", super, nil)
	var mbs []Mailbox
	decodeData(t, rec, &mbs)
	if len(mbs) != 1 || mbs[0].Domain != "other.com" {
		t.Errorf("super filter = %+v, want only other.com", mbs)
	}
}

// ── Overview ─────────────────────────────────────────────────────

func TestOverviewCounts(t *testing.T) {
	env := newEnv(t)
	super := superToken(t, env)
	env.db.Create(&Domain{Domain: "example.com", Active: true})
	env.db.Create(&Domain{Domain: "other.com", Active: true})
	env.db.Create(&Mailbox{Username: "a@example.com", Domain: "example.com", LocalPart: "a", Active: true})
	env.db.Create(&Mailbox{Username: "b@other.com", Domain: "other.com", LocalPart: "b", Active: true})
	env.db.Create(&Alias{Address: "i@example.com", Goto: "a@example.com", Domain: "example.com", Active: true})

	rec := env.do(t, http.MethodGet, "/api/v1/admin/overview", super, nil)
	var ov overviewResponse
	decodeData(t, rec, &ov)
	if ov.Domains != 2 || ov.Accounts != 2 || ov.Aliases != 1 {
		t.Errorf("super overview = %+v, want 2/2/1", ov)
	}
	if ov.Version != nil || ov.Servers != nil || ov.Queue != nil {
		t.Errorf("n/a fields should be null, got %+v", ov)
	}

	// domain_admin sees only its scope
	da := daToken(t, env, "da@x", "example.com")
	rec = env.do(t, http.MethodGet, "/api/v1/admin/overview", da, nil)
	decodeData(t, rec, &ov)
	if ov.Domains != 1 || ov.Accounts != 1 || ov.Aliases != 1 {
		t.Errorf("da overview = %+v, want 1/1/1 scoped", ov)
	}
}
