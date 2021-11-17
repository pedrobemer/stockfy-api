package user

import (
	"errors"
	"stockfyApi/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {

	type test struct {
		uid         string
		email       string
		displayName string
		// password            string
		userType            string
		expectedUserCreated *entity.Users
		expectedError       error
	}

	tests := []test{
		{
			uid:         "93avpow384",
			email:       "test@gmail.com",
			displayName: "Test Name",
			userType:    "normal",
			expectedUserCreated: &entity.Users{
				Id:       "39148-38149v-jk48",
				Username: "Test Name",
				Email:    "test@gmail.com",
				Uid:      "93avpow384",
				Type:     "normal",
			},
			expectedError: nil,
		},
		{
			uid:                 "93avpow384",
			email:               "ERROR_USER_REPOSITORY",
			displayName:         "Test Name",
			userType:            "normal",
			expectedUserCreated: nil,
			expectedError:       errors.New("Unknown error in the user repository"),
		},
		{
			email:               "test@gmail.com",
			displayName:         "Test Name",
			userType:            "normal",
			expectedUserCreated: nil,
			expectedError:       entity.ErrInvalidUserUidBlank,
		},
	}

	mockedRepo := NewMockRepo()
	mockedExternalApi := NewExternalApi()

	assetApp := NewApplication(mockedRepo, mockedExternalApi)

	for _, testCase := range tests {
		userCreated, err := assetApp.CreateUser(testCase.uid, testCase.email,
			testCase.displayName, testCase.userType)
		assert.Equal(t, testCase.expectedError, err)
		assert.Equal(t, testCase.expectedUserCreated, userCreated)
	}

}

func TestDeleteUser(t *testing.T) {
	type test struct {
		userUid          string
		expectedUserInfo *entity.UserInfo
		expectedError    error
	}

	tests := []test{
		{
			userUid: "8qjd340",
			expectedUserInfo: &entity.UserInfo{
				DisplayName: "Test Name",
				Email:       "test@gmail.com",
				UID:         "8qjd340",
			},
			expectedError: nil,
		},
		{
			userUid:          "ERROR_USER_REPOSITORY",
			expectedUserInfo: nil,
			expectedError:    errors.New("Unknown delete error in the user repository"),
		},
		{
			userUid:          "Invalid",
			expectedUserInfo: nil,
			expectedError:    errors.New("Database Interface error"),
		},
	}

	mockedRepo := NewMockRepo()
	mockedExternalApi := NewExternalApi()
	assetApp := NewApplication(mockedRepo, mockedExternalApi)

	for _, testCase := range tests {
		deletedUser, err := assetApp.DeleteUser(testCase.userUid)
		assert.Equal(t, testCase.expectedError, err)
		assert.Equal(t, testCase.expectedUserInfo, deletedUser)
	}

}

func TestUpdateUser(t *testing.T) {
	type test struct {
		userUid            string
		email              string
		displayName        string
		password           string
		expectedUserUpdate *entity.Users
		expectedError      error
	}

	tests := []test{
		{
			userUid:     "49qadkd0",
			email:       "test@gmail.com",
			displayName: "Test Name",
			password:    "test",
			expectedUserUpdate: &entity.Users{
				Id:       "391ahb4",
				Username: "Test Name",
				Email:    "test@gmail.com",
				Uid:      "49qadkd0",
				Type:     "normal",
			},
			expectedError: nil,
		},
		{
			userUid:  "49qadkd0",
			email:    "test2@gmail.com",
			password: "test",
			expectedUserUpdate: &entity.Users{
				Id:       "391ahb4",
				Username: "Test Name",
				Email:    "test2@gmail.com",
				Uid:      "49qadkd0",
				Type:     "normal",
			},
			expectedError: nil,
		},
		{
			userUid:     "49qadkd0",
			displayName: "Test Name 2",
			password:    "test",
			expectedUserUpdate: &entity.Users{
				Id:       "391ahb4",
				Username: "Test Name 2",
				Email:    "test@gmail.com",
				Uid:      "49qadkd0",
				Type:     "normal",
			},
			expectedError: nil,
		},
		{
			userUid:            "49qadkd0",
			displayName:        "ERROR_USER_FIREBASE",
			email:              "test@gmail.com",
			password:           "test",
			expectedUserUpdate: nil,
			expectedError:      errors.New("Unknown update error in the user repository"),
		},
		{
			userUid:            "49qadkd0",
			displayName:        "ERROR_USER_REPOSITORY",
			email:              "test@gmail.com",
			password:           "test",
			expectedUserUpdate: nil,
			expectedError:      errors.New("Unknown update error in the user repository"),
		},
	}

	mockedRepo := NewMockRepo()
	mockedExternalApi := NewExternalApi()
	assetApp := NewApplication(mockedRepo, mockedExternalApi)

	for _, testCase := range tests {
		userUpdated, err := assetApp.UpdateUser(testCase.userUid, testCase.email,
			testCase.displayName, testCase.password)
		assert.Equal(t, testCase.expectedUserUpdate, userUpdated)
		assert.Equal(t, testCase.expectedError, err)
	}
}

func TestSearchUser(t *testing.T) {
	type test struct {
		userUid          string
		expectedUserInfo *entity.Users
		expectedError    error
	}

	tests := []test{
		{
			userUid: "TestID",
			expectedUserInfo: &entity.Users{
				Uid:      "TestID",
				Email:    "test@gmail.com",
				Username: "Test Name",
				Type:     "normal",
			},
			expectedError: nil,
		},
		{
			userUid:          "Invalid",
			expectedUserInfo: nil,
			expectedError:    entity.ErrInvalidUserSearch,
		},
	}

	mockedRepo := NewMockRepo()
	assetApp := NewApplication(mockedRepo, nil)

	for _, testCase := range tests {
		searchedUser, err := assetApp.SearchUser(testCase.userUid)
		assert.Equal(t, testCase.expectedUserInfo, searchedUser)
		assert.Equal(t, testCase.expectedError, err)
	}
}

func TestUserCreate(t *testing.T) {
	type test struct {
		email            string
		password         string
		displayName      string
		expectedUserInfo *entity.UserInfo
		expectedError    error
	}

	tests := []test{
		{
			email:       "test@gmail.com",
			password:    "testando",
			displayName: "Test Name",
			expectedUserInfo: &entity.UserInfo{
				DisplayName: "Test Name",
				Email:       "test@gmail.com",
				UID:         "abj39as$$",
			},
			expectedError: nil,
		},
		{
			email:            "Error",
			password:         "testando",
			displayName:      "Test Name",
			expectedUserInfo: nil,
			expectedError:    errors.New("Error Mock Firebase"),
		},
	}

	mockedExtApi := NewExternalApi()
	assetApp := NewApplication(nil, mockedExtApi)

	for _, testCase := range tests {
		userInfo, err := assetApp.UserCreate(testCase.email, testCase.password,
			testCase.displayName)
		assert.Equal(t, testCase.expectedError, err)
		assert.Equal(t, testCase.expectedUserInfo, userInfo)
	}
}

func TestUserCreateCustomToken(t *testing.T) {
	expectedCustomToken := "194nc4850d"

	mockedExtApi := NewExternalApi()
	assetApp := NewApplication(nil, mockedExtApi)

	customToken, err := assetApp.UserCreateCustomToken("38qdasja")

	assert.Nil(t, err)
	assert.Equal(t, expectedCustomToken, customToken)
}

func TestUserRequestIdToken(t *testing.T) {
	type test struct {
		webKey            string
		customToken       string
		expectedTokenInfo *entity.ReqIdToken
		expectedError     error
	}

	tests := []test{
		{
			webKey:            "TestKey",
			customToken:       "49292",
			expectedTokenInfo: nil,
			expectedError:     entity.ErrInvalidUserToken,
		},
		{
			webKey:      "TestKey",
			customToken: "1acn49",
			expectedTokenInfo: &entity.ReqIdToken{
				Token:              "a419148a",
				RequestSecureToken: true,
				Kind:               "aa93q8",
				IdToken:            "294akfsnf49",
				IsNewUser:          false,
			},
			expectedError: nil,
		},
	}

	mockedExtApi := NewExternalApi()
	assetApp := NewApplication(nil, mockedExtApi)

	for _, testCase := range tests {
		userTokenInfo, err := assetApp.UserRequestIdToken(testCase.webKey,
			testCase.customToken)
		assert.Equal(t, testCase.expectedError, err)
		assert.Equal(t, testCase.expectedTokenInfo, userTokenInfo)

	}
}

func TestUserSendVerificationEmail(t *testing.T) {
	type test struct {
		webKey              string
		userIdToken         string
		expectedApiResponse entity.EmailVerificationResponse
		expectedError       error
	}

	err := map[string]interface{}{
		"code": 400,
		"errors": struct {
			domain  string
			message string
			reason  string
		}{"global", "INVALID_ID_TOKEN", "invalid"},
		"message": "INVALID_ID_TOKEN",
	}

	tests := []test{
		{
			webKey:      "TestKey",
			userIdToken: "ak4jaf49",
			expectedApiResponse: entity.EmailVerificationResponse{
				UserIdToken: "ak4jaf49",
				Email:       "test@gmail.com",
				Error:       nil,
			},
			expectedError: nil,
		},
		{
			webKey:      "TestKey",
			userIdToken: "Invalid",
			expectedApiResponse: entity.EmailVerificationResponse{
				UserIdToken: "Invalid",
				Error:       err,
			},
			expectedError: entity.ErrInvalidUserSendEmail,
		},
	}

	mockedExtApi := NewExternalApi()
	userApp := NewApplication(nil, mockedExtApi)

	for _, testCase := range tests {
		emailResponse, err := userApp.UserSendVerificationEmail(testCase.webKey,
			testCase.userIdToken)

		assert.Equal(t, testCase.expectedError, err)
		assert.Equal(t, testCase.expectedApiResponse, emailResponse)

	}

}

func TestUserSendForgotPasswordEmail(t *testing.T) {
	type test struct {
		webKey              string
		email               string
		expectedApiResponse entity.EmailForgotPasswordResponse
		expectedError       error
	}

	err := map[string]interface{}{
		"code": 400,
		"errors": struct {
			domain  string
			message string
			reason  string
		}{"global", "EMAIL_NOT_FOUND", "invalid"},
		"message": "EMAIL_NOT_FOUND",
	}

	tests := []test{
		{
			webKey: "9948cdi49ac",
			email:  "test@gmail.com",
			expectedApiResponse: entity.EmailForgotPasswordResponse{
				Email: "test@gmail.com",
				Error: nil,
			},
			expectedError: nil,
		},
		{
			webKey: "9948cdi49ac",
			email:  "Invalid",
			expectedApiResponse: entity.EmailForgotPasswordResponse{
				Email: "Invalid",
				Error: err,
			},
			expectedError: entity.ErrInvalidUserSendEmail,
		},
	}

	mockedExtApi := NewExternalApi()
	userApp := NewApplication(nil, mockedExtApi)

	for _, testCase := range tests {
		emailResetPasswdResp, err := userApp.UserSendForgotPasswordEmail(
			testCase.webKey, testCase.email)

		assert.Equal(t, testCase.expectedApiResponse, emailResetPasswdResp)
		assert.Equal(t, testCase.expectedError, err)
	}

}

func TestUserTokenVerification(t *testing.T) {
	type test struct {
		idToken               string
		expectedUserTokenInfo *entity.UserTokenInfo
		expectedError         error
	}

	tests := []test{
		{
			idToken:               "INVALID_ID_TOKEN",
			expectedUserTokenInfo: nil,
			expectedError:         errors.New("INVALID_ID_TOKEN"),
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			expectedUserTokenInfo: &entity.UserTokenInfo{
				Email:         "test@email.com",
				EmailVerified: true,
				UserID:        "TestUserID",
			},
			expectedError: nil,
		},
	}

	mockedExtApi := NewExternalApi()
	userApp := NewApplication(nil, mockedExtApi)

	for _, testCase := range tests {
		userTokenInfo, err := userApp.UserTokenVerification(testCase.idToken)

		assert.Equal(t, testCase.expectedError, err)
		assert.Equal(t, testCase.expectedUserTokenInfo, userTokenInfo)
	}

}

func TestUserLogin(t *testing.T) {
	type test struct {
		email             string
		password          string
		expectedUserLogin *entity.UserLoginResponse
		expectedError     error
	}

	tests := []test{
		{
			email:             "",
			password:          "",
			expectedUserLogin: nil,
			expectedError:     errors.New("INVALID_EMAIL"),
		},
		{
			email:             "test@email.com",
			password:          "",
			expectedUserLogin: nil,
			expectedError:     errors.New("MISSING_PASSWORD"),
		},
		{
			email:             "UNKNOWN_EMAIL",
			password:          "test",
			expectedUserLogin: nil,
			expectedError:     errors.New("EMAIL_NOT_FOUND"),
		},
		{
			email:             "test@email.com",
			password:          "WRONG_PASSWORD",
			expectedUserLogin: nil,
			expectedError:     errors.New("INVALID_PASSWORD"),
		},
		{
			email:    "test@email.com",
			password: "test",
			expectedUserLogin: &entity.UserLoginResponse{
				Email:        "test@email.com",
				DisplayName:  "Test User Name",
				IdToken:      "ValidIdToken",
				RefreshToken: "ValidRefreshToken",
				Expiration:   "3600",
				Error:        nil,
			},
			expectedError: nil,
		},
	}

	mockedExtApi := NewExternalApi()
	userApp := NewApplication(nil, mockedExtApi)

	for _, testCase := range tests {
		userLoginResponse, err := userApp.UserLogin("", testCase.email,
			testCase.password)

		assert.Equal(t, testCase.expectedError, err)
		assert.Equal(t, testCase.expectedUserLogin, userLoginResponse)
	}
}

func TestUserRefreshIdToken(t *testing.T) {
	type test struct {
		refreshToken            string
		expectedUserRefreshInfo *entity.UserRefreshTokenResponse
		expectedError           error
	}

	tests := []test{
		{
			refreshToken:            "",
			expectedUserRefreshInfo: nil,
			expectedError:           errors.New("MISSING_REFRESH_TOKEN"),
		},
		{
			refreshToken:            "WRONG_REFRESH_TOKEN",
			expectedUserRefreshInfo: nil,
			expectedError:           errors.New("INVALID_REFRESH_TOKEN"),
		},
		{
			refreshToken: "ValidRefreshToken",
			expectedUserRefreshInfo: &entity.UserRefreshTokenResponse{
				RefreshToken: "ValidRefreshToken",
				IdToken:      "ValidIdToken",
				UserUid:      "TestUserUID",
				TokenType:    "Bearer",
				Expiration:   "3600",
				Error:        nil,
			},
			expectedError: nil,
		},
	}

	mockedExtApi := NewExternalApi()
	userApp := NewApplication(nil, mockedExtApi)

	for _, testCase := range tests {
		refreshTokenInfo, err := userApp.UserRefreshIdToken("",
			testCase.refreshToken)

		assert.Equal(t, testCase.expectedError, err)
		assert.Equal(t, testCase.expectedUserRefreshInfo, refreshTokenInfo)
	}

}

func TestUserLoginOAuth2(t *testing.T) {
	type test struct {
		oauth2Idtoken          string
		providerId             string
		requestUri             string
		expectedUserInfoOAuth2 *entity.UserInfoOAuth2
		expectedError          error
	}

	tests := []test{
		{
			oauth2Idtoken:          "ERROR_IDP_RESPONSE",
			providerId:             "google.com",
			requestUri:             "http://test.com",
			expectedUserInfoOAuth2: nil,
			expectedError:          errors.New("INVALID_IDP_RESPONSE"),
		},
		{
			oauth2Idtoken: "NEW_USER",
			providerId:    "google.com",
			requestUri:    "http://test.com",
			expectedUserInfoOAuth2: &entity.UserInfoOAuth2{
				IdToken:       "ValidIdTokenWithoutPrivilegedUser",
				OAuthIdToken:  "NEW_USER",
				Email:         "test@email.com",
				EmailVerified: true,
				Fullname:      "Test Name",
				UserUid:       "TestUID",
				RefreshToken:  "ValidRefreshToken",
				Expiration:    "3600",
				IsNewUser:     true,
				Error:         nil,
			},
			expectedError: nil,
		},
		{
			oauth2Idtoken: "NO_NEW_USER",
			providerId:    "google.com",
			requestUri:    "http://test.com",
			expectedUserInfoOAuth2: &entity.UserInfoOAuth2{
				IdToken:       "ValidIdTokenWithoutPrivilegedUser",
				OAuthIdToken:  "NO_NEW_USER",
				Email:         "test@email.com",
				EmailVerified: true,
				Fullname:      "Test Name",
				UserUid:       "TestUID",
				RefreshToken:  "ValidRefreshToken",
				Expiration:    "3600",
				IsNewUser:     false,
				Error:         nil,
			},
			expectedError: nil,
		},
	}

	mockedExtApi := NewExternalApi()
	userApp := NewApplication(nil, mockedExtApi)

	for _, testCase := range tests {
		userInfoOAuth2, err := userApp.UserLoginOAuth2("",
			testCase.oauth2Idtoken, testCase.providerId, testCase.requestUri)

		assert.Equal(t, testCase.expectedError, err)
		assert.Equal(t, testCase.expectedUserInfoOAuth2, userInfoOAuth2)
	}

}
