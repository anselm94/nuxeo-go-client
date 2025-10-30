package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/anselm94/nuxeo-go-client"
	nuxeoauth "github.com/anselm94/nuxeo-go-client/auth"
)

type CustomAuthenticator struct{}

func (ca *CustomAuthenticator) GetAuthHeaders(ctx context.Context, req *http.Request) map[string]string {
	return map[string]string{
		"X-Custom-Auth": "custom-auth-value",
	}
}

func main() {
	// 1. initialize context

	ctx := context.Background()

	// 2. setup authentication
	var authenticator nuxeo.Authenticator

	// 2.a Basic Authenticator
	authenticator = nuxeoauth.NewBasicAuthenticator("Administrator", "Administrator")

	// ####
	// #### 2.b Bearer Authenticator ####
	// ####
	// authenticator = nuxeoauth.NewBearerAuthenticator("your-bearer-token")

	// ####
	// #### 2.c Token Authenticator ####
	// ####
	// authenticator = nuxeoauth.NewTokenAuthenticator("your-token")

	// ####
	// #### 2.d.1 OAuth2 Authenticator (Authorization Code Flow) ####
	// ####
	// authCodeFlowOptions := nuxeoauth.NewOAuth2AuthorizationCodeOptions("your-client-id", "your-client-secret", "your-redirect-uri")
	// authCodeOauthAuthenticator := nuxeoauth.NewOAuth2Authenticator(authCodeFlowOptions, "https://demo.nuxeo.com/nuxeo")
	// authenticator = authCodeOauthAuthenticator

	// authURL := authCodeOauthAuthenticator.AuthCodeUrl(ctx)
	// fmt.Printf("Visit the URL for the auth dialog: %v", authURL)

	// // After obtaining the auth code from the redirect URL
	// authCode := "authorization-code-from-callback"
	// err := authCodeOauthAuthenticator.SetAuthCode(ctx, authCode)
	// if err != nil {
	// 	panic(err)
	// }

	// ####
	// #### 2.d.2 OAuth2 Authenticator (Client Credentials Flow) ####
	// ####
	// clientCredFlowOptions := nuxeoauth.NewOAuth2ClientCredentialsOptions("your-client-id", "your-client-secret")
	// clientCredOauthAuthenticator := nuxeoauth.NewOAuth2Authenticator(clientCredFlowOptions, "https://demo.nuxeo.com/nuxeo")
	// authenticator = clientCredOauthAuthenticator

	// ####
	// #### 2.d.3 OAuth2 Authenticator (Jwt Token Flow) ####
	// ####
	// jwtFlowOptions := nuxeoauth.NewOAuth2JwtOptions("your-jwt-token")
	// jwtOauthAuthenticator := nuxeoauth.NewOAuth2Authenticator(jwtFlowOptions, "https://demo.nuxeo.com/nuxeo")
	// authenticator = jwtOauthAuthenticator

	// ####
	// #### 2.e Custom Authenticator ####
	// ####
	// authenticator = &CustomAuthenticator{}

	// 3. Initialize Nuxeo client
	nuxeoClientOptions := nuxeo.DefaultNuxeoClientOptions()
	nuxeoClientOptions.Authenticator = authenticator

	// 4. create Nuxeo client
	nuxeoClient := nuxeo.NewClient("https://demo.nuxeo.com/nuxeo", &nuxeoClientOptions)

	// 5. use Nuxeo client
	serverVersion, err := nuxeoClient.ServerVersion(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Nuxeo Server Version: %d.%d.%d\n", serverVersion.Major, serverVersion.Minor, serverVersion.Patch)

	currentUser, err := nuxeoClient.CurrentUser(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Current User:", currentUser.Id)
	fmt.Println("E-Mail:", fmt.Sprintf("%v", currentUser.Email()))
}
