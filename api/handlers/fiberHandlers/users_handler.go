package fiberHandlers

import (
	"reflect"
	"stockfyApi/api/presenter"
	"stockfyApi/entity"
	"stockfyApi/externalApi/oauth2"
	"stockfyApi/token"
	"stockfyApi/usecases"
	"stockfyApi/usecases/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UsersApi struct {
	ApplicationLogic    usecases.Applications
	FirebaseWebKey      string
	GoogleOAuth2        oauth2.GoogleOAuth2
	FacebookOAuth2      oauth2.FacebookOAuth2
	StateTokenInterface token.Maker
	TokenMaker          func(symmetricKey string) (token.Maker, error)
	StateUsername       string
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

func (f *UsersApi) SignUp(c *fiber.Ctx) error {

	var signUpUser presenter.SignUpBody

	if err := c.BodyParser(&signUpUser); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   entity.ErrInvalidApiBody.Error(),
			"code":    400,
		})
	}

	// Create the user on Firebase
	user, err := f.ApplicationLogic.UserApp.UserCreate(signUpUser.Email,
		signUpUser.Password, signUpUser.DisplayName)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   err.Error(),
			"code":    400,
		})
	}

	// Create User in our database
	_, err = f.ApplicationLogic.UserApp.CreateUser(user.UID, user.Email,
		user.DisplayName, "normal")
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiInternalError.Error(),
			"error":   err.Error(),
			"code":    500,
		})
	}

	// Create Custom token for the user with a specific UID
	token, err := f.ApplicationLogic.UserApp.UserCreateCustomToken(user.UID)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiInternalError.Error(),
			"error":   err.Error(),
			"code":    500,
		})
	}

	// Request a ID token for Firebase based on the custom token
	userIdToken, err := f.ApplicationLogic.UserApp.UserRequestIdToken(
		f.FirebaseWebKey, token)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiInternalError.Error(),
			"error":   err.Error(),
			"code":    500,
		})
	}

	// Sent Email verification for every new user created
	_, err = f.ApplicationLogic.UserApp.UserSendVerificationEmail(
		f.FirebaseWebKey, userIdToken.IdToken)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   err.Error(),
			"code":    400,
		})
	}

	userApiReturn := presenter.ConvertUserToUserApiReturn(user.Email,
		user.DisplayName)

	err = c.JSON(&fiber.Map{
		"success":  true,
		"userInfo": userApiReturn,
		"message":  "User was registered successfully",
	})

	return err
}

func (f *UsersApi) SignIn(c *fiber.Ctx) error {
	var userLogin presenter.SignInBody

	if err := c.BodyParser(&userLogin); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   entity.ErrInvalidApiBody.Error(),
			"code":    400,
		})
	}

	userLoginResponse, err := f.ApplicationLogic.UserApp.UserLogin(
		f.FirebaseWebKey, userLogin.Email, userLogin.Password)

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   err.Error(),
			"code":    400,
		})
	}

	userLoginApiReturn := presenter.ConvertUserLoginToUserLoginApiReturn(
		userLoginResponse.Email, userLoginResponse.DisplayName,
		userLoginResponse.IdToken, userLoginResponse.RefreshToken,
		userLoginResponse.Expiration)

	err = c.JSON(&fiber.Map{
		"success":  true,
		"userInfo": userLoginApiReturn,
		"message":  "User login was successful",
	})

	return nil
}

// First phase from the OAuth2 authetication process. The goal here is for the
// user to login using third-party accounts and give consent to get its user
// profile information with the authorization code.
func (f *UsersApi) SignInOAuth(c *fiber.Ctx) error {

	// Specify which kind of token maker, we will use to create our state
	// variable, with a new random symmetric key every time.
	tokenMaker, _ := f.TokenMaker(utils.RandString(32))

	// Every time a new request arrives, this handler will create a random
	// username.
	f.StateUsername = utils.RandString(12)

	// Create the state token. Every request has a unique symmetric and username
	// value.
	stateToken, _ := tokenMaker.CreateToken(f.StateUsername, time.Minute)

	// Save the current token maker for the redirect URI from the grant
	// authorization process
	f.StateTokenInterface = tokenMaker

	switch c.Query("type") {
	case "google":
		authorizationUrl := f.GoogleOAuth2.Interface.GrantAuthorizationUrl(
			stateToken)

		return c.Redirect(authorizationUrl)
	case "facebook":
		authorizationUrl := f.FacebookOAuth2.Interface.GrantAuthorizationUrl(
			stateToken)

		return c.Redirect(authorizationUrl)
	default:
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   entity.ErrInvalidApiQueryLoginType.Error(),
			"code":    400,
		})
	}

}

// This is the second phase for the OAuth2 authentication process. Here we will
// exchange the authorization code obtained in the first phase to get the token.
// with this token we will be able to login in our API.
func (f *UsersApi) OAuth2Redirect(c *fiber.Ctx) error {
	var userInfo *entity.UserInfoOAuth2
	var err error

	if c.Query("state") == "" {
		return c.Status(403).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiAuthorization.Error(),
			"error":   entity.ErrInvalidApiQueryStateBlank.Error(),
			"code":    403,
		})
	}

	receivedPayload, err := f.StateTokenInterface.VerifyToken(c.Query("state"))
	if err != nil {
		return c.Status(403).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiAuthorization.Error(),
			"error":   entity.ErrInvalidApiQueryState.Error() + err.Error(),
			"code":    403,
		})
	}

	if receivedPayload.Username != f.StateUsername || f.StateUsername == "" {
		return c.Status(403).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiAuthorization.Error(),
			"error":   entity.ErrInvalidApiQueryStateDoesNotMatch.Error(),
			"code":    403,
		})
	}

	// After the verification that the state has a valid token (our API was
	// responsible for sending this request, we return the StateUsername to its
	// default value, which is a empty string. So, the only way that this handler
	// does not return an error is when the SignInOAuth handlers is called first.
	// In this case, even if for some reason a malicius person steals the OAuth2
	// code, he will not be able to request for this handler.
	f.StateUsername = ""

	switch c.Params("company") {
	case "google":
		if c.Query("code") == "" {
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": entity.ErrMessageApiRequest.Error(),
				"error":   entity.ErrInvalidApiQueryOAuth2CodeBlank.Error(),
				"code":    400,
			})
		}
		googleUserInfo, err := f.GoogleOAuth2.Interface.GrantAccessToken(
			c.Query("code"))
		if googleUserInfo.Error != "" {
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": entity.ErrMessageApiRequest.Error(),
				"error":   strings.ToUpper(googleUserInfo.Error),
				"code":    400,
			})
		}

		// Login in the Firebase with the OAuth information
		userInfo, err = f.ApplicationLogic.UserApp.UserLoginOAuth2(
			f.FirebaseWebKey, googleUserInfo.IdToken, "google.com",
			f.GoogleOAuth2.Config.RedirectURI)
		if err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": entity.ErrMessageApiRequest.Error(),
				"error":   err.Error(),
				"code":    400,
			})

		}
		break
	case "facebook":
		if c.Query("code") == "" {
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": entity.ErrMessageApiRequest.Error(),
				"error":   entity.ErrInvalidApiQueryOAuth2CodeBlank.Error(),
				"code":    400,
			})
		}

		facebookUserInfo, err := f.FacebookOAuth2.Interface.GrantAccessToken(
			c.Query("code"))
		if facebookUserInfo.Error != nil {
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": entity.ErrMessageApiRequest.Error(),
				"error":   "ERROR",
				"code":    400,
			})
		}

		// Login in the Firebase with the OAuth information
		userInfo, err = f.ApplicationLogic.UserApp.UserLoginOAuth2(
			f.FirebaseWebKey, facebookUserInfo.AccessToken, "facebook.com",
			f.FacebookOAuth2.Config.RedirectURI)
		if err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": entity.ErrMessageApiRequest.Error(),
				"error":   err.Error(),
				"code":    400,
			})

		}
		break
	default:
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   entity.ErrInvalidApiParamsCompany.Error(),
			"code":    400,
		})
	}

	// Verify if the user already exists in our database. If not, we need to
	// create.
	if userInfo.IsNewUser == true {
		// Create User in our database
		_, err = f.ApplicationLogic.UserApp.CreateUser(userInfo.UserUid,
			userInfo.Email, userInfo.Fullname, "normal")
		if err != nil {
			return c.Status(500).JSON(&fiber.Map{
				"success": false,
				"message": entity.ErrMessageApiInternalError.Error(),
				"error":   err.Error(),
				"code":    500,
			})
		}
	}

	userLoginApiReturn := presenter.ConvertUserLoginToUserLoginApiReturn(
		userInfo.Email, userInfo.Fullname, userInfo.IdToken,
		userInfo.RefreshToken, userInfo.Expiration)

	err = c.JSON(&fiber.Map{
		"success":  true,
		"userInfo": userLoginApiReturn,
		"message":  "User login was successful",
	})

	return nil
}

func (f *UsersApi) ForgotPassword(c *fiber.Ctx) error {

	var passwordResetEmail presenter.ForgotPasswordBody

	if err := c.BodyParser(&passwordResetEmail); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   entity.ErrInvalidApiBody.Error(),
			"code":    400,
		})
	}

	// Send Email to reset password
	emailForgotPassResp, err := f.ApplicationLogic.UserApp.
		UserSendForgotPasswordEmail(f.FirebaseWebKey, passwordResetEmail.Email)
	if err != nil {
		if err.Error() == "EMAIL_NOT_FOUND" {
			return c.Status(404).JSON(&fiber.Map{
				"success": false,
				"message": entity.ErrMessageApiEmail.Error(),
				"code":    404,
			})
		} else {
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": entity.ErrMessageApiRequest.Error(),
				"error":   err.Error(),
				"code":    400,
			})
		}
	}

	err = c.JSON(&fiber.Map{
		"success":  true,
		"userInfo": emailForgotPassResp,
		"message":  "The email for password reset was sent successfully",
	})

	return err
}

func (f *UsersApi) RefreshIdToken(c *fiber.Ctx) error {
	var userRefreshToken presenter.UserRefreshIdTokenBody

	if err := c.BodyParser(&userRefreshToken); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   entity.ErrInvalidApiBody.Error(),
			"code":    400,
		})
	}

	refreshTokenResponse, err := f.ApplicationLogic.UserApp.UserRefreshIdToken(
		f.FirebaseWebKey, userRefreshToken.RefreshToken)

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   err.Error(),
			"code":    400,
		})
	}

	refreshTokenApiReturn := presenter.
		ConvertUserRefreshTokenToUserRefreshTokenApiReturn(
			refreshTokenResponse.RefreshToken, refreshTokenResponse.IdToken,
			refreshTokenResponse.TokenType, refreshTokenResponse.Expiration)

	err = c.JSON(&fiber.Map{
		"success":          true,
		"userRefreshToken": refreshTokenApiReturn,
		"message":          "The token was updated successfully",
	})

	return err
}

func (f *UsersApi) DeleteUser(c *fiber.Ctx) error {

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	deletedUser, err := f.ApplicationLogic.UserApp.DeleteUser(userId.String())
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   err.Error(),
			"code":    400,
		})
	}

	err = c.JSON(&fiber.Map{
		"success":  true,
		"userInfo": deletedUser,
		"message":  "User was deleted successfully",
	})

	return err
}

func (f *UsersApi) UpdateUserInfo(c *fiber.Ctx) error {

	var userInfoUpdate presenter.SignUpBody

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	if err := c.BodyParser(&userInfoUpdate); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   entity.ErrInvalidApiBody.Error(),
			"code":    400,
		})
	}

	userUpdated, err := f.ApplicationLogic.UserApp.UpdateUser(userId.String(),
		userInfoUpdate.Email, userInfoUpdate.DisplayName,
		userInfoUpdate.Password)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   err.Error(),
			"code":    400,
		})
	}

	userApiReturn := presenter.ConvertUserToUserApiReturn(userUpdated.Email,
		userInfoUpdate.DisplayName)

	err = c.JSON(&fiber.Map{
		"success":  true,
		"userInfo": userApiReturn,
		"message":  "User information was updated successfully",
	})

	return err
}
