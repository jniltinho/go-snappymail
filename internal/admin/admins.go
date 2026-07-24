package admin

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// adminView is the safe JSON representation of an Admin (no password) plus its
// managed domains.
type adminView struct {
	Username   string   `json:"username"`
	Superadmin bool     `json:"superadmin"`
	Active     bool     `json:"active"`
	Domains    []string `json:"domains,omitempty"`
}

// adminCreateRequest is the JSON body for creating an admin.
type adminCreateRequest struct {
	Username   string   `json:"username"`
	Password   string   `json:"password"`
	Superadmin bool     `json:"superadmin"`
	Domains    []string `json:"domains"` // managed domains for a domain_admin
	Active     *bool    `json:"active"`
}

// adminUpdateRequest updates an admin (password/domains optional).
type adminUpdateRequest struct {
	Password   string   `json:"password"`
	Superadmin *bool    `json:"superadmin"`
	Domains    []string `json:"domains"`
	Active     *bool    `json:"active"`
}

// errInvalidDomain marks a domain link request that names an invalid or
// nonexistent domain — mapped to 400 by the callers.
var errInvalidDomain = errors.New("admin: invalid or nonexistent domain")

// linkDomains creates active domain_admin rows for each domain, validating that
// each is well-formed and exists. Runs inside the caller's transaction so a bad
// domain rolls the whole change back.
func linkDomains(tx *gorm.DB, username string, domains []string) error {
	for _, d := range domains {
		if !validDomain(d) {
			return errInvalidDomain
		}
		var count int64
		if err := tx.Model(&Domain{}).Where("domain = ?", d).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return errInvalidDomain
		}
		if err := tx.Create(&DomainAdmin{Username: username, Domain: d, Active: true}).Error; err != nil {
			return err
		}
	}
	return nil
}

// lastActiveSuperadmin reports whether username is the only active superadmin
// left — used to refuse a demotion/deactivation/deletion that would lock every
// superadmin out of the panel.
func lastActiveSuperadmin(tx *gorm.DB, username string) (bool, error) {
	var others int64
	if err := tx.Model(&Admin{}).
		Where("superadmin = ? AND active = ? AND username <> ?", true, true, username).
		Count(&others).Error; err != nil {
		return false, err
	}
	return others == 0, nil
}

func (h *Handlers) toAdminView(a Admin) adminView {
	v := adminView{Username: a.Username, Superadmin: a.Superadmin, Active: a.Active}
	if !a.Superadmin {
		var das []DomainAdmin
		h.db.Where("username = ?", a.Username).Find(&das)
		for _, da := range das {
			v.Domains = append(v.Domains, da.Domain)
		}
	}
	return v
}

// ListAdmins lists panel admins. Superadmin only (admins:read).
func (h *Handlers) ListAdmins(c *echo.Context) error {
	if _, allowed := requirePermission(c, PermAdminsRead, ""); !allowed {
		return nil
	}
	var admins []Admin
	if err := h.db.Order("username").Find(&admins).Error; err != nil {
		return fail(c, http.StatusInternalServerError, "list admins failed")
	}
	views := make([]adminView, 0, len(admins))
	for _, a := range admins {
		views = append(views, h.toAdminView(a))
	}
	return ok(c, views)
}

// GetAdmin returns a single admin. Superadmin only.
func (h *Handlers) GetAdmin(c *echo.Context) error {
	if _, allowed := requirePermission(c, PermAdminsRead, ""); !allowed {
		return nil
	}
	username := c.Param("username")
	var a Admin
	if err := h.db.First(&a, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fail(c, http.StatusNotFound, "admin not found")
		}
		return fail(c, http.StatusInternalServerError, "get admin failed")
	}
	return ok(c, h.toAdminView(a))
}

// CreateAdmin creates a panel admin and its domain links. Superadmin only.
func (h *Handlers) CreateAdmin(c *echo.Context) error {
	if _, allowed := requirePermission(c, PermAdminsWrite, ""); !allowed {
		return nil
	}
	var req adminCreateRequest
	if err := c.Bind(&req); err != nil {
		return fail(c, http.StatusBadRequest, "invalid request body")
	}
	if req.Username == "" || len(req.Password) < 8 {
		return fail(c, http.StatusBadRequest, "username and an 8+ char password are required")
	}
	var count int64
	if err := h.db.Model(&Admin{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		return fail(c, http.StatusInternalServerError, "create admin failed")
	}
	if count > 0 {
		return fail(c, http.StatusConflict, "admin already exists")
	}
	hash, err := HashPassword(req.Password)
	if err != nil {
		return fail(c, http.StatusInternalServerError, "create admin failed")
	}
	a := Admin{Username: req.Username, Password: hash, Superadmin: req.Superadmin, Active: true}
	if req.Active != nil {
		a.Active = *req.Active
	}
	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&a).Error; err != nil {
			return err
		}
		if !req.Superadmin {
			return linkDomains(tx, a.Username, req.Domains)
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, errInvalidDomain) {
			return fail(c, http.StatusBadRequest, "one or more domains are invalid or do not exist")
		}
		return fail(c, http.StatusInternalServerError, "create admin failed")
	}
	return created(c, h.toAdminView(a))
}

// UpdateAdmin updates an admin (password/role/domains/active). Superadmin only.
func (h *Handlers) UpdateAdmin(c *echo.Context) error {
	if _, allowed := requirePermission(c, PermAdminsWrite, ""); !allowed {
		return nil
	}
	username := c.Param("username")
	var req adminUpdateRequest
	if err := c.Bind(&req); err != nil {
		return fail(c, http.StatusBadRequest, "invalid request body")
	}
	var a Admin
	if err := h.db.First(&a, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fail(c, http.StatusNotFound, "admin not found")
		}
		return fail(c, http.StatusInternalServerError, "update admin failed")
	}
	updates := map[string]any{}
	if req.Superadmin != nil {
		updates["superadmin"] = *req.Superadmin
	}
	if req.Active != nil {
		updates["active"] = *req.Active
	}
	if req.Password != "" {
		if len(req.Password) < 8 {
			return fail(c, http.StatusBadRequest, "password must be at least 8 characters")
		}
		hash, err := HashPassword(req.Password)
		if err != nil {
			return fail(c, http.StatusInternalServerError, "update admin failed")
		}
		updates["password"] = hash
	}
	// Refuse a change that would strip the last active superadmin of its access.
	demote := req.Superadmin != nil && !*req.Superadmin
	deactivate := req.Active != nil && !*req.Active
	if a.Superadmin && a.Active && (demote || deactivate) {
		last, err := lastActiveSuperadmin(h.db, username)
		if err != nil {
			return fail(c, http.StatusInternalServerError, "update admin failed")
		}
		if last {
			return fail(c, http.StatusConflict, "cannot demote or deactivate the last active superadmin")
		}
	}
	err := h.db.Transaction(func(tx *gorm.DB) error {
		if len(updates) > 0 {
			if err := tx.Model(&a).Updates(updates).Error; err != nil {
				return err
			}
		}
		// Replace domain links when Domains is provided.
		if req.Domains != nil {
			if err := tx.Delete(&DomainAdmin{}, "username = ?", username).Error; err != nil {
				return err
			}
			return linkDomains(tx, username, req.Domains)
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, errInvalidDomain) {
			return fail(c, http.StatusBadRequest, "one or more domains are invalid or do not exist")
		}
		return fail(c, http.StatusInternalServerError, "update admin failed")
	}
	if err := h.db.First(&a, "username = ?", username).Error; err != nil {
		return fail(c, http.StatusInternalServerError, "update admin failed")
	}
	return ok(c, h.toAdminView(a))
}

// DeleteAdmin removes an admin and its domain links. Superadmin only.
func (h *Handlers) DeleteAdmin(c *echo.Context) error {
	if _, allowed := requirePermission(c, PermAdminsWrite, ""); !allowed {
		return nil
	}
	username := c.Param("username")
	var a Admin
	if err := h.db.First(&a, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fail(c, http.StatusNotFound, "admin not found")
		}
		return fail(c, http.StatusInternalServerError, "delete admin failed")
	}
	if a.Superadmin && a.Active {
		last, err := lastActiveSuperadmin(h.db, username)
		if err != nil {
			return fail(c, http.StatusInternalServerError, "delete admin failed")
		}
		if last {
			return fail(c, http.StatusConflict, "cannot delete the last active superadmin")
		}
	}
	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&DomainAdmin{}, "username = ?", username).Error; err != nil {
			return err
		}
		return tx.Delete(&Admin{}, "username = ?", username).Error
	})
	if err != nil {
		return fail(c, http.StatusInternalServerError, "delete admin failed")
	}
	return ok(c, map[string]string{"deleted": username})
}
