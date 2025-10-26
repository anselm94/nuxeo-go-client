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
	DirectoryName string           `json:"directoryName"`
	ID            string           `json:"id"`
	Properties    map[string]Field `json:"properties"`
}

func NewDirectoryEntry(id string) entityDirectoryEntry {
	return entityDirectoryEntry{
		ID:         id,
		Properties: make(map[string]Field),
	}
}

func (d *entityDirectoryEntry) Id() string {
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

func (d entityDirectoryEntry) Label() string {
	if val, ok := d.Properties[DirectoryPropertyLabel]; ok {
		if labelStr, err := val.String(); err == nil && labelStr != nil {
			return *labelStr
		}
	}
	return ""
}

func (d entityDirectoryEntry) Ordering() float64 {
	if val, ok := d.Properties[DirectoryPropertyOrdering]; ok {
		if orderingFloat, err := val.Float(); err == nil && orderingFloat != nil {
			return *orderingFloat
		}
	}
	return 0
}

func (d entityDirectoryEntry) Obsolete() float64 {
	if val, ok := d.Properties[DirectoryPropertyObsolete]; ok {
		if obsoleteFloat, err := val.Float(); err == nil && obsoleteFloat != nil {
			return *obsoleteFloat
		}
	}
	return 0
}

func (d entityDirectoryEntry) Property(key string) Field {
	return d.Properties[key]
}

func (d *entityDirectoryEntry) SetProperty(key string, value any) {
	fieldVal, _ := NewField(value)
	d.Properties[key] = fieldVal
}

type entityDirectoryEntries entities[entityDirectoryEntry]
