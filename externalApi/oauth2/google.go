package oauth2

import (
	"net/url"
	"stockfyApi/client"
	"strings"
)

// Interface for the functions responsible for the OAuth2
type GoogleOAuthInterface interface {
	GrantAuthorizationUrl() string
	GrantAccessToken(authCode string) (GoogleOAuthResp, error)
}

// Information returned by the OAuth2 when our application receives the access
// token from the Google Server.
type GoogleOAuthResp struct {
	AccessToken      string `json:"access_token,omitempty"`
	IdToken          string `json:"id_token,omitempty"`
	Expiration       int    `json:"expires_in,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	Scope            string `json:"scope,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

// Google OAuth2 configuration
type ConfigGoogleOAuth2 struct {
	ClientID              string
	ClientSecret          string
	RedirectURI           string
	Scope                 []string
	AuthorizationEndpoint string
	TokenEndpoint         string
}

type GoogleOAuth2 struct {
	Interface GoogleOAuthInterface
	Config    ConfigGoogleOAuth2
}

// Returns the struct responsible to point to the functions to apply OAuth2 for
// Google.
func GoogleOAuthConfig(clientId string, clientSecret string, redirectURI string,
	scope []string, authEndpoint string, tokenEndpoint string) *ConfigGoogleOAuth2 {
	return &ConfigGoogleOAuth2{
		ClientID:              clientId,
		ClientSecret:          clientSecret,
		RedirectURI:           redirectURI,
		Scope:                 scope,
		AuthorizationEndpoint: authEndpoint,
		TokenEndpoint:         tokenEndpoint,
	}
}

func appendScope(scopes []string) string {
	var appendScope string

	for _, scope := range scopes {
		appendScope = appendScope + scope + " "
	}

	return appendScope
}

// Generates the OAuth2 URL to grant authorization for our application to use
// user information from Google.
func (g *ConfigGoogleOAuth2) GrantAuthorizationUrl() string {

	URL, _ := url.Parse(g.AuthorizationEndpoint)

	scope := appendScope(g.Scope)

	parameters := url.Values{}
	parameters.Add("client_id", g.ClientID)
	parameters.Add("scope", scope)
	parameters.Add("redirect_uri", g.RedirectURI)
	parameters.Add("response_type", "code")

	URL.RawQuery = parameters.Encode()
	authorizationUrl := URL.String()

	return authorizationUrl
}

// Generates the OAuth Access Token for the application to access the user
// resources granted in the scope.
func (g *ConfigGoogleOAuth2) GrantAccessToken(authCode string) (GoogleOAuthResp, error) {
	var googleOAuthInfo GoogleOAuthResp

	dataUrlFormMap := url.Values{
		"code":          {authCode},
		"client_id":     {g.ClientID},
		"client_secret": {g.ClientSecret},
		"redirect_uri":  {g.RedirectURI},
		"grant_type":    {"authorization_code"},
	}
	dataUrlFormStr := dataUrlFormMap.Encode()

	client.RequestAndAssignToBody("POST", g.TokenEndpoint,
		"application/x-www-form-urlencoded", strings.NewReader(dataUrlFormStr),
		&googleOAuthInfo)

	return googleOAuthInfo, nil
}
