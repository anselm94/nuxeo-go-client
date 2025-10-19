package nuxeo

import "testing"

func TestNewUser(t *testing.T) {
	props := map[string]any{"emailVerified": true}
	user := NewUser("jdoe", "jdoe@example.com", props)
	if user.Username != "jdoe" {
		t.Errorf("Username got %q, want %q", user.Username, "jdoe")
	}
	if user.Email != "jdoe@example.com" {
		t.Errorf("Email got %q, want %q", user.Email, "jdoe@example.com")
	}
	if user.Properties["emailVerified"] != true {
		t.Errorf("Properties[emailVerified] got %v, want %v", user.Properties["emailVerified"], true)
	}
}
