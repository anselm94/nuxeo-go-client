package nuxeo

import "testing"

func TestNewGroup(t *testing.T) {
	props := map[string]any{"role": "admin"}
	group := NewGroup("admins", props)
	if group.Name != "admins" {
		t.Errorf("Name got %q, want %q", group.Name, "admins")
	}
	if group.Properties["role"] != "admin" {
		t.Errorf("Properties[role] got %v, want %v", group.Properties["role"], "admin")
	}
}
