package nuxeo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"
)

// Operation represents a Nuxeo operation.
type Operation struct {
	request        *NuxeoRequest
	logger         slog.Logger
	automationId   string
	inputDocuments []string
	inputBlobs     []io.Reader
	context        map[string]string
	params         map[string]string
}

type payloadOperation struct {
	Params  map[string]string `json:"params,omitempty"`
	Context map[string]string `json:"context,omitempty"`
	Input   any               `json:"input,omitempty"`
}

// NewOperation creates a new Operation instance.
func (c *NuxeoClient) NewOperation(ctx context.Context, automationId string, options *NuxeoRequestOption) *Operation {
	return &Operation{
		request:      c.NewRequest(ctx).SetNuxeoOption(options),
		logger:       *c.logger,
		automationId: automationId,
		params:       make(map[string]string),
		context:      make(map[string]string),
	}
}

// SetInput sets the input for the operation.
func (o *Operation) SetInputDocument(docIdOrPath string) *Operation {
	o.inputDocuments = []string{
		docIdOrPath,
	}
	return o
}

func (o *Operation) SetInputDocumentList(docIdsOrPaths []string) *Operation {
	o.inputDocuments = docIdsOrPaths
	return o
}

func (o *Operation) SetInputBlob(blob io.Reader) *Operation {
	o.inputBlobs = []io.Reader{
		blob,
	}
	return o
}

func (o *Operation) SetInputBlobList(blobs []io.Reader) *Operation {
	o.inputBlobs = blobs
	return o
}

func (o *Operation) SetContext(key string, value string) *Operation {
	o.context[key] = value
	return o
}

// SetParam sets a parameter for the operation.
func (o *Operation) SetParam(key string, value any) *Operation {
	switch v := value.(type) {
	case string:
		o.params[key] = v
	case int, int32, int64:
		o.params[key] = fmt.Sprintf("%d", v)
	case float32, float64:
		o.params[key] = fmt.Sprintf("%f", v)
	case time.Time:
		o.params[key] = v.Format(time.RFC3339)
	case bool:
		o.params[key] = fmt.Sprintf("%t", v)
	default:
		o.params[key] = fmt.Sprintf("%v", v)
	}
	return o
}

func (o *Operation) SetParams(params map[string]any) *Operation {
	for key, val := range params {
		o.SetParam(key, val)
	}
	return o
}

// Execute runs the operation using the client.
func (o *Operation) Execute() (io.ReadCloser, error) {
	if len(o.inputBlobs) > 0 {
		return o.executeViaMultipart()
	}
	return o.executeViaJson()
}

func (o *Operation) executeViaJson() (io.ReadCloser, error) {
	o.request.SetDoNotParseResponse(true)

	payload := payloadOperation{
		Params:  o.params,
		Context: o.context,
	}
	if len(o.inputDocuments) == 1 {
		payload.Input = "doc:" + o.inputDocuments[0]
	} else if len(o.inputDocuments) > 1 {
		payload.Input = "docs:" + strings.Join(o.inputDocuments, ",")
	}
	o.request.SetBody(payload)

	res, err := o.request.Post("/site/automation/" + o.automationId)
	if err != nil || res.StatusCode() != 200 {
		o.logger.Error("Failed to execute operation", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to execute operation: %d %w", res.StatusCode(), err)
	}

	return res.Body, err
}

func (o *Operation) executeViaMultipart() (io.ReadCloser, error) {
	o.request.SetContentType("multipart/related")
	o.request.SetDoNotParseResponse(true)

	// add json payload as `application/json+nxrequest` part
	payloadBytes, _ := json.Marshal(payloadOperation{
		Params:  o.params,
		Context: o.context,
	})
	o.request.SetMultipartField("root", "", "application/json+nxrequest", bytes.NewReader(payloadBytes))

	// add input documents one by one
	for i, blob := range o.inputBlobs {
		fieldName := fmt.Sprintf("input-%d", i+1)
		o.request.SetMultipartField(fieldName, fmt.Sprintf("blob%d", i+1), "application/octet-stream", blob)
	}

	res, err := o.request.Post("/site/automation/" + o.automationId)
	if err != nil || res.StatusCode() != 200 {
		o.logger.Error("Failed to execute operation with blobs", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to execute operation with blobs: %d %w", res.StatusCode(), err)
	}

	return res.Body, nil
}
