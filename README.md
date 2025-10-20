# Nuxeo Go Client

A modern, idiomatic Go SDK for the Nuxeo Platform REST API.  
Designed for composability, testability, and developer experience.

## Features

- Synchronous, context-aware API
- Explicit error handling
- Pluggable authentication (Basic, Token, OAuth2, etc.)
- Document CRUD, queries, and property management
- Automation operations
- Blob upload/download (batch and single)
- User, group, directory, workflow management
- Extensible via interfaces and composition
- Table-driven tests and examples

## Installation

```bash
go get github.com/anselm94/nuxeo
```

## Quick Start

```go
import (
    "context"
    "fmt"
    "github.com/anselm94/nuxeo"
)

// Example custom logger
 type myLogger struct{}
 func (l myLogger) Printf(format string, v ...any) {
     fmt.Printf("LOG: "+format+"\n", v...)
 }

func main() {
    ctx := context.Background()
    client, err := nuxeo.NewClient(ctx,
        nuxeo.WithBaseURL("https://nuxeo.example.com"),
        nuxeo.WithUser("Administrator"),
        nuxeo.WithPassword("secret"),
        nuxeo.WithLogger(myLogger{}), // Custom logger
    )
    if err != nil {
        panic(err)
    }

    repo := nuxeo.NewRepository("default")
    docs, err := repo.QueryDocuments("SELECT * FROM Document WHERE ecm:parentId IS NULL")
    if err != nil {
        panic(err)
    }
    for _, doc := range docs {
        fmt.Printf(" - %s (%s)\n", doc.Properties["dc:title"], doc.ID)
    }
}
```

## Authentication

The Nuxeo Go client supports multiple authentication strategies, matching those available on the Nuxeo Platform. You specify the method and credentials when constructing the client.

### Supported Methods

- Basic: Username and password
- Portal: Username and portal secret
- Token: Server-issued authentication token
- Bearer/OAuth2: OAuth2 access token or full token object (with automatic refresh)
- Custom: Register your own authenticator

### Basic Authentication

```go
client := nuxeo.NewClient(ctx,
    nuxeo.WithBaseURL("https://nuxeo.example.com"),
    nuxeo.WithAuth(nuxeo.AuthInfo{
        Method:   "basic",
        User:     "Administrator",
        Password: "Administrator",
    }),
)
```

### Portal Authentication

```go
client := nuxeo.NewClient(ctx,
    nuxeo.WithBaseURL("https://nuxeo.example.com"),
    nuxeo.WithAuth(nuxeo.AuthInfo{
        Method: "portal",
        User:   "joe",
        Secret: "shared-secret-from-server",
    }),
)
```

### Token Authentication

```go
client := nuxeo.NewClient(ctx,
    nuxeo.WithBaseURL("https://nuxeo.example.com"),
    nuxeo.WithAuth(nuxeo.AuthInfo{
        Method: "token",
        Token:  "a-token",
    }),
)
```

To request a token from the server using your credentials:

```go
token, err := nuxeo.RequestAuthenticationToken(ctx, "https://nuxeo.example.com", nuxeo.AuthInfo{
    Method:   "basic",
    User:     "Administrator",
    Password: "Administrator",
}, "My App", "deviceUID", "deviceName", "rw")
if err != nil {
    panic(err)
}
client := nuxeo.NewClient(ctx,
    nuxeo.WithBaseURL("https://nuxeo.example.com"),
    nuxeo.WithAuth(nuxeo.AuthInfo{
        Method: "token",
        Token:  token,
    }),
)
```

### OAuth2 / Bearer Token Authentication

```go
client := nuxeo.NewClient(ctx,
    nuxeo.WithBaseURL("https://nuxeo.example.com"),
    nuxeo.WithAuth(nuxeo.AuthInfo{
        Method: "bearerToken", // or "bearer"
        Bearer: "access_token", // or use OAuth2 field for full token object
        ClientID: "my-app",     // optional, for automatic refresh
        ClientSecret: "my-secret", // required if the client defines a secret
        OAuth2: &nuxeo.OAuth2Token{
            AccessToken:  "access_token",
            RefreshToken: "refresh_token",
            TokenType:    "bearer",
            ExpiresIn:    3600,
        },
    }),
)
```

#### OAuth2 & JWT Grant Helpers

Generate an authorization URL (using oauth2.Config):

```go
import (
    "golang.org/x/oauth2"
    "github.com/anselm94/nuxeo"
)

config := &oauth2.Config{
    ClientID:     "my-app",
    ClientSecret: "my-secret",
    Endpoint: oauth2.Endpoint{
        AuthURL:  "https://nuxeo.example.com/oauth2/authorize",
        TokenURL: "https://nuxeo.example.com/oauth2/token",
    },
    Scopes: []string{"openid"},
}
url := nuxeo.GetAuthorizationURL(config, "xyz")
fmt.Println(url)
```

Exchange an authorization code for an access token:

```go
token, err := nuxeo.FetchAccessTokenFromAuthorizationCode(ctx, config, "AUTH_CODE")
if err != nil {
    panic(err)
}
```

JWT Bearer Grant:

```go
jwt := "<your-jwt-token>"
token, err := nuxeo.FetchAccessTokenFromJWTToken(ctx, config, jwt)
if err != nil {
    panic(err)
}
client := nuxeo.NewClient(ctx,
    nuxeo.WithBaseURL("https://nuxeo.example.com"),
    nuxeo.WithAuth(nuxeo.AuthInfo{
        Method:      "oauth2",
        OAuth2Token: token,
        JWT:         jwt, // optional, for refresh
    }),
)
```

Refresh an access token:

```go
token, err := nuxeo.RefreshAccessToken(ctx, config, token)
if err != nil {
    panic(err)
}
```

### Migration Note

If you previously used manual token structs or custom expiry logic, migrate to using `*oauth2.Token` and the provided helpers. All expiry and refresh logic is now handled automatically and idiomatically via the Go standard library and SDK helpers.

### Custom Authenticators

You can register your own authentication strategy by implementing the `Authenticator` interface:

```go
type MyAuth struct{}
func (a *MyAuth) ComputeAuthenticationHeaders(auth nuxeo.AuthInfo) map[string]string { ... }
func (a *MyAuth) AuthenticateURL(url string, auth nuxeo.AuthInfo) string { ... }
func (a *MyAuth) CanRefreshAuthentication() bool { ... }
func (a *MyAuth) RefreshAuthentication(ctx context.Context, baseURL string, auth nuxeo.AuthInfo) (nuxeo.AuthInfo, error) { ... }

nuxeo.RegisterAuthenticator("myMethod", &MyAuth{})
```

## Document Operations

```go
doc, err := repo.Fetch(ctx, "some-id")
doc.Set(map[string]any{"dc:title": "New Title"})
err = doc.Save(ctx)
```

## Blob Upload

```go
file, _ := os.Open("example.pdf")
info, _ := file.Stat()
blob := nuxeo.NewBlob(file, "example.pdf", "application/pdf", info.Size())

batch := client.BatchUpload()
blobs, err := batch.Upload(ctx, blob)
```

## Automation Operations

```go
op := client.Operation("Blob.Convert")
op.SetInput(doc).SetParam("converter", "text")
result, err := op.Execute(ctx)
```

## User, Group, Directory, Workflow

```go
user, err := client.Users().Fetch(ctx, "jsmith")
group, err := client.Groups().Fetch(ctx, "administrators")
dir := client.Directory("countries")
entries, err := dir.FetchAll(ctx)
wf, err := client.Workflows("default").Start(ctx, "SerialDocumentReview", nil)
```

## Error Handling & Thread Safety

All errors are returned as the last value.  
No panics for normal errors.

`NuxeoClient` is safe for concurrent use if no mutable state is shared between goroutines. Configuration should be set at construction time.

## Testing

Table-driven tests and mocks provided.  
Run all tests:

```bash
go test ./...
```

## Logging & Instrumentation

- Inject a custom logger via `WithLogger` to capture SDK events.
- Register hooks for request/response instrumentation (metrics, tracing).

## Advanced Usage

- Custom authenticators via interface
- Context for cancellation and timeouts
- Batch uploads, automation, workflows
- Request/response hooks for metrics/tracing

## Contributing

See [docs/design](docs/design/) for architecture and DX guidelines.  
PRs and issues welcome!

## License

MIT
