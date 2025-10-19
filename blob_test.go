package nuxeo

import "testing"

func TestNewBlob(t *testing.T) {
	data := []byte("hello world")
	blob := NewBlob("file.txt", "text/plain", data)
	if blob.Filename != "file.txt" {
		t.Errorf("Filename got %q, want %q", blob.Filename, "file.txt")
	}
	if blob.MimeType != "text/plain" {
		t.Errorf("MimeType got %q, want %q", blob.MimeType, "text/plain")
	}
	if blob.Length != int64(len(data)) {
		t.Errorf("Length got %d, want %d", blob.Length, len(data))
	}
	if string(blob.Data) != "hello world" {
		t.Errorf("Data got %q, want %q", string(blob.Data), "hello world")
	}
}
