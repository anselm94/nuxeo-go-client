package nuxeo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"
)

// Operation represents a Nuxeo operation.
type Operation struct {
	Input   string            `json:"input,omitempty"`
	Params  map[string]string `json:"params,omitempty"`
	Context map[string]string `json:"context,omitempty"`

	// internal

	automationId   string
	inputDocuments []string
	inputBlobs     []io.Reader

	request *NuxeoRequest
	logger  slog.Logger
}

// SetInput sets the input for the operation.
func (o *Operation) SetDocumentInput(docIdOrPath string) *Operation {
	o.inputDocuments = []string{
		docIdOrPath,
	}
	return o
}

func (o *Operation) SetDocumentListInput(docIdsOrPaths []string) *Operation {
	o.inputDocuments = docIdsOrPaths
	return o
}

func (o *Operation) SetBlobInput(blob io.Reader) *Operation {
	o.inputBlobs = []io.Reader{
		blob,
	}
	return o
}

func (o *Operation) SetBlobListInput(blobs []io.Reader) *Operation {
	o.inputBlobs = blobs
	return o
}

func (o *Operation) SetContext(key string, value string) *Operation {
	o.Context[key] = value
	return o
}

// SetParam sets a parameter for the operation.
func (o *Operation) SetParam(key string, value any) *Operation {
	switch v := value.(type) {
	case string:
		o.Params[key] = v
	case int, int32, int64:
		o.Params[key] = fmt.Sprintf("%d", v)
	case float32, float64:
		o.Params[key] = fmt.Sprintf("%f", v)
	case time.Time:
		o.Params[key] = v.Format(time.RFC3339)
	case bool:
		o.Params[key] = fmt.Sprintf("%t", v)
	default:
		o.Params[key] = fmt.Sprintf("%v", v)
	}
	return o
}

func (o *Operation) SetParams(params map[string]any) *Operation {
	for key, val := range params {
		o.SetParam(key, val)
	}
	return o
}

// ExecuteInto executes the operation and decodes the response into out.
func (o *Operation) ExecuteInto(out any) error {
	res, err := o.Execute()
	if err != nil {
		return err
	}
	defer res.Close()

	return json.NewDecoder(res).Decode(out)
}

// Execute runs the operation using the client.
func (o *Operation) Execute() (io.ReadCloser, error) {
	o.request.SetDoNotParseResponse(true)

	// decide execution method based on presence of blobs
	if len(o.inputBlobs) > 0 {
		return o.executeViaMultipart()
	} else {
		return o.executeViaJson()
	}
}

func (o *Operation) executeViaJson() (io.ReadCloser, error) {
	// compute Input field based on inputDocuments
	if len(o.inputDocuments) == 1 {
		o.Input = "doc:" + o.inputDocuments[0]
	} else if len(o.inputDocuments) > 1 {
		o.Input = "docs:" + strings.Join(o.inputDocuments, ",")
	}

	o.request.SetBody(o)

	res, err := o.request.Post("/site/automation/" + o.automationId)
	if err != nil || res.StatusCode() != 200 {
		o.logger.Error("Failed to execute operation", "error", err, "status", res.StatusCode())
		return nil, fmt.Errorf("failed to execute operation: %d %w", res.StatusCode(), err)
	} else if res.StatusCode() == 204 {
		return nil, nil
	}

	return res.Body, err
}

func (o *Operation) executeViaMultipart() (io.ReadCloser, error) {
	o.request.SetContentType("multipart/related")

	// add json payload as `application/json+nxrequest` part
	payloadBytes, _ := json.Marshal(o)
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
	} else if res.StatusCode() == 204 {
		return nil, nil
	}

	return res.Body, nil
}
