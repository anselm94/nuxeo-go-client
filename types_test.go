package nuxeo

import "testing"

func TestDocumentTypeConstants(t *testing.T) {
	if DocumentTypeFile != "File" {
		t.Errorf("DocumentTypeFile got %q, want %q", DocumentTypeFile, "File")
	}
	if DocumentTypeFolder != "Folder" {
		t.Errorf("DocumentTypeFolder got %q, want %q", DocumentTypeFolder, "Folder")
	}
	if DocumentTypeWorkspace != "Workspace" {
		t.Errorf("DocumentTypeWorkspace got %q, want %q", DocumentTypeWorkspace, "Workspace")
	}
}
