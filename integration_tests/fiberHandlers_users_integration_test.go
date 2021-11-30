package integration_tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"stockfyApi/api/handlers/fiberHandlers"
	"stockfyApi/api/middleware"
	"stockfyApi/api/presenter"
	"stockfyApi/database/postgresql"
	"stockfyApi/entity"
	"stockfyApi/usecases"
	"stockfyApi/usecases/user"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	var err error
	err = godotenv.Load(os.ExpandEnv("../database-dev.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}

	os.Exit(m.Run())
}

func connectDatabase() *pgx.Conn {

	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")
	DB_PORT := os.Getenv("DB_PORT")
	DB_HOST := os.Getenv("DB_HOST")

	dbinfo := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)

	DBpool, err := pgx.Connect(context.Background(), dbinfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return DBpool
}

func TestFiberHandlersIntegrationTestSignUp(t *testing.T) {

	type body struct {
		Success  bool                     `json:"success"`
		Message  string                   `json:"message"`
		Error    string                   `json:"error"`
		Code     int                      `json:"code"`
		UserInfo *presenter.UserApiReturn `json:"userInfo"`
	}

	type test struct {
		contentType      string
		email            string
		password         string
		displayName      string
		expectedResponse body
	}

	tests := []test{
		{
			contentType: "application/pdf",
			email:       "",
			password:    "PasswdTest",
			displayName: "Test Username",
			expectedResponse: body{
				Success:  false,
				Message:  entity.ErrMessageApiRequest.Error(),
				Error:    entity.ErrInvalidApiBody.Error(),
				Code:     400,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			email:       "",
			password:    "PasswdTest",
			displayName: "Test Username",
			expectedResponse: body{
				Success:  false,
				Message:  entity.ErrMessageApiRequest.Error(),
				Error:    errors.New("email must be a non-empty string").Error(),
				Code:     400,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			email:       "test@email.com",
			password:    "",
			displayName: "Test Username",
			expectedResponse: body{
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
			email:       "test@email.com",
			password:    "PasswdTest",
			displayName: "",
			expectedResponse: body{
				Success:  false,
				Message:  entity.ErrMessageApiRequest.Error(),
				Error:    errors.New("display name must be a non-empty string").Error(),
				Code:     400,
				UserInfo: nil,
			},
		},
		{
			contentType: "application/json",
			email:       "test@email.com",
			password:    "testPasswd",
			displayName: "Test Name",
			expectedResponse: body{
				Success: true,
				Message: "User was registered successfully",
				Code:    200,
				UserInfo: &presenter.UserApiReturn{
					Email:       "test@email.com",
					DisplayName: "Test Name",
				},
			},
		},
		{
			contentType: "application/json",
			email:       "test@email.com",
			password:    "testPasswd",
			displayName: "Test Name",
			expectedResponse: body{
				Success:  false,
				Message:  entity.ErrMessageApiInternalError.Error(),
				Error:    "entity.CreateUser: scany: rows final error: ERROR: duplicate key value violates unique constraint \"users_pk\" (SQLSTATE 23505)",
				Code:     500,
				UserInfo: nil,
			},
		},
	}

	DBpool := connectDatabase()

	dbInterfaces := postgresql.NewPostgresInstance(DBpool)
	firebaseInterface := user.NewExternalApi()

	applicationLogics := usecases.NewApplications(dbInterfaces,
		firebaseInterface)

	users := fiberHandlers.UsersApi{
		ApplicationLogic: *applicationLogics,
	}

	app := fiber.New()
	api := app.Group("/api")
	api.Post("/signup", users.SignUp)

	for _, testCase := range tests {
		bodyResponse := body{}
		bodyRequestStruct := presenter.SignUpBody{
			Email:       testCase.email,
			Password:    testCase.password,
			DisplayName: testCase.displayName,
		}

		resp, _ := fiberHandlers.MockHttpRequest(app, "POST", "/api/signup",
			testCase.contentType, "", bodyRequestStruct)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &bodyResponse)
		bodyResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResponse, bodyResponse)

	}

}

func TestFiberHandlersIntegrationTestUpdateUser(t *testing.T) {
	type body struct {
		Success  bool                     `json:"success"`
		Message  string                   `json:"message"`
		Error    string                   `json:"error"`
		Code     int                      `json:"code"`
		UserInfo *presenter.UserApiReturn `json:"userInfo"`
	}

	type test struct {
		idToken          string
		contentType      string
		email            string
		password         string
		displayName      string
		expectedResponse body
	}

	tests := []test{
		{
			idToken:     "INVALID_ID_TOKEN",
			contentType: "application/json",
			expectedResponse: body{
				Code:     401,
				Success:  false,
				Message:  entity.ErrMessageApiAuthentication.Error(),
				Error:    "",
				UserInfo: nil,
			},
		},
		{
			idToken:     "TestNormalID",
			contentType: "application/pdf",
			email:       "Test Name",
			password:    "PasswdTest",
			displayName: "Test Name",
			expectedResponse: body{
				Success:  false,
				Message:  entity.ErrMessageApiRequest.Error(),
				Error:    entity.ErrInvalidApiBody.Error(),
				Code:     400,
				UserInfo: nil,
			},
		},
		{
			idToken:     "TestNormalID",
			displayName: "ERROR_USER_FIREBASE",
			email:       "test@email.com",
			password:    "PasswdChange",
			contentType: "application/json",
			expectedResponse: body{
				Code:     400,
				Success:  false,
				Message:  entity.ErrMessageApiRequest.Error(),
				Error:    errors.New("Unknown update error in the user repository").Error(),
				UserInfo: nil,
			},
		},
		{
			idToken:     "TestNormalID",
			displayName: "Test Name Change",
			email:       "test@email.com",
			password:    "PasswdChange",
			contentType: "application/json",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "User information was updated successfully",
				Error:   "",
				UserInfo: &presenter.UserApiReturn{
					Email:       "test@email.com",
					DisplayName: "Test Name Change",
				},
			},
		},
	}

	DBpool := connectDatabase()

	dbInterfaces := postgresql.NewPostgresInstance(DBpool)
	firebaseInterface := user.NewExternalApi()

	applicationLogics := usecases.NewApplications(dbInterfaces,
		firebaseInterface)

	users := fiberHandlers.UsersApi{
		ApplicationLogic: *applicationLogics,
	}

	// Mock HTTP request
	app := fiber.New()
	api := app.Group("/api")
	api.Use(middleware.NewFiberMiddleware(middleware.FiberMiddleware{
		UserAuthentication: applicationLogics.UserApp,
		ErrorHandler: func(c *fiber.Ctx, e error) error {
			var err error
			c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": entity.ErrMessageApiAuthentication.Error(),
				"code":    401,
			})

			return err
		},
		ContextKey: "user",
	}))
	api.Put("/update-user", users.UpdateUserInfo)

	for _, testCase := range tests {
		jsonResponse := body{}
		bodyRequestStruct := presenter.SignUpBody{
			Email:       testCase.email,
			Password:    testCase.password,
			DisplayName: testCase.displayName,
		}
		resp, _ := fiberHandlers.MockHttpRequest(app, "PUT", "/api/update-user",
			testCase.contentType, testCase.idToken, bodyRequestStruct)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResponse, jsonResponse)
	}

}

func TestFiberHandlersIntegrationTestDeleteUser(t *testing.T) {
	type body struct {
		Success  bool             `json:"success"`
		Message  string           `json:"message"`
		Error    string           `json:"error"`
		Code     int              `json:"code"`
		UserInfo *entity.UserInfo `json:"userInfo"`
	}

	type test struct {
		idToken          string
		contentType      string
		expectedResponse body
	}

	tests := []test{
		{
			idToken:     "INVALID_ID_TOKEN",
			contentType: "application/json",
			expectedResponse: body{
				Code:     401,
				Success:  false,
				Message:  entity.ErrMessageApiAuthentication.Error(),
				Error:    "",
				UserInfo: nil,
			},
		},
		{
			idToken:     "Invalid",
			contentType: "application/json",
			expectedResponse: body{
				Code:     400,
				Success:  false,
				Message:  entity.ErrMessageApiRequest.Error(),
				Error:    errors.New("Database Interface error").Error(),
				UserInfo: nil,
			},
		},
		{
			idToken:     "TestNormalID",
			contentType: "application/json",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "User was deleted successfully",
				Error:   "",
				UserInfo: &entity.UserInfo{
					UID:         "TestNormalID",
					Email:       "test@email.com",
					DisplayName: "Test Name",
				},
			},
		},
	}

	DBpool := connectDatabase()

	dbInterfaces := postgresql.NewPostgresInstance(DBpool)
	firebaseInterface := user.NewExternalApi()

	applicationLogics := usecases.NewApplications(dbInterfaces,
		firebaseInterface)

	users := fiberHandlers.UsersApi{
		ApplicationLogic: *applicationLogics,
	}

	// Mock HTTP request
	app := fiber.New()
	api := app.Group("/api")
	api.Use(middleware.NewFiberMiddleware(middleware.FiberMiddleware{
		UserAuthentication: applicationLogics.UserApp,
		ErrorHandler: func(c *fiber.Ctx, e error) error {
			var err error
			c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": entity.ErrMessageApiAuthentication.Error(),
				"code":    401,
			})

			return err
		},
		ContextKey: "user",
	}))
	api.Delete("/delete-user", users.DeleteUser)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := fiberHandlers.MockHttpRequest(app, "DELETE", "/api/delete-user",
			testCase.contentType, testCase.idToken, nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResponse, jsonResponse)
	}

}
