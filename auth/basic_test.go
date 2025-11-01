package nuxeoauth

import (
	"encoding/base64"
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
	headers := a.GetAuthHeaders(nil)
	cred := base64.StdEncoding.EncodeToString([]byte("user:pass"))
	expected := map[string]string{"Authorization": "Basic " + cred}
	if !reflect.DeepEqual(headers, expected) {
		t.Errorf("Expected headers %v, got %v", expected, headers)
	}
}

func TestGetAuthHeaders_EmptyUsername(t *testing.T) {
	a := NewBasicAuthenticator("", "pass")
	headers := a.GetAuthHeaders(nil)
	if len(headers) != 0 {
		t.Errorf("Expected empty headers, got %v", headers)
	}
}

func TestGetAuthHeaders_EmptyPassword(t *testing.T) {
	a := NewBasicAuthenticator("user", "")
	headers := a.GetAuthHeaders(nil)
	if len(headers) != 0 {
		t.Errorf("Expected empty headers, got %v", headers)
	}
}

func TestGetAuthHeaders_Base64Encoding(t *testing.T) {
	a := NewBasicAuthenticator("foo", "bar")
	headers := a.GetAuthHeaders(nil)
	expectedCred := base64.StdEncoding.EncodeToString([]byte("foo:bar"))
	expected := "Basic " + expectedCred
	if val, ok := headers["Authorization"]; !ok || val != expected {
		t.Errorf("Expected Authorization header '%s', got '%s'", expected, val)
	}
}
