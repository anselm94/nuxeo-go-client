package nuxeo

// entityGroup represents a Nuxeo group.
type entityGroup struct {
	entity
	Id           string           `json:"id"`
	Properties   map[string]Field `json:"properties"`
	MemberUsers  []string         `json:"memberUsers"`
	MemberGroups []string         `json:"memberGroups"`
	ParentGroups []string         `json:"parentGroups"`
}

func NewGroup(groupId string) *entityGroup {
	return &entityGroup{
		entity: entity{
			EntityType: EntityTypeGroup,
		},
		Id: groupId,
	}
}

type entityGroups paginableEntities[entityGroup]
