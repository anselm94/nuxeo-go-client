package nuxeo

type ExtendedGroup struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	Url   string `json:"url"`
}

// User represents a Nuxeo user.
type User struct {
	entity
	Id              string          `json:"id"`
	IsAdministrator bool            `json:"isAdministrator"`
	IsAnonymous     bool            `json:"isAnonymous"`
	Properties      map[string]any  `json:"properties"`
	ExtendedGroups  []ExtendedGroup `json:"extendedGroups,omitempty"`
}

func (u *User) IdOrUsername() string {
	if u.Id != "" {
		return u.Id
	}
	return u.Username()
}

func (u *User) Username() string {
	return u.Properties[UserPropertyUsername].(string)
}

func (u *User) Password() string {
	return u.Properties[UserPropertyPassword].(string)
}

func (u *User) FirstName() string {
	return u.Properties[UserPropertyFirstName].(string)
}

func (u *User) LastName() string {
	return u.Properties[UserPropertyLastName].(string)
}

func (u *User) Email() string {
	return u.Properties[UserPropertyEmail].(string)
}

func (u *User) Groups() []string {
	groups, ok := u.Properties[UserPropertyGroups].([]any)
	if !ok {
		return []string{}
	}
	result := make([]string, len(groups))
	for i, g := range groups {
		result[i] = g.(string)
	}
	return result
}

func (u *User) Company() string {
	return u.Properties[UserPropertyCompany].(string)
}

func (u *User) TenantId() string {
	return u.Properties[UserPropertyTenantId].(string)
}

type Users paginableEntities[User]
