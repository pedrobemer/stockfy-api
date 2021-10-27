package user

import (
	"errors"
	"stockfyApi/entity"
)

type MockApplication struct {
}

func NewMockApplication() *MockApplication {
	return &MockApplication{}
}

func (a *MockApplication) CreateUser(uid string, email string, displayName string,
	userType string) (*[]entity.Users, error) {

	userInfo, err := entity.NewUser(uid, displayName, email, userType)
	if err != nil {
		return nil, err
	}

	return &[]entity.Users{*userInfo}, nil
}

func (a *MockApplication) DeleteUser(userUid string) (*entity.Users, error) {

	return &entity.Users{
		Id:       "TestId",
		Uid:      "TestUID",
		Username: "Test Name",
		Email:    "test@email.com",
		Type:     "normal",
	}, nil

}

func (a *MockApplication) UpdateUser(userUid string, email string,
	displayName string) (*entity.Users, error) {
	updateInfo, err := entity.NewUser(userUid, displayName, email, "normal")
	if err != nil {
		return nil, err
	}

	return &entity.Users{
		Id:       "TestId",
		Uid:      updateInfo.Uid,
		Username: updateInfo.Username,
		Email:    updateInfo.Email,
		Type:     updateInfo.Type,
	}, nil
}

func (a *MockApplication) SearchUser(userUid string) (*entity.Users, error) {
	if userUid == "INVALID_UID" {
		return nil, entity.ErrInvalidUserSearch
	} else if userUid == "USER_WITH_PRIVILEGE" {
		return &entity.Users{
			Id:       "TestId",
			Uid:      userUid,
			Username: "Test Name",
			Email:    "test@email.com",
			Type:     "admin",
		}, nil
	} else {
		return &entity.Users{
			Id:       "TestId",
			Uid:      userUid,
			Username: "Test Name",
			Email:    "test@email.com",
			Type:     "normal",
		}, nil
	}
}

func (a *MockApplication) UserCreate(email string, password string,
	displayName string) (*entity.UserInfo, error) {

	if email == "" {
		return nil, errors.New("email must be a non-empty string")
	} else if password == "" {
		return nil, errors.New("password must be a string at least 6 characters long")
	} else if displayName == "" {
		return nil, errors.New("display name must be a non-empty string")
	} else if displayName == "WRONG_CUSTOM_TOKEN" {
		return &entity.UserInfo{
			UID:         "INVALID_USER_UID",
			Email:       email,
			DisplayName: displayName,
		}, nil
	} else if displayName == "WRONG_ID_TOKEN" {
		return &entity.UserInfo{
			UID:         "INVALID_ID_TOKEN",
			Email:       email,
			DisplayName: displayName,
		}, nil
	} else if displayName == "WRONG_EMAIL_VERIFICATION" {
		return &entity.UserInfo{
			UID:         "INVALID_SEND_EMAIL",
			Email:       email,
			DisplayName: displayName,
		}, nil
	} else if displayName == "WRONG_USER_INFO" {
		return &entity.UserInfo{
			UID:         "INVALID_USER_INFO",
			Email:       "",
			DisplayName: displayName,
		}, nil
	} else {
		return &entity.UserInfo{
			UID:         "TestUID",
			Email:       email,
			DisplayName: displayName,
		}, nil
	}
}

func (a *MockApplication) UserCreateCustomToken(userUid string) (string, error) {
	if userUid == "INVALID_USER_UID" {
		return "", errors.New("Some Error")
	} else if userUid == "INVALID_ID_TOKEN" {
		return "INVALID_CUSTOM_TOKEN", nil
	} else if userUid == "INVALID_SEND_EMAIL" {
		return "INVALID_SEND_EMAIL", nil
	} else {
		return "validCustomToken", nil

	}
}

func (a *MockApplication) UserRequestIdToken(webKey string, customToken string) (
	*entity.ReqIdToken, error) {
	if customToken == "INVALID_CUSTOM_TOKEN" {
		return nil, entity.ErrInvalidUserToken
	} else if customToken == "INVALID_SEND_EMAIL" {
		return &entity.ReqIdToken{
			Token:              "TestToken",
			RequestSecureToken: true,
			Kind:               "TestKind",
			IdToken:            "INVALID_USER_UID_TOKEN",
			IsNewUser:          false,
		}, nil
	} else {
		return &entity.ReqIdToken{
			Token:              "TestToken",
			RequestSecureToken: true,
			Kind:               "TestKind",
			IdToken:            "ValidBearerToken",
			IsNewUser:          false,
		}, nil
	}
}

func (a *MockApplication) UserSendVerificationEmail(webKey, userIdToken string) (
	entity.EmailVerificationResponse, error) {
	if userIdToken == "INVALID_USER_UID_TOKEN" {
		return entity.EmailVerificationResponse{
			Error: map[string]interface{}{
				"code":    400,
				"message": "INVALID_ID_TOKEN",
				"errors": map[string]string{
					"message": "INVALID_ID_TOKEN",
					"domain":  "global",
					"reason":  "invalid",
				},
			},
		}, errors.New("INVALID_ID_TOKEN")
	}

	return entity.EmailVerificationResponse{
		UserIdToken: userIdToken,
		Email:       "test@email.com",
		Error:       nil,
	}, nil
}

func (a *MockApplication) UserSendForgotPasswordEmail(webKey string, email string) (
	entity.EmailForgotPasswordResponse, error) {
	if email == "INVALID_EMAIL" {
		return entity.EmailForgotPasswordResponse{
			Error: map[string]interface{}{
				"code":    400,
				"message": "EMAIL_NOT_FOUND",
				"errors": map[string]string{
					"message": "EMAIL_NOT_FOUND",
					"domain":  "global",
					"reason":  "invalid",
				},
			},
		}, errors.New("EMAIL_NOT_FOUND")
	} else if email == "" {
		return entity.EmailForgotPasswordResponse{
			Error: map[string]interface{}{
				"code":    400,
				"message": "MISSING_EMAIL",
				"errors": map[string]string{
					"message": "MISSING_EMAIL",
					"domain":  "global",
					"reason":  "invalid",
				},
			},
		}, errors.New("MISSING_EMAIL")
	}

	return entity.EmailForgotPasswordResponse{
		Email: email,
		Error: nil,
	}, nil
}

func (a *MockApplication) UserDelete(userUid string) (*entity.UserInfo, error) {
	return &entity.UserInfo{
		UID:         "TestUID",
		Email:       "test@email.com",
		DisplayName: "Test Name",
	}, nil

}

func (a *MockApplication) UserUpdateInfo(userUid string, email string,
	password string, displayName string) (*entity.UserInfo, error) {
	return &entity.UserInfo{
		Email:       email,
		DisplayName: displayName,
		UID:         userUid,
	}, nil
}

func (a *MockApplication) UserTokenVerification(idToken string) (
	*entity.UserTokenInfo, error) {

	if idToken == "ValidIdTokenPrivilegeUser" {
		return &entity.UserTokenInfo{
			UserID:        "USER_WITH_PRIVILEGE",
			Email:         "test@email.com",
			EmailVerified: true,
		}, nil
	} else if idToken == "ValidIdTokenWithoutPrivilegedUser" {
		return &entity.UserTokenInfo{
			UserID:        "USER_WITHOUT_PRIVILEGE",
			Email:         "test@email.com",
			EmailVerified: true,
		}, nil
	} else if idToken == "ValidIdTokenWithoutEmailVerification" {
		return &entity.UserTokenInfo{
			UserID:        "Unverified User UID",
			Email:         "test@email.com",
			EmailVerified: false,
		}, nil
	} else {
		return nil, errors.New("Invalid Token")
	}
}
