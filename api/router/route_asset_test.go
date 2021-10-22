package router

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"stockfyApi/api/handlers/fiberHandlers"
	"stockfyApi/api/middleware"
	"stockfyApi/api/presenter"
	"stockfyApi/entity"
	"stockfyApi/usecases"
	"stockfyApi/usecases/logicApi"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestApiAssetGet(t *testing.T) {
	type body struct {
		Success bool                      `json:"success"`
		Message string                    `json:"message"`
		Error   string                    `json:"error"`
		Code    int                       `json:"code"`
		Asset   *presenter.AssetApiReturn `json:"asset"`
	}

	type test struct {
		idToken      string
		path         string
		expectedResp body
	}

	layOut := "2006-01-02"
	dateFormatted, _ := time.Parse(layOut, "2021-10-01")
	tests := []test{
		{
			idToken: "ValidIdTokenWithoutEmailVerification",
			path:    "TEST3?withOrders=true&withOrderResume=true",
			expectedResp: body{
				Code:    401,
				Success: false,
				Message: entity.ErrMessageApiAuthentication.Error(),
				Asset:   nil,
				Error:   "",
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			path:    "TEST3?withOrders=error&withOrderResume=true",
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Asset:   nil,
				Error:   entity.ErrInvalidApiQueryWithOrders.Error(),
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			path:    "TEST3?withOrders=true&withOrderResume=error",
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Asset:   nil,
				Error:   entity.ErrInvalidApiQueryWithOrderResume.Error(),
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			path:    "UNKNOWN_SYMBOL?withOrders=true&withOrderResume=true",
			expectedResp: body{
				Code:    404,
				Success: false,
				Message: entity.ErrMessageApiAssetSymbolUser.Error(),
				Asset:   nil,
				Error:   "",
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			path:    "ERROR_REPOSITORY?withOrders=true&withOrderResume=true",
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Asset:   nil,
				Error:   errors.New("Unknown repository error").Error(),
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			path:    "TEST3?withOrders=true&withOrderResume=true",
			expectedResp: body{
				Code:    200,
				Success: true,
				Message: "Asset information returned successfully",
				Asset: &presenter.AssetApiReturn{
					Id:         "TestID",
					Symbol:     "TEST3",
					Preference: "TestPref",
					Fullname:   "Test Name",
					AssetType: &presenter.AssetType{
						Id:      "TestAssetTypeID",
						Type:    "ETF",
						Name:    "Test ETF",
						Country: "BR",
					},
					Sector: &presenter.Sector{
						Id:   "TestSectorID",
						Name: "Test Sector",
					},
					Orders: &[]presenter.OrderApiReturn{
						{
							Id:        "Order1",
							Quantity:  2,
							Price:     29.29,
							Currency:  "USD",
							OrderType: "Dividendos",
							Date:      dateFormatted,
							Brokerage: &presenter.Brokerage{
								Id:      "BrokerageID",
								Name:    "Test Broker",
								Country: "US",
							},
						},
						{
							Id:        "Order2",
							Quantity:  2,
							Price:     29.29,
							Currency:  "USD",
							OrderType: "Dividendos",
							Date:      dateFormatted,
							Brokerage: &presenter.Brokerage{
								Id:      "BrokerageID",
								Name:    "Test Broker",
								Country: "US",
							},
						},
					},
					OrderInfos: &presenter.OrderInfos{
						TotalQuantity:        4,
						WeightedAdjPrice:     28.20,
						WeightedAveragePrice: 29.29,
					},
				},
				Error: "",
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()

	// Declare Sector Application Logic
	asset := fiberHandlers.AssetApi{
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
	api.Get("/asset/:symbol", asset.GetAsset)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "GET", "/api/asset/"+testCase.path,
			testCase.idToken, nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}
}

func TestApiAssetPost(t *testing.T) {

	type body struct {
		Success bool                      `json:"success"`
		Message string                    `json:"message"`
		Error   string                    `json:"error"`
		Code    int                       `json:"code"`
		Asset   *presenter.AssetApiReturn `json:"asset"`
	}

	type bodyRequest struct {
		AssetType string `json:"assetType"`
		Symbol    string `json:"symbol"`
		Fullname  string `json:"fullname"`
		Country   string `json:"country"`
	}

	type test struct {
		idToken      string
		bodyReq      bodyRequest
		expectedResp body
	}

	tests := []test{
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			bodyReq: bodyRequest{
				AssetType: "ETF",
				Symbol:    "TEST11",
				Fullname:  "Test Company",
				Country:   "BR",
			},
			expectedResp: body{
				Code:    403,
				Success: false,
				Message: entity.ErrMessageApiAuthorization.Error(),
				Asset:   nil,
				Error:   entity.ErrInvalidUserAdminPrivilege.Error(),
			},
		},
		{
			idToken: "ValidIdTokenPrivilegeUser",
			bodyReq: bodyRequest{
				AssetType: "ETF",
				Symbol:    "SYMBOL_EXIST",
				Fullname:  "Test Company",
				Country:   "BR",
			},
			expectedResp: body{
				Code:    403,
				Success: false,
				Message: entity.ErrMessageApiAuthorization.Error(),
				Asset:   nil,
				Error:   entity.ErrInvalidAssetSymbolExist.Error(),
			},
		},
		{
			idToken: "ValidIdTokenPrivilegeUser",
			bodyReq: bodyRequest{
				AssetType: "ETF",
				Symbol:    "TEST11",
				Fullname:  "Test Company",
				Country:   "ERROR",
			},
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Asset:   nil,
				Error:   entity.ErrInvalidCountryCode.Error(),
			},
		},
		{
			idToken: "ValidIdTokenPrivilegeUser",
			bodyReq: bodyRequest{
				AssetType: "ETF",
				Symbol:    "UNKNOWN_SYMBOL",
				Fullname:  "Test Company",
				Country:   "BR",
			},
			expectedResp: body{
				Code:    404,
				Success: false,
				Message: entity.ErrMessageApiAssetSymbolUser.Error(),
				Asset:   nil,
				Error:   "",
			},
		},
		{
			idToken: "ValidIdTokenPrivilegeUser",
			bodyReq: bodyRequest{
				AssetType: "ETF",
				Symbol:    "ERROR_SECTOR_REPO",
				Fullname:  "Test Company",
				Country:   "BR",
			},
			expectedResp: body{
				Code:    500,
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Asset:   nil,
				Error:   errors.New("Unknown sector repository error").Error(),
			},
		},
		{
			idToken: "ValidIdTokenPrivilegeUser",
			bodyReq: bodyRequest{
				AssetType: "ETF",
				Symbol:    "ERROR_ASSETTYPE_REPO",
				Fullname:  "Test Company",
				Country:   "BR",
			},
			expectedResp: body{
				Code:    500,
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Asset:   nil,
				Error:   errors.New("Unknown asset type repository error").Error(),
			},
		},
		{
			idToken: "ValidIdTokenPrivilegeUser",
			bodyReq: bodyRequest{
				AssetType: "ETF",
				Symbol:    "ERROR_ASSET_REPO",
				Fullname:  "Test Company",
				Country:   "BR",
			},
			expectedResp: body{
				Code:    500,
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Asset:   nil,
				Error:   errors.New("Unknown asset repository error").Error(),
			},
		},
		{
			idToken: "ValidIdTokenPrivilegeUser",
			bodyReq: bodyRequest{
				AssetType: "ETF",
				Symbol:    "TEST11",
				Fullname:  "Test Company",
				Country:   "BR",
			},
			expectedResp: body{
				Code:    200,
				Success: true,
				Message: "Asset creation was sucessful",
				Asset: &presenter.AssetApiReturn{
					Id:         "TestID",
					Symbol:     "TEST11",
					Preference: "TestPref",
					Fullname:   "Test Name",
				},
				Error: "",
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()
	logicApi := logicApi.NewMockApplication()

	// Declare Sector Application Logic
	asset := fiberHandlers.AssetApi{
		ApplicationLogic: *usecases,
		LogicApi:         logicApi,
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
	api.Post("/asset", asset.CreateAsset)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "POST", "/api/asset", testCase.idToken,
			testCase.bodyReq)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}
}

func TestApiAssetDelete(t *testing.T) {

	type body struct {
		Success bool                      `json:"success"`
		Message string                    `json:"message"`
		Error   string                    `json:"error"`
		Code    int                       `json:"code"`
		Asset   *presenter.AssetApiReturn `json:"asset"`
	}

	type test struct {
		idToken      string
		path         string
		expectedResp body
	}

	assetApiReturn := presenter.AssetApiReturn{
		Id:         "TestID",
		Symbol:     "TEST3",
		Preference: "TestPref",
		Fullname:   "Test Name",
		AssetType: &presenter.AssetType{
			Id:      "TestAssetTypeID",
			Type:    "STOCK",
			Country: "US",
			Name:    "Test ASTY Name",
		},
		Sector: &presenter.Sector{
			Id:   "TestSectorID",
			Name: "Test Sector",
		},
	}

	tests := []test{
		{
			idToken: "ValidIdTokenWithoutEmailVerification",
			path:    "TEST11?myUser=true",
			expectedResp: body{
				Code:    401,
				Success: false,
				Message: entity.ErrMessageApiAuthentication.Error(),
				Error:   "",
				Asset:   nil,
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			path:    "TEST11",
			expectedResp: body{
				Code:    403,
				Success: false,
				Message: entity.ErrMessageApiAuthorization.Error(),
				Error:   entity.ErrInvalidUserAdminPrivilege.Error(),
				Asset:   nil,
			},
		},
		{
			idToken: "ValidIdTokenPrivilegeUser",
			path:    "ERROR_ASSET_REPO",
			expectedResp: body{
				Code:    500,
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Error:   errors.New("Unknown asset repository error").Error(),
				Asset:   nil,
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			path:    "ERROR_ASSET_REPO?myUser=true",
			expectedResp: body{
				Code:    500,
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Error:   errors.New("Unknown asset repository error").Error(),
				Asset:   nil,
			},
		},
		{
			idToken: "ValidIdTokenPrivilegeUser",
			path:    "UNKNOWN_SYMBOL",
			expectedResp: body{
				Code:    404,
				Success: false,
				Message: entity.ErrMessageApiAssetSymbolUser.Error(),
				Error:   "",
				Asset:   nil,
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			path:    "UNKNOWN_SYMBOL?myUser=true",
			expectedResp: body{
				Code:    404,
				Success: false,
				Message: entity.ErrMessageApiAssetSymbolUser.Error(),
				Error:   "",
				Asset:   nil,
			},
		},
		{
			idToken: "ValidIdTokenPrivilegeUser",
			path:    "ERROR_ASSETUSER_REPO",
			expectedResp: body{
				Code:    500,
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Error:   errors.New("Unknown asset user repository error").Error(),
				Asset:   nil,
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			path:    "ERROR_ASSETUSER_REPO?myUser=true",
			expectedResp: body{
				Code:    500,
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Error:   errors.New("Unknown asset user repository error").Error(),
				Asset:   nil,
			},
		},
		{
			idToken: "ValidIdTokenPrivilegeUser",
			path:    "ERROR_ORDERS_REPO",
			expectedResp: body{
				Code:    500,
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Error:   errors.New("Unknown orders repository error").Error(),
				Asset:   nil,
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			path:    "ERROR_ORDERS_REPO?myUser=true",
			expectedResp: body{
				Code:    500,
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Error:   errors.New("Unknown orders repository error").Error(),
				Asset:   nil,
			},
		},
		{
			idToken: "ValidIdTokenPrivilegeUser",
			path:    "ERROR_EARNINGS_REPO",
			expectedResp: body{
				Code:    500,
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Error:   errors.New("Unknown earnings repository error").Error(),
				Asset:   nil,
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			path:    "ERROR_EARNINGS_REPO?myUser=true",
			expectedResp: body{
				Code:    500,
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Error:   errors.New("Unknown earnings repository error").Error(),
				Asset:   nil,
			},
		},
		{
			idToken: "ValidIdTokenPrivilegeUser",
			path:    "TEST3",
			expectedResp: body{
				Code:    200,
				Success: true,
				Message: "Asset was deleted successfuly",
				Error:   "",
				Asset:   &assetApiReturn,
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			path:    "TEST3?myUser=true",
			expectedResp: body{
				Code:    200,
				Success: true,
				Message: "Asset was deleted successfuly",
				Error:   "",
				Asset:   &assetApiReturn,
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()
	logicApi := logicApi.NewMockApplication()

	// Declare Sector Application Logic
	asset := fiberHandlers.AssetApi{
		ApplicationLogic: *usecases,
		LogicApi:         logicApi,
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
	api.Delete("/asset/:symbol", asset.DeleteAsset)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "DELETE", "/api/asset/"+testCase.path,
			testCase.idToken, nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}
}
