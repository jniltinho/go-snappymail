package handler

import "testing"

func TestSplitAddrs(t *testing.T) {
	tests := []struct {
		in   string
		want []string
	}{
		{"", nil},
		{"a@x.com", []string{"a@x.com"}},
		{"a@x.com, b@y.com", []string{"a@x.com", "b@y.com"}},
		{"  spaced@test.local  ", []string{"spaced@test.local"}},
	}

	for _, tt := range tests {
		got := splitAddrs(tt.in)
		if len(got) != len(tt.want) {
			t.Fatalf("splitAddrs(%q) = %v, want %v", tt.in, got, tt.want)
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Fatalf("splitAddrs(%q)[%d] = %q, want %q", tt.in, i, got[i], tt.want[i])
			}
		}
	}
}
