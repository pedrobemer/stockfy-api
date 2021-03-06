package fiberHandlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"stockfyApi/api/middleware"
	"stockfyApi/api/presenter"
	"stockfyApi/entity"
	"stockfyApi/usecases"
	"stockfyApi/usecases/logicApi"
	"testing"

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

	dateFormatted := entity.StringToTime("2021-10-01")
	tests := []test{
		{
			idToken: "ValidIdTokenWithoutEmailVerification",
			path:    "TEST3?withOrders=true&withOrderResume=true&withPrice=true",
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
			path:    "TEST3?withOrders=error&withOrderResume=true&withPrice=true",
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
			path:    "TEST3?withOrders=true&withOrderResume=error&withPrice=true",
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
			path:    "TEST3?withOrders=true&withOrderResume=true&withPrice=error",
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Asset:   nil,
				Error:   entity.ErrInvalidApiQueryWithPrice.Error(),
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
			path:    "INVALID_SYMBOL?withOrders=true&withOrderResume=true",
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Asset:   nil,
				Error:   entity.ErrInvalidAssetSymbol.Error(),
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			path:    "ERROR_ASSET_REPOSITORY?withOrders=true&withOrderResume=true",
			expectedResp: body{
				Code:    500,
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Asset:   nil,
				Error:   errors.New("Unknown error in the asset repository").Error(),
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			path:    "TEST3?withOrders=true&withOrderResume=true&withPrice=true",
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
					Orders: []presenter.OrderApiReturn{
						{
							Id:        "Order1",
							Quantity:  2,
							Price:     29.29,
							Currency:  "USD",
							OrderType: "buy",
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
							OrderType: "buy",
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
						WeightedAdjPrice:     29.29,
						WeightedAveragePrice: 29.29,
					},
					Price: &presenter.AssetPrice{
						OpenPrice:   200.19,
						ActualPrice: 199.98,
					},
				},
				Error: "",
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			path:    "TEST3?withOrders=true&withOrderResume=false&withPrice=false",
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
					Orders: []presenter.OrderApiReturn{
						{
							Id:        "Order1",
							Quantity:  2,
							Price:     29.29,
							Currency:  "USD",
							OrderType: "buy",
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
							OrderType: "buy",
							Date:      dateFormatted,
							Brokerage: &presenter.Brokerage{
								Id:      "BrokerageID",
								Name:    "Test Broker",
								Country: "US",
							},
						},
					},
				},
				Error: "",
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			path:    "TEST3?withOrders=false&withOrderResume=true&withPrice=false",
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
					OrderInfos: &presenter.OrderInfos{
						TotalQuantity:        4,
						WeightedAdjPrice:     29.29,
						WeightedAveragePrice: 29.29,
					},
				},
				Error: "",
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			path:    "TEST3?withOrders=false&withOrderResume=false&withPrice=true",
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
					Price: &presenter.AssetPrice{
						OpenPrice:   200.19,
						ActualPrice: 199.98,
					},
				},
				Error: "",
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()
	logicApi := logicApi.NewMockApplication(*usecases)

	// Declare Sector Application Logic
	asset := AssetApi{
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
	api.Get("/asset/:symbol", asset.GetAsset)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "GET", "/api/asset/"+testCase.path,
			"application/json", testCase.idToken, nil)

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
		contentType  string
		bodyReq      bodyRequest
		expectedResp body
	}

	tests := []test{
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
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
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/json",
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
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/pdf",
			bodyReq: bodyRequest{
				AssetType: "ETF",
				Symbol:    "TEST11",
				Fullname:  "Test Company",
				Country:   "BR",
			},
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Asset:   nil,
				Error:   entity.ErrInvalidApiBody.Error(),
			},
		},
		{
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/json",
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
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/json",
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
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/json",
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
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/json",
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
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/json",
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
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/json",
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
	logicApi := logicApi.NewMockApplication(*usecases)

	// Declare Sector Application Logic
	asset := AssetApi{
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
		resp, _ := MockHttpRequest(app, "POST", "/api/asset", testCase.contentType,
			testCase.idToken, testCase.bodyReq)

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
	logicApi := logicApi.NewMockApplication(*usecases)

	// Declare Sector Application Logic
	asset := AssetApi{
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
			"application/json", testCase.idToken, nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}
}

func TestApiAssetLookup(t *testing.T) {
	type body struct {
		Success      bool                 `json:"success"`
		Message      string               `json:"message"`
		Error        string               `json:"error"`
		Code         int                  `json:"code"`
		SymbolLookup *entity.SymbolLookup `json:"symbolLookup"`
	}

	type test struct {
		idToken      string
		symbol       string
		country      string
		expectedResp body
	}

	tests := []test{
		{
			idToken: "ValidIdTokenWithoutEmailVerification",
			symbol:  "TEST3",
			country: "BR",
			expectedResp: body{
				Code:         401,
				Success:      false,
				Message:      entity.ErrMessageApiAuthentication.Error(),
				SymbolLookup: nil,
				Error:        "",
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			symbol:  "UNKNOWN_SYMBOL",
			country: "BR",
			expectedResp: body{
				Code:         404,
				Success:      false,
				Message:      entity.ErrInvalidAssetSymbol.Error(),
				SymbolLookup: nil,
				Error:        "",
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			symbol:  "",
			country: "BR",
			expectedResp: body{
				Code:         400,
				Success:      false,
				Message:      entity.ErrMessageApiRequest.Error(),
				SymbolLookup: nil,
				Error:        entity.ErrInvalidApiQuerySymbolBlank.Error(),
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			symbol:  "TEST3",
			country: "ERROR",
			expectedResp: body{
				Code:         400,
				Success:      false,
				Message:      entity.ErrMessageApiRequest.Error(),
				SymbolLookup: nil,
				Error:        entity.ErrInvalidCountryCode.Error(),
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			symbol:  "TEST3",
			country: "BR",
			expectedResp: body{
				Code:    200,
				Success: true,
				Message: "Symbol Lookup returned successfully",
				SymbolLookup: &entity.SymbolLookup{
					Fullname: "Test Name",
					Symbol:   "TEST3",
					Type:     "ETP",
				},
				Error: "",
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			symbol:  "TEST3",
			country: "US",
			expectedResp: body{
				Code:    200,
				Success: true,
				Message: "Symbol Lookup returned successfully",
				SymbolLookup: &entity.SymbolLookup{
					Fullname: "Test Name",
					Symbol:   "TEST3",
					Type:     "ETP",
				},
				Error: "",
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()
	logicApi := logicApi.NewMockApplication(*usecases)

	// Declare Sector Application Logic
	asset := AssetApi{
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
	api.Get("/asset-lookup", asset.GetSymbolLookup)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "GET", "/api/asset-lookup?symbol="+
			testCase.symbol+"&country="+testCase.country, "application/json",
			testCase.idToken, nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}
}

func TestApiAssetGetPrice(t *testing.T) {
	type body struct {
		Success     bool                `json:"success"`
		Message     string              `json:"message"`
		Error       string              `json:"error"`
		Code        int                 `json:"code"`
		SymbolPrice *entity.SymbolPrice `json:"symbolPrice"`
	}

	type test struct {
		idToken      string
		symbol       string
		country      string
		expectedResp body
	}

	tests := []test{
		{
			idToken: "ValidIdTokenWithoutEmailVerification",
			symbol:  "TEST3",
			country: "BR",
			expectedResp: body{
				Code:        401,
				Success:     false,
				Message:     entity.ErrMessageApiAuthentication.Error(),
				SymbolPrice: nil,
				Error:       "",
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			symbol:  "UNKNOWN_SYMBOL",
			country: "BR",
			expectedResp: body{
				Code:        404,
				Success:     false,
				Message:     entity.ErrInvalidAssetSymbol.Error(),
				SymbolPrice: nil,
				Error:       "",
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			symbol:  "",
			country: "BR",
			expectedResp: body{
				Code:        400,
				Success:     false,
				Message:     entity.ErrMessageApiRequest.Error(),
				SymbolPrice: nil,
				Error:       entity.ErrInvalidApiQuerySymbolBlank.Error(),
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			symbol:  "TEST3",
			country: "ERROR",
			expectedResp: body{
				Code:        400,
				Success:     false,
				Message:     entity.ErrMessageApiRequest.Error(),
				SymbolPrice: nil,
				Error:       entity.ErrInvalidCountryCode.Error(),
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			symbol:  "TEST3",
			country: "BR",
			expectedResp: body{
				Code:    200,
				Success: true,
				Message: "Symbol Price returned successfully",
				SymbolPrice: &entity.SymbolPrice{
					Symbol:         "TEST3",
					CurrentPrice:   29.29,
					LowPrice:       28.00,
					HighPrice:      29.89,
					OpenPrice:      29.29,
					PrevClosePrice: 29.29,
					MarketCap:      1018388,
				},
				Error: "",
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()
	logicApi := logicApi.NewMockApplication(*usecases)

	// Declare Sector Application Logic
	asset := AssetApi{
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
	api.Get("/asset-price", asset.GetSymbolPrice)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "GET", "/api/asset-price?symbol="+
			testCase.symbol+"&country="+testCase.country, "application/json",
			testCase.idToken, nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}
}
