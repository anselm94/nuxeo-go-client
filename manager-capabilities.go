package nuxeo

import (
	"context"
	"log/slog"

	"github.com/anselm94/nuxeo-go-client/internal"
)

// capabilitiesManager provides access to the Nuxeo /capabilities REST API endpoint.
// See: https://doc.nuxeo.com/rest-api/1/capabilities-endpoint/
//
// This manager allows fetching server, cluster, and repository capabilities.
type capabilitiesManager struct {
	client *NuxeoClient
	logger *slog.Logger
}

// FetchCapabilities retrieves the server capabilities from the Nuxeo /capabilities endpoint.
//
// Endpoint: GET /capabilities
// Returns an entityCapabilities struct describing server, cluster, and repository features.
// See: https://doc.nuxeo.com/rest-api/1/capabilities-endpoint/
//
// Example usage:
//
//	capabilities, err := cm.FetchCapabilities(ctx)
func (cm *capabilitiesManager) FetchCapabilities(ctx context.Context) (*Capabilities, error) {
	path := internal.PathApiV1 + "/capabilities"
	res, err := cm.client.NewRequest(ctx, nil).SetResult(&Capabilities{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		cm.logger.Error("Failed to fetch capabilities", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Capabilities), nil
}
