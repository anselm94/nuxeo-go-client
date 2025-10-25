package nuxeo

type Directory struct {
	entity
	Name    string `json:"name"`
	Schema  string `json:"schema"`
	IdField string `json:"idField"`
	Parent  string `json:"parent"`
}

type Directories entities[Directory]

type DirectoryEntry struct {
	entity
	DirectoryName string         `json:"directoryName"`
	Id            string         `json:"id"`
	Properties    map[string]any `json:"properties"`
}

func (d DirectoryEntry) Label() string {
	return d.Properties[DirectoryPropertyLabel].(string)
}

func (d DirectoryEntry) Ordering() float64 {
	return d.Properties[DirectoryPropertyOrdering].(float64)
}

func (d DirectoryEntry) Obsolete() float64 {
	return d.Properties[DirectoryPropertyObsolete].(float64)
}

func (d DirectoryEntry) Property(key string) any {
	return d.Properties[key]
}

type DirectoryEntries entities[DirectoryEntry]
