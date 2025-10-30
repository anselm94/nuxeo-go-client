package nuxeo

// Directory represents a Nuxeo Directory entity.
// See: https://doc.nuxeo.com/rest-api/1/directory-entity-type/
type Directory struct {
	entity
	Name    string `json:"name"`
	Schema  string `json:"schema"`
	IdField string `json:"idField"`
	Parent  string `json:"parent"`
}

// Directories is a collection of Directory entities returned by GET /directory.
// See: https://doc.nuxeo.com/rest-api/1/directory-endpoint/
type Directories entities[Directory]

// DirectoryEntry represents a Nuxeo Directory Entry entity.
// See: https://doc.nuxeo.com/rest-api/1/directory-entry-entity-type/
type DirectoryEntry struct {
	entity
	DirectoryName string           `json:"directoryName"`
	ID            string           `json:"id"`
	Properties    map[string]Field `json:"properties"`
}

// NewDirectoryEntry creates a new Directory Entry with the given id.
// The Properties map is initialized empty.
func NewDirectoryEntry(id string) DirectoryEntry {
	return DirectoryEntry{
		ID:         id,
		Properties: make(map[string]Field),
	}
}

// Id returns the entry's id, falling back to the 'id' property if not set.
func (d *DirectoryEntry) Id() string {
	if d.ID != "" {
		return d.ID
	}
	if val, ok := d.Properties[DirectoryPropertyId]; ok {
		if idStr, err := val.String(); err == nil && idStr != nil {
			return *idStr
		}
	}
	return ""
}

// Label returns the entry's label property, if present.
func (d DirectoryEntry) Label() string {
	if val, ok := d.Properties[DirectoryPropertyLabel]; ok {
		if labelStr, err := val.String(); err == nil && labelStr != nil {
			return *labelStr
		}
	}
	return ""
}

// Ordering returns the entry's ordering property, if present.
func (d DirectoryEntry) Ordering() float64 {
	if val, ok := d.Properties[DirectoryPropertyOrdering]; ok {
		if orderingFloat, err := val.Float(); err == nil && orderingFloat != nil {
			return *orderingFloat
		}
	}
	return 0
}

// Obsolete returns the entry's obsolete property, if present.
func (d DirectoryEntry) Obsolete() float64 {
	if val, ok := d.Properties[DirectoryPropertyObsolete]; ok {
		if obsoleteFloat, err := val.Float(); err == nil && obsoleteFloat != nil {
			return *obsoleteFloat
		}
	}
	return 0
}

// Property returns the value of the given property key for the entry.
func (d DirectoryEntry) Property(key string) (Field, bool) {
	value, found := d.Properties[key]
	return value, found
}

// SetProperty sets the value of the given property key for the entry.
func (d *DirectoryEntry) SetProperty(key string, value Field) {
	d.Properties[key] = value
}

// DirectoryEntries is a collection of Directory Entry entities returned by GET /directory/{directoryName}.
// See: https://doc.nuxeo.com/rest-api/1/directory-endpoint/
type DirectoryEntries entities[DirectoryEntry]
