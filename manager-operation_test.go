package nuxeo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"
)

// mockOperation returns a simple operation for testing.
func mockOperation() *operation {
	return NewOperation("Document.Fetch").SetInputDocumentId("123")
}

func TestOperationManager_ExecuteInto(t *testing.T) {
	type output struct {
		Foo string `json:"foo"`
	}
	tests := []struct {
		name        string
		mockBody    any
		mockStatus  int
		mockErr     error
		expectErr   bool
		expectValue string
	}{
		{
			name:        "success JSON decode",
			mockBody:    output{Foo: "bar"},
			mockStatus:  200,
			expectErr:   false,
			expectValue: "bar",
		},
		{
			name:       "error from Execute",
			mockBody:   nil,
			mockStatus: 500,
			mockErr:    errors.New("network error"),
			expectErr:  true,
		},
		{
			name:       "malformed JSON",
			mockBody:   "{not-json}",
			mockStatus: 200,
			expectErr:  true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := newMockNuxeoClient(func(req *http.Request) (*http.Response, error) {
				var body io.ReadCloser
				switch v := tt.mockBody.(type) {
				case string:
					body = io.NopCloser(strings.NewReader(v))
				case nil:
					body = io.NopCloser(bytes.NewBuffer(nil))
				default:
					body = testMarshalBody(t, v)
				}
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       body,
				}, tt.mockErr
			})
			manager := &operationManager{client: client, logger: slog.Default()}
			var out output
			err := manager.ExecuteInto(context.Background(), *mockOperation(), &out, nil)
			if tt.expectErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectErr && out.Foo != tt.expectValue {
				t.Errorf("expected Foo=%q, got %q", tt.expectValue, out.Foo)
			}
		})
	}
}

func TestOperationManager_Execute_JSON(t *testing.T) {
	t.Run("returns body for 200", func(t *testing.T) {
		t.Parallel()
		client := newMockNuxeoClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`{"foo":"bar"}`)),
			}, nil
		})
		manager := &operationManager{client: client, logger: slog.Default()}
		op := mockOperation()
		rc, err := manager.Execute(context.Background(), *op, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer rc.Close()
		var out map[string]string
		if err := json.NewDecoder(rc).Decode(&out); err != nil {
			t.Fatalf("decode error: %v", err)
		}
		if out["foo"] != "bar" {
			t.Errorf("expected foo=bar, got %v", out)
		}
	})

	t.Run("returns nil for 204", func(t *testing.T) {
		t.Parallel()
		client := newMockNuxeoClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 204,
				Body:       io.NopCloser(bytes.NewBuffer(nil)),
			}, nil
		})
		manager := &operationManager{client: client, logger: slog.Default()}
		op := mockOperation().SetVoidOperation(true)
		rc, err := manager.Execute(context.Background(), *op, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rc != nil {
			t.Errorf("expected nil body for 204, got %v", rc)
		}
	})

	t.Run("error response", func(t *testing.T) {
		t.Parallel()
		client := newMockNuxeoClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(strings.NewReader(`{"error":"fail"}`)),
			}, errors.New("fail")
		})
		manager := &operationManager{client: client, logger: slog.Default()}
		op := mockOperation()
		_, err := manager.Execute(context.Background(), *op, nil)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}

func TestOperationManager_Execute_Multipart(t *testing.T) {
	t.Run("multipart with one blob", func(t *testing.T) {
		t.Parallel()
		client := newMockNuxeoClient(func(req *http.Request) (*http.Response, error) {
			// read and discard the multipart body for test purposes
			_, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
			}, nil
		})
		manager := &operationManager{client: client, logger: slog.Default()}
		blob := NewBlob("file.txt", "text/plain", 160, io.NopCloser(strings.NewReader("hello")))
		op := NewOperation("Blob.Attach").SetInputBlob(*blob)
		var out map[string]any
		err := manager.ExecuteInto(context.Background(), *op, &out, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if out["ok"] != true {
			t.Errorf("expected ok=true, got %v", out)
		}
	})

	t.Run("multipart with multiple blobs", func(t *testing.T) {
		t.Parallel()
		client := newMockNuxeoClient(func(req *http.Request) (*http.Response, error) {
			// read and discard the multipart body for test purposes
			_, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
			}, nil
		})
		manager := &operationManager{client: client, logger: slog.Default()}
		blob1 := NewBlob("file.txt", "text/plain", 160, io.NopCloser(strings.NewReader("hello")))
		blob2 := NewBlob("image.png", "image/png", 2048, io.NopCloser(strings.NewReader("...binary data...")))
		op := NewOperation("Blob.Attach").SetInputBlobs(*blob1, *blob2)
		var out map[string]any
		err := manager.ExecuteInto(context.Background(), *op, &out, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if out["ok"] != true {
			t.Errorf("expected ok=true, got %v", out)
		}
	})

	t.Run("returns nil for 204", func(t *testing.T) {
		t.Parallel()
		client := newMockNuxeoClient(func(req *http.Request) (*http.Response, error) {
			// read and discard the multipart body for test purposes
			_, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			return &http.Response{
				StatusCode: 204,
				Body:       io.NopCloser(bytes.NewBuffer(nil)),
			}, nil
		})
		manager := &operationManager{client: client, logger: slog.Default()}
		blob1 := NewBlob("file.txt", "text/plain", 160, io.NopCloser(strings.NewReader("hello")))
		op := mockOperation().SetVoidOperation(true).SetInputBlob(*blob1)
		rc, err := manager.Execute(context.Background(), *op, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rc != nil {
			t.Errorf("expected nil body for 204, got %v", rc)
		}
	})

	t.Run("error response", func(t *testing.T) {
		t.Parallel()
		client := newMockNuxeoClient(func(req *http.Request) (*http.Response, error) {
			// read and discard the multipart body for test purposes
			_, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			return &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(strings.NewReader(`{"ok":false}`)),
			}, nil
		})
		manager := &operationManager{client: client, logger: slog.Default()}
		blob1 := NewBlob("file.txt", "text/plain", 160, io.NopCloser(strings.NewReader("hello")))
		op := NewOperation("Blob.Attach").SetInputBlob(*blob1)
		var out map[string]any
		err := manager.ExecuteInto(context.Background(), *op, &out, nil)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if out["ok"] == true {
			t.Errorf("expected ok=true, got %v", out)
		}
	})
}

func TestOperationManager_FetchOperation(t *testing.T) {
	t.Run("returns nil (not implemented)", func(t *testing.T) {
		t.Parallel()
		client := newMockNuxeoClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`{}`)),
			}, nil
		})
		manager := &operationManager{client: client, logger: slog.Default()}
		op, err := manager.FetchOperation(context.Background(), "Document.Fetch")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if op != nil {
			t.Errorf("expected nil, got %v", op)
		}
	})
}
