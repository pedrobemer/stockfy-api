package fiberHandlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"stockfyApi/externalApi/oauth2"
	"strings"
)

// Google OAuth2 configuration
type MockConfigFacebookOAuth2 struct {
	ClientID              string
	ClientSecret          string
	RedirectURI           string
	Scope                 []string
	AuthorizationEndpoint string
	TokenEndpoint         string
}

// Returns the struct responsible to point to the functions to apply OAuth2 for
// Google.
func MockFacebookOAuthConfig(clientId string, clientSecret string,
	redirectURI string, scope []string, authEndpoint string,
	tokenEndpoint string) *MockConfigFacebookOAuth2 {
	return &MockConfigFacebookOAuth2{
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
func (f *MockConfigFacebookOAuth2) GrantAuthorizationUrl(state string) string {

	URL, _ := url.Parse(f.AuthorizationEndpoint)

	scope := appendScope(f.Scope)

	parameters := url.Values{}
	parameters.Add("client_id", f.ClientID)
	parameters.Add("scope", scope)
	parameters.Add("redirect_uri", f.RedirectURI)
	parameters.Add("response_type", "code")
	parameters.Add("state", state)

	URL.RawQuery = parameters.Encode()
	authorizationUrl := URL.String()

	return authorizationUrl
}

// Generates the OAuth Access Token for the application to access the user
// resources granted in the scope.
func (f *MockConfigFacebookOAuth2) GrantAccessToken(authCode string) (
	oauth2.FacebookOAuthResp, error) {
	var facebookOAuthInfo oauth2.FacebookOAuthResp

	URL, _ := url.Parse(f.TokenEndpoint)

	dataUrlFormMap := url.Values{
		"code":          {authCode},
		"client_id":     {f.ClientID},
		"client_secret": {f.ClientSecret},
		"redirect_uri":  {f.RedirectURI},
	}

	URL.RawQuery = dataUrlFormMap.Encode()
	tokenUrl := URL.String()

	// client.RequestAndAssignToBody("GET", tokenUrl, "", nil, &facebookOAuthInfo)
	// fmt.Println(facebookOAuthInfo)
	mockHttp := MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {

			type respBodyStruct struct {
				AccessToken string                 `json:"access_token,omitempty"`
				Expiration  int                    `json:"expires_in,omitempty"`
				TokenType   string                 `json:"token_type,omitempty"`
				Error       map[string]interface{} `json:"error,omitempty"`
			}

			var code string
			bodyResp := respBodyStruct{}

			// Treat body from the request to get the code value from the URL query
			urlQuery := strings.Split(req.URL.RawQuery, "&")
			for _, query := range urlQuery {
				queryParams := strings.Split(string(query), "=")

				if queryParams[0] == "code" {
					code = queryParams[1]
				}
			}

			// If the code is invalid then return error, else returns the information
			// with the access token from the OAuth2
			switch code {
			case "INVALID_CODE":
				bodyResp = respBodyStruct{
					Error: map[string]interface{}{
						"message":       "This authorization code has expired.",
						"type":          "OAuthException",
						"code":          100,
						"error_subcode": 36007,
						"fbtrace_id":    "AbMyb8fkwlpI97SBNg7eWoY",
					},
				}
				break
			default:
				bodyResp = respBodyStruct{
					AccessToken: code,
					Expiration:  3600,
					TokenType:   "TestTokenType",
				}
			}

			bodyByte, _ := json.Marshal(bodyResp)

			respHeader := http.Header{
				"Content-Type": {"application/json"},
			}
			return &http.Response{
				Status:     "200 OK",
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     respHeader,
				Body:       ioutil.NopCloser(bytes.NewReader(bodyByte)),
				Request:    req,
			}, nil
		},
	}

	resp, _ := mockHttp.MockHttpOutsideRequest("GET", tokenUrl, "", nil)

	bodyByte, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(bodyByte, &facebookOAuthInfo)

	return facebookOAuthInfo, nil
}
