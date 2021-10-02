package fiberHandlers

import (
	"fmt"
	"reflect"
	"stockfyApi/api/presenter"
	"stockfyApi/usecases"

	"github.com/gofiber/fiber/v2"
)

type FirebaseApi struct {
	ApplicationLogic usecases.Applications
	FirebaseWebKey   string
}

type emailVer struct {
	RequestType string `json:"requestType,omitempty"`
	IdToken     string `json:"idToken,omitempty"`
	Email       string `json:"email,omitempty"`
}

type passwordReset struct {
	RequestType string `json:"requestType,omitempty"`
	Email       string `json:"email,omitempty"`
}

func (firebaseAuth *FirebaseApi) SignUp(c *fiber.Ctx) error {
	var err error
	var signUpUser presenter.SignUpBody

	if err := c.BodyParser(&signUpUser); err != nil {
		fmt.Println(err)
	}

	// Create the user on Firebase
	user, err := firebaseAuth.ApplicationLogic.UserApp.UserCreate(signUpUser.Email,
		signUpUser.Password, signUpUser.DisplayName)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	// Create Custom token for the user with a specific UID
	token, err := firebaseAuth.ApplicationLogic.UserApp.UserCreateCustomToken(
		user.UID)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	// Request a ID token for Firebase BASED on the custom token
	userIdToken, err := firebaseAuth.ApplicationLogic.UserApp.UserRequestIdToken(
		firebaseAuth.FirebaseWebKey, token)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	// Sent Email verification for every new user created
	emailVerificationResp, err := firebaseAuth.ApplicationLogic.UserApp.
		UserSendVerificationEmail(firebaseAuth.FirebaseWebKey, userIdToken.IdToken)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
			"error":   emailVerificationResp.Error,
		})
	}

	// Create User in our database
	_, err = firebaseAuth.ApplicationLogic.UserApp.CreateUser(
		user.UID, user.Email, user.DisplayName, "normal")
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	userApiReturn := presenter.ConvertUserToUserApiReturn(user.Email,
		user.DisplayName)

	if err := c.JSON(&fiber.Map{
		"success":  true,
		"userInfo": userApiReturn,
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
	var passwordResetEmail presenter.ForgotPasswordBody

	if err := c.BodyParser(&passwordResetEmail); err != nil {
		fmt.Println(err)
	}
	fmt.Println(passwordResetEmail)

	// Send Email to reset password
	emailForgotPassResp, err := firebaseAuth.ApplicationLogic.UserApp.
		UserSendForgotPasswordEmail(firebaseAuth.FirebaseWebKey,
			passwordResetEmail.Email)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
			"error":   emailForgotPassResp.Error,
		})
	}

	if err := c.JSON(&fiber.Map{
		"success":        true,
		"passwordUpdate": emailForgotPassResp,
		"message":        "The email for password reset was successfully sent",
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

	// Delete User from Firebase
	deletedUser, err := firebaseAuth.ApplicationLogic.UserApp.UserDelete(
		userId.String())
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	firebaseAuth.ApplicationLogic.UserApp.DeleteUser(userId.String())

	if err := c.JSON(&fiber.Map{
		"success":  true,
		"userInfo": deletedUser,
		"message":  "The user was deleted successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err
}

// func (firebaseAuth *FirebaseApi) UpdateUserInfo(c *fiber.Ctx) error {
// 	var err error
// 	var userInfoUpdate database.SignUpBodyPost
// 	var userDb database.UserDatabase

// 	userInfo := c.Context().Value("user")
// 	userId := reflect.ValueOf(userInfo).FieldByName("userID")

// 	if err := c.BodyParser(&userInfoUpdate); err != nil {
// 		fmt.Println(err)
// 	}

// 	params := (&auth.UserToUpdate{})

// 	if userInfoUpdate.DisplayName != "" {
// 		params.DisplayName(userInfoUpdate.DisplayName)
// 	}
// 	if userInfoUpdate.Email != "" {
// 		params.DisplayName(userInfoUpdate.Email)
// 	}
// 	if userInfoUpdate.Password != "" {
// 		params.Password(userInfoUpdate.Password)
// 	}
// 	fmt.Println(params)

// 	userRecord, err := firebaseAuth.FirebaseAuth.UpdateUser(context.Background(),
// 		userId.String(), params)
// 	if err != nil {
// 		return c.Status(404).JSON(&fiber.Map{
// 			"success": false,
// 			"message": err.Error(),
// 		})
// 	}

// 	userDb.Email = userRecord.Email
// 	userDb.Uid = userRecord.UID
// 	userDb.Username = userRecord.DisplayName
// 	database.UpdateUser(firebaseAuth.Db, userDb)

// 	if err := c.JSON(&fiber.Map{
// 		"success":  true,
// 		"userInfo": userRecord,
// 		"message":  "The user information was updated successfully",
// 	}); err != nil {
// 		return c.Status(500).JSON(&fiber.Map{
// 			"success": false,
// 			"message": err,
// 		})
// 	}

// 	return err
// }
