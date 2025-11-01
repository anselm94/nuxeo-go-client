package nuxeoauth

import "resty.dev/v3"

/**
TokenAuthenticator implements token-based authentication

Example:
```go
import (
	"github.com/anselm94/nuxeo-go-client/auth"
)

authenticator := nuxeoauth.NewTokenAuthenticator("your-token")
```
*/

// TokenAuthenticator implements token-based authentication.
type TokenAuthenticator struct {
	token string
}

func NewTokenAuthenticator(token string) *TokenAuthenticator {
	return &TokenAuthenticator{
		token: token,
	}
}

func (a *TokenAuthenticator) GetAuthHeaders(req *resty.Request) map[string]string {
	headers := make(map[string]string)
	if a.token != "" {
		headers["X-Authentication-Token"] = a.token
	}
	return headers
}
