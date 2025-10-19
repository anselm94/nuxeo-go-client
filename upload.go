package nuxeo

// BatchUpload represents a batch upload session.
type BatchUpload struct {
	ID    string
	Blobs []*BatchBlob
}

// BatchBlob represents a blob in a batch upload.
type BatchBlob struct {
	ID   string
	Blob *Blob
}

// NewBatchUpload creates a new BatchUpload instance.
func NewBatchUpload(id string) *BatchUpload {
	return &BatchUpload{
		ID:    id,
		Blobs: []*BatchBlob{},
	}
}

// AddBlob adds a blob to the batch upload.
func (b *BatchUpload) AddBlob(blob *Blob) {
	batchBlob := &BatchBlob{ID: blob.Filename, Blob: blob}
	b.Blobs = append(b.Blobs, batchBlob)
}
