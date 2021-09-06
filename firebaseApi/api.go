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
