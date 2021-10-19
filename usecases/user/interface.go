package user

import (
	"stockfyApi/entity"
)

type Repository interface {
	Create(signUp entity.Users) ([]entity.Users, error)
	Delete(firebaseUid string) ([]entity.Users, error)
	Update(userInfo entity.Users) ([]entity.Users, error)
	Search(userUid string) ([]entity.Users, error)
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
	UpdateUserInfo(usedUid string, email string, password string,
		displayName string) (entity.UserInfo, error)
	VerifyIDToken(idToken string) (entity.UserTokenInfo, error)
}

type UseCases interface {
	CreateUser(uid string, email string, displayName string, userType string) (
		*[]entity.Users, error)
	DeleteUser(userUid string) (*entity.Users, error)
	UpdateUser(userUid string, email string, displayName string) (*entity.Users,
		error)
	SearchUser(userUid string) (*entity.Users, error)
	UserCreate(email string, password string, displayName string) (
		*entity.UserInfo, error)
	UserCreateCustomToken(userUid string) (string, error)
	UserRequestIdToken(webKey string, customToken string) (*entity.ReqIdToken,
		error)
	UserSendVerificationEmail(webKey, userIdToken string) (
		entity.EmailVerificationResponse, error)
	UserSendForgotPasswordEmail(webKey string, email string) (
		entity.EmailForgotPasswordResponse, error)
	UserDelete(userUid string) (*entity.UserInfo, error)
	UserUpdateInfo(userUid string, email string, password string,
		displayName string) (*entity.UserInfo, error)
	UserTokenVerification(idToken string) (*entity.UserTokenInfo, error)
}
