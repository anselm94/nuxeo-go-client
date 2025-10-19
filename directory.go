package nuxeo

// DirectoryEntry represents an entry in a Nuxeo directory.
type DirectoryEntry struct {
	ID         string
	Properties map[string]any
}

// Directory represents a Nuxeo directory.
type Directory struct {
	Name string
}

// NewDirectory creates a new Directory instance.
func NewDirectory(name string) *Directory {
	return &Directory{Name: name}
}

// NewDirectoryEntry creates a new DirectoryEntry instance.
func NewDirectoryEntry(id string, props map[string]any) *DirectoryEntry {
	return &DirectoryEntry{
		ID:         id,
		Properties: props,
	}
}
