package nuxeo

import "context"

// Repository represents a Nuxeo repository.
type Repository struct {
	name string

	// internal

	client *NuxeoClient
}

///////////////////
//// DOCUMENTS ////
///////////////////

func (r *Repository) FetchDocumentRoot(ctx context.Context, options *nuxeoRequestOptions) (*Document, error) {
	return nil, nil
}

func (r *Repository) FetchDocumentById(ctx context.Context, documentID string, options *nuxeoRequestOptions) (*Document, error) {
	return nil, nil
}

func (r *Repository) FetchDocumentByPath(ctx context.Context, documentPath string, options *nuxeoRequestOptions) (*Document, error) {
	return nil, nil
}

func (r *Repository) CreateDocumentById(ctx context.Context, parentId string, document Document, options *nuxeoRequestOptions) (*Document, error) {
	return nil, nil
}

func (r *Repository) CreateDocumentByPath(ctx context.Context, parentPath string, document Document, options *nuxeoRequestOptions) (*Document, error) {
	return nil, nil
}

func (r *Repository) UpdateDocument(ctx context.Context, documentId string, document Document, options *nuxeoRequestOptions) (*Document, error) {
	return nil, nil
}

func (r *Repository) DeleteDocument(ctx context.Context, documentId string) error {
	return nil
}

///////////////
//// QUERY ////
///////////////

func (r *Repository) Query(ctx context.Context, query string, queryParams []string, paginationOptions *SortedPaginationOptions, options *nuxeoRequestOptions) (*Documents, error) {
	return nil, nil
}

func (r *Repository) QueryByProvider(ctx context.Context, providerName string, queryParams []string, namedQueryParams map[string]string, paginationOptions *SortedPaginationOptions, options *nuxeoRequestOptions) (*Documents, error) {
	return nil, nil
}

///////////////
//// AUDIT ////
///////////////

func (r *Repository) FetchAuditByPath(ctx context.Context, documentPath string, options *nuxeoRequestOptions) (*Audit, error) {
	return nil, nil
}

func (r *Repository) FetchAuditById(ctx context.Context, documentId string, options *nuxeoRequestOptions) (*Audit, error) {
	return nil, nil
}

/////////////
//// ACP ////
/////////////

func (r *Repository) FetchPermissionsByPath(ctx context.Context, documentPath string, options *nuxeoRequestOptions) (*ACP, error) {
	return nil, nil
}

func (r *Repository) FetchPermissionsById(ctx context.Context, documentId string, options *nuxeoRequestOptions) (*ACP, error) {
	return nil, nil
}

//////////////////
//// CHILDREN ////
//////////////////

func (r *Repository) FetchChildrenByPath(ctx context.Context, parentPath string, options *nuxeoRequestOptions) (*Documents, error) {
	return nil, nil
}

func (r *Repository) FetchChildrenById(ctx context.Context, parentId string, options *nuxeoRequestOptions) (*Documents, error) {
	return nil, nil
}

///////////////
//// BLOBS ////
///////////////

func (r *Repository) StreamBlobByPath(ctx context.Context, documentPath string, blobXPath string, options *nuxeoRequestOptions) (*Blob, error) {
	return nil, nil
}

func (r *Repository) StreamBlobById(ctx context.Context, documentId string, blobXPath string, options *nuxeoRequestOptions) (*Blob, error) {
	return nil, nil
}

///////////////////
//// WORKFLOWS ////
///////////////////

func (r *Repository) StartWorkflowInstanceWithDocId(ctx context.Context, documentId string, workflow Workflow, options *nuxeoRequestOptions) (*Workflow, error) {
	return nil, nil
}

func (r *Repository) StartWorkflowInstanceWithDocPath(ctx context.Context, documentPath string, workflow Workflow, options *nuxeoRequestOptions) (*Workflow, error) {
	return nil, nil
}

func (r *Repository) FetchWorkflowInstancesByDocId(ctx context.Context, documentId string, options *nuxeoRequestOptions) (*Workflows, error) {
	return nil, nil
}

func (r *Repository) FetchWorkflowInstancesByDocPath(ctx context.Context, documentPath string, options *nuxeoRequestOptions) (*Workflows, error) {
	return nil, nil
}

func (r *Repository) FetchWorkflowInstance(ctx context.Context, workflowInstanceId string, options *nuxeoRequestOptions) (*Workflow, error) {
	return nil, nil
}

func (r *Repository) CancelWorkflowInstance(ctx context.Context, workflowInstanceId string) error {
	return nil
}

func (r *Repository) FetchWorkflowInstanceGraph(ctx context.Context, workflowInstanceId string, options *nuxeoRequestOptions) (*WorkflowGraph, error) {
	return nil, nil
}

func (r *Repository) FetchWorkflowModel(ctx context.Context, workflowModelName string, options *nuxeoRequestOptions) (*Workflow, error) {
	return nil, nil
}

func (r *Repository) FetchWorkflowModelGraph(ctx context.Context, workflowModelName string, options *nuxeoRequestOptions) (*WorkflowGraph, error) {
	return nil, nil
}

func (r *Repository) FetchWorkflowModels(ctx context.Context, options *nuxeoRequestOptions) (*Workflows, error) {
	return nil, nil
}

/////////////////////
//// WEB ADAPTER ////
/////////////////////

func (r *Repository) CreateForAdapter(ctx context.Context, documentId string, adapter string, pathSuffix string, queryParams []string, payload any, options *nuxeoRequestOptions) (*any, error) {
	return nil, nil
}

func (r *Repository) FetchForAdapter(ctx context.Context, documentId string, adapter string, pathSuffix string, queryParams []string, options *nuxeoRequestOptions) (*any, error) {
	return nil, nil
}

func (r *Repository) UpdateForAdapter(ctx context.Context, documentId string, adapter string, pathSuffix string, queryParams []string, payload any, options *nuxeoRequestOptions) (*any, error) {
	return nil, nil
}

func (r *Repository) DeleteForAdapter(ctx context.Context, documentId string, adapter string, pathSuffix string, queryParams []string) error {
	return nil
}
