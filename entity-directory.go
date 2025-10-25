package nuxeo

type entityDirectory struct {
	entity
	Name    string `json:"name"`
	Schema  string `json:"schema"`
	IdField string `json:"idField"`
	Parent  string `json:"parent"`
}

type entityDirectories entities[entityDirectory]

type entityDirectoryEntry struct {
	entity
	DirectoryName string         `json:"directoryName"`
	ID            string         `json:"id"`
	Properties    map[string]any `json:"properties"`
}

func NewDirectoryEntry(id string) entityDirectoryEntry {
	return entityDirectoryEntry{
		ID:         id,
		Properties: make(map[string]any),
	}
}

func (d *entityDirectoryEntry) Id() string {
	if d.ID != "" {
		return d.ID
	}
	if val, ok := d.Properties[DirectoryPropertyId]; ok {
		return val.(string)
	}
	return ""
}

func (d entityDirectoryEntry) Label() string {
	return d.Properties[DirectoryPropertyLabel].(string)
}

func (d entityDirectoryEntry) Ordering() float64 {
	return d.Properties[DirectoryPropertyOrdering].(float64)
}

func (d entityDirectoryEntry) Obsolete() float64 {
	return d.Properties[DirectoryPropertyObsolete].(float64)
}

func (d entityDirectoryEntry) Property(key string) any {
	return d.Properties[key]
}

type entityDirectoryEntries entities[entityDirectoryEntry]
