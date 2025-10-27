package nuxeo

import (
	"testing"
)

func TestNewACP(t *testing.T) {
	t.Parallel()
	acp := NewACP()
	if acp == nil {
		t.Fatal("NewACP returned nil")
	}
	if acp.entity.EntityType != EntityTypeACP {
		t.Errorf("EntityType: got %q, want %q", acp.entity.EntityType, EntityTypeACP)
	}
	if len(acp.ACLs) != 0 {
		t.Errorf("Expected empty ACLs, got %d", len(acp.ACLs))
	}
}

func TestNewACL(t *testing.T) {
	t.Parallel()
	acl := NewACL("local")
	if acl == nil {
		t.Fatal("NewACL returned nil")
	}
	if acl.Name != "local" {
		t.Errorf("Name: got %q, want %q", acl.Name, "local")
	}
	if len(acl.ACEs) != 0 {
		t.Errorf("Expected empty ACEs, got %d", len(acl.ACEs))
	}
}

func TestNewACE(t *testing.T) {
	t.Parallel()
	ace := NewACE("bob", "Read", true)
	if ace == nil {
		t.Fatal("NewACE returned nil")
	}
	if ace.Username != "bob" {
		t.Errorf("Username: got %q, want %q", ace.Username, "bob")
	}
	if ace.Permission != "Read" {
		t.Errorf("Permission: got %q, want %q", ace.Permission, "Read")
	}
	if ace.Granted != true {
		t.Errorf("Granted: got %v, want true", ace.Granted)
	}
	// Other fields should be zero values
	if ace.ID != "" {
		t.Errorf("ID: got %q, want empty", ace.ID)
	}
	if ace.Creator != "" {
		t.Errorf("Creator: got %q, want empty", ace.Creator)
	}
	if ace.Begin != nil {
		t.Errorf("Begin: got %v, want nil", ace.Begin)
	}
	if ace.End != nil {
		t.Errorf("End: got %v, want nil", ace.End)
	}
	if ace.Status != "" {
		t.Errorf("Status: got %q, want empty", ace.Status)
	}
}
