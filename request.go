package nuxeo

import (
	"strconv"
	"strings"

	"resty.dev/v3"
)

///////////////////////////////
//// NUXEO REQUEST OPTIONS ////
///////////////////////////////

type nuxeoRequestOptions struct {
	RepositoryName      string
	CustomHeaders       map[string]string
	Enrichers           map[string][]string
	FetchProperties     map[string][]string
	TranslateProperties map[string][]string
	Schemas             []string
	Depth               int
	Version             string
	TransactionTimeout  int
	HttpTimeout         int
}

func NewNuxeoRequestOptions() *nuxeoRequestOptions {
	return &nuxeoRequestOptions{
		Enrichers:           make(map[string][]string),
		FetchProperties:     make(map[string][]string),
		TranslateProperties: make(map[string][]string),
	}
}

func (o *nuxeoRequestOptions) SetRepositoryName(name string) *nuxeoRequestOptions {
	o.RepositoryName = name
	return o
}

func (o *nuxeoRequestOptions) SetHeader(key string, value string) *nuxeoRequestOptions {
	o.CustomHeaders[key] = value
	return o
}

func (o *nuxeoRequestOptions) SetTransactionTimeout(timeout int) *nuxeoRequestOptions {
	o.TransactionTimeout = timeout
	return o
}

func (o *nuxeoRequestOptions) SetHttpTimeout(timeout int) *nuxeoRequestOptions {
	o.HttpTimeout = timeout
	return o
}

func (o *nuxeoRequestOptions) SetEnricher(entityType string, values []string) *nuxeoRequestOptions {
	o.Enrichers[entityType] = values
	return o
}

func (o *nuxeoRequestOptions) SetEnricherForDocument(values []string) *nuxeoRequestOptions {
	return o.SetEnricher("document", values)
}

func (o *nuxeoRequestOptions) SetEnricherForUser(values []string) *nuxeoRequestOptions {
	return o.SetEnricher("user", values)
}

func (o *nuxeoRequestOptions) SetFetchProperties(entityType string, values []string) *nuxeoRequestOptions {
	o.FetchProperties[entityType] = values
	return o
}

func (o *nuxeoRequestOptions) SetFetchPropertiesForDirectory(values []string) *nuxeoRequestOptions {
	return o.SetFetchProperties("directory", values)
}

func (o *nuxeoRequestOptions) SetFetchPropertiesForDocument(values []string) *nuxeoRequestOptions {
	return o.SetFetchProperties("document", values)
}

func (o *nuxeoRequestOptions) SetFetchPropertiesForGroup(values []string) *nuxeoRequestOptions {
	return o.SetFetchProperties("group", values)
}

func (o *nuxeoRequestOptions) SetFetchPropertiesForTask(values []string) *nuxeoRequestOptions {
	return o.SetFetchProperties("task", values)
}

func (o *nuxeoRequestOptions) SetFetchPropertiesForWorkflow(values []string) *nuxeoRequestOptions {
	return o.SetFetchProperties("workflow", values)
}

func (o *nuxeoRequestOptions) SetTranslatedProperties(entityType string, values []string) *nuxeoRequestOptions {
	o.TranslateProperties[entityType] = values
	return o
}

func (o *nuxeoRequestOptions) SetTranslatedPropertiesForDirectory(values []string) *nuxeoRequestOptions {
	return o.SetTranslatedProperties("directory", values)
}

func (o *nuxeoRequestOptions) SetSchemas(schemas []string) *nuxeoRequestOptions {
	o.Schemas = schemas
	return o
}

func (o *nuxeoRequestOptions) SetDepth(depth int) *nuxeoRequestOptions {
	o.Depth = depth
	return o
}

func (o *nuxeoRequestOptions) SetVersion(version string) *nuxeoRequestOptions {
	o.Version = version
	return o
}

///////////////////////
//// NUXEO REQUEST ////
///////////////////////

type nuxeoRequest struct {
	*resty.Request
}

func (r *nuxeoRequest) setNuxeoOption(options *nuxeoRequestOptions) *nuxeoRequest {
	if options == nil {
		return r
	}

	// repository name as header
	if options.RepositoryName != "" {
		r.SetHeader(HeaderXRepository, options.RepositoryName)
	}

	// set custom headers
	for key, value := range options.CustomHeaders {
		r.SetHeader(key, value)
	}

	// Set enrichers as headers
	for key, values := range options.Enrichers {
		r.SetHeader("enrichers-"+key, strings.Join(values, ","))
	}

	// Set fetch properties as headers
	for key, values := range options.FetchProperties {
		r.SetHeader("fetch-"+key, strings.Join(values, ","))
	}

	// Set translate properties as headers
	for key, values := range options.TranslateProperties {
		r.SetHeader("translate-"+key, strings.Join(values, ","))
	}

	// Set schemas as header
	if len(options.Schemas) > 0 {
		r.SetHeader(HeaderProperties, strings.Join(options.Schemas, ","))
	}

	// Set depth as header
	if options.Depth > 0 {
		r.SetHeader(HeaderDepth, strconv.Itoa(options.Depth))
	}

	// set version as header
	if options.Version != "" {
		r.SetHeader(HeaderXVersioningOption, options.Version)
	}

	// Set transaction timeout as header
	if options.TransactionTimeout > 0 {
		r.SetHeader(HeaderNuxeoTxTimeout, strconv.Itoa(options.TransactionTimeout))
	}

	// Set HTTP timeout as header
	if options.TransactionTimeout > 0 && options.HttpTimeout == 0 {
		// make the http timeout a bit longer than the transaction timeout
		options.HttpTimeout = options.TransactionTimeout + 5
	}
	if options.HttpTimeout > 0 {
		r.SetHeader(HeaderTimeout, strconv.Itoa(options.HttpTimeout))
	}

	return r
}
