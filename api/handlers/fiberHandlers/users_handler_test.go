package fiberHandlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"stockfyApi/api/presenter"
	"stockfyApi/entity"
	"stockfyApi/usecases"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestApiUsersSignUp(t *testing.T) {

	type body struct {
		Success  bool                     `json:"success"`
		Message  string                   `json:"message"`
		Error    string                   `json:"error"`
		Code     int                      `json:"code"`
		UserInfo *presenter.UserApiReturn `json:"userInfo"`
	}

	type test struct {
		contentType  string
		bodyReq      presenter.SignUpBody
		expectedResp body
	}

	tests := []test{
		{
			contentType: "application/pdf",
			bodyReq: presenter.SignUpBody{
				Email:       "",
				Password:    "PasswdTest",
				DisplayName: "Test Username",
			},
			expectedResp: body{
				Success:  false,
				Message:  entity.ErrMessageApiRequest.Error(),
				Error:    entity.ErrInvalidApiBody.Error(),
				Code:     400,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.SignUpBody{
				Email:       "",
				Password:    "PasswdTest",
				DisplayName: "Test Username",
			},
			expectedResp: body{
				Success:  false,
				Message:  entity.ErrMessageApiRequest.Error(),
				Error:    errors.New("email must be a non-empty string").Error(),
				Code:     400,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.SignUpBody{
				Email:       "test@email.com",
				Password:    "",
				DisplayName: "Test Username",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error: errors.New("password must be a string at least 6 characters long").
					Error(),
				Code:     400,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.SignUpBody{
				Email:       "test@email.com",
				Password:    "PasswdTest",
				DisplayName: "",
			},
			expectedResp: body{
				Success:  false,
				Message:  entity.ErrMessageApiRequest.Error(),
				Error:    errors.New("display name must be a non-empty string").Error(),
				Code:     400,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.SignUpBody{
				Email:       "test@email.com",
				Password:    "PasswdTest",
				DisplayName: "",
			},
			expectedResp: body{
				Success:  false,
				Message:  entity.ErrMessageApiRequest.Error(),
				Error:    errors.New("display name must be a non-empty string").Error(),
				Code:     400,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.SignUpBody{
				Email:       "test@email.com",
				Password:    "PasswdTest",
				DisplayName: "WRONG_CUSTOM_TOKEN",
			},
			expectedResp: body{
				Success:  false,
				Message:  entity.ErrMessageApiInternalError.Error(),
				Error:    errors.New("Some Error").Error(),
				Code:     500,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.SignUpBody{
				Email:       "test@email.com",
				Password:    "PasswdTest",
				DisplayName: "WRONG_ID_TOKEN",
			},
			expectedResp: body{
				Success:  false,
				Message:  entity.ErrMessageApiInternalError.Error(),
				Error:    entity.ErrInvalidUserToken.Error(),
				Code:     500,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.SignUpBody{
				Email:       "test@email.com",
				Password:    "PasswdTest",
				DisplayName: "WRONG_EMAIL_VERIFICATION",
			},
			expectedResp: body{
				Success:  false,
				Message:  entity.ErrMessageApiRequest.Error(),
				Error:    "INVALID_ID_TOKEN",
				Code:     400,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.SignUpBody{
				Email:       "test@email.com",
				Password:    "PasswdTest",
				DisplayName: "WRONG_USER_INFO",
			},
			expectedResp: body{
				Success:  false,
				Message:  entity.ErrMessageApiInternalError.Error(),
				Error:    entity.ErrInvalidUserEmailBlank.Error(),
				Code:     500,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.SignUpBody{
				Email:       "test@email.com",
				Password:    "PasswdTest",
				DisplayName: "Test Username",
			},
			expectedResp: body{
				Success: true,
				Message: "User was registered successfully",
				Error:   "",
				Code:    200,
				UserInfo: &presenter.UserApiReturn{
					Email:       "test@email.com",
					DisplayName: "Test Username",
				},
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()

	// Declare Sector Application Logic
	users := FirebaseApi{
		ApplicationLogic: *usecases,
	}

	// Mock HTTP request
	app := fiber.New()
	api := app.Group("/api")
	api.Post("/signup", users.SignUp)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "POST", "/api/signup",
			testCase.contentType, "", testCase.bodyReq)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}

}

func TestApiForgotPassword(t *testing.T) {
	type body struct {
		Success     bool                                `json:"success"`
		Message     string                              `json:"message"`
		Error       *map[string]interface{}             `json:"error"`
		ErrorString string                              `json:"error"`
		Code        int                                 `json:"code"`
		UserInfo    *entity.EmailForgotPasswordResponse `json:"userInfo"`
	}

	type test struct {
		contentType  string
		bodyReq      presenter.ForgotPasswordBody
		expectedResp body
	}

	tests := []test{
		{
			contentType: "application/pdf",
			bodyReq: presenter.ForgotPasswordBody{
				Email: "",
			},
			expectedResp: body{
				Success:     false,
				Message:     entity.ErrMessageApiRequest.Error(),
				ErrorString: entity.ErrInvalidApiBody.Error(),
				Code:        400,
				UserInfo:    nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.ForgotPasswordBody{
				Email: "",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error: &map[string]interface{}{
					"code":    400,
					"message": "MISSING_EMAIL",
					"errors": map[string]interface{}{
						"message": "MISSING_EMAIL",
						"domain":  "global",
						"reason":  "invalid",
					},
				},
				Code:     400,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.ForgotPasswordBody{
				Email: "INVALID_EMAIL",
			},
			expectedResp: body{
				Success:  false,
				Message:  entity.ErrMessageApiEmail.Error(),
				Error:    nil,
				Code:     404,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			bodyReq: presenter.ForgotPasswordBody{
				Email: "test@email.com",
			},
			expectedResp: body{
				Success: true,
				Message: "The email for password reset was sent successfully",
				Error:   nil,
				Code:    200,
				UserInfo: &entity.EmailForgotPasswordResponse{
					Email: "test@email.com",
				},
			},
		},
		// {
		// 	bodyReq: presenter.ForgotPasswordBody{
		// 		Email: "test@email.com",
		// 	},
		// 	expectedResp: body{
		// 		Success:  false,
		// 		Message:  entity.ErrMessageApiRequest.Error(),
		// 		Error:    errors.New("display name must be a non-empty string").Error(),
		// 		Code:     400,
		// 		UserInfo: nil,
		// 	},
		// },
		// {
		// 	bodyReq: presenter.ForgotPasswordBody{
		// 		Email: "test@email.com",
		// 	},
		// 	expectedResp: body{
		// 		Success: true,
		// 		Message: "User was registered successfully",
		// 		Error:   "",
		// 		Code:    200,
		// 		UserInfo: &presenter.UserApiReturn{
		// 			Email:       "test@email.com",
		// 			DisplayName: "Test Username",
		// 		},
		// 	},
		// },
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()

	// Declare Sector Application Logic
	users := FirebaseApi{
		ApplicationLogic: *usecases,
	}

	// Mock HTTP request
	app := fiber.New()
	api := app.Group("/api")
	api.Post("/forgot-password", users.ForgotPassword)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "POST", "/api/forgot-password",
			testCase.contentType, "", testCase.bodyReq)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp.Code, jsonResponse.Code)
		assert.Equal(t, testCase.expectedResp.Success, jsonResponse.Success)
		assert.Equal(t, testCase.expectedResp.Message, jsonResponse.Message)
		if jsonResponse.Error != nil {
			errorExpResp := *testCase.expectedResp.Error
			errorResp := *jsonResponse.Error
			assert.Equal(t, errorExpResp["message"], errorResp["message"])
		}

	}

}
