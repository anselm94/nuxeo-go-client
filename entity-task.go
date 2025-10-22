package nuxeo

import "time"

// Task represents a workflow task.
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
	Created                time.Time     `json:"created"`
	DueDate                time.Time     `json:"dueDate"`
	NodeName               string        `json:"nodeName"`
	TargetDocumentIds      []Document    `json:"targetDocumentIds"` // TODO: json unmarshal documents from strings
	Actors                 []User        `json:"actors"`            // TODO: json unmarshal users from { "id": "username" }
	DelegatedActors        []User        `json:"delegatedActors"`   // TODO: json unmarshal users from { "id": "username" }
	Comments               []TaskComment `json:"comments"`
	Variables              TaskVariables `json:"variables"`
	TaskInfo               TaskInfo      `json:"taskInfo"`
}

type Tasks entities[Task]

type TaskComment struct {
	Author string    `json:"author"`
	Text   string    `json:"text"`
	Date   time.Time `json:"date"`
}

type TaskVariables struct {
	Comment      string    `json:"comment"`
	Assignees    []string  `json:"assignees"`
	EndDate      time.Time `json:"end_date"`
	Participants []string  `json:"participants"`
}

type TaskInfo struct {
	AllowTaskReassignment bool           `json:"allowTaskReassignment"`
	TaskActions           []TaskInfoItem `json:"taskActions"`
	LayoutResource        TaskInfoItem   `json:"layoutResource"`
	Schemas               []TaskInfoItem `json:"schemas"`
}

type TaskInfoItem struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Label string `json:"label"`
}
