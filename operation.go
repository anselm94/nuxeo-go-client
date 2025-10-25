package nuxeo

import (
	"fmt"
	"strings"
	"time"
)

///////////////////////////
//// OPERATION PAYLOAD ////
///////////////////////////

type operationPayload struct {
	Input   string            `json:"input,omitempty"`
	Params  map[string]string `json:"params,omitempty"`
	Context map[string]string `json:"context,omitempty"`
}

func NewOperation(operationId string) *operation {
	return &operation{
		operationId:      operationId,
		params:           make(map[string]string),
		context:          make(map[string]string),
		inputDocumentIds: make([]string, 0),
		inputBlobs:       make([]Blob, 0),
	}
}

// operation represents a Nuxeo operation.
type operation struct {
	operationId      string
	inputDocumentIds []string
	inputBlobs       []Blob
	params           map[string]string
	context          map[string]string
	isVoid           bool
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
