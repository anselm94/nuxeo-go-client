package nuxeoauth

import "resty.dev/v3"

/**
BearerAuthenticator implements Bearer token authentication.

Example:
```go
import (
	"github.com/anselm94/nuxeo-go-client/auth"
)

authenticator := nuxeoauth.NewBearerAuthenticator("your-bearer-token")
```
*/

// BearerAuthenticator implements Bearer token authentication.
type BearerAuthenticator struct {
	token string
}

// NewBearerAuthenticator creates a new BearerAuthenticator with the given token.
func NewBearerAuthenticator(token string) *BearerAuthenticator {
	return &BearerAuthenticator{
		token: token,
	}
}

func (a *BearerAuthenticator) GetAuthHeaders(req *resty.Request) map[string]string {
	headers := make(map[string]string)
	if a.token != "" {
		headers["Authorization"] = "Bearer " + a.token
	}
	return headers
}
