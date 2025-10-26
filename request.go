package nuxeo

import (
	"strconv"
	"strings"

	"github.com/anselm94/nuxeo-go-client/internal"
	"resty.dev/v3"
)

/////////////////////////////////
//// NUXEO REQUEST OPTIONS ////
/////////////////////////////////

// nuxeoRequestOptions configures options for a Nuxeo API request, such as repository, headers, enrichers, schemas, and timeouts.
type nuxeoRequestOptions struct {
	repositoryName      string
	customHeaders       map[string]string
	enrichers           map[string][]string
	fetchProperties     map[string][]string
	translateProperties map[string][]string
	schemas             []string
	depth               int
	version             string
	transactionTimeout  int
	httpTimeout         int
}

// NewNuxeoRequestOptions creates a new nuxeoRequestOptions with initialized maps.
func NewNuxeoRequestOptions() *nuxeoRequestOptions {
	return &nuxeoRequestOptions{
		enrichers:           make(map[string][]string),
		fetchProperties:     make(map[string][]string),
		translateProperties: make(map[string][]string),
	}
}

// SetRepositoryName sets the repository name for the request.
func (o *nuxeoRequestOptions) SetRepositoryName(name string) *nuxeoRequestOptions {
	o.repositoryName = name
	return o
}

// SetHeader adds a custom header to the request.
func (o *nuxeoRequestOptions) SetHeader(key string, value string) *nuxeoRequestOptions {
	o.customHeaders[key] = value
	return o
}

// SetTransactionTimeout sets the transaction timeout (in seconds).
func (o *nuxeoRequestOptions) SetTransactionTimeout(timeout int) *nuxeoRequestOptions {
	o.transactionTimeout = timeout
	return o
}

// SetHttpTimeout sets the HTTP timeout (in seconds).
func (o *nuxeoRequestOptions) SetHttpTimeout(timeout int) *nuxeoRequestOptions {
	o.httpTimeout = timeout
	return o
}

// SetEnricher sets enrichers for a given entity type.
func (o *nuxeoRequestOptions) SetEnricher(entityType string, values []string) *nuxeoRequestOptions {
	o.enrichers[entityType] = values
	return o
}

// SetEnricherForDocument sets enrichers for document entities.
func (o *nuxeoRequestOptions) SetEnricherForDocument(values []string) *nuxeoRequestOptions {
	return o.SetEnricher("document", values)
}

// SetEnricherForUser sets enrichers for user entities.
func (o *nuxeoRequestOptions) SetEnricherForUser(values []string) *nuxeoRequestOptions {
	return o.SetEnricher("user", values)
}

// SetFetchProperties sets fetch properties for a given entity type.
func (o *nuxeoRequestOptions) SetFetchProperties(entityType string, values []string) *nuxeoRequestOptions {
	o.fetchProperties[entityType] = values
	return o
}

// SetFetchPropertiesForDirectory sets fetch properties for directory entities.
func (o *nuxeoRequestOptions) SetFetchPropertiesForDirectory(values []string) *nuxeoRequestOptions {
	return o.SetFetchProperties("directory", values)
}

// SetFetchPropertiesForDocument sets fetch properties for document entities.
func (o *nuxeoRequestOptions) SetFetchPropertiesForDocument(values []string) *nuxeoRequestOptions {
	return o.SetFetchProperties("document", values)
}

// SetFetchPropertiesForGroup sets fetch properties for group entities.
func (o *nuxeoRequestOptions) SetFetchPropertiesForGroup(values []string) *nuxeoRequestOptions {
	return o.SetFetchProperties("group", values)
}

// SetFetchPropertiesForTask sets fetch properties for task entities.
func (o *nuxeoRequestOptions) SetFetchPropertiesForTask(values []string) *nuxeoRequestOptions {
	return o.SetFetchProperties("task", values)
}

// SetFetchPropertiesForWorkflow sets fetch properties for workflow entities.
func (o *nuxeoRequestOptions) SetFetchPropertiesForWorkflow(values []string) *nuxeoRequestOptions {
	return o.SetFetchProperties("workflow", values)
}

// SetTranslatedProperties sets translated properties for a given entity type.
func (o *nuxeoRequestOptions) SetTranslatedProperties(entityType string, values []string) *nuxeoRequestOptions {
	o.translateProperties[entityType] = values
	return o
}

// SetTranslatedPropertiesForDirectory sets translated properties for directory entities.
func (o *nuxeoRequestOptions) SetTranslatedPropertiesForDirectory(values []string) *nuxeoRequestOptions {
	return o.SetTranslatedProperties("directory", values)
}

// SetSchemas sets the schemas to fetch for the request.
func (o *nuxeoRequestOptions) SetSchemas(schemas []string) *nuxeoRequestOptions {
	o.schemas = schemas
	return o
}

// SetDepth sets the depth for the request (e.g., for document children).
func (o *nuxeoRequestOptions) SetDepth(depth int) *nuxeoRequestOptions {
	o.depth = depth
	return o
}

// SetVersion sets the versioning option for the request.
func (o *nuxeoRequestOptions) SetVersion(version string) *nuxeoRequestOptions {
	o.version = version
	return o
}

///////////////////////
//// NUXEO REQUEST ////
///////////////////////

// nuxeoRequest wraps a Resty request for Nuxeo API calls.
type nuxeoRequest struct {
	*resty.Request
}

// setNuxeoOption applies nuxeoRequestOptions to the request as headers and parameters.
func (r *nuxeoRequest) setNuxeoOption(options *nuxeoRequestOptions) *nuxeoRequest {
	if options == nil {
		return r
	}

	// repository name as header
	if options.repositoryName != "" {
		r.SetHeader(internal.HeaderXRepository, options.repositoryName)
	}

	// set custom headers
	for key, value := range options.customHeaders {
		r.SetHeader(key, value)
	}

	// Set enrichers as headers
	for key, values := range options.enrichers {
		r.SetHeader("enrichers-"+key, strings.Join(values, ","))
	}

	// Set fetch properties as headers
	for key, values := range options.fetchProperties {
		r.SetHeader("fetch-"+key, strings.Join(values, ","))
	}

	// Set translate properties as headers
	for key, values := range options.translateProperties {
		r.SetHeader("translate-"+key, strings.Join(values, ","))
	}

	// Set schemas as header
	if len(options.schemas) > 0 {
		r.SetHeader(internal.HeaderProperties, strings.Join(options.schemas, ","))
	}

	// Set depth as header
	if options.depth > 0 {
		r.SetHeader(internal.HeaderDepth, strconv.Itoa(options.depth))
	}

	// set version as header
	if options.version != "" {
		r.SetHeader(internal.HeaderXVersioningOption, options.version)
	}

	// Set transaction timeout as header
	if options.transactionTimeout > 0 {
		r.SetHeader(internal.HeaderNuxeoTxTimeout, strconv.Itoa(options.transactionTimeout))
	}

	// Set HTTP timeout as header
	if options.transactionTimeout > 0 && options.httpTimeout == 0 {
		// make the http timeout a bit longer than the transaction timeout
		options.httpTimeout = options.transactionTimeout + 5
	}
	if options.httpTimeout > 0 {
		r.SetHeader(internal.HeaderTimeout, strconv.Itoa(options.httpTimeout))
	}

	return r
}
