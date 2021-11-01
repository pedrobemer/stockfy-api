package fiberHandlers

import (
	"encoding/json"
	"io/ioutil"
	"stockfyApi/api/middleware"
	"stockfyApi/entity"
	"stockfyApi/usecases"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestApiFinnhubGetSymbol(t *testing.T) {

	type body struct {
		Success      bool                 `json:"success"`
		Message      string               `json:"message"`
		Error        string               `json:"error"`
		Code         int                  `json:"code"`
		SymbolLookup *entity.SymbolLookup `json:"symbolLookup"`
	}

	type test struct {
		idToken      string
		contentType  string
		pathQuery    string
		expectedResp body
	}

	tests := []test{
		{
			idToken:     "ValidIdTokenWithoutEmailVerification",
			contentType: "application/json",
			pathQuery:   "?symbol=ITUB4&country=BR",
			expectedResp: body{
				Code:         401,
				Success:      false,
				Message:      entity.ErrMessageApiAuthentication.Error(),
				Error:        "",
				SymbolLookup: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			pathQuery:   "?country=BR",
			expectedResp: body{
				Code:         400,
				Success:      false,
				Message:      entity.ErrMessageApiRequest.Error(),
				Error:        entity.ErrInvalidApiQuerySymbolBlank.Error(),
				SymbolLookup: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			pathQuery:   "?symbol=ITUB4&country=ERROR",
			expectedResp: body{
				Code:         400,
				Success:      false,
				Message:      entity.ErrMessageApiRequest.Error(),
				Error:        entity.ErrInvalidCountryCode.Error(),
				SymbolLookup: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			pathQuery:   "?symbol=UNKNOWN_SYMBOL&country=BR",
			expectedResp: body{
				Code:         404,
				Success:      false,
				Message:      entity.ErrInvalidAssetSymbol.Error(),
				Error:        "",
				SymbolLookup: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			pathQuery:   "?symbol=TEST3&country=BR",
			expectedResp: body{
				Code:    200,
				Success: true,
				Message: "Symbol Lookup via Finnhub returned successfully",
				Error:   "",
				SymbolLookup: &entity.SymbolLookup{
					Fullname: "Test Name",
					Symbol:   "TEST3",
					Type:     "ETP",
				},
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()
	// logicApi := logicApi.NewMockApplication(*usecases)

	// Declare Sector Application Logic
	finnhub := FinnhubApi{
		ApplicationLogic: *usecases,
		// ApiLogic:         nil,
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
	api.Get("/finnhub/symbol-lookup", finnhub.GetSymbol)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "GET", "/api/finnhub/symbol-lookup"+
			testCase.pathQuery, testCase.contentType, testCase.idToken, nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}
}

func TestApiFinnhubGetSymbolPrice(t *testing.T) {
	type body struct {
		Success     bool                `json:"success"`
		Message     string              `json:"message"`
		Error       string              `json:"error"`
		Code        int                 `json:"code"`
		SymbolPrice *entity.SymbolPrice `json:"symbolPrice"`
	}

	type test struct {
		idToken      string
		contentType  string
		pathQuery    string
		expectedResp body
	}

	tests := []test{
		{
			idToken:     "ValidIdTokenWithoutEmailVerification",
			contentType: "application/json",
			pathQuery:   "?symbol=ITUB4&country=BR",
			expectedResp: body{
				Code:        401,
				Success:     false,
				Message:     entity.ErrMessageApiAuthentication.Error(),
				Error:       "",
				SymbolPrice: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			pathQuery:   "?symbol=ITUB4&country=ERROR",
			expectedResp: body{
				Code:        400,
				Success:     false,
				Message:     entity.ErrMessageApiRequest.Error(),
				Error:       entity.ErrInvalidCountryCode.Error(),
				SymbolPrice: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			pathQuery:   "?country=BR",
			expectedResp: body{
				Code:        400,
				Success:     false,
				Message:     entity.ErrMessageApiRequest.Error(),
				Error:       entity.ErrInvalidApiQuerySymbolBlank.Error(),
				SymbolPrice: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			pathQuery:   "?symbol=UNKNOWN_SYMBOL&country=BR",
			expectedResp: body{
				Code:        404,
				Success:     false,
				Message:     entity.ErrInvalidAssetSymbol.Error(),
				Error:       "",
				SymbolPrice: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			pathQuery:   "?symbol=TEST3&country=BR",
			expectedResp: body{
				Code:    200,
				Success: true,
				Message: "Symbol Price via Finnhub returned successfully",
				Error:   "",
				SymbolPrice: &entity.SymbolPrice{
					CurrentPrice:   29.29,
					LowPrice:       28.00,
					HighPrice:      29.89,
					OpenPrice:      29.29,
					PrevClosePrice: 29.29,
					MarketCap:      1018388,
				},
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()
	// logicApi := logicApi.NewMockApplication(*usecases)

	// Declare Sector Application Logic
	finnhub := FinnhubApi{
		ApplicationLogic: *usecases,
		// ApiLogic:         nil,
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
	api.Get("/finnhub/symbol-price", finnhub.GetSymbolPrice)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "GET", "/api/finnhub/symbol-price"+
			testCase.pathQuery, testCase.contentType, testCase.idToken, nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}
}
