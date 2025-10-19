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

Supports multiple strategies:
- Basic
- Token
- Bearer
- Portal
- OAuth2

Example:

```go
client := nuxeo.NewClient("https://nuxeo.example.com", nuxeo.AuthInfo{
    Method: "oauth2",
    OAuth2Token: nuxeo.OAuth2Token{
        AccessToken: "your-access-token",
    },
})
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
