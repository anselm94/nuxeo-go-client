package nuxeo

import (
	"context"
	"testing"
)

func TestNewOperation_ParamsInput(t *testing.T) {
	op := NewOperation("Document.Fetch")
	if op.ID != "Document.Fetch" {
		t.Errorf("ID got %q, want %q", op.ID, "Document.Fetch")
	}
	op.SetParam("key", "value")
	if op.Params["key"] != "value" {
		t.Errorf("SetParam got %v, want %v", op.Params["key"], "value")
	}
	op.SetInput("input-data")
	if op.Input != "input-data" {
		t.Errorf("SetInput got %v, want %v", op.Input, "input-data")
	}
}

func TestOperation_Execute_Stub(t *testing.T) {
	op := NewOperation("Document.Fetch")
	ctx := context.Background()
	client := &NuxeoClient{}
	result, err := op.Execute(ctx, client)
	if result != nil {
		t.Errorf("Execute got %v, want nil (stub)", result)
	}
	if err != nil {
		t.Errorf("Execute err got %v, want nil (stub)", err)
	}
}
