package nuxeo

// Task represents a Nuxeo workflow task instance.
//
// A task is a step in a workflow assigned to one or more users, with metadata describing its state, actors, related documents, comments, and variables.
// See: https://doc.nuxeo.com/rest-api/1/task-entity-type/
// Task models the REST API 'task' entity-type.
// Fields map directly to Nuxeo's JSON representation.
type Task struct {
	entity
	Id                     string        `json:"id"`
	Name                   string        `json:"name"`
	WorkflowInstanceId     string        `json:"workflowInstanceId"`
	WorkflowModelName      string        `json:"workflowModelName"`
	WorkflowInitiator      User          `json:"workflowInitiator"` // TODO: json unmarshal user from string
	WorkflowTitle          string        `json:"workflowTitle"`
	WorkflowLifeCycleState string        `json:"workflowLifeCycleState"`
	GraphResource          string        `json:"graphResource"`
	State                  string        `json:"state"`
	Directive              string        `json:"directive"`
	Created                *ISO8601Time  `json:"created"`
	DueDate                *ISO8601Time  `json:"dueDate"`
	NodeName               string        `json:"nodeName"`
	TargetDocumentIds      []Document    `json:"targetDocumentIds"` // TODO: json unmarshal documents from strings
	Actors                 []User        `json:"actors"`            // TODO: json unmarshal users from { "id": "username" }
	DelegatedActors        []User        `json:"delegatedActors"`   // TODO: json unmarshal users from { "id": "username" }
	Comments               []TaskComment `json:"comments"`
	Variables              TaskVariables `json:"variables"`
	TaskInfo               TaskInfo      `json:"taskInfo"`
}

// NewTask creates a new entityTask with the given ID and sets the EntityType to 'task'.
func NewTask(id string) *Task {
	return &Task{
		entity: entity{
			EntityType: EntityTypeTask,
		},
		Id: id,
	}
}

// Tasks is a slice wrapper for multiple entityTask objects, as returned by Nuxeo task queries.
type Tasks entities[Task]

// TaskComment represents a comment on a workflow task, including author, text, and date.
type TaskComment struct {
	Author string       `json:"author"`
	Text   string       `json:"text"`
	Date   *ISO8601Time `json:"date"`
}

// TaskVariables holds custom variables for a workflow task, such as comment, assignees, end date, and participants.
type TaskVariables struct {
	Comment      string       `json:"comment"`
	Assignees    []string     `json:"assignees"`
	EndDate      *ISO8601Time `json:"end_date"`
	Participants []string     `json:"participants"`
}

// TaskInfo provides metadata about a workflow task, including allowed actions, reassignment, layout, and schemas.
type TaskInfo struct {
	AllowTaskReassignment bool           `json:"allowTaskReassignment"`
	TaskActions           []TaskInfoItem `json:"taskActions"`
	LayoutResource        TaskInfoItem   `json:"layoutResource"`
	Schemas               []TaskInfoItem `json:"schemas"`
}

// TaskInfoItem describes an actionable item or resource for a workflow task, such as an action, layout, or schema.
type TaskInfoItem struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Label string `json:"label"`
}
