package nuxeo

// Blob represents a binary object in Nuxeo.
type Blob struct {
	Filename string
	MimeType string
	Length   int64
	Data     []byte
}

// NewBlob creates a new Blob instance.
func NewBlob(filename, mimeType string, data []byte) *Blob {
	return &Blob{
		Filename: filename,
		MimeType: mimeType,
		Length:   int64(len(data)),
		Data:     data,
	}
}
