package nuxeo

import (
	"context"
	"testing"
)

func TestBasicAuth_Authenticate(t *testing.T) {
	ctx := context.Background()
	client := &NuxeoClient{}
	auth := &BasicAuth{User: "admin", Password: "secret"}
	err := auth.Authenticate(ctx, client)
	if err != nil {
		t.Errorf("BasicAuth.Authenticate returned error: %v", err)
	}
}

func TestTokenAuth_Authenticate(t *testing.T) {
	ctx := context.Background()
	client := &NuxeoClient{}
	auth := &TokenAuth{Token: "token123"}
	err := auth.Authenticate(ctx, client)
	if err != nil {
		t.Errorf("TokenAuth.Authenticate returned error: %v", err)
	}
}
