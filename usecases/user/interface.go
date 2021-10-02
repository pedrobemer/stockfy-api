package user

import (
	"stockfyApi/entity"
)

type Repository interface {
	Create(signUp entity.Users) ([]entity.Users, error)
	Delete(firebaseUid string) ([]entity.Users, error)
}

type ExternalUserDatabase interface {
	CreateUser(email string, password string, displayName string) (
		*entity.UserInfo, error)
	DeleteUser(userId string) (*entity.UserInfo, error)
	CustomToken(userUid string) (string, error)
	RequestIdToken(webKey string, customToken string) entity.ReqIdToken
	SendVerificationEmail(webKey string, userIdToken string) entity.
		EmailVerificationResponse
	SendForgotPasswordEmail(webKey string, email string) entity.
		EmailForgotPasswordResponse
}
