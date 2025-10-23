package nuxeo

import (
	"context"
	"log/slog"
)

type CapabilitiesManager struct {
	client *NuxeoClient
	logger *slog.Logger
}

func (cm *CapabilitiesManager) FetchCapabilities(ctx context.Context) (*Capabilities, error) {
	path := apiV1 + "/capabilities"
	res, err := cm.client.NewRequest(ctx, nil).SetResult(&Capabilities{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		cm.logger.Error("Failed to fetch capabilities", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Capabilities), nil
}
