package user

import (
	"stockfyApi/entity"
)

type Application struct {
	repo    Repository
	extRepo ExternalUserDatabase
}

//NewApplication create new use case
func NewApplication(r Repository, externalRepo ExternalUserDatabase) *Application {
	return &Application{
		repo:    r,
		extRepo: externalRepo,
	}
}

// Create User in our Repository (database)
func (a *Application) CreateUser(uid string, email string, displayName string,
	userType string) (*[]entity.Users, error) {
	userInfo, err := entity.NewUser(uid, displayName, email, userType)
	if err != nil {
		return nil, err
	}

	userCreated, err := a.repo.Create(*userInfo)

	return &userCreated, err
}

// Create User in Firebase
func (a *Application) UserCreate(email string, password string,
	displayName string) (*entity.UserInfo, error) {

	userInfo, err := a.extRepo.CreateUser(email, password, displayName)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

// Create a custom token for a user with a given UID
func (a *Application) UserCreateCustomToken(userUid string) (string, error) {
	return a.extRepo.CustomToken(userUid)
}

// User with a customToken request for the firebase API its id Token for
// authentication
func (a *Application) UserRequestIdToken(webKey string, customToken string) (
	*entity.ReqIdToken, error) {
	userTokenInfo := a.extRepo.RequestIdToken(webKey, customToken)
	if userTokenInfo.IdToken == "" {
		return nil, entity.ErrInvalidUserToken
	}

	return &userTokenInfo, nil
}

// Send a verification email for the user with correspondent id Token
func (a *Application) UserSendVerificationEmail(webKey, userIdToken string) (
	entity.EmailVerificationResponse, error) {
	apiResponse := a.extRepo.SendVerificationEmail(webKey, userIdToken)
	if apiResponse.Error != nil {
		return apiResponse, entity.ErrInvalidUserEmailVerification
	}

	return apiResponse, nil
}
