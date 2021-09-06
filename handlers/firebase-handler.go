package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"stockfyApi/client"
	"stockfyApi/database"
	"stockfyApi/firebaseApi"

	"firebase.google.com/go/auth"
	"github.com/gofiber/fiber/v2"
)

type FirebaseApi struct {
	Db             database.PgxIface
	FirebaseAuth   *auth.Client
	FirebaseWebKey string
}

type emailVer struct {
	RequestType string `json:"requestType,omitempty"`
	IdToken     string `json:"idToken,omitempty"`
	Email       string `json:"email,omitempty"`
}

type reqIdToken struct {
	Token              string `json:"token,omitempty"`
	RequestSecureToken bool   `json:"requestSecureToken,omitempty"`
	Kind               string `json:"kind,omitempty"`
	IdToken            string `json:"idToken,omitempty"`
	IsNewUser          bool   `json:"isNewUser,omitempty"`
}

func (firebaseAuth *FirebaseApi) SignUp(c *fiber.Ctx) error {
	var err error
	var signUpUser database.SignUpBodyPost
	var bodyRespEmail emailVer
	var bodyRespIdToken reqIdToken

	if err := c.BodyParser(&signUpUser); err != nil {
		fmt.Println(err)
	}
	fmt.Println(signUpUser)

	authApi := firebaseApi.Firebase{Auth: firebaseAuth.FirebaseAuth}

	// Create the user on Firebase
	user, err := authApi.CreateUser(signUpUser)
	fmt.Println(user)

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	// Create Custom Token for email verification
	token, err := authApi.Auth.CustomToken(context.Background(), user.UID)

	// Request a ID token for Firebase
	url := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key=" +
		firebaseAuth.FirebaseWebKey
	bodyByte, _ := json.Marshal(reqIdToken{Token: token,
		RequestSecureToken: true})
	bodyReader := bytes.NewReader(bodyByte)
	client.RequestAndAssignToBody("POST", url, bodyReader, &bodyRespIdToken)

	// Sent Email verification for every new user created
	url = "https://identitytoolkit.googleapis.com/v1/accounts:sendOobCode?key=" +
		firebaseAuth.FirebaseWebKey
	bodyByte, _ = json.Marshal(emailVer{RequestType: "VERIFY_EMAIL",
		IdToken: bodyRespIdToken.IdToken})
	bodyReader = bytes.NewReader(bodyByte)
	client.RequestAndAssignToBody("POST", url, bodyReader, &bodyRespEmail)

	if bodyRespEmail.Email != user.Email {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "Email does not match",
		})
	}

	if err := c.JSON(&fiber.Map{
		"success":  true,
		"userInfo": user,
		"message":  "User registered successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err
}
