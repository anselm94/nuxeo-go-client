package nuxeo

import (
	"testing"
)

func TestNewDirectoryEntry(t *testing.T) {
	entry := NewDirectoryEntry("test-id")
	if entry.ID != "test-id" {
		t.Errorf("ID = %q, want %q", entry.ID, "test-id")
	}
	if entry.Properties == nil {
		t.Errorf("Properties should be initialized, got nil")
	}
}

func TestEntityDirectoryEntry_Id(t *testing.T) {
	cases := []struct {
		name   string
		id     string
		propId any
		expect string
	}{
		{"ID field set", "id123", nil, "id123"},
		{"ID field empty, property set", "", "id456", "id456"},
		{"ID field empty, property not set", "", nil, ""},
	}
	for _, tc := range cases {
		d := NewDirectoryEntry(tc.id)
		if tc.propId != nil {
			d.SetProperty(DirectoryPropertyId, tc.propId)
		}
		if got := d.Id(); got != tc.expect {
			t.Errorf("%s: Id() = %q, want %q", tc.name, got, tc.expect)
		}
	}
}

func TestEntityDirectoryEntry_Label(t *testing.T) {
	d := NewDirectoryEntry("")
	d.SetProperty(DirectoryPropertyLabel, "labelValue")
	if got := d.Label(); got != "labelValue" {
		t.Errorf("Label() = %q, want %q", got, "labelValue")
	}
	// No label property
	d2 := NewDirectoryEntry("")
	if got := d2.Label(); got != "" {
		t.Errorf("Label() = %q, want empty string", got)
	}
}

func TestEntityDirectoryEntry_Ordering(t *testing.T) {
	d := NewDirectoryEntry("")
	d.SetProperty(DirectoryPropertyOrdering, 42.5)
	if got := d.Ordering(); got != 42.5 {
		t.Errorf("Ordering() = %v, want 42.5", got)
	}
	// No ordering property
	d2 := NewDirectoryEntry("")
	if got := d2.Ordering(); got != 0 {
		t.Errorf("Ordering() = %v, want 0", got)
	}
}

func TestEntityDirectoryEntry_Obsolete(t *testing.T) {
	d := NewDirectoryEntry("")
	d.SetProperty(DirectoryPropertyObsolete, 1.0)
	if got := d.Obsolete(); got != 1.0 {
		t.Errorf("Obsolete() = %v, want 1.0", got)
	}
	// No obsolete property
	d2 := NewDirectoryEntry("")
	if got := d2.Obsolete(); got != 0 {
		t.Errorf("Obsolete() = %v, want 0", got)
	}
}

func TestEntityDirectoryEntry_Property(t *testing.T) {
	d := NewDirectoryEntry("")
	d.SetProperty("foo", "bar")
	val, found := d.Property("foo")
	if !found {
		t.Errorf("Property() did not find key 'foo'")
	}
	str, err := val.String()
	if err != nil || str == nil || *str != "bar" {
		t.Errorf("Property() value = %v, want 'bar'", str)
	}
	_, found2 := d.Property("missing")
	if found2 {
		t.Errorf("Property() found unexpected key 'missing'")
	}
}

func TestEntityDirectoryEntry_SetProperty(t *testing.T) {
	d := NewDirectoryEntry("")
	d.SetProperty("foo", "bar")
	val, found := d.Property("foo")
	if !found {
		t.Errorf("SetProperty() did not set key 'foo'")
	}
	str, err := val.String()
	if err != nil || str == nil || *str != "bar" {
		t.Errorf("SetProperty() value = %v, want 'bar'", str)
	}
	// Overwrite property
	d.SetProperty("foo", "baz")
	val2, _ := d.Property("foo")
	str2, _ := val2.String()
	if str2 == nil || *str2 != "baz" {
		t.Errorf("SetProperty() overwrite value = %v, want 'baz'", str2)
	}
}
