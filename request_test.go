package nuxeo

import (
	"reflect"
	"resty.dev/v3"
	"testing"
)

// Test NewNuxeoRequestOptions initializes maps correctly
func TestNewNuxeoRequestOptions_InitializesMaps(t *testing.T) {
	opts := NewNuxeoRequestOptions()
	if opts.enrichers == nil || opts.fetchProperties == nil || opts.translateProperties == nil {
		t.Errorf("Expected maps to be initialized, got nil")
	}
}

// Table-driven tests for all setters
func TestNuxeoRequestOptions_Setters(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		setup func(*nuxeoRequestOptions)
		check func(*nuxeoRequestOptions, *testing.T)
	}{
		{
			name:  "SetRepositoryName",
			setup: func(o *nuxeoRequestOptions) { o.SetRepositoryName("repo1") },
			check: func(o *nuxeoRequestOptions, t *testing.T) {
				if o.repositoryName != "repo1" {
					t.Errorf("Expected repositoryName 'repo1', got '%s'", o.repositoryName)
				}
			},
		},
		{
			name:  "SetHeader",
			setup: func(o *nuxeoRequestOptions) { o.customHeaders = make(map[string]string); o.SetHeader("X-Test", "val") },
			check: func(o *nuxeoRequestOptions, t *testing.T) {
				if o.customHeaders["X-Test"] != "val" {
					t.Errorf("Expected customHeaders['X-Test']='val', got '%s'", o.customHeaders["X-Test"])
				}
			},
		},
		{
			name:  "SetEnricher",
			setup: func(o *nuxeoRequestOptions) { o.SetEnricher("document", []string{"foo", "bar"}) },
			check: func(o *nuxeoRequestOptions, t *testing.T) {
				want := []string{"foo", "bar"}
				if !reflect.DeepEqual(o.enrichers["document"], want) {
					t.Errorf("Expected enrichers['document']=%v, got %v", want, o.enrichers["document"])
				}
			},
		},
		{
			name:  "SetFetchProperties",
			setup: func(o *nuxeoRequestOptions) { o.SetFetchProperties("group", []string{"a", "b"}) },
			check: func(o *nuxeoRequestOptions, t *testing.T) {
				want := []string{"a", "b"}
				if !reflect.DeepEqual(o.fetchProperties["group"], want) {
					t.Errorf("Expected fetchProperties['group']=%v, got %v", want, o.fetchProperties["group"])
				}
			},
		},
		{
			name:  "SetTranslatedProperties",
			setup: func(o *nuxeoRequestOptions) { o.SetTranslatedProperties("directory", []string{"x"}) },
			check: func(o *nuxeoRequestOptions, t *testing.T) {
				want := []string{"x"}
				if !reflect.DeepEqual(o.translateProperties["directory"], want) {
					t.Errorf("Expected translateProperties['directory']=%v, got %v", want, o.translateProperties["directory"])
				}
			},
		},
		{
			name:  "SetSchemas",
			setup: func(o *nuxeoRequestOptions) { o.SetSchemas([]string{"schema1", "schema2"}) },
			check: func(o *nuxeoRequestOptions, t *testing.T) {
				want := []string{"schema1", "schema2"}
				if !reflect.DeepEqual(o.schemas, want) {
					t.Errorf("Expected schemas=%v, got %v", want, o.schemas)
				}
			},
		},
		{
			name:  "SetDepth",
			setup: func(o *nuxeoRequestOptions) { o.SetDepth(3) },
			check: func(o *nuxeoRequestOptions, t *testing.T) {
				if o.depth != 3 {
					t.Errorf("Expected depth=3, got %d", o.depth)
				}
			},
		},
		{
			name:  "SetVersion",
			setup: func(o *nuxeoRequestOptions) { o.SetVersion("major") },
			check: func(o *nuxeoRequestOptions, t *testing.T) {
				if o.version != "major" {
					t.Errorf("Expected version='major', got '%s'", o.version)
				}
			},
		},
		{
			name:  "SetTransactionTimeout",
			setup: func(o *nuxeoRequestOptions) { o.SetTransactionTimeout(42) },
			check: func(o *nuxeoRequestOptions, t *testing.T) {
				if o.transactionTimeout != 42 {
					t.Errorf("Expected transactionTimeout=42, got %d", o.transactionTimeout)
				}
			},
		},
		{
			name:  "SetHttpTimeout",
			setup: func(o *nuxeoRequestOptions) { o.SetHttpTimeout(99) },
			check: func(o *nuxeoRequestOptions, t *testing.T) {
				if o.httpTimeout != 99 {
					t.Errorf("Expected httpTimeout=99, got %d", o.httpTimeout)
				}
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			opts := NewNuxeoRequestOptions()
			tc.setup(opts)
			tc.check(opts, t)
		})
	}
}

// Test setNuxeoOption applies all headers correctly
func TestNuxeoRequest_setNuxeoOption_AllHeaders(t *testing.T) {
	req := resty.New().R()
	nr := &nuxeoRequest{Request: req}
	opts := NewNuxeoRequestOptions()
	opts.SetRepositoryName("repo")
	opts.SetHeader("X-Custom", "custom-value")
	opts.SetEnricher("document", []string{"e1", "e2"})
	opts.SetFetchProperties("document", []string{"fp1"})
	opts.SetTranslatedProperties("directory", []string{"tp1"})
	opts.SetSchemas([]string{"s1", "s2"})
	opts.SetDepth(2)
	opts.SetVersion("v1")
	opts.SetTransactionTimeout(10)
	opts.SetHttpTimeout(15)

	nr.setNuxeoOption(opts)

	headers := req.Header

	cases := []struct {
		key  string
		want string
	}{
		{"X-NXRepository", "repo"},
		{"X-Custom", "custom-value"},
		{"enrichers-document", "e1,e2"},
		{"fetch-document", "fp1"},
		{"translate-directory", "tp1"},
		{"properties", "s1,s2"},
		{"depth", "2"},
		{"X-Versioning-Option", "v1"},
		{"Nuxeo-Transaction-Timeout", "10"},
		{"timeout", "15"},
	}
	for _, tc := range cases {
		if got := headers.Get(tc.key); got != tc.want {
			t.Errorf("Header %q: want %q, got %q", tc.key, tc.want, got)
		}
	}
}

// Test setNuxeoOption with nil options
func TestNuxeoRequest_setNuxeoOption_NilOptions(t *testing.T) {
	req := resty.New().R()
	nr := &nuxeoRequest{Request: req}
	ret := nr.setNuxeoOption(nil)
	if ret != nr {
		t.Errorf("Expected receiver to be returned unchanged when options is nil")
	}
}

// Test transaction/http timeout logic
func TestNuxeoRequest_setNuxeoOption_TimeoutLogic(t *testing.T) {
	req := resty.New().R()
	nr := &nuxeoRequest{Request: req}
	opts := NewNuxeoRequestOptions()
	opts.SetTransactionTimeout(20)
	// httpTimeout not set, should be set to transactionTimeout+5
	nr.setNuxeoOption(opts)
	if got := req.Header.Get("timeout"); got != "25" {
		t.Errorf("Expected timeout header to be '25', got '%s'", got)
	}
}
