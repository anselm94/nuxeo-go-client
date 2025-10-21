package nuxeo

// Group represents a Nuxeo group.
type Group struct {
	Name       string
	Properties map[string]any
}

type Groups PaginatedEntities[Group]
