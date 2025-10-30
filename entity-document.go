package nuxeo

import (
	"slices"
)

// EntityDocument represents a Nuxeo document entity, including metadata, properties, facets, and file content.
// It maps to the Nuxeo REST API document model.
type Document struct {
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
func NewDocument(documentType string, name string) *Document {
	return &Document{
		entity: entity{
			EntityType: EntityTypeDocument,
		},
		Type: documentType,
		Name: name,
		Properties: map[string]Field{
			DocumentPropertyDCTitle: NewStringField(name),
		},
	}
}

// HasFacet returns true if the document has the specified facet.
func (d *Document) HasFacet(facet string) bool {
	return slices.Contains(d.Facets, facet)
}

// IsFolder returns true if the document is folderish (can contain children).
func (d *Document) IsFolder() bool {
	return d.HasFacet("Folderish")
}

// IsCollection returns true if the document is a collection.
func (d *Document) IsCollection() bool {
	return d.HasFacet("Collection")
}

// IsCollectable returns true if the document can be collected (not a collection member).
func (d *Document) IsCollectable() bool {
	return d.HasFacet("NotCollectionMember")
}

// Property returns the value of the specified property key.
func (d *Document) Property(key string) (Field, bool) {
	value, found := d.Properties[key]
	return value, found
}

// SetProperty sets the value of the specified property key.
func (d *Document) SetProperty(key string, value Field) {
	d.Properties[key] = value
}

// FileContent returns the main file Blob of the document, if present.
func (d *Document) FileContent() *blob {
	if fieldBlob, ok := d.Properties[DocumentPropertyFileContent]; ok {
		var blob blob
		if err := fieldBlob.Complex(&blob); err == nil {
			return &blob
		}
	}
	return nil
}

// UploadInfo represents the upload information for a blob in Nuxeo.
type UploadInfo struct {
	Batch  string `json:"upload-batch"`
	FileId string `json:"upload-fileId"`
}

type uploadFileInfo struct {
	File UploadInfo `json:"file"`
}

// SetUploadInfoProperty sets the upload information for the document's file content property such as "file:content", "files:files", etc.
func (d *Document) SetUploadInfoProperty(key string, infos ...UploadInfo) {
	if len(infos) == 1 {
		uploadInfoFld, _ := NewComplexField(infos[0])
		d.SetProperty(key, uploadInfoFld)
	} else if len(infos) > 1 {
		uploadInfos := make([]uploadFileInfo, len(infos))
		for i, bi := range infos {
			uploadInfos[i] = uploadFileInfo{File: bi}
		}
		uploadInfoFld, _ := NewComplexField(uploadInfos)
		d.SetProperty(key, uploadInfoFld)
	}
}

// Thumbnail returns the thumbnail Blob of the document, if present.
func (d *Document) Thumbnail() *blob {
	if fieldBlob, ok := d.Properties[DocumentPropertyThumbThumbnail]; ok {
		var blob blob
		if err := fieldBlob.Complex(&blob); err == nil {
			return &blob
		}
	}
	return nil
}

// EntityDocuments is a paginated collection of EntityDocument objects.
type Documents paginableEntities[Document]
