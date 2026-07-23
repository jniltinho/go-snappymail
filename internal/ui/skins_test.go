package ui

import "testing"

func TestNormalizeSkin(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"", "snappymail"},
		{"snappymail", "snappymail"},
		{"snappymail-default", "snappymail"},
		{"Gmail", "gmail"},
		{"google", "gmail"},
		{"outlook", "outlook"},
		{"microsoft", "outlook"},
		{"unknown-brand", "snappymail"},
	}
	for _, tt := range tests {
		if got := NormalizeSkin(tt.in); got != tt.want {
			t.Fatalf("NormalizeSkin(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestCatalogMatchesAvailableSkins(t *testing.T) {
	cat := Catalog()
	ids := AvailableSkins()
	if len(cat) != len(ids) {
		t.Fatalf("catalog len %d != ids len %d", len(cat), len(ids))
	}
	for i, s := range cat {
		if s.ID != ids[i] {
			t.Fatalf("catalog[%d].ID = %q, ids[%d] = %q", i, s.ID, i, ids[i])
		}
		if s.Label == "" {
			t.Fatalf("catalog[%d] missing label", i)
		}
	}
}

func TestSkinByID(t *testing.T) {
	s, ok := SkinByID("gmail")
	if !ok || s.Label != "Gmail" {
		t.Fatalf("SkinByID(gmail) = %+v, %v", s, ok)
	}
}
