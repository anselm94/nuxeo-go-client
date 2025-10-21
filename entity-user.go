package nuxeo

// User represents a Nuxeo user.
type User struct {
	EntityType string         `json:"entity-type"`
	Id         string         `json:"id"`
	Properties map[string]any `json:"properties"`
}
