package ui

import "strings"

// Known skin identifiers (layout + visual identity).
// To add a skin: docs/skins.md  ·  make new-skin ID=<name>
const (
	SkinSnappyMail = "snappymail"
	SkinGmail      = "gmail"
	SkinOutlook    = "outlook"
)

var available = []string{SkinSnappyMail, SkinGmail, SkinOutlook}

// AvailableSkins returns skin IDs the frontend may offer (some may be stubs until implemented).
func AvailableSkins() []string {
	out := make([]string, len(available))
	copy(out, available)
	return out
}

// NormalizeSkin maps config values to a known skin id; unknown values fall back to snappymail.
func NormalizeSkin(raw string) string {
	s := strings.ToLower(strings.TrimSpace(raw))
	switch s {
	case "", "snappymail", "snappymail-default", "default":
		return SkinSnappyMail
	case "gmail", "google":
		return SkinGmail
	case "outlook", "office", "microsoft":
		return SkinOutlook
	default:
		return SkinSnappyMail
	}
}
