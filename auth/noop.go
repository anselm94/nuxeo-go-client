package auth

import (
	"context"
	"net/http"
)

// NoOpAuthenticator implements no-op placeholder authentication.
type NoOpAuthenticator struct {
}

// NewNoOpAuthenticator creates a new NoOpAuthenticator.
func NewNoOpAuthenticator() *NoOpAuthenticator {
	return &NoOpAuthenticator{}
}

func (a *NoOpAuthenticator) GetAuthHeaders(ctx context.Context, req *http.Request) map[string]string {
	// No operation performed
	return map[string]string{}
}
