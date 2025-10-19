package nuxeo

import (
	"context"
	"net/http"
	"testing"
)

func TestNewRequest_Builder(t *testing.T) {
	req := NewRequest("GET", "http://localhost:8080/nuxeo/api", nil)
	if req.Method != "GET" {
		t.Errorf("Method got %q, want %q", req.Method, "GET")
	}
	if req.URL != "http://localhost:8080/nuxeo/api" {
		t.Errorf("URL got %q, want %q", req.URL, "http://localhost:8080/nuxeo/api")
	}
	if len(req.Headers) != 0 {
		t.Errorf("Headers got %v, want empty map", req.Headers)
	}
}

func TestRequest_Do_MockClient(t *testing.T) {
	ctx := context.Background()
	req := NewRequest("GET", "http://localhost:8080/nuxeo/api", nil)
	client := &http.Client{}
	// This will fail without a real server, but should not panic
	hook := testHook{}
	logger := testLogger{}
	_, err := req.Do(ctx, client, logger, hook)
	if err == nil {
		t.Errorf("Do got nil error, want error due to no server")
	}
}

type testLogger struct{}

func (l testLogger) Printf(format string, v ...any) {}

type testHook struct{}

func (h testHook) BeforeRequest(method, url string)             {}
func (h testHook) AfterResponse(method, url string, status int) {}
