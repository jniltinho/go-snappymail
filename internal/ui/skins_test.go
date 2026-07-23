package ui

import "testing"

func TestNormalizeSkin(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"", SkinSnappyMail},
		{"snappymail", SkinSnappyMail},
		{"snappymail-default", SkinSnappyMail},
		{"Gmail", SkinGmail},
		{"outlook", SkinOutlook},
		{"microsoft", SkinOutlook},
		{"unknown-brand", SkinSnappyMail},
	}
	for _, tt := range tests {
		if got := NormalizeSkin(tt.in); got != tt.want {
			t.Fatalf("NormalizeSkin(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
