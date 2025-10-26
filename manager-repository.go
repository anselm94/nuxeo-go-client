package nuxeo

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/anselm94/nuxeo-go-client/internal"
)

// repository represents a Nuxeo repository.
type repository struct {
	name string

	// internal

	client *NuxeoClient
	logger *slog.Logger
}

func (r *repository) Name() string {
	return r.name
}

///////////////////
//// DOCUMENTS ////
///////////////////

func (r *repository) FetchDocumentRoot(ctx context.Context, options *nuxeoRequestOptions) (*entityDocument, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path/"
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityDocument{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch document root", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityDocument), nil
}

func (r *repository) FetchDocumentById(ctx context.Context, documentID string, options *nuxeoRequestOptions) (*entityDocument, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentID)
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityDocument{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch document by ID", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityDocument), nil
}

func (r *repository) FetchDocumentByPath(ctx context.Context, documentPath string, options *nuxeoRequestOptions) (*entityDocument, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path" + url.PathEscape(documentPath)
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityDocument{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch document by path", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityDocument), nil
}

func (r *repository) CreateDocumentById(ctx context.Context, parentId string, doc entityDocument, options *nuxeoRequestOptions) (*entityDocument, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(parentId)
	res, err := r.client.NewRequest(ctx, options).SetBody(doc).SetResult(&entityDocument{}).SetError(&nuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to create document by ID", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityDocument), nil
}

func (r *repository) CreateDocumentByPath(ctx context.Context, parentPath string, document entityDocument, options *nuxeoRequestOptions) (*entityDocument, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path" + url.PathEscape(parentPath)
	res, err := r.client.NewRequest(ctx, options).SetBody(document).SetResult(&entityDocument{}).SetError(&nuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to create document by path", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityDocument), nil
}

func (r *repository) UpdateDocument(ctx context.Context, documentId string, document entityDocument, options *nuxeoRequestOptions) (*entityDocument, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId)
	res, err := r.client.NewRequest(ctx, options).SetBody(document).SetResult(&entityDocument{}).SetError(&nuxeoError{}).Put(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to update document", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityDocument), nil
}

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

func (r *repository) FetchAuditByPath(ctx context.Context, documentPath string, options *nuxeoRequestOptions) (*entityAudit, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path" + documentPath + "/@audit"
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityAudit{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch audit by path", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityAudit), nil
}

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

func (r *repository) FetchPermissionsByPath(ctx context.Context, documentPath string, options *nuxeoRequestOptions) (*entityACP, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path" + documentPath + "/@acl"
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityACP{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch permissions by path", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityACP), nil
}

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

func (r *repository) FetchChildrenByPath(ctx context.Context, parentPath string, options *nuxeoRequestOptions) (*entityDocuments, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path" + parentPath + "/@children"
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityDocuments{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch children by path", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityDocuments), nil
}

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

func (r *repository) StreamBlobByPath(ctx context.Context, documentPath string, blobXPath string, options *nuxeoRequestOptions) (*Blob, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path" + documentPath + "/@blob/" + url.PathEscape(blobXPath)
	res, err := r.client.NewRequest(ctx, options).SetError(&nuxeoError{}).Get(path)

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

func (r *repository) StreamBlobById(ctx context.Context, documentId string, blobXPath string, options *nuxeoRequestOptions) (*Blob, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId) + "/@blob/" + url.PathEscape(blobXPath)
	res, err := r.client.NewRequest(ctx, options).SetError(&nuxeoError{}).Get(path)

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

func (r *repository) StartWorkflowInstanceWithDocId(ctx context.Context, documentId string, workflow entityWorkflow, options *nuxeoRequestOptions) (*entityWorkflow, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId) + "/@workflow"
	res, err := r.client.NewRequest(ctx, options).SetBody(workflow).SetResult(&entityWorkflow{}).SetError(&nuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to start workflow instance with document ID", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityWorkflow), nil
}

func (r *repository) StartWorkflowInstanceWithDocPath(ctx context.Context, documentPath string, workflow entityWorkflow, options *nuxeoRequestOptions) (*entityWorkflow, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path" + documentPath + "/@workflow"
	res, err := r.client.NewRequest(ctx, options).SetBody(workflow).SetResult(&entityWorkflow{}).SetError(&nuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to start workflow instance with document path", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityWorkflow), nil
}

func (r *repository) FetchWorkflowInstancesByDocId(ctx context.Context, documentId string, options *nuxeoRequestOptions) (*entityWorkflows, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId) + "/@workflow"
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityWorkflows{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch workflow instances by document ID", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityWorkflows), nil
}

func (r *repository) FetchWorkflowInstancesByDocPath(ctx context.Context, documentPath string, options *nuxeoRequestOptions) (*entityWorkflows, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/path" + documentPath + "/@workflow"
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityWorkflows{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch workflow instances by document path", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityWorkflows), nil
}

func (r *repository) FetchWorkflowInstance(ctx context.Context, workflowInstanceId string, options *nuxeoRequestOptions) (*entityWorkflow, error) {
	path := internal.PathApiV1 + "/workflow/" + url.PathEscape(workflowInstanceId)
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityWorkflow{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch workflow instance", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityWorkflow), nil
}

func (r *repository) CancelWorkflowInstance(ctx context.Context, workflowInstanceId string) error {
	path := internal.PathApiV1 + "/workflow/" + url.PathEscape(workflowInstanceId)
	res, err := r.client.NewRequest(ctx, nil).SetError(&nuxeoError{}).Delete(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to cancel workflow instance", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (r *repository) FetchWorkflowInstanceGraph(ctx context.Context, workflowInstanceId string, options *nuxeoRequestOptions) (*entityWorkflowGraph, error) {
	path := internal.PathApiV1 + "/workflow/" + url.PathEscape(workflowInstanceId) + "/graph"
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityWorkflowGraph{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch workflow instance graph", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityWorkflowGraph), nil
}

func (r *repository) FetchWorkflowModel(ctx context.Context, workflowModelName string, options *nuxeoRequestOptions) (*entityWorkflow, error) {
	path := internal.PathApiV1 + "/workflowModel/" + url.PathEscape(workflowModelName)
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityWorkflow{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch workflow model", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityWorkflow), nil
}

func (r *repository) FetchWorkflowModelGraph(ctx context.Context, workflowModelName string, options *nuxeoRequestOptions) (*entityWorkflowGraph, error) {
	path := internal.PathApiV1 + "/workflowModel/" + url.PathEscape(workflowModelName) + "/graph"
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityWorkflowGraph{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch workflow model graph", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityWorkflowGraph), nil
}

func (r *repository) FetchWorkflowModels(ctx context.Context, options *nuxeoRequestOptions) (*entityWorkflows, error) {
	path := internal.PathApiV1 + "/workflowModel"
	res, err := r.client.NewRequest(ctx, options).SetResult(&entityWorkflows{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch workflow models", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityWorkflows), nil
}

/////////////////////
//// WEB ADAPTER ////
/////////////////////

func (r *repository) CreateForAdapter(ctx context.Context, documentId string, adapter string, pathSuffix string, queryParams []string, payload any, options *nuxeoRequestOptions) (*http.Response, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId) + "/@" + url.PathEscape(adapter) + "/" + pathSuffix
	params := url.Values{}
	for _, qp := range queryParams {
		params.Add("queryParams", qp)
	}
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var result any
	res, err := r.client.NewRequest(ctx, options).SetBody(payload).SetResult(&result).SetError(&nuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to create for adapter", slog.String("error", err.Error()))
		return nil, err
	}
	return res.RawResponse, nil
}

func (r *repository) FetchForAdapter(ctx context.Context, documentId string, adapter string, pathSuffix string, queryParams []string, options *nuxeoRequestOptions) (*http.Response, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId) + "/@" + url.PathEscape(adapter) + "/" + pathSuffix
	params := url.Values{}
	for _, qp := range queryParams {
		params.Add("queryParams", qp)
	}
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var result any
	res, err := r.client.NewRequest(ctx, options).SetResult(&result).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to fetch for adapter", slog.String("error", err.Error()))
		return nil, err
	}
	return res.RawResponse, nil
}

func (r *repository) UpdateForAdapter(ctx context.Context, documentId string, adapter string, pathSuffix string, queryParams []string, payload any, options *nuxeoRequestOptions) (*http.Response, error) {
	path := internal.PathApiV1 + "/repo/" + url.PathEscape(r.name) + "/id/" + url.PathEscape(documentId) + "/@" + url.PathEscape(adapter) + "/" + pathSuffix
	params := url.Values{}
	for _, qp := range queryParams {
		params.Add("queryParams", qp)
	}
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var result any
	res, err := r.client.NewRequest(ctx, options).SetBody(payload).SetResult(&result).SetError(&nuxeoError{}).Put(path)

	if err := handleNuxeoError(err, res); err != nil {
		r.logger.Error("Failed to update for adapter", slog.String("error", err.Error()))
		return nil, err
	}
	return res.RawResponse, nil
}

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
