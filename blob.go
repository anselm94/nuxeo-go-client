package nuxeo

import "io"

// Blob represents a binary object in Nuxeo.
type Blob struct {
	Filename string        `json:"name"`
	MimeType string        `json:"mime-type"`
	Length   int64         `json:"length"`
	Stream   io.ReadCloser `json:"-"`

	Encoding        string `json:"encoding,omitempty"`
	DigestAlgorithm string `json:"digestAlgorithm"`
	Digest          string `json:"digest"`
	Data            string `json:"data"`
	BlobUrl         string `json:"blobUrl"`
}

type Blobs entities[Blob]
