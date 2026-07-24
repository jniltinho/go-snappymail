package admin

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// domainWriteRequest is the JSON body for create/update domain.
type domainWriteRequest struct {
	Domain      string `json:"domain"`
	Description string `json:"description"`
	MaxQuota    int64  `json:"maxquota"`
	Active      *bool  `json:"active"`
}

// ListDomains returns domains visible to the caller (all for superadmin, scoped
// for domain_admin).
func (h *Handlers) ListDomains(c *echo.Context) error {
	claims, allowed := requirePermission(c, PermDomainsRead, "")
	if !allowed {
		return nil
	}
	q := h.db.Model(&Domain{})
	if !claims.Superadmin {
		if len(claims.Domains) == 0 {
			return ok(c, []Domain{})
		}
		q = q.Where("domain IN ?", claims.Domains)
	}
	var domains []Domain
	if err := q.Order("domain").Find(&domains).Error; err != nil {
		return fail(c, http.StatusInternalServerError, "list domains failed")
	}
	return ok(c, domains)
}

// GetDomain returns a single domain by name (scoped).
func (h *Handlers) GetDomain(c *echo.Context) error {
	name := c.Param("domain")
	if _, allowed := requirePermission(c, PermDomainsRead, name); !allowed {
		return nil
	}
	var d Domain
	if err := h.db.First(&d, "domain = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fail(c, http.StatusNotFound, "domain not found")
		}
		return fail(c, http.StatusInternalServerError, "get domain failed")
	}
	return ok(c, d)
}

// CreateDomain creates a domain. Superadmin only (domains:write).
func (h *Handlers) CreateDomain(c *echo.Context) error {
	var req domainWriteRequest
	if err := c.Bind(&req); err != nil {
		return fail(c, http.StatusBadRequest, "invalid request body")
	}
	if _, allowed := requirePermission(c, PermDomainsWrite, req.Domain); !allowed {
		return nil
	}
	if !validDomain(req.Domain) {
		return fail(c, http.StatusBadRequest, "invalid domain name")
	}

	var count int64
	if err := h.db.Model(&Domain{}).Where("domain = ?", req.Domain).Count(&count).Error; err != nil {
		return fail(c, http.StatusInternalServerError, "create domain failed")
	}
	if count > 0 {
		return fail(c, http.StatusConflict, "domain already exists")
	}

	d := Domain{
		Domain:      req.Domain,
		Description: req.Description,
		MaxQuota:    req.MaxQuota,
		Active:      true,
	}
	if req.Active != nil {
		d.Active = *req.Active
	}
	if err := h.db.Create(&d).Error; err != nil {
		return fail(c, http.StatusInternalServerError, "create domain failed")
	}
	return created(c, d)
}

// UpdateDomain updates a domain's editable fields. Superadmin only.
func (h *Handlers) UpdateDomain(c *echo.Context) error {
	name := c.Param("domain")
	if _, allowed := requirePermission(c, PermDomainsWrite, name); !allowed {
		return nil
	}
	var req domainWriteRequest
	if err := c.Bind(&req); err != nil {
		return fail(c, http.StatusBadRequest, "invalid request body")
	}

	var d Domain
	if err := h.db.First(&d, "domain = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fail(c, http.StatusNotFound, "domain not found")
		}
		return fail(c, http.StatusInternalServerError, "update domain failed")
	}
	updates := map[string]any{
		"description": req.Description,
		"maxquota":    req.MaxQuota,
	}
	if req.Active != nil {
		updates["active"] = *req.Active
	}
	if err := h.db.Model(&d).Updates(updates).Error; err != nil {
		return fail(c, http.StatusInternalServerError, "update domain failed")
	}
	h.db.First(&d, "domain = ?", name)
	return ok(c, d)
}

// DeleteDomain removes a domain and its mailboxes/aliases/admin links in a
// transaction. Superadmin only.
func (h *Handlers) DeleteDomain(c *echo.Context) error {
	name := c.Param("domain")
	if _, allowed := requirePermission(c, PermDomainsWrite, name); !allowed {
		return nil
	}
	var d Domain
	if err := h.db.First(&d, "domain = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fail(c, http.StatusNotFound, "domain not found")
		}
		return fail(c, http.StatusInternalServerError, "delete domain failed")
	}
	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Mailbox{}, "domain = ?", name).Error; err != nil {
			return err
		}
		if err := tx.Delete(&Alias{}, "domain = ?", name).Error; err != nil {
			return err
		}
		if err := tx.Delete(&DomainAdmin{}, "domain = ?", name).Error; err != nil {
			return err
		}
		return tx.Delete(&Domain{}, "domain = ?", name).Error
	})
	if err != nil {
		return fail(c, http.StatusInternalServerError, "delete domain failed")
	}
	return ok(c, map[string]string{"deleted": name})
}
