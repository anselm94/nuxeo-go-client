package nuxeo

import "testing"

func TestNewDocument_Properties(t *testing.T) {
	doc := NewDocument("123", "File", "/default-domain/workspace", map[string]any{"dc:title": "Test"})
	if doc.ID != "123" {
		t.Errorf("ID got %q, want %q", doc.ID, "123")
	}
	if doc.Type != "File" {
		t.Errorf("Type got %q, want %q", doc.Type, "File")
	}
	if doc.Path != "/default-domain/workspace" {
		t.Errorf("Path got %q, want %q", doc.Path, "/default-domain/workspace")
	}
	if doc.GetProperty("dc:title") != "Test" {
		t.Errorf("GetProperty got %v, want %v", doc.GetProperty("dc:title"), "Test")
	}
	doc.SetProperty("dc:description", "A doc")
	if doc.GetProperty("dc:description") != "A doc" {
		t.Errorf("SetProperty/GetProperty got %v, want %v", doc.GetProperty("dc:description"), "A doc")
	}
}
