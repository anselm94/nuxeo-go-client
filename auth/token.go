package auth

/**
TokenAuthenticator implements token-based authentication

Example:
```go
import (
	"context"
	"github.com/anselm94/nuxeo-go-client"
	"github.com/anselm94/nuxeo-go-client/auth"
)

ctx := context.Background()
authenticator := NewTokenAuthenticator("your-token")
client, err := nuxeo.NewClient(ctx,
	nuxeo.WithBaseURL("https://nuxeo.example.com/nuxeo"),
	nuxeo.WithAuthenticator(authenticator),
)
if err != nil {
	panic(err)
}
// Use client...
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
