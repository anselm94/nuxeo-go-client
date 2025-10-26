package nuxeo

// entityTask represents a Nuxeo workflow task instance.
//
// A task is a step in a workflow assigned to one or more users, with metadata describing its state, actors, related documents, comments, and variables.
// See: https://doc.nuxeo.com/rest-api/1/task-entity-type/
// entityTask models the REST API 'task' entity-type.
// Fields map directly to Nuxeo's JSON representation.
type entityTask struct {
	entity
	Id                     string              `json:"id"`
	Name                   string              `json:"name"`
	WorkflowInstanceId     string              `json:"workflowInstanceId"`
	WorkflowModelName      string              `json:"workflowModelName"`
	WorkflowInitiator      entityUser          `json:"workflowInitiator"` // TODO: json unmarshal user from string
	WorkflowTitle          string              `json:"workflowTitle"`
	WorkflowLifeCycleState string              `json:"workflowLifeCycleState"`
	GraphResource          string              `json:"graphResource"`
	State                  string              `json:"state"`
	Directive              string              `json:"directive"`
	Created                *ISO8601Time        `json:"created"`
	DueDate                *ISO8601Time        `json:"dueDate"`
	NodeName               string              `json:"nodeName"`
	TargetDocumentIds      []entityDocument    `json:"targetDocumentIds"` // TODO: json unmarshal documents from strings
	Actors                 []entityUser        `json:"actors"`            // TODO: json unmarshal users from { "id": "username" }
	DelegatedActors        []entityUser        `json:"delegatedActors"`   // TODO: json unmarshal users from { "id": "username" }
	Comments               []entityTaskComment `json:"comments"`
	Variables              entityTaskVariables `json:"variables"`
	TaskInfo               entityTaskInfo      `json:"taskInfo"`
}

// NewTask creates a new entityTask with the given ID and sets the EntityType to 'task'.
func NewTask(id string) *entityTask {
	return &entityTask{
		entity: entity{
			EntityType: EntityTypeTask,
		},
		Id: id,
	}
}

// entityTasks is a slice wrapper for multiple entityTask objects, as returned by Nuxeo task queries.
type entityTasks entities[entityTask]

// entityTaskComment represents a comment on a workflow task, including author, text, and date.
type entityTaskComment struct {
	Author string       `json:"author"`
	Text   string       `json:"text"`
	Date   *ISO8601Time `json:"date"`
}

// entityTaskVariables holds custom variables for a workflow task, such as comment, assignees, end date, and participants.
type entityTaskVariables struct {
	Comment      string       `json:"comment"`
	Assignees    []string     `json:"assignees"`
	EndDate      *ISO8601Time `json:"end_date"`
	Participants []string     `json:"participants"`
}

// entityTaskInfo provides metadata about a workflow task, including allowed actions, reassignment, layout, and schemas.
type entityTaskInfo struct {
	AllowTaskReassignment bool                 `json:"allowTaskReassignment"`
	TaskActions           []entityTaskInfoItem `json:"taskActions"`
	LayoutResource        entityTaskInfoItem   `json:"layoutResource"`
	Schemas               []entityTaskInfoItem `json:"schemas"`
}

// entityTaskInfoItem describes an actionable item or resource for a workflow task, such as an action, layout, or schema.
type entityTaskInfoItem struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Label string `json:"label"`
}
