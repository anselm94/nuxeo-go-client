package nuxeo

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
	Properties      map[string]any        `json:"properties"`
	ExtendedGroups  []entityExtendedGroup `json:"extendedGroups,omitempty"`
}

func NewUser(username string) *entityUser {
	properties := make(map[string]any)
	properties[UserPropertyUsername] = username // Set username by default
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
	return u.Properties[UserPropertyUsername].(string)
}

func (u *entityUser) Password() string {
	return u.Properties[UserPropertyPassword].(string)
}

func (u *entityUser) FirstName() string {
	return u.Properties[UserPropertyFirstName].(string)
}

func (u *entityUser) LastName() string {
	return u.Properties[UserPropertyLastName].(string)
}

func (u *entityUser) Email() string {
	return u.Properties[UserPropertyEmail].(string)
}

func (u *entityUser) Groups() []string {
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

func (u *entityUser) Company() string {
	return u.Properties[UserPropertyCompany].(string)
}

func (u *entityUser) TenantId() string {
	return u.Properties[UserPropertyTenantId].(string)
}

func (u *entityUser) Property(key string) any {
	return u.Properties[key]
}

func (u *entityUser) SetProperty(key string, value any) {
	u.Properties[key] = value
}

type entityUsers paginableEntities[entityUser]
