package nuxeo

import "testing"

func TestNewBatchUpload_AddBlob(t *testing.T) {
	batch := NewBatchUpload("batch1")
	blob := NewBlob("file.txt", "text/plain", []byte("data"))
	batch.AddBlob(blob)
	if batch.ID != "batch1" {
		t.Errorf("ID got %q, want %q", batch.ID, "batch1")
	}
	if len(batch.Blobs) != 1 {
		t.Errorf("Blobs len got %d, want 1", len(batch.Blobs))
	}
	if batch.Blobs[0].ID != "file.txt" {
		t.Errorf("BatchBlob ID got %q, want %q", batch.Blobs[0].ID, "file.txt")
	}
	if string(batch.Blobs[0].Blob.Data) != "data" {
		t.Errorf("BatchBlob Data got %q, want %q", string(batch.Blobs[0].Blob.Data), "data")
	}
}
