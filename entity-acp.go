package nuxeo

// ACP represents an Access Control Policy (ACP) for a Nuxeo document.
// It contains a list of ACLs (Access Control Lists) that define permissions for users and groups.
type ACP struct {
	entity

	ACLs []ACL `json:"acl"`
}

// NewACP creates a new ACP with the ACP entity type set.
func NewACP() *ACP {
	return &ACP{
		entity: entity{
			EntityType: EntityTypeACP,
		},
	}
}

// ACL represents an Access Control List (ACL) within an ACP.
// An ACL contains a name and a list of ACEs (Access Control Entries).
type ACL struct {
	Name string `json:"name"`
	ACEs []ACE  `json:"aces"`
}

// NewACL creates a new ACL with the specified name.
func NewACL(name string) *ACL {
	return &ACL{
		Name: name,
	}
}

// ACE represents an Access Control Entry (ACE) within an ACL.
// An ACE defines a permission grant or denial for a user or group, with optional time bounds and status.
type ACE struct {
	ID         string       `json:"id"`
	Username   string       `json:"username"`
	Permission string       `json:"permission"`
	Granted    bool         `json:"granted"`
	Creator    string       `json:"creator"`
	Begin      *ISO8601Time `json:"begin"`
	End        *ISO8601Time `json:"end"`
	Status     string       `json:"status"`
}

// NewACE creates a new ACE for the specified username, permission, and grant status.
func NewACE(username, permission string, granted bool) *ACE {
	return &ACE{
		Username:   username,
		Permission: permission,
		Granted:    granted,
	}
}
