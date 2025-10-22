package nuxeo

import (
	"time"
)

type ACP struct {
	entity

	ACLs []ACL `json:"acl"`
}

type ACL struct {
	Name string `json:"name"`
	ACEs []ACE  `json:"ace"`
}

type ACE struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	Permission string    `json:"permission"`
	Granted    bool      `json:"granted"`
	Creator    string    `json:"creator"`
	Begin      time.Time `json:"begin"`
	End        time.Time `json:"end"`
	Status     string    `json:"status"`
}
