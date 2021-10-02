package user

import (
	"errors"
	"stockfyApi/entity"

	"firebase.google.com/go/auth"
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

// func (m *MockExternal) DeleteUser(userId string) (*auth.UserRecord, error) {

// }

func (m *MockExternal) DeleteUser(userId string) (*auth.UserRecord, error) {
	return &auth.UserRecord{}, nil
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

func (m *MockExternal) SendVerificationEmail(webKey string, userIdToken string) entity.
	EmailVerificationResponse {

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
		}
	}

	return entity.EmailVerificationResponse{
		UserIdToken: userIdToken,
		Email:       "test@gmail.com",
		Error:       nil,
	}
}
