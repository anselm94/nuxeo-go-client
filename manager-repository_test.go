package nuxeo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"testing"
)

// TestRepository_Name verifies that the Name() method returns the correct repository name.
func TestRepository_Name(t *testing.T) {
	repo := &repository{name: "default"}
	got := repo.Name()
	want := "default"
	if got != want {
		t.Errorf("repository.Name() = %v, want %v", got, want)
	}
}

// newTestRepository returns a repository manager with a mock client.
func newTestRepository(respond func(req *http.Request) (*http.Response, error)) *repository {
	client := newMockNuxeoClient(respond)
	return &repository{
		name:   "default",
		client: client,
		logger: slog.Default(),
	}
}

func TestRepository_FetchDocumentRoot(t *testing.T) {
	tests := []struct {
		name       string
		mockResp   any
		mockStatus int
		mockErr    error
		wantErr    bool
		wantDocID  string
	}{
		{
			name:       "success",
			mockResp:   &Document{ID: "rootdoc"},
			mockStatus: 200,
			wantDocID:  "rootdoc",
		},
		{
			name:       "not found",
			mockResp:   &NuxeoError{Message: "not found"},
			mockStatus: 404,
			wantErr:    true,
		},
		{
			name:    "http error",
			mockErr: errors.New("network error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				body, _ := json.Marshal(tt.mockResp)
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			})
			got, err := repo.FetchDocumentRoot(context.Background(), nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchDocumentRoot() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got.ID != tt.wantDocID {
				t.Errorf("FetchDocumentRoot() got.ID = %v, want %v", got.ID, tt.wantDocID)
			}
		})
	}
}

func TestRepository_FetchDocumentById(t *testing.T) {
	tests := []struct {
		name       string
		docID      string
		mockResp   any
		mockStatus int
		mockErr    error
		wantErr    bool
		wantDocID  string
	}{
		{
			name:       "success",
			docID:      "doc123",
			mockResp:   &Document{ID: "doc123"},
			mockStatus: 200,
			wantDocID:  "doc123",
		},
		{
			name:       "not found",
			docID:      "missing",
			mockResp:   &NuxeoError{Message: "not found"},
			mockStatus: 404,
			wantErr:    true,
		},
		{
			name:    "http error",
			docID:   "err",
			mockErr: errors.New("network error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				body, _ := json.Marshal(tt.mockResp)
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			})
			got, err := repo.FetchDocumentById(context.Background(), tt.docID, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchDocumentById() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got.ID != tt.wantDocID {
				t.Errorf("FetchDocumentById() got.ID = %v, want %v", got.ID, tt.wantDocID)
			}
		})
	}
}

func TestRepository_FetchDocumentByPath(t *testing.T) {
	tests := []struct {
		name       string
		docPath    string
		mockResp   any
		mockStatus int
		mockErr    error
		wantErr    bool
		wantDocID  string
	}{
		{
			name:       "success",
			docPath:    "/default-domain/workspaces",
			mockResp:   &Document{ID: "docpath123"},
			mockStatus: 200,
			wantDocID:  "docpath123",
		},
		{
			name:       "not found",
			docPath:    "/missing",
			mockResp:   &NuxeoError{Message: "not found"},
			mockStatus: 404,
			wantErr:    true,
		},
		{
			name:    "http error",
			docPath: "/err",
			mockErr: errors.New("network error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				body, _ := json.Marshal(tt.mockResp)
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			})
			got, err := repo.FetchDocumentByPath(context.Background(), tt.docPath, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchDocumentByPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got.ID != tt.wantDocID {
				t.Errorf("FetchDocumentByPath() got.ID = %v, want %v", got.ID, tt.wantDocID)
			}
		})
	}
}

func TestRepository_CreateDocumentById(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		doc := Document{ID: "newdoc"}
		repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&doc)
			return &http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		})
		got, err := repo.CreateDocumentById(context.Background(), "parent123", doc, nil)
		if err != nil {
			t.Fatalf("CreateDocumentById() error = %v, want nil", err)
		}
		if got.ID != "newdoc" {
			t.Errorf("CreateDocumentById() got.ID = %v, want newdoc", got.ID)
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&NuxeoError{Message: "bad request"})
			return &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		})
		_, err := repo.CreateDocumentById(context.Background(), "parent123", Document{ID: "fail"}, nil)
		if err == nil {
			t.Errorf("CreateDocumentById() error = nil, want error")
		}
	})
}

func TestRepository_CreateDocumentByPath(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		doc := Document{ID: "createdbyPath"}
		repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&doc)
			return &http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		})
		got, err := repo.CreateDocumentByPath(context.Background(), "/default-domain/workspaces", doc, nil)
		if err != nil {
			t.Fatalf("CreateDocumentByPath() error = %v, want nil", err)
		}
		if got.ID != "createdbyPath" {
			t.Errorf("CreateDocumentByPath() got.ID = %v, want createdbyPath", got.ID)
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&NuxeoError{Message: "bad request"})
			return &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		})
		_, err := repo.CreateDocumentByPath(context.Background(), "/default-domain/workspaces", Document{ID: "fail"}, nil)
		if err == nil {
			t.Errorf("CreateDocumentByPath() error = nil, want error")
		}
	})
}

func TestRepository_UpdateDocument(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		doc := Document{ID: "updatedDoc"}
		repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&doc)
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		})
		got, err := repo.UpdateDocument(context.Background(), "updatedDoc", doc, nil)
		if err != nil {
			t.Fatalf("UpdateDocument() error = %v, want nil", err)
		}
		if got.ID != "updatedDoc" {
			t.Errorf("UpdateDocument() got.ID = %v, want updatedDoc", got.ID)
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&NuxeoError{Message: "update failed"})
			return &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		})
		_, err := repo.UpdateDocument(context.Background(), "failDoc", Document{ID: "failDoc"}, nil)
		if err == nil {
			t.Errorf("UpdateDocument() error = nil, want error")
		}
	})
}

func TestRepository_DeleteDocument(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 204,
				Body:       io.NopCloser(bytes.NewReader([]byte{})),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		})
		err := repo.DeleteDocument(context.Background(), "doc123")
		if err != nil {
			t.Errorf("DeleteDocument() error = %v, want nil", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&NuxeoError{Message: "not found"})
			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		})
		err := repo.DeleteDocument(context.Background(), "ghost")
		if err == nil {
			t.Errorf("DeleteDocument() error = nil, want error")
		}
	})
}

func TestRepository_Query(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		docs := Documents{}
		docs.Entries = []Document{{ID: "doc1"}, {ID: "doc2"}}
		repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&docs)
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		})
		got, err := repo.Query(context.Background(), "SELECT * FROM Document", nil, nil, nil)
		if err != nil {
			t.Fatalf("Query() error = %v, want nil", err)
		}
		if len(got.Entries) != 2 {
			t.Errorf("Query() got %d entries, want 2", len(got.Entries))
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&NuxeoError{Message: "bad query"})
			return &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		})
		_, err := repo.Query(context.Background(), "bad", nil, nil, nil)
		if err == nil {
			t.Errorf("Query() error = nil, want error")
		}
	})
}

func TestRepository_QueryByProvider(t *testing.T) {
	tests := []struct {
		name             string
		providerName     string
		queryParams      []string
		namedQueryParams map[string]string
		mockResp         any
		mockStatus       int
		mockErr          error
		wantErr          bool
		wantDocIDs       []string
	}{
		{
			name:             "success with docs",
			providerName:     "myProvider",
			queryParams:      []string{"foo=bar"},
			namedQueryParams: map[string]string{"nxql": "SELECT * FROM Document"},
			mockResp: &Documents{
				Entries: []Document{{ID: "docA"}, {ID: "docB"}},
			},
			mockStatus: 200,
			wantDocIDs: []string{"docA", "docB"},
		},
		{
			name:             "empty result",
			providerName:     "emptyProvider",
			queryParams:      nil,
			namedQueryParams: map[string]string{},
			mockResp: &Documents{
				Entries: []Document{},
			},
			mockStatus: 200,
			wantDocIDs: []string{},
		},
		{
			name:             "provider not found",
			providerName:     "missingProvider",
			queryParams:      nil,
			namedQueryParams: map[string]string{},
			mockResp:         &NuxeoError{Message: "provider not found"},
			mockStatus:       404,
			wantErr:          true,
		},
		{
			name:         "http error",
			providerName: "errProvider",
			mockErr:      errors.New("network error"),
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				body, _ := json.Marshal(tt.mockResp)
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			})
			got, err := repo.QueryByProvider(
				context.Background(),
				tt.providerName,
				tt.queryParams,
				tt.namedQueryParams,
				nil,
				nil,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryByProvider() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				var gotIDs []string
				for _, doc := range got.Entries {
					gotIDs = append(gotIDs, doc.ID)
				}
				if len(gotIDs) != len(tt.wantDocIDs) {
					t.Errorf("QueryByProvider() got %d entries, want %d", len(gotIDs), len(tt.wantDocIDs))
				}
				for i, wantID := range tt.wantDocIDs {
					if gotIDs[i] != wantID {
						t.Errorf("QueryByProvider() entry[%d].ID = %v, want %v", i, gotIDs[i], wantID)
					}
				}
			}
		})
	}
}

func TestRepository_FetchAuditByPath(t *testing.T) {
	tests := []struct {
		name        string
		docPath     string
		mockResp    any
		mockStatus  int
		mockErr     error
		wantErr     bool
		wantAuditID int
	}{
		{
			name:        "success",
			docPath:     "/default-domain/workspaces",
			mockResp:    &Audit{Entries: []AuditLogEntry{{ID: 101}}},
			mockStatus:  200,
			wantAuditID: 101,
		},
		{
			name:       "not found",
			docPath:    "/missing",
			mockResp:   &NuxeoError{Message: "not found"},
			mockStatus: 404,
			wantErr:    true,
		},
		{
			name:    "http error",
			docPath: "/err",
			mockErr: errors.New("network error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				body, _ := json.Marshal(tt.mockResp)
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			})
			got, err := repo.FetchAuditByPath(context.Background(), tt.docPath, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchAuditByPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && len(got.Entries) > 0 && got.Entries[0].ID != tt.wantAuditID {
				t.Errorf("FetchAuditByPath() got.Entries[0].ID = %v, want %v", got.Entries[0].ID, tt.wantAuditID)
			}
		})
	}
}

func TestRepository_FetchAuditById(t *testing.T) {
	tests := []struct {
		name        string
		docID       string
		mockResp    any
		mockStatus  int
		mockErr     error
		wantErr     bool
		wantAuditID int
	}{
		{
			name:        "success",
			docID:       "doc123",
			mockResp:    &Audit{Entries: []AuditLogEntry{{ID: 202}}},
			mockStatus:  200,
			wantAuditID: 202,
		},
		{
			name:       "not found",
			docID:      "missing",
			mockResp:   &NuxeoError{Message: "not found"},
			mockStatus: 404,
			wantErr:    true,
		},
		{
			name:    "http error",
			docID:   "err",
			mockErr: errors.New("network error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				body, _ := json.Marshal(tt.mockResp)
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			})
			got, err := repo.FetchAuditById(context.Background(), tt.docID, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchAuditById() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && len(got.Entries) > 0 && got.Entries[0].ID != tt.wantAuditID {
				t.Errorf("FetchAuditById() got.Entries[0].ID = %v, want %v", got.Entries[0].ID, tt.wantAuditID)
			}
		})
	}
}

func TestRepository_FetchPermissionsByPath(t *testing.T) {
	tests := []struct {
		name       string
		docPath    string
		mockResp   any
		mockStatus int
		mockErr    error
		wantErr    bool
		wantACLs   int
		wantACE    string
	}{
		{
			name:       "success",
			docPath:    "/default-domain/workspaces",
			mockResp:   &ACP{ACLs: []ACL{{Name: "local", ACEs: []ACE{{Username: "bob"}}}}},
			mockStatus: 200,
			wantACLs:   1,
			wantACE:    "bob",
		},
		{
			name:       "not found",
			docPath:    "/missing",
			mockResp:   &NuxeoError{Message: "not found"},
			mockStatus: 404,
			wantErr:    true,
		},
		{
			name:    "http error",
			docPath: "/err",
			mockErr: errors.New("network error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				body, _ := json.Marshal(tt.mockResp)
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			})
			got, err := repo.FetchPermissionsByPath(context.Background(), tt.docPath, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchPermissionsByPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				if len(got.ACLs) != tt.wantACLs {
					t.Errorf("FetchPermissionsByPath() got %d ACLs, want %d", len(got.ACLs), tt.wantACLs)
				}
				if tt.wantACLs > 0 && len(got.ACLs[0].ACEs) > 0 && got.ACLs[0].ACEs[0].Username != tt.wantACE {
					t.Errorf("FetchPermissionsByPath() ACE Username = %v, want %v", got.ACLs[0].ACEs[0].Username, tt.wantACE)
				}
			}
		})
	}
}

func TestRepository_FetchPermissionsById(t *testing.T) {
	tests := []struct {
		name       string
		docID      string
		mockResp   any
		mockStatus int
		mockErr    error
		wantErr    bool
		wantACLs   int
		wantACE    string
	}{
		{
			name:       "success",
			docID:      "doc123",
			mockResp:   &ACP{ACLs: []ACL{{Name: "local", ACEs: []ACE{{Username: "alice"}}}}},
			mockStatus: 200,
			wantACLs:   1,
			wantACE:    "alice",
		},
		{
			name:       "not found",
			docID:      "missing",
			mockResp:   &NuxeoError{Message: "not found"},
			mockStatus: 404,
			wantErr:    true,
		},
		{
			name:    "http error",
			docID:   "err",
			mockErr: errors.New("network error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				body, _ := json.Marshal(tt.mockResp)
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			})
			got, err := repo.FetchPermissionsById(context.Background(), tt.docID, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchPermissionsById() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				if len(got.ACLs) != tt.wantACLs {
					t.Errorf("FetchPermissionsById() got %d ACLs, want %d", len(got.ACLs), tt.wantACLs)
				}
				if tt.wantACLs > 0 && len(got.ACLs[0].ACEs) > 0 && got.ACLs[0].ACEs[0].Username != tt.wantACE {
					t.Errorf("FetchPermissionsById() ACE Username = %v, want %v", got.ACLs[0].ACEs[0].Username, tt.wantACE)
				}
			}
		})
	}
}

func TestRepository_FetchChildrenByPath(t *testing.T) {
	tests := []struct {
		name       string
		parentPath string
		mockResp   any
		mockStatus int
		mockErr    error
		wantErr    bool
		wantDocIDs []string
	}{
		{
			name:       "success",
			parentPath: "/default-domain/workspaces",
			mockResp:   &Documents{Entries: []Document{{ID: "child1"}, {ID: "child2"}}},
			mockStatus: 200,
			wantDocIDs: []string{"child1", "child2"},
		},
		{
			name:       "not found",
			parentPath: "/missing",
			mockResp:   &NuxeoError{Message: "not found"},
			mockStatus: 404,
			wantErr:    true,
		},
		{
			name:       "http error",
			parentPath: "/err",
			mockErr:    errors.New("network error"),
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				body, _ := json.Marshal(tt.mockResp)
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			})
			got, err := repo.FetchChildrenByPath(context.Background(), tt.parentPath, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchChildrenByPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				var gotIDs []string
				for _, doc := range got.Entries {
					gotIDs = append(gotIDs, doc.ID)
				}
				if len(gotIDs) != len(tt.wantDocIDs) {
					t.Errorf("FetchChildrenByPath() got %d entries, want %d", len(gotIDs), len(tt.wantDocIDs))
				}
				for i, wantID := range tt.wantDocIDs {
					if gotIDs[i] != wantID {
						t.Errorf("FetchChildrenByPath() entry[%d].ID = %v, want %v", i, gotIDs[i], wantID)
					}
				}
			}
		})
	}
}

func TestRepository_FetchChildrenById(t *testing.T) {
	tests := []struct {
		name       string
		parentId   string
		mockResp   any
		mockStatus int
		mockErr    error
		wantErr    bool
		wantDocIDs []string
	}{
		{
			name:       "success",
			parentId:   "parent123",
			mockResp:   &Documents{Entries: []Document{{ID: "childA"}, {ID: "childB"}}},
			mockStatus: 200,
			wantDocIDs: []string{"childA", "childB"},
		},
		{
			name:       "not found",
			parentId:   "missing",
			mockResp:   &NuxeoError{Message: "not found"},
			mockStatus: 404,
			wantErr:    true,
		},
		{
			name:     "http error",
			parentId: "err",
			mockErr:  errors.New("network error"),
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				body, _ := json.Marshal(tt.mockResp)
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			})
			got, err := repo.FetchChildrenById(context.Background(), tt.parentId, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchChildrenById() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				var gotIDs []string
				for _, doc := range got.Entries {
					gotIDs = append(gotIDs, doc.ID)
				}
				if len(gotIDs) != len(tt.wantDocIDs) {
					t.Errorf("FetchChildrenById() got %d entries, want %d", len(gotIDs), len(tt.wantDocIDs))
				}
				for i, wantID := range tt.wantDocIDs {
					if gotIDs[i] != wantID {
						t.Errorf("FetchChildrenById() entry[%d].ID = %v, want %v", i, gotIDs[i], wantID)
					}
				}
			}
		})
	}
}

func TestRepository_StreamBlobByPath(t *testing.T) {
	tests := []struct {
		name        string
		docPath     string
		blobXPath   string
		mockBody    []byte
		mockStatus  int
		mockHeaders map[string]string
		mockErr     error
		wantErr     bool
		wantName    string
		wantMime    string
		wantLength  string
	}{
		{
			name:       "success",
			docPath:    "/default-domain/workspaces",
			blobXPath:  "file:content",
			mockBody:   []byte("blobdata"),
			mockStatus: 200,
			mockHeaders: map[string]string{
				"Content-Disposition": `attachment; filename="testfile.txt"`,
				"Content-Type":        "text/plain",
				"Content-Length":      "8",
			},
			wantName:   "testfile.txt",
			wantMime:   "text/plain",
			wantLength: "8",
		},
		{
			name:      "http error",
			docPath:   "/err",
			blobXPath: "file:content",
			mockErr:   errors.New("network error"),
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				h := http.Header{}
				for k, v := range tt.mockHeaders {
					h.Set(k, v)
				}
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewReader(tt.mockBody)),
					Header:     h,
				}, nil
			})
			blob, err := repo.StreamBlobByPath(context.Background(), tt.docPath, tt.blobXPath, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("StreamBlobByPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				if blob.Filename != tt.wantName {
					t.Errorf("StreamBlobByPath() Filename = %v, want %v", blob.Filename, tt.wantName)
				}
				if blob.MimeType != tt.wantMime {
					t.Errorf("StreamBlobByPath() MimeType = %v, want %v", blob.MimeType, tt.wantMime)
				}
				if blob.Length != tt.wantLength {
					t.Errorf("StreamBlobByPath() Length = %v, want %v", blob.Length, tt.wantLength)
				}
				b, _ := io.ReadAll(blob)
				if !bytes.Equal(b, tt.mockBody) {
					t.Errorf("StreamBlobByPath() stream body = %v, want %v", b, tt.mockBody)
				}
			}
		})
	}
}

func TestRepository_StreamBlobById(t *testing.T) {
	tests := []struct {
		name        string
		docID       string
		blobXPath   string
		mockBody    []byte
		mockStatus  int
		mockHeaders map[string]string
		mockErr     error
		wantErr     bool
		wantName    string
		wantMime    string
		wantLength  string
	}{
		{
			name:       "success",
			docID:      "doc123",
			blobXPath:  "file:content",
			mockBody:   []byte("blobdata"),
			mockStatus: 200,
			mockHeaders: map[string]string{
				"Content-Disposition": `attachment; filename="testfile.txt"`,
				"Content-Type":        "text/plain",
				"Content-Length":      "8",
			},
			wantName:   "testfile.txt",
			wantMime:   "text/plain",
			wantLength: "8",
		},
		{
			name:      "http error",
			docID:     "err",
			blobXPath: "file:content",
			mockErr:   errors.New("network error"),
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				h := http.Header{}
				for k, v := range tt.mockHeaders {
					h.Set(k, v)
				}
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewReader(tt.mockBody)),
					Header:     h,
				}, nil
			})
			blob, err := repo.StreamBlobById(context.Background(), tt.docID, tt.blobXPath, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("StreamBlobById() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				if blob.Filename != tt.wantName {
					t.Errorf("StreamBlobById() Filename = %v, want %v", blob.Filename, tt.wantName)
				}
				if blob.MimeType != tt.wantMime {
					t.Errorf("StreamBlobById() MimeType = %v, want %v", blob.MimeType, tt.wantMime)
				}
				if blob.Length != tt.wantLength {
					t.Errorf("StreamBlobById() Length = %v, want %v", blob.Length, tt.wantLength)
				}
				b, _ := io.ReadAll(blob)
				if !bytes.Equal(b, tt.mockBody) {
					t.Errorf("StreamBlobById() stream body = %v, want %v", b, tt.mockBody)
				}
			}
		})
	}
}

func TestRepository_StartWorkflowInstanceWithDocId(t *testing.T) {
	tests := []struct {
		name       string
		docID      string
		workflow   Workflow
		mockResp   any
		mockStatus int
		mockErr    error
		wantErr    bool
		wantID     string
	}{
		{
			name:       "success",
			docID:      "doc123",
			workflow:   Workflow{Id: "wf1"},
			mockResp:   &Workflow{Id: "wf1"},
			mockStatus: 201,
			wantID:     "wf1",
		},
		{
			name:       "bad request",
			docID:      "doc123",
			workflow:   Workflow{Id: "wf1"},
			mockResp:   &NuxeoError{Message: "bad request"},
			mockStatus: 400,
			wantErr:    true,
		},
		{
			name:     "http error",
			docID:    "doc123",
			workflow: Workflow{Id: "wf1"},
			mockErr:  errors.New("network error"),
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				body, _ := json.Marshal(tt.mockResp)
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			})
			got, err := repo.StartWorkflowInstanceWithDocId(context.Background(), tt.docID, tt.workflow, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("StartWorkflowInstanceWithDocId() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got.Id != tt.wantID {
				t.Errorf("StartWorkflowInstanceWithDocId() got.Id = %v, want %v", got.Id, tt.wantID)
			}
		})
	}
}

func TestRepository_StartWorkflowInstanceWithDocPath(t *testing.T) {
	tests := []struct {
		name       string
		docPath    string
		workflow   Workflow
		mockResp   any
		mockStatus int
		mockErr    error
		wantErr    bool
		wantID     string
	}{
		{
			name:       "success",
			docPath:    "/default-domain/workspaces",
			workflow:   Workflow{Id: "wf2"},
			mockResp:   &Workflow{Id: "wf2"},
			mockStatus: 201,
			wantID:     "wf2",
		},
		{
			name:       "bad request",
			docPath:    "/default-domain/workspaces",
			workflow:   Workflow{Id: "wf2"},
			mockResp:   &NuxeoError{Message: "bad request"},
			mockStatus: 400,
			wantErr:    true,
		},
		{
			name:     "http error",
			docPath:  "/default-domain/workspaces",
			workflow: Workflow{Id: "wf2"},
			mockErr:  errors.New("network error"),
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				body, _ := json.Marshal(tt.mockResp)
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			})
			got, err := repo.StartWorkflowInstanceWithDocPath(context.Background(), tt.docPath, tt.workflow, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("StartWorkflowInstanceWithDocPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got.Id != tt.wantID {
				t.Errorf("StartWorkflowInstanceWithDocPath() got.Id = %v, want %v", got.Id, tt.wantID)
			}
		})
	}
}

func TestRepository_FetchWorkflowInstancesByDocId(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		wfs := Workflows{Entries: []Workflow{{Id: "wf1"}, {Id: "wf2"}}}
		repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&wfs)
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		})
		got, err := repo.FetchWorkflowInstancesByDocId(context.Background(), "doc123", nil)
		if err != nil {
			t.Fatalf("FetchWorkflowInstancesByDocId() error = %v, want nil", err)
		}
		if len(got.Entries) != 2 {
			t.Errorf("FetchWorkflowInstancesByDocId() got %d entries, want 2", len(got.Entries))
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&NuxeoError{Message: "not found"})
			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		})
		_, err := repo.FetchWorkflowInstancesByDocId(context.Background(), "ghost", nil)
		if err == nil {
			t.Errorf("FetchWorkflowInstancesByDocId() error = nil, want error")
		}
	})
}

func TestRepository_FetchWorkflowInstancesByDocPath(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		wfs := Workflows{Entries: []Workflow{{Id: "wfA"}, {Id: "wfB"}}}
		repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&wfs)
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		})
		got, err := repo.FetchWorkflowInstancesByDocPath(context.Background(), "/default-domain/workspaces", nil)
		if err != nil {
			t.Fatalf("FetchWorkflowInstancesByDocPath() error = %v, want nil", err)
		}
		if len(got.Entries) != 2 {
			t.Errorf("FetchWorkflowInstancesByDocPath() got %d entries, want 2", len(got.Entries))
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&NuxeoError{Message: "not found"})
			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		})
		_, err := repo.FetchWorkflowInstancesByDocPath(context.Background(), "/ghost", nil)
		if err == nil {
			t.Errorf("FetchWorkflowInstancesByDocPath() error = nil, want error")
		}
	})
}

func TestRepository_FetchWorkflowInstance(t *testing.T) {
	tests := []struct {
		name       string
		wfID       string
		mockResp   any
		mockStatus int
		mockErr    error
		wantErr    bool
		wantID     string
	}{
		{
			name:       "success",
			wfID:       "wf123",
			mockResp:   &Workflow{Id: "wf123"},
			mockStatus: 200,
			wantID:     "wf123",
		},
		{
			name:       "not found",
			wfID:       "missing",
			mockResp:   &NuxeoError{Message: "not found"},
			mockStatus: 404,
			wantErr:    true,
		},
		{
			name:    "http error",
			wfID:    "err",
			mockErr: errors.New("network error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				body, _ := json.Marshal(tt.mockResp)
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			})
			got, err := repo.FetchWorkflowInstance(context.Background(), tt.wfID, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchWorkflowInstance() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got.Id != tt.wantID {
				t.Errorf("FetchWorkflowInstance() got.Id = %v, want %v", got.Id, tt.wantID)
			}
		})
	}
}

func TestRepository_CancelWorkflowInstance(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 204,
				Body:       io.NopCloser(bytes.NewReader([]byte{})),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		})
		err := repo.CancelWorkflowInstance(context.Background(), "wf123")
		if err != nil {
			t.Errorf("CancelWorkflowInstance() error = %v, want nil", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&NuxeoError{Message: "not found"})
			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		})
		err := repo.CancelWorkflowInstance(context.Background(), "missing")
		if err == nil {
			t.Errorf("CancelWorkflowInstance() error = nil, want error")
		}
	})
}

func TestRepository_FetchWorkflowInstanceGraph(t *testing.T) {
	fieldId := Field("graph1")
	tests := []struct {
		name        string
		wfID        string
		mockResp    any
		mockStatus  int
		mockErr     error
		wantErr     bool
		wantGraphID string
	}{
		{
			name: "success",
			wfID: "wf123",
			mockResp: &WorkflowGraph{Nodes: map[string]Field{
				"Id": fieldId,
			}},
			mockStatus:  200,
			wantGraphID: "graph1",
		},
		{
			name:       "not found",
			wfID:       "missing",
			mockResp:   &NuxeoError{Message: "not found"},
			mockStatus: 404,
			wantErr:    true,
		},
		{
			name:    "http error",
			wfID:    "err",
			mockErr: errors.New("network error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				body, _ := json.Marshal(tt.mockResp)
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			})
			got, err := repo.FetchWorkflowInstanceGraph(context.Background(), tt.wfID, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchWorkflowInstanceGraph() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				if field, ok := got.Nodes["Id"]; ok {
					if fieldValue, _ := field.String(); fieldValue != nil && *fieldValue != tt.wantGraphID {
						t.Errorf("FetchWorkflowInstanceGraph() got.Id = %v, want %v", *fieldValue, tt.wantGraphID)
					}
				}
			}
		})
	}
}

func TestRepository_FetchWorkflowModel(t *testing.T) {
	tests := []struct {
		name       string
		modelName  string
		mockResp   any
		mockStatus int
		mockErr    error
		wantErr    bool
		wantID     string
	}{
		{
			name:       "success",
			modelName:  "modelA",
			mockResp:   &Workflow{Id: "modelA"},
			mockStatus: 200,
			wantID:     "modelA",
		},
		{
			name:       "not found",
			modelName:  "missing",
			mockResp:   &NuxeoError{Message: "not found"},
			mockStatus: 404,
			wantErr:    true,
		},
		{
			name:      "http error",
			modelName: "err",
			mockErr:   errors.New("network error"),
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				body, _ := json.Marshal(tt.mockResp)
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			})
			got, err := repo.FetchWorkflowModel(context.Background(), tt.modelName, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchWorkflowModel() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got.Id != tt.wantID {
				t.Errorf("FetchWorkflowModel() got.Id = %v, want %v", got.Id, tt.wantID)
			}
		})
	}
}

func TestRepository_FetchWorkflowModelGraph(t *testing.T) {
	tests := []struct {
		name        string
		modelName   string
		mockResp    any
		mockStatus  int
		mockErr     error
		wantErr     bool
		wantGraphID string
	}{
		{
			name:      "success",
			modelName: "modelA",
			mockResp: &WorkflowGraph{Nodes: map[string]Field{
				"Id": NewStringField("graphA"),
			}},
			mockStatus:  200,
			wantGraphID: "graphA",
		},
		{
			name:       "not found",
			modelName:  "missing",
			mockResp:   &NuxeoError{Message: "not found"},
			mockStatus: 404,
			wantErr:    true,
		},
		{
			name:      "http error",
			modelName: "err",
			mockErr:   errors.New("network error"),
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				body, _ := json.Marshal(tt.mockResp)
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			})
			got, err := repo.FetchWorkflowModelGraph(context.Background(), tt.modelName, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchWorkflowModelGraph() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				if field, ok := got.Nodes["Id"]; ok {
					if fieldValue, _ := field.String(); fieldValue != nil && *fieldValue != tt.wantGraphID {
						t.Errorf("FetchWorkflowModelGraph() got.Id = %v, want %v", *fieldValue, tt.wantGraphID)
					}
				}
			}
		})
	}
}

func TestRepository_FetchWorkflowModels(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		models := Workflows{Entries: []Workflow{{Id: "modelA"}, {Id: "modelB"}}}
		repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&models)
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		})
		got, err := repo.FetchWorkflowModels(context.Background(), nil)
		if err != nil {
			t.Fatalf("FetchWorkflowModels() error = %v, want nil", err)
		}
		if len(got.Entries) != 2 {
			t.Errorf("FetchWorkflowModels() got %d entries, want 2", len(got.Entries))
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&NuxeoError{Message: "not found"})
			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		})
		_, err := repo.FetchWorkflowModels(context.Background(), nil)
		if err == nil {
			t.Errorf("FetchWorkflowModels() error = nil, want error")
		}
	})
}

func TestRepository_CreateForAdapter(t *testing.T) {
	tests := []struct {
		name        string
		docID       string
		adapter     string
		pathSuffix  string
		queryParams []string
		payload     any
		mockStatus  int
		mockBody    []byte
		mockErr     error
		wantErr     bool
		wantStatus  int
	}{
		{
			name:        "success",
			docID:       "doc123",
			adapter:     "custom",
			pathSuffix:  "action",
			queryParams: []string{"foo=bar"},
			payload:     map[string]string{"key": "value"},
			mockStatus:  201,
			mockBody:    []byte(`{"result":"ok"}`),
			wantStatus:  201,
		},
		{
			name:    "http error",
			docID:   "doc123",
			adapter: "custom",
			mockErr: errors.New("network error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewReader(tt.mockBody)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			})
			resp, err := repo.CreateForAdapter(context.Background(), tt.docID, tt.adapter, tt.pathSuffix, tt.queryParams, tt.payload, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateForAdapter() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && resp.StatusCode != tt.wantStatus {
				t.Errorf("CreateForAdapter() StatusCode = %v, want %v", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}

func TestRepository_FetchForAdapter(t *testing.T) {
	tests := []struct {
		name        string
		docID       string
		adapter     string
		pathSuffix  string
		queryParams []string
		mockStatus  int
		mockBody    []byte
		mockErr     error
		wantErr     bool
		wantStatus  int
	}{
		{
			name:        "success",
			docID:       "doc123",
			adapter:     "custom",
			pathSuffix:  "info",
			queryParams: []string{"foo=bar"},
			mockStatus:  200,
			mockBody:    []byte(`{"result":"ok"}`),
			wantStatus:  200,
		},
		{
			name:    "http error",
			docID:   "doc123",
			adapter: "custom",
			mockErr: errors.New("network error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewReader(tt.mockBody)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			})
			resp, err := repo.FetchForAdapter(context.Background(), tt.docID, tt.adapter, tt.pathSuffix, tt.queryParams, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchForAdapter() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && resp.StatusCode != tt.wantStatus {
				t.Errorf("FetchForAdapter() StatusCode = %v, want %v", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}

func TestRepository_UpdateForAdapter(t *testing.T) {
	tests := []struct {
		name        string
		docID       string
		adapter     string
		pathSuffix  string
		queryParams []string
		payload     any
		mockStatus  int
		mockBody    []byte
		mockErr     error
		wantErr     bool
		wantStatus  int
	}{
		{
			name:        "success",
			docID:       "doc123",
			adapter:     "custom",
			pathSuffix:  "update",
			queryParams: []string{"foo=bar"},
			payload:     map[string]string{"update": "yes"},
			mockStatus:  200,
			mockBody:    []byte(`{"result":"updated"}`),
			wantStatus:  200,
		},
		{
			name:    "http error",
			docID:   "doc123",
			adapter: "custom",
			mockErr: errors.New("network error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewReader(tt.mockBody)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			})
			resp, err := repo.UpdateForAdapter(context.Background(), tt.docID, tt.adapter, tt.pathSuffix, tt.queryParams, tt.payload, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateForAdapter() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && resp.StatusCode != tt.wantStatus {
				t.Errorf("UpdateForAdapter() StatusCode = %v, want %v", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}

func TestRepository_DeleteForAdapter(t *testing.T) {
	tests := []struct {
		name        string
		docID       string
		adapter     string
		pathSuffix  string
		queryParams []string
		mockStatus  int
		mockBody    []byte
		mockErr     error
		wantErr     bool
		wantStatus  int
	}{
		{
			name:        "success",
			docID:       "doc123",
			adapter:     "custom",
			pathSuffix:  "delete",
			queryParams: []string{"foo=bar"},
			mockStatus:  204,
			mockBody:    []byte{},
			wantStatus:  204,
		},
		{
			name:    "http error",
			docID:   "doc123",
			adapter: "custom",
			mockErr: errors.New("network error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newTestRepository(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewReader(tt.mockBody)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			})
			resp, err := repo.DeleteForAdapter(context.Background(), tt.docID, tt.adapter, tt.pathSuffix, tt.queryParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteForAdapter() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && resp.StatusCode != tt.wantStatus {
				t.Errorf("DeleteForAdapter() StatusCode = %v, want %v", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}
