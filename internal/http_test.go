package internal

import (
	"net/http"
	"testing"
)

func TestNewHTTPClient(t *testing.T) {
	client := NewHTTPClient()
	if client == nil {
		t.Errorf("NewHTTPClient returned nil")
	}
	_, ok := any(client).(*http.Client)
	if !ok {
		t.Errorf("NewHTTPClient did not return *http.Client")
	}
}
