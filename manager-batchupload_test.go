package nuxeo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"
)

func TestBatchUploadManager_CreateBatch(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		respond   func(req *http.Request) (*http.Response, error)
		wantErr   bool
		wantBatch *batchUpload
	}{
		{
			name: "success",
			respond: func(req *http.Request) (*http.Response, error) {
				body := `{"name":"file1","batchId":"batch1"}`
				return &http.Response{
					StatusCode: 201,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
				}, nil
			},
			wantErr:   false,
			wantBatch: &batchUpload{Name: "file1", BatchId: "batch1"},
		},
		{
			name: "client error",
			respond: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("client error")
			},
			wantErr:   true,
			wantBatch: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newMockNuxeoClient(tc.respond)
			bum := &batchUploadManager{
				client: client,
				logger: slog.Default(),
			}
			batch, err := bum.CreateBatch(context.Background(), &nuxeoRequestOptions{})
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tc.wantBatch != nil && batch != nil {
				if batch.Name != tc.wantBatch.Name || batch.BatchId != tc.wantBatch.BatchId {
					t.Errorf("unexpected batch: got %+v, want %+v", batch, tc.wantBatch)
				}
			} else if tc.wantBatch != batch {
				if tc.wantBatch == nil && batch != nil {
					t.Errorf("expected nil batch, got %+v", batch)
				}
				if tc.wantBatch != nil && batch == nil {
					t.Errorf("expected batch %+v, got nil", tc.wantBatch)
				}
			}
		})
	}
}

func TestBatchUploadManager_FetchBatchUploads(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		respond     func(req *http.Request) (*http.Response, error)
		wantErr     bool
		wantBatches *[]batchUpload
	}{
		{
			name: "success",
			respond: func(req *http.Request) (*http.Response, error) {
				body := `[{"name":"file1","batchId":"batch1"},{"name":"file2","batchId":"batch1"}]`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
				}, nil
			},
			wantErr:     false,
			wantBatches: &[]batchUpload{{Name: "file1", BatchId: "batch1"}, {Name: "file2", BatchId: "batch1"}},
		},
		{
			name: "client error",
			respond: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("client error")
			},
			wantErr:     true,
			wantBatches: nil,
		},
		{
			name: "empty batch",
			respond: func(req *http.Request) (*http.Response, error) {
				body := `[]`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
				}, nil
			},
			wantErr:     false,
			wantBatches: &[]batchUpload{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newMockNuxeoClient(tc.respond)
			bum := &batchUploadManager{
				client: client,
				logger: slog.Default(),
			}
			batches, err := bum.FetchBatchUploads(context.Background(), "batch1", &nuxeoRequestOptions{})
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tc.wantBatches != nil && batches != nil {
				if len(*batches) != len(*tc.wantBatches) {
					t.Errorf("unexpected batch count: got %d, want %d", len(*batches), len(*tc.wantBatches))
				}
				for i := range *batches {
					if (*batches)[i].Name != (*tc.wantBatches)[i].Name || (*batches)[i].BatchId != (*tc.wantBatches)[i].BatchId {
						t.Errorf("unexpected batch at %d: got %+v, want %+v", i, (*batches)[i], (*tc.wantBatches)[i])
					}
				}
			} else if tc.wantBatches != batches {
				if tc.wantBatches == nil && batches != nil {
					t.Errorf("expected nil batches, got %+v", batches)
				}
				if tc.wantBatches != nil && batches == nil {
					t.Errorf("expected batches %+v, got nil", tc.wantBatches)
				}
			}
		})
	}
}

func TestBatchUploadManager_FetchBatchUpload(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		respond   func(req *http.Request) (*http.Response, error)
		wantErr   bool
		wantBatch *batchUpload
	}{
		{
			name: "success",
			respond: func(req *http.Request) (*http.Response, error) {
				body := `{"name":"file1","batchId":"batch1","fileIdx":"0"}`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
			wantErr:   false,
			wantBatch: &batchUpload{Name: "file1", BatchId: "batch1", FileIdx: "0"},
		},
		{
			name: "client error",
			respond: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("client error")
			},
			wantErr:   true,
			wantBatch: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newMockNuxeoClient(tc.respond)
			bum := &batchUploadManager{client: client, logger: slog.Default()}
			batch, err := bum.FetchBatchUpload(context.Background(), "batch1", "0", &nuxeoRequestOptions{})
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tc.wantBatch != nil && batch != nil {
				if batch.Name != tc.wantBatch.Name || batch.BatchId != tc.wantBatch.BatchId || batch.FileIdx != tc.wantBatch.FileIdx {
					t.Errorf("unexpected batch: got %+v, want %+v", batch, tc.wantBatch)
				}
			} else if tc.wantBatch != batch {
				t.Errorf("expected batch %+v, got %+v", tc.wantBatch, batch)
			}
		})
	}
}

func TestBatchUploadManager_CancelBatch(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		respond func(req *http.Request) (*http.Response, error)
		wantErr bool
	}{
		{
			name: "success",
			respond: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 204,
					Body:       io.NopCloser(strings.NewReader("")),
				}, nil
			},
			wantErr: false,
		},
		{
			name: "client error",
			respond: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("client error")
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newMockNuxeoClient(tc.respond)
			bum := &batchUploadManager{client: client, logger: slog.Default()}
			err := bum.CancelBatch(context.Background(), "batch1", &nuxeoRequestOptions{})
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestBatchUploadManager_ExecuteBatchUploads(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		respond    func(req *http.Request) (*http.Response, error)
		wantErr    bool
		wantResult string
	}{
		{
			name: "success",
			respond: func(req *http.Request) (*http.Response, error) {
				body := `"executed"`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
			wantErr:    false,
			wantResult: "executed",
		},
		{
			name: "client error",
			respond: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("client error")
			},
			wantErr:    true,
			wantResult: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newMockNuxeoClient(tc.respond)
			bum := &batchUploadManager{client: client, logger: slog.Default()}
			out := new(string)
			op := operation{operationId: "opId"}
			result, err := bum.ExecuteBatchUploads(context.Background(), "batch1", op, out, &nuxeoRequestOptions{})
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tc.wantErr && result != nil {
				if s, ok := result.(*string); ok {
					if *s != tc.wantResult {
						t.Errorf("unexpected result: got %v, want %v", *s, tc.wantResult)
					}
				} else {
					t.Errorf("unexpected result type: %T", result)
				}
			}
		})
	}
}

func TestBatchUploadManager_ExecuteBatchUpload(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		respond    func(req *http.Request) (*http.Response, error)
		wantErr    bool
		wantResult string
	}{
		{
			name: "success",
			respond: func(req *http.Request) (*http.Response, error) {
				body := `"executed"`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
			wantErr:    false,
			wantResult: "executed",
		},
		{
			name: "client error",
			respond: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("client error")
			},
			wantErr:    true,
			wantResult: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newMockNuxeoClient(tc.respond)
			bum := &batchUploadManager{client: client, logger: slog.Default()}
			out := new(string)
			op := operation{operationId: "opId"}
			result, err := bum.ExecuteBatchUpload(context.Background(), "batch1", "0", op, out, &nuxeoRequestOptions{})
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tc.wantErr && result != nil {
				if s, ok := result.(*string); ok {
					if *s != tc.wantResult {
						t.Errorf("unexpected result: got %v, want %v", *s, tc.wantResult)
					}
				} else {
					t.Errorf("unexpected result type: %T", result)
				}
			}
		})
	}
}

func TestBatchUploadManager_Upload(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		respond     func(req *http.Request) (*http.Response, error)
		wantErr     bool
		wantBatch   *batchUpload
		wantHeaders map[string]string
	}{
		{
			name: "success",
			respond: func(req *http.Request) (*http.Response, error) {
				// Check required headers
				expected := map[string]string{
					"X-Upload-Type": "normal",
					"X-File-Name":   "file1",
					"X-File-Type":   "application/pdf",
					"X-File-Size":   "123",
				}
				for k, v := range expected {
					if got := req.Header.Get(k); got != v {
						return nil, fmt.Errorf("header %s: got %q, want %q", k, got, v)
					}
				}
				body := `{"name":"file1","batchId":"batch1","fileIdx":"0"}`
				return &http.Response{
					StatusCode: 201,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
			wantErr:   false,
			wantBatch: &batchUpload{Name: "file1", BatchId: "batch1", FileIdx: "0"},
			wantHeaders: map[string]string{
				"X-Upload-Type": "normal",
				"X-File-Name":   "file1",
				"X-File-Type":   "application/pdf",
				"X-File-Size":   "123",
			},
		},
		{
			name: "client error",
			respond: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("client error")
			},
			wantErr:     true,
			wantBatch:   nil,
			wantHeaders: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newMockNuxeoClient(tc.respond)
			bum := &batchUploadManager{client: client, logger: slog.Default()}
			opts := NewUploadOptions("file1", 123, "application/pdf")
			batch, err := bum.Upload(context.Background(), "batch1", "0", opts, strings.NewReader("blobdata"), &nuxeoRequestOptions{})
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tc.wantBatch != nil && batch != nil {
				if batch.Name != tc.wantBatch.Name || batch.BatchId != tc.wantBatch.BatchId || batch.FileIdx != tc.wantBatch.FileIdx {
					t.Errorf("unexpected batch: got %+v, want %+v", batch, tc.wantBatch)
				}
			} else if tc.wantBatch != batch {
				t.Errorf("expected batch %+v, got %+v", tc.wantBatch, batch)
			}
			// Header checks are performed in the mock for "success" case.
		})
	}
}

func TestBatchUploadManager_UploadAsChunk(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		respond   func(req *http.Request) (*http.Response, error)
		wantErr   bool
		wantBatch *batchUpload
	}{
		{
			name: "success",
			respond: func(req *http.Request) (*http.Response, error) {
				// Check required chunked headers
				expected := map[string]string{
					"X-Upload-Type":        "chunked",
					"X-File-Name":          "file1",
					"X-File-Type":          "application/pdf",
					"X-File-Size":          "123",
					"X-Upload-Chunk-Index": "0",
					"X-Upload-Chunk-Count": "1",
				}
				for k, v := range expected {
					if got := req.Header.Get(k); got != v {
						return nil, fmt.Errorf("header %s: got %q, want %q", k, got, v)
					}
				}
				body := `{"name":"file1","batchId":"batch1","fileIdx":"0"}`
				return &http.Response{
					StatusCode: 201,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
			wantErr:   false,
			wantBatch: &batchUpload{Name: "file1", BatchId: "batch1", FileIdx: "0"},
		},
		{
			name: "client error",
			respond: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("client error")
			},
			wantErr:   true,
			wantBatch: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newMockNuxeoClient(tc.respond)
			bum := &batchUploadManager{client: client, logger: slog.Default()}
			opts := uploadOptions{
				fileName:         "file1",
				fileSize:         123,
				fileMimeType:     "application/pdf",
				uploadChunkIndex: 0,
				totalChunkCount:  1,
			}
			batch, err := bum.UploadAsChunk(context.Background(), "batch1", "0", strings.NewReader("chunkdata"), opts, &nuxeoRequestOptions{})
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tc.wantBatch != nil && batch != nil {
				if batch.Name != tc.wantBatch.Name || batch.BatchId != tc.wantBatch.BatchId || batch.FileIdx != tc.wantBatch.FileIdx {
					t.Errorf("unexpected batch: got %+v, want %+v", batch, tc.wantBatch)
				}
			} else if tc.wantBatch != batch {
				t.Errorf("expected batch %+v, got %+v", tc.wantBatch, batch)
			}
			// Header checks are performed in the mock for "success" case.
		})
	}
}
