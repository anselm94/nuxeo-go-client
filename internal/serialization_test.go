package internal

import (
	"testing"
)

func TestMarshalUnmarshal(t *testing.T) {
	data := map[string]any{"foo": "bar"}
	jsonBytes, err := Marshal(data)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	var out map[string]any
	err = Unmarshal(jsonBytes, &out)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if out["foo"] != "bar" {
		t.Errorf("Unmarshal got %v, want %v", out["foo"], "bar")
	}
}
