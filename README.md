# Nuxeo Go Client

The Nuxeo Go Client is an _unofficial_ modern, idiomatic Go SDK for the Nuxeo Automation and REST API.

[![Go Reference](https://pkg.go.dev/badge/github.com/anselm94/nuxeo-go-client.svg)](https://pkg.go.dev/github.com/anselm94/nuxeo-go-client)
[![Go Report Card](https://goreportcard.com/badge/github.com/anselm94/nuxeo-go-client)](https://goreportcard.com/report/github.com/anselm94/nuxeo-go-client)
[![Tests](https://github.com/anselm94/nuxeo-go-client/actions/workflows/go.yml/badge.svg)](https://github.com/anselm94/nuxeo-go-client/actions/workflows/go.yml)

## Features

- 🚀 Modern, idiomatic Go API
- 🔒 Multiple authentication strategies (Basic, Token, OAuth2, Portal, Custom)
- 📄 Document CRUD, queries, and property management
- 🤖 Automation operations made easy
- 📦 Blob upload/download (batch & single)
- 👥 User, group, directory, and workflow management
- 🧩 Extensible via interfaces and composition
- 🧪 Table-driven unit tests & examples
- 🧹 Explicit error handling, no panics for normal errors
- 🏃‍♂️ Safe for concurrent use
- 🛠️ Custom logging & instrumentation hooks

## Installation

```bash
go get github.com/anselm94/nuxeo-go-client
```

## Quick Start

Following is a simple example that connects to a Nuxeo server using basic authentication, queries documents, and prints their titles and IDs.

```go
package main

import (
	"context"
	"fmt"

	"github.com/anselm94/nuxeo-go-client"
	nuxeoauth "github.com/anselm94/nuxeo-go-client/auth"
)

func main() {
	ctx := context.Background()

	// client options
	nuxeoClientOptions := nuxeo.DefaultNuxeoClientOptions()

	// basic authenticator
	nuxeoClientOptions.Authenticator = nuxeoauth.NewBasicAuthenticator("Administrator", "Administrator")

	// initialize client
	nuxeoClient := nuxeo.NewClient("<base-url>/nuxeo", &nuxeoClientOptions)

	// get repository
	repo := nuxeoClient.Repository() // or nuxeoClient.RepositoryWithName("default")

	// query documents
	docs, err := repo.Query(ctx, "SELECT * FROM Document", nil, &nuxeo.SortedPaginationOptions{
		CurrentPageIndex: 0,
		PageSize:         5,
		SortBy:           "dc:created",
		SortOrder:        nuxeo.SortOrderDesc,
	}, nil)
	if err != nil {
		panic(err)
	}

	for _, doc := range docs.Entries {
		fmt.Printf(" - %s (%s)\n", doc.Title, doc.ID)
	}
}
```

## Authentication

The client supports multiple authentication strategies out-of-the-box via `github.com/anselm94/nuxeo-go-client/auth` package:

- **Basic**: Username and password
- **Token**: Server-issued authentication token
- **OAuth2 (Auth Code, Client Credentials, JWT)**: OAuth2 grant flows (with automatic refresh)
- **Custom**: Implement your own authenticator

Check out the example - [`examples/nuxeo-auth`](./examples/nuxeo-auth/) for usage patterns.

### 1. Basic Authentication

```go
import (
	"github.com/anselm94/nuxeo-go-client"
	nuxeoauth "github.com/anselm94/nuxeo-go-client/auth"
)

// client options
nuxeoClientOptions := nuxeo.DefaultNuxeoClientOptions()

// basic authenticator
nuxeoClientOptions.Authenticator = nuxeoauth.NewBasicAuthenticator("<username>", "<password>")
```

### 2. Token Authentication

```go
import (
	"github.com/anselm94/nuxeo-go-client"
	nuxeoauth "github.com/anselm94/nuxeo-go-client/auth"
)

// client options
nuxeoClientOptions := nuxeo.DefaultNuxeoClientOptions()

// token authenticator
nuxeoClientOptions.Authenticator = nuxeoauth.NewTokenAuthenticator("<your-token>")
```

### 3. OAuth2 Authentication

#### a. Authorization Code Grant

```go
import (
	"github.com/anselm94/nuxeo-go-client"
	nuxeoauth "github.com/anselm94/nuxeo-go-client/auth"
)

// client options
nuxeoClientOptions := nuxeo.DefaultNuxeoClientOptions()

// authorization code flow authenticator
authCodeFlowOptions := nuxeoauth.NewOAuth2AuthorizationCodeOptions("your-client-id", "your-client-secret", "your-redirect-uri")
authenticator := nuxeoauth.NewOAuth2Authenticator(authCodeFlowOptions, "<base-url>/nuxeo")
nuxeoClientOptions.Authenticator = authenticator

// User visits the authorization URL
authURL := authenticator.AuthCodeUrl(ctx)
fmt.Printf("Visit the URL for the auth dialog: %v", authURL)

// ... handle redirect and get authorization code ...

// After obtaining the auth code from the redirect URL
authCode := "authorization-code-from-callback"
err := authenticator.SetAuthCode(ctx, authCode)
```

#### b. Client Credentials Grant

```go
import (
	"github.com/anselm94/nuxeo-go-client"
	nuxeoauth "github.com/anselm94/nuxeo-go-client/auth"
)

// client options
nuxeoClientOptions := nuxeo.DefaultNuxeoClientOptions()

// client credentials flow authenticator
clientCredentialsOptions := nuxeoauth.NewOAuth2ClientCredentialsOptions("your-client-id", "your-client-secret")
nuxeoClientOptions.Authenticator = nuxeoauth.NewOAuth2Authenticator(clientCredentialsOptions, "<base-url>/nuxeo")
```

#### c. JWT Bearer Grant

```go
import (
	"github.com/anselm94/nuxeo-go-client"
	nuxeoauth "github.com/anselm94/nuxeo-go-client/auth"
)

// client options
nuxeoClientOptions := nuxeo.DefaultNuxeoClientOptions()

jwtOptions := nuxeoauth.NewOAuth2JWTOptions("your-jwt-token")
nuxeoClientOptions.Authenticator = nuxeoauth.NewOAuth2Authenticator(jwtOptions, "<base-url>/nuxeo")
```

### 4. Custom Authentication

Implement the `nuxeo.Authenticator` interface to create a custom authenticator.

```go
import (
	"context"
	"github.com/anselm94/nuxeo-go-client"
)

type CustomAuthenticator struct {
	// Add any necessary fields
}

func (a *CustomAuthenticator) GetAuthHeaders(ctx context.Context, req *http.Request) map[string]string {
	// Implement your custom authentication logic to return headers
	return map[string]string{
		"X-Custom-Auth": "custom-auth-value",
	}
}

// client options
nuxeoClientOptions := nuxeo.DefaultNuxeoClientOptions()

// use the custom authenticator
nuxeoClientOptions.Authenticator = &CustomAuthenticator{}
```

## Document Operations

```go
repo := nuxeoClient.Repository()

// Fetch a document by ID, optionally specifying properties from all applicable schemas
doc, err := repo.FetchDocumentById(ctx, "your-doc-id", nuxeo.NewNuxeoRequestOptions().SetSchemas([]string{"*"}))
if err != nil {
	panic(err)
}
fmt.Println("Title:", doc.Title)

// Query documents using NXQL
docs, err := repo.Query(ctx, "SELECT * FROM Document WHERE ecm:primaryType = 'File'", nil, &nuxeo.SortedPaginationOptions{
	CurrentPageIndex: 0,
	PageSize:         10,
	SortBy:           "dc:created",
	SortOrder:        nuxeo.SortOrderDesc,
}, nil)
if err != nil {
	panic(err)
}

for _, doc := range docs.Entries {
	fmt.Printf("%s (%s)\n", doc.Title, doc.ID)
}
```

## Automation Operations

```go
operationManager := nuxeoClient.OperationManager()

// Run a Nuxeo Automation operation
op := nuxeo.NewOperation("Document.Query")
op.SetParam("query", "SELECT * FROM Document WHERE dc:creator = 'Administrator'")
var result nuxeo.Documents
err := operationManager.ExecuteInto(ctx, *op, &result, nil)
if err != nil {
	panic(err)
}

for _, doc := range result.Entries {
	fmt.Println(doc.Title)
}
```

## Blob Operations

### 1. Uploading Blobs and Creating Documents

```go
repo := nuxeoClient.Repository()
uploadManager := nuxeoClient.BatchUploadManager()

// create a batch
batch, err := uploadManager.CreateBatch(ctx, nil)
if err != nil {
	panic(err)
}

// read the file along with its size
file, err := os.Open("example.pdf")
if err != nil {
	panic(err)
}
defer file.Close()
fileLength, _ := file.Stat()

// upload the file to the batch
uploadOpts := nuxeo.NewUploadOptions("example.pdf", fileLength.Size(), "application/pdf")
batchUploadInfo, err := uploadManager.Upload(ctx, batch.BatchId, "0", uploadOpts, file, nil)
if err != nil {
	panic(err)
}

// create a document with the uploaded blob
newDoc := nuxeo.NewDocument("File", "example.pdf")
newDoc.SetProperty(nuxeo.DocumentPropertyDCDescription, nuxeo.NewStringField("An example PDF file"))
newDoc.SetUploadInfoProperty(nuxeo.DocumentPropertyFileContent, nuxeo.UploadInfo{
	Batch:  batchUploadInfo.BatchId,
	FileId: batchUploadInfo.FileIdx,
})
createdDoc, err := repo.CreateDocumentById(ctx, "<parent id>", *newDoc, nil)
if err != nil {
	panic(err)
}
fmt.Println("Created document:", createdDoc.Title)
```

### 2. Downloading Blobs

```go
repo := nuxeoClient.Repository()

// download the file content
blob, err := repo.StreamBlobById(ctx, "<document id>", nuxeo.DocumentPropertyFileContent, nil)
if err != nil {
	panic(err)
}
defer blob.Close()

// create a local file
outFile, err := os.Create("downloaded_example.pdf")
if err != nil {
	panic(err)
}
defer outFile.Close()

// copy the content from the blob to the local file
_, err = io.Copy(outFile, blob)
if err != nil {
	panic(err)
}
fmt.Println("Downloaded file content to: downloaded_example.pdf")
```

## User, Group Operations

```go
userManager := nuxeoClient.UserManager()

// Fetch a user and a group
user, err := userManager.FetchUser(ctx, "Administrator", nil)
if err != nil {
	panic(err)
}
fmt.Println("User Email:", user.Email())

fetchGroupOptions := nuxeo.NewNuxeoRequestOptions().SetFetchPropertiesForGroup([]string{
	nuxeo.FetchPropertyGroupMemberUsers,
})
group, err := userManager.FetchGroup(ctx, "administrators", fetchGroupOptions)
if err != nil {
	panic(err)
}
fmt.Println("Group Members:", group.MemberUsers)
```

## Directory Operations

```go
directoryManager := nuxeoClient.DirectoryManager()

// List directories and fetch entries
dirs, err := directoryManager.FetchDirectories(ctx, nil)
if err != nil {
	panic(err)
}
fmt.Println("Directories:", dirs)

entries, err := directoryManager.FetchDirectoryEntries(ctx, "continent", nil, nil)
if err != nil {
	panic(err)
}
for _, entry := range entries.Entries {
	fmt.Println("Entry ID:", entry.ID)
}
```

## Workflow Operations

```go
// List tasks
taskManager := nuxeoClient.TaskManager()

tasks, err := taskManager.FetchTasks(ctx, "Administrator", "", "", nil)
if err != nil {
	panic(err)
}
for _, task := range tasks.Entries {
	fmt.Println("Task:", task.Name, task.State)
}
```

## Error Handling & Thread Safety

All errors are returned as the last value. No panics for normal errors. `NuxeoClient` is safe for concurrent use.

```go
// handle error as NuxeoError
if nuxeoErr, ok := err.(*NuxeoError); ok {
	// Handle NuxeoError specifically
	fmt.Println("NuxeoError:", nuxeoErr.Message)
} else {
	// Handle other errors
	fmt.Println("Error:", err)
}
```

## Testing

Run all tests:

```bash
go test ./...
```

## Logging

Inject a custom `slog.Logger` handler to capture structured logs.

```go
import (
    "log/slog"

    "github.com/anselm94/nuxeo-go-client"
)

nuxeoClientOptions := nuxeo.DefaultNuxeoClientOptions()

// set custom logger or any slog adapter - see https://github.com/go-slog/awesome-slog?tab=readme-ov-file#adapters
nuxeoClientOptions.Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
    // ... other options
}))

// initialize client
nuxeoClient := nuxeo.NewClient("<base-url>/nuxeo", &nuxeoClientOptions)
```

## Instrumentation

Register before-request and after-response middlewares for metrics, tracing, etc. The middlewares are [`resty.RequestMiddleware`](https://pkg.go.dev/github.com/go-resty/resty/v2#RequestMiddleware) and [`resty.ResponseMiddleware`](https://pkg.go.dev/github.com/go-resty/resty/v2#ResponseMiddleware) types from [go-resty/resty](https://github.com/go-resty/resty).

```go
import (
    "github.com/anselm94/nuxeo-go-client"
    "github.com/go-resty/resty/v2"
)

nuxeoClientOptions := nuxeo.DefaultNuxeoClientOptions()

// before request middleware
nuxeoClientOptions.BeforeRequestMiddleware = func(c *resty.Client, req *resty.Request) error {
    // custom logic before request is sent
    return nil
}

// after response middleware
nuxeoClientOptions.AfterResponseMiddleware = func(c *resty.Client, resp *resty.Response) error {
    // custom logic after response is received
    return nil
}

// initialize client
nuxeoClient := nuxeo.NewClient("<base-url>/nuxeo", &nuxeoClientOptions)
```

## Advanced Usage

- Custom authenticators via interface
- Context for cancellation and timeouts
- Batch uploads, automation, workflows
- Request/response hooks for metrics/tracing

## Project Structure

```
├── auth/                # Authentication strategies
├── examples/            # Example programs
├── internal/            # Internal helpers
├── entity-*.go          # Domain entities
├── manager-*.go         # Managers for repository, batch upload, etc.
├── operation.go         # Automation operations
├── blob.go              # Blob/file upload/download
├── nuxeo.go             # Main client implementation
├── errors.go            # Error types and handling
├── constants.go         # API constants
├── *_test.go            # Table-driven unit tests
└── README.md            # This file
```

## Contributing

PRs and issues welcome!

## License

[MIT License](./LICENSE)
