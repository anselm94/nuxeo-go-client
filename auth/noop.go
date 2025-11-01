package nuxeoauth

import (
	"resty.dev/v3"
)

// NoOpAuthenticator implements no-op placeholder authentication.
type NoOpAuthenticator struct {
}

// NewNoOpAuthenticator creates a new NoOpAuthenticator.
func NewNoOpAuthenticator() *NoOpAuthenticator {
	return &NoOpAuthenticator{}
}

func (a *NoOpAuthenticator) GetAuthHeaders(req *resty.Request) map[string]string {
	// No operation performed
	return map[string]string{}
}
