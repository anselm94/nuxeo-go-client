package nuxeoauth

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

import (
	"context"
	"net/http"
)

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

func (a *BearerAuthenticator) GetAuthHeaders(ctx context.Context, req *http.Request) map[string]string {
	headers := make(map[string]string)
	if a.token != "" {
		headers["Authorization"] = "Bearer " + a.token
	}
	return headers
}
