package nuxeo

import (
	"reflect"
	"testing"
)

func TestNewGroup(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		groupId string
	}{
		{"normal id", "admins"},
		{"empty id", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGroup(tc.groupId)
			if g.Id != tc.groupId {
				t.Errorf("Id mismatch: got %q, want %q", g.Id, tc.groupId)
			}
			if g.entity.EntityType != EntityTypeGroup {
				t.Errorf("EntityType mismatch: got %q, want %q", g.entity.EntityType, EntityTypeGroup)
			}
		})
	}
}

func TestEntityGroupFields(t *testing.T) {
	t.Parallel()
	props := map[string]Field{}
	f, err := NewField("value")
	if err != nil {
		t.Fatalf("NewField error: %v", err)
	}
	props["custom"] = f
	g := &entityGroup{
		entity:       entity{EntityType: EntityTypeGroup},
		Id:           "testgroup",
		Properties:   props,
		MemberUsers:  []string{"user1", "user2"},
		MemberGroups: []string{"groupA"},
		ParentGroups: []string{"parent1"},
	}
	if g.Id != "testgroup" {
		t.Errorf("Id mismatch: got %q, want %q", g.Id, "testgroup")
	}
	if !reflect.DeepEqual(g.MemberUsers, []string{"user1", "user2"}) {
		t.Errorf("MemberUsers mismatch: got %v", g.MemberUsers)
	}
	if !reflect.DeepEqual(g.MemberGroups, []string{"groupA"}) {
		t.Errorf("MemberGroups mismatch: got %v", g.MemberGroups)
	}
	if !reflect.DeepEqual(g.ParentGroups, []string{"parent1"}) {
		t.Errorf("ParentGroups mismatch: got %v", g.ParentGroups)
	}
	if val, ok := g.Properties["custom"]; !ok {
		t.Error("custom property missing")
	} else {
		str, err := val.String()
		if err != nil || str == nil || *str != "value" {
			t.Errorf("custom property value: got %v, want 'value'", str)
		}
	}
}
