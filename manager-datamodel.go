package nuxeo

import (
	"context"
	"log/slog"

	"github.com/anselm94/nuxeo/internal"
)

type dataModelManager struct {
	// internal

	client *NuxeoClient
	logger *slog.Logger
}

func (dmm *dataModelManager) FetchTypes(ctx context.Context) (*entityDocTypes, error) {
	path := internal.PathApiV1 + "/config/types"
	res, err := dmm.client.NewRequest(ctx, nil).SetResult(&entityDocTypes{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		dmm.logger.Error("Failed to fetch document types", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityDocTypes), nil
}

func (dmm *dataModelManager) FetchType(ctx context.Context, name string) (*entityDocType, error) {
	path := internal.PathApiV1 + "/config/types/" + name
	res, err := dmm.client.NewRequest(ctx, nil).SetResult(&entityDocType{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		dmm.logger.Error("Failed to fetch document type", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityDocType), nil
}

func (dmm *dataModelManager) FetchSchemas(ctx context.Context) (*entitySchemas, error) {
	path := internal.PathApiV1 + "/config/schemas"
	res, err := dmm.client.NewRequest(ctx, nil).SetResult(&entitySchemas{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		dmm.logger.Error("Failed to fetch schemas", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entitySchemas), nil
}

func (dmm *dataModelManager) FetchSchema(ctx context.Context, name string) (*entitySchema, error) {
	path := internal.PathApiV1 + "/config/schemas/" + name
	res, err := dmm.client.NewRequest(ctx, nil).SetResult(&entitySchema{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		dmm.logger.Error("Failed to fetch schema", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entitySchema), nil
}

func (dmm *dataModelManager) FetchFacets(ctx context.Context) (*entityFacets, error) {
	path := internal.PathApiV1 + "/config/facets"
	res, err := dmm.client.NewRequest(ctx, nil).SetResult(&entityFacets{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		dmm.logger.Error("Failed to fetch facets", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityFacets), nil
}

func (dmm *dataModelManager) FetchFacet(ctx context.Context, name string) (*entityFacet, error) {
	path := internal.PathApiV1 + "/config/facets/" + name
	res, err := dmm.client.NewRequest(ctx, nil).SetResult(&entityFacet{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		dmm.logger.Error("Failed to fetch facet", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityFacet), nil
}
