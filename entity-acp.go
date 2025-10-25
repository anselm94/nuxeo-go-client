package nuxeo

type entityACP struct {
	entity

	ACLs []entityACL `json:"acl"`
}

func NewACP() *entityACP {
	return &entityACP{
		entity: entity{
			EntityType: EntityTypeACP,
		},
	}
}

type entityACL struct {
	Name string      `json:"name"`
	ACEs []entityACE `json:"ace"`
}

func NewACL(name string) *entityACL {
	return &entityACL{
		Name: name,
	}
}

type entityACE struct {
	ID         string       `json:"id"`
	Username   string       `json:"username"`
	Permission string       `json:"permission"`
	Granted    bool         `json:"granted"`
	Creator    string       `json:"creator"`
	Begin      *ISO8601Time `json:"begin"`
	End        *ISO8601Time `json:"end"`
	Status     string       `json:"status"`
}

func NewACE(username, permission string, granted bool) *entityACE {
	return &entityACE{
		Username:   username,
		Permission: permission,
		Granted:    granted,
	}
}
