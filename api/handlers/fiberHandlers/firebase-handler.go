package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
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

type passwordReset struct {
	RequestType string `json:"requestType,omitempty"`
	Email       string `json:"email,omitempty"`
}

func (firebaseAuth *FirebaseApi) SignUp(c *fiber.Ctx) error {
	var err error
	var signUpUser database.SignUpBodyPost
	var bodyRespEmail emailVer
	var bodyRespIdToken reqIdToken
	var userDb database.UserDatabase

	if err := c.BodyParser(&signUpUser); err != nil {
		fmt.Println(err)
	}

	authApi := firebaseApi.Firebase{Auth: firebaseAuth.FirebaseAuth}

	// Create the user on Firebase
	user, err := authApi.CreateUser(signUpUser)
	fmt.Println(user)

	if err != nil {
		return c.Status(409).JSON(&fiber.Map{
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

	userDb.Email = signUpUser.Email
	userDb.Uid = user.UID
	userDb.Username = user.DisplayName

	_, err = database.CreateUser(firebaseAuth.Db, userDb)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "User not saved on the database",
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

func (firebaseAuth *FirebaseApi) ForgotPassword(c *fiber.Ctx) error {
	var err error
	var passwordResetEmail passwordReset
	var bodyRespPassReset passwordReset

	if err := c.BodyParser(&passwordResetEmail); err != nil {
		fmt.Println(err)
	}
	fmt.Println(passwordResetEmail)

	// Request a ID token for Firebase
	url := "https://identitytoolkit.googleapis.com/v1/accounts:sendOobCode?key=" +
		firebaseAuth.FirebaseWebKey
	bodyByte, _ := json.Marshal(passwordReset{RequestType: "PASSWORD_RESET",
		Email: passwordResetEmail.Email})
	bodyReader := bytes.NewReader(bodyByte)
	client.RequestAndAssignToBody("POST", url, bodyReader, &bodyRespPassReset)

	if err := c.JSON(&fiber.Map{
		"success":  true,
		"userInfo": bodyRespPassReset.Email,
		"message":  "The email for password reset was successfully sent",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err
}

func (firebaseAuth *FirebaseApi) DeleteUser(c *fiber.Ctx) error {
	var err error

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	authApi := firebaseApi.Firebase{Auth: firebaseAuth.FirebaseAuth}

	userRecord, err := authApi.DeleteUser(userId.String())
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	database.DeleteUser(firebaseAuth.Db, userId.String())

	if err := c.JSON(&fiber.Map{
		"success":  true,
		"userInfo": userRecord,
		"message":  "The user was deleted successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err
}

func (firebaseAuth *FirebaseApi) UpdateUserInfo(c *fiber.Ctx) error {
	var err error
	var userInfoUpdate database.SignUpBodyPost
	var userDb database.UserDatabase

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	if err := c.BodyParser(&userInfoUpdate); err != nil {
		fmt.Println(err)
	}

	params := (&auth.UserToUpdate{})

	if userInfoUpdate.DisplayName != "" {
		params.DisplayName(userInfoUpdate.DisplayName)
	}
	if userInfoUpdate.Email != "" {
		params.DisplayName(userInfoUpdate.Email)
	}
	if userInfoUpdate.Password != "" {
		params.Password(userInfoUpdate.Password)
	}
	fmt.Println(params)

	userRecord, err := firebaseAuth.FirebaseAuth.UpdateUser(context.Background(),
		userId.String(), params)
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	userDb.Email = userRecord.Email
	userDb.Uid = userRecord.UID
	userDb.Username = userRecord.DisplayName
	database.UpdateUser(firebaseAuth.Db, userDb)

	if err := c.JSON(&fiber.Map{
		"success":  true,
		"userInfo": userRecord,
		"message":  "The user information was updated successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err
}
