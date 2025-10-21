package nuxeo

// Group represents a Nuxeo group.
type Group struct {
	Name       string
	Properties map[string]any
}

// NewGroup creates a new Group instance.
func NewGroup(name string, props map[string]any) *Group {
	return &Group{
		Name:       name,
		Properties: props,
	}
}
