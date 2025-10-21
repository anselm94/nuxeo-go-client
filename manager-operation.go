package nuxeo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
)

type OperationManager struct {
	client *NuxeoClient
	logger *slog.Logger
}

func (om *OperationManager) NewOperation(operationId string) *Operation {
	return &Operation{
		operationId:    operationId,
		params:         make(map[string]string),
		context:        make(map[string]string),
		inputDocuments: make([]string, 0),
		inputBlobs:     make([]Blob, 0),
		logger:         *om.logger,
	}
}

// ExecuteInto executes the operation and decodes the response into out.
func (o *OperationManager) ExecuteInto(ctx context.Context, operation *Operation, out any) error {
	res, err := o.Execute(ctx, operation)
	if err != nil {
		return err
	}
	defer res.Close()

	return json.NewDecoder(res).Decode(out)
}

// Execute runs the operation using the client.
func (o *OperationManager) Execute(ctx context.Context, operation *Operation) (io.ReadCloser, error) {
	// decide execution method based on presence of blobs
	if len(operation.Blobs()) > 0 {
		return o.executeViaMultipart(ctx, operation)
	} else {
		return o.executeViaJson(ctx, operation)
	}
}

func (o *OperationManager) executeViaJson(ctx context.Context, operation *Operation) (io.ReadCloser, error) {
	request := o.client.NewRequest(ctx)
	request.SetDoNotParseResponse(true)

	request.SetBody(operation.Payload())

	res, err := request.Post("/site/automation/" + operation.operationId)
	if err != nil || res.StatusCode() != 200 {
		o.logger.Error("Failed to execute operation", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to execute operation: %d %w", res.StatusCode(), err)
	} else if res.StatusCode() == 204 {
		return nil, nil
	}

	return res.Body, err
}

func (o *OperationManager) executeViaMultipart(ctx context.Context, operation *Operation) (io.ReadCloser, error) {
	request := o.client.NewRequest(ctx)
	request.SetDoNotParseResponse(true)

	request.SetContentType("multipart/related")

	// add json payload as `application/json+nxrequest` part
	payloadBytes, _ := json.Marshal(operation.Payload())
	request.SetMultipartField("root", "", "application/json+nxrequest", bytes.NewReader(payloadBytes))

	// add input documents one by one
	for i, blob := range operation.Blobs() {
		fieldName := fmt.Sprintf("input-%d", i+1)
		request.SetMultipartField(fieldName, fmt.Sprintf("blob%d", i+1), "application/octet-stream", blob.Stream)
	}

	res, err := request.Post("/site/automation/" + operation.operationId)
	if err != nil || res.StatusCode() != 200 {
		o.logger.Error("Failed to execute operation with blobs", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to execute operation with blobs: %d %w", res.StatusCode(), err)
	} else if res.StatusCode() == 204 {
		return nil, nil
	}

	return res.Body, nil
}

func (o *OperationManager) FetchOperation(ctx context.Context, operationId string) (*OperationPayload, error) {
	return nil, nil
}
