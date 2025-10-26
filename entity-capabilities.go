package nuxeo

// entityCapabilities represents the server capabilities as returned by the Nuxeo REST API /capabilities endpoint.
// See: https://doc.nuxeo.com/rest-api/1/capabilities-endpoint/
type entityCapabilities struct {
	entity
	// Server contains Nuxeo server distribution information.
	Server struct {
		DistributionName    string `json:"distributionName"`    // Name of the Nuxeo distribution (e.g., "lts").
		DistributionVersion string `json:"distributionVersion"` // Version of the Nuxeo distribution (e.g., "2021.39.3").
		DistributionServer  string `json:"distributionServer"`  // Server type (e.g., "tomcat").
	} `json:"server"`
	// Cluster contains cluster-related capabilities.
	Cluster struct {
		Enabled bool   `json:"enabled"` // True if clustering is enabled.
		NodeID  string `json:"nodeId"`  // Unique node identifier.
	} `json:"cluster"`
	// Repository contains repository-specific capabilities, keyed by repository name.
	Repository map[string]struct {
		QueryBlobKeys bool `json:"queryBlobKeys"` // True if querying blob keys is supported.
	} `json:"repository"`
}
