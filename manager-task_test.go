package nuxeo

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"testing"
)

// helper to create a taskManager with a mock client
func newTestTaskManager(respond func(req *http.Request) (*http.Response, error)) *taskManager {
	return &taskManager{
		client: newMockNuxeoClient(respond),
		logger: slog.Default(),
	}
}

func TestTaskManager_FetchTasks(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name    string
		resp    any
		status  int
		wantErr bool
	}{
		{
			name:    "success",
			resp:    Tasks{Entries: []Task{{entity: entity{EntityType: "task"}}}},
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
			tm := newTestTaskManager(func(req *http.Request) (*http.Response, error) {
				body := testMarshalBody(t, tc.resp)
				return &http.Response{
					StatusCode: tc.status,
					Body:       body,
					Header:     make(http.Header),
				}, nil
			})
			got, err := tm.FetchTasks(context.Background(), "user1", "wf1", "model1", nil)
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

func TestTaskManager_FetchTask(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name    string
		resp    any
		status  int
		wantErr bool
	}{
		{
			name:    "success",
			resp:    Task{entity: entity{EntityType: "task"}},
			status:  http.StatusOK,
			wantErr: false,
		},
		{
			name:    "not found",
			resp:    NuxeoError{Message: "not found"},
			status:  http.StatusNotFound,
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tm := newTestTaskManager(func(req *http.Request) (*http.Response, error) {
				body := testMarshalBody(t, tc.resp)
				return &http.Response{
					StatusCode: tc.status,
					Body:       body,
					Header:     make(http.Header),
				}, nil
			})
			got, err := tm.FetchTask(context.Background(), "task1", nil)
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

func TestTaskManager_ReassignTask(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name    string
		resp    any
		status  int
		wantErr bool
	}{
		{
			name:    "success",
			resp:    Task{entity: entity{EntityType: "task"}},
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
			tm := newTestTaskManager(func(req *http.Request) (*http.Response, error) {
				body := testMarshalBody(t, tc.resp)
				return &http.Response{
					StatusCode: tc.status,
					Body:       body,
					Header:     make(http.Header),
				}, nil
			})
			got, err := tm.ReassignTask(context.Background(), "task1", "actor1", "comment", nil)
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

func TestTaskManager_DelegateTask(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name    string
		resp    any
		status  int
		wantErr bool
	}{
		{
			name:    "success",
			resp:    Task{entity: entity{EntityType: "task"}},
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
			tm := newTestTaskManager(func(req *http.Request) (*http.Response, error) {
				body := testMarshalBody(t, tc.resp)
				return &http.Response{
					StatusCode: tc.status,
					Body:       body,
					Header:     make(http.Header),
				}, nil
			})
			got, err := tm.DelegateTask(context.Background(), "task1", "actor2", "comment", nil)
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

func TestTaskManager_CompleteTask(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name    string
		resp    any
		status  int
		wantErr bool
	}{
		{
			name:    "success",
			resp:    Task{entity: entity{EntityType: "task"}},
			status:  http.StatusOK,
			wantErr: false,
		},
		{
			name:    "error response",
			resp:    NuxeoError{Message: "fail"},
			status:  http.StatusBadRequest,
			wantErr: true,
		},
		{
			name:    "network error",
			resp:    nil,
			status:  http.StatusOK,
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tm := newTestTaskManager(func(req *http.Request) (*http.Response, error) {
				if tc.name == "network error" {
					return nil, errors.New("network fail")
				}
				body := testMarshalBody(t, tc.resp)
				return &http.Response{
					StatusCode: tc.status,
					Body:       body,
					Header:     make(http.Header),
				}, nil
			})
			req := TaskCompletionRequest{
				Id:        "task1",
				Comment:   "done",
				Variables: map[string]Field{"foo": Field(json.RawMessage(`"bar"`))},
			}
			got, err := tm.CompleteTask(context.Background(), "task1", "approve", req, nil)
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
