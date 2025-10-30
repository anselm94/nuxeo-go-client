package nuxeo

import (
	"reflect"
	"testing"
)

func TestNewUser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		username string
	}{
		{"normal username", "jdoe"},
		{"empty username", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewUser(tc.username)
			if u.Id != tc.username {
				t.Errorf("Id mismatch: got %q, want %q", u.Id, tc.username)
			}
			if u.entity.EntityType != EntityTypeUser {
				t.Errorf("EntityType mismatch: got %q, want %q", u.entity.EntityType, EntityTypeUser)
			}
			if val, ok := u.Properties[UserPropertyUsername]; !ok {
				t.Error("username property missing")
			} else {
				str, err := val.String()
				if err != nil || str == nil || *str != tc.username {
					t.Errorf("username property value: got %v, want %q", str, tc.username)
				}
			}
		})
	}
}

func TestUserGetters(t *testing.T) {
	t.Parallel()
	u := NewUser("jdoe")
	// Set all properties
	props := map[string]Field{
		UserPropertyPassword:  NewStringField("secret"),
		UserPropertyFirstName: NewStringField("John"),
		UserPropertyLastName:  NewStringField("Doe"),
		UserPropertyEmail:     NewStringField("jdoe@example.com"),
		UserPropertyGroups:    NewStringListField([]string{"admins", "users"}),
		UserPropertyCompany:   NewStringField("Acme"),
		UserPropertyTenantId:  NewStringField("tenant1"),
	}
	for k, v := range props {
		u.SetProperty(k, v)
	}
	if got := u.IdOrUsername(); got != "jdoe" {
		t.Errorf("IdOrUsername: got %q, want 'jdoe'", got)
	}
	if got := u.Username(); got != "jdoe" {
		t.Errorf("Username: got %q, want 'jdoe'", got)
	}
	if got := u.Password(); got != "secret" {
		t.Errorf("Password: got %q, want 'secret'", got)
	}
	if got := u.FirstName(); got != "John" {
		t.Errorf("FirstName: got %q, want 'John'", got)
	}
	if got := u.LastName(); got != "Doe" {
		t.Errorf("LastName: got %q, want 'Doe'", got)
	}
	if got := u.Email(); got != "jdoe@example.com" {
		t.Errorf("Email: got %q, want 'jdoe@example.com'", got)
	}
	if got := u.Groups(); !reflect.DeepEqual(got, []string{"admins", "users"}) {
		t.Errorf("Groups: got %v, want [admins users]", got)
	}
	if got := u.Company(); got != "Acme" {
		t.Errorf("Company: got %q, want 'Acme'", got)
	}
	if got := u.TenantId(); got != "tenant1" {
		t.Errorf("TenantId: got %q, want 'tenant1'", got)
	}
}

func TestUserGettersMissingProperties(t *testing.T) {
	t.Parallel()
	u := NewUser("jdoe")
	// Remove all properties except username
	for k := range u.Properties {
		if k != UserPropertyUsername {
			delete(u.Properties, k)
		}
	}
	if got := u.Password(); got != "" {
		t.Errorf("Password missing: got %q, want empty", got)
	}
	if got := u.FirstName(); got != "" {
		t.Errorf("FirstName missing: got %q, want empty", got)
	}
	if got := u.LastName(); got != "" {
		t.Errorf("LastName missing: got %q, want empty", got)
	}
	if got := u.Email(); got != "" {
		t.Errorf("Email missing: got %q, want empty", got)
	}
	if got := u.Groups(); got != nil {
		t.Errorf("Groups missing: got %v, want nil", got)
	}
	if got := u.Company(); got != "" {
		t.Errorf("Company missing: got %q, want empty", got)
	}
	if got := u.TenantId(); got != "" {
		t.Errorf("TenantId missing: got %q, want empty", got)
	}
}

func TestUserPropertyAndSetProperty(t *testing.T) {
	t.Parallel()
	u := NewUser("jdoe")
	// Set custom property
	u.SetProperty("custom", NewStringField("customval"))
	f, ok := u.Property("custom")
	if !ok {
		t.Error("custom property not found")
	}
	str, err := f.String()
	if err != nil || str == nil || *str != "customval" {
		t.Errorf("custom property value: got %v, want 'customval'", str)
	}
}

func TestUserExtendedGroups(t *testing.T) {
	t.Parallel()
	groups := []ExtendedGroup{
		{Name: "admins", Label: "Admins", Url: "/group/admins"},
		{Name: "users", Label: "Users", Url: "/group/users"},
	}
	u := NewUser("jdoe")
	u.ExtendedGroups = groups
	if !reflect.DeepEqual(u.ExtendedGroups, groups) {
		t.Errorf("ExtendedGroups mismatch: got %v, want %v", u.ExtendedGroups, groups)
	}
}
