package nuxeo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

// ExecuteInto executes a Nuxeo Automation operation and decodes the response into the provided output variable.
//
// The operation is executed via the REST endpoint `/site/automation/{operationId}`.
// If the operation returns a JSON response, it is decoded into 'out'.
// For blob responses, use Execute and handle the stream directly.
// See: https://doc.nuxeo.com/rest-api/1/automation-endpoint/#operation-execution-response
func (o *operationManager) ExecuteInto(ctx context.Context, operation operation, out any, requestOptions *nuxeoRequestOptions) error {
	res, err := o.Execute(ctx, operation, requestOptions)
	if err != nil {
		return err
	}
	defer res.Close()

	return json.NewDecoder(res).Decode(out)
}

// Execute runs a Nuxeo Automation operation and returns the raw response stream.
//
// The operation is executed via the REST endpoint `/site/automation/{operationId}`.
// If the operation includes blobs, the request is sent as multipart/related; otherwise, as JSON.
// Caller is responsible for closing the returned stream.
// See: https://doc.nuxeo.com/rest-api/1/automation-endpoint/#executing-operations
func (o *operationManager) Execute(ctx context.Context, operation operation, requestOptions *nuxeoRequestOptions) (io.ReadCloser, error) {
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
func (o *operationManager) executeViaJson(ctx context.Context, operation operation, requestOptions *nuxeoRequestOptions) (io.ReadCloser, error) {
	path := "/site/automation/" + url.PathEscape(operation.operationId)

	request := o.client.NewRequest(ctx, requestOptions).SetError(&nuxeoError{})
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

	return res.Body, err
}

// executeViaMultipart sends the operation request as multipart/related to the Automation endpoint.
//
// Used when the operation includes blob input. The JSON payload is sent as the first part, followed by each blob.
// See: https://doc.nuxeo.com/rest-api/1/automation-endpoint/#taking-a-blob-as-input
func (o *operationManager) executeViaMultipart(ctx context.Context, operation operation, requestOptions *nuxeoRequestOptions) (io.ReadCloser, error) {
	path := "/site/automation/" + url.PathEscape(operation.operationId)

	request := o.client.NewRequest(ctx, requestOptions).SetError(&nuxeoError{})
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
	case 0:
		return nil, fmt.Errorf("no blobs to send in multipart request")
	case 1:
		// single blob input
		blob := operation.blobs()[0]
		request.SetMultipartField("input", blob.Filename, blob.MimeType, blob.Stream)
	default:
		// multiple blob inputs
		for i, blob := range operation.blobs() {
			fieldName := fmt.Sprintf("input-%d", i+1)
			request.SetMultipartField(fieldName, blob.Filename, blob.MimeType, blob.Stream)
		}
	}

	res, err := request.Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		o.logger.Error("Failed to execute operation with blobs", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to execute operation with blobs: %d %w", res.StatusCode(), err)
	} else if res.StatusCode() == 204 {
		return nil, nil
	}

	return res.Body, nil
}

// FetchOperation retrieves the description of an Automation operation from the server.
//
// This performs a GET request to `/site/automation/{operationId}` and returns the operation metadata.
// See: https://doc.nuxeo.com/rest-api/1/automation-endpoint/#getting-the-automation-service
func (o *operationManager) FetchOperation(ctx context.Context, operationId string) (*operationPayload, error) {
	return nil, nil
}
