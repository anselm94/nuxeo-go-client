package nuxeo

import (
	"slices"
)

// EntityDocument represents a Nuxeo document entity, including metadata, properties, facets, and file content.
// It maps to the Nuxeo REST API document model.
type entityDocument struct {
	entity
	Repository                  string           `json:"repository"`
	ID                          string           `json:"uid"`
	Path                        string           `json:"path"`
	Type                        string           `json:"type"`
	State                       string           `json:"state"`
	ParentRef                   string           `json:"parentRef"`
	IsCheckedOut                bool             `json:"isCheckedOut"`
	IsVersion                   bool             `json:"isVersion"`
	IsProxy                     bool             `json:"isProxy"`
	ProxyTargetId               string           `json:"proxyTargetId"`
	VersionableId               string           `json:"versionableId"`
	ChangeToken                 string           `json:"changeToken"`
	IsTrashed                   bool             `json:"isTrashed"`
	Title                       string           `json:"title"`
	Name                        string           `json:"name"`
	VersionLabel                string           `json:"versionLabel"`
	LockOwner                   string           `json:"lockOwner"`
	LockCreated                 string           `json:"lockCreated"`
	LastModified                *ISO8601Time     `json:"lastModified"`
	IsRecord                    bool             `json:"isRecord"`
	RetainUntil                 string           `json:"retainUntil"`
	HasLegalHold                bool             `json:"hasLegalHold"`
	IsUnderRetentionOrLegalHold bool             `json:"isUnderRetentionOrLegalHold"`
	Properties                  map[string]Field `json:"properties"`
	Facets                      []string         `json:"facets"`
}

// NewDocument creates a new EntityDocument with the specified type and name.
// Sets EntityType to "document" and initializes properties.
func NewDocument(documentType string, name string) *entityDocument {
	return &entityDocument{
		entity: entity{
			EntityType: EntityTypeDocument,
		},
		Type:       documentType,
		Name:       name,
		Properties: make(map[string]Field),
	}
}

// HasFacet returns true if the document has the specified facet.
func (d *entityDocument) HasFacet(facet string) bool {
	return slices.Contains(d.Facets, facet)
}

// IsFolder returns true if the document is folderish (can contain children).
func (d *entityDocument) IsFolder() bool {
	return d.HasFacet("Folderish")
}

// IsCollection returns true if the document is a collection.
func (d *entityDocument) IsCollection() bool {
	return d.HasFacet("Collection")
}

// IsCollectable returns true if the document can be collected (not a collection member).
func (d *entityDocument) IsCollectable() bool {
	return d.HasFacet("NotCollectionMember")
}

// Property returns the value of the specified property key.
func (d *entityDocument) Property(key string) (Field, bool) {
	value, found := d.Properties[key]
	return value, found
}

// SetProperty sets the value of the specified property key.
func (d *entityDocument) SetProperty(key string, value any) {
	fieldVal, _ := NewField(value)
	d.Properties[key] = fieldVal
}

// FileContent returns the main file Blob of the document, if present.
func (d *entityDocument) FileContent() *blob {
	if fieldBlob, ok := d.Properties[DocumentPropertyFileContent]; ok {
		var blob blob
		if err := fieldBlob.Complex(&blob); err == nil {
			return &blob
		}
	}
	return nil
}

// Thumbnail returns the thumbnail Blob of the document, if present.
func (d *entityDocument) Thumbnail() *blob {
	if fieldBlob, ok := d.Properties[DocumentPropertyThumbThumbnail]; ok {
		var blob blob
		if err := fieldBlob.Complex(&blob); err == nil {
			return &blob
		}
	}
	return nil
}

// EntityDocuments is a paginated collection of EntityDocument objects.
type entityDocuments paginableEntities[entityDocument]
