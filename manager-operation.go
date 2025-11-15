package nuxeo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/anselm94/nuxeo-go-client/internal"
)

// operationManager provides methods to execute Nuxeo Automation operations via the REST API.
//
// See: https://doc.nuxeo.com/rest-api/1/automation-endpoint/
type operationManager struct {
	client *NuxeoClient
	logger *slog.Logger
}

// Execute runs a Nuxeo Automation operation and returns the operation response.
//
// The operation is executed via the REST endpoint `/site/automation/{operationId}`.
// See: https://doc.nuxeo.com/rest-api/1/automation-endpoint/#executing-operations
func (o *operationManager) Execute(ctx context.Context, operation operation, requestOptions *nuxeoRequestOptions) (*operationResponse, error) {
	// decide execution method based on presence of blobs
	if len(operation.blobs()) > 0 {
		return o.executeViaMultipart(ctx, operation, requestOptions)
	} else {
		return o.executeViaJson(ctx, operation, requestOptions)
	}
}

// executeViaJson sends the operation request as JSON to the Automation endpoint.
//
// Used when the operation does not include blob input. Handles void operations and error responses.
// See: https://doc.nuxeo.com/rest-api/1/automation-endpoint/#executing-operations
func (o *operationManager) executeViaJson(ctx context.Context, operation operation, requestOptions *nuxeoRequestOptions) (*operationResponse, error) {
	path := "/site/automation/" + url.PathEscape(operation.operationId)

	request := o.client.NewRequest(ctx, requestOptions).SetError(&NuxeoError{})
	request.SetDoNotParseResponse(true)
	request.SetBody(operation.payload())

	if operation.isVoid {
		request.SetHeader(internal.HeaderXVoidOperation, "true")
	}

	res, err := request.Post(path)
	if err := handleNuxeoError(err, res); err != nil {
		o.logger.Error("Failed to execute operation", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to execute operation: %d %w", res.StatusCode(), err)
	} else if res.StatusCode() == 204 {
		return nil, nil
	}

	return newOperationResponse(res), nil
}

// executeViaMultipart sends the operation request as multipart/related to the Automation endpoint.
//
// Used when the operation includes blob input. The JSON payload is sent as the first part, followed by each blob.
// See: https://doc.nuxeo.com/rest-api/1/automation-endpoint/#taking-a-blob-as-input
func (o *operationManager) executeViaMultipart(ctx context.Context, operation operation, requestOptions *nuxeoRequestOptions) (*operationResponse, error) {
	path := "/site/automation/" + url.PathEscape(operation.operationId)

	request := o.client.NewRequest(ctx, requestOptions).SetError(&NuxeoError{})
	request.SetDoNotParseResponse(true)
	request.SetHeader("Accept", "application/json, */*")

	if operation.isVoid {
		request.SetHeader(internal.HeaderXVoidOperation, "true")
	}

	// request.SetContentType("multipart/related")

	// add json payload as `application/json` part
	payloadBytes, _ := json.Marshal(operation.payload())
	request.SetMultipartField("request", "", "application/json", bytes.NewReader(payloadBytes))

	// add blobs as subsequent parts
	switch len(operation.blobs()) {
	case 1:
		// single blob input
		blob := operation.blobs()[0]
		request.SetMultipartField("input", blob.Filename, blob.MimeType, blob)
	default:
		// multiple blob inputs
		for i, blob := range operation.blobs() {
			fieldName := fmt.Sprintf("input-%d", i+1)
			request.SetMultipartField(fieldName, blob.Filename, blob.MimeType, blob)
		}
	}

	res, err := request.Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		o.logger.Error("Failed to execute operation with blobs", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to execute operation with blobs: %d %w", res.StatusCode(), err)
	}
	return newOperationResponse(res), nil
}
