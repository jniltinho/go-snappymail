package session

import (
	"testing"
	"time"
)

func resetSessionStore(t *testing.T) {
	t.Helper()
	storeMu.Lock()
	store = map[string]*IMAPSession{}
	storeMu.Unlock()
}

func TestEncryptDecryptPassword(t *testing.T) {
	t.Parallel()

	key := "test-secret-key-32-chars!!"
	tests := []struct {
		name      string
		plaintext string
	}{
		{name: "simple password", plaintext: "Password1@"},
		{name: "empty password", plaintext: ""},
		{name: "unicode password", plaintext: "päss🔑"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			enc, err := EncryptPassword(tt.plaintext, key)
			if err != nil {
				t.Fatalf("EncryptPassword() error = %v", err)
			}
			if enc == "" && tt.plaintext != "" {
				t.Fatal("expected non-empty ciphertext")
			}
			got, err := DecryptPassword(enc, key)
			if err != nil {
				t.Fatalf("DecryptPassword() error = %v", err)
			}
			if got != tt.plaintext {
				t.Fatalf("roundtrip = %q, want %q", got, tt.plaintext)
			}
		})
	}
}

func TestDecryptPasswordWrongKey(t *testing.T) {
	enc, err := EncryptPassword("secret", "key-one-32-characters-long!!!!!")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := DecryptPassword(enc, "key-two-32-characters-long!!!!!"); err == nil {
		t.Fatal("expected error with wrong key")
	}
}

func TestSessionSetGetDelete(t *testing.T) {
	resetSessionStore(t)

	s := &IMAPSession{
		IMAPHost: "mail.example.com",
		IMAPPort: 993,
		Username: "user@example.com",
	}
	if err := s.SetPassword("pass", "test-secret-key-32-chars!!"); err != nil {
		t.Fatal(err)
	}

	Set("sess-1", s)
	got := Get("sess-1")
	if got == nil || got.Username != "user@example.com" {
		t.Fatalf("Get() = %+v", got)
	}

	pass, err := got.Password("test-secret-key-32-chars!!")
	if err != nil || pass != "pass" {
		t.Fatalf("Password() = %q, err = %v", pass, err)
	}

	Delete("sess-1")
	if Get("sess-1") != nil {
		t.Fatal("expected session to be deleted")
	}
}

func TestCleanupRemovesIdleSessions(t *testing.T) {
	resetSessionStore(t)

	s := &IMAPSession{Username: "idle@example.com", LastUsed: time.Now().Add(-2 * time.Hour)}
	storeMu.Lock()
	store["idle"] = s
	storeMu.Unlock()

	Cleanup(time.Hour)
	if Get("idle") != nil {
		t.Fatal("expected idle session to be removed")
	}
}
