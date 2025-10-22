package nuxeo

import (
	"slices"
	"time"
)

// Document represents a Nuxeo document entity.
type Document struct {
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
	LastModified                time.Time      `json:"lastModified"`
	IsRecord                    bool           `json:"isRecord"`
	RetainUntil                 string         `json:"retainUntil"`
	HasLegalHold                bool           `json:"hasLegalHold"`
	IsUnderRetentionOrLegalHold bool           `json:"isUnderRetentionOrLegalHold"`
	Properties                  map[string]any `json:"properties"`
	Facets                      []string       `json:"facets"`
}

func (d *Document) HasFacet(facet string) bool {
	return slices.Contains(d.Facets, facet)
}

func (d *Document) IsFolder() bool {
	return d.HasFacet("Folderish")
}

func (d *Document) IsCollection() bool {
	return d.HasFacet("Collection")
}

func (d *Document) IsCollectable() bool {
	return d.HasFacet("NotCollectionMember")
}

type Documents paginableEntities[Document]
