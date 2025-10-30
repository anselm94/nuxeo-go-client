package nuxeo

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"log/slog"
)

// helper to create a directoryManager with a mock client
func newTestDirectoryManager(respond func(req *http.Request) (*http.Response, error)) *directoryManager {
	return &directoryManager{
		client: newMockNuxeoClient(respond),
		logger: slog.Default(),
	}
}

func TestDirectoryManager_FetchDirectories(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name    string
		resp    any
		status  int
		wantErr bool
	}{
		{
			name:    "success",
			resp:    Directories{Entries: []Directory{{Name: "foo"}}},
			status:  http.StatusOK,
			wantErr: false,
		},
		{
			name:    "error response",
			resp:    NuxeoError{Message: "fail"},
			status:  http.StatusBadRequest,
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dm := newTestDirectoryManager(func(req *http.Request) (*http.Response, error) {
				body := testMarshalBody(t, tc.resp)
				return &http.Response{
					StatusCode: tc.status,
					Body:       body,
					Header:     make(http.Header),
				}, nil
			})
			got, err := dm.FetchDirectories(context.Background(), nil)
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tc.wantErr && got == nil {
				t.Errorf("expected result, got nil")
			}
		})
	}
}

func TestDirectoryManager_FetchDirectoryEntries(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name    string
		resp    any
		status  int
		wantErr bool
	}{
		{
			name:    "success",
			resp:    DirectoryEntries{Entries: []DirectoryEntry{{ID: "id1"}}},
			status:  http.StatusOK,
			wantErr: false,
		},
		{
			name:    "error response",
			resp:    NuxeoError{Message: "fail"},
			status:  http.StatusBadRequest,
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dm := newTestDirectoryManager(func(req *http.Request) (*http.Response, error) {
				body := testMarshalBody(t, tc.resp)
				return &http.Response{
					StatusCode: tc.status,
					Body:       body,
					Header:     make(http.Header),
				}, nil
			})
			got, err := dm.FetchDirectoryEntries(context.Background(), "foo", nil, nil)
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tc.wantErr && got == nil {
				t.Errorf("expected result, got nil")
			}
		})
	}
}

func TestDirectoryManager_CreateDirectoryEntry(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name    string
		resp    any
		status  int
		wantErr bool
	}{
		{
			name:    "success",
			resp:    DirectoryEntry{ID: "id1"},
			status:  http.StatusCreated,
			wantErr: false,
		},
		{
			name:    "error response",
			resp:    NuxeoError{Message: "fail"},
			status:  http.StatusBadRequest,
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dm := newTestDirectoryManager(func(req *http.Request) (*http.Response, error) {
				body := testMarshalBody(t, tc.resp)
				return &http.Response{
					StatusCode: tc.status,
					Body:       body,
					Header:     make(http.Header),
				}, nil
			})
			entry := DirectoryEntry{ID: "id1"}
			got, err := dm.CreateDirectoryEntry(context.Background(), "foo", entry, nil)
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tc.wantErr && got == nil {
				t.Errorf("expected result, got nil")
			}
		})
	}
}

func TestDirectoryManager_FetchDirectoryEntry(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name    string
		resp    any
		status  int
		wantErr bool
	}{
		{
			name:    "success",
			resp:    DirectoryEntry{ID: "id1"},
			status:  http.StatusOK,
			wantErr: false,
		},
		{
			name:    "error response",
			resp:    NuxeoError{Message: "fail"},
			status:  http.StatusNotFound,
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dm := newTestDirectoryManager(func(req *http.Request) (*http.Response, error) {
				body := testMarshalBody(t, tc.resp)
				return &http.Response{
					StatusCode: tc.status,
					Body:       body,
					Header:     make(http.Header),
				}, nil
			})
			got, err := dm.FetchDirectoryEntry(context.Background(), "foo", "id1", nil)
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tc.wantErr && got == nil {
				t.Errorf("expected result, got nil")
			}
		})
	}
}

func TestDirectoryManager_UpdateDirectoryEntry(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name    string
		resp    any
		status  int
		wantErr bool
	}{
		{
			name:    "success",
			resp:    DirectoryEntry{ID: "id1"},
			status:  http.StatusOK,
			wantErr: false,
		},
		{
			name:    "error response",
			resp:    NuxeoError{Message: "fail"},
			status:  http.StatusBadRequest,
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dm := newTestDirectoryManager(func(req *http.Request) (*http.Response, error) {
				body := testMarshalBody(t, tc.resp)
				return &http.Response{
					StatusCode: tc.status,
					Body:       body,
					Header:     make(http.Header),
				}, nil
			})
			entry := DirectoryEntry{ID: "id1"}
			got, err := dm.UpdateDirectoryEntry(context.Background(), "foo", "id1", entry, nil)
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tc.wantErr && got == nil {
				t.Errorf("expected result, got nil")
			}
		})
	}
}

func TestDirectoryManager_DeleteDirectoryEntry(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name    string
		status  int
		wantErr bool
		respErr error
	}{
		{
			name:    "success",
			status:  http.StatusNoContent,
			wantErr: false,
		},
		{
			name:    "error response",
			status:  http.StatusBadRequest,
			wantErr: true,
		},
		{
			name:    "network error",
			status:  http.StatusOK,
			wantErr: true,
			respErr: errors.New("network fail"),
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dm := newTestDirectoryManager(func(req *http.Request) (*http.Response, error) {
				if tc.respErr != nil {
					return nil, tc.respErr
				}
				body := testMarshalBody(t, NuxeoError{Message: "fail"})
				return &http.Response{
					StatusCode: tc.status,
					Body:       body,
					Header:     make(http.Header),
				}, nil
			})
			err := dm.DeleteDirectoryEntry(context.Background(), "foo", "id1", nil)
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
