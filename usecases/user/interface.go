package user

import (
	"stockfyApi/entity"

	"firebase.google.com/go/auth"
)

type Repository interface {
	Create(signUp entity.Users) ([]entity.Users, error)
}

type ExternalUserDatabase interface {
	CreateUser(email string, password string, displayName string) (
		*entity.UserInfo, error)
	DeleteUser(userId string) (*auth.UserRecord, error)
	CustomToken(userUid string) (string, error)
	RequestIdToken(webKey string, customToken string) entity.ReqIdToken
	SendVerificationEmail(webKey string, userIdToken string) entity.
		EmailVerificationResponse
}
