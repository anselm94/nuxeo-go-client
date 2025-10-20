package nuxeo

import (
	"context"
	"strings"

	"resty.dev/v3"
)

type NuxeoRequestOption struct {
	RepositoryName      string
	Schemas             []string
	Enrichers           map[string][]string
	FetchProperties     map[string][]string
	TranslateProperties map[string][]string
	Depth               int
	HttpTimeout         int
	TransactionTimeout  int
}

type NuxeoRequest struct {
	*resty.Request
}

func (c *NuxeoClient) NewRequest(ctx context.Context) *NuxeoRequest {
	req := c.restClient.R().SetContext(ctx)

	return &NuxeoRequest{
		Request: req,
	}
}

func (r *NuxeoRequest) SetNuxeoOption(options *NuxeoRequestOption) *NuxeoRequest {
	if options == nil {
		return r
	}

	// repository name as header
	if options.RepositoryName != "" {
		r.SetHeader("X-NXRepository", options.RepositoryName)
	}

	// Set schemas as header
	if len(options.Schemas) > 0 {
		r.SetHeader("properties", strings.Join(options.Schemas, ","))
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

	// Set depth as header
	if options.Depth > 0 {
		r.SetHeader("depth", string(options.Depth))
	}

	// Set transaction timeout as header
	if options.TransactionTimeout > 0 {
		r.SetHeader("Nuxeo-Transaction-Timeout", string(options.TransactionTimeout))
	}

	// Set HTTP timeout as header
	if options.TransactionTimeout > 0 && options.HttpTimeout == 0 {
		// make the http timeout a bit longer than the transaction timeout
		options.HttpTimeout = options.TransactionTimeout + 5
	}
	if options.HttpTimeout > 0 {
		r.SetHeader("timeout", string(options.HttpTimeout))
	}

	return r
}
