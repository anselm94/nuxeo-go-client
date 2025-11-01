package nuxeoauth

import (
	"context"
	"net/http"
	"testing"
)

func TestNewNoOpAuthenticator(t *testing.T) {
	auth := NewNoOpAuthenticator()
	if auth == nil {
		t.Fatal("NewNoOpAuthenticator() returned nil")
	}
}

func TestNoOpAuthenticator_GetAuthHeaders(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name string
		ctx  context.Context
		req  *http.Request
	}{
		{"nil context and request", nil, nil},
		{"background context, nil request", context.Background(), nil},
		{"background context, empty request", context.Background(), &http.Request{}},
	}

	auth := NewNoOpAuthenticator()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			headers := auth.GetAuthHeaders(nil)
			if headers == nil {
				t.Errorf("Expected empty map, got nil")
			}
			if len(headers) != 0 {
				t.Errorf("Expected empty map, got: %v", headers)
			}
		})
	}
}
