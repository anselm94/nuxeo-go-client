package nuxeoauth

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
import (
	"context"
	"net/http"
)

// TokenAuthenticator implements token-based authentication.
type TokenAuthenticator struct {
	token string
}

func NewTokenAuthenticator(token string) *TokenAuthenticator {
	return &TokenAuthenticator{
		token: token,
	}
}

func (a *TokenAuthenticator) GetAuthHeaders(ctx context.Context, req *http.Request) map[string]string {
	headers := make(map[string]string)
	if a.token != "" {
		headers["X-Authentication-Token"] = a.token
	}
	return headers
}
