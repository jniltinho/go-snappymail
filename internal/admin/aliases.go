package admin

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// aliasWriteRequest is the JSON body for create/update alias. Goto is a
// comma-separated list of destination addresses.
type aliasWriteRequest struct {
	Address string `json:"address"`
	Goto    string `json:"goto"`
	Active  *bool  `json:"active"`
}

// validGoto checks that every destination in a comma-separated Goto is a valid
// email address and that the list is non-empty.
func validGoto(gotoList string) bool {
	parts := strings.Split(gotoList, ",")
	n := 0
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if !validEmail(p) {
			return false
		}
		n++
	}
	return n > 0
}

// ListAliases returns aliases visible to the caller (scoped by domain).
func (h *Handlers) ListAliases(c *echo.Context) error {
	cl, ok2 := claimsFrom(c)
	if !ok2 {
		return fail(c, http.StatusUnauthorized, "not authenticated")
	}
	q := h.db.Model(&Alias{})
	if d := c.QueryParam("domain"); d != "" {
		if !cl.Can(PermAliasesRead, d) {
			return fail(c, http.StatusForbidden, "forbidden")
		}
		q = q.Where("domain = ?", d)
	} else if !cl.Superadmin {
		if len(cl.Domains) == 0 {
			return ok(c, []Alias{})
		}
		q = q.Where("domain IN ?", cl.Domains)
	}
	var aliases []Alias
	if err := q.Order("address").Find(&aliases).Error; err != nil {
		return fail(c, http.StatusInternalServerError, "list aliases failed")
	}
	return ok(c, aliases)
}

// GetAlias returns a single alias (scoped by the row's domain column).
func (h *Handlers) GetAlias(c *echo.Context) error {
	address := c.Param("address")
	claims, ok2 := claimsFrom(c)
	if !ok2 {
		return fail(c, http.StatusUnauthorized, "not authenticated")
	}
	var a Alias
	if err := h.db.First(&a, "address = ?", address).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fail(c, http.StatusNotFound, "alias not found")
		}
		return fail(c, http.StatusInternalServerError, "get alias failed")
	}
	if !claims.Can(PermAliasesRead, a.Domain) {
		return fail(c, http.StatusForbidden, "forbidden")
	}
	return ok(c, a)
}

// CreateAlias creates an alias under an existing domain the caller manages.
func (h *Handlers) CreateAlias(c *echo.Context) error {
	var req aliasWriteRequest
	if err := c.Bind(&req); err != nil {
		return fail(c, http.StatusBadRequest, "invalid request body")
	}
	if !validEmail(req.Address) {
		return fail(c, http.StatusBadRequest, "invalid alias address")
	}
	if !validGoto(req.Goto) {
		return fail(c, http.StatusBadRequest, "invalid destination(s)")
	}
	dom := domainOf(req.Address)
	if _, allowed := requirePermission(c, PermAliasesWrite, dom); !allowed {
		return nil
	}
	var dcount int64
	if err := h.db.Model(&Domain{}).Where("domain = ?", dom).Count(&dcount).Error; err != nil {
		return fail(c, http.StatusInternalServerError, "create alias failed")
	}
	if dcount == 0 {
		return fail(c, http.StatusBadRequest, "domain does not exist")
	}
	var acount int64
	h.db.Model(&Alias{}).Where("address = ?", req.Address).Count(&acount)
	if acount > 0 {
		return fail(c, http.StatusConflict, "alias already exists")
	}
	a := Alias{Address: req.Address, Goto: req.Goto, Domain: dom, Active: true}
	if req.Active != nil {
		a.Active = *req.Active
	}
	if err := h.db.Create(&a).Error; err != nil {
		return fail(c, http.StatusInternalServerError, "create alias failed")
	}
	return created(c, a)
}

// UpdateAlias updates an alias (scoped).
func (h *Handlers) UpdateAlias(c *echo.Context) error {
	address := c.Param("address")
	claims, ok2 := claimsFrom(c)
	if !ok2 {
		return fail(c, http.StatusUnauthorized, "not authenticated")
	}
	var req aliasWriteRequest
	if err := c.Bind(&req); err != nil {
		return fail(c, http.StatusBadRequest, "invalid request body")
	}
	if req.Goto != "" && !validGoto(req.Goto) {
		return fail(c, http.StatusBadRequest, "invalid destination(s)")
	}
	var a Alias
	if err := h.db.First(&a, "address = ?", address).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fail(c, http.StatusNotFound, "alias not found")
		}
		return fail(c, http.StatusInternalServerError, "update alias failed")
	}
	if !claims.Can(PermAliasesWrite, a.Domain) {
		return fail(c, http.StatusForbidden, "forbidden")
	}
	updates := map[string]any{}
	if req.Goto != "" {
		updates["goto"] = req.Goto
	}
	if req.Active != nil {
		updates["active"] = *req.Active
	}
	if len(updates) > 0 {
		if err := h.db.Model(&a).Updates(updates).Error; err != nil {
			return fail(c, http.StatusInternalServerError, "update alias failed")
		}
	}
	h.db.First(&a, "address = ?", address)
	return ok(c, a)
}

// DeleteAlias removes an alias (scoped).
func (h *Handlers) DeleteAlias(c *echo.Context) error {
	address := c.Param("address")
	claims, ok2 := claimsFrom(c)
	if !ok2 {
		return fail(c, http.StatusUnauthorized, "not authenticated")
	}
	var a Alias
	if err := h.db.First(&a, "address = ?", address).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fail(c, http.StatusNotFound, "alias not found")
		}
		return fail(c, http.StatusInternalServerError, "delete alias failed")
	}
	if !claims.Can(PermAliasesWrite, a.Domain) {
		return fail(c, http.StatusForbidden, "forbidden")
	}
	if err := h.db.Delete(&Alias{}, "address = ?", address).Error; err != nil {
		return fail(c, http.StatusInternalServerError, "delete alias failed")
	}
	return ok(c, map[string]string{"deleted": address})
}
