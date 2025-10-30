package nuxeo

import (
	"io"
	"strconv"
)

// blob represents a binary object in Nuxeo, typically used for file uploads and document properties.
// Fields map to Nuxeo's blob JSON structure.
//
// - Filename: Name of the file
// - MimeType: MIME type of the file
// - Length: Size of the file in bytes
// - Stream: File data as io.ReadCloser (not serialized)
// - Encoding, DigestAlgorithm, Digest, Data, BlobUrl: Only present when blob is a document property
//
// Used for uploading files, retrieving blobs from documents, and batch upload operations.
type blob struct {
	io.ReadCloser
	Filename string `json:"name"`
	MimeType string `json:"mime-type"`
	Length   string `json:"length"`

	// (Readonly) Encoding
	Encoding string `json:"encoding,omitempty"`
	// (Readonly) Digest Algorithm
	DigestAlgorithm string `json:"digestAlgorithm"`
	// (Readonly) Digest
	Digest string `json:"digest"`
	// (Readonly) Data URL
	Data string `json:"data"`
	// (Readonly) Blob URL
	BlobUrl string `json:"blobUrl"`
}

// NewBlob creates a new Blob instance with the specified filename, MIME type, length, and data stream.
func NewBlob(filename, mimeType string, length int64, stream io.ReadCloser) *blob {
	return &blob{
		ReadCloser: stream,
		Filename:   filename,
		MimeType:   mimeType,
		Length:     strconv.FormatInt(length, 10),
	}
}

func (b *blob) Size() int64 {
	size, err := strconv.ParseInt(b.Length, 10, 64)
	if err != nil {
		return 0
	}
	return size
}
