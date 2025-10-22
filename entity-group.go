package nuxeo

// Group represents a Nuxeo group.
type Group struct {
	entity
	Id           string         `json:"id"`
	Properties   map[string]any `json:"properties"`
	MemberUsers  []string       `json:"memberUsers"`
	MemberGroups []string       `json:"memberGroups"`
	ParentGroups []string       `json:"parentGroups"`
}

type Groups paginableEntities[Group]
