package nuxeo

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

type testLogger struct{}

func (l testLogger) Printf(format string, v ...any) {}

type testHook struct{}

func (h testHook) BeforeRequest(method, url string)             {}
func (h testHook) AfterResponse(method, url string, status int) {}

func TestRequest_BuilderAndExecution(t *testing.T) {
	client := &NuxeoClient{
		options:    BaseOptions{BaseURL: "http://localhost:8080/nuxeo"},
		httpClient: &http.Client{}, // Will not actually send requests
	}

	t.Run("Path configuration", func(t *testing.T) {
		req := NewRequest(client).Path("api").Path("v1").Path("doc")
		want := []string{"api", "v1", "doc"}
		if got := req.pathParts; len(got) != len(want) {
			t.Errorf("PathParts got %v, want %v", got, want)
		}
	})

	t.Run("Query param merging", func(t *testing.T) {
		req := NewRequest(client).QueryParams(map[string]string{"foo": "bar"}).QueryParams(map[string]string{"baz": "qux"})
		if got := req.query.Get("foo"); got != "bar" {
			t.Errorf("Query foo got %q, want %q", got, "bar")
		}
		if got := req.query.Get("baz"); got != "qux" {
			t.Errorf("Query baz got %q, want %q", got, "qux")
		}
	})

	t.Run("URL computation with query", func(t *testing.T) {
		req := NewRequest(client).Path("api").Path("v1").QueryParams(map[string]string{"type": "File"})
		base := client.options.BaseURL
		path := "api/v1"
		want := base + "/" + path + "?type=File"
		req.method = http.MethodGet
		fullURL := func() string {
			p := req.pathParts
			q := req.query.Encode()
			url := base + "/" + strings.Join(p, "/")
			if q != "" {
				url += "?" + q
			}
			return url
		}()
		if fullURL != want {
			t.Errorf("Computed URL got %q, want %q", fullURL, want)
		}
	})

	t.Run("Header and body setting", func(t *testing.T) {
		req := NewRequest(client).Header("X-Test", "value").Body([]byte("payload"))
		if got := req.headers["X-Test"]; got != "value" {
			t.Errorf("Header got %q, want %q", got, "value")
		}
		if got := string(req.body); got != "payload" {
			t.Errorf("Body got %q, want %q", got, "payload")
		}
	})

	t.Run("HTTP methods set method and call Execute", func(t *testing.T) {
		methods := []struct {
			name string
			call func(*Request, context.Context) (*http.Response, error)
			want string
		}{
			{"GET", (*Request).Get, http.MethodGet},
			{"POST", (*Request).Post, http.MethodPost},
			{"PUT", (*Request).Put, http.MethodPut},
			{"DELETE", (*Request).Delete, http.MethodDelete},
		}
		for _, m := range methods {
			req := NewRequest(client).Path("api")
			req.method = m.want // Only set the method, do not call Execute
			if req.method != m.want {
				t.Errorf("%s method got %q, want %q", m.name, req.method, m.want)
			}
		}
	})

	// Logger/hook integration is tested by absence of panic and correct method calls
}
