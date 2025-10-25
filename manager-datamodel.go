package nuxeo

import (
	"context"
	"log/slog"

	"github.com/anselm94/nuxeo/internal"
)

type DataModelManager struct {
	// internal

	client *NuxeoClient
	logger *slog.Logger
}

func (dmm *DataModelManager) FetchTypes(ctx context.Context) (*DocTypes, error) {
	path := internal.PathApiV1 + "/config/types"
	res, err := dmm.client.NewRequest(ctx, nil).SetResult(&DocTypes{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		dmm.logger.Error("Failed to fetch document types", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*DocTypes), nil
}

func (dmm *DataModelManager) FetchType(ctx context.Context, name string) (*DocType, error) {
	path := internal.PathApiV1 + "/config/types/" + name
	res, err := dmm.client.NewRequest(ctx, nil).SetResult(&DocType{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		dmm.logger.Error("Failed to fetch document type", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*DocType), nil
}

func (dmm *DataModelManager) FetchSchemas(ctx context.Context) (*Schemas, error) {
	path := internal.PathApiV1 + "/config/schemas"
	res, err := dmm.client.NewRequest(ctx, nil).SetResult(&Schemas{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		dmm.logger.Error("Failed to fetch schemas", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Schemas), nil
}

func (dmm *DataModelManager) FetchSchema(ctx context.Context, name string) (*Schema, error) {
	path := internal.PathApiV1 + "/config/schemas/" + name
	res, err := dmm.client.NewRequest(ctx, nil).SetResult(&Schema{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		dmm.logger.Error("Failed to fetch schema", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Schema), nil
}

func (dmm *DataModelManager) FetchFacets(ctx context.Context) (*Facets, error) {
	path := internal.PathApiV1 + "/config/facets"
	res, err := dmm.client.NewRequest(ctx, nil).SetResult(&Facets{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		dmm.logger.Error("Failed to fetch facets", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Facets), nil
}

func (dmm *DataModelManager) FetchFacet(ctx context.Context, name string) (*Facet, error) {
	path := internal.PathApiV1 + "/config/facets/" + name
	res, err := dmm.client.NewRequest(ctx, nil).SetResult(&Facet{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		dmm.logger.Error("Failed to fetch facet", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Facet), nil
}
