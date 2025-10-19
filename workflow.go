package nuxeo

// Workflow represents a Nuxeo workflow.
type Workflow struct {
	ID         string
	Name       string
	Properties map[string]any
}

// NewWorkflow creates a new Workflow instance.
func NewWorkflow(id, name string, props map[string]any) *Workflow {
	return &Workflow{
		ID:         id,
		Name:       name,
		Properties: props,
	}
}

// Task represents a workflow task.
type Task struct {
	ID         string
	Name       string
	Properties map[string]any
}

// NewTask creates a new Task instance.
func NewTask(id, name string, props map[string]any) *Task {
	return &Task{
		ID:         id,
		Name:       name,
		Properties: props,
	}
}
