package nuxeo

import (
	"context"
	"log/slog"

	"github.com/anselm94/nuxeo-go-client/internal"
)

type capabilitiesManager struct {
	client *NuxeoClient
	logger *slog.Logger
}

func (cm *capabilitiesManager) FetchCapabilities(ctx context.Context) (*entityCapabilities, error) {
	path := internal.PathApiV1 + "/capabilities"
	res, err := cm.client.NewRequest(ctx, nil).SetResult(&entityCapabilities{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		cm.logger.Error("Failed to fetch capabilities", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityCapabilities), nil
}
