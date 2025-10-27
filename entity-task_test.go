package nuxeo

import (
	"reflect"
	"testing"
)

func TestNewTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		id   string
	}{
		{"normal id", "task-123"},
		{"empty id", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			task := NewTask(tc.id)
			if task.Id != tc.id {
				t.Errorf("Id mismatch: got %q, want %q", task.Id, tc.id)
			}
			if task.entity.EntityType != EntityTypeTask {
				t.Errorf("EntityType mismatch: got %q, want %q", task.entity.EntityType, EntityTypeTask)
			}
			// All other fields should be zero values
			if task.Name != "" || task.WorkflowInstanceId != "" || task.WorkflowModelName != "" || task.WorkflowTitle != "" || task.State != "" {
				t.Errorf("Unexpected non-zero string fields in new task")
			}
			if task.Created != nil || task.DueDate != nil {
				t.Errorf("Expected Created/DueDate to be nil")
			}
			if len(task.TargetDocumentIds) != 0 || len(task.Actors) != 0 || len(task.DelegatedActors) != 0 || len(task.Comments) != 0 {
				t.Errorf("Expected slices to be empty")
			}
		})
	}
}

func TestEntityTaskComment(t *testing.T) {
	t.Parallel()

	date := (*ISO8601Time)(nil)
	comment := entityTaskComment{Author: "bob", Text: "review", Date: date}
	if comment.Author != "bob" {
		t.Errorf("Author mismatch")
	}
	if comment.Text != "review" {
		t.Errorf("Text mismatch")
	}
	if comment.Date != nil {
		t.Errorf("Date should be nil")
	}
}

func TestEntityTaskVariables(t *testing.T) {
	t.Parallel()

	endDate := (*ISO8601Time)(nil)
	vars := entityTaskVariables{
		Comment:      "needs work",
		Assignees:    []string{"alice", "bob"},
		EndDate:      endDate,
		Participants: []string{"alice"},
	}
	if vars.Comment != "needs work" {
		t.Errorf("Comment mismatch")
	}
	if !reflect.DeepEqual(vars.Assignees, []string{"alice", "bob"}) {
		t.Errorf("Assignees mismatch")
	}
	if vars.EndDate != nil {
		t.Errorf("EndDate should be nil")
	}
	if !reflect.DeepEqual(vars.Participants, []string{"alice"}) {
		t.Errorf("Participants mismatch")
	}
}

func TestEntityTaskInfoAndItem(t *testing.T) {
	t.Parallel()

	actions := []entityTaskInfoItem{{Name: "approve", Url: "/approve", Label: "Approve"}}
	layout := entityTaskInfoItem{Name: "main", Url: "/layout", Label: "Main"}
	schemas := []entityTaskInfoItem{{Name: "schema1", Url: "/s1", Label: "Schema1"}}
	info := entityTaskInfo{
		AllowTaskReassignment: true,
		TaskActions:           actions,
		LayoutResource:        layout,
		Schemas:               schemas,
	}
	if !info.AllowTaskReassignment {
		t.Errorf("AllowTaskReassignment should be true")
	}
	if !reflect.DeepEqual(info.TaskActions, actions) {
		t.Errorf("TaskActions mismatch")
	}
	if info.LayoutResource != layout {
		t.Errorf("LayoutResource mismatch")
	}
	if !reflect.DeepEqual(info.Schemas, schemas) {
		t.Errorf("Schemas mismatch")
	}
	// Edge case: empty slices
	infoEmpty := entityTaskInfo{}
	if len(infoEmpty.TaskActions) != 0 || len(infoEmpty.Schemas) != 0 {
		t.Errorf("Expected empty slices in zero value info")
	}
}
