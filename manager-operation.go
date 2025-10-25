package nuxeo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/anselm94/nuxeo/internal"
)

type OperationManager struct {
	client *NuxeoClient
	logger *slog.Logger
}

func (om *OperationManager) NewOperation(operationId string) *operation {
	return &operation{
		operationId:      operationId,
		params:           make(map[string]string),
		context:          make(map[string]string),
		inputDocumentIds: make([]string, 0),
		inputBlobs:       make([]Blob, 0),
		logger:           *om.logger,
	}
}

// ExecuteInto executes the operation and decodes the response into out.
func (o *OperationManager) ExecuteInto(ctx context.Context, operation operation, requestOptions *nuxeoRequestOptions, out any) error {
	res, err := o.Execute(ctx, operation, requestOptions)
	if err != nil {
		return err
	}
	defer res.Close()

	return json.NewDecoder(res).Decode(out)
}

// Execute runs the operation using the client.
func (o *OperationManager) Execute(ctx context.Context, operation operation, requestOptions *nuxeoRequestOptions) (io.ReadCloser, error) {
	// decide execution method based on presence of blobs
	if len(operation.blobs()) > 0 {
		return o.executeViaMultipart(ctx, operation, requestOptions)
	} else {
		return o.executeViaJson(ctx, operation, requestOptions)
	}
}

func (o *OperationManager) executeViaJson(ctx context.Context, operation operation, requestOptions *nuxeoRequestOptions) (io.ReadCloser, error) {
	request := o.client.NewRequest(ctx, requestOptions)
	request.SetDoNotParseResponse(true)

	if operation.isVoid {
		request.SetHeader(internal.HeaderXVoidOperation, "true")
	}

	request.SetBody(operation.payload())

	path := "/site/automation/" + url.PathEscape(operation.operationId)
	res, err := request.Post(path)
	if err != nil || res.StatusCode() != 200 {
		o.logger.Error("Failed to execute operation", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to execute operation: %d %w", res.StatusCode(), err)
	} else if res.StatusCode() == 204 {
		return nil, nil
	}

	return res.Body, err
}

func (o *OperationManager) executeViaMultipart(ctx context.Context, operation operation, requestOptions *nuxeoRequestOptions) (io.ReadCloser, error) {
	request := o.client.NewRequest(ctx, requestOptions)
	request.SetDoNotParseResponse(true)

	if operation.isVoid {
		request.SetHeader(internal.HeaderXVoidOperation, "true")
	}

	request.SetContentType("multipart/related")

	// add json payload as `application/json+nxrequest` part
	payloadBytes, _ := json.Marshal(operation.payload())
	request.SetMultipartField("root", "", "application/json+nxrequest", bytes.NewReader(payloadBytes))

	// add input documents one by one
	for i, blob := range operation.blobs() {
		fieldName := fmt.Sprintf("input-%d", i+1)
		request.SetMultipartField(fieldName, fmt.Sprintf("blob%d", i+1), "application/octet-stream", blob.Stream)
	}

	path := "/site/automation/" + url.PathEscape(operation.operationId)
	res, err := request.Post(path)
	if err != nil || res.StatusCode() != 200 {
		o.logger.Error("Failed to execute operation with blobs", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to execute operation with blobs: %d %w", res.StatusCode(), err)
	} else if res.StatusCode() == 204 {
		return nil, nil
	}

	return res.Body, nil
}

func (o *OperationManager) FetchOperation(ctx context.Context, operationId string) (*operationPayload, error) {
	return nil, nil
}

///////////////////////////
//// OPERATION PAYLOAD ////
///////////////////////////

type operationPayload struct {
	Input   string            `json:"input,omitempty"`
	Params  map[string]string `json:"params,omitempty"`
	Context map[string]string `json:"context,omitempty"`
}

// operation represents a Nuxeo operation.
type operation struct {
	operationId      string
	inputDocumentIds []string
	inputBlobs       []Blob
	params           map[string]string
	context          map[string]string
	isVoid           bool

	request *nuxeoRequest
	logger  slog.Logger
}

// SetInput sets the input for the operation.
func (o *operation) SetInputDocumentId(docIdOrPath string) *operation {
	o.inputDocumentIds = []string{
		docIdOrPath,
	}
	return o
}

func (o *operation) SetInputDocumentIds(docIdsOrPaths []string) *operation {
	o.inputDocumentIds = docIdsOrPaths
	return o
}

func (o *operation) SetInputBlob(blob Blob) *operation {
	o.inputBlobs = []Blob{
		blob,
	}
	return o
}

func (o *operation) SetInputBlobs(blobs []Blob) *operation {
	o.inputBlobs = blobs
	return o
}

func (o *operation) SetContext(key string, value string) *operation {
	o.context[key] = value
	return o
}

// SetParam sets a parameter for the operation.
func (o *operation) SetParam(key string, value any) *operation {
	switch v := value.(type) {
	case string:
		o.params[key] = v
	case int, int32, int64:
		o.params[key] = fmt.Sprintf("%d", v)
	case float32, float64:
		o.params[key] = fmt.Sprintf("%f", v)
	case time.Time:
		o.params[key] = v.Format(ISO8601TimeLayout)
	case bool:
		o.params[key] = fmt.Sprintf("%t", v)
	default:
		o.params[key] = fmt.Sprintf("%v", v)
	}
	return o
}

func (o *operation) SetParams(params map[string]any) *operation {
	for key, val := range params {
		o.SetParam(key, val)
	}
	return o
}

func (o *operation) SetVoidOperation(isVoid bool) *operation {
	o.isVoid = isVoid
	return o
}

func (o *operation) payload() *operationPayload {
	payload := &operationPayload{
		Params:  o.params,
		Context: o.context,
	}
	// compute Input field based on inputDocuments
	if len(o.inputDocumentIds) == 1 {
		payload.Input = "doc:" + o.inputDocumentIds[0]
	} else if len(o.inputDocumentIds) > 1 {
		payload.Input = "docs:" + strings.Join(o.inputDocumentIds, ",")
	}
	return payload
}

func (o *operation) blobs() []Blob {
	return o.inputBlobs
}
