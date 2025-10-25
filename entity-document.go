package nuxeo

import (
	"slices"
)

// entityDocument represents a Nuxeo document entity.
type entityDocument struct {
	entity
	Repository                  string         `json:"repository"`
	ID                          string         `json:"uid"`
	Path                        string         `json:"path"`
	Type                        string         `json:"type"`
	State                       string         `json:"state"`
	ParentRef                   string         `json:"parentRef"`
	IsCheckedOut                bool           `json:"isCheckedOut"`
	IsVersion                   bool           `json:"isVersion"`
	IsProxy                     bool           `json:"isProxy"`
	ProxyTargetId               string         `json:"proxyTargetId"`
	VersionableId               string         `json:"versionableId"`
	ChangeToken                 string         `json:"changeToken"`
	IsTrashed                   bool           `json:"isTrashed"`
	Title                       string         `json:"title"`
	VersionLabel                string         `json:"versionLabel"`
	LockOwner                   string         `json:"lockOwner"`
	LockCreated                 string         `json:"lockCreated"`
	LastModified                *ISO8601Time   `json:"lastModified"`
	IsRecord                    bool           `json:"isRecord"`
	RetainUntil                 string         `json:"retainUntil"`
	HasLegalHold                bool           `json:"hasLegalHold"`
	IsUnderRetentionOrLegalHold bool           `json:"isUnderRetentionOrLegalHold"`
	Properties                  map[string]any `json:"properties"`
	Facets                      []string       `json:"facets"`
}

func NewDocument(documentType string, name string) *entityDocument {
	return &entityDocument{
		entity: entity{
			EntityType: EntityTypeDocument,
		},
		Type:  documentType,
		Title: name,
	}
}

func (d *entityDocument) HasFacet(facet string) bool {
	return slices.Contains(d.Facets, facet)
}

func (d *entityDocument) IsFolder() bool {
	return d.HasFacet("Folderish")
}

func (d *entityDocument) IsCollection() bool {
	return d.HasFacet("Collection")
}

func (d *entityDocument) IsCollectable() bool {
	return d.HasFacet("NotCollectionMember")
}

type entityDocuments paginableEntities[entityDocument]
