package nuxeo

// EntityGroup represents a Nuxeo group entity.
// It includes group properties, member users, member groups, and parent groups.
type Group struct {
	entity
	Id           string           `json:"id"`
	Properties   map[string]Field `json:"properties"`
	MemberUsers  []string         `json:"memberUsers"`
	MemberGroups []string         `json:"memberGroups"`
	ParentGroups []string         `json:"parentGroups"`
}

// NewGroup creates a new EntityGroup with the given group ID.
// The group ID is set as both the Id and the "id" property.
func NewGroup(groupId string) *Group {
	return &Group{
		entity: entity{
			EntityType: EntityTypeGroup,
		},
		Id: groupId,
	}
}

type Groups paginableEntities[Group]
