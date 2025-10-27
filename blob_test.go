package nuxeo

import (
	"io"
	"strings"
	"testing"
)

// dummyReadCloser is a simple io.ReadCloser for testing
// It wraps a strings.Reader and implements Close as a no-op

type dummyReadCloser struct {
	io.Reader
}

func (d *dummyReadCloser) Close() error { return nil }

func TestNewBlob(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		filename string
		mimeType string
		length   int64
		stream   io.ReadCloser
		want     blob
	}{
		{
			name:     "basic",
			filename: "file.txt",
			mimeType: "text/plain",
			length:   123,
			stream:   &dummyReadCloser{strings.NewReader("data")},
			want: blob{
				Filename: "file.txt",
				MimeType: "text/plain",
				Length:   "123",
				Stream:   &dummyReadCloser{strings.NewReader("data")},
			},
		},
		{
			name:     "zero length",
			filename: "empty.bin",
			mimeType: "application/octet-stream",
			length:   0,
			stream:   nil,
			want: blob{
				Filename: "empty.bin",
				MimeType: "application/octet-stream",
				Length:   "0",
				Stream:   nil,
			},
		},
		{
			name:     "negative length",
			filename: "bad.bin",
			mimeType: "application/octet-stream",
			length:   -42,
			stream:   nil,
			want: blob{
				Filename: "bad.bin",
				MimeType: "application/octet-stream",
				Length:   "-42",
				Stream:   nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := NewBlob(tc.filename, tc.mimeType, tc.length, tc.stream)
			if b.Filename != tc.want.Filename {
				t.Errorf("Filename: got %q, want %q", b.Filename, tc.want.Filename)
			}
			if b.MimeType != tc.want.MimeType {
				t.Errorf("MimeType: got %q, want %q", b.MimeType, tc.want.MimeType)
			}
			if b.Length != tc.want.Length {
				t.Errorf("Length: got %q, want %q", b.Length, tc.want.Length)
			}
			if (b.Stream == nil) != (tc.want.Stream == nil) {
				t.Errorf("Stream: got %v, want %v", b.Stream, tc.want.Stream)
			}
		})
	}
}
