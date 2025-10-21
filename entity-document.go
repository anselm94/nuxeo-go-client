package nuxeo

import (
	"context"
	"fmt"
	"mime"
	"slices"
	"time"
)

type AuditEntry struct {
	EntityType    string    `json:"entity-type"`
	ID            int       `json:"id"`
	Category      string    `json:"category"`
	PrincipalName string    `json:"principalName"`
	Comment       string    `json:"comment"`
	DocLifeCycle  string    `json:"docLifeCycle"`
	DocPath       string    `json:"docPath"`
	DocType       string    `json:"docType"`
	DocUUID       string    `json:"docUUID"`
	EventID       string    `json:"eventId"`
	RepositoryID  string    `json:"repositoryId"`
	EventDate     time.Time `json:"eventDate"`
	LogDate       time.Time `json:"logDate"`
	Extended      struct {
		BlobFilename   string `json:"blobFilename"`
		DownloadReason string `json:"downloadReason"`
		ClientReason   string `json:"clientReason"`
		BlobXPath      string `json:"blobXPath"`
	} `json:"extended"`
}

type ACL struct {
	Name string `json:"name"`
	ACEs []ACE  `json:"aces"`
}

type ACE struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	ExternalUser bool   `json:"externalUser"`
	Permission   string `json:"permission"`
	Granted      bool   `json:"granted"`
	Creator      string `json:"creator"`
	Begin        string `json:"begin"`
	End          string `json:"end"`
	Status       string `json:"status"`
}

type DocumentContextParameters struct {
	ACLs        []ACL        `json:"acls,omitempty"`
	Audit       []AuditEntry `json:"audit,omitempty"`
	Breadcrumb  DocumentList `json:"breadcrumb,omitempty"`
	Children    DocumentList `json:"children,omitempty"`
	Collections []any        `json:"collections,omitempty"`
	DocumentUrl string       `json:"documentUrl,omitempty"`
	Favorites   struct {
		IsFavorite bool `json:"isFavorite,omitempty"`
	} `json:"favorites,omitempty"`
	FirstAccessibleAncestor Document `json:"firstAccessibleAncestor,omitempty"`
	HasContent              bool     `json:"hasContent,omitempty"`
	HasFolderishChild       bool     `json:"hasFolderishContent,omitempty"`
	PendingTasks            []any    `json:"pendingTasks,omitempty"`
	Permissions             []string `json:"permissions,omitempty"`
	Preview                 struct {
		URL string `json:"url,omitempty"`
	} `json:"preview,omitempty"`
	Publications struct {
		ResultsCount int `json:"resultsCount,omitempty"`
	} `json:"publications,omitempty"`
	Subscriptions struct {
		IsSubscribed bool `json:"isSubscribed,omitempty"`
	}
	RunningWorkflows        []any `json:"runningWorkflows,omitempty"`
	SubscribedNotifications []any `json:"subscribedNotifications,omitempty"`
	SubTypes                []any `json:"subtypes,omitempty"`
	Tags                    []any `json:"tags,omitempty"`
	Thumbnail               struct {
		URL string `json:"url,omitempty"`
	} `json:"thumbnail,omitempty"`
	UserVisiblePermissions []string `json:"userVisiblePermissions,omitempty"`
}

type DocumentList struct {
	EntityType string     `json:"entity-type"`
	Entries    []Document `json:"entries"`
}

// Document represents a Nuxeo document entity.
type Document struct {
	EntityType                  string         `json:"entity-type"`
	Repository                  string         `json:"repository"`
	ID                          string         `json:"uid"`
	Path                        string         `json:"path"`
	Type                        string         `json:"type"`
	State                       string         `json:"state"`
	ParentRef                   string         `json:"parentRef"`
	IsCheckedOut                bool           `json:"isCheckedOut"`
	IsRecord                    bool           `json:"isRecord"`
	RetainUntil                 string         `json:"retainUntil"`
	HasLegalHold                bool           `json:"hasLegalHold"`
	IsUnderRetentionOrLegalHold bool           `json:"isUnderRetentionOrLegalHold"`
	IsVersion                   bool           `json:"isVersion"`
	IsProxy                     bool           `json:"isProxy"`
	ChangeToken                 string         `json:"changeToken"`
	IsTrashed                   bool           `json:"isTrashed"`
	Title                       string         `json:"title"`
	LastModified                string         `json:"lastModified"`
	Properties                  map[string]any `json:"properties"`
	Facets                      []string       `json:"facets"`
	Schemas                     []Schema       `json:"schemas"`
	ContextParameters           map[string]any `json:"contextParameters"`
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

func (c *NuxeoClient) FetchBlob(ctx context.Context, documentId string, xPath string, options *NuxeoRequestOptions) (*Blob, error) {
	if xPath == "" {
		xPath = "blobholder:0"
	}
	res, err := c.NewRequest(ctx).SetNuxeoOption(options).Get("/api/v1/id/" + documentId + "/@blob/" + xPath)
	if err != nil || res.StatusCode() != 200 {
		c.logger.Error("Failed to fetch blob", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to fetch blob: %d %w", res.StatusCode(), err)
	}

	// extract content type
	contentType := res.Header().Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// extract filename from Content-Disposition header
	disposition := res.Header().Get("Content-Disposition")
	filename := ""
	if disposition != "" {
		_, params, err := mime.ParseMediaType(disposition)
		if err == nil {
			filename = params["filename"]
		}
	}

	return &Blob{
		Filename: filename,
		MimeType: contentType,
		Data:     res.Body,
	}, nil
}

// func (c *NuxeoClient) MoveDocument(ctx context.Context, documentId string, destinationId string, newName string, options *NuxeoRequestOptions) (*Document, error) {
// 	movedDoc := &Document{}

// 	err := c.NewOperation(ctx, "Document.Move", options).
// 		SetDocumentInput(documentId).
// 		SetParam("name", newName).
// 		SetParam("target", destinationId).
// 		ExecuteInto(movedDoc)

// 	if err != nil {
// 		c.logger.Error("Failed to move document", "error", err)
// 		return nil, fmt.Errorf("failed to move document: %w", err)
// 	}
// 	return movedDoc, nil
// }

// func (c *NuxeoClient) FollowTransitionForDocument(ctx context.Context, documentId string, transitionName string, options *NuxeoRequestOptions) (*Document, error) {
// 	updatedDoc := &Document{}

// 	err := c.NewOperation(ctx, "Document.FollowLifecycleTransition", options).
// 		SetDocumentInput(documentId).
// 		SetParam("value", transitionName).
// 		ExecuteInto(updatedDoc)

// 	if err != nil {
// 		c.logger.Error("Failed to follow transition for document", "error", err)
// 		return nil, fmt.Errorf("failed to follow transition for document: %w", err)
// 	}
// 	return updatedDoc, nil
// }

// Configuration options for the conversion.
// At least one of the 'converter', 'type' or 'format' option must be defined
type DocumentConversionOptions struct {
	// The Blob xpath. Default to the main blob 'blobholder:0'
	XPath string
	// Named converter to use
	ConverterName string
	// The destination mime type, such as 'application/pdf'
	MimeType string
	// The destination format, such as 'pdf'
	Format string
}

// func (c *NuxeoClient) ConvertDocument(ctx context.Context, documentId string, conversionOptions DocumentConversionOptions, options *NuxeoRequestOption) error {

// }
