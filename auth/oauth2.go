package auth

/**
OAuth2Authenticator implements OAuth2 authentication with support for Authorization Grant Flow, Client Credentials Flow, and JWT.

Example:
```go
import (
	"context"
	"fmt"
	"github.com/anselm94/nuxeo-go-client"
	"github.com/anselm94/nuxeo-go-client/auth"
)

ctx := context.Background()

//// Using Authorization Grant Flow ////

authInfo := auth.OAuth2Info{
	ClientID:     "your-client-id",
	ClientSecret: "your-client-secret",
	RedirectURL:  "https://your-redirect-url.com/callback",
}
authenticator := auth.NewOAuth2Authenticator(authInfo, "https://nuxeo.example.com/nuxeo")
authURL := authenticator.AuthCodeUrl(ctx)
fmt.Printf("Visit the URL for the auth dialog: %v", authURL)

// After obtaining the auth code from the redirect URL
authCode := "authorization-code-from-callback"
err := authenticator.SetAuthCode(ctx, authCode)
if err != nil {
	panic(err)
}

//// Using Client Credentials Flow ////

authInfo := auth.OAuth2Info{
	ClientID:     "your-client-id",
	ClientSecret: "your-client-secret",
}
authenticator := auth.NewOAuth2Authenticator(authInfo, "https://nuxeo.example.com/nuxeo")

//// Using JWT Flow ////

authInfo := auth.OAuth2Info{
	JwtToken: "your-jwt-token",
}
authenticator := auth.NewOAuth2Authenticator(authInfo, "https://nuxeo.example.com/nuxeo")

//// Creating Nuxeo Client ////

client, err := nuxeo.NewClient(ctx,
	nuxeo.WithBaseURL("https://nuxeo.example.com/nuxeo"),
	nuxeo.WithAuthenticator(authenticator),
)
if err != nil {
	panic(err)
}
// Use client...
*/

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	oauth2TokenPath = "/oauth2/token"
	oauth2AuthPath  = "/oauth2/authorize"

	OAuth2AuthorizeState = "nuxeo-go-client-state"
)

type OAuth2Info struct {
	// Authorization Grant Flow + Client Credentials

	// OAuth2 Client ID
	ClientID string
	// OAuth2 Client Secret
	ClientSecret string

	// Authorization Grant Flow (only)

	// OAuth2 Redirect URL
	RedirectURL string

	// JWT Flow

	// JWT Token
	JwtToken string
}

// OAuth2Authenticator delegates to OAuth2Token and supports refresh.
type OAuth2Authenticator struct {
	authInfo OAuth2Info
	baseUrl  string

	tokenSource oauth2.TokenSource
}

func NewOAuth2Authenticator(authInfo OAuth2Info, baseUrl string) *OAuth2Authenticator {
	// either token or client credentials should be provided
	if authInfo.JwtToken == "" && (authInfo.ClientID == "" || authInfo.ClientSecret == "") {
		return nil
	}

	return &OAuth2Authenticator{
		authInfo: authInfo,
		baseUrl:  baseUrl,
	}
}

func (a *OAuth2Authenticator) GetAuthHeaders(ctx context.Context, req *http.Request) map[string]string {
	headers := make(map[string]string)
	if a.tokenSource != nil {
		if token, err := a.tokenSource.Token(); err == nil && token != nil {
			headers["Authorization"] = "Bearer " + token.AccessToken
		}
	}
	return headers
}

func (a *OAuth2Authenticator) getAuthGrantFlowConfig() oauth2.Config {
	return oauth2.Config{
		Scopes:      []string{},
		RedirectURL: a.authInfo.RedirectURL,
		Endpoint: oauth2.Endpoint{
			TokenURL: a.baseUrl + oauth2TokenPath,
			AuthURL:  a.baseUrl + oauth2AuthPath,
		},
	}
}

func (a *OAuth2Authenticator) getClientCredentialsFlowConfig() clientcredentials.Config {
	return clientcredentials.Config{
		ClientID:     a.authInfo.ClientID,
		ClientSecret: a.authInfo.ClientSecret,
		Scopes:       []string{},
		TokenURL:     a.baseUrl + oauth2TokenPath,
	}
}

func (a *OAuth2Authenticator) GetTokenSource(ctx context.Context) oauth2.TokenSource {
	// if token source is already set, return it
	if a.tokenSource != nil {
		return a.tokenSource
	}

	// if JWT token is provided, return static token source
	if a.authInfo.JwtToken != "" {
		staticSource := oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: a.authInfo.JwtToken,
		})
		a.tokenSource = staticSource
	} else {
		// if redirect URL is not provided, use client credentials flow
		if a.authInfo.RedirectURL == "" {
			config := a.getClientCredentialsFlowConfig()
			a.tokenSource = config.TokenSource(ctx)
		}
		// otherwise, use auth grant flow (tokenSource will be set after auth code exchange)
	}

	return a.tokenSource
}

func (a *OAuth2Authenticator) AuthCodeUrl(ctx context.Context) string {
	config := a.getAuthGrantFlowConfig()
	return config.AuthCodeURL(OAuth2AuthorizeState, oauth2.AccessTypeOffline)
}

func (a *OAuth2Authenticator) SetAuthCode(ctx context.Context, code string) error {
	config := a.getAuthGrantFlowConfig()

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("Error exchanging auth code for token: %w", err)
	}

	// also update the token source
	a.tokenSource = config.TokenSource(ctx, token)
	return nil
}
