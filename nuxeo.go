// Package nuxeo provides a Go client for the Nuxeo API.
// Entry point: NuxeoClient
package nuxeo

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/anselm94/nuxeo/auth"
	"github.com/anselm94/nuxeo/hook"
	"resty.dev/v3"
)

// Option configures the NuxeoClient.
type Option func(*NuxeoClient)

// Authenticator defines the interface for authentication strategies.
type Authenticator interface {
	GetAuthHeaders(ctx context.Context, req *http.Request) map[string]string
}

// Hook allows instrumentation of requests and responses for metrics/tracing.
type Hook interface {
	BeforeRequest(method, url string)
	AfterResponse(method, url string, status int)
}

func WithLogger(logger slog.Logger) Option {
	return func(c *NuxeoClient) {
		c.logger = &logger
	}
}

func WithHook(hook Hook) Option {
	return func(c *NuxeoClient) {
		c.hook = hook
	}
}

func WithAuthenticator(authenticator Authenticator) Option {
	return func(c *NuxeoClient) {
		c.authenticator = authenticator
	}
}

func WithTimeOut(timeout time.Duration) Option {
	return func(c *NuxeoClient) {
		c.timeout = timeout
	}
}

type NuxeoClient struct {
	logger        *slog.Logger
	authenticator Authenticator
	hook          Hook
	timeout       time.Duration
	restClient    *resty.Client
}

func NewClient(baseUrl string, opts ...Option) *NuxeoClient {
	client := &NuxeoClient{}

	for _, opt := range opts {
		opt(client)
	}

	// setup resty client
	client.restClient = resty.New()
	client.restClient.SetBaseURL(baseUrl)

	// setup authenticator
	if client.authenticator == nil {
		client.authenticator = auth.NewNoOpAuthenticator()
	}
	client.restClient.AddRequestMiddleware(func(c *resty.Client, r *resty.Request) error {
		headers := client.authenticator.GetAuthHeaders(r.Context(), r.RawRequest)
		for k, v := range headers {
			client.restClient.SetHeader(k, v)
		}
		return nil
	})

	// setup hooks
	if client.hook == nil {
		client.hook = hook.NewNoOpHook()
	}
	client.restClient.AddRequestMiddleware(func(c *resty.Client, r *resty.Request) error {
		client.hook.BeforeRequest(r.Method, r.URL)
		return nil
	})
	client.restClient.AddResponseMiddleware(func(c *resty.Client, r *resty.Response) error {
		client.hook.AfterResponse(r.Request.Method, r.Request.URL, r.StatusCode())
		return nil
	})

	// setup logger
	if client.logger == nil {
		client.logger = slog.Default()
	}

	// setup connection config
	if client.timeout == 0 {
		client.timeout = 30 * time.Second // default timeout
	}
	client.restClient.SetTimeout(client.timeout)

	return client
}

func (c *NuxeoClient) Close() error {
	return c.restClient.Close()
}
