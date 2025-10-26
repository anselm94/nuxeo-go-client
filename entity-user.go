package nuxeo

import "encoding/json"

type entityExtendedGroup struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	Url   string `json:"url"`
}

// entityUser represents a Nuxeo user.
type entityUser struct {
	entity
	Id              string                `json:"id"`
	IsAdministrator bool                  `json:"isAdministrator"`
	IsAnonymous     bool                  `json:"isAnonymous"`
	Properties      map[string]Field      `json:"properties"`
	ExtendedGroups  []entityExtendedGroup `json:"extendedGroups,omitempty"`
}

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

func (u *entityUser) IdOrUsername() string {
	if u.Id != "" {
		return u.Id
	}
	return u.Username()
}

func (u *entityUser) Username() string {
	if val, err := u.Properties[UserPropertyUsername].String(); err == nil && val != nil {
		return *val
	}
	return ""
}

func (u *entityUser) Password() string {
	if val, err := u.Properties[UserPropertyPassword].String(); err == nil && val != nil {
		return *val
	}
	return ""
}

func (u *entityUser) FirstName() string {
	if val, err := u.Properties[UserPropertyFirstName].String(); err == nil && val != nil {
		return *val
	}
	return ""
}

func (u *entityUser) LastName() string {
	if val, err := u.Properties[UserPropertyLastName].String(); err == nil && val != nil {
		return *val
	}
	return ""
}

func (u *entityUser) Email() string {
	if val, err := u.Properties[UserPropertyEmail].String(); err == nil && val != nil {
		return *val
	}
	return ""
}

func (u *entityUser) Groups() []string {
	if val, err := u.Properties[UserPropertyGroups].StringList(); err == nil {
		return val
	}
	return nil
}

func (u *entityUser) Company() string {
	if val, err := u.Properties[UserPropertyCompany].String(); err == nil && val != nil {
		return *val
	}
	return ""
}

func (u *entityUser) TenantId() string {
	if val, err := u.Properties[UserPropertyTenantId].String(); err == nil && val != nil {
		return *val
	}
	return ""
}

func (u *entityUser) Property(key string) Field {
	return u.Properties[key]
}

func (u *entityUser) SetProperty(key string, value any) {
	if fieldValue, err := json.Marshal(value); err == nil {
		u.Properties[key] = fieldValue
	}
}

type entityUsers paginableEntities[entityUser]
