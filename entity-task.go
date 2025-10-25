package nuxeo

// entityTask represents a workflow task.
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

func NewTask(id string) *entityTask {
	return &entityTask{
		entity: entity{
			EntityType: EntityTypeTask,
		},
		Id: id,
	}
}

type entityTasks entities[entityTask]

type entityTaskComment struct {
	Author string       `json:"author"`
	Text   string       `json:"text"`
	Date   *ISO8601Time `json:"date"`
}

type entityTaskVariables struct {
	Comment      string       `json:"comment"`
	Assignees    []string     `json:"assignees"`
	EndDate      *ISO8601Time `json:"end_date"`
	Participants []string     `json:"participants"`
}

type entityTaskInfo struct {
	AllowTaskReassignment bool                 `json:"allowTaskReassignment"`
	TaskActions           []entityTaskInfoItem `json:"taskActions"`
	LayoutResource        entityTaskInfoItem   `json:"layoutResource"`
	Schemas               []entityTaskInfoItem `json:"schemas"`
}

type entityTaskInfoItem struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Label string `json:"label"`
}
