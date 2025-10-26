package nuxeo

// EntityACP represents an Access Control Policy (ACP) for a Nuxeo document.
// It contains a list of ACLs (Access Control Lists) that define permissions for users and groups.
type entityACP struct {
	entity

	ACLs []entityACL `json:"acl"`
}

// NewACP creates a new EntityACP with the ACP entity type set.
func NewACP() *entityACP {
	return &entityACP{
		entity: entity{
			EntityType: EntityTypeACP,
		},
	}
}

// EntityACL represents an Access Control List (ACL) within an ACP.
// An ACL contains a name and a list of ACEs (Access Control Entries).
type entityACL struct {
	Name string      `json:"name"`
	ACEs []entityACE `json:"ace"`
}

// NewACL creates a new EntityACL with the specified name.
func NewACL(name string) *entityACL {
	return &entityACL{
		Name: name,
	}
}

// EntityACE represents an Access Control Entry (ACE) within an ACL.
// An ACE defines a permission grant or denial for a user or group, with optional time bounds and status.
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

// NewACE creates a new EntityACE for the specified username, permission, and grant status.
func NewACE(username, permission string, granted bool) *entityACE {
	return &entityACE{
		Username:   username,
		Permission: permission,
		Granted:    granted,
	}
}
