package fiberHandlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"stockfyApi/api/middleware"
	"stockfyApi/api/presenter"
	"stockfyApi/entity"
	"stockfyApi/externalApi/oauth2"
	"stockfyApi/usecases"
	"stockfyApi/usecases/utils"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestApiUsersSignUp(t *testing.T) {

	type body struct {
		Success  bool                     `json:"success"`
		Message  string                   `json:"message"`
		Error    string                   `json:"error"`
		Code     int                      `json:"code"`
		UserInfo *presenter.UserApiReturn `json:"userInfo"`
	}

	type test struct {
		contentType  string
		bodyReq      presenter.SignUpBody
		expectedResp body
	}

	tests := []test{
		{
			contentType: "application/pdf",
			bodyReq: presenter.SignUpBody{
				Email:       "",
				Password:    "PasswdTest",
				DisplayName: "Test Username",
			},
			expectedResp: body{
				Success:  false,
				Message:  entity.ErrMessageApiRequest.Error(),
				Error:    entity.ErrInvalidApiBody.Error(),
				Code:     400,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.SignUpBody{
				Email:       "",
				Password:    "PasswdTest",
				DisplayName: "Test Username",
			},
			expectedResp: body{
				Success:  false,
				Message:  entity.ErrMessageApiRequest.Error(),
				Error:    errors.New("email must be a non-empty string").Error(),
				Code:     400,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.SignUpBody{
				Email:       "test@email.com",
				Password:    "",
				DisplayName: "Test Username",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error: errors.New("password must be a string at least 6 characters long").
					Error(),
				Code:     400,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.SignUpBody{
				Email:       "test@email.com",
				Password:    "PasswdTest",
				DisplayName: "",
			},
			expectedResp: body{
				Success:  false,
				Message:  entity.ErrMessageApiRequest.Error(),
				Error:    errors.New("display name must be a non-empty string").Error(),
				Code:     400,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.SignUpBody{
				Email:       "test@email.com",
				Password:    "PasswdTest",
				DisplayName: "",
			},
			expectedResp: body{
				Success:  false,
				Message:  entity.ErrMessageApiRequest.Error(),
				Error:    errors.New("display name must be a non-empty string").Error(),
				Code:     400,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.SignUpBody{
				Email:       "test@email.com",
				Password:    "PasswdTest",
				DisplayName: "WRONG_CUSTOM_TOKEN",
			},
			expectedResp: body{
				Success:  false,
				Message:  entity.ErrMessageApiInternalError.Error(),
				Error:    errors.New("Some Error").Error(),
				Code:     500,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.SignUpBody{
				Email:       "test@email.com",
				Password:    "PasswdTest",
				DisplayName: "WRONG_ID_TOKEN",
			},
			expectedResp: body{
				Success:  false,
				Message:  entity.ErrMessageApiInternalError.Error(),
				Error:    entity.ErrInvalidUserToken.Error(),
				Code:     500,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.SignUpBody{
				Email:       "test@email.com",
				Password:    "PasswdTest",
				DisplayName: "WRONG_EMAIL_VERIFICATION",
			},
			expectedResp: body{
				Success:  false,
				Message:  entity.ErrMessageApiRequest.Error(),
				Error:    "INVALID_ID_TOKEN",
				Code:     400,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.SignUpBody{
				Email:       "test@email.com",
				Password:    "PasswdTest",
				DisplayName: "WRONG_USER_INFO",
			},
			expectedResp: body{
				Success:  false,
				Message:  entity.ErrMessageApiInternalError.Error(),
				Error:    entity.ErrInvalidUserEmailBlank.Error(),
				Code:     500,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.SignUpBody{
				Email:       "test@email.com",
				Password:    "PasswdTest",
				DisplayName: "Test Username",
			},
			expectedResp: body{
				Success: true,
				Message: "User was registered successfully",
				Error:   "",
				Code:    200,
				UserInfo: &presenter.UserApiReturn{
					Email:       "test@email.com",
					DisplayName: "Test Username",
				},
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()

	// Declare Sector Application Logic
	users := UsersApi{
		ApplicationLogic: *usecases,
	}

	// Mock HTTP request
	app := fiber.New()
	api := app.Group("/api")
	api.Post("/signup", users.SignUp)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "POST", "/api/signup",
			testCase.contentType, "", testCase.bodyReq)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}

}

func TestApiUsersSignIn(t *testing.T) {

	type body struct {
		Success  bool                          `json:"success"`
		Message  string                        `json:"message"`
		Error    string                        `json:"error"`
		Code     int                           `json:"code"`
		UserInfo *presenter.UserLoginApiReturn `json:"userInfo"`
	}

	type test struct {
		contentType  string
		bodyReq      presenter.SignInBody
		expectedResp body
	}

	tests := []test{
		{
			contentType: "application/pdf",
			bodyReq: presenter.SignInBody{
				Email:    "",
				Password: "PasswdTest",
			},
			expectedResp: body{
				Success:  false,
				Message:  entity.ErrMessageApiRequest.Error(),
				Error:    entity.ErrInvalidApiBody.Error(),
				Code:     400,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.SignInBody{
				Email:    "",
				Password: "PasswdTest",
			},
			expectedResp: body{
				Success:  false,
				Message:  entity.ErrMessageApiRequest.Error(),
				Error:    "INVALID_EMAIL",
				Code:     400,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.SignInBody{
				Email:    "test@email.com",
				Password: "",
			},
			expectedResp: body{
				Success:  false,
				Message:  entity.ErrMessageApiRequest.Error(),
				Error:    "MISSING_PASSWORD",
				Code:     400,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.SignInBody{
				Email:    "test@email.com",
				Password: "PasswdTest",
			},
			expectedResp: body{
				Success: true,
				Message: "User login was successful",
				Error:   "",
				Code:    200,
				UserInfo: &presenter.UserLoginApiReturn{
					Email:        "test@email.com",
					DisplayName:  "Test User Name",
					IdToken:      "ValidIdToken",
					RefreshToken: "ValidRefreshToken",
					Expiration:   "3600",
				},
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()

	// Declare Sector Application Logic
	users := UsersApi{
		ApplicationLogic: *usecases,
	}

	// Mock HTTP request
	app := fiber.New()
	api := app.Group("/api")
	api.Post("/signin", users.SignIn)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "POST", "/api/signin",
			testCase.contentType, "", testCase.bodyReq)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}

}

func TestApiUsersSignInOAuth(t *testing.T) {

	clientId := "TestClientID"
	clientSecret := "TestClientSecret"
	redirectUri := "http://localhost:3000/api/signin/oauth2/google"
	googleScope := []string{
		"https://accounts.google.com/o/oauth2/auth",
		"https://oauth2.googleapis.com/token",
	}
	facebookScope := []string{
		"email",
		"public_profile",
	}
	googleAuthEndpoint := "https://accounts.google.com/o/oauth2/auth"
	facebookAuthEndpoint := "https://www.facebook.com/v12.0/dialog/oauth"
	googleTokenEndpoint := "https://oauth2.googleapis.com/token"
	facebookTokenEndpoint := "https://graph.facebook.com/v12.0/oauth/access_token"

	googleOAuth2Config := MockGoogleOAuthConfig(clientId, clientSecret,
		redirectUri, googleScope, googleAuthEndpoint, googleTokenEndpoint)
	facebookOAuth2Config := MockFacebookOAuthConfig(clientId, clientSecret,
		redirectUri, facebookScope, facebookAuthEndpoint, facebookTokenEndpoint)

	type body struct {
		Success  bool                          `json:"success"`
		Message  string                        `json:"message"`
		Error    string                        `json:"error"`
		Code     int                           `json:"code"`
		UserInfo *presenter.UserLoginApiReturn `json:"userInfo"`
	}

	type test struct {
		contentType  string
		urlQuery     string
		expectedResp body
		expectedURL  string
	}

	maker, _ := MockNewPasetoMaker("")
	pasetoToken, _ := maker.CreateToken("", time.Minute)

	tests := []test{
		{
			contentType: "application/json",
			urlQuery:    "type=google",
			expectedResp: body{
				Code: 302,
			},
			expectedURL: googleOAuth2Config.GrantAuthorizationUrl(pasetoToken),
		},
		{
			contentType: "application/json",
			urlQuery:    "type=facebook",
			expectedResp: body{
				Code: 302,
			},
			expectedURL: facebookOAuth2Config.GrantAuthorizationUrl(pasetoToken),
		},
		{
			contentType: "application/json",
			urlQuery:    "type=ERROR",
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidApiQueryLoginType.Error(),
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()

	// Declare Sector Application Logic
	users := UsersApi{
		ApplicationLogic: *usecases,
		GoogleOAuth2: oauth2.GoogleOAuth2{
			Interface: googleOAuth2Config,
			Config: oauth2.ConfigGoogleOAuth2{
				ClientID:              googleOAuth2Config.ClientID,
				ClientSecret:          googleOAuth2Config.ClientSecret,
				RedirectURI:           googleOAuth2Config.RedirectURI,
				Scope:                 googleOAuth2Config.Scope,
				AuthorizationEndpoint: googleOAuth2Config.AuthorizationEndpoint,
				TokenEndpoint:         googleOAuth2Config.TokenEndpoint,
			},
		},
		FacebookOAuth2: oauth2.FacebookOAuth2{
			Interface: facebookOAuth2Config,
			Config: oauth2.ConfigFacebookOAuth2{
				ClientID:              facebookOAuth2Config.ClientID,
				ClientSecret:          facebookOAuth2Config.ClientSecret,
				RedirectURI:           facebookOAuth2Config.RedirectURI,
				Scope:                 facebookOAuth2Config.Scope,
				AuthorizationEndpoint: facebookOAuth2Config.AuthorizationEndpoint,
				TokenEndpoint:         facebookOAuth2Config.TokenEndpoint,
			},
		},
		TokenMaker: MockNewPasetoMaker,
	}

	// Mock HTTP request
	app := fiber.New()
	api := app.Group("/api")
	api.Get("/signin", users.SignInOAuth)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "GET", "/api/signin?"+testCase.urlQuery,
			testCase.contentType, "", nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		if url, _ := resp.Location(); url != nil {
			assert.Equal(t, testCase.expectedURL, url.String())
		} else {
			assert.Equal(t, testCase.expectedURL, "")
		}

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}

}

func TestApiUsersGoogleOAuth2Redirect(t *testing.T) {

	// This function is the mock function to get the access token from the
	// Google OAuth2

	clientId := "TestClientID"
	clientSecret := "TestClientSecret"
	redirectUri := "http://localhost:3000/api/signin/oauth2/google"
	scope := []string{
		"https://accounts.google.com/o/oauth2/auth",
		"https://oauth2.googleapis.com/token",
	}
	authorizationEnpoint := "https://accounts.google.com/o/oauth2/auth"
	tokenEndpoint := "https://oauth2.googleapis.com/token"
	googleOAuth2Config := MockGoogleOAuthConfig(clientId, clientSecret,
		redirectUri, scope, authorizationEnpoint, tokenEndpoint)

	type body struct {
		Success  bool                          `json:"success"`
		Message  string                        `json:"message"`
		Error    string                        `json:"error"`
		Code     int                           `json:"code"`
		UserInfo *presenter.UserLoginApiReturn `json:"userInfo"`
	}

	type test struct {
		contentType   string
		urlParams     string
		urlQuery      string
		stateUsername string
		expectedResp  body
	}

	tests := []test{
		{
			contentType:   "application/json",
			urlParams:     "ERROR",
			stateUsername: "VALID_USERNAME",
			urlQuery:      "code=Test&state=VALID_USERNAME",
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidApiParamsCompany.Error(),
			},
		},
		{
			contentType:   "application/json",
			urlParams:     "google",
			urlQuery:      "code=Test",
			stateUsername: "VALID_USERNAME",
			expectedResp: body{
				Code:    403,
				Success: false,
				Message: entity.ErrMessageApiAuthorization.Error(),
				Error:   entity.ErrInvalidApiQueryStateBlank.Error(),
			},
		},
		{
			contentType:   "application/json",
			urlParams:     "google",
			urlQuery:      "code=Test&state=INVALID_TOKEN",
			stateUsername: "INVALID_TOKEN",
			expectedResp: body{
				Code:    403,
				Success: false,
				Message: entity.ErrMessageApiAuthorization.Error(),
				Error:   entity.ErrInvalidApiQueryState.Error() + "invalid token",
			},
		},
		{
			contentType:   "application/json",
			urlParams:     "google",
			urlQuery:      "code=Test&state=EXPIRED_TOKEN",
			stateUsername: "EXPIRED_TOKEN",
			expectedResp: body{
				Code:    403,
				Success: false,
				Message: entity.ErrMessageApiAuthorization.Error(),
				Error: entity.ErrInvalidApiQueryState.Error() +
					entity.ErrExpiredToken.Error(),
			},
		},
		{
			contentType:   "application/json",
			urlParams:     "google",
			urlQuery:      "code=Test&state=WRONG_STATE",
			stateUsername: "VALID_STATE",
			expectedResp: body{
				Code:    403,
				Success: false,
				Message: entity.ErrMessageApiAuthorization.Error(),
				Error:   entity.ErrInvalidApiQueryStateDoesNotMatch.Error(),
			},
		},
		{
			contentType:   "application/json",
			urlParams:     "google",
			urlQuery:      "state=VALID_USERNAME",
			stateUsername: "VALID_USERNAME",
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidApiQueryOAuth2CodeBlank.Error(),
			},
		},
		{
			contentType:   "application/json",
			urlParams:     "google",
			urlQuery:      "code=INVALID_CODE&state=VALID_USERNAME",
			stateUsername: "VALID_USERNAME",
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   "INVALID_GRANT",
			},
		},
		{
			contentType:   "application/json",
			urlParams:     "google",
			urlQuery:      "code=ERROR_IDP_RESPONSE&state=VALID_USERNAME",
			stateUsername: "VALID_USERNAME",
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   errors.New("INVALID_IDP_RESPONSE").Error(),
			},
		},
		{
			contentType:   "application/json",
			urlParams:     "google",
			urlQuery:      "code=NEW_USER_WITHOUT_EMAIL&state=VALID_USERNAME",
			stateUsername: "VALID_USERNAME",
			expectedResp: body{
				Code:    500,
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Error:   entity.ErrInvalidUserEmailBlank.Error(),
			},
		},
		{
			contentType:   "application/json",
			urlParams:     "google",
			urlQuery:      "code=TestCode&state=VALID_USERNAME",
			stateUsername: "VALID_USERNAME",
			expectedResp: body{
				Code:    200,
				Success: true,
				Message: "User login was successful",
				Error:   "",
				UserInfo: &presenter.UserLoginApiReturn{
					Email:        "test@email.com",
					DisplayName:  "Test Name",
					IdToken:      "ValidIdTokenWithoutPrivilegedUser",
					RefreshToken: "ValidRefreshToken",
					Expiration:   "3600",
				},
			},
		},
	}

	tokenMaker, err := MockNewPasetoMaker(utils.RandString(32))
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenMaker)

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()

	// Declare Sector Application Logic
	users := UsersApi{
		ApplicationLogic: *usecases,
		GoogleOAuth2: oauth2.GoogleOAuth2{
			Interface: googleOAuth2Config,
			Config: oauth2.ConfigGoogleOAuth2{
				ClientID:              googleOAuth2Config.ClientID,
				ClientSecret:          googleOAuth2Config.ClientSecret,
				RedirectURI:           googleOAuth2Config.RedirectURI,
				Scope:                 googleOAuth2Config.Scope,
				AuthorizationEndpoint: googleOAuth2Config.AuthorizationEndpoint,
				TokenEndpoint:         googleOAuth2Config.TokenEndpoint,
			},
		},
		TokenMaker:          MockNewPasetoMaker,
		StateTokenInterface: tokenMaker,
	}

	// Mock HTTP request
	app := fiber.New()
	api := app.Group("/api")
	api.Get("/signin/oauth2/:company", users.OAuth2Redirect)

	for _, testCase := range tests {
		users.StateUsername = testCase.stateUsername
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "GET", "/api/signin/oauth2/"+
			testCase.urlParams+"?"+testCase.urlQuery, testCase.contentType, "",
			nil)
		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}

}

func TestApiUsersFacebookOAuth2Redirect(t *testing.T) {

	// This function is the mock function to get the access token from the
	// Google OAuth2

	clientId := "TestClientID"
	clientSecret := "TestClientSecret"
	redirectUri := "http://localhost:3000/api/signin/oauth2/facebook"
	scope := []string{
		"email",
		"public_profile",
	}
	authorizationEnpoint := "https://www.facebook.com/v12.0/dialog/oauth"
	tokenEndpoint := "https://graph.facebook.com/v12.0/oauth/access_token"
	facebookOAuth2Config := MockFacebookOAuthConfig(clientId, clientSecret,
		redirectUri, scope, authorizationEnpoint, tokenEndpoint)

	type body struct {
		Success  bool                          `json:"success"`
		Message  string                        `json:"message"`
		Error    string                        `json:"error"`
		Code     int                           `json:"code"`
		UserInfo *presenter.UserLoginApiReturn `json:"userInfo"`
	}

	type test struct {
		contentType   string
		stateUsername string
		urlParams     string
		urlQuery      string
		expectedResp  body
	}

	tests := []test{
		{
			contentType:   "application/json",
			urlParams:     "facebook",
			urlQuery:      "state=VALID_USERNAME",
			stateUsername: "VALID_USERNAME",
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidApiQueryOAuth2CodeBlank.Error(),
			},
		},
		{
			contentType:   "application/json",
			urlParams:     "facebook",
			urlQuery:      "code=INVALID_CODE&state=VALID_USERNAME",
			stateUsername: "VALID_USERNAME",
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   "ERROR",
			},
		},
		{
			contentType:   "application/json",
			urlParams:     "facebook",
			urlQuery:      "code=ERROR_IDP_RESPONSE&state=VALID_USERNAME",
			stateUsername: "VALID_USERNAME",
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   errors.New("INVALID_IDP_RESPONSE").Error(),
			},
		},
		{
			contentType:   "application/json",
			urlParams:     "facebook",
			urlQuery:      "code=TestCode&state=VALID_USERNAME",
			stateUsername: "VALID_USERNAME",
			expectedResp: body{
				Code:    200,
				Success: true,
				Message: "User login was successful",
				Error:   "",
				UserInfo: &presenter.UserLoginApiReturn{
					Email:        "test@email.com",
					DisplayName:  "Test Name",
					IdToken:      "ValidIdTokenWithoutPrivilegedUser",
					RefreshToken: "ValidRefreshToken",
					Expiration:   "3600",
				},
			},
		},
	}

	// Mock Token Maker
	tokenMaker, err := MockNewPasetoMaker(utils.RandString(32))
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenMaker)

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()

	// Declare Sector Application Logic
	users := UsersApi{
		ApplicationLogic: *usecases,
		FacebookOAuth2: oauth2.FacebookOAuth2{
			Interface: facebookOAuth2Config,
			Config: oauth2.ConfigFacebookOAuth2{
				ClientID:              facebookOAuth2Config.ClientID,
				ClientSecret:          facebookOAuth2Config.ClientSecret,
				RedirectURI:           facebookOAuth2Config.RedirectURI,
				Scope:                 facebookOAuth2Config.Scope,
				AuthorizationEndpoint: facebookOAuth2Config.AuthorizationEndpoint,
				TokenEndpoint:         facebookOAuth2Config.TokenEndpoint,
			},
		},
		TokenMaker:          MockNewPasetoMaker,
		StateTokenInterface: tokenMaker,
	}

	// Mock HTTP request
	app := fiber.New()
	api := app.Group("/api")
	api.Get("/signin/oauth2/:company", users.OAuth2Redirect)

	for _, testCase := range tests {
		users.StateUsername = testCase.stateUsername
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "GET", "/api/signin/oauth2/"+
			testCase.urlParams+"?"+testCase.urlQuery, testCase.contentType, "",
			nil)
		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}

}

func TestApiForgotPassword(t *testing.T) {
	type body struct {
		Success  bool                                `json:"success"`
		Message  string                              `json:"message"`
		Error    string                              `json:"error"`
		Code     int                                 `json:"code"`
		UserInfo *entity.EmailForgotPasswordResponse `json:"userInfo"`
	}

	type test struct {
		contentType  string
		bodyReq      presenter.ForgotPasswordBody
		expectedResp body
	}

	tests := []test{
		{
			contentType: "application/pdf",
			bodyReq: presenter.ForgotPasswordBody{
				Email: "",
			},
			expectedResp: body{
				Success:  false,
				Message:  entity.ErrMessageApiRequest.Error(),
				Error:    entity.ErrInvalidApiBody.Error(),
				Code:     400,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.ForgotPasswordBody{
				Email: "",
			},
			expectedResp: body{
				Success:  false,
				Message:  entity.ErrMessageApiRequest.Error(),
				Error:    "MISSING_EMAIL",
				Code:     400,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.ForgotPasswordBody{
				Email: "INVALID_EMAIL",
			},
			expectedResp: body{
				Success:  false,
				Message:  entity.ErrMessageApiEmail.Error(),
				Error:    "",
				Code:     404,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.ForgotPasswordBody{
				Email: "test@email.com",
			},
			expectedResp: body{
				Success: true,
				Message: "The email for password reset was sent successfully",
				Error:   "",
				Code:    200,
				UserInfo: &entity.EmailForgotPasswordResponse{
					Email: "test@email.com",
				},
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()

	// Declare Sector Application Logic
	users := UsersApi{
		ApplicationLogic: *usecases,
	}

	// Mock HTTP request
	app := fiber.New()
	api := app.Group("/api")
	api.Post("/forgot-password", users.ForgotPassword)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "POST", "/api/forgot-password",
			testCase.contentType, "", testCase.bodyReq)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}

}

func TestApiRefreshIdToken(t *testing.T) {
	type body struct {
		Success   bool                                 `json:"success"`
		Message   string                               `json:"message"`
		Error     string                               `json:"error"`
		Code      int                                  `json:"code"`
		UserToken *presenter.UserRefreshTokenApiReturn `json:"userRefreshToken"`
	}

	type test struct {
		idToken      string
		contentType  string
		bodyReq      presenter.UserRefreshIdTokenBody
		expectedResp body
	}

	tests := []test{
		{
			idToken:     "ValidIdTokenWithoutEmailVerification",
			contentType: "application/json",
			expectedResp: body{
				Code:      401,
				Success:   false,
				Message:   entity.ErrMessageApiAuthentication.Error(),
				Error:     "",
				UserToken: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/pdf",
			expectedResp: body{
				Code:      400,
				Success:   false,
				Message:   entity.ErrMessageApiRequest.Error(),
				Error:     entity.ErrInvalidApiBody.Error(),
				UserToken: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			bodyReq: presenter.UserRefreshIdTokenBody{
				RefreshToken: "",
			},
			expectedResp: body{
				Code:      400,
				Success:   false,
				Message:   entity.ErrMessageApiRequest.Error(),
				Error:     "MISSING_REFRESH_TOKEN",
				UserToken: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			bodyReq: presenter.UserRefreshIdTokenBody{
				RefreshToken: "UNKNOWN_REFRESH_TOKEN",
			},
			expectedResp: body{
				Code:      400,
				Success:   false,
				Message:   entity.ErrMessageApiRequest.Error(),
				Error:     "INVALID_REFRESH_TOKEN",
				UserToken: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			bodyReq: presenter.UserRefreshIdTokenBody{
				RefreshToken: "ValidRefreshToken",
			},
			expectedResp: body{
				Code:    200,
				Success: true,
				Message: "The token was updated successfully",
				Error:   "",
				UserToken: &presenter.UserRefreshTokenApiReturn{
					RefreshToken: "ValidRefreshToken",
					IdToken:      "ValidIdToken",
					TokenType:    "Bearer",
					Expiration:   "3600",
				},
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()
	// logicApi := logicApi.NewMockApplication(*usecases)

	// Declare Sector Application Logic
	users := UsersApi{
		ApplicationLogic: *usecases,
	}

	// Mock HTTP request
	app := fiber.New()
	api := app.Group("/api")
	api.Use(middleware.NewFiberMiddleware(middleware.FiberMiddleware{
		UserAuthentication: usecases.UserApp,
		ErrorHandler: func(c *fiber.Ctx, e error) error {
			var err error
			c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": entity.ErrMessageApiAuthentication.Error(),
				"code":    401,
			})

			return err
		},
		ContextKey: "user",
	}))
	api.Post("/refresh-token", users.RefreshIdToken)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "POST", "/api/refresh-token",
			testCase.contentType, testCase.idToken, testCase.bodyReq)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}
}

func TestApiDeleteUser(t *testing.T) {
	type body struct {
		Success  bool             `json:"success"`
		Message  string           `json:"message"`
		Error    string           `json:"error"`
		Code     int              `json:"code"`
		UserInfo *entity.UserInfo `json:"userInfo"`
	}

	type test struct {
		idToken      string
		contentType  string
		expectedResp body
	}

	tests := []test{
		{
			idToken:     "ValidIdTokenWithoutEmailVerification",
			contentType: "application/json",
			expectedResp: body{
				Code:     401,
				Success:  false,
				Message:  entity.ErrMessageApiAuthentication.Error(),
				Error:    "",
				UserInfo: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutRegister",
			contentType: "application/json",
			expectedResp: body{
				Code:     400,
				Success:  false,
				Message:  entity.ErrMessageApiRequest.Error(),
				Error:    errors.New("Invalid user UID").Error(),
				UserInfo: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			expectedResp: body{
				Code:    200,
				Success: true,
				Message: "User was deleted successfully",
				Error:   "",
				UserInfo: &entity.UserInfo{
					UID:         "TestUID",
					Email:       "test@email.com",
					DisplayName: "Test Name",
				},
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()
	// logicApi := logicApi.NewMockApplication(*usecases)

	// Declare Sector Application Logic
	users := UsersApi{
		ApplicationLogic: *usecases,
	}

	// Mock HTTP request
	app := fiber.New()
	api := app.Group("/api")
	api.Use(middleware.NewFiberMiddleware(middleware.FiberMiddleware{
		UserAuthentication: usecases.UserApp,
		ErrorHandler: func(c *fiber.Ctx, e error) error {
			var err error
			c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": entity.ErrMessageApiAuthentication.Error(),
				"code":    401,
			})

			return err
		},
		ContextKey: "user",
	}))
	api.Delete("/delete-user", users.DeleteUser)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "DELETE", "/api/delete-user",
			testCase.contentType, testCase.idToken, nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}

}

func TestApiUpdateUserInfo(t *testing.T) {
	type body struct {
		Success  bool                     `json:"success"`
		Message  string                   `json:"message"`
		Error    string                   `json:"error"`
		Code     int                      `json:"code"`
		UserInfo *presenter.UserApiReturn `json:"userInfo"`
	}

	type test struct {
		idToken      string
		contentType  string
		bodyReq      presenter.SignUpBody
		expectedResp body
	}

	tests := []test{
		{
			idToken:     "ValidIdTokenWithoutEmailVerification",
			contentType: "application/json",
			expectedResp: body{
				Code:     401,
				Success:  false,
				Message:  entity.ErrMessageApiAuthentication.Error(),
				Error:    "",
				UserInfo: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/pdf",
			expectedResp: body{
				Code:     400,
				Success:  false,
				Message:  entity.ErrMessageApiRequest.Error(),
				Error:    entity.ErrInvalidApiBody.Error(),
				UserInfo: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutRegister",
			contentType: "application/json",
			expectedResp: body{
				Code:     400,
				Success:  false,
				Message:  entity.ErrMessageApiRequest.Error(),
				Error:    errors.New("INVALID_USER_UID").Error(),
				UserInfo: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			bodyReq: presenter.SignUpBody{
				Email:       "testChange@email.com",
				DisplayName: "Test Change",
			},
			expectedResp: body{
				Code:    200,
				Success: true,
				Message: "User information was updated successfully",
				Error:   "",
				UserInfo: &presenter.UserApiReturn{
					Email:       "testChange@email.com",
					DisplayName: "Test Change",
				},
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()
	// logicApi := logicApi.NewMockApplication(*usecases)

	// Declare Sector Application Logic
	users := UsersApi{
		ApplicationLogic: *usecases,
	}

	// Mock HTTP request
	app := fiber.New()
	api := app.Group("/api")
	api.Use(middleware.NewFiberMiddleware(middleware.FiberMiddleware{
		UserAuthentication: usecases.UserApp,
		ErrorHandler: func(c *fiber.Ctx, e error) error {
			var err error
			c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": entity.ErrMessageApiAuthentication.Error(),
				"code":    401,
			})

			return err
		},
		ContextKey: "user",
	}))
	api.Put("/update-user", users.UpdateUserInfo)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "PUT", "/api/update-user",
			testCase.contentType, testCase.idToken, testCase.bodyReq)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}
}
