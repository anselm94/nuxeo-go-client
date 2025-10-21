package nuxeo

import (
	"fmt"
	"log/slog"
	"strings"
	"time"
)

type OperationPayload struct {
	Input   string            `json:"input,omitempty"`
	Params  map[string]string `json:"params,omitempty"`
	Context map[string]string `json:"context,omitempty"`
}

// Operation represents a Nuxeo operation.
type Operation struct {
	operationId    string
	inputDocuments []string
	inputBlobs     []Blob
	params         map[string]string
	context        map[string]string

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

func (o *Operation) SetBlobInput(blob Blob) *Operation {
	o.inputBlobs = []Blob{
		blob,
	}
	return o
}

func (o *Operation) SetBlobListInput(blobs []Blob) *Operation {
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

func (o *Operation) Payload() *OperationPayload {
	payload := &OperationPayload{
		Params:  o.params,
		Context: o.context,
	}
	// compute Input field based on inputDocuments
	if len(o.inputDocuments) == 1 {
		payload.Input = "doc:" + o.inputDocuments[0]
	} else if len(o.inputDocuments) > 1 {
		payload.Input = "docs:" + strings.Join(o.inputDocuments, ",")
	}
	return payload
}

func (o *Operation) Blobs() []Blob {
	return o.inputBlobs
}
