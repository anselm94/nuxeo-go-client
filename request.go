package nuxeo

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

/**
Nuxeo Request Builder and Executor to build and execute GET and POST requests, including path, query parameters, headers, and body.

Example:

```go
ctx := context.Background()
client, err := nuxeo.NewClient(ctx,
    nuxeo.WithBaseURL("https://nuxeo.example.com/nuxeo"),
    nuxeo.WithUser("Administrator"),
    nuxeo.WithPassword("password"),
)
if err != nil {
    panic(err)
}
// Build a GET request to fetch documents of type "File"
resp, err := nuxeo.NewRequest(client).
    Path("api").
    Path("v1").
    Path("doc").
    QueryParams(map[string]string{"type": "File", "pageSize": "10"}).
    Header("Accept", "application/json").
    Get(ctx)
if err != nil {
    fmt.Printf("Request failed: %v\n", err)
    return
}
defer resp.Body.Close()
// Handle resp.Body (JSON decode, etc.)
// Build a POST request with a JSON body
jsonBody := []byte(`{"entity-type":"document","type":"File","name":"myfile"}`)
resp, err = nuxeo.NewRequest(client).
    Path("api").
    Path("v1").
    Path("doc").
    Header("Content-Type", "application/json").
    Body(jsonBody).
    Post(ctx)
if err != nil {
    fmt.Printf("POST failed: %v\n", err)
    return
}
defer resp.Body.Close()
// Handle resp.Body

// You can also use Put/Delete similarly:
// nuxeo.NewRequest(client).Path("api").Path("v1").Path("doc").Put(ctx)
// nuxeo.NewRequest(client).Path("api").Path("v1").Path("doc").Delete(ctx)
```
*/

// Request builds and executes REST requests on a Nuxeo Platform instance.
// Chainable for path and query param configuration. Use HTTP methods to execute.
type Request struct {
	client    *NuxeoClient
	pathParts []string
	query     url.Values
	method    string
	headers   map[string]string
	body      []byte
}

// NewRequest creates a new Request builder for the Nuxeo API.
func NewRequest(c *NuxeoClient) *Request {
	return &Request{
		client:    c,
		pathParts: []string{},
		query:     url.Values{},
		headers:   map[string]string{},
	}
}

// Path appends a path segment to the request path. Chainable.
func (r *Request) Path(segment string) *Request {
	if segment != "" {
		r.pathParts = append(r.pathParts, segment)
	}
	return r
}

// QueryParams merges new query parameters. Chainable.
func (r *Request) QueryParams(params map[string]string) *Request {
	for k, v := range params {
		r.query.Set(k, v)
	}
	return r
}

// Header sets a header. Chainable.
func (r *Request) Header(key, value string) *Request {
	if key != "" {
		r.headers[key] = value
	}
	return r
}

// Body sets the request body. Chainable.
func (r *Request) Body(b []byte) *Request {
	r.body = b
	return r
}

// Get executes a GET request.
func (r *Request) Get(ctx context.Context) (*http.Response, error) {
	r.method = http.MethodGet
	return r.Execute(ctx)
}

// Post executes a POST request.
func (r *Request) Post(ctx context.Context) (*http.Response, error) {
	r.method = http.MethodPost
	return r.Execute(ctx)
}

// Put executes a PUT request.
func (r *Request) Put(ctx context.Context) (*http.Response, error) {
	r.method = http.MethodPut
	return r.Execute(ctx)
}

// Delete executes a DELETE request.
func (r *Request) Delete(ctx context.Context) (*http.Response, error) {
	r.method = http.MethodDelete
	return r.Execute(ctx)
}

// Execute builds the final URL and options, then executes the request via Nuxeo's HTTP method.
func (r *Request) Execute(ctx context.Context) (*http.Response, error) {
	base := r.client.options.BaseURL
	path := strings.Join(r.pathParts, "/")
	fullURL := fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(path, "/"))
	if len(r.query) > 0 {
		fullURL += "?" + r.query.Encode()
	}
	if r.client.logger != nil {
		r.client.logger.Printf("Request: %s %s", r.method, fullURL)
	}
	if r.client.hook != nil {
		r.client.hook.BeforeRequest(r.method, fullURL)
	}
	var bodyReader *bytes.Reader
	if r.body != nil {
		bodyReader = bytes.NewReader(r.body)
	}
	req, err := http.NewRequestWithContext(ctx, r.method, fullURL, bodyReader)
	if err != nil {
		if r.client.logger != nil {
			r.client.logger.Printf("Error creating request: %v", err)
		}
		return nil, err
	}
	for k, v := range r.headers {
		req.Header.Set(k, v)
	}
	resp, err := r.client.httpClient.Do(req)
	if r.client.hook != nil && resp != nil {
		r.client.hook.AfterResponse(r.method, fullURL, resp.StatusCode)
	}
	if err != nil && r.client.logger != nil {
		r.client.logger.Printf("Error executing request: %v", err)
	}
	return resp, err
}

// DX: Example usage
//   resp, err := nuxeo.NewRequest().Path("api").Path("v1").Path("doc").QueryParams(map[string]string{"type": "File"}).Get(ctx)
//   // resp, err := nuxeo.NewRequest().Path("directory").Path("user").Get(ctx)
