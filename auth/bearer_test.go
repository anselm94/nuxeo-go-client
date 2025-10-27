package nuxeoauth

import (
	"context"
	"net/http"
	"testing"
)

func TestNewBearerAuthenticator(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected string
	}{
		{"Normal token", "abc123", "abc123"},
		{"Empty token", "", ""},
		{"Token with spaces", "foo bar", "foo bar"},
		{"Token with special chars", "tok$%^&*()", "tok$%^&*()"},
		{"Token with leading/trailing spaces", "  token  ", "  token  "},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			auth := NewBearerAuthenticator(tt.token)
			if auth.token != tt.expected {
				t.Errorf("expected token %q, got %q", tt.expected, auth.token)
			}
		})
	}
}

func TestBearerAuthenticator_GetAuthHeaders(t *testing.T) {
	tests := []struct {
		name           string
		token          string
		expectedHeader string
		expectHeader   bool
	}{
		{"Normal token", "abc123", "Bearer abc123", true},
		{"Empty token", "", "", false},
		{"Token with spaces", "foo bar", "Bearer foo bar", true},
		{"Token with special chars", "tok$%^&*()", "Bearer tok$%^&*()", true},
		{"Token with leading/trailing spaces", "  token  ", "Bearer   token  ", true},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			auth := NewBearerAuthenticator(tt.token)
			headers := auth.GetAuthHeaders(context.Background(), &http.Request{})
			val, ok := headers["Authorization"]
			if tt.expectHeader {
				if !ok {
					t.Errorf("expected Authorization header, got none")
				} else if val != tt.expectedHeader {
					t.Errorf("expected header value %q, got %q", tt.expectedHeader, val)
				}
			} else {
				if ok {
					t.Errorf("expected no Authorization header, got %q", val)
				}
			}
		})
	}
}
