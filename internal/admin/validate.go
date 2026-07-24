package admin

import (
	"regexp"
	"strings"
)

// domainRe and emailRe are pragmatic validators — enough to reject obviously
// invalid input at the API boundary without pretending to fully parse RFC 5321.
var (
	domainRe = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$`)
	emailRe  = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$`)
)

func validDomain(s string) bool {
	return len(s) <= 255 && domainRe.MatchString(s)
}

func validEmail(s string) bool {
	return len(s) <= 255 && emailRe.MatchString(s)
}

// domainOf returns the domain part of an email address (lowercased), or "".
func domainOf(email string) string {
	i := strings.LastIndexByte(email, '@')
	if i < 0 {
		return ""
	}
	return strings.ToLower(email[i+1:])
}

// localPartOf returns the local part of an email address, or "".
func localPartOf(email string) string {
	i := strings.LastIndexByte(email, '@')
	if i < 0 {
		return ""
	}
	return email[:i]
}
