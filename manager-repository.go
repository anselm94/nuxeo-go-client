package nuxeo

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/anselm94/nuxeo-go-client/internal"
)

// repository provides access to a Nuxeo repository and its document, query, audit, permission, workflow, and adapter APIs.
// It encapsulates repository-specific operations and maintains a reference to the client and logger.
//
// Use repository methods to interact with documents, run queries, manage audits, permissions, children, blobs, workflows, and invoke web adapters.
// See Nuxeo REST API documentation: https://doc.nuxeo.com/rest-api/repository-endpoint/
type repository struct {
	name string

	// internal

	client *NuxeoClient
	logger *slog.Logger
}

// Name returns the repository's name.
// This is the logical name used in API paths (e.g., "default").
func (r *repository) Name() string {
	return r.name
}

///////////////////
//// DOCUMENTS ////
///////////////////

// FetchDocumentRoot retrieves the root document of the repository.
// Maps to GET /api/v1/repo/{repo}/path/
// Returns the root entityDocument or error.
func (r *repository) FetchDocumentRoot(ctx context.Context, options *nuxeoRequestOptions) (*entityDocument, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path/"
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityDocument{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch document root", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityDocument), nil
}

// FetchDocumentById retrieves a document by its unique ID.
// Maps to GET /api/v1/repo/{repo}/id/{id}
// Returns the entityDocument or error.
func (r *repository) FetchDocumentById(ctx context.Context, documentID string, options *nuxeoRequestOptions) (*entityDocument, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentID)
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityDocument{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch document by ID", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityDocument), nil
}

// FetchDocumentByPath retrieves a document by its repository path.
// Maps to GET /api/v1/repo/{repo}/path/{path}
// Returns the entityDocument or error.
func (r *repository) FetchDocumentByPath(ctx context.Context, documentPath string, options *nuxeoRequestOptions) (*entityDocument, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path" + url.PathEscape(documentPath)
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityDocument{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch document by path", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityDocument), nil
}

// CreateDocumentById creates a new document under a parent document specified by ID.
// Maps to POST /api/v1/repo/{repo}/id/{parentId}
// Returns the created entityDocument or error.
func (r *repository) CreateDocumentById(ctx context.Context, parentId string, doc entityDocument, options *nuxeoRequestOptions) (*entityDocument, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(parentId)
	res, err := r.client.NewRequest(ctx, options).SetBody(doc).SetResult(&entityDocument{}).SetError(&nuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to create document by ID", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityDocument), nil
}

// CreateDocumentByPath creates a new document under a parent document specified by repository path.
// Maps to POST /api/v1/repo/{repo}/path/{parentPath}
// Returns the created entityDocument or error.
func (r *repository) CreateDocumentByPath(ctx context.Context, parentPath string, document entityDocument, options *nuxeoRequestOptions) (*entityDocument, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path" + url.PathEscape(parentPath)
	res, err := r.client.NewRequest(ctx, options).SetBody(document).SetResult(&entityDocument{}).SetError(&nuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to create document by path", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityDocument), nil
}

// UpdateDocument updates an existing document by its ID.
// Maps to PUT /api/v1/repo/{repo}/id/{id}
// Returns the updated entityDocument or error.
func (r *repository) UpdateDocument(ctx context.Context, documentId string, document entityDocument, options *nuxeoRequestOptions) (*entityDocument, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId)
	res, err := r.client.NewRequest(ctx, options).SetBody(document).SetResult(&entityDocument{}).SetError(&nuxeoError{}).Put(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to update document", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityDocument), nil
}

// DeleteDocument deletes a document by its ID.
// Maps to DELETE /api/v1/repo/{repo}/id/{id}
// Returns error if deletion fails.
func (r *repository) DeleteDocument(ctx context.Context, documentId string) error {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId)
	res, err := r.client.NewRequest(ctx, nil).SetError(&nuxeoError{}).Delete(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to delete document", slog.String("error", err.Error()))
		return err
	}
	return nil
}

///////////////
//// QUERY ////
///////////////

// Query executes a NXQL query against the repository.
// Maps to GET /api/v1/query
// Accepts NXQL query string, query parameters, pagination options, and request options.
// Returns entityDocuments (list of documents) or error.
func (r *repository) Query(ctx context.Context, query string, queryParams []string, paginationOptions *SortedPaginationOptions, options *nuxeoRequestOptions) (*entityDocuments, error) {
	path := internal.PathApiV1 + "/query"

	params := url.Values{}
	params.Add("query", query)
	for _, qp := range queryParams {
		params.Add("queryParams", qp)
	}
	params = internal.MergeUrlValues(params, paginationOptions.QueryParams())
	path += "?" + params.Encode()

	res, err := r.client.NewRequest(ctx, options).SetResult(&entityDocuments{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch documents", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityDocuments), nil
}

// QueryByProvider executes a named query provider against the repository.
// Maps to GET /api/v1/query/{providerName}
// Accepts provider name, query parameters, named query parameters, pagination options, and request options.
// Returns entityDocuments or error.
func (r *repository) QueryByProvider(ctx context.Context, providerName string, queryParams []string, namedQueryParams map[string]string, paginationOptions *SortedPaginationOptions, options *nuxeoRequestOptions) (*entityDocuments, error) {
	path := internal.PathApiV1 + "/query/" + url.PathEscape(providerName)

	params := url.Values{}
	for k, v := range namedQueryParams {
		params.Add(k, v)
	}
	for _, qp := range queryParams {
		params.Add("queryParams", qp)
	}
	params = internal.MergeUrlValues(params, paginationOptions.QueryParams())

	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	res, err := r.client.NewRequest(ctx, options).SetResult(&entityDocuments{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch documents by provider", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityDocuments), nil
}

///////////////
//// AUDIT ////
///////////////

// FetchAuditByPath retrieves audit logs for a document by its repository path.
// Maps to GET /api/v1/repo/{repo}/path/{path}/@audit
// Returns entityAudit or error.
func (r *repository) FetchAuditByPath(ctx context.Context, documentPath string, options *nuxeoRequestOptions) (*entityAudit, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path" + documentPath + "/@audit"
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityAudit{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch audit by path", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityAudit), nil
}

// FetchAuditById retrieves audit logs for a document by its ID.
// Maps to GET /api/v1/repo/{repo}/id/{id}/@audit
// Returns entityAudit or error.
func (r *repository) FetchAuditById(ctx context.Context, documentId string, options *nuxeoRequestOptions) (*entityAudit, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId) + "/@audit"
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityAudit{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch audit by ID", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityAudit), nil
}

/////////////
//// ACP ////
/////////////

// FetchPermissionsByPath retrieves permissions (ACLs) for a document by its repository path.
// Maps to GET /api/v1/repo/{repo}/path/{path}/@acl
// Returns entityACP or error.
func (r *repository) FetchPermissionsByPath(ctx context.Context, documentPath string, options *nuxeoRequestOptions) (*entityACP, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path" + documentPath + "/@acl"
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityACP{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch permissions by path", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityACP), nil
}

// FetchPermissionsById retrieves permissions (ACLs) for a document by its ID.
// Maps to GET /api/v1/repo/{repo}/id/{id}/@acl
// Returns entityACP or error.
func (r *repository) FetchPermissionsById(ctx context.Context, documentId string, options *nuxeoRequestOptions) (*entityACP, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId) + "/@acl"
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityACP{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch permissions by ID", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityACP), nil
}

//////////////////
//// CHILDREN ////
//////////////////

// FetchChildrenByPath retrieves child documents under a parent specified by repository path.
// Maps to GET /api/v1/repo/{repo}/path/{parentPath}/@children
// Returns entityDocuments or error.
func (r *repository) FetchChildrenByPath(ctx context.Context, parentPath string, options *nuxeoRequestOptions) (*entityDocuments, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path" + parentPath + "/@children"
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityDocuments{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch children by path", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityDocuments), nil
}

// FetchChildrenById retrieves child documents under a parent specified by document ID.
// Maps to GET /api/v1/repo/{repo}/id/{parentId}/@children
// Returns entityDocuments or error.
func (r *repository) FetchChildrenById(ctx context.Context, parentId string, options *nuxeoRequestOptions) (*entityDocuments, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(parentId) + "/@children"
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityDocuments{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch children by ID", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityDocuments), nil
}

///////////////
//// BLOBS ////
///////////////

// StreamBlobByPath streams a blob from a document specified by repository path and blob XPath.
// Maps to GET /api/v1/repo/{repo}/path/{path}/@blob/{xpath}
// Returns Blob (stream, filename, mimetype, length) or error.
func (r *repository) StreamBlobByPath(ctx context.Context, documentPath string, blobXPath string, options *nuxeoRequestOptions) (*blob, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path" + documentPath + "/@blob/" + url.PathEscape(blobXPath)
	res, err := r.client.NewRequest(ctx, options).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to stream blob by path", slog.String("error", err.Error()))
		return nil, err
	}
	return &blob{
		Filename: internal.GetStreamFilenameFrom(res),
		MimeType: internal.GetStreamContentTypeFrom(res),
		Stream:   res.Body,
		Length:   strconv.Itoa(internal.GetStreamContentLengthFrom(res)),
	}, nil
}

// StreamBlobById streams a blob from a document specified by ID and blob XPath.
// Maps to GET /api/v1/repo/{repo}/id/{id}/@blob/{xpath}
// Returns Blob (stream, filename, mimetype, length) or error.
func (r *repository) StreamBlobById(ctx context.Context, documentId string, blobXPath string, options *nuxeoRequestOptions) (*blob, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId) + "/@blob/" + url.PathEscape(blobXPath)
	res, err := r.client.NewRequest(ctx, options).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to stream blob by ID", slog.String("error", err.Error()))
		return nil, err
	}
	return &blob{
		Filename: internal.GetStreamFilenameFrom(res),
		MimeType: internal.GetStreamContentTypeFrom(res),
		Stream:   res.Body,
		Length:   strconv.Itoa(internal.GetStreamContentLengthFrom(res)),
	}, nil
}

///////////////////
//// WORKFLOWS ////
///////////////////

// StartWorkflowInstanceWithDocId starts a workflow instance for a document specified by ID.
// Maps to POST /api/v1/repo/{repo}/id/{id}/@workflow
// Returns the started entityWorkflow or error.
func (r *repository) StartWorkflowInstanceWithDocId(ctx context.Context, documentId string, workflow entityWorkflow, options *nuxeoRequestOptions) (*entityWorkflow, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId) + "/@workflow"
	res, err := r.client.NewRequest(ctx, options).SetBody(workflow).SetResult(&entityWorkflow{}).SetError(&nuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to start workflow instance with document ID", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityWorkflow), nil
}

// StartWorkflowInstanceWithDocPath starts a workflow instance for a document specified by repository path.
// Maps to POST /api/v1/repo/{repo}/path/{path}/@workflow
// Returns the started entityWorkflow or error.
func (r *repository) StartWorkflowInstanceWithDocPath(ctx context.Context, documentPath string, workflow entityWorkflow, options *nuxeoRequestOptions) (*entityWorkflow, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path" + documentPath + "/@workflow"
	res, err := r.client.NewRequest(ctx, options).SetBody(workflow).SetResult(&entityWorkflow{}).SetError(&nuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to start workflow instance with document path", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityWorkflow), nil
}

// FetchWorkflowInstancesByDocId retrieves workflow instances for a document specified by ID.
// Maps to GET /api/v1/repo/{repo}/id/{id}/@workflow
// Returns entityWorkflows (list) or error.
func (r *repository) FetchWorkflowInstancesByDocId(ctx context.Context, documentId string, options *nuxeoRequestOptions) (*entityWorkflows, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId) + "/@workflow"
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityWorkflows{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch workflow instances by document ID", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityWorkflows), nil
}

// FetchWorkflowInstancesByDocPath retrieves workflow instances for a document specified by repository path.
// Maps to GET /api/v1/repo/{repo}/path/{path}/@workflow
// Returns entityWorkflows (list) or error.
func (r *repository) FetchWorkflowInstancesByDocPath(ctx context.Context, documentPath string, options *nuxeoRequestOptions) (*entityWorkflows, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path" + documentPath + "/@workflow"
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityWorkflows{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch workflow instances by document path", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityWorkflows), nil
}

// FetchWorkflowInstance retrieves a workflow instance by its ID.
// Maps to GET /api/v1/workflow/{workflowInstanceId}
// Returns entityWorkflow or error.
func (r *repository) FetchWorkflowInstance(ctx context.Context, workflowInstanceId string, options *nuxeoRequestOptions) (*entityWorkflow, error) {
	path := internal.PathApiV1 + "/workflow/" + url.PathEscape(workflowInstanceId)
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityWorkflow{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch workflow instance", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityWorkflow), nil
}

// CancelWorkflowInstance cancels a workflow instance by its ID.
// Maps to DELETE /api/v1/workflow/{workflowInstanceId}
// Returns error if cancellation fails.
func (r *repository) CancelWorkflowInstance(ctx context.Context, workflowInstanceId string) error {
	path := internal.PathApiV1 + "/workflow/" + url.PathEscape(workflowInstanceId)
	res, err := r.client.NewRequest(ctx, nil).SetError(&nuxeoError{}).Delete(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to cancel workflow instance", slog.String("error", err.Error()))
		return err
	}
	return nil
}

// FetchWorkflowInstanceGraph retrieves the graph of a workflow instance by its ID.
// Maps to GET /api/v1/workflow/{workflowInstanceId}/graph
// Returns entityWorkflowGraph or error.
func (r *repository) FetchWorkflowInstanceGraph(ctx context.Context, workflowInstanceId string, options *nuxeoRequestOptions) (*entityWorkflowGraph, error) {
	path := internal.PathApiV1 + "/workflow/" + url.PathEscape(workflowInstanceId) + "/graph"
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityWorkflowGraph{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch workflow instance graph", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityWorkflowGraph), nil
}

// FetchWorkflowModel retrieves a workflow model by its name.
// Maps to GET /api/v1/workflowModel/{workflowModelName}
// Returns entityWorkflow or error.
func (r *repository) FetchWorkflowModel(ctx context.Context, workflowModelName string, options *nuxeoRequestOptions) (*entityWorkflow, error) {
	path := internal.PathApiV1 + "/workflowModel/" + url.PathEscape(workflowModelName)
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityWorkflow{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch workflow model", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityWorkflow), nil
}

// FetchWorkflowModelGraph retrieves the graph of a workflow model by its name.
// Maps to GET /api/v1/workflowModel/{workflowModelName}/graph
// Returns entityWorkflowGraph or error.
func (r *repository) FetchWorkflowModelGraph(ctx context.Context, workflowModelName string, options *nuxeoRequestOptions) (*entityWorkflowGraph, error) {
	path := internal.PathApiV1 + "/workflowModel/" + url.PathEscape(workflowModelName) + "/graph"
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityWorkflowGraph{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch workflow model graph", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityWorkflowGraph), nil
}

// FetchWorkflowModels retrieves all workflow models available in the repository.
// Maps to GET /api/v1/workflowModel
// Returns entityWorkflows (list) or error.
func (r *repository) FetchWorkflowModels(ctx context.Context, options *nuxeoRequestOptions) (*entityWorkflows, error) {
	path := internal.PathApiV1 + "/workflowModel"
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityWorkflows{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch workflow models", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityWorkflows), nil
}

///////////////////
//// WEB ADAPTER ////
////////////////////

// The following methods allow invoking custom web adapters for documents, supporting extension points for business logic and integrations.
// See Nuxeo Web Adapter documentation for details: https://doc.nuxeo.com/rest-api/web-adapter/
//
// CreateForAdapter invokes a custom web adapter for a document by its ID using POST.
// Maps to POST /api/v1/repo/{repo}/id/{id}/@{adapter}/{pathSuffix}
// See https://doc.nuxeo.com/rest-api/web-adapter/
// Returns the raw HTTP response or error.
func (r *repository) CreateForAdapter(ctx context.Context, documentId string, adapter string, pathSuffix string, queryParams []string, payload any, options *nuxeoRequestOptions) (*http.Response, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId) + "/@" + url.PathEscape(adapter) + "/" + pathSuffix

	params := url.Values{}
	for _, qp := range queryParams {
		params.Add("queryParams", qp)
	}
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	res, err := r.client.NewRequest(ctx, options).SetBody(payload).SetError(&nuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to create for adapter", slog.String("error", err.Error()))
		return nil, err
	}
	return res.RawResponse, nil
}

// FetchForAdapter invokes a custom web adapter for a document by its ID using GET.
// Maps to GET /api/v1/repo/{repo}/id/{id}/@{adapter}/{pathSuffix}
// See https://doc.nuxeo.com/rest-api/web-adapter/
// Returns the raw HTTP response or error.
func (r *repository) FetchForAdapter(ctx context.Context, documentId string, adapter string, pathSuffix string, queryParams []string, options *nuxeoRequestOptions) (*http.Response, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId) + "/@" + url.PathEscape(adapter) + "/" + pathSuffix

	params := url.Values{}
	for _, qp := range queryParams {
		params.Add("queryParams", qp)
	}
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	res, err := r.client.NewRequest(ctx, options).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch for adapter", slog.String("error", err.Error()))
		return nil, err
	}
	return res.RawResponse, nil
}

// UpdateForAdapter invokes a custom web adapter for a document by its ID using PUT.
// Maps to PUT /api/v1/repo/{repo}/id/{id}/@{adapter}/{pathSuffix}
// See https://doc.nuxeo.com/rest-api/web-adapter/
// Returns the raw HTTP response or error.
func (r *repository) UpdateForAdapter(ctx context.Context, documentId string, adapter string, pathSuffix string, queryParams []string, payload any, options *nuxeoRequestOptions) (*http.Response, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId) + "/@" + url.PathEscape(adapter) + "/" + pathSuffix

	params := url.Values{}
	for _, qp := range queryParams {
		params.Add("queryParams", qp)
	}
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	res, err := r.client.NewRequest(ctx, options).SetBody(payload).SetError(&nuxeoError{}).Put(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to update for adapter", slog.String("error", err.Error()))
		return nil, err
	}
	return res.RawResponse, nil
}

// DeleteForAdapter invokes a custom web adapter for a document by its ID using DELETE.
// Maps to DELETE /api/v1/repo/{repo}/id/{id}/@{adapter}/{pathSuffix}
// See https://doc.nuxeo.com/rest-api/web-adapter/
// Returns the raw HTTP response or error.
func (r *repository) DeleteForAdapter(ctx context.Context, documentId string, adapter string, pathSuffix string, queryParams []string) (*http.Response, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId) + "/@" + url.PathEscape(adapter) + "/" + pathSuffix

	params := url.Values{}
	for _, qp := range queryParams {
		params.Add("queryParams", qp)
	}
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	res, err := r.client.NewRequest(ctx, nil).SetError(&nuxeoError{}).Delete(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to delete for adapter", slog.String("error", err.Error()))
		return nil, err
	}
	return res.RawResponse, nil
}
