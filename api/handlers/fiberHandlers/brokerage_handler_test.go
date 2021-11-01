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

func TestApiGetBrokerageFirms(t *testing.T) {
	type body struct {
		Success   bool                  `json:"success"`
		Message   string                `json:"message"`
		Error     string                `json:"error"`
		Code      int                   `json:"code"`
		Brokerage []presenter.Brokerage `json:"brokerage"`
	}

	type test struct {
		idToken      string
		contentType  string
		path         string
		expectedResp body
	}

	tests := []test{
		{
			idToken:     "ValidIdTokenWithoutEmailVerification",
			contentType: "application/json",
			expectedResp: body{
				Code:      401,
				Success:   false,
				Message:   entity.ErrMessageApiAuthentication.Error(),
				Error:     "",
				Brokerage: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			path:        "country=ERROR",
			expectedResp: body{
				Code:      400,
				Success:   false,
				Message:   entity.ErrMessageApiRequest.Error(),
				Error:     entity.ErrInvalidCountryCode.Error(),
				Brokerage: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			path:        "country=ERROR_BROKERAGE_SEARCH",
			expectedResp: body{
				Code:      500,
				Success:   false,
				Message:   entity.ErrMessageApiInternalError.Error(),
				Error:     errors.New("Unknown error in the brokerage repository").Error(),
				Brokerage: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			path:        "country=US",
			expectedResp: body{
				Code:    200,
				Success: true,
				Message: "Returned successfully the brokerage firms information",
				Error:   "",
				Brokerage: []presenter.Brokerage{
					{
						Id:      "TestBrokerageID1",
						Name:    "Test " + "US" + " 1",
						Country: "US",
					},
					{
						Id:      "TestBrokerageID1",
						Name:    "Test " + "US" + " 2",
						Country: "US",
					},
				},
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			path:        "",
			expectedResp: body{
				Code:    200,
				Success: true,
				Message: "Returned successfully the brokerage firms information",
				Error:   "",
				Brokerage: []presenter.Brokerage{
					{
						Id:      "TestBrokerageID1",
						Name:    "Test US 1",
						Country: "US",
					},
					{
						Id:      "TestBrokerageID2",
						Name:    "Test US 2",
						Country: "US",
					},
					{
						Id:      "TestBrokerageID3",
						Name:    "Test BR 1",
						Country: "BR",
					},
				},
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()
	// logicApi := logicApi.NewMockApplication(*usecases)

	// Declare Sector Application Logic
	brokerage := BrokerageApi{
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
	api.Get("/brokerage", brokerage.GetBrokerageFirms)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "GET", "/api/brokerage?"+testCase.path,
			testCase.contentType, testCase.idToken, nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}

}

func TestApiGetBrokerageFirm(t *testing.T) {
	type body struct {
		Success   bool                 `json:"success"`
		Message   string               `json:"message"`
		Error     string               `json:"error"`
		Code      int                  `json:"code"`
		Brokerage *presenter.Brokerage `json:"brokerage"`
	}

	type test struct {
		idToken      string
		contentType  string
		name         string
		expectedResp body
	}

	tests := []test{
		{
			idToken:     "ValidIdTokenWithoutEmailVerification",
			contentType: "application/json",
			expectedResp: body{
				Code:      401,
				Success:   false,
				Message:   entity.ErrMessageApiAuthentication.Error(),
				Error:     "",
				Brokerage: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			name:        "UNKNOWN_BROKERAGE",
			expectedResp: body{
				Code:      404,
				Success:   false,
				Message:   entity.ErrInvalidBrokerageNameSearch.Error(),
				Error:     "",
				Brokerage: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			name:        "ERROR_BROKERAGE_SEARCH",
			expectedResp: body{
				Code:      500,
				Success:   false,
				Message:   entity.ErrMessageApiInternalError.Error(),
				Error:     errors.New("Unknown error in the brokerage repository").Error(),
				Brokerage: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			name:        "TestBrokerage",
			expectedResp: body{
				Code:    200,
				Success: true,
				Message: "Brokerage firm information returned successfully",
				Error:   "",
				Brokerage: &presenter.Brokerage{
					Id:      "TestBrokerageID1",
					Name:    "TestBrokerage",
					Country: "US",
				},
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()
	// logicApi := logicApi.NewMockApplication(*usecases)

	// Declare Sector Application Logic
	brokerage := BrokerageApi{
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
	api.Get("/brokerage/:name", brokerage.GetBrokerageFirm)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "GET", "/api/brokerage/"+testCase.name,
			testCase.contentType, testCase.idToken, nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}
}
