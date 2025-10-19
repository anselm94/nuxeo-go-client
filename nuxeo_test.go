package nuxeo

import (
	"context"
	"testing"
)

func TestNewClient_OptionChaining(t *testing.T) {
	ctx := context.Background()
	client, err := NewClient(ctx,
		WithBaseURL("http://localhost:8080/nuxeo"),
		WithUser("admin"),
		WithPassword("secret"),
		WithToken("token123"),
		WithTimeout(30),
	)
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	if client.options.BaseURL != "http://localhost:8080/nuxeo" {
		t.Errorf("BaseURL got %q, want %q", client.options.BaseURL, "http://localhost:8080/nuxeo")
	}
	if client.options.User != "admin" {
		t.Errorf("User got %q, want %q", client.options.User, "admin")
	}
	if client.options.Password != "secret" {
		t.Errorf("Password got %q, want %q", client.options.Password, "secret")
	}
	if client.options.Token != "token123" {
		t.Errorf("Token got %q, want %q", client.options.Token, "token123")
	}
	if client.options.Timeout != 30 {
		t.Errorf("Timeout got %d, want %d", client.options.Timeout, 30)
	}
}
