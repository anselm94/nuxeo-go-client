package nuxeo

// entityWorkflow represents a Nuxeo workflow.
type entityWorkflow struct {
	entity
	Id                  string           `json:"id"`
	Name                string           `json:"name"`
	Title               string           `json:"title"`
	State               string           `json:"state"`
	WorkflowModelName   string           `json:"workflowModelName"`
	Initiator           entityUser       `json:"initiator"`           // TODO: JSON unmarshal string to user
	AttachedDocumentIds []entityDocument `json:"attachedDocumentIds"` // TODO: JSON unmarshal string to document
	Variables           map[string]any   `json:"variables"`
	GraphResource       string           `json:"graphResource"`
}

func NewWorkflow(id string) *entityWorkflow {
	return &entityWorkflow{
		entity: entity{
			EntityType: EntityTypeWorkflow,
		},
		Id: id,
	}
}

type entityWorkflows entities[entityWorkflow]

type entityWorkflowGraph struct {
	entity
	Nodes       map[string]any `json:"nodes"`
	Transitions map[string]any `json:"transitions"`
}
