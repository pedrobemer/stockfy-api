package fiberHandlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"stockfyApi/api/middleware"
	"stockfyApi/api/presenter"
	"stockfyApi/entity"
	"stockfyApi/usecases"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestApiSectorGet(t *testing.T) {

	type body struct {
		Success bool           `json:"success"`
		Message string         `json:"message"`
		Error   string         `json:"error"`
		Code    int            `json:"code"`
		Sector  *entity.Sector `json:"sector"`
	}

	type test struct {
		idToken      string
		sectorName   string
		expectedResp body
	}

	tests := []test{
		{
			idToken:    "ValidIdTokenWithoutPrivilegedUser",
			sectorName: "ERROR_NAME",
			expectedResp: body{
				Code:    500,
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Error:   errors.New("Some Error").Error(),
				Sector:  nil,
			},
		},
		{
			idToken:    "ValidIdTokenWithoutPrivilegedUser",
			sectorName: "INVALID_NAME",
			expectedResp: body{
				Code:    404,
				Success: false,
				Message: entity.ErrMessageApiSectorName.Error(),
				Error:   entity.ErrInvalidSectorSearchName.Error(),
				Sector:  nil,
			},
		},
		{
			idToken:    "ValidIdTokenWithoutPrivilegedUser",
			sectorName: "Test",
			expectedResp: body{
				Code:    200,
				Success: true,
				Message: "Sector information returned successfully",
				Error:   "",
				Sector: &entity.Sector{
					Id:   "TestID",
					Name: "Test",
				},
			},
		},
		{
			idToken:    "ValidIdTokenWithoutEmailVerification",
			sectorName: "Test",
			expectedResp: body{
				Code:    401,
				Success: false,
				Message: entity.ErrMessageApiAuthentication.Error(),
				Sector:  nil,
				Error:   "",
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()

	// Declare Sector Application Logic
	sector := SectorApi{
		ApplicationLogic: *usecases,
	}

	// Mock HTTP request
	app := fiber.New()
	api := app.Group("/api")
	api.Use(middleware.NewFiberMiddleware(middleware.FiberMiddleware{
		UserAuthentication: usecases.UserApp,
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
	api.Get("/sector/:sector", sector.GetSector)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "GET", "/api/sector/"+testCase.sectorName,
			"application/json", testCase.idToken, nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}

}

func TestApiSectorCreateSector(t *testing.T) {

	type body struct {
		Success bool              `json:"success"`
		Message string            `json:"message"`
		Error   string            `json:"error"`
		Code    int               `json:"code"`
		Sector  *presenter.Sector `json:"sector"`
	}

	type bodyRequest struct {
		Sector string `json:"sector"`
	}

	type test struct {
		idToken      string
		bodyRequest  bodyRequest
		expectedResp body
	}

	tests := []test{
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			bodyRequest: bodyRequest{
				Sector: "Test Sector",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiAuthorization.Error(),
				Error:   entity.ErrInvalidUserAdminPrivilege.Error(),
				Code:    403,
				Sector:  nil,
			},
		},
		{
			idToken: "ValidIdTokenPrivilegeUser",
			bodyRequest: bodyRequest{
				Sector: "Test Sector",
			},
			expectedResp: body{
				Success: true,
				Message: "Sector creation was successful",
				Error:   "",
				Code:    200,
				Sector: &presenter.Sector{
					Id:   "TestID",
					Name: "Test Sector",
				},
			},
		},
		{
			idToken: "ValidIdTokenPrivilegeUser",
			bodyRequest: bodyRequest{
				Sector: "ERROR_SECTOR",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Error:   errors.New("Some Error").Error(),
				Code:    500,
				Sector:  nil,
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()

	// Declare Sector Application Logic
	sector := SectorApi{
		ApplicationLogic: *usecases,
	}

	// Mock HTTP request
	app := fiber.New()
	api := app.Group("/api")
	api.Use(middleware.NewFiberMiddleware(middleware.FiberMiddleware{
		UserAuthentication: usecases.UserApp,
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
	api.Post("/sector", sector.CreateSector)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "POST", "/api/sector", "application/json",
			testCase.idToken, testCase.bodyRequest)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}

}
