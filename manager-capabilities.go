package nuxeo

import (
	"context"
	"fmt"
	"log/slog"
)

type CapabilitiesManager struct {
	client *NuxeoClient
	logger *slog.Logger
}

func (cm *CapabilitiesManager) FetchCapabilities(ctx context.Context) (*Capabilities, error) {
	capabilities := &Capabilities{}
	res, err := cm.client.NewRequest(ctx, nil).SetResult(capabilities).Get("/api/v1/capabilities")

	if err != nil || res.StatusCode() != 200 {
		cm.logger.Error("Failed to get server capabilities", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to get server capabilities: %d %w", res.StatusCode(), err)
	}

	return capabilities, nil
}
