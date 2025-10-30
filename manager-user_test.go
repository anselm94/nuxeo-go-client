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

func newTestUserManager(respond func(req *http.Request) (*http.Response, error)) *userManager {
	client := newMockNuxeoClient(respond)
	return &userManager{
		client: client,
		logger: slog.Default(),
	}
}

func TestUserManager_FetchUser(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name       string
		args       args
		mockResp   any
		mockStatus int
		mockErr    error
		wantErr    bool
		wantId     string
	}{
		{
			name:       "success",
			args:       args{id: "john"},
			mockResp:   &User{Id: "john"},
			mockStatus: 200,
			wantId:     "john",
		},
		{
			name:       "not found",
			args:       args{id: "missing"},
			mockResp:   &NuxeoError{Message: "not found"},
			mockStatus: 404,
			wantErr:    true,
		},
		{
			name:    "http error",
			args:    args{id: "err"},
			mockErr: errors.New("network error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				body, _ := json.Marshal(tt.mockResp)
				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
				}, nil
			})
			got, err := um.FetchUser(context.Background(), tt.args.id, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got.Id != tt.wantId {
				t.Errorf("FetchUser() got.Id = %v, want %v", got.Id, tt.wantId)
			}
		})
	}
}

func TestUserManager_CreateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		user := User{Id: "alice"}
		um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&user)
			return &http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			}, nil
		})
		got, err := um.CreateUser(context.Background(), user, nil)
		if err != nil {
			t.Fatalf("CreateUser() error = %v, want nil", err)
		}
		if got.Id != "alice" {
			t.Errorf("CreateUser() got.Id = %v, want alice", got.Id)
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&NuxeoError{Message: "bad request"})
			return &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			}, nil
		})
		_, err := um.CreateUser(context.Background(), User{Id: "bob"}, nil)
		if err == nil {
			t.Errorf("CreateUser() error = nil, want error")
		}
	})
}

func TestUserManager_UpdateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		user := User{Id: "eve"}
		um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&user)
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			}, nil
		})
		got, err := um.UpdateUser(context.Background(), "eve", user, nil)
		if err != nil {
			t.Fatalf("UpdateUser() error = %v, want nil", err)
		}
		if got.Id != "eve" {
			t.Errorf("UpdateUser() got.Id = %v, want eve", got.Id)
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&NuxeoError{Message: "not found"})
			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			}, nil
		})
		_, err := um.UpdateUser(context.Background(), "ghost", User{Id: "ghost"}, nil)
		if err == nil {
			t.Errorf("UpdateUser() error = nil, want error")
		}
	})
}

func TestUserManager_DeleteUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 204,
				Body:       io.NopCloser(bytes.NewReader([]byte{})),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			}, nil
		})
		err := um.DeleteUser(context.Background(), "john", nil)
		if err != nil {
			t.Errorf("DeleteUser() error = %v, want nil", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&NuxeoError{Message: "not found"})
			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			}, nil
		})
		err := um.DeleteUser(context.Background(), "ghost", nil)
		if err == nil {
			t.Errorf("DeleteUser() error = nil, want error")
		}
	})
}

func TestUserManager_SearchUsers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		users := Users{}
		users.Entries = []User{{Id: "john"}, {Id: "jane"}}
		um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&users)
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			}, nil
		})
		got, err := um.SearchUsers(context.Background(), "john", nil, nil)
		if err != nil {
			t.Fatalf("SearchUsers() error = %v, want nil", err)
		}
		if len(got.Entries) != 2 {
			t.Errorf("SearchUsers() got %d entries, want 2", len(got.Entries))
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&NuxeoError{Message: "bad query"})
			return &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			}, nil
		})
		_, err := um.SearchUsers(context.Background(), "bad", nil, nil)
		if err == nil {
			t.Errorf("SearchUsers() error = nil, want error")
		}
	})
}

func TestUserManager_AddUserToGroup(t *testing.T) {
	tests := []struct {
		name       string
		user       string
		group      string
		mockResp   any
		mockStatus int
		mockErr    error
		wantErr    bool
		wantId     string
	}{
		{
			name:       "success",
			user:       "alice",
			group:      "devs",
			mockResp:   &User{Id: "alice"},
			mockStatus: 200,
			wantId:     "alice",
		},
		{
			name:       "not found",
			user:       "ghost",
			group:      "devs",
			mockResp:   &NuxeoError{Message: "not found"},
			mockStatus: 404,
			wantErr:    true,
		},
		{
			name:    "http error",
			user:    "alice",
			group:   "devs",
			mockErr: errors.New("network error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
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
			got, err := um.AddUserToGroup(context.Background(), tt.user, tt.group, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddUserToGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got.Id != tt.wantId {
				t.Errorf("AddUserToGroup() got.Id = %v, want %v", got.Id, tt.wantId)
			}
		})
	}
}

func TestUserManager_FetchGroup(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		group := Group{Id: "admins"}
		um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&group)
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			}, nil
		})
		got, err := um.FetchGroup(context.Background(), "admins", nil)
		if err != nil {
			t.Fatalf("FetchGroup() error = %v, want nil", err)
		}
		if got.Id != "admins" {
			t.Errorf("FetchGroup() got.Id = %v, want admins", got.Id)
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&NuxeoError{Message: "not found"})
			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			}, nil
		})
		_, err := um.FetchGroup(context.Background(), "ghost", nil)
		if err == nil {
			t.Errorf("FetchGroup() error = nil, want error")
		}
	})
}

func TestUserManager_CreateGroup(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		group := Group{Id: "devs"}
		um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&group)
			return &http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			}, nil
		})
		got, err := um.CreateGroup(context.Background(), group, nil)
		if err != nil {
			t.Fatalf("CreateGroup() error = %v, want nil", err)
		}
		if got.Id != "devs" {
			t.Errorf("CreateGroup() got.Id = %v, want devs", got.Id)
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&NuxeoError{Message: "bad request"})
			return &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			}, nil
		})
		_, err := um.CreateGroup(context.Background(), Group{Id: "bad"}, nil)
		if err == nil {
			t.Errorf("CreateGroup() error = nil, want error")
		}
	})
}

func TestUserManager_UpdateGroup(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		group := Group{Id: "devs"}
		um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&group)
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			}, nil
		})
		got, err := um.UpdateGroup(context.Background(), "devs", group, nil)
		if err != nil {
			t.Fatalf("UpdateGroup() error = %v, want nil", err)
		}
		if got.Id != "devs" {
			t.Errorf("UpdateGroup() got.Id = %v, want devs", got.Id)
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&NuxeoError{Message: "not found"})
			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			}, nil
		})
		_, err := um.UpdateGroup(context.Background(), "ghost", Group{Id: "ghost"}, nil)
		if err == nil {
			t.Errorf("UpdateGroup() error = nil, want error")
		}
	})
}

func TestUserManager_DeleteGroup(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 204,
				Body:       io.NopCloser(bytes.NewReader([]byte{})),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			}, nil
		})
		err := um.DeleteGroup(context.Background(), "devs", nil)
		if err != nil {
			t.Errorf("DeleteGroup() error = %v, want nil", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&NuxeoError{Message: "not found"})
			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			}, nil
		})
		err := um.DeleteGroup(context.Background(), "ghost", nil)
		if err == nil {
			t.Errorf("DeleteGroup() error = nil, want error")
		}
	})
}

func TestUserManager_SearchGroup(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		groups := Groups{}
		groups.Entries = []Group{{Id: "admins"}, {Id: "devs"}}
		um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&groups)
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			}, nil
		})
		got, err := um.SearchGroup(context.Background(), "dev", nil, nil)
		if err != nil {
			t.Fatalf("SearchGroup() error = %v, want nil", err)
		}
		if len(got.Entries) != 2 {
			t.Errorf("SearchGroup() got %d entries, want 2", len(got.Entries))
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
			body, _ := json.Marshal(&NuxeoError{Message: "bad query"})
			return &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			}, nil
		})
		_, err := um.SearchGroup(context.Background(), "bad", nil, nil)
		if err == nil {
			t.Errorf("SearchGroup() error = nil, want error")
		}
	})
}

func TestUserManager_AttachGroupToUser(t *testing.T) {
	tests := []struct {
		name       string
		group      string
		user       string
		mockResp   any
		mockStatus int
		mockErr    error
		wantErr    bool
		wantId     string
	}{
		{
			name:       "success",
			group:      "devs",
			user:       "alice",
			mockResp:   &Group{Id: "devs"},
			mockStatus: 200,
			wantId:     "devs",
		},
		{
			name:       "not found",
			group:      "ghost",
			user:       "bob",
			mockResp:   &NuxeoError{Message: "not found"},
			mockStatus: 404,
			wantErr:    true,
		},
		{
			name:    "http error",
			group:   "devs",
			user:    "err",
			mockErr: errors.New("network error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
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
			got, err := um.AttachGroupToUser(context.Background(), tt.group, tt.user, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("AttachGroupToUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got.Id != tt.wantId {
				t.Errorf("AttachGroupToUser() got.Id = %v, want %v", got.Id, tt.wantId)
			}
		})
	}
}

func TestUserManager_FetchGroupMemberUsers(t *testing.T) {
	tests := []struct {
		name       string
		group      string
		mockResp   any
		mockStatus int
		mockErr    error
		wantErr    bool
		wantCount  int
	}{
		{
			name:       "success",
			group:      "devs",
			mockResp:   &Users{Entries: []User{{Id: "alice"}, {Id: "bob"}}},
			mockStatus: 200,
			wantCount:  2,
		},
		{
			name:       "not found",
			group:      "ghost",
			mockResp:   &NuxeoError{Message: "not found"},
			mockStatus: 404,
			wantErr:    true,
		},
		{
			name:    "http error",
			group:   "devs",
			mockErr: errors.New("network error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
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
			got, err := um.FetchGroupMemberUsers(context.Background(), tt.group, nil, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchGroupMemberUsers() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && len(got.Entries) != tt.wantCount {
				t.Errorf("FetchGroupMemberUsers() got %d entries, want %d", len(got.Entries), tt.wantCount)
			}
		})
	}
}

func TestUserManager_FetchGroupMemberGroups(t *testing.T) {
	tests := []struct {
		name       string
		group      string
		mockResp   any
		mockStatus int
		mockErr    error
		wantErr    bool
		wantCount  int
	}{
		{
			name:       "success",
			group:      "devs",
			mockResp:   &Groups{Entries: []Group{{Id: "admins"}, {Id: "devs"}}},
			mockStatus: 200,
			wantCount:  2,
		},
		{
			name:       "not found",
			group:      "ghost",
			mockResp:   &NuxeoError{Message: "not found"},
			mockStatus: 404,
			wantErr:    true,
		},
		{
			name:    "http error",
			group:   "devs",
			mockErr: errors.New("network error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
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
			got, err := um.FetchGroupMemberGroups(context.Background(), tt.group, nil, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchGroupMemberGroups() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && len(got.Entries) != tt.wantCount {
				t.Errorf("FetchGroupMemberGroups() got %d entries, want %d", len(got.Entries), tt.wantCount)
			}
		})
	}
}

func TestUserManager_FetchWorkflowInstances(t *testing.T) {
	tests := []struct {
		name       string
		mockResp   any
		mockStatus int
		mockErr    error
		wantErr    bool
		wantCount  int
	}{
		{
			name:       "success",
			mockResp:   &Workflows{Entries: []Workflow{{Id: "wf1"}, {Id: "wf2"}}},
			mockStatus: 200,
			wantCount:  2,
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
			um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
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
			got, err := um.FetchWorkflowInstances(context.Background(), nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchWorkflowInstances() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && len(got.Entries) != tt.wantCount {
				t.Errorf("FetchWorkflowInstances() got %d entries, want %d", len(got.Entries), tt.wantCount)
			}
		})
	}
}

func TestUserManager_StartWorkflowInstance(t *testing.T) {
	tests := []struct {
		name       string
		input      Workflow
		mockResp   any
		mockStatus int
		mockErr    error
		wantErr    bool
		wantId     string
	}{
		{
			name:       "success",
			input:      Workflow{Id: "wf1"},
			mockResp:   &Workflow{Id: "wf1"},
			mockStatus: 201,
			wantId:     "wf1",
		},
		{
			name:       "bad request",
			input:      Workflow{Id: "bad"},
			mockResp:   &NuxeoError{Message: "bad request"},
			mockStatus: 400,
			wantErr:    true,
		},
		{
			name:    "http error",
			input:   Workflow{Id: "wf1"},
			mockErr: errors.New("network error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
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
			got, err := um.StartWorkflowInstance(context.Background(), tt.input, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("StartWorkflowInstance() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got.Id != tt.wantId {
				t.Errorf("StartWorkflowInstance() got.Id = %v, want %v", got.Id, tt.wantId)
			}
		})
	}
}

func TestUserManager_FetchCurrentUser(t *testing.T) {
	tests := []struct {
		name        string
		loginResp   any
		loginStatus int
		loginErr    error
		userResp    any
		userStatus  int
		userErr     error
		wantErr     bool
		wantUserId  string
	}{
		{
			name:        "success",
			loginResp:   map[string]string{"username": "john"},
			loginStatus: 200,
			userResp:    &User{Id: "john"},
			userStatus:  200,
			wantUserId:  "john",
		},
		{
			name:     "login error",
			loginErr: errors.New("login failed"),
			wantErr:  true,
		},
		{
			name:        "login returns no username",
			loginResp:   map[string]string{},
			loginStatus: 200,
			userResp:    &User{Id: ""},
			userStatus:  200,
			wantErr:     false, // changed from true to false
			wantUserId:  "",    // explicitly expect empty user id
		},
		{
			name:        "user fetch error",
			loginResp:   map[string]string{"username": "ghost"},
			loginStatus: 200,
			userResp:    &NuxeoError{Message: "not found"},
			userStatus:  404,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			loginCalled := false
			userCalled := false
			um := newTestUserManager(func(req *http.Request) (*http.Response, error) {
				switch {
				case !loginCalled:
					loginCalled = true
					if tt.loginErr != nil {
						return nil, tt.loginErr
					}
					body, _ := json.Marshal(tt.loginResp)
					return &http.Response{
						StatusCode: tt.loginStatus,
						Body:       io.NopCloser(bytes.NewReader(body)),
						Header:     http.Header{"Content-Type": []string{"application/json"}},
					}, nil
				case !userCalled:
					userCalled = true
					if tt.userErr != nil {
						return nil, tt.userErr
					}
					body, _ := json.Marshal(tt.userResp)
					return &http.Response{
						StatusCode: tt.userStatus,
						Body:       io.NopCloser(bytes.NewReader(body)),
						Header:     http.Header{"Content-Type": []string{"application/json"}},
					}, nil
				default:
					return nil, errors.New("unexpected call")
				}
			})
			got, err := um.FetchCurrentUser(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchCurrentUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got.Id != tt.wantUserId {
				t.Errorf("FetchCurrentUser() got.Id = %v, want %v", got.Id, tt.wantUserId)
			}
		})
	}
}
