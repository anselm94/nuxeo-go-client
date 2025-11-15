package nuxeo

import (
	"bytes"
	"context"
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
		resp, err := manager.Execute(context.Background(), *op, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatalf("expected response, got nil")
		}
		var out map[string]string
		if err := resp.As(&out); err != nil {
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
		resp, err := manager.Execute(context.Background(), *op, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp != nil {
			t.Errorf("expected nil body for 204, got %v", resp)
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

	t.Run("malformed JSON", func(t *testing.T) {
		t.Parallel()
		client := newMockNuxeoClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("{not-json}")),
			}, nil
		})
		manager := &operationManager{client: client, logger: slog.Default()}
		op := mockOperation()
		resp, err := manager.Execute(context.Background(), *op, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var out map[string]string
		err = resp.As(&out)
		if err == nil {
			t.Errorf("expected decode error, got nil")
		}
	})
}

// ...existing code...

func TestOperationManager_Execute_Multipart(t *testing.T) {
	t.Run("multipart with one blob", func(t *testing.T) {
		t.Parallel()
		client := newMockNuxeoClient(func(req *http.Request) (*http.Response, error) {
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
		resp, err := manager.Execute(context.Background(), *op, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatalf("expected response, got nil")
		}
		var out map[string]any
		if err := resp.As(&out); err != nil {
			t.Fatalf("decode error: %v", err)
		}
		if out["ok"] != true {
			t.Errorf("expected ok=true, got %v", out)
		}
	})

	t.Run("multipart with multiple blobs", func(t *testing.T) {
		t.Parallel()
		client := newMockNuxeoClient(func(req *http.Request) (*http.Response, error) {
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
		resp, err := manager.Execute(context.Background(), *op, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatalf("expected response, got nil")
		}
		var out map[string]any
		if err := resp.As(&out); err != nil {
			t.Fatalf("decode error: %v", err)
		}
		if out["ok"] != true {
			t.Errorf("expected ok=true, got %v", out)
		}
	})

	t.Run("returns nil for 204", func(t *testing.T) {
		t.Parallel()
		client := newMockNuxeoClient(func(req *http.Request) (*http.Response, error) {
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
		resp, err := manager.Execute(context.Background(), *op, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var respBody bytes.Buffer
		io.Copy(&respBody, resp.res.Body)
		if respBody.Len() != 0 {
			t.Errorf("expected nil body for 204, got %v", resp)
		}
	})

	t.Run("error response", func(t *testing.T) {
		t.Parallel()
		client := newMockNuxeoClient(func(req *http.Request) (*http.Response, error) {
			_, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			return &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(strings.NewReader(`{"ok":false}`)),
			}, errors.New("fail")
		})
		manager := &operationManager{client: client, logger: slog.Default()}
		blob1 := NewBlob("file.txt", "text/plain", 160, io.NopCloser(strings.NewReader("hello")))
		op := NewOperation("Blob.Attach").SetInputBlob(*blob1)
		_, err := manager.Execute(context.Background(), *op, nil)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}
