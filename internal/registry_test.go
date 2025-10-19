package internal

import "testing"

func TestRegistry_SetGet(t *testing.T) {
	r := NewRegistry()
	r.Set("foo", 42)
	v := r.Get("foo")
	if v != 42 {
		t.Errorf("Registry.Get got %v, want %v", v, 42)
	}
}
