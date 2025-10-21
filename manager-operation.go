package nuxeo

import (
	"context"
	"log/slog"
)

type OperationManager struct {
	client *NuxeoClient
	logger *slog.Logger
}

func (om *OperationManager) NewOperation(ctx context.Context, automationId string, options *NuxeoRequestOptions) *Operation {
	return &Operation{
		automationId: automationId,
		Params:       make(map[string]string),
		Context:      make(map[string]string),
		request:      om.client.NewRequest(ctx).SetNuxeoOption(options),
		logger:       *om.logger,
	}
}
