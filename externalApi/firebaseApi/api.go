package firebaseApi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"stockfyApi/client"
	"stockfyApi/entity"
	"strings"

	"firebase.google.com/go/auth"
)

type Firebase struct {
	Auth *auth.Client
}

func NewFirebase(auth *auth.Client) *Firebase {
	return &Firebase{
		Auth: auth,
	}
}

// Create user in the Firebase authentication database using the Firebase SDK
func (authClient *Firebase) CreateUser(email string, password string,
	displayName string) (*entity.UserInfo, error) {

	params := (&auth.UserToCreate{}).Email(email).
		EmailVerified(false).Password(password).
		DisplayName(displayName)

	user, err := authClient.Auth.CreateUser(context.Background(), params)
	if err != nil {
		var errorMsg string
		var splittedError []string
		if strings.Contains(err.Error(), "message") {
			splittedError = strings.Fields(err.Error())
			errorMsg = strings.ReplaceAll(splittedError[11], "\"", "")

		}
		errorMsg = err.Error()

		return nil, errors.New(errorMsg)
	}
	userInfo := entity.ConvertUserInfo(user.Email, user.DisplayName, user.UID)

	return &userInfo, err
}

// Delete user in the Firebase authentication database using the Firebase SDK
func (authClient *Firebase) DeleteUser(userUid string) (*entity.UserInfo, error) {
	var err error

	userInfo, _ := authClient.Auth.GetUser(context.Background(), userUid)
	err = authClient.Auth.DeleteUser(context.Background(), userUid)
	if err != nil {
		return nil, err
	}

	deletedUserInfo := entity.ConvertUserInfo(userInfo.Email,
		userInfo.DisplayName, userUid)

	return &deletedUserInfo, err
}

// Create a custom token for a given user based on its correspodent UID. This
// implementation uses the Firebase SDK
func (authClient *Firebase) CustomToken(userUid string) (string, error) {
	return authClient.Auth.CustomToken(context.Background(), userUid)
}

// Request a ID token based on a custom token. This implementation uses the REST
// API from the Firebase project.
func (authClient *Firebase) RequestIdToken(webKey string, customToken string) entity.
	ReqIdToken {

	var responseReqIdToken entity.ReqIdToken

	url := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key=" +
		webKey
	bodyByte, _ := json.Marshal(entity.ReqIdToken{Token: customToken,
		RequestSecureToken: true})
	bodyReader := bytes.NewReader(bodyByte)
	client.RequestAndAssignToBody("POST", url, "application/json", bodyReader,
		&responseReqIdToken)

	return responseReqIdToken
}

// Send a Email for verification for a user with the specified ID token. This
// implementation uses the REST API from the Firebase project.
func (authClient *Firebase) SendVerificationEmail(webKey string,
	userIdToken string) (entity.EmailVerificationResponse, error) {

	var emailResponse entity.EmailVerificationResponse

	url := "https://identitytoolkit.googleapis.com/v1/accounts:sendOobCode?key=" +
		webKey

	bodyByte, _ := json.Marshal(EmailVerificationParams{
		RequestType: "VERIFY_EMAIL", IdToken: userIdToken})
	bodyReader := bytes.NewReader(bodyByte)
	client.RequestAndAssignToBody("POST", url, "application/json", bodyReader,
		&emailResponse)
	if emailResponse.Error != nil {
		errorMap := emailResponse.Error["errors"]
		errorString := entity.InterfaceToString(errorMap)
		splittedError := strings.Fields(errorString)
		errorMsg := strings.ReplaceAll(splittedError[1], "message:", "")

		return emailResponse, errors.New(errorMsg)

	}

	emailResponse.UserIdToken = userIdToken

	return emailResponse, nil
}

// Send a email to update the password. This implementation uses the REST API
// from the Firebase project.
func (authClient *Firebase) SendForgotPasswordEmail(webKey string,
	email string) (entity.EmailForgotPasswordResponse, error) {

	var emailPassResetResponse entity.EmailForgotPasswordResponse

	url := "https://identitytoolkit.googleapis.com/v1/accounts:sendOobCode?key=" +
		webKey

	bodyByte, _ := json.Marshal(PasswordReset{RequestType: "PASSWORD_RESET",
		Email: email})
	bodyReader := bytes.NewReader(bodyByte)
	client.RequestAndAssignToBody("POST", url, "application/json", bodyReader,
		&emailPassResetResponse)
	if emailPassResetResponse.Error != nil {
		errorMap := emailPassResetResponse.Error["errors"]
		errorString := entity.InterfaceToString(errorMap)
		splittedError := strings.Fields(errorString)
		errorMsg := strings.ReplaceAll(splittedError[1], "message:", "")

		return emailPassResetResponse, errors.New(errorMsg)
	}

	return emailPassResetResponse, nil
}

func (authClient *Firebase) UpdateUserInfo(usedUid string, email string,
	password string, displayName string) (entity.UserInfo, error) {

	params := (&auth.UserToUpdate{})

	if displayName != "" {
		params.DisplayName(displayName)
	}
	if email != "" {
		params.DisplayName(email)
	}
	if password != "" {
		params.Password(password)
	}

	userUpdateInfo, err := authClient.Auth.UpdateUser(context.Background(),
		usedUid, params)

	useInfo := entity.ConvertUserInfo(userUpdateInfo.Email,
		userUpdateInfo.DisplayName, userUpdateInfo.UID)

	return useInfo, err
}

func (authClient *Firebase) VerifyIDToken(idToken string) (entity.UserTokenInfo,
	error) {
	firebaseToken, err := authClient.Auth.VerifyIDToken(context.Background(),
		idToken)
	if err != nil {
		return entity.UserTokenInfo{}, err
	}

	userTokenInfo := entity.ConvertUserTokenInfo(firebaseToken.Claims["user_id"].(string),
		firebaseToken.Claims["email"].(string), firebaseToken.Claims["email_verified"].(bool))

	return userTokenInfo, nil
}

func (authClient *Firebase) UserLogin(webKey string, email string,
	password string) (entity.UserLoginResponse, error) {
	var loginResponse entity.UserLoginResponse

	url := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?" +
		"key=" + webKey

	bodyByte, _ := json.Marshal(UserLogin{
		Email: email, Password: password, ReturnSecureToken: true})
	bodyReader := bytes.NewReader(bodyByte)
	client.RequestAndAssignToBody("POST", url, "application/json", bodyReader,
		&loginResponse)

	if loginResponse.Error != nil {
		errorMap := loginResponse.Error["errors"]
		errorString := entity.InterfaceToString(errorMap)
		splittedError := strings.Fields(errorString)
		errorMsg := strings.ReplaceAll(splittedError[1], "message:", "")

		return loginResponse, errors.New(errorMsg)
	}

	return loginResponse, nil
}

func (authClient *Firebase) UserRefreshIdToken(webKey string,
	refreshToken string) (entity.UserRefreshTokenResponse, error) {

	var refreshTokenResponse entity.UserRefreshTokenResponse

	urlReq := "https://securetoken.googleapis.com/v1/token?key=" + webKey

	dataUrlFormMap := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	}

	dataUrlFormStr := dataUrlFormMap.Encode()
	client.RequestAndAssignToBody("POST", urlReq, "application/x-www-form-urlencoded",
		strings.NewReader(dataUrlFormStr), &refreshTokenResponse)

	if refreshTokenResponse.Error != nil {
		errorInterface := refreshTokenResponse.Error["message"]
		errorString := entity.InterfaceToString(errorInterface)

		return refreshTokenResponse, errors.New(errorString)
	}

	return refreshTokenResponse, nil
}

func (authClient *Firebase) UserLoginOAuth2(webKey string, idToken string,
	providerId string, requestUri string) (entity.UserInfoOAuth2, error) {
	var oauthUserInfo entity.UserInfoOAuth2
	var postBody string

	url := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithIdp?key=" +
		webKey

	if providerId == "google.com" {
		postBody = "id_token=" + idToken + "&providerId=" + providerId
	} else {
		postBody = "access_token=" + idToken + "&providerId=" + providerId
	}

	bodyByte, _ := json.Marshal(UserLoginOAuth2{
		PostBody:            postBody,
		RequestUri:          requestUri,
		ReturnIdpCredential: true,
		ReturnSecureToken:   true,
	})
	bodyReader := bytes.NewReader(bodyByte)

	client.RequestAndAssignToBody("POST", url, "application/json", bodyReader,
		&oauthUserInfo)

	if oauthUserInfo.Error != nil {
		errorInterface := oauthUserInfo.Error["message"]
		errorString := entity.InterfaceToString(errorInterface)

		return oauthUserInfo, errors.New(errorString)
	}

	return oauthUserInfo, nil

}
