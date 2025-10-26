package nuxeo

import "encoding/json"

// EntityExtendedGroup represents an extended group for a Nuxeo user.
// It contains the group's name, label, and URL.
type entityExtendedGroup struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	Url   string `json:"url"`
}

// EntityUser represents a Nuxeo user entity.
// It includes user properties, administrator status, and extended group memberships.
type entityUser struct {
	entity
	Id              string                `json:"id"`
	IsAdministrator bool                  `json:"isAdministrator"`
	IsAnonymous     bool                  `json:"isAnonymous"`
	Properties      map[string]Field      `json:"properties"`
	ExtendedGroups  []entityExtendedGroup `json:"extendedGroups,omitempty"`
}

// NewUser creates a new EntityUser with the given username.
// The username is set as both the Id and the "username" property.
func NewUser(username string) *entityUser {
	properties := make(map[string]Field)
	// Set username property
	if fieldUsername, err := NewField(username); err == nil {
		properties[UserPropertyUsername] = fieldUsername
	}
	return &entityUser{
		entity: entity{
			EntityType: EntityTypeUser,
		},
		Id:         username,
		Properties: properties,
	}
}

// IdOrUsername returns the user's Id if set, otherwise the username property.
func (u *entityUser) IdOrUsername() string {
	if u.Id != "" {
		return u.Id
	}
	return u.Username()
}

// Username returns the user's username property.
func (u *entityUser) Username() string {
	if val, err := u.Properties[UserPropertyUsername].String(); err == nil && val != nil {
		return *val
	}
	return ""
}

// Password returns the user's password property.
func (u *entityUser) Password() string {
	if val, err := u.Properties[UserPropertyPassword].String(); err == nil && val != nil {
		return *val
	}
	return ""
}

// FirstName returns the user's first name property.
func (u *entityUser) FirstName() string {
	if val, err := u.Properties[UserPropertyFirstName].String(); err == nil && val != nil {
		return *val
	}
	return ""
}

// LastName returns the user's last name property.
func (u *entityUser) LastName() string {
	if val, err := u.Properties[UserPropertyLastName].String(); err == nil && val != nil {
		return *val
	}
	return ""
}

// Email returns the user's email property.
func (u *entityUser) Email() string {
	if val, err := u.Properties[UserPropertyEmail].String(); err == nil && val != nil {
		return *val
	}
	return ""
}

// Groups returns the list of group names the user belongs to.
func (u *entityUser) Groups() []string {
	if val, err := u.Properties[UserPropertyGroups].StringList(); err == nil {
		return val
	}
	return nil
}

// Company returns the user's company property.
func (u *entityUser) Company() string {
	if val, err := u.Properties[UserPropertyCompany].String(); err == nil && val != nil {
		return *val
	}
	return ""
}

// TenantId returns the user's tenant ID property.
func (u *entityUser) TenantId() string {
	if val, err := u.Properties[UserPropertyTenantId].String(); err == nil && val != nil {
		return *val
	}
	return ""
}

// Property returns the Field value for the given property key.
func (u *entityUser) Property(key string) Field {
	return u.Properties[key]
}

// SetProperty sets the value for the given property key.
func (u *entityUser) SetProperty(key string, value any) {
	if fieldValue, err := json.Marshal(value); err == nil {
		u.Properties[key] = fieldValue
	}
}

type entityUsers paginableEntities[entityUser]
