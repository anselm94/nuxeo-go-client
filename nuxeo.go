// Package nuxeo provides a Go client for the Nuxeo API.
// Entry point: NuxeoClient
package nuxeo

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/anselm94/nuxeo/auth"
	"resty.dev/v3"
)

// NuxeoClientOption configures the NuxeoClient.
type NuxeoClientOption func(*NuxeoClient)

// Authenticator defines the interface for authentication strategies.
type Authenticator interface {
	GetAuthHeaders(ctx context.Context, req *http.Request) map[string]string
}

func WithLogger(logger slog.Logger) NuxeoClientOption {
	return func(c *NuxeoClient) {
		c.logger = &logger
	}
}

func WithBeforeRequestMiddleware(middleware resty.RequestMiddleware) NuxeoClientOption {
	return func(c *NuxeoClient) {
		c.middlewareBeforeRequest = &middleware
	}
}

func WithAfterResponseMiddleware(middleware resty.ResponseMiddleware) NuxeoClientOption {
	return func(c *NuxeoClient) {
		c.middlewareAfterResponse = &middleware
	}
}

func WithAuthenticator(authenticator Authenticator) NuxeoClientOption {
	return func(c *NuxeoClient) {
		c.authenticator = authenticator
	}
}

func WithTimeOut(timeout time.Duration) NuxeoClientOption {
	return func(c *NuxeoClient) {
		c.timeout = timeout
	}
}

type NuxeoClient struct {
	logger                  *slog.Logger
	middlewareBeforeRequest *resty.RequestMiddleware
	middlewareAfterResponse *resty.ResponseMiddleware
	authenticator           Authenticator

	// config

	timeout time.Duration

	// internal

	restClient *resty.Client
}

func NewClient(baseUrl string, opts ...NuxeoClientOption) *NuxeoClient {
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

	// setup middlewares
	if client.middlewareBeforeRequest != nil {
		client.restClient.AddRequestMiddleware(*client.middlewareBeforeRequest)
	}
	if client.middlewareAfterResponse != nil {
		client.restClient.AddResponseMiddleware(*client.middlewareAfterResponse)
	}

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

func (c *NuxeoClient) CapabilitiesManager(ctx context.Context) *CapabilitiesManager {
	return &CapabilitiesManager{
		client: c,
		logger: c.logger,
	}
}

func (c *NuxeoClient) Repository() *Repository {
	return &Repository{
		name:   RepositoryDefault,
		client: c,
	}
}

func (c *NuxeoClient) RepositoryWithName(name string) *Repository {
	return &Repository{
		name:   name,
		client: c,
	}
}

func (c *NuxeoClient) OperationManager() *OperationManager {
	return &OperationManager{
		client: c,
		logger: c.logger,
	}
}

func (c *NuxeoClient) UserManager() *UserManager {
	return &UserManager{
		client: c,
		logger: c.logger,
	}
}

func (c *NuxeoClient) DirectoryManager() *DirectoryManager {
	return &DirectoryManager{
		client: c,
		logger: c.logger,
	}
}

func (c *NuxeoClient) TaskManager() *TaskManager {
	return &TaskManager{
		client: c,
		logger: c.logger,
	}
}

func (c *NuxeoClient) BatchUploadManager() *BatchUploadManager {
	return &BatchUploadManager{
		client: c,
		logger: c.logger,
	}
}

func (c *NuxeoClient) ConfigManager() *ConfigManager {
	return &ConfigManager{
		client: c,
		logger: c.logger,
	}
}

////////////////////////
//// COMMON METHODS ////
////////////////////////

func (c *NuxeoClient) CurrentUser(ctx context.Context) (*User, error) {
	return c.UserManager().FetchCurrentUser(ctx)
}

// ServerVersion represents the Nuxeo server version.
type ServerVersion struct {
	Major int
	Minor int
	Patch int
}

func (c *NuxeoClient) ServerVersion(ctx context.Context) (*ServerVersion, error) {
	version := &ServerVersion{}

	// first get the server version from capabilities
	capabilities, err := c.CapabilitiesManager(ctx).FetchCapabilities(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get server version: %w", err)
	}

	// parse version string
	parts := strings.Split(capabilities.Server.DistributionVersion, ".")
	if len(parts) > 0 {
		fmt.Sscanf(parts[0], "%d", &version.Major)
	}
	if len(parts) > 1 {
		fmt.Sscanf(parts[1], "%d", &version.Minor)
	}
	if len(parts) > 2 {
		fmt.Sscanf(parts[2], "%d", &version.Patch)
	}
	return version, nil
}
