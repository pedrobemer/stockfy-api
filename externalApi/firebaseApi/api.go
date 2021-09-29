package firebaseApi

import (
	"context"
	"fmt"
	"stockfyApi/database"

	"firebase.google.com/go/auth"
)

type Firebase struct {
	Auth *auth.Client
}

func (authClient *Firebase) CreateUser(
	signUpUser database.SignUpBodyPost) (*auth.UserRecord, error) {

	params := (&auth.UserToCreate{}).Email(signUpUser.Email).
		EmailVerified(false).Password(signUpUser.Password).
		DisplayName(signUpUser.DisplayName)

	user, err := authClient.Auth.CreateUser(context.Background(), params)
	fmt.Println(user)
	fmt.Println(err)

	return user, err
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
