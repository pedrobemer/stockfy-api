package oauth2

import (
	"fmt"
	"net/url"
	"stockfyApi/client"
)

// Interface for the functions responsible for the OAuth2
type FacebookOAuthInterface interface {
	GrantAuthorizationUrl() string
	GrantAccessToken(authCode string) (FacebookOAuthResp, error)
}

// Google OAuth2 configuration
type ConfigFacebookOAuth2 struct {
	ClientID              string
	ClientSecret          string
	RedirectURI           string
	Scope                 []string
	AuthorizationEndpoint string
	TokenEndpoint         string
}

type FacebookOAuth2 struct {
	Interface FacebookOAuthInterface
	Config    ConfigFacebookOAuth2
}

// Information returned by the OAuth2 when our application receives the access
// token from the Google Server.
type FacebookOAuthResp struct {
	AccessToken string `json:"access_token,omitempty"`
	// IdToken          string `json:"id_token,omitempty"`
	Expiration int    `json:"expires_in,omitempty"`
	TokenType  string `json:"token_type,omitempty"`
	// Scope            string `json:"scope,omitempty"`
	RefreshToken string                 `json:"refresh_token,omitempty"`
	Error        map[string]interface{} `json:"error,omitempty"`
}

// Returns the struct responsible to point to the functions to apply OAuth2 for
// Google.
func FacebookOAuthConfig(clientId string, clientSecret string,
	redirectURI string, scope []string, authEndpoint string,
	tokenEndpoint string) *ConfigFacebookOAuth2 {
	return &ConfigFacebookOAuth2{
		ClientID:              clientId,
		ClientSecret:          clientSecret,
		RedirectURI:           redirectURI,
		Scope:                 scope,
		AuthorizationEndpoint: authEndpoint,
		TokenEndpoint:         tokenEndpoint,
	}
}

// Generates the OAuth2 URL to grant authorization for our application to use
// user information from Google.
func (f *ConfigFacebookOAuth2) GrantAuthorizationUrl() string {

	URL, _ := url.Parse(f.AuthorizationEndpoint)

	scope := appendScope(f.Scope)

	parameters := url.Values{}
	parameters.Add("client_id", f.ClientID)
	parameters.Add("scope", scope)
	parameters.Add("redirect_uri", f.RedirectURI)
	parameters.Add("response_type", "code")
	parameters.Add("state", "test")

	URL.RawQuery = parameters.Encode()
	authorizationUrl := URL.String()

	return authorizationUrl
}

// Generates the OAuth Access Token for the application to access the user
// resources granted in the scope.
func (f *ConfigFacebookOAuth2) GrantAccessToken(authCode string) (
	FacebookOAuthResp, error) {
	var facebookOAuthInfo FacebookOAuthResp

	URL, _ := url.Parse(f.TokenEndpoint)

	dataUrlFormMap := url.Values{
		"code":          {authCode},
		"client_id":     {f.ClientID},
		"client_secret": {f.ClientSecret},
		"redirect_uri":  {f.RedirectURI},
	}

	URL.RawQuery = dataUrlFormMap.Encode()
	tokenUrl := URL.String()

	client.RequestAndAssignToBody("GET", tokenUrl, "", nil, &facebookOAuthInfo)
	fmt.Println(facebookOAuthInfo)

	return facebookOAuthInfo, nil
}
