package nuxeoauth

/**
OAuth2Authenticator implements OAuth2 authentication with support for Authorization Grant Flow, Client Credentials Flow, and JWT.

Example:
```go
import (
	"fmt"
	"github.com/anselm94/nuxeo-go-client/auth"
)

ctx := context.Background()

//// Using Authorization Grant Flow ////

oauth2Option := nuxeoauth.NewOAuth2AuthorizationGrantOptions("your-client-id", "your-client-secret", "https://your-redirect-url.com/callback")
authenticator := nuxeoauth.NewOAuth2Authenticator(oauth2Option, "https://nuxeo.example.com/nuxeo")
authURL := authenticator.AuthCodeUrl(ctx)
fmt.Printf("Visit the URL for the auth dialog: %v", authURL)

// After obtaining the auth code from the redirect URL
authCode := "authorization-code-from-callback"
err := authenticator.SetAuthCode(ctx, authCode)
if err != nil {
	panic(err)
}

//// Using Client Credentials Flow ////

oauth2Option := nuxeoauth.NewOAuth2ClientCredentialsOptions("your-client-id", "your-client-secret")
authenticator := nuxeoauth.NewOAuth2Authenticator(oauth2Option, "https://nuxeo.example.com/nuxeo")

//// Using JWT Flow ////

oauth2Option := nuxeoauth.NewOAuth2JwtOptions("your-jwt-token")
authenticator := nuxeoauth.NewOAuth2Authenticator(oauth2Option, "https://nuxeo.example.com/nuxeo")

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

type OAuth2Options struct {
	// OAuth2 Client ID
	clientId string
	// OAuth2 Client Secret
	clientSecret string

	// OAuth2 Redirect URL
	redirectUrl string

	// JWT Token
	jwtToken string
}

// Create OAuth2Info for Client Credentials Flow
func NewOAuth2ClientCredentialsOptions(clientId string, clientSecret string) OAuth2Options {
	return OAuth2Options{
		clientId:     clientId,
		clientSecret: clientSecret,
	}
}

// Create OAuth2Info for Authorization Code Grant Flow
func NewOAuth2AuthorizationCodeOptions(clientId string, clientSecret string, redirectUrl string) OAuth2Options {
	return OAuth2Options{
		clientId:     clientId,
		clientSecret: clientSecret,
		redirectUrl:  redirectUrl,
	}
}

// Create OAuth2Info for JWT Flow
func NewOAuth2JwtOptions(jwtToken string) OAuth2Options {
	return OAuth2Options{
		jwtToken: jwtToken,
	}
}

// OAuth2Authenticator delegates to OAuth2Token and supports refresh.
type OAuth2Authenticator struct {
	options OAuth2Options
	baseUrl string

	tokenSource oauth2.TokenSource
}

// NewOAuth2Authenticator creates a new OAuth2Authenticator with the given OAuth2Info and Nuxeo base URL (with `/nuxeo`).
func NewOAuth2Authenticator(options OAuth2Options, baseUrl string) *OAuth2Authenticator {
	// either token or client credentials should be provided
	if options.jwtToken == "" && (options.clientId == "" || options.clientSecret == "") {
		return nil
	}

	return &OAuth2Authenticator{
		options: options,
		baseUrl: baseUrl,
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
		RedirectURL: a.options.redirectUrl,
		Endpoint: oauth2.Endpoint{
			TokenURL: a.baseUrl + oauth2TokenPath,
			AuthURL:  a.baseUrl + oauth2AuthPath,
		},
	}
}

func (a *OAuth2Authenticator) getClientCredentialsFlowConfig() clientcredentials.Config {
	return clientcredentials.Config{
		ClientID:     a.options.clientId,
		ClientSecret: a.options.clientSecret,
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
	if a.options.jwtToken != "" {
		staticSource := oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: a.options.jwtToken,
		})
		a.tokenSource = staticSource
	} else {
		// if redirect URL is not provided, use client credentials flow
		if a.options.redirectUrl == "" {
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
