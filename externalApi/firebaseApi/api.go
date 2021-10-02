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

func (authClient *Firebase) DeleteUser(userId string) (*auth.UserRecord, error) {
	var err error

	userInfo, err := authClient.Auth.GetUser(context.Background(), userId)
	err = authClient.Auth.DeleteUser(context.Background(), userId)
	if err != nil {
		return nil, err
	}

	return userInfo, err
}

func (authClient *Firebase) CustomToken(userUid string) (string, error) {
	return authClient.Auth.CustomToken(context.Background(), userUid)
}

func (authClient *Firebase) RequestIdToken(webKey string, customToken string) entity.ReqIdToken {

	var responseReqIdToken entity.ReqIdToken

	url := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key=" +
		webKey
	bodyByte, _ := json.Marshal(entity.ReqIdToken{Token: customToken,
		RequestSecureToken: true})
	bodyReader := bytes.NewReader(bodyByte)
	client.RequestAndAssignToBody("POST", url, bodyReader, &responseReqIdToken)

	return responseReqIdToken
}

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
