package nuxeo

import "testing"

func TestNewWorkflow(t *testing.T) {
	props := map[string]any{"status": "active"}
	wf := NewWorkflow("wf1", "Approval", props)
	if wf.ID != "wf1" {
		t.Errorf("ID got %q, want %q", wf.ID, "wf1")
	}
	if wf.Name != "Approval" {
		t.Errorf("Name got %q, want %q", wf.Name, "Approval")
	}
	if wf.Properties["status"] != "active" {
		t.Errorf("Properties[status] got %v, want %v", wf.Properties["status"], "active")
	}
}

func TestNewTask(t *testing.T) {
	props := map[string]any{"priority": "high"}
	task := NewTask("t1", "Review", props)
	if task.ID != "t1" {
		t.Errorf("ID got %q, want %q", task.ID, "t1")
	}
	if task.Name != "Review" {
		t.Errorf("Name got %q, want %q", task.Name, "Review")
	}
	if task.Properties["priority"] != "high" {
		t.Errorf("Properties[priority] got %v, want %v", task.Properties["priority"], "high")
	}
}
