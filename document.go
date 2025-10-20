package nuxeo

// Document represents a Nuxeo document entity.
type Document struct {
	ID           string         `json:"uid"`
	Path         string         `json:"path"`
	Type         string         `json:"type"`
	State        string         `json:"state"`
	Title        string         `json:"title"`
	LastModified string         `json:"lastModified"`
	Properties   map[string]any `json:"properties"`
}

type DocumentList struct {
	EntityType string     `json:"entity-type"`
	Entries    []Document `json:"entries"`
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
