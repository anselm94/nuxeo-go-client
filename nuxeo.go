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

	"github.com/anselm94/nuxeo-go-client/auth"
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

func WithHeader(key string, value string) NuxeoClientOption {
	return func(c *NuxeoClient) {
		c.headers[key] = value
	}
}

type NuxeoClient struct {
	logger                  *slog.Logger
	middlewareBeforeRequest *resty.RequestMiddleware
	middlewareAfterResponse *resty.ResponseMiddleware
	authenticator           Authenticator

	// config

	timeout time.Duration
	headers map[string]string

	// internal

	restClient *resty.Client
}

func NewClient(baseUrl string, opts ...NuxeoClientOption) *NuxeoClient {
	client := &NuxeoClient{
		headers: make(map[string]string),
	}

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

func (c *NuxeoClient) NewRequest(ctx context.Context, options *nuxeoRequestOptions) *nuxeoRequest {
	req := &nuxeoRequest{
		Request: c.restClient.R().SetContext(ctx),
	}

	// set headers from client
	for k, v := range c.headers {
		req.SetHeader(k, v)
	}

	return req.setNuxeoOption(options)
}

func (c *NuxeoClient) CapabilitiesManager(ctx context.Context) *capabilitiesManager {
	return &capabilitiesManager{
		client: c,
		logger: c.logger,
	}
}

func (c *NuxeoClient) Repository() *repository {
	return &repository{
		name:   RepositoryDefault,
		client: c,
		logger: c.logger,
	}
}

func (c *NuxeoClient) RepositoryWithName(name string) *repository {
	return &repository{
		name:   name,
		client: c,
	}
}

func (c *NuxeoClient) OperationManager() *operationManager {
	return &operationManager{
		client: c,
		logger: c.logger,
	}
}

func (c *NuxeoClient) UserManager() *userManager {
	return &userManager{
		client: c,
		logger: c.logger,
	}
}

func (c *NuxeoClient) DirectoryManager() *directoryManager {
	return &directoryManager{
		client: c,
		logger: c.logger,
	}
}

func (c *NuxeoClient) TaskManager() *taskManager {
	return &taskManager{
		client: c,
		logger: c.logger,
	}
}

func (c *NuxeoClient) BatchUploadManager() *batchUploadManager {
	return &batchUploadManager{
		client: c,
		logger: c.logger,
	}
}

func (c *NuxeoClient) DataModelManager() *dataModelManager {
	return &dataModelManager{
		client: c,
		logger: c.logger,
	}
}

/////////////////
//// METHODS ////
/////////////////

func (c *NuxeoClient) CurrentUser(ctx context.Context) (*entityUser, error) {
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
