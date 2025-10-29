package nuxeo

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"
)

func TestDataModelManager_FetchTypes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		respond  func(req *http.Request) (*http.Response, error)
		wantErr  bool
		wantType string
	}{
		{
			name: "success",
			respond: func(req *http.Request) (*http.Response, error) {
				body := `{
					"docTypes": {
						"File": {
							"parent": "Document",
							"facets": ["Versionable"],
							"schemas": ["file", "common"]
						}
					},
					"schemas": {
						"file": {
							"@prefix": {"DataType": "file"},
							"content": "blob"
						},
						"common": {
							"@prefix": {"DataType": "common"},
							"title": "string"
						}
					}
				}`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
			wantErr:  false,
			wantType: "File",
		},
		{
			name: "client error",
			respond: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("client error")
			},
			wantErr:  true,
			wantType: "",
		},
		{
			name: "invalid json",
			respond: func(req *http.Request) (*http.Response, error) {
				body := `invalid`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
			wantErr:  true,
			wantType: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newMockNuxeoClient(tc.respond)
			dmm := &dataModelManager{client: client, logger: slog.Default()}
			types, err := dmm.FetchTypes(context.Background())
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tc.wantErr && types != nil {
				if _, ok := types.DocTypes[tc.wantType]; !ok {
					t.Errorf("expected docType %q in result, got %+v", tc.wantType, types.DocTypes)
				}
			}
		})
	}
}

func TestDataModelManager_FetchType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		respond  func(req *http.Request) (*http.Response, error)
		wantErr  bool
		wantName string
	}{
		{
			name: "success",
			respond: func(req *http.Request) (*http.Response, error) {
				body := `{
					"name": "File",
					"parent": "Document",
					"facets": ["Versionable"],
					"schemas": []
				}`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
			wantErr:  false,
			wantName: "File",
		},
		{
			name: "client error",
			respond: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("client error")
			},
			wantErr:  true,
			wantName: "",
		},
		{
			name: "invalid json",
			respond: func(req *http.Request) (*http.Response, error) {
				body := `invalid`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
			wantErr:  true,
			wantName: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newMockNuxeoClient(tc.respond)
			dmm := &dataModelManager{client: client, logger: slog.Default()}
			docType, err := dmm.FetchType(context.Background(), "File")
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tc.wantErr && docType != nil {
				if docType.Name != tc.wantName {
					t.Errorf("unexpected docType name: got %q, want %q", docType.Name, tc.wantName)
				}
			}
		})
	}
}

func TestDataModelManager_FetchSchemas(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		respond   func(req *http.Request) (*http.Response, error)
		wantErr   bool
		wantCount int
	}{
		{
			name: "success",
			respond: func(req *http.Request) (*http.Response, error) {
				body := `[{"name":"file","prefix":"file","fields":{"content":"blob"}},{"name":"common","prefix":"common","fields":{"title":"string"}}]`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name: "client error",
			respond: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("client error")
			},
			wantErr:   true,
			wantCount: 0,
		},
		{
			name: "invalid json",
			respond: func(req *http.Request) (*http.Response, error) {
				body := `invalid`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
			wantErr:   true,
			wantCount: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newMockNuxeoClient(tc.respond)
			dmm := &dataModelManager{client: client, logger: slog.Default()}
			schemas, err := dmm.FetchSchemas(context.Background())
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tc.wantErr && schemas != nil {
				if len(*schemas) != tc.wantCount {
					t.Errorf("unexpected schema count: got %d, want %d", len(*schemas), tc.wantCount)
				}
			}
		})
	}
}

func TestDataModelManager_FetchSchema(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		respond   func(req *http.Request) (*http.Response, error)
		wantErr   bool
		wantName  string
		wantField string
	}{
		{
			name: "success",
			respond: func(req *http.Request) (*http.Response, error) {
				body := `{"name":"file","prefix":"file","fields":{"content":"blob"}}`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
			wantErr:   false,
			wantName:  "file",
			wantField: "content",
		},
		{
			name: "client error",
			respond: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("client error")
			},
			wantErr:   true,
			wantName:  "",
			wantField: "",
		},
		{
			name: "invalid json",
			respond: func(req *http.Request) (*http.Response, error) {
				body := `invalid`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
			wantErr:   true,
			wantName:  "",
			wantField: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newMockNuxeoClient(tc.respond)
			dmm := &dataModelManager{client: client, logger: slog.Default()}
			schema, err := dmm.FetchSchema(context.Background(), "file")
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tc.wantErr && schema != nil {
				if schema.Name != tc.wantName {
					t.Errorf("unexpected schema name: got %q, want %q", schema.Name, tc.wantName)
				}
				if _, ok := schema.Fields[tc.wantField]; !ok {
					t.Errorf("expected field %q in schema, got %+v", tc.wantField, schema.Fields)
				}
			}
		})
	}
}

func TestDataModelManager_FetchFacets(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		respond   func(req *http.Request) (*http.Response, error)
		wantErr   bool
		wantCount int
	}{
		{
			name: "success",
			respond: func(req *http.Request) (*http.Response, error) {
				body := `[{"name":"Versionable","schemas":[{"name":"common","prefix":"common","fields":{"title":"string"}}]}]`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
			wantErr:   false,
			wantCount: 1,
		},
		{
			name: "client error",
			respond: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("client error")
			},
			wantErr:   true,
			wantCount: 0,
		},
		{
			name: "invalid json",
			respond: func(req *http.Request) (*http.Response, error) {
				body := `invalid`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
			wantErr:   true,
			wantCount: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newMockNuxeoClient(tc.respond)
			dmm := &dataModelManager{client: client, logger: slog.Default()}
			facets, err := dmm.FetchFacets(context.Background())
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tc.wantErr && facets != nil {
				if len(*facets) != tc.wantCount {
					t.Errorf("unexpected facet count: got %d, want %d", len(*facets), tc.wantCount)
				}
			}
		})
	}
}

func TestDataModelManager_FetchFacet(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		respond   func(req *http.Request) (*http.Response, error)
		wantErr   bool
		wantName  string
		wantField string
	}{
		{
			name: "success",
			respond: func(req *http.Request) (*http.Response, error) {
				body := `{"name":"Versionable","schemas":[{"name":"common","prefix":"common","fields":{"title":"string"}}]}`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
			wantErr:   false,
			wantName:  "Versionable",
			wantField: "title",
		},
		{
			name: "client error",
			respond: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("client error")
			},
			wantErr:   true,
			wantName:  "",
			wantField: "",
		},
		{
			name: "invalid json",
			respond: func(req *http.Request) (*http.Response, error) {
				body := `invalid`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
			wantErr:   true,
			wantName:  "",
			wantField: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newMockNuxeoClient(tc.respond)
			dmm := &dataModelManager{client: client, logger: slog.Default()}
			facet, err := dmm.FetchFacet(context.Background(), "Versionable")
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tc.wantErr && facet != nil {
				if facet.Name != tc.wantName {
					t.Errorf("unexpected facet name: got %q, want %q", facet.Name, tc.wantName)
				}
				if len(facet.Schemas) > 0 {
					if _, ok := facet.Schemas[0].Fields[tc.wantField]; !ok {
						t.Errorf("expected field %q in facet schema, got %+v", tc.wantField, facet.Schemas[0].Fields)
					}
				}
			}
		})
	}
}
