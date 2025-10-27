package nuxeoauth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/oauth2"
)

type mockTokenSource struct {
	token *oauth2.Token
	err   error
}

func (m *mockTokenSource) Token() (*oauth2.Token, error) {
	return m.token, m.err
}

func TestOAuth2OptionsConstructors(t *testing.T) {
	t.Run("AuthorizationCodeOptions", func(t *testing.T) {
		opt := NewOAuth2AuthorizationCodeOptions("cid", "csecret", "redir")
		if opt.clientId != "cid" || opt.clientSecret != "csecret" || opt.redirectUrl != "redir" {
			t.Errorf("AuthorizationCodeOptions fields not set correctly: %+v", opt)
		}
	})
	t.Run("ClientCredentialsOptions", func(t *testing.T) {
		opt := NewOAuth2ClientCredentialsOptions("cid", "csecret")
		if opt.clientId != "cid" || opt.clientSecret != "csecret" || opt.redirectUrl != "" {
			t.Errorf("ClientCredentialsOptions fields not set correctly: %+v", opt)
		}
	})
	t.Run("JwtOptions", func(t *testing.T) {
		opt := NewOAuth2JwtOptions("jwt-token")
		if opt.jwtToken != "jwt-token" {
			t.Errorf("JwtOptions field not set correctly: %+v", opt)
		}
	})
}

func TestNewOAuth2Authenticator(t *testing.T) {
	t.Run("JWT", func(t *testing.T) {
		auth := NewOAuth2Authenticator(NewOAuth2JwtOptions("jwt"), "http://base")
		if auth == nil {
			t.Fatal("Expected authenticator, got nil")
		}
		if auth.options.jwtToken != "jwt" {
			t.Errorf("JWT token not set correctly")
		}
	})
	t.Run("ClientCredentials", func(t *testing.T) {
		auth := NewOAuth2Authenticator(NewOAuth2ClientCredentialsOptions("cid", "csecret"), "http://base")
		if auth == nil {
			t.Fatal("Expected authenticator, got nil")
		}
		if auth.options.clientId != "cid" || auth.options.clientSecret != "csecret" {
			t.Errorf("Client credentials not set correctly")
		}
	})
	t.Run("AuthorizationCode", func(t *testing.T) {
		auth := NewOAuth2Authenticator(NewOAuth2AuthorizationCodeOptions("cid", "csecret", "redir"), "http://base")
		if auth == nil {
			t.Fatal("Expected authenticator, got nil")
		}
		if auth.options.redirectUrl != "redir" {
			t.Errorf("Redirect URL not set correctly")
		}
	})
	t.Run("InvalidOptions", func(t *testing.T) {
		auth := NewOAuth2Authenticator(OAuth2Options{}, "http://base")
		if auth != nil {
			t.Error("Expected nil for invalid options")
		}
	})
}

func TestGetAuthHeaders_JWT(t *testing.T) {
	auth := NewOAuth2Authenticator(NewOAuth2JwtOptions("jwt-token"), "http://base")
	ts := auth.GetTokenSource(context.Background())
	tok, err := ts.Token()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	_ = tok
	auth.tokenSource = ts
	headers := auth.GetAuthHeaders(context.Background(), &http.Request{})
	if headers["Authorization"] != "Bearer jwt-token" {
		t.Errorf("Expected Authorization header with jwt-token, got: %v", headers)
	}
}

func TestGetAuthHeaders_ClientCredentials(t *testing.T) {
	auth := NewOAuth2Authenticator(NewOAuth2ClientCredentialsOptions("cid", "csecret"), "http://base")
	tok := &oauth2.Token{AccessToken: "access-token"}
	auth.tokenSource = &mockTokenSource{token: tok}
	headers := auth.GetAuthHeaders(context.Background(), &http.Request{})
	if headers["Authorization"] != "Bearer access-token" {
		t.Errorf("Expected Authorization header with access-token, got: %v", headers)
	}
}

func TestGetAuthHeaders_NoTokenSource(t *testing.T) {
	auth := NewOAuth2Authenticator(NewOAuth2ClientCredentialsOptions("cid", "csecret"), "http://base")
	headers := auth.GetAuthHeaders(context.Background(), &http.Request{})
	if len(headers) != 0 {
		t.Errorf("Expected empty headers when no tokenSource, got: %v", headers)
	}
}

func TestAuthCodeUrl(t *testing.T) {
	auth := NewOAuth2Authenticator(NewOAuth2AuthorizationCodeOptions("cid", "csecret", "redir"), "http://base")
	url := auth.AuthCodeUrl(context.Background())
	if url == "" || url == "http://base/oauth2/authorize" {
		t.Errorf("Expected non-empty auth code URL, got: %v", url)
	}
}

func TestSetAuthCode_Success(t *testing.T) {
	// Simulate OAuth2 server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth2/token" && r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"access_token": "tok",
				"token_type":   "bearer",
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	auth := NewOAuth2Authenticator(NewOAuth2AuthorizationCodeOptions("cid", "csecret", "redir"), ts.URL)
	err := auth.SetAuthCode(context.Background(), "good-code")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	tsToken := auth.tokenSource
	if tsToken == nil {
		t.Error("Expected tokenSource to be set after SetAuthCode")
	}
	// Try to get token from tokenSource
	_, err = tsToken.Token()
	if err != nil {
		t.Errorf("Expected no error from tokenSource, got: %v", err)
	}
}

func TestSetAuthCode_Error(t *testing.T) {
	// Simulate OAuth2 server returning error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth2/token" && r.Method == http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "invalid_grant",
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	auth := NewOAuth2Authenticator(NewOAuth2AuthorizationCodeOptions("cid", "csecret", "redir"), ts.URL)
	err := auth.SetAuthCode(context.Background(), "bad-code")
	if err == nil {
		t.Error("Expected error for bad code, got nil")
	}
}

func TestGetTokenSource_ClientCredentials(t *testing.T) {
	// Simulate OAuth2 server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth2/token" && r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"access_token": "client-token",
				"token_type":   "bearer",
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	auth := NewOAuth2Authenticator(NewOAuth2ClientCredentialsOptions("cid", "csecret"), ts.URL)
	tsToken := auth.GetTokenSource(context.Background())
	_, err := tsToken.Token()
	if err != nil {
		t.Errorf("Expected no error from tokenSource, got: %v", err)
	}
}

func TestGetTokenSource_AlreadySet(t *testing.T) {
	auth := NewOAuth2Authenticator(NewOAuth2ClientCredentialsOptions("cid", "csecret"), "http://base")
	// Set a custom tokenSource
	customTS := &mockTokenSource{token: &oauth2.Token{AccessToken: "preset-token"}}
	auth.tokenSource = customTS

	// Call GetTokenSource, should return the same tokenSource
	ts := auth.GetTokenSource(context.Background())
	if ts != customTS {
		t.Errorf("Expected GetTokenSource to return the already set tokenSource")
	}

	// Ensure token is as expected
	token, err := ts.Token()
	if err != nil {
		t.Fatalf("Unexpected error from tokenSource: %v", err)
	}
	if token.AccessToken != "preset-token" {
		t.Errorf("Expected tokenSource to return preset-token, got: %v", token.AccessToken)
	}
}
