package nuxeo

import "io"

// Blob represents a binary object in Nuxeo.
type Blob struct {
	Filename string        `json:"name"`
	MimeType string        `json:"mime-type"`
	Length   string        `json:"length"`
	Stream   io.ReadCloser `json:"-"`

	// (Present only as Document property blob) Encoding
	Encoding string `json:"encoding,omitempty"`
	// (Present only as Document property blob) Digest Algorithm
	DigestAlgorithm string `json:"digestAlgorithm"`
	// (Present only as Document property blob) Digest
	Digest string `json:"digest"`
	// (Present only as Document property blob) Data URL
	Data string `json:"data"`
	// (Present only as Document property blob) Blob URL
	BlobUrl string `json:"blobUrl"`
}

type Blobs entities[Blob]
