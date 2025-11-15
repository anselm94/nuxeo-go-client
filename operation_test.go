package nuxeo

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewOperation(t *testing.T) {
	op := NewOperation("Document.Fetch")
	if op == nil {
		t.Fatal("NewOperation returned nil")
	}
	if op.operationId != "Document.Fetch" {
		t.Errorf("operationId = %q, want %q", op.operationId, "Document.Fetch")
	}
	if len(op.inputDocumentIds) != 0 {
		t.Errorf("inputDocumentIds not empty")
	}
	if len(op.inputBlobs) != 0 {
		t.Errorf("inputBlobs not empty")
	}
	if len(op.context) != 0 {
		t.Errorf("context not empty")
	}
	if op.isVoid {
		t.Errorf("isVoid should be false by default")
	}
}

func TestSetInputDocumentId(t *testing.T) {
	op := NewOperation("op")
	op.SetInputDocumentId("/path/to/doc")
	if !reflect.DeepEqual(op.inputDocumentIds, []string{"/path/to/doc"}) {
		t.Errorf("inputDocumentIds = %v, want %v", op.inputDocumentIds, []string{"/path/to/doc"})
	}
}

func TestSetInputDocumentIds(t *testing.T) {
	op := NewOperation("op")
	ids := []string{"doc1", "doc2"}
	op.SetInputDocumentIds(ids...)
	if !reflect.DeepEqual(op.inputDocumentIds, ids) {
		t.Errorf("inputDocumentIds = %v, want %v", op.inputDocumentIds, ids)
	}
}

func TestSetInputBlob(t *testing.T) {
	op := NewOperation("op")
	b := blob{Filename: "f.txt", MimeType: "text/plain", Length: "123"}
	op.SetInputBlob(b)
	if len(op.inputBlobs) != 1 || op.inputBlobs[0].Filename != "f.txt" {
		t.Errorf("inputBlobs = %v, want blob with Filename 'f.txt'", op.inputBlobs)
	}
}

func TestSetInputBlobs(t *testing.T) {
	op := NewOperation("op")
	blobs := []blob{{Filename: "a"}, {Filename: "b"}}
	op.SetInputBlobs(blobs...)
	if !reflect.DeepEqual(op.inputBlobs, blobs) {
		t.Errorf("inputBlobs = %v, want %v", op.inputBlobs, blobs)
	}
}

func TestSetContext(t *testing.T) {
	op := NewOperation("op")
	op.SetContext("foo", "bar")
	if v, ok := op.context["foo"]; !ok || v != "bar" {
		t.Errorf("context[foo] = %v, want 'bar'", v)
	}
}

func TestSetParam_AllTypes(t *testing.T) {
	op := NewOperation("op")
	op.SetParam("str", "hello")
	op.SetParam("int", 42)
	op.SetParam("float", 3.14)
	op.SetParam("bool", true)
	op.SetParam("other", []string{"x", "y"})
	if op.params["str"] != "hello" {
		t.Errorf("params[str] = %v, want 'hello'", op.params["str"])
	}
	if op.params["int"] != "42" {
		t.Errorf("params[int] = %v, want '42'", op.params["int"])
	}
	if op.params["float"] == "" || op.params["float"][:3] != "3.1" {
		t.Errorf("params[float] = %v, want prefix '3.14'", op.params["float"])
	}
	if op.params["bool"] != "true" {
		t.Errorf("params[bool] = %v, want 'true'", op.params["bool"])
	}
	if op.params["other"] != "[x y]" {
		t.Errorf("params[other] = %v, want '[x y]'", op.params["other"])
	}
}

func TestSetParams_MixedTypes(t *testing.T) {
	op := NewOperation("op")
	params := map[string]any{
		"a": "A",
		"b": 1,
		"c": false,
	}
	op.SetParams(params)
	if op.params["a"] != "A" || op.params["b"] != "1" || op.params["c"] != "false" {
		t.Errorf("params = %v, want map with string/int/bool values", op.params)
	}
}

func TestSetVoidOperation(t *testing.T) {
	op := NewOperation("op")
	op.SetVoidOperation(true)
	if !op.isVoid {
		t.Errorf("isVoid = %v, want true", op.isVoid)
	}
	op.SetVoidOperation(false)
	if op.isVoid {
		t.Errorf("isVoid = %v, want false", op.isVoid)
	}
}

func TestPayload_InputVariants(t *testing.T) {
	cases := []struct {
		name      string
		inputDocs []string
		wantInput string
	}{
		{"none", nil, ""},
		{"single", []string{"doc1"}, "doc:doc1"},
		{"multi", []string{"doc1", "doc2"}, "docs:doc1,doc2"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			op := NewOperation("op")
			op.inputDocumentIds = tc.inputDocs
			pl := op.payload()
			if pl.Input != tc.wantInput {
				t.Errorf("payload.Input = %q, want %q", pl.Input, tc.wantInput)
			}
		})
	}
}

func TestPayload_ParamsAndContext(t *testing.T) {
	op := NewOperation("op")
	op.SetParam("foo", "bar")
	op.SetContext("baz", "qux")
	pl := op.payload()
	if pl.Params["foo"] != "bar" {
		t.Errorf("payload.Params[foo] = %v, want 'bar'", pl.Params["foo"])
	}
	if pl.Context["baz"] != "qux" {
		t.Errorf("payload.Context[baz] = %v, want 'qux'", pl.Context["baz"])
	}
}

func TestBlobs(t *testing.T) {
	op := NewOperation("op")
	b1 := blob{Filename: "a"}
	b2 := blob{Filename: "b"}
	op.SetInputBlobs(b1, b2)
	blobs := op.blobs()
	if !reflect.DeepEqual(blobs, []blob{b1, b2}) {
		t.Errorf("blobs() = %v, want %v", blobs, []blob{b1, b2})
	}
}

func TestEdgeCases_EmptyParamsContext(t *testing.T) {
	op := NewOperation("op")
	pl := op.payload()
	if len(pl.Params) != 0 {
		t.Errorf("payload.Params not empty: %v", pl.Params)
	}
	if len(pl.Context) != 0 {
		t.Errorf("payload.Context not empty: %v", pl.Context)
	}
}

func TestOperationResponse_IsVoid(t *testing.T) {
	resp := mockRestyResponse(204, http.Header{"Content-Type": []string{"application/json"}}, nil)
	opr := newOperationResponse(resp)
	if !opr.IsVoid() {
		t.Errorf("IsVoid() = false, want true for 204 response")
	}
	resp2 := mockRestyResponse(200, http.Header{"Content-Type": []string{"application/json"}}, testMarshalBody(t, []byte(`{"foo":"bar"}`)))
	opr2 := newOperationResponse(resp2)
	if opr2.IsVoid() {
		t.Errorf("IsVoid() = true, want false for 200 response")
	}
}

func TestOperationResponse_As(t *testing.T) {
	body := map[string]string{"foo": "bar"}
	resp := mockRestyResponse(200, http.Header{"Content-Type": []string{"application/json"}}, testMarshalBody(t, body))
	opr := newOperationResponse(resp)
	var v map[string]string
	err := opr.As(&v)
	if err != nil {
		t.Fatalf("As() error = %v", err)
	}
	if v["foo"] != "bar" {
		t.Errorf("As() = %v, want foo=bar", v)
	}
}

func TestOperationResponse_AsDocument(t *testing.T) {
	body := map[string]string{"entity-type": "document", "uid": "123"}
	resp := mockRestyResponse(200, http.Header{"Content-Type": []string{"application/json"}}, testMarshalBody(t, body))
	opr := newOperationResponse(resp)
	doc, err := opr.AsDocument()
	if err != nil {
		t.Fatalf("AsDocument() error = %v", err)
	}
	if doc == nil || doc.ID != "123" {
		t.Errorf("AsDocument() = %v, want UID=123", doc)
	}
	// Not JSON
	resp2 := mockRestyResponse(200, http.Header{"Content-Type": []string{"text/plain"}}, testMarshalBody(t, "not json"))
	opr2 := newOperationResponse(resp2)
	_, err2 := opr2.AsDocument()
	if err2 == nil {
		t.Errorf("AsDocument() error = nil, want error for non-JSON")
	}
}

func TestOperationResponse_AsDocumentList(t *testing.T) {
	body := Documents{
		Entries: []Document{
			{ID: "1"},
			{ID: "2"},
		},
	}
	resp := mockRestyResponse(200, http.Header{"Content-Type": []string{"application/json"}}, testMarshalBody(t, body))
	opr := newOperationResponse(resp)
	docs, err := opr.AsDocumentList()
	if err != nil {
		t.Fatalf("AsDocumentList() error = %v", err)
	}
	if len(docs.Entries) != 2 || docs.Entries[0].ID != "1" || docs.Entries[1].ID != "2" {
		t.Errorf("AsDocumentList() = %v, want 2 docs with IDs 1,2", docs)
	}
	// Not JSON
	resp2 := mockRestyResponse(200, http.Header{"Content-Type": []string{"text/plain"}}, testMarshalBody(t, "not json"))
	opr2 := newOperationResponse(resp2)
	_, err2 := opr2.AsDocumentList()
	if err2 == nil {
		t.Errorf("AsDocumentList() error = nil, want error for non-JSON")
	}
}

func TestOperationResponse_AsBlob(t *testing.T) {
	body := []byte("blobdata")
	resp := mockRestyResponse(200, http.Header{"Content-Type": []string{"application/octet-stream"}, "Content-Disposition": []string{`attachment; filename="file.bin"`}}, testMarshalBody(t, body))
	opr := newOperationResponse(resp)
	bl, err := opr.AsBlob()
	if err != nil {
		t.Fatalf("AsBlob() error = %v", err)
	}
	if bl.Filename == "" {
		t.Errorf("AsBlob() Filename empty")
	}
	if bl.MimeType != "application/octet-stream" {
		t.Errorf("AsBlob() MimeType = %v, want application/octet-stream", bl.MimeType)
	}
	if bl.ReadCloser == nil {
		t.Errorf("AsBlob() ReadCloser is nil")
	}
}

func TestOperationResponse_AsBlobList_NotMultipart(t *testing.T) {
	body := []byte("not multipart")
	resp := mockRestyResponse(200, http.Header{"Content-Type": []string{"application/json"}}, testMarshalBody(t, body))
	opr := newOperationResponse(resp)
	seq, err := opr.AsBlobList()
	if err == nil || seq != nil {
		t.Errorf("AsBlobList() = %v, err = %v, want error for non-multipart", seq, err)
	}
}
