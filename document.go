package nuxeo

// Document represents a Nuxeo document entity.
type Document struct {
	ID         string
	Type       string
	Path       string
	Properties map[string]any
}

// NewDocument creates a new Document instance.
func NewDocument(id, docType, path string, props map[string]any) *Document {
	return &Document{
		ID:         id,
		Type:       docType,
		Path:       path,
		Properties: props,
	}
}

// SetProperty sets a property on the document.
func (d *Document) SetProperty(key string, value any) {
	d.Properties[key] = value
}

// GetProperty retrieves a property from the document.
func (d *Document) GetProperty(key string) any {
	return d.Properties[key]
}
