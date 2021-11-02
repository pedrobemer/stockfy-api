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

	userCreated := []entity.Users{
		{
			Id:       "39148-38149v-jk48",
			Username: "Test Name",
			Email:    "test@gmail.com",
			Uid:      "93avpow384",
			Type:     "normal",
		},
	}

	return userCreated, nil
}

func (m *MockDb) Delete(firebaseUid string) ([]entity.Users, error) {

	deletedUser := []entity.Users{
		{
			Id:       "39148-38149v-jk48",
			Username: "Test Name",
			Email:    "test@gmail.com",
			Uid:      "93avpow384",
			Type:     "normal",
		},
	}
	return deletedUser, nil
}

func (m *MockDb) Update(userInfo entity.Users) ([]entity.Users, error) {
	return []entity.Users{
		{
			Id:       "391ahb4",
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
	if email == "Error" {
		return nil, errors.New("Error Mock Firebase")
	}

	return &entity.UserInfo{
		DisplayName: "Test Name",
		Email:       "test@gmail.com",
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
	return entity.UserTokenInfo{}, nil
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
