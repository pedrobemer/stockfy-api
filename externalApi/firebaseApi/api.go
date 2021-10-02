package firebaseApi

import (
	"bytes"
	"context"
	"encoding/json"
	"stockfyApi/client"
	"stockfyApi/entity"

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
		return nil, err
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
	client.RequestAndAssignToBody("POST", url, bodyReader, &responseReqIdToken)

	return responseReqIdToken
}

// Send a Email for verification for a user with the specified ID token. This
// implementation uses the REST API from the Firebase project.
func (authClient *Firebase) SendVerificationEmail(webKey string,
	userIdToken string) entity.EmailVerificationResponse {

	var emailResponse entity.EmailVerificationResponse

	url := "https://identitytoolkit.googleapis.com/v1/accounts:sendOobCode?key=" +
		webKey

	bodyByte, _ := json.Marshal(EmailVerificationParams{
		RequestType: "VERIFY_EMAIL", IdToken: userIdToken})
	bodyReader := bytes.NewReader(bodyByte)
	client.RequestAndAssignToBody("POST", url, bodyReader, &emailResponse)

	emailResponse.UserIdToken = userIdToken

	return emailResponse
}

// Send a email to update the password. This implementation uses the REST API
// from the Firebase project.
func (authClient *Firebase) SendForgotPasswordEmail(webKey string,
	email string) entity.EmailForgotPasswordResponse {

	var emailPassResetResponse entity.EmailForgotPasswordResponse

	url := "https://identitytoolkit.googleapis.com/v1/accounts:sendOobCode?key=" +
		webKey

	bodyByte, _ := json.Marshal(PasswordReset{RequestType: "PASSWORD_RESET",
		Email: email})
	bodyReader := bytes.NewReader(bodyByte)
	client.RequestAndAssignToBody("POST", url, bodyReader, &emailPassResetResponse)

	return emailPassResetResponse
}
