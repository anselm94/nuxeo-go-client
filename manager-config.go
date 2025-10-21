package nuxeo

import (
	"context"
	"fmt"
	"log/slog"
)

type ConfigManager struct {

	// internal

	client *NuxeoClient
	logger *slog.Logger
}

func (cm *ConfigManager) FetchTypes(ctx context.Context) (*DocTypes, error) {
	res, err := cm.client.NewRequest(ctx).SetResult(&DocTypes{}).Get("/api/v1/config/types")

	if err != nil || res.StatusCode() != 200 {
		cm.logger.Error("Failed to fetch document types", "error", err)
		return nil, fmt.Errorf("failed to fetch document types: %w", err)
	}

	return res.Result().(*DocTypes), nil
}

func (cm *ConfigManager) FetchType(ctx context.Context, name string) (*DocType, error) {
	res, err := cm.client.NewRequest(ctx).SetResult(&DocType{}).Get("/api/v1/config/types/" + name)

	if err != nil || res.StatusCode() != 200 {
		cm.logger.Error("Failed to fetch document type", "error", err)
		return nil, fmt.Errorf("failed to fetch document type: %w", err)
	}

	return res.Result().(*DocType), nil
}

func (cm *ConfigManager) FetchSchemas(ctx context.Context) (*Schemas, error) {
	res, err := cm.client.NewRequest(ctx).SetResult(&Schemas{}).Get("/api/v1/config/schemas")

	if err != nil || res.StatusCode() != 200 {
		cm.logger.Error("Failed to fetch schemas", "error", err)
		return nil, fmt.Errorf("failed to fetch schemas: %w", err)
	}

	return res.Result().(*Schemas), nil
}

func (cm *ConfigManager) FetchSchema(ctx context.Context, name string) (*Schema, error) {
	res, err := cm.client.NewRequest(ctx).SetResult(&Schema{}).Get("/api/v1/config/schemas/" + name)

	if err != nil || res.StatusCode() != 200 {
		cm.logger.Error("Failed to fetch schema", "error", err)
		return nil, fmt.Errorf("failed to fetch schema: %w", err)
	}

	return res.Result().(*Schema), nil
}

func (cm *ConfigManager) FetchFacets(ctx context.Context) (*Facets, error) {
	res, err := cm.client.NewRequest(ctx).SetResult(&Facets{}).Get("/api/v1/config/facets")

	if err != nil || res.StatusCode() != 200 {
		cm.logger.Error("Failed to fetch facets", "error", err)
		return nil, fmt.Errorf("failed to fetch facets: %w", err)
	}

	return res.Result().(*Facets), nil
}

func (cm *ConfigManager) FetchFacet(ctx context.Context, name string) (*Facet, error) {
	res, err := cm.client.NewRequest(ctx).SetResult(&Facet{}).Get("/api/v1/config/facets/" + name)

	if err != nil || res.StatusCode() != 200 {
		cm.logger.Error("Failed to fetch facet", "error", err)
		return nil, fmt.Errorf("failed to fetch facet: %w", err)
	}

	return res.Result().(*Facet), nil
}
