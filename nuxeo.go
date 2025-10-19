// Package nuxeo provides a Go client for the Nuxeo API.
// Entry point: NuxeoClient
package nuxeo

import (
	"context"
)

// Option configures the NuxeoClient.
type Option func(*NuxeoClient)

// NuxeoClient is the main entry point for interacting with the Nuxeo API.
type Logger interface {
	Printf(format string, v ...any)
}

// Hook allows instrumentation of requests and responses for metrics/tracing.
type Hook interface {
	BeforeRequest(method, url string)
	AfterResponse(method, url string, status int)
}

type NuxeoClient struct {
	options BaseOptions
	logger  Logger
	hook    Hook
	// TODO: Add config, auth, http client, etc.
}

// NewClient creates a new NuxeoClient with the given options.
// WithLogger sets a custom logger for the client.
func WithLogger(logger Logger) Option {
	return func(c *NuxeoClient) {
		c.logger = logger
	}
}

// WithHook sets a custom hook for metrics/tracing.
func WithHook(hook Hook) Option {
	return func(c *NuxeoClient) {
		c.hook = hook
	}
}

func NewClient(ctx context.Context, opts ...Option) (*NuxeoClient, error) {
	client := &NuxeoClient{}
	for _, opt := range opts {
		opt(client)
	}
	if client.logger == nil {
		client.logger = defaultLogger{}
	}
	// TODO: Implement client construction and option handling
	return client, nil
}

// defaultLogger is a basic logger using the standard library.
type defaultLogger struct{}

func (l defaultLogger) Printf(format string, v ...any) {
	// No-op or use log.Printf(format, v...) if desired
}
