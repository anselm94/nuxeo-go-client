package nuxeo

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
