package nuxeo

import "testing"

func TestNewRepository(t *testing.T) {
	repo := NewRepository("default")
	if repo.Name != "default" {
		t.Errorf("Name got %q, want %q", repo.Name, "default")
	}
}

func TestRepository_GetDocument_Stub(t *testing.T) {
	repo := NewRepository("default")
	doc, err := repo.GetDocument("123")
	if doc != nil {
		t.Errorf("GetDocument got %v, want nil (stub)", doc)
	}
	if err != nil {
		t.Errorf("GetDocument err got %v, want nil (stub)", err)
	}
}

func TestRepository_QueryDocuments_Stub(t *testing.T) {
	repo := NewRepository("default")
	docs, err := repo.QueryDocuments("SELECT * FROM Document")
	if docs != nil {
		t.Errorf("QueryDocuments got %v, want nil (stub)", docs)
	}
	if err != nil {
		t.Errorf("QueryDocuments err got %v, want nil (stub)", err)
	}
}
