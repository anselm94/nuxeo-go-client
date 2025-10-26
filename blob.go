package nuxeo

import "io"

// Blob represents a binary object in Nuxeo, typically used for file uploads and document properties.
// Fields map to Nuxeo's Blob JSON structure.
//
// - Filename: Name of the file
// - MimeType: MIME type of the file
// - Length: Size of the file in bytes
// - Stream: File data as io.ReadCloser (not serialized)
// - Encoding, DigestAlgorithm, Digest, Data, BlobUrl: Only present when Blob is a document property
//
// Used for uploading files, retrieving blobs from documents, and batch upload operations.
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

// Blobs is a slice of Blob objects, used for representing multiple blobs in Nuxeo responses.
type Blobs entities[Blob]
