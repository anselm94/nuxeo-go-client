package nuxeo

import (
	"context"
)

// AuthInfo holds authentication details for the NuxeoClient.
type AuthInfo struct {
	User     string
	Password string
	Token    string
}

// Authenticator defines the interface for authentication strategies.
type Authenticator interface {
	Authenticate(ctx context.Context, client *NuxeoClient) error
}

// BasicAuth implements basic username/password authentication.
type BasicAuth struct {
	User     string
	Password string
}

func (a *BasicAuth) Authenticate(ctx context.Context, client *NuxeoClient) error {
	// TODO: Implement basic authentication logic
	return nil
}

// TokenAuth implements token-based authentication.
type TokenAuth struct {
	Token string
}

func (a *TokenAuth) Authenticate(ctx context.Context, client *NuxeoClient) error {
	// TODO: Implement token authentication logic
	return nil
}
