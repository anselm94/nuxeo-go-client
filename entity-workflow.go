package nuxeo

// Workflow represents a Nuxeo workflow.
type Workflow struct {
	entity
	Id                  string         `json:"id"`
	Name                string         `json:"name"`
	Title               string         `json:"title"`
	State               string         `json:"state"`
	WorkflowModelName   string         `json:"workflowModelName"`
	Initiator           User           `json:"initiator"`           // TODO: JSON unmarshal string to user
	AttachedDocumentIds []Document     `json:"attachedDocumentIds"` // TODO: JSON unmarshal string to document
	Variables           map[string]any `json:"variables"`
	GraphResource       string         `json:"graphResource"`
}

type Workflows entities[Workflow]

type WorkflowGraph struct {
	entity
	Nodes       map[string]any `json:"nodes"`
	Transitions map[string]any `json:"transitions"`
}
