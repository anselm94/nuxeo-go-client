package nuxeo

import (
	"context"
	"log/slog"
)

type ConfigManager struct {

	// internal

	client *NuxeoClient
	logger *slog.Logger
}

func (cm *ConfigManager) FetchTypes(ctx context.Context) (*DocTypes, error) {
	res, err := cm.client.NewRequest(ctx, nil).SetResult(&DocTypes{}).SetError(&NuxeoError{}).Get(apiV1 + "/config/types")

	if err := handleNuxeoError(err, res); err != nil {
		cm.logger.Error("Failed to fetch document types", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*DocTypes), nil
}

func (cm *ConfigManager) FetchType(ctx context.Context, name string) (*DocType, error) {
	res, err := cm.client.NewRequest(ctx, nil).SetResult(&DocType{}).SetError(&NuxeoError{}).Get(apiV1 + "/config/types/" + name)

	if err := handleNuxeoError(err, res); err != nil {
		cm.logger.Error("Failed to fetch document type", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*DocType), nil
}

func (cm *ConfigManager) FetchSchemas(ctx context.Context) (*Schemas, error) {
	res, err := cm.client.NewRequest(ctx, nil).SetResult(&Schemas{}).SetError(&NuxeoError{}).Get(apiV1 + "/config/schemas")

	if err := handleNuxeoError(err, res); err != nil {
		cm.logger.Error("Failed to fetch schemas", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Schemas), nil
}

func (cm *ConfigManager) FetchSchema(ctx context.Context, name string) (*Schema, error) {
	res, err := cm.client.NewRequest(ctx, nil).SetResult(&Schema{}).SetError(&NuxeoError{}).Get(apiV1 + "/config/schemas/" + name)

	if err := handleNuxeoError(err, res); err != nil {
		cm.logger.Error("Failed to fetch schema", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Schema), nil
}

func (cm *ConfigManager) FetchFacets(ctx context.Context) (*Facets, error) {
	res, err := cm.client.NewRequest(ctx, nil).SetResult(&Facets{}).SetError(&NuxeoError{}).Get(apiV1 + "/config/facets")

	if err := handleNuxeoError(err, res); err != nil {
		cm.logger.Error("Failed to fetch facets", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Facets), nil
}

func (cm *ConfigManager) FetchFacet(ctx context.Context, name string) (*Facet, error) {
	res, err := cm.client.NewRequest(ctx, nil).SetResult(&Facet{}).SetError(&NuxeoError{}).Get(apiV1 + "/config/facets/" + name)

	if err := handleNuxeoError(err, res); err != nil {
		cm.logger.Error("Failed to fetch facet", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Facet), nil
}
