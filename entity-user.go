package nuxeo

// EntityExtendedGroup represents an extended group for a Nuxeo user.
// It contains the group's name, label, and URL.
type ExtendedGroup struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	Url   string `json:"url"`
}

// EntityUser represents a Nuxeo user entity.
// It includes user properties, administrator status, and extended group memberships.
type User struct {
	entity
	Id              string           `json:"id"`
	IsAdministrator bool             `json:"isAdministrator"`
	IsAnonymous     bool             `json:"isAnonymous"`
	Properties      map[string]Field `json:"properties"`
	ExtendedGroups  []ExtendedGroup  `json:"extendedGroups,omitempty"`
}

// NewUser creates a new EntityUser with the given username.
// The username is set as both the Id and the "username" property.
func NewUser(username string) *User {
	return &User{
		entity: entity{
			EntityType: EntityTypeUser,
		},
		Id: username,
		Properties: map[string]Field{
			UserPropertyUsername: NewStringField(username),
		},
	}
}

// IdOrUsername returns the user's Id if set, otherwise the username property.
func (u *User) IdOrUsername() string {
	if u.Id != "" {
		return u.Id
	}
	return u.Username()
}

// Username returns the user's username property.
func (u *User) Username() string {
	if val, err := u.Properties[UserPropertyUsername].String(); err == nil && val != nil {
		return *val
	}
	return ""
}

// Password returns the user's password property.
func (u *User) Password() string {
	if val, err := u.Properties[UserPropertyPassword].String(); err == nil && val != nil {
		return *val
	}
	return ""
}

// FirstName returns the user's first name property.
func (u *User) FirstName() string {
	if val, err := u.Properties[UserPropertyFirstName].String(); err == nil && val != nil {
		return *val
	}
	return ""
}

// LastName returns the user's last name property.
func (u *User) LastName() string {
	if val, err := u.Properties[UserPropertyLastName].String(); err == nil && val != nil {
		return *val
	}
	return ""
}

// Email returns the user's email property.
func (u *User) Email() string {
	if val, err := u.Properties[UserPropertyEmail].String(); err == nil && val != nil {
		return *val
	}
	return ""
}

// Groups returns the list of group names the user belongs to.
func (u *User) Groups() []string {
	if val, err := u.Properties[UserPropertyGroups].StringList(); err == nil {
		return val
	}
	return nil
}

// Company returns the user's company property.
func (u *User) Company() string {
	if val, err := u.Properties[UserPropertyCompany].String(); err == nil && val != nil {
		return *val
	}
	return ""
}

// TenantId returns the user's tenant ID property.
func (u *User) TenantId() string {
	if val, err := u.Properties[UserPropertyTenantId].String(); err == nil && val != nil {
		return *val
	}
	return ""
}

// Property returns the Field value for the given property key.
func (u *User) Property(key string) (Field, bool) {
	value, found := u.Properties[key]
	return value, found
}

// SetProperty sets the value for the given property key.
func (u *User) SetProperty(key string, value Field) {
	u.Properties[key] = value
}

type Users paginableEntities[User]
