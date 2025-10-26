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

type operationManager struct {
	client *NuxeoClient
	logger *slog.Logger
}

// ExecuteInto executes the operation and decodes the response into out.
func (o *operationManager) ExecuteInto(ctx context.Context, operation operation, requestOptions *nuxeoRequestOptions, out any) error {
	res, err := o.Execute(ctx, operation, requestOptions)
	if err != nil {
		return err
	}
	defer res.Close()

	return json.NewDecoder(res).Decode(out)
}

// Execute runs the operation using the client.
func (o *operationManager) Execute(ctx context.Context, operation operation, requestOptions *nuxeoRequestOptions) (io.ReadCloser, error) {
	// decide execution method based on presence of blobs
	if len(operation.blobs()) > 0 {
		return o.executeViaMultipart(ctx, operation, requestOptions)
	} else {
		return o.executeViaJson(ctx, operation, requestOptions)
	}
}

func (o *operationManager) executeViaJson(ctx context.Context, operation operation, requestOptions *nuxeoRequestOptions) (io.ReadCloser, error) {
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

func (o *operationManager) executeViaMultipart(ctx context.Context, operation operation, requestOptions *nuxeoRequestOptions) (io.ReadCloser, error) {
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

func (o *operationManager) FetchOperation(ctx context.Context, operationId string) (*operationPayload, error) {
	return nil, nil
}
