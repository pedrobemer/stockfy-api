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

func (a *Application) DeleteUser(userUid string) (*entity.Users, error) {
	user, err := a.repo.Delete(userUid)
	if err != nil {
		return nil, err
	}

	return &user[0], err
}

func (a *Application) UpdateUser(userUid string, email string,
	displayName string) (*entity.Users, error) {
	updateInfo, err := entity.NewUser(userUid, displayName, email, "normal")
	if err != nil {
		return nil, err
	}

	updateUser, err := a.repo.Update(*updateInfo)
	if err != nil {
		return nil, err
	}

	return &updateUser[0], nil
}

func (a *Application) SearchUser(userUid string) (*entity.Users, error) {
	searchedUser, err := a.repo.Search(userUid)
	if err != nil {
		return nil, err
	}

	return &searchedUser[0], nil
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

	apiResponse, err := a.extRepo.SendVerificationEmail(webKey, userIdToken)
	if err != nil {
		return apiResponse, entity.ErrInvalidUserSendEmail
	}

	return apiResponse, nil
}

// Send a email to recover the forgot password based on the given email
func (a *Application) UserSendForgotPasswordEmail(webKey string, email string) (
	entity.EmailForgotPasswordResponse, error) {
	apiResponse, err := a.extRepo.SendForgotPasswordEmail(webKey, email)
	if err != nil {
		return apiResponse, entity.ErrInvalidUserSendEmail
	}

	return apiResponse, nil
}

// Delete the user based on its UID
func (a *Application) UserDelete(userUid string) (*entity.UserInfo, error) {
	deletedUser, err := a.extRepo.DeleteUser(userUid)
	if err != nil {
		return nil, err
	}

	return deletedUser, nil
}

// Update the user information such as email, password and displayName based on
// its UID
func (a *Application) UserUpdateInfo(userUid string, email string,
	password string, displayName string) (*entity.UserInfo, error) {
	updateUserInfo, err := a.extRepo.UpdateUserInfo(userUid, email, password,
		displayName)
	if err != nil {
		return nil, err
	}

	return &updateUserInfo, nil
}

func (a *Application) UserTokenVerification(idToken string) (*entity.UserTokenInfo,
	error) {
	userTokenInfo, err := a.extRepo.VerifyIDToken(idToken)
	if err != nil {
		return nil, err
	}

	return &userTokenInfo, nil
}

func (a *Application) UserLogin(webKey string, email string, password string) (
	*entity.UserLoginResponse, error) {

	userLoginResponse, err := a.extRepo.UserLogin(webKey, email, password)

	if err != nil {
		return nil, err
	}

	return &userLoginResponse, nil
}

func (a *Application) UserRefreshIdToken(webKey string, refreshToken string) (
	*entity.UserRefreshTokenResponse, error) {

	refreshTokenResp, err := a.extRepo.UserRefreshIdToken(webKey, refreshToken)

	if err != nil {
		return nil, err
	}

	return &refreshTokenResp, nil
}

func (a *Application) UserLoginOAuth2(webKey string, oauthIdToken string,
	providerId string, requestUri string) (*entity.UserInfoOAuth2, error) {

	userInfoOAuth2, err := a.extRepo.UserLoginOAuth2(webKey, oauthIdToken,
		providerId, requestUri)

	if err != nil {
		return nil, err
	}

	return &userInfoOAuth2, nil
}
