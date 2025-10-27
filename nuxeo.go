// Package nuxeo provides a Go client for the Nuxeo API.
// Entry point: NuxeoClient
package nuxeo

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"strings"
	"time"

	nuxeoauth "github.com/anselm94/nuxeo-go-client/auth"
	"resty.dev/v3"
)

///////////////////////
//// Authenticator ////
///////////////////////

// Authenticator defines the interface for authentication strategies.
type Authenticator interface {
	// GetAuthHeaders returns authentication headers for the request context.
	GetAuthHeaders(ctx context.Context, req *http.Request) map[string]string
}

//////////////////////////////
//// Nuxeo Client Options ////
//////////////////////////////

type nuxeoClientOptions struct {
	Authenticator           Authenticator
	Logger                  *slog.Logger
	BeforeRequestMiddleware resty.RequestMiddleware
	AfterResponseMiddleware resty.ResponseMiddleware
	Timeout                 time.Duration
	CustomHeaders           map[string]string
}

func DefaultNuxeoClientOptions() nuxeoClientOptions {
	return nuxeoClientOptions{
		Authenticator:           nuxeoauth.NewNoOpAuthenticator(),
		Logger:                  slog.Default(),
		BeforeRequestMiddleware: nil,
		AfterResponseMiddleware: nil,
		Timeout:                 30 * time.Second,
		CustomHeaders:           make(map[string]string),
	}
}

//////////////////////
//// Nuxeo Client ////
//////////////////////

// NuxeoClient is the main client for interacting with the Nuxeo API.
type NuxeoClient struct {
	authenticator           Authenticator
	logger                  *slog.Logger
	middlewareBeforeRequest *resty.RequestMiddleware
	middlewareAfterResponse *resty.ResponseMiddleware

	// config

	timeout time.Duration
	headers map[string]string

	// internal

	restClient *resty.Client
}

// NewClient creates a new NuxeoClient for the given base URL and options.
func NewClient(baseUrl string, options *nuxeoClientOptions) *NuxeoClient {
	client := &NuxeoClient{
		headers: make(map[string]string),
	}

	// default options if nil
	if options == nil {
		defaultOptions := DefaultNuxeoClientOptions()
		options = &defaultOptions
	}

	// apply options
	client.authenticator = options.Authenticator
	client.logger = options.Logger
	if options.BeforeRequestMiddleware != nil {
		client.middlewareBeforeRequest = &options.BeforeRequestMiddleware
	}
	if options.AfterResponseMiddleware != nil {
		client.middlewareAfterResponse = &options.AfterResponseMiddleware
	}
	client.timeout = options.Timeout
	maps.Copy(client.headers, options.CustomHeaders)

	// setup resty client
	client.restClient = resty.New()
	client.restClient.SetBaseURL(baseUrl)

	// setup authenticator
	if client.authenticator == nil {
		client.authenticator = nuxeoauth.NewNoOpAuthenticator()
	}
	client.restClient.AddRequestMiddleware(func(c *resty.Client, r *resty.Request) error {
		headers := client.authenticator.GetAuthHeaders(r.Context(), r.RawRequest)
		for k, v := range headers {
			r.SetHeader(k, v)
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

	// setup connection config
	if client.timeout == 0 {
		client.timeout = 30 * time.Second // default timeout
	}
	client.restClient.SetTimeout(client.timeout)

	return client
}

// Close releases resources held by the NuxeoClient.
func (c *NuxeoClient) Close() error {
	return c.restClient.Close()
}

// NewRequest creates a new Nuxeo API request with the given context and options.
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

// CapabilitiesManager returns a manager for Nuxeo server capabilities.
func (c *NuxeoClient) CapabilitiesManager(ctx context.Context) *capabilitiesManager {
	return &capabilitiesManager{
		client: c,
		logger: c.logger,
	}
}

// Repository returns a manager for the default Nuxeo repository.
func (c *NuxeoClient) Repository() *repository {
	return &repository{
		name:   RepositoryDefault,
		client: c,
		logger: c.logger,
	}
}

// RepositoryWithName returns a manager for the specified Nuxeo repository.
func (c *NuxeoClient) RepositoryWithName(name string) *repository {
	return &repository{
		name:   name,
		client: c,
	}
}

// OperationManager returns a manager for Nuxeo automation operations.
func (c *NuxeoClient) OperationManager() *operationManager {
	return &operationManager{
		client: c,
		logger: c.logger,
	}
}

// UserManager returns a manager for Nuxeo users.
func (c *NuxeoClient) UserManager() *userManager {
	return &userManager{
		client: c,
		logger: c.logger,
	}
}

// DirectoryManager returns a manager for Nuxeo directories.
func (c *NuxeoClient) DirectoryManager() *directoryManager {
	return &directoryManager{
		client: c,
		logger: c.logger,
	}
}

// TaskManager returns a manager for Nuxeo tasks.
func (c *NuxeoClient) TaskManager() *taskManager {
	return &taskManager{
		client: c,
		logger: c.logger,
	}
}

// BatchUploadManager returns a manager for Nuxeo batch uploads.
func (c *NuxeoClient) BatchUploadManager() *batchUploadManager {
	return &batchUploadManager{
		client: c,
		logger: c.logger,
	}
}

// DataModelManager returns a manager for Nuxeo data models.
func (c *NuxeoClient) DataModelManager() *dataModelManager {
	return &dataModelManager{
		client: c,
		logger: c.logger,
	}
}

/////////////////
//// METHODS ////
/////////////////

// CurrentUser returns the current authenticated Nuxeo user.
func (c *NuxeoClient) CurrentUser(ctx context.Context) (*entityUser, error) {
	return c.UserManager().FetchCurrentUser(ctx)
}

// ServerVersion represents the Nuxeo server version.
type ServerVersion struct {
	Major int // Major version number
	Minor int // Minor version number
	Patch int // Patch version number
}

// ServerVersion returns the Nuxeo server version as a struct.
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
