package nuxeo

import (
	"context"
	"log/slog"

	"github.com/anselm94/nuxeo-go-client/internal"
)

// dataModelManager provides methods to introspect Nuxeo Server data model configuration.
// It allows fetching document types, schemas, and facets via the REST API endpoints:
// See: https://doc.nuxeo.com/rest-api/1/data-model-endpoint/
type dataModelManager struct {
	client *NuxeoClient // Nuxeo API client
	logger *slog.Logger // Logger for error reporting
}

// FetchTypes retrieves all document types contributed to the Nuxeo Server.
// Endpoint: GET /config/types
// Returns a collection of document type entities.
// See: https://doc.nuxeo.com/rest-api/1/data-model-endpoint/#get-all-document-types
func (dmm *dataModelManager) FetchTypes(ctx context.Context) (*DocTypes, error) {
	path := internal.PathApiV1 + "/config/types"
	res, err := dmm.client.NewRequest(ctx, nil).SetResult(&DocTypes{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		dmm.logger.Error("Failed to fetch document types", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*DocTypes), nil
}

// FetchType retrieves a single document type entity by name.
// Endpoint: GET /config/types/{DOC_TYPE}
// Returns the document type entity for the given DOC_TYPE.
// See: https://doc.nuxeo.com/rest-api/1/data-model-endpoint/#get-a-document-type
func (dmm *dataModelManager) FetchType(ctx context.Context, name string) (*DocType, error) {
	path := internal.PathApiV1 + "/config/types/" + name
	res, err := dmm.client.NewRequest(ctx, nil).SetResult(&DocType{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		dmm.logger.Error("Failed to fetch document type", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*DocType), nil
}

// FetchSchemas retrieves all schemas contributed to the Nuxeo Server.
// Endpoint: GET /config/schemas
// Returns a collection of schema entities.
// See: https://doc.nuxeo.com/rest-api/1/data-model-endpoint/#get-all-schemas
func (dmm *dataModelManager) FetchSchemas(ctx context.Context) (*Schemas, error) {
	path := internal.PathApiV1 + "/config/schemas"
	res, err := dmm.client.NewRequest(ctx, nil).SetResult(&Schemas{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		dmm.logger.Error("Failed to fetch schemas", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Schemas), nil
}

// FetchSchema retrieves a single schema entity by name.
// Endpoint: GET /config/schemas/{SCHEMA_NAME}
// Returns the schema entity for the given SCHEMA_NAME.
// See: https://doc.nuxeo.com/rest-api/1/data-model-endpoint/#get-a-schema
func (dmm *dataModelManager) FetchSchema(ctx context.Context, name string) (*Schema, error) {
	path := internal.PathApiV1 + "/config/schemas/" + name
	res, err := dmm.client.NewRequest(ctx, nil).SetResult(&Schema{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		dmm.logger.Error("Failed to fetch schema", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Schema), nil
}

// FetchFacets retrieves all facets contributed to the Nuxeo Server.
// Endpoint: GET /config/facets
// Returns a collection of facet entities.
// See: https://doc.nuxeo.com/rest-api/1/data-model-endpoint/#get-all-facets
func (dmm *dataModelManager) FetchFacets(ctx context.Context) (*Facets, error) {
	path := internal.PathApiV1 + "/config/facets"
	res, err := dmm.client.NewRequest(ctx, nil).SetResult(&Facets{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		dmm.logger.Error("Failed to fetch facets", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Facets), nil
}

// FetchFacet retrieves a single facet entity by name.
// Endpoint: GET /config/facets/{FACET_NAME}
// Returns the facet entity for the given FACET_NAME.
// See: https://doc.nuxeo.com/rest-api/1/data-model-endpoint/#get-a-facet
func (dmm *dataModelManager) FetchFacet(ctx context.Context, name string) (*Facet, error) {
	path := internal.PathApiV1 + "/config/facets/" + name
	res, err := dmm.client.NewRequest(ctx, nil).SetResult(&Facet{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		dmm.logger.Error("Failed to fetch facet", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Facet), nil
}
