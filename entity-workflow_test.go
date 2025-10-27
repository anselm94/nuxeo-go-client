package nuxeo

import (
	"reflect"
	"testing"
)

func TestNewWorkflow(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		id   string
	}{
		{"normal id", "wf-123"},
		{"empty id", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			wf := NewWorkflow(tc.id)
			if wf.Id != tc.id {
				t.Errorf("Id mismatch: got %q, want %q", wf.Id, tc.id)
			}
			if wf.entity.EntityType != EntityTypeWorkflow {
				t.Errorf("EntityType mismatch: got %q, want %q", wf.entity.EntityType, EntityTypeWorkflow)
			}
			// All other fields should be zero values
			if wf.Name != "" || wf.Title != "" || wf.State != "" || wf.WorkflowModelName != "" || wf.GraphResource != "" {
				t.Errorf("Unexpected non-zero string fields in new workflow")
			}
			if len(wf.AttachedDocumentIds) != 0 {
				t.Errorf("Expected AttachedDocumentIds to be empty")
			}
			if len(wf.Variables) != 0 {
				t.Errorf("Expected Variables to be nil or empty")
			}
		})
	}
}

func TestEntityWorkflowGraph(t *testing.T) {
	t.Parallel()

	nodes := map[string]Field{"n1": Field("1"), "n2": Field("2")}
	transitions := map[string]Field{"t1": Field("A"), "t2": Field("B")}
	graph := entityWorkflowGraph{
		entity:      entity{EntityType: EntityTypeGraph},
		Nodes:       nodes,
		Transitions: transitions,
	}
	if graph.entity.EntityType != EntityTypeGraph {
		t.Errorf("EntityType mismatch")
	}
	if !reflect.DeepEqual(graph.Nodes, nodes) {
		t.Errorf("Nodes mismatch")
	}
	if !reflect.DeepEqual(graph.Transitions, transitions) {
		t.Errorf("Transitions mismatch")
	}
	// Edge case: empty maps
	emptyGraph := entityWorkflowGraph{}
	if len(emptyGraph.Nodes) != 0 || len(emptyGraph.Transitions) != 0 {
		t.Errorf("Expected empty maps in zero value graph")
	}
}
