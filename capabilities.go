package nuxeo

import (
	"context"
	"fmt"
)

type Capabilities struct {
	EntityType string `json:"entity-type"`
	Server     struct {
		DistributionName    string `json:"distributionName"`
		DistributionVersion string `json:"distributionVersion"`
		DistributionServer  string `json:"distributionServer"`
	} `json:"server"`
	Cluster struct {
		Enabled bool   `json:"enabled"`
		NodeID  string `json:"nodeId"`
	} `json:"cluster"`
	Repository map[string]struct {
		QueryBlobKeys bool `json:"queryBlobKeys"`
	} `json:"repository"`
}

func (c *NuxeoClient) Capabilities(ctx context.Context) (*Capabilities, error) {
	capabilities := &Capabilities{}
	res, err := c.NewRequest(ctx).SetResult(capabilities).Get("/api/v1/capabilities")

	if err != nil || res.StatusCode() != 200 {
		c.logger.Error("Failed to get server capabilities", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to get server capabilities: %d %w", res.StatusCode(), err)
	}

	return capabilities, nil
}
