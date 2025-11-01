package nuxeoauth

import (
	"testing"
)

func TestNewTokenAuthenticator(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name  string
		token string
	}{
		{"empty token", ""},
		{"normal token", "abc123"},
		{"whitespace token", "   "},
		{"special chars", "!@#$%^&*()_+-="},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			auth := NewTokenAuthenticator(tc.token)
			if auth == nil {
				t.Fatal("NewTokenAuthenticator returned nil")
			}
			if auth.token != tc.token {
				t.Errorf("Expected token '%s', got '%s'", tc.token, auth.token)
			}
		})
	}
}

func TestTokenAuthenticator_GetAuthHeaders(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		token    string
		expected map[string]string
	}{
		{"empty token", "", map[string]string{}},
		{"normal token", "abc123", map[string]string{"X-Authentication-Token": "abc123"}},
		{"whitespace token", "   ", map[string]string{"X-Authentication-Token": "   "}},
		{"special chars", "!@#$%^&*()_+-=", map[string]string{"X-Authentication-Token": "!@#$%^&*()_+-="}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			auth := NewTokenAuthenticator(tc.token)
			headers := auth.GetAuthHeaders(nil)
			if len(headers) != len(tc.expected) {
				t.Errorf("Expected %d headers, got %d", len(tc.expected), len(headers))
			}
			for k, v := range tc.expected {
				if headers[k] != v {
					t.Errorf("Expected header %s=%s, got %s", k, v, headers[k])
				}
			}
			for k := range headers {
				if _, ok := tc.expected[k]; !ok {
					t.Errorf("Unexpected header: %s", k)
				}
			}
		})
	}
}
