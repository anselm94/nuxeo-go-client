package nuxeo

import (
	"context"
	"net/http"
)

// Request represents an HTTP request to the Nuxeo API.
type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
}

// NewRequest creates a new Request for the Nuxeo API.
func NewRequest(method, url string, body []byte) *Request {
	return &Request{
		Method:  method,
		URL:     url,
		Headers: make(map[string]string),
		Body:    body,
	}
}

// Do executes the HTTP request using the provided context and http.Client.
func (r *Request) Do(ctx context.Context, client *http.Client, logger Logger, hook Hook) (*http.Response, error) {
	if hook != nil {
		hook.BeforeRequest(r.Method, r.URL)
	}
	logger.Printf("Request: %s %s", r.Method, r.URL)
	req, err := http.NewRequestWithContext(ctx, r.Method, r.URL, nil)
	if err != nil {
		logger.Printf("Error creating request: %v", err)
		return nil, err
	}
	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}
	// TODO: Set request body if present
	resp, err := client.Do(req)
	if hook != nil && resp != nil {
		hook.AfterResponse(r.Method, r.URL, resp.StatusCode)
	}
	if err != nil {
		logger.Printf("Error executing request: %v", err)
	}
	return resp, err
}
