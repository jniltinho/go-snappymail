package admin

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// mailboxCreateRequest is the JSON body for creating a mailbox.
type mailboxCreateRequest struct {
	Username string `json:"username"` // full email address
	Password string `json:"password"`
	Name     string `json:"name"`
	Quota    int64  `json:"quota"`
	Active   *bool  `json:"active"`
}

// mailboxUpdateRequest updates editable fields (password optional).
type mailboxUpdateRequest struct {
	Password string `json:"password"`
	Name     string `json:"name"`
	Quota    int64  `json:"quota"`
	Active   *bool  `json:"active"`
}

// ListMailboxes returns mailboxes visible to the caller. An optional ?domain=
// narrows results; domain_admins are always constrained to their own domains.
func (h *Handlers) ListMailboxes(c *echo.Context) error {
	cl, ok2 := claimsFrom(c)
	if !ok2 {
		return fail(c, http.StatusUnauthorized, "not authenticated")
	}
	q := h.db.Model(&Mailbox{})
	if d := c.QueryParam("domain"); d != "" {
		if !cl.Can(PermMailboxesRead, d) {
			return fail(c, http.StatusForbidden, "forbidden")
		}
		q = q.Where("domain = ?", d)
	} else if !cl.Superadmin {
		if len(cl.Domains) == 0 {
			return ok(c, []Mailbox{})
		}
		q = q.Where("domain IN ?", cl.Domains)
	}
	var mbs []Mailbox
	if err := q.Order("username").Find(&mbs).Error; err != nil {
		return fail(c, http.StatusInternalServerError, "list mailboxes failed")
	}
	return ok(c, mbs)
}

// GetMailbox returns a single mailbox (scoped). Authorization uses the row's
// actual domain column after loading — not the address in the URL — so a
// non-canonical row cannot be reached with an in-scope-looking param.
func (h *Handlers) GetMailbox(c *echo.Context) error {
	username := c.Param("username")
	claims, ok2 := claimsFrom(c)
	if !ok2 {
		return fail(c, http.StatusUnauthorized, "not authenticated")
	}
	var mb Mailbox
	if err := h.db.First(&mb, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fail(c, http.StatusNotFound, "mailbox not found")
		}
		return fail(c, http.StatusInternalServerError, "get mailbox failed")
	}
	if !claims.Can(PermMailboxesRead, mb.Domain) {
		return fail(c, http.StatusForbidden, "forbidden")
	}
	return ok(c, mb)
}

// CreateMailbox creates a mailbox under an existing domain the caller manages.
func (h *Handlers) CreateMailbox(c *echo.Context) error {
	var req mailboxCreateRequest
	if err := c.Bind(&req); err != nil {
		return fail(c, http.StatusBadRequest, "invalid request body")
	}
	if !validEmail(req.Username) {
		return fail(c, http.StatusBadRequest, "invalid email address")
	}
	dom := domainOf(req.Username)
	if _, allowed := requirePermission(c, PermMailboxesWrite, dom); !allowed {
		return nil
	}
	if len(req.Password) < 8 {
		return fail(c, http.StatusBadRequest, "password must be at least 8 characters")
	}

	// Domain must exist (no orphan mailboxes).
	var dcount int64
	if err := h.db.Model(&Domain{}).Where("domain = ?", dom).Count(&dcount).Error; err != nil {
		return fail(c, http.StatusInternalServerError, "create mailbox failed")
	}
	if dcount == 0 {
		return fail(c, http.StatusBadRequest, "domain does not exist")
	}
	var mcount int64
	h.db.Model(&Mailbox{}).Where("username = ?", req.Username).Count(&mcount)
	if mcount > 0 {
		return fail(c, http.StatusConflict, "mailbox already exists")
	}

	hash, err := HashMailboxPassword(req.Password)
	if err != nil {
		return fail(c, http.StatusInternalServerError, "create mailbox failed")
	}
	mb := Mailbox{
		Username:  req.Username,
		Password:  hash,
		Name:      req.Name,
		Quota:     req.Quota,
		LocalPart: localPartOf(req.Username),
		Domain:    dom,
		Maildir:   dom + "/" + req.Username + "/",
		Active:    true,
	}
	if req.Active != nil {
		mb.Active = *req.Active
	}
	if err := h.db.Create(&mb).Error; err != nil {
		return fail(c, http.StatusInternalServerError, "create mailbox failed")
	}
	return created(c, mb)
}

// UpdateMailbox updates a mailbox (scoped). Password changes are re-hashed.
func (h *Handlers) UpdateMailbox(c *echo.Context) error {
	username := c.Param("username")
	claims, ok2 := claimsFrom(c)
	if !ok2 {
		return fail(c, http.StatusUnauthorized, "not authenticated")
	}
	var req mailboxUpdateRequest
	if err := c.Bind(&req); err != nil {
		return fail(c, http.StatusBadRequest, "invalid request body")
	}
	var mb Mailbox
	if err := h.db.First(&mb, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fail(c, http.StatusNotFound, "mailbox not found")
		}
		return fail(c, http.StatusInternalServerError, "update mailbox failed")
	}
	if !claims.Can(PermMailboxesWrite, mb.Domain) {
		return fail(c, http.StatusForbidden, "forbidden")
	}
	updates := map[string]any{"name": req.Name, "quota": req.Quota}
	if req.Active != nil {
		updates["active"] = *req.Active
	}
	if req.Password != "" {
		if len(req.Password) < 8 {
			return fail(c, http.StatusBadRequest, "password must be at least 8 characters")
		}
		hash, err := HashMailboxPassword(req.Password)
		if err != nil {
			return fail(c, http.StatusInternalServerError, "update mailbox failed")
		}
		updates["password"] = hash
	}
	if err := h.db.Model(&mb).Updates(updates).Error; err != nil {
		return fail(c, http.StatusInternalServerError, "update mailbox failed")
	}
	h.db.First(&mb, "username = ?", username)
	return ok(c, mb)
}

// DeleteMailbox removes a mailbox (scoped).
func (h *Handlers) DeleteMailbox(c *echo.Context) error {
	username := c.Param("username")
	claims, ok2 := claimsFrom(c)
	if !ok2 {
		return fail(c, http.StatusUnauthorized, "not authenticated")
	}
	var mb Mailbox
	if err := h.db.First(&mb, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fail(c, http.StatusNotFound, "mailbox not found")
		}
		return fail(c, http.StatusInternalServerError, "delete mailbox failed")
	}
	if !claims.Can(PermMailboxesWrite, mb.Domain) {
		return fail(c, http.StatusForbidden, "forbidden")
	}
	if err := h.db.Delete(&Mailbox{}, "username = ?", username).Error; err != nil {
		return fail(c, http.StatusInternalServerError, "delete mailbox failed")
	}
	return ok(c, map[string]string{"deleted": username})
}
