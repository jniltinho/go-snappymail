package ui

import (
	"slices"
	"strings"
)

// Skin describes a webmail UI skin exposed to the SPA and config.toml.
type Skin struct {
	ID      string   `json:"id"`
	Label   string   `json:"label"`
	Ready   bool     `json:"ready"`
	Aliases []string `json:"-"`
}

// catalog is the authoritative skin list (server + docs + validate-skins.sh).
// To add a skin: make new-skin ID=<name> --register
//
// catalog-begin
var catalog = []Skin{
	{
		ID:      "snappymail",
		Label:   "SnappyMail",
		Ready:   true,
		Aliases: []string{"default", "snappymail-default", "theme-default"},
	},
	{
		ID:      "gmail",
		Label:   "Gmail",
		Ready:   false,
		Aliases: []string{"google"},
	},
	{
		ID:      "outlook",
		Label:   "Outlook",
		Ready:   true,
		Aliases: []string{"office", "microsoft"},
	},
	{
		ID:      "carbonio",
		Label:   "Carbonio",
		Ready:   true,
		Aliases: []string{"zextras"},
	},
} // catalog-end

const defaultSkinID = "snappymail"

// Catalog returns skin metadata for API consumers.
func Catalog() []Skin {
	out := make([]Skin, len(catalog))
	copy(out, catalog)
	return out
}

// AvailableSkins returns skin ids only (backward compatible).
func AvailableSkins() []string {
	ids := make([]string, len(catalog))
	for i, s := range catalog {
		ids[i] = s.ID
	}
	return ids
}

// NormalizeSkin maps config.toml values and aliases to a catalog id.
// Unknown values fall back to the default skin.
func NormalizeSkin(raw string) string {
	s := strings.ToLower(strings.TrimSpace(raw))
	if s == "" {
		return defaultSkinID
	}
	for _, skin := range catalog {
		if s == skin.ID || slices.Contains(skin.Aliases, s) {
			return skin.ID
		}
	}
	return defaultSkinID
}

// SkinByID returns catalog entry or false.
func SkinByID(id string) (Skin, bool) {
	id = NormalizeSkin(id)
	for _, s := range catalog {
		if s.ID == id {
			return s, true
		}
	}
	return Skin{}, false
}
