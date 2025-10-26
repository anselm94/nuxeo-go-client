package nuxeo

import (
	"fmt"
	"strings"
	"time"
)

///////////////////////////
//// OPERATION PAYLOAD ////
///////////////////////////

// See Nuxeo Automation REST API: https://doc.nuxeo.com/rest-api/1/automation-endpoint/
// This file provides types and helpers for building and executing Automation operations via the REST API.

// operationPayload defines the JSON body for Nuxeo Automation requests.
//
// See: https://doc.nuxeo.com/rest-api/1/automation-endpoint/#executing-operations
//
// Fields:
//
//	Input   - Reference to the input object (e.g. "doc:/path", "docs:...", or omitted for void)
//	Params  - Operation parameters (all as strings)
//	Context - Context variables for the operation
//
// For blob input, use multipart/related requests with this payload as the first part.
type operationPayload struct {
	Input   string            `json:"input,omitempty"`
	Params  map[string]string `json:"params,omitempty"`
	Context map[string]string `json:"context,omitempty"`
}

// NewOperation creates a new Nuxeo Automation operation request.
//
// The operationId should match the ID of the operation or chain as listed in the Automation service registry.
// See: https://doc.nuxeo.com/rest-api/1/automation-endpoint/#getting-the-automation-service
func NewOperation(operationId string) *operation {
	return &operation{
		operationId:      operationId,
		params:           make(map[string]string),
		context:          make(map[string]string),
		inputDocumentIds: make([]string, 0),
		inputBlobs:       make([]Blob, 0),
	}
}

// Operation represents a Nuxeo Automation operation request.
//
// Operations are executed via the REST endpoint `/site/automation/{operationId}`.
// See: https://doc.nuxeo.com/rest-api/1/automation-endpoint/
//
// The operation can take as input a document, documents, blob(s), or void, and parameters/context as strings.
// For blob input, use multipart/related requests; for documents, use JSON.
// Use SetVoidOperation(true) to avoid downloading blob responses (see X-NXVoidOperation header).
//
// Example usage:
//
//	op := NewOperation("Document.Fetch").SetInputDocumentId("/default-domain/workspaces/myws/file")
//	op.SetParam("xpath", "file:content")
//	op.SetContext("myVar", "value")
type operation struct {
	operationId      string
	inputDocumentIds []string
	inputBlobs       []Blob
	params           map[string]string
	context          map[string]string
	isVoid           bool
}

// SetInput sets the input for the operation.
// SetInputDocumentId sets the input for the operation to a single document.
//
// The input can be a document UID or absolute path. The value is encoded as "doc:{idOrPath}" in the request payload.
// See: https://doc.nuxeo.com/rest-api/1/automation-endpoint/#request-input
func (o *operation) SetInputDocumentId(docIdOrPath string) *operation {
	o.inputDocumentIds = []string{
		docIdOrPath,
	}
	return o
}

// SetInputDocumentIds sets the input for the operation to a list of documents.
//
// Each entry can be a document UID or absolute path. The value is encoded as "docs:{idOrPath1},{idOrPath2},..." in the request payload.
// See: https://doc.nuxeo.com/rest-api/1/automation-endpoint/#request-input
func (o *operation) SetInputDocumentIds(docIdsOrPaths []string) *operation {
	o.inputDocumentIds = docIdsOrPaths
	return o
}

// SetInputBlob sets the input for the operation to a single blob.
//
// For blob input, the request will be sent as multipart/related with the JSON payload as the first part and the blob as the second part.
// See: https://doc.nuxeo.com/rest-api/1/automation-endpoint/#request-input
func (o *operation) SetInputBlob(blob Blob) *operation {
	o.inputBlobs = []Blob{
		blob,
	}
	return o
}

// SetInputBlobs sets the input for the operation to a list of blobs.
//
// For blob list input, the request will be sent as multipart/related with the JSON payload as the first part and each blob as a subsequent part.
// See: https://doc.nuxeo.com/rest-api/1/automation-endpoint/#request-input
func (o *operation) SetInputBlobs(blobs []Blob) *operation {
	o.inputBlobs = blobs
	return o
}

// SetContext sets a context variable for the operation request.
//
// Context variables are available to the operation or chain during execution.
// See: https://doc.nuxeo.com/rest-api/1/automation-endpoint/#executing-operations
func (o *operation) SetContext(key string, value string) *operation {
	o.context[key] = value
	return o
}

// SetParam sets a parameter for the operation request.
//
// All parameters are encoded as strings in the payload, but can represent numbers, dates, booleans, EL expressions, etc.
// See: https://doc.nuxeo.com/rest-api/1/automation-endpoint/#request-parameter-types
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

// SetParams sets multiple parameters for the operation request.
//
// All parameters are encoded as strings in the payload, but can represent numbers, dates, booleans, EL expressions, etc.
// See: https://doc.nuxeo.com/rest-api/1/automation-endpoint/#request-parameter-types
func (o *operation) SetParams(params map[string]any) *operation {
	for key, val := range params {
		o.SetParam(key, val)
	}
	return o
}

// SetVoidOperation marks the operation as void, indicating no output is expected.
//
// This sets the X-NXVoidOperation header, which avoids downloading blob responses.
// See: https://doc.nuxeo.com/rest-api/1/automation-endpoint/#operation-execution-response
func (o *operation) SetVoidOperation(isVoid bool) *operation {
	o.isVoid = isVoid
	return o
}

// payload builds the operationPayload for the request body.
//
// The Input field is computed based on the input documents set; for blobs, use multipart requests.
// See: https://doc.nuxeo.com/rest-api/1/automation-endpoint/#executing-operations
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

// blobs returns the input blobs for the operation, if any.
//
// Used to determine if the request should be sent as multipart/related.
func (o *operation) blobs() []Blob {
	return o.inputBlobs
}
