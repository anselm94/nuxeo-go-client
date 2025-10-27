package nuxeo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	nuxeoauth "github.com/anselm94/nuxeo-go-client/auth"
	resty "resty.dev/v3"
)

// --- Test for DefaultNuxeoClientOptions ---
func TestDefaultNuxeoClientOptions(t *testing.T) {
	opts := DefaultNuxeoClientOptions()

	if opts.Authenticator == nil {
		t.Errorf("Authenticator should not be nil")
	}
	if _, ok := opts.Authenticator.(*nuxeoauth.NoOpAuthenticator); !ok {
		t.Errorf("Authenticator should be NoOpAuthenticator, got %T", opts.Authenticator)
	}
	if opts.Logger == nil {
		t.Errorf("Logger should not be nil")
	}
	if opts.BeforeRequestMiddleware != nil {
		t.Errorf("BeforeRequestMiddleware should be nil by default")
	}
	if opts.AfterResponseMiddleware != nil {
		t.Errorf("AfterResponseMiddleware should be nil by default")
	}
	if opts.Timeout != 30*time.Second {
		t.Errorf("Timeout should be 30s by default, got %v", opts.Timeout)
	}
	if opts.CustomHeaders == nil {
		t.Errorf("CustomHeaders should not be nil")
	}
	if len(opts.CustomHeaders) != 0 {
		t.Errorf("CustomHeaders should be empty by default")
	}
}

// --- Minimal mock Authenticator ---
type mockAuthenticator struct{}

func (m *mockAuthenticator) GetAuthHeaders(ctx context.Context, req *http.Request) map[string]string {
	return map[string]string{"X-Mock": "mocked"}
}

func TestNewClient_DefaultsAndCustomOptions(t *testing.T) {
	t.Run("nil options uses defaults", func(t *testing.T) {
		client := NewClient("http://localhost", nil)
		if client.authenticator == nil {
			t.Errorf("authenticator should not be nil")
		}
		if client.logger == nil {
			t.Errorf("logger should not be nil")
		}
		if client.timeout != 30*time.Second {
			t.Errorf("timeout should be 30s by default, got %v", client.timeout)
		}
		if client.headers == nil {
			t.Errorf("headers should not be nil")
		}
		if len(client.headers) != 0 {
			t.Errorf("headers should be empty by default")
		}
	})

	t.Run("custom options are propagated", func(t *testing.T) {
		customHeaders := map[string]string{"X-Test": "value"}
		customTimeout := 10 * time.Second
		mockAuth := &mockAuthenticator{}
		opts := &nuxeoClientOptions{
			Authenticator:           mockAuth,
			Logger:                  nil,
			BeforeRequestMiddleware: nil,
			AfterResponseMiddleware: nil,
			Timeout:                 customTimeout,
			CustomHeaders:           customHeaders,
		}
		client := NewClient("http://localhost", opts)
		if client.authenticator != mockAuth {
			t.Errorf("custom authenticator not set")
		}
		if client.timeout != customTimeout {
			t.Errorf("custom timeout not set")
		}
		if client.headers["X-Test"] != "value" {
			t.Errorf("custom header not set")
		}
	})
}

func TestNuxeoClient_Close(t *testing.T) {
	client := &NuxeoClient{restClient: resty.New()}
	err := client.Close()
	if err != nil {
		t.Errorf("Close() should not return error, got: %v", err)
	}
}

func TestNuxeoClient_ManagerGetters(t *testing.T) {
	client := NewClient("http://localhost", nil)
	ctx := context.Background()

	if m := client.CapabilitiesManager(ctx); m == nil {
		t.Errorf("CapabilitiesManager should not be nil")
	}
	if m := client.Repository(); m == nil {
		t.Errorf("Repository should not be nil")
	}
	if m := client.RepositoryWithName("repo"); m == nil {
		t.Errorf("RepositoryWithName should not be nil")
	}
	if m := client.OperationManager(); m == nil {
		t.Errorf("OperationManager should not be nil")
	}
	if m := client.UserManager(); m == nil {
		t.Errorf("UserManager should not be nil")
	}
	if m := client.DirectoryManager(); m == nil {
		t.Errorf("DirectoryManager should not be nil")
	}
	if m := client.TaskManager(); m == nil {
		t.Errorf("TaskManager should not be nil")
	}
	if m := client.BatchUploadManager(); m == nil {
		t.Errorf("BatchUploadManager should not be nil")
	}
	if m := client.DataModelManager(); m == nil {
		t.Errorf("DataModelManager should not be nil")
	}
}

// --- Test for NuxeoClient.NewRequest ---
// (see below for extended tests)

// --- Test for Authenticator header injection ---
func TestNuxeoClient_AuthenticatorHeaders(t *testing.T) {
	ts := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Mock") != "mocked" {
			t.Errorf("Authenticator header not injected, got: %v", r.Header.Get("X-Mock"))
			for k, v := range r.Header {
				t.Logf("Header: %s=%v", k, v)
			}
		}
		w.WriteHeader(http.StatusOK)
	})

	testServer := httptest.NewServer(ts)
	defer testServer.Close()

	client := NewClient(testServer.URL, &nuxeoClientOptions{
		Authenticator: &mockAuthenticator{},
	})
	resp, err := client.NewRequest(context.Background(), nil).Request.Get("/")
	if err != nil {
		t.Errorf("Request failed: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		t.Errorf("Expected 200 OK, got: %v", resp.StatusCode())
	}
}

// --- Test for Authenticator header injection for fallback ---
func TestNuxeoClient_AuthenticatorHeadersFallback(t *testing.T) {
	ts := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	testServer := httptest.NewServer(ts)
	defer testServer.Close()

	client := NewClient(testServer.URL, &nuxeoClientOptions{Authenticator: nil})
	resp, err := client.NewRequest(context.Background(), nil).Request.Get("/")
	if err != nil {
		t.Errorf("Request failed: %v", err)
	}
	if resp.Request.Header.Get("X-Mock") != "" {
		t.Errorf("NoOpAuthenticator should not inject headers, got: %v", resp.Request.Header.Get("X-Mock"))
	}
}

// --- Test for Before/After Middleware invocation ---
func TestNuxeoClient_MiddlewareInvocation(t *testing.T) {
	beforeCalled := false
	before := func(c *resty.Client, r *resty.Request) error {
		beforeCalled = true
		return nil
	}
	after := func(c *resty.Client, resp *resty.Response) error {
		// Can't reliably test afterCalled without a real server
		return nil
	}
	client := NewClient("http://localhost", &nuxeoClientOptions{
		BeforeRequestMiddleware: before,
		AfterResponseMiddleware: after,
	})
	// Trigger a request (will fail, but middleware should be called)
	_, _ = client.NewRequest(context.Background(), nil).Request.Get("/notfound")
	if !beforeCalled {
		t.Errorf("BeforeRequestMiddleware not called")
	}
	// AfterResponseMiddleware is only called on response, so we can't easily test without a real server
	// But we can check that it's set
	if client.middlewareAfterResponse == nil {
		t.Errorf("AfterResponseMiddleware not set")
	}
}

// --- Test for timeout enforcement ---
func TestNuxeoClient_TimeoutEnforcement(t *testing.T) {
	client := NewClient("http://localhost", &nuxeoClientOptions{Timeout: 1 * time.Millisecond})
	start := time.Now()
	_, err := client.restClient.R().Get("http://10.255.255.1") // unroutable IP, should timeout
	elapsed := time.Since(start)
	if err == nil {
		t.Errorf("Expected timeout error, got nil")
	} else {
		if !isTimeoutError(err) {
			t.Errorf("Expected timeout error, got: %v", err)
		}
	}
	if elapsed > 100*time.Millisecond {
		t.Errorf("Request did not timeout quickly, elapsed: %v", elapsed)
	}
}

// isTimeoutError checks if the error is a timeout error
func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	type timeout interface {
		Timeout() bool
	}
	if terr, ok := err.(timeout); ok && terr.Timeout() {
		return true
	}
	if err.Error() != "" && (containsString(err.Error(), "timeout") || containsString(err.Error(), "i/o timeout")) {
		return true
	}
	return false
}

func containsString(s, substr string) bool {
	return len(substr) > 0 && len(s) > 0 && (s == substr || (len(s) > len(substr) && (s[:len(substr)] == substr || containsString(s[1:], substr))))
}

// --- Test for CurrentUser and ServerVersion error handling ---
func TestNuxeoClient_CurrentUser_ServerVersion_Error(t *testing.T) {
	client := NewClient("http://localhost", nil)
	ctx := context.Background()
	_, err := client.CurrentUser(ctx)
	if err == nil {
		t.Errorf("Expected error from CurrentUser on mock client, got nil")
	}
	_, err = client.ServerVersion(ctx)
	if err == nil {
		t.Errorf("Expected error from ServerVersion on mock client, got nil")
	}
}

// --- Test for invalid base URL ---
func TestNuxeoClient_InvalidBaseURL(t *testing.T) {
	client := NewClient("://invalid-url", nil)
	resp, err := client.restClient.R().Get("/test")
	if err == nil {
		t.Errorf("Expected error for invalid base URL, got nil")
	}
	if resp != nil {
		t.Errorf("Expected nil response for invalid base URL, got: %v", resp)
	}
}

func TestNuxeoClient_NewRequest(t *testing.T) {
	client := NewClient("http://localhost", &nuxeoClientOptions{
		CustomHeaders: map[string]string{"X-Client": "client-value"},
	})

	t.Run("sets context and client headers", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "key", "val")
		req := client.NewRequest(ctx, nil)
		if req.Request.Context().Value("key") != "val" {
			t.Errorf("context not propagated")
		}
		if req.Request.Header.Get("X-Client") != "client-value" {
			t.Errorf("client header not set")
		}
	})

	t.Run("sets request custom headers and overrides client headers", func(t *testing.T) {
		opts := NewNuxeoRequestOptions().SetHeader("X-Client", "override").SetHeader("X-Req", "req-value")
		req := client.NewRequest(context.Background(), opts)
		if req.Request.Header.Get("X-Client") != "override" {
			t.Errorf("request header should override client header")
		}
		if req.Request.Header.Get("X-Req") != "req-value" {
			t.Errorf("request custom header not set")
		}
	})

	t.Run("nil options does not panic", func(t *testing.T) {
		req := client.NewRequest(context.Background(), nil)
		if req == nil {
			t.Errorf("NewRequest should not return nil")
		}
	})

	t.Run("empty client and request headers", func(t *testing.T) {
		client := NewClient("http://localhost", nil)
		opts := NewNuxeoRequestOptions()
		req := client.NewRequest(context.Background(), opts)
		if len(req.Request.Header) != 0 {
			t.Errorf("headers should be empty")
		}
	})
}
