package nuxeo

type ACP struct {
	entity

	ACLs []ACL `json:"acl"`
}

type ACL struct {
	Name string `json:"name"`
	ACEs []ACE  `json:"ace"`
}

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
