package nuxeo

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/anselm94/nuxeo/internal"
)

// Repository represents a Nuxeo repository.
type Repository struct {
	name string

	// internal

	client *NuxeoClient
	logger *slog.Logger
}

func (r *Repository) Name() string {
	return r.name
}

///////////////////
//// DOCUMENTS ////
///////////////////

func (r *Repository) FetchDocumentRoot(ctx context.Context, options *nuxeoRequestOptions) (*Document, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path/"
	res, err := r.client.NewRequest(ctx, options).SetResult(&Document{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch document root", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Document), nil
}

func (r *Repository) FetchDocumentById(ctx context.Context, documentID string, options *nuxeoRequestOptions) (*Document, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentID)
	res, err := r.client.NewRequest(ctx, options).SetResult(&Document{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch document by ID", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Document), nil
}

func (r *Repository) FetchDocumentByPath(ctx context.Context, documentPath string, options *nuxeoRequestOptions) (*Document, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path" + url.PathEscape(documentPath)
	res, err := r.client.NewRequest(ctx, options).SetResult(&Document{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch document by path", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Document), nil
}

func (r *Repository) CreateDocumentById(ctx context.Context, parentId string, document Document, options *nuxeoRequestOptions) (*Document, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(parentId)
	res, err := r.client.NewRequest(ctx, options).SetBody(document).SetResult(&Document{}).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to create document by ID", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Document), nil
}

func (r *Repository) CreateDocumentByPath(ctx context.Context, parentPath string, document Document, options *nuxeoRequestOptions) (*Document, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path" + url.PathEscape(parentPath)
	res, err := r.client.NewRequest(ctx, options).SetBody(document).SetResult(&Document{}).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to create document by path", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Document), nil
}

func (r *Repository) UpdateDocument(ctx context.Context, documentId string, document Document, options *nuxeoRequestOptions) (*Document, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId)
	res, err := r.client.NewRequest(ctx, options).SetBody(document).SetResult(&Document{}).SetError(&NuxeoError{}).Put(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to update document", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Document), nil
}

func (r *Repository) DeleteDocument(ctx context.Context, documentId string) error {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId)
	res, err := r.client.NewRequest(ctx, nil).SetError(&NuxeoError{}).Delete(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to delete document", slog.String("error", err.Error()))
		return err
	}
	return nil
}

///////////////
//// QUERY ////
///////////////

func (r *Repository) Query(ctx context.Context, query string, queryParams []string, paginationOptions *SortedPaginationOptions, options *nuxeoRequestOptions) (*Documents, error) {
	path := internal.PathApiV1 + "/query"

	params := url.Values{}
	params.Add("query", query)
	for _, qp := range queryParams {
		params.Add("queryParams", qp)
	}
	params = internal.MergeUrlValues(params, paginationOptions.QueryParams())
	path += "?" + params.Encode()

	res, err := r.client.NewRequest(ctx, options).SetResult(&Documents{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch documents", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Documents), nil
}

func (r *Repository) QueryByProvider(ctx context.Context, providerName string, queryParams []string, namedQueryParams map[string]string, paginationOptions *SortedPaginationOptions, options *nuxeoRequestOptions) (*Documents, error) {
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

	res, err := r.client.NewRequest(ctx, options).SetResult(&Documents{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch documents by provider", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Documents), nil
}

///////////////
//// AUDIT ////
///////////////

func (r *Repository) FetchAuditByPath(ctx context.Context, documentPath string, options *nuxeoRequestOptions) (*Audit, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path" + documentPath + "/@audit"
	res, err := r.client.NewRequest(ctx, options).SetResult(&Audit{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch audit by path", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Audit), nil
}

func (r *Repository) FetchAuditById(ctx context.Context, documentId string, options *nuxeoRequestOptions) (*Audit, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId) + "/@audit"
	res, err := r.client.NewRequest(ctx, options).SetResult(&Audit{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch audit by ID", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Audit), nil
}

/////////////
//// ACP ////
/////////////

func (r *Repository) FetchPermissionsByPath(ctx context.Context, documentPath string, options *nuxeoRequestOptions) (*ACP, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path" + documentPath + "/@acl"
	res, err := r.client.NewRequest(ctx, options).SetResult(&ACP{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch permissions by path", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*ACP), nil
}

func (r *Repository) FetchPermissionsById(ctx context.Context, documentId string, options *nuxeoRequestOptions) (*ACP, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId) + "/@acl"
	res, err := r.client.NewRequest(ctx, options).SetResult(&ACP{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch permissions by ID", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*ACP), nil
}

//////////////////
//// CHILDREN ////
//////////////////

func (r *Repository) FetchChildrenByPath(ctx context.Context, parentPath string, options *nuxeoRequestOptions) (*Documents, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path" + parentPath + "/@children"
	res, err := r.client.NewRequest(ctx, options).SetResult(&Documents{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch children by path", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Documents), nil
}

func (r *Repository) FetchChildrenById(ctx context.Context, parentId string, options *nuxeoRequestOptions) (*Documents, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(parentId) + "/@children"
	res, err := r.client.NewRequest(ctx, options).SetResult(&Documents{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch children by ID", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Documents), nil
}

///////////////
//// BLOBS ////
///////////////

func (r *Repository) StreamBlobByPath(ctx context.Context, documentPath string, blobXPath string, options *nuxeoRequestOptions) (*Blob, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path" + documentPath + "/@blob/" + url.PathEscape(blobXPath)
	res, err := r.client.NewRequest(ctx, options).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to stream blob by path", slog.String("error", err.Error()))
		return nil, err
	}
	return &Blob{
		Filename: internal.GetStreamFilenameFrom(res),
		MimeType: internal.GetStreamContentTypeFrom(res),
		Stream:   res.Body,
		Length:   internal.GetStreamContentLengthFrom(res),
	}, nil
}

func (r *Repository) StreamBlobById(ctx context.Context, documentId string, blobXPath string, options *nuxeoRequestOptions) (*Blob, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId) + "/@blob/" + url.PathEscape(blobXPath)
	res, err := r.client.NewRequest(ctx, options).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to stream blob by ID", slog.String("error", err.Error()))
		return nil, err
	}
	return &Blob{
		Filename: internal.GetStreamFilenameFrom(res),
		MimeType: internal.GetStreamContentTypeFrom(res),
		Stream:   res.Body,
		Length:   internal.GetStreamContentLengthFrom(res),
	}, nil
}

///////////////////
//// WORKFLOWS ////
///////////////////

func (r *Repository) StartWorkflowInstanceWithDocId(ctx context.Context, documentId string, workflow Workflow, options *nuxeoRequestOptions) (*Workflow, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId) + "/@workflow"
	res, err := r.client.NewRequest(ctx, options).SetBody(workflow).SetResult(&Workflow{}).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to start workflow instance with document ID", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Workflow), nil
}

func (r *Repository) StartWorkflowInstanceWithDocPath(ctx context.Context, documentPath string, workflow Workflow, options *nuxeoRequestOptions) (*Workflow, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path" + documentPath + "/@workflow"
	res, err := r.client.NewRequest(ctx, options).SetBody(workflow).SetResult(&Workflow{}).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to start workflow instance with document path", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Workflow), nil
}

func (r *Repository) FetchWorkflowInstancesByDocId(ctx context.Context, documentId string, options *nuxeoRequestOptions) (*Workflows, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId) + "/@workflow"
	res, err := r.client.NewRequest(ctx, options).SetResult(&Workflows{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch workflow instances by document ID", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Workflows), nil
}

func (r *Repository) FetchWorkflowInstancesByDocPath(ctx context.Context, documentPath string, options *nuxeoRequestOptions) (*Workflows, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path" + documentPath + "/@workflow"
	res, err := r.client.NewRequest(ctx, options).SetResult(&Workflows{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch workflow instances by document path", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Workflows), nil
}

func (r *Repository) FetchWorkflowInstance(ctx context.Context, workflowInstanceId string, options *nuxeoRequestOptions) (*Workflow, error) {
	path := internal.PathApiV1 + "/workflow/" + url.PathEscape(workflowInstanceId)
	res, err := r.client.NewRequest(ctx, options).SetResult(&Workflow{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch workflow instance", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Workflow), nil
}

func (r *Repository) CancelWorkflowInstance(ctx context.Context, workflowInstanceId string) error {
	path := internal.PathApiV1 + "/workflow/" + url.PathEscape(workflowInstanceId)
	res, err := r.client.NewRequest(ctx, nil).SetError(&NuxeoError{}).Delete(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to cancel workflow instance", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (r *Repository) FetchWorkflowInstanceGraph(ctx context.Context, workflowInstanceId string, options *nuxeoRequestOptions) (*WorkflowGraph, error) {
	path := internal.PathApiV1 + "/workflow/" + url.PathEscape(workflowInstanceId) + "/graph"
	res, err := r.client.NewRequest(ctx, options).SetResult(&WorkflowGraph{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch workflow instance graph", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*WorkflowGraph), nil
}

func (r *Repository) FetchWorkflowModel(ctx context.Context, workflowModelName string, options *nuxeoRequestOptions) (*Workflow, error) {
	path := internal.PathApiV1 + "/workflowModel/" + url.PathEscape(workflowModelName)
	res, err := r.client.NewRequest(ctx, options).SetResult(&Workflow{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch workflow model", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Workflow), nil
}

func (r *Repository) FetchWorkflowModelGraph(ctx context.Context, workflowModelName string, options *nuxeoRequestOptions) (*WorkflowGraph, error) {
	path := internal.PathApiV1 + "/workflowModel/" + url.PathEscape(workflowModelName) + "/graph"
	res, err := r.client.NewRequest(ctx, options).SetResult(&WorkflowGraph{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch workflow model graph", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*WorkflowGraph), nil
}

func (r *Repository) FetchWorkflowModels(ctx context.Context, options *nuxeoRequestOptions) (*Workflows, error) {
	path := internal.PathApiV1 + "/workflowModel"
	res, err := r.client.NewRequest(ctx, options).SetResult(&Workflows{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch workflow models", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Workflows), nil
}

/////////////////////
//// WEB ADAPTER ////
/////////////////////

func (r *Repository) CreateForAdapter(ctx context.Context, documentId string, adapter string, pathSuffix string, queryParams []string, payload any, options *nuxeoRequestOptions) (*http.Response, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId) + "/@" + url.PathEscape(adapter) + "/" + pathSuffix
	params := url.Values{}
	for _, qp := range queryParams {
		params.Add("queryParams", qp)
	}
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var result any
	res, err := r.client.NewRequest(ctx, options).SetBody(payload).SetResult(&result).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to create for adapter", slog.String("error", err.Error()))
		return nil, err
	}
	return res.RawResponse, nil
}

func (r *Repository) FetchForAdapter(ctx context.Context, documentId string, adapter string, pathSuffix string, queryParams []string, options *nuxeoRequestOptions) (*http.Response, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId) + "/@" + url.PathEscape(adapter) + "/" + pathSuffix
	params := url.Values{}
	for _, qp := range queryParams {
		params.Add("queryParams", qp)
	}
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var result any
	res, err := r.client.NewRequest(ctx, options).SetResult(&result).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch for adapter", slog.String("error", err.Error()))
		return nil, err
	}
	return res.RawResponse, nil
}

func (r *Repository) UpdateForAdapter(ctx context.Context, documentId string, adapter string, pathSuffix string, queryParams []string, payload any, options *nuxeoRequestOptions) (*http.Response, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId) + "/@" + url.PathEscape(adapter) + "/" + pathSuffix
	params := url.Values{}
	for _, qp := range queryParams {
		params.Add("queryParams", qp)
	}
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var result any
	res, err := r.client.NewRequest(ctx, options).SetBody(payload).SetResult(&result).SetError(&NuxeoError{}).Put(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to update for adapter", slog.String("error", err.Error()))
		return nil, err
	}
	return res.RawResponse, nil
}

func (r *Repository) DeleteForAdapter(ctx context.Context, documentId string, adapter string, pathSuffix string, queryParams []string) (*http.Response, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId) + "/@" + url.PathEscape(adapter) + "/" + pathSuffix
	params := url.Values{}
	for _, qp := range queryParams {
		params.Add("queryParams", qp)
	}
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	res, err := r.client.NewRequest(ctx, nil).SetError(&NuxeoError{}).Delete(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to delete for adapter", slog.String("error", err.Error()))
		return nil, err
	}
	return res.RawResponse, nil
}
