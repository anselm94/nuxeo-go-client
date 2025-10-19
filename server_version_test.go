package nuxeo

import "testing"

func TestParseServerVersion(t *testing.T) {
	v := ParseServerVersion("10.2.3")
	if v.Major != 10 {
		t.Errorf("Major got %d, want %d", v.Major, 10)
	}
	if v.Minor != 2 {
		t.Errorf("Minor got %d, want %d", v.Minor, 2)
	}
	if v.Patch != 3 {
		t.Errorf("Patch got %d, want %d", v.Patch, 3)
	}
}
