package user

import (
	"errors"
	"stockfyApi/entity"
)

type MockDb struct {
}

type MockExternal struct {
}

func NewMockRepo() *MockDb {
	return &MockDb{}
}

func NewExternalApi() *MockExternal {
	return &MockExternal{}
}

func (m *MockDb) Create(signUp entity.Users) ([]entity.Users, error) {

	if signUp.Email == "ERROR_USER_REPOSITORY" {
		return nil, errors.New("Unknown error in the user repository")
	}

	userCreated := []entity.Users{
		{
			Username: signUp.Username,
			Email:    signUp.Email,
			Uid:      signUp.Uid,
			Type:     signUp.Type,
		},
	}

	return userCreated, nil
}

func (m *MockDb) Delete(firebaseUid string) ([]entity.Users, error) {

	if firebaseUid == "ERROR_USER_REPOSITORY" {
		return nil, errors.New("Unknown delete error in the user repository")
	}

	deletedUser := []entity.Users{
		{
			Username: "Test Name",
			Email:    "test@gmail.com",
			Uid:      firebaseUid,
			Type:     "normal",
		},
	}
	return deletedUser, nil
}

func (m *MockDb) Update(userInfo entity.Users) ([]entity.Users, error) {
	if userInfo.Username == "ERROR_USER_REPOSITORY" {
		return nil, errors.New("Unknown update error in the user repository")
	}

	return []entity.Users{
		{
			Uid:      userInfo.Uid,
			Username: userInfo.Username,
			Email:    userInfo.Email,
			Type:     userInfo.Type,
		},
	}, nil
}

func (m *MockDb) Search(userUid string) ([]entity.Users, error) {
	if userUid == "Invalid" {
		return nil, entity.ErrInvalidUserSearch
	}

	return []entity.Users{
		{
			Uid:      "TestID",
			Email:    "test@gmail.com",
			Username: "Test Name",
			Type:     "normal",
		},
	}, nil
}

func (m *MockExternal) CreateUser(email string, password string,
	displayName string) (*entity.UserInfo, error) {

	if displayName == "" {
		return nil, errors.New("display name must be a non-empty string")
	}

	if len(password) < 6 {
		return nil, errors.New("password must be a string at least 6 characters long")
	}

	if email == "Error" {
		return nil, errors.New("Error Mock Firebase")
	}

	if email == "" {
		return nil, errors.New("email must be a non-empty string")
	}

	return &entity.UserInfo{
		DisplayName: displayName,
		Email:       email,
		UID:         "abj39as$$",
	}, nil
}

func (m *MockExternal) DeleteUser(userId string) (*entity.UserInfo, error) {
	if userId == "Invalid" {
		return nil, errors.New("Database Interface error")
	}

	return &entity.UserInfo{
		DisplayName: "Test Name",
		Email:       "test@gmail.com",
		UID:         userId,
	}, nil
}

func (m *MockExternal) CustomToken(userUid string) (string, error) {
	return "194nc4850d", nil
}
func (m *MockExternal) RequestIdToken(webKey string, customToken string) entity.ReqIdToken {
	if customToken == "49292" {
		return entity.ReqIdToken{}
	}

	return entity.ReqIdToken{
		Token:              "a419148a",
		RequestSecureToken: true,
		Kind:               "aa93q8",
		IdToken:            "294akfsnf49",
		IsNewUser:          false,
	}
}

func (m *MockExternal) SendVerificationEmail(webKey string, userIdToken string) (
	entity.EmailVerificationResponse, error) {

	if userIdToken == "Invalid" {

		err := map[string]interface{}{
			"code": 400,
			"errors": struct {
				domain  string
				message string
				reason  string
			}{"global", "INVALID_ID_TOKEN", "invalid"},
			"message": "INVALID_ID_TOKEN",
		}
		return entity.EmailVerificationResponse{
			UserIdToken: userIdToken,
			Error:       err,
		}, errors.New("INVALID_ID_TOKEN")
	}

	return entity.EmailVerificationResponse{
		UserIdToken: userIdToken,
		Email:       "test@gmail.com",
		Error:       nil,
	}, nil
}

func (m *MockExternal) SendForgotPasswordEmail(webKey string, email string) (
	entity.EmailForgotPasswordResponse, error) {

	if email == "Invalid" {

		err := map[string]interface{}{
			"code": 400,
			"errors": struct {
				domain  string
				message string
				reason  string
			}{"global", "EMAIL_NOT_FOUND", "invalid"},
			"message": "EMAIL_NOT_FOUND",
		}
		return entity.EmailForgotPasswordResponse{
			Email: email,
			Error: err,
		}, errors.New("EMAIL_NOT_FOUND")
	}

	return entity.EmailForgotPasswordResponse{
		Email: email,
		Error: nil,
	}, nil
}

func (m *MockExternal) UpdateUserInfo(usedUid string, email string,
	password string, displayName string) (entity.UserInfo, error) {
	var emailParams, nameParams string

	if displayName == "ERROR_USER_FIREBASE" {
		return entity.UserInfo{}, errors.New("Unknown update error in the user repository")
	}

	if displayName != "" {
		nameParams = displayName
	} else {
		nameParams = "Test Name"
	}

	if email != "" {
		emailParams = email
	} else {
		emailParams = "test@gmail.com"
	}

	return entity.UserInfo{
		UID:         usedUid,
		DisplayName: nameParams,
		Email:       emailParams,
	}, nil
}

func (m *MockExternal) VerifyIDToken(idToken string) (entity.UserTokenInfo, error) {

	if idToken == "INVALID_ID_TOKEN" {
		return entity.UserTokenInfo{}, errors.New("INVALID_ID_TOKEN")
	}

	return entity.UserTokenInfo{
		Email:         "test@email.com",
		EmailVerified: true,
		UserID:        "TestUserID",
	}, nil
}

func (m *MockExternal) UserLogin(webKey string, email string,
	password string) (entity.UserLoginResponse, error) {

	if email == "" {
		return entity.UserLoginResponse{}, errors.New("INVALID_EMAIL")
	}

	if password == "" {
		return entity.UserLoginResponse{}, errors.New("MISSING_PASSWORD")
	}

	if email == "UNKNOWN_EMAIL" {
		return entity.UserLoginResponse{}, errors.New("EMAIL_NOT_FOUND")
	}

	if password == "WRONG_PASSWORD" {
		return entity.UserLoginResponse{}, errors.New("INVALID_PASSWORD")
	}

	return entity.UserLoginResponse{
		Email:        email,
		DisplayName:  "Test User Name",
		IdToken:      "ValidIdToken",
		RefreshToken: "ValidRefreshToken",
		Expiration:   "3600",
		Error:        nil,
	}, nil
}

func (m *MockExternal) UserRefreshIdToken(webKey string,
	refreshToken string) (entity.UserRefreshTokenResponse, error) {

	if refreshToken == "" {
		return entity.UserRefreshTokenResponse{}, errors.New("MISSING_REFRESH_TOKEN")
	}

	if refreshToken == "WRONG_REFRESH_TOKEN" {
		return entity.UserRefreshTokenResponse{}, errors.New("INVALID_REFRESH_TOKEN")
	}

	return entity.UserRefreshTokenResponse{
		RefreshToken: refreshToken,
		IdToken:      "ValidIdToken",
		UserUid:      "TestUserUID",
		TokenType:    "Bearer",
		Expiration:   "3600",
		Error:        nil,
	}, nil
}

func (m *MockExternal) UserLoginOAuth2(webKey string, idToken string,
	providerId string, requestUri string) (entity.UserInfoOAuth2, error) {

	var isNewUser bool

	switch idToken {
	case "ERROR_IDP_RESPONSE":
		return entity.UserInfoOAuth2{}, errors.New("INVALID_IDP_RESPONSE")
	case "NEW_USER":
		isNewUser = true
		break
	default:
		isNewUser = false
	}

	return entity.UserInfoOAuth2{
		IdToken:       "ValidIdTokenWithoutPrivilegedUser",
		OAuthIdToken:  idToken,
		Email:         "test@email.com",
		EmailVerified: true,
		Fullname:      "Test Name",
		UserUid:       "TestUID",
		RefreshToken:  "ValidRefreshToken",
		Expiration:    "3600",
		IsNewUser:     isNewUser,
		Error:         nil,
	}, nil

}
