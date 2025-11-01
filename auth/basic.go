package nuxeoauth

/**
BasicAuthenticator implements basic username/password authentication.

Example:
```go
import (
	"github.com/anselm94/nuxeo-go-client/auth"
)

authenticator := nuxeoauth.NewBasicAuthenticator("Administrator", "password")
```
*/

import (
	"encoding/base64"

	"resty.dev/v3"
)

// BasicAuthenticator implements basic username/password authentication.
type BasicAuthenticator struct {
	username string
	password string
}

// NewBasicAuthenticator creates a new BasicAuthenticator with the given username and password.
func NewBasicAuthenticator(username string, password string) *BasicAuthenticator {
	return &BasicAuthenticator{
		username: username,
		password: password,
	}
}

func (a *BasicAuthenticator) GetAuthHeaders(req *resty.Request) map[string]string {
	headers := make(map[string]string)
	if a.username != "" && a.password != "" {
		cred := a.username + ":" + a.password
		headers["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(cred))
	}
	return headers
}
