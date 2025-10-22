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

type DirectoryEntries entities[DirectoryEntry]
