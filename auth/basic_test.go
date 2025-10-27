package nuxeoauth

import (
	"context"
	"encoding/base64"
	"net/http"
	"reflect"
	"testing"
)

func TestNewBasicAuthenticator(t *testing.T) {
	a := NewBasicAuthenticator("user", "pass")
	if a.username != "user" {
		t.Errorf("Expected username 'user', got '%s'", a.username)
	}
	if a.password != "pass" {
		t.Errorf("Expected password 'pass', got '%s'", a.password)
	}
}

func TestGetAuthHeaders_ValidCredentials(t *testing.T) {
	a := NewBasicAuthenticator("user", "pass")
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	headers := a.GetAuthHeaders(context.Background(), req)
	cred := base64.StdEncoding.EncodeToString([]byte("user:pass"))
	expected := map[string]string{"Authorization": "Basic " + cred}
	if !reflect.DeepEqual(headers, expected) {
		t.Errorf("Expected headers %v, got %v", expected, headers)
	}
}

func TestGetAuthHeaders_EmptyUsername(t *testing.T) {
	a := NewBasicAuthenticator("", "pass")
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	headers := a.GetAuthHeaders(context.Background(), req)
	if len(headers) != 0 {
		t.Errorf("Expected empty headers, got %v", headers)
	}
}

func TestGetAuthHeaders_EmptyPassword(t *testing.T) {
	a := NewBasicAuthenticator("user", "")
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	headers := a.GetAuthHeaders(context.Background(), req)
	if len(headers) != 0 {
		t.Errorf("Expected empty headers, got %v", headers)
	}
}

func TestGetAuthHeaders_Base64Encoding(t *testing.T) {
	a := NewBasicAuthenticator("foo", "bar")
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	headers := a.GetAuthHeaders(context.Background(), req)
	expectedCred := base64.StdEncoding.EncodeToString([]byte("foo:bar"))
	expected := "Basic " + expectedCred
	if val, ok := headers["Authorization"]; !ok || val != expected {
		t.Errorf("Expected Authorization header '%s', got '%s'", expected, val)
	}
}
