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

type MockGoogleOAuth2 struct {
	ClientID              string
	ClientSecret          string
	RedirectURI           string
	Scope                 []string
	AuthorizationEndpoint string
	TokenEndpoint         string
}

func MockGoogleOAuthConfig(clientId string, clientSecret string, redirectURI string,
	scope []string, authEndpoint string, tokenEndpoint string) *MockGoogleOAuth2 {
	return &MockGoogleOAuth2{
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
func (g *MockGoogleOAuth2) GrantAuthorizationUrl(state string) string {

	URL, _ := url.Parse(g.AuthorizationEndpoint)

	scope := appendScope(g.Scope)

	parameters := url.Values{}
	parameters.Add("client_id", g.ClientID)
	parameters.Add("scope", scope)
	parameters.Add("redirect_uri", g.RedirectURI)
	parameters.Add("response_type", "code")
	parameters.Add("state", state)

	URL.RawQuery = parameters.Encode()
	authorizationUrl := URL.String()

	return authorizationUrl
}

// Generates the OAuth Access Token for the application to access the user
// resources granted in the scope.
func (g *MockGoogleOAuth2) GrantAccessToken(authCode string) (
	oauth2.GoogleOAuthResp, error) {
	var googleOAuthInfo oauth2.GoogleOAuthResp

	dataUrlFormMap := url.Values{
		"code":          {authCode},
		"client_id":     {g.ClientID},
		"client_secret": {g.ClientSecret},
		"redirect_uri":  {g.RedirectURI},
		"grant_type":    {"authorization_code"},
	}
	dataUrlFormStr := dataUrlFormMap.Encode()

	// client.RequestAndAssignToBody("POST", g.TokenEndpoint,
	// 	"application/x-www-form-urlencoded", strings.NewReader(dataUrlFormStr),
	// 	&googleOAuthInfo)

	mockHttp := MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {

			type respBodyStruct struct {
				AccessToken      string `json:"access_token,omitempty"`
				IdToken          string `json:"id_token,omitempty"`
				Expiration       int    `json:"expires_in,omitempty"`
				TokenType        string `json:"token_type,omitempty"`
				Scope            string `json:"scope,omitempty"`
				RefreshToken     string `json:"refresh_token,omitempty"`
				Error            string `json:"error,omitempty"`
				ErrorDescription string `json:"error_description,omitempty"`
			}

			var code string
			bodyResp := respBodyStruct{}

			// Read body from the request
			body, _ := ioutil.ReadAll(req.Body)

			// Treat body from the request to get the code value from the URL query
			bodyQuery := strings.Split(string(body), "&")
			for _, query := range bodyQuery {
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
					Error:            "invalid_grant",
					ErrorDescription: "Bad Request",
				}
				break
			default:
				bodyResp = respBodyStruct{
					AccessToken:  "TestAccessToken",
					IdToken:      code,
					Expiration:   3600,
					TokenType:    "TestTokenType",
					Scope:        "TestScope",
					RefreshToken: "TestRefreshToken",
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
	resp, _ := mockHttp.MockHttpOutsideRequest("POST", g.TokenEndpoint,
		"application/x-www-form-urlencoded", strings.NewReader(dataUrlFormStr))

	bodyByte, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(bodyByte, &googleOAuthInfo)

	return googleOAuthInfo, nil
}
