package nuxeo

import "io"

// Blob represents a binary object in Nuxeo.
type Blob struct {
	Filename string
	MimeType string
	Length   int64
	Stream   io.ReadCloser
}
