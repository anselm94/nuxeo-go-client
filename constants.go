package nuxeo

// Headers

const (
	HeaderAuthorization        = "Authorization"
	HeaderContentType          = "Content-Type"
	HeaderContentLength        = "Content-Length"
	HeaderDepth                = "depth"
	HeaderTimeout              = "timeout"
	HeaderProperties           = "properties"
	HeaderNuxeoTxTimeout       = "Nuxeo-Transaction-Timeout"
	HeaderNxUser               = "NX_USER"
	HeaderNxToken              = "NX_TOKEN"
	HeaderNxRd                 = "NX_RD"
	HeaderNxTs                 = "NX_TS"
	HeaderNxEsSync             = "nx_es_sync"
	HeaderUserAgent            = "User-Agent"
	HeaderXAuthenticationToken = "X-Authentication-Token"
	HeaderXRepository          = "X-NXRepository"
	HeaderXProperties          = "X-NXproperties"
	HeaderXVoidOperation       = "X-NXVoidOperation"
	HeaderXVersioningOption    = "X-Versioning-Option"
)

const (
	HeaderValueOctetStream = "application/octet-stream"
)

// Entity Types

const (
	EntityTypeACP              = "acls"
	EntityTypeAnnotation       = "annotation"
	EntityTypeAnnotations      = "annotations"
	EntityTypeCapabilities     = "capabilities"
	EntityTypeComment          = "comment"
	EntityTypeComments         = "comments"
	EntityTypeAudit            = "logEntries"
	EntityTypeBlobs            = "blobs"
	EntityTypeDirectories      = "directories"
	EntityTypeDirectory        = "directory"
	EntityTypeDirectoryEntries = "directoryEntries"
	EntityTypeDirectoryEntry   = "directoryEntry"
	EntityTypeDocument         = "document"
	EntityTypeDocuments        = "documents"
	EntityTypeDocType          = "docType"
	EntityTypeDocTypes         = "docTypes"
	EntityTypeException        = "exception"
	EntityTypeFacet            = "facet"
	EntityTypeGraph            = "graph"
	EntityTypeGroup            = "group"
	EntityTypeGroups           = "groups"
	EntityTypeLogEntry         = "logEntry"
	EntityTypeLogin            = "login"
	EntityTypeOperation        = "operation"
	EntityTypeRecordSet        = "recordSet"
	EntityTypeSchema           = "schema"
	EntityTypeString           = "string"
	EntityTypeTask             = "task"
	EntityTypeTasks            = "tasks"
	EntityTypeUser             = "user"
	EntityTypeUsers            = "users"
	EntityTypeWorkflow         = "workflow"
	EntityTypeWorkflows        = "workflows"
)

// Repository

const (
	RepositoryDefault = "default"
)

// User

const (
	UserPropertyFirstName = "firstName"
	UserPropertyLastName  = "lastName"
	UserPropertyEmail     = "email"
	UserPropertyGroups    = "groups"
	UserPropertyUsername  = "username"
	UserPropertyCompany   = "company"
	UserPropertyPassword  = "password"
	UserPropertyTenantId  = "tenantId"
)

// ACL

const (
	AclLocal   = "local"
	AclInherit = "inherited"
)

// Document

const (
	DocumentStateDeleted = "deleted"
)

const (
	DocumentPropertyFileContent = "file:content"
)

// Directory

const (
	DirectoryPropertyId       = "id"
	DirectoryPropertyLabel    = "label"
	DirectoryPropertyOrdering = "ordering"
	DirectoryPropertyObsolete = "obsolete"
)

// Operations

const (
	OperationBlobAttachOnDocument       = "Blob.AttachOnDocument"
	OperationDirectoryEntries           = "Directory.Entries"
	OperationDocumentAddPermission      = "Document.AddPermission"
	OperationDocumentRemovePermission   = "Document.RemovePermission"
	OperationDocumentRemoveProxies      = "Document.RemoveProxies"
	OperationDocumentCheckIn            = "Document.CheckIn"
	OperationDocumentGetLastVersion     = "Document.GetLastVersion"
	OperationDocumentGetBlob            = "Document.GetBlob"
	OperationDocumentGetBlobs           = "Document.GetBlobs"
	OperationDocumentGetBlobsByProperty = "Document.GetBlobsByProperty"
	OperationDocumentTrash              = "Document.Trash"
	OperationDocumentUntrash            = "Document.Untrash"
	OperationDocumentUpdate             = "Document.Update"
	OperationEsWaitForIndexing          = "Elasticsearch.WaitForIndexing"
	OperationRepositoryGetDocument      = "Repository.GetDocument"
)
