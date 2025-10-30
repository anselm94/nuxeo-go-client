package nuxeo

// Workflow represents a Nuxeo workflow instance or model.
//
// A workflow tracks the state, initiator, attached documents, variables, and graph for a business process.
// See: https://doc.nuxeo.com/rest-api/1/workflow-entity-type/
// Workflow models the REST API 'workflow' entity-type.
// Fields map directly to Nuxeo's JSON representation.
type Workflow struct {
	entity
	Id                  string           `json:"id"`
	Name                string           `json:"name"`
	Title               string           `json:"title"`
	State               string           `json:"state"`
	WorkflowModelName   string           `json:"workflowModelName"`
	Initiator           User             `json:"initiator"`           // TODO: JSON unmarshal string to user
	AttachedDocumentIds []Document       `json:"attachedDocumentIds"` // TODO: JSON unmarshal string to document
	Variables           map[string]Field `json:"variables"`
	GraphResource       string           `json:"graphResource"`
}

// NewWorkflow creates a new entityWorkflow with the given ID and sets the EntityType to 'workflow'.
func NewWorkflow(id string) *Workflow {
	return &Workflow{
		entity: entity{
			EntityType: EntityTypeWorkflow,
		},
		Id: id,
	}
}

// Workflows is a slice wrapper for multiple entityWorkflow objects, as returned by Nuxeo workflow queries.
type Workflows entities[Workflow]

// WorkflowGraph represents the graph structure of a workflow instance or model, including nodes and transitions.
type WorkflowGraph struct {
	entity
	Nodes       map[string]Field `json:"nodes"`
	Transitions map[string]Field `json:"transitions"`
}
