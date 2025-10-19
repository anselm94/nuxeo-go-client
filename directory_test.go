package nuxeo

import "testing"

func TestNewDirectory(t *testing.T) {
	dir := NewDirectory("userDirectory")
	if dir.Name != "userDirectory" {
		t.Errorf("Name got %q, want %q", dir.Name, "userDirectory")
	}
}

func TestNewDirectoryEntry(t *testing.T) {
	props := map[string]any{"label": "Manager"}
	entry := NewDirectoryEntry("mgr", props)
	if entry.ID != "mgr" {
		t.Errorf("ID got %q, want %q", entry.ID, "mgr")
	}
	if entry.Properties["label"] != "Manager" {
		t.Errorf("Properties[label] got %v, want %v", entry.Properties["label"], "Manager")
	}
}
