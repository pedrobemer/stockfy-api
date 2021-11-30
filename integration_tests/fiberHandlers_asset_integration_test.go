package integration_tests

import (
	"encoding/json"
	"io/ioutil"
	"stockfyApi/api/handlers/fiberHandlers"
	"stockfyApi/api/middleware"
	"stockfyApi/api/presenter"
	"stockfyApi/database/postgresql"
	"stockfyApi/entity"
	externalapi "stockfyApi/externalApi"
	"stockfyApi/externalApi/alphaVantage"
	"stockfyApi/externalApi/finnhub"
	"stockfyApi/usecases"
	"stockfyApi/usecases/logicApi"
	"stockfyApi/usecases/user"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

func configureAssetApp(dbpool *pgx.Conn) (fiberHandlers.AssetApi,
	usecases.Applications) {

	dbInterfaces := postgresql.NewPostgresInstance(dbpool)
	firebaseInterface := user.NewExternalApi()

	applicationLogics := usecases.NewApplications(dbInterfaces,
		firebaseInterface)

	mockFinnhubClient := finnhub.MockClient{
		Client: fiberHandlers.MockClient{
			DoFunc: mockDoFuncFinnhubVerifySymbol,
		},
	}

	mockAlphaClient := alphaVantage.MockClient{
		Client: fiberHandlers.MockClient{
			DoFunc: mockDoFuncAlphaVerifySymbol,
		},
	}

	finnhubInterface := finnhub.NewFinnhubApi("Test",
		mockFinnhubClient.HttpOutsideClientRequest)
	alphaInterface := alphaVantage.NewAlphaVantageApi("Test",
		mockAlphaClient.HttpOutsideClientRequest)

	externalInterface := externalapi.ThirdPartyInterfaces{
		FinnhubApi:      finnhubInterface,
		AlphaVantageApi: alphaInterface,
	}
	logicApiUseCases := logicApi.NewApplication(*applicationLogics,
		externalInterface)

	assets := fiberHandlers.AssetApi{
		ApplicationLogic: *applicationLogics,
		LogicApi:         logicApiUseCases,
	}

	return assets, *applicationLogics
}

func TestFiberHandlersIntegrationTestCreateAsset(t *testing.T) {

	type body struct {
		Success bool                      `json:"success"`
		Message string                    `json:"message"`
		Error   string                    `json:"error"`
		Code    int                       `json:"code"`
		Asset   *presenter.AssetApiReturn `json:"asset"`
	}

	type test struct {
		idToken          string
		contentType      string
		assetType        string
		symbol           string
		fullname         string
		country          string
		expectedResponse body
	}

	tests := []test{
		{
			idToken:     "INVALID_ID_TOKEN",
			contentType: "application/json",
			expectedResponse: body{
				Code:    401,
				Success: false,
				Message: entity.ErrMessageApiAuthentication.Error(),
				Error:   "",
				Asset:   nil,
			},
		},
		{
			idToken:     "TestNoAdminID",
			contentType: "application/json",
			expectedResponse: body{
				Code:    403,
				Success: false,
				Message: entity.ErrMessageApiAuthorization.Error(),
				Error:   entity.ErrInvalidUserAdminPrivilege.Error(),
				Asset:   nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/pdf",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidApiBody.Error(),
				Asset:   nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			symbol:      "ITUB4",
			assetType:   "STOCK",
			country:     "BR",
			fullname:    "Itaú Unibanco Holding S.A",
			expectedResponse: body{
				Code:    403,
				Success: false,
				Message: entity.ErrMessageApiAuthorization.Error(),
				Error:   entity.ErrInvalidAssetSymbolExist.Error(),
				Asset:   nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			symbol:      "KNRI11",
			assetType:   "STOCK",
			country:     "ERROR",
			fullname:    "Itaú Unibanco Holding S.A",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidCountryCode.Error(),
				Asset:   nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			symbol:      "ERRO3",
			assetType:   "STOCK",
			country:     "BR",
			fullname:    "Itaú Unibanco Holding S.A",
			expectedResponse: body{
				Code:    404,
				Success: false,
				Message: entity.ErrMessageApiAssetSymbolUser.Error(),
				Error:   "",
				Asset:   nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			symbol:      "ERRO3",
			assetType:   "STOCK",
			country:     "US",
			fullname:    "Itaú Unibanco Holding S.A",
			expectedResponse: body{
				Code:    404,
				Success: false,
				Message: entity.ErrMessageApiAssetSymbolUser.Error(),
				Error:   "",
				Asset:   nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			symbol:      "ERRO3",
			assetType:   "STOCK",
			country:     "US",
			fullname:    "Itaú Unibanco Holding S.A",
			expectedResponse: body{
				Code:    404,
				Success: false,
				Message: entity.ErrMessageApiAssetSymbolUser.Error(),
				Error:   "",
				Asset:   nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			symbol:      "FLRY3",
			assetType:   "STOCK",
			country:     "BR",
			fullname:    "Fleury S.A",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Asset creation was sucessful",
				Error:   "",
				Asset: &presenter.AssetApiReturn{
					Preference: "ON",
					Fullname:   "Fleury S.A",
					Symbol:     "FLRY3",
				},
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			symbol:      "AMT",
			assetType:   "REIT",
			country:     "US",
			fullname:    "American Tower Corp",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Asset creation was sucessful",
				Error:   "",
				Asset: &presenter.AssetApiReturn{
					Preference: "",
					Fullname:   "American Tower Corp",
					Symbol:     "AMT",
				},
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			symbol:      "VTI",
			assetType:   "ETF",
			country:     "US",
			fullname:    "Vanguard Total Stock Mkt ETF",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Asset creation was sucessful",
				Error:   "",
				Asset: &presenter.AssetApiReturn{
					Preference: "",
					Fullname:   "Vanguard Total Stock Mkt ETF",
					Symbol:     "VTI",
				},
			},
		},
	}

	DBpool := connectDatabase()

	assets, applicationsLogics := configureAssetApp(DBpool)

	app := fiber.New()
	api := app.Group("/api")
	api.Use(middleware.NewFiberMiddleware(middleware.FiberMiddleware{
		UserAuthentication: applicationsLogics.UserApp,
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
	api.Post("/asset", assets.CreateAsset)

	for _, testCase := range tests {
		bodyResponse := body{}
		bodyRequestStruct := presenter.AssetBody{
			AssetType: testCase.assetType,
			Symbol:    testCase.symbol,
			Fullname:  testCase.fullname,
			Country:   testCase.country,
		}

		resp, _ := fiberHandlers.MockHttpRequest(app, "POST", "/api/asset",
			testCase.contentType, testCase.idToken, bodyRequestStruct)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &bodyResponse)
		bodyResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResponse.Code, bodyResponse.Code)
		assert.Equal(t, testCase.expectedResponse.Success, bodyResponse.Success)
		assert.Equal(t, testCase.expectedResponse.Message, bodyResponse.Message)
		assert.Equal(t, testCase.expectedResponse.Error, bodyResponse.Error)

		if testCase.expectedResponse.Asset != nil {
			assert.Equal(t, testCase.expectedResponse.Asset.Symbol,
				bodyResponse.Asset.Symbol)
			assert.Equal(t, testCase.expectedResponse.Asset.Preference,
				bodyResponse.Asset.Preference)
			assert.Equal(t, testCase.expectedResponse.Asset.Fullname,
				bodyResponse.Asset.Fullname)
		} else {
			assert.Nil(t, testCase.expectedResponse.Asset)
		}

	}

}

func TestFiberHandlersIntegrationTestGetAsset(t *testing.T) {
	type body struct {
		Success bool                      `json:"success"`
		Message string                    `json:"message"`
		Error   string                    `json:"error"`
		Code    int                       `json:"code"`
		Asset   *presenter.AssetApiReturn `json:"asset"`
	}

	type test struct {
		idToken              string
		symbol               string
		withOrdersQuery      string
		withOrderResumeQuery string
		expectedResponse     body
	}

	tests := []test{
		{
			idToken:              "INVALID_ID_TOKEN",
			symbol:               "FLRY3",
			withOrdersQuery:      "false",
			withOrderResumeQuery: "false",
			expectedResponse: body{
				Code:    401,
				Success: false,
				Message: entity.ErrMessageApiAuthentication.Error(),
				Error:   "",
				Asset:   nil,
			},
		},
		{
			idToken:              "TestAdminID",
			symbol:               "FLRY3",
			withOrdersQuery:      "false",
			withOrderResumeQuery: "error",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidApiQueryWithOrderResume.Error(),
				Asset:   nil,
			},
		},
		{
			idToken:              "TestAdminID",
			symbol:               "FLRY3",
			withOrdersQuery:      "error",
			withOrderResumeQuery: "false",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidApiQueryWithOrders.Error(),
				Asset:   nil,
			},
		},
		{
			idToken:              "TestAdminID",
			symbol:               "FLRY3",
			withOrdersQuery:      "false",
			withOrderResumeQuery: "false",
			expectedResponse: body{
				Code:    404,
				Success: false,
				Message: entity.ErrMessageApiAssetSymbolUser.Error(),
				Error:   "",
				Asset:   nil,
			},
		},
		{
			idToken:              "TestAdminID",
			symbol:               "ITUB4",
			withOrdersQuery:      "false",
			withOrderResumeQuery: "false",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Asset information returned successfully",
				Error:   "",
				Asset: &presenter.AssetApiReturn{
					Symbol:     "ITUB4",
					Preference: "PN",
					Fullname:   "Itau Unibanco Holding S.A",
					AssetType: &presenter.AssetType{
						Type:    "STOCK",
						Country: "BR",
						Name:    "Ações Brasil",
					},
					Sector: &presenter.Sector{
						Name: "Finances",
					},
				},
			},
		},
		{
			idToken:              "TestAdminID",
			symbol:               "ITUB4",
			withOrdersQuery:      "false",
			withOrderResumeQuery: "true",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Asset information returned successfully",
				Error:   "",
				Asset: &presenter.AssetApiReturn{
					Symbol:     "ITUB4",
					Preference: "PN",
					Fullname:   "Itau Unibanco Holding S.A",
					AssetType: &presenter.AssetType{
						Type:    "STOCK",
						Country: "BR",
						Name:    "Ações Brasil",
					},
					Sector: &presenter.Sector{
						Name: "Finances",
					},
					OrderInfos: &presenter.OrderInfos{
						TotalQuantity:        29,
						WeightedAveragePrice: 24.005483870967744,
						WeightedAdjPrice:     24.343793103448277,
					},
				},
			},
		},
		{
			idToken:              "TestAdminID",
			symbol:               "ITUB4",
			withOrdersQuery:      "true",
			withOrderResumeQuery: "false",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Asset information returned successfully",
				Error:   "",
				Asset: &presenter.AssetApiReturn{
					Symbol:     "ITUB4",
					Preference: "PN",
					Fullname:   "Itau Unibanco Holding S.A",
					AssetType: &presenter.AssetType{
						Type:    "STOCK",
						Country: "BR",
						Name:    "Ações Brasil",
					},
					Sector: &presenter.Sector{
						Name: "Finances",
					},
				},
			},
		},
		{
			idToken:              "TestAdminID",
			symbol:               "ITUB4",
			withOrdersQuery:      "true",
			withOrderResumeQuery: "true",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Asset information returned successfully",
				Error:   "",
				Asset: &presenter.AssetApiReturn{
					Symbol:     "ITUB4",
					Preference: "PN",
					Fullname:   "Itau Unibanco Holding S.A",
					AssetType: &presenter.AssetType{
						Type:    "STOCK",
						Country: "BR",
						Name:    "Ações Brasil",
					},
					Sector: &presenter.Sector{
						Name: "Finances",
					},
					OrderInfos: &presenter.OrderInfos{
						TotalQuantity:        29,
						WeightedAveragePrice: 24.005483870967744,
						WeightedAdjPrice:     24.343793103448277,
					},
				},
			},
		},
	}

	DBpool := connectDatabase()
	assets, applicationLogics := configureAssetApp(DBpool)

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
	api.Get("/asset/:symbol", assets.GetAsset)

	for _, testCase := range tests {
		bodyResponse := body{}

		resp, _ := fiberHandlers.MockHttpRequest(app, "GET", "/api/asset/"+
			testCase.symbol+"?withOrders="+testCase.withOrdersQuery+
			"&withOrderResume="+testCase.withOrderResumeQuery, "",
			testCase.idToken, nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &bodyResponse)
		bodyResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResponse.Code, bodyResponse.Code)
		assert.Equal(t, testCase.expectedResponse.Success, bodyResponse.Success)
		assert.Equal(t, testCase.expectedResponse.Message, bodyResponse.Message)
		assert.Equal(t, testCase.expectedResponse.Error, bodyResponse.Error)

		if testCase.withOrderResumeQuery == "true" {
			assert.Equal(t, testCase.expectedResponse.Asset.OrderInfos,
				bodyResponse.Asset.OrderInfos)
		}

		if testCase.withOrdersQuery == "true" {
			assert.NotNil(t, bodyResponse.Asset.Orders)
		}

		if testCase.expectedResponse.Asset != nil {
			assert.Equal(t, testCase.expectedResponse.Asset.Symbol,
				bodyResponse.Asset.Symbol)
			assert.Equal(t, testCase.expectedResponse.Asset.Preference,
				bodyResponse.Asset.Preference)
			assert.Equal(t, testCase.expectedResponse.Asset.Fullname,
				bodyResponse.Asset.Fullname)
			assert.Equal(t, testCase.expectedResponse.Asset.Sector.Name,
				bodyResponse.Asset.Sector.Name)
			assert.Equal(t, testCase.expectedResponse.Asset.AssetType.Type,
				bodyResponse.Asset.AssetType.Type)
			assert.Equal(t, testCase.expectedResponse.Asset.AssetType.Country,
				bodyResponse.Asset.AssetType.Country)
			assert.Equal(t, testCase.expectedResponse.Asset.AssetType.Name,
				bodyResponse.Asset.AssetType.Name)
		} else {
			assert.Nil(t, bodyResponse.Asset)
		}
	}
}

func TestFiberHandlersIntegrationTestDeleteAssetWithMyUserFalse(t *testing.T) {

	type body struct {
		Success bool                      `json:"success"`
		Message string                    `json:"message"`
		Error   string                    `json:"error"`
		Code    int                       `json:"code"`
		Asset   *presenter.AssetApiReturn `json:"asset"`
	}

	type test struct {
		idToken          string
		symbol           string
		myUserQuery      string
		expectedResponse body
	}

	tests := []test{
		{
			idToken:     "INVALID_ID_TOKEN",
			symbol:      "FLRY3",
			myUserQuery: "false",
			expectedResponse: body{
				Code:    401,
				Success: false,
				Message: entity.ErrMessageApiAuthentication.Error(),
				Error:   "",
				Asset:   nil,
			},
		},
		{
			idToken:     "TestNoAdminID",
			symbol:      "FLRY3",
			myUserQuery: "false",
			expectedResponse: body{
				Code:    403,
				Success: false,
				Message: entity.ErrMessageApiAuthorization.Error(),
				Error:   entity.ErrInvalidUserAdminPrivilege.Error(),
				Asset:   nil,
			},
		},
		{
			idToken:     "TestAdminID",
			symbol:      "ERRO3",
			myUserQuery: "false",
			expectedResponse: body{
				Code:    404,
				Success: false,
				Message: entity.ErrMessageApiAssetSymbolUser.Error(),
				Error:   "",
				Asset:   nil,
			},
		},
		{
			idToken:     "TestAdminID",
			symbol:      "FLRY3",
			myUserQuery: "false",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Asset was deleted successfuly",
				Error:   "",
				Asset: &presenter.AssetApiReturn{
					Symbol:     "FLRY3",
					Fullname:   "Fleury S.A",
					Preference: "ON",
					AssetType: &presenter.AssetType{
						Country: "BR",
						Type:    "STOCK",
						Name:    "Ações Brasil",
					},
					Sector: &presenter.Sector{
						Name: "Health Care",
					},
				},
			},
		},
		{
			idToken:     "TestAdminID",
			symbol:      "AMT",
			myUserQuery: "false",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Asset was deleted successfuly",
				Error:   "",
				Asset: &presenter.AssetApiReturn{
					Symbol:     "AMT",
					Fullname:   "American Tower Corp",
					Preference: "",
					AssetType: &presenter.AssetType{
						Country: "US",
						Type:    "REIT",
						Name:    "REITs",
					},
					Sector: &presenter.Sector{
						Name: "Real Estate",
					},
				},
			},
		},
		{
			idToken:     "TestAdminID",
			symbol:      "VTI",
			myUserQuery: "false",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Asset was deleted successfuly",
				Error:   "",
				Asset: &presenter.AssetApiReturn{
					Symbol:     "VTI",
					Fullname:   "Vanguard Total Stock Mkt ETF",
					Preference: "",
					AssetType: &presenter.AssetType{
						Country: "US",
						Type:    "ETF",
						Name:    "ETFs EUA",
					},
					Sector: &presenter.Sector{
						Name: "Blend",
					},
				},
			},
		},
	}

	DBpool := connectDatabase()

	assets, applicationLogics := configureAssetApp(DBpool)

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
	api.Delete("/asset/:symbol", assets.DeleteAsset)

	for _, testCase := range tests {
		bodyResponse := body{}

		resp, _ := fiberHandlers.MockHttpRequest(app, "DELETE", "/api/asset/"+
			testCase.symbol+"?myUser="+testCase.myUserQuery, "", testCase.idToken,
			nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &bodyResponse)
		bodyResponse.Code = resp.StatusCode
		if testCase.expectedResponse.Asset != nil {
			assert.Equal(t, testCase.expectedResponse.Asset.Symbol,
				bodyResponse.Asset.Symbol)
			assert.Equal(t, testCase.expectedResponse.Asset.Preference,
				bodyResponse.Asset.Preference)
			assert.Equal(t, testCase.expectedResponse.Asset.Fullname,
				bodyResponse.Asset.Fullname)
			assert.Equal(t, testCase.expectedResponse.Asset.Sector.Name,
				bodyResponse.Asset.Sector.Name)
			assert.Equal(t, testCase.expectedResponse.Asset.AssetType.Type,
				bodyResponse.Asset.AssetType.Type)
			assert.Equal(t, testCase.expectedResponse.Asset.AssetType.Country,
				bodyResponse.Asset.AssetType.Country)
			assert.Equal(t, testCase.expectedResponse.Asset.AssetType.Name,
				bodyResponse.Asset.AssetType.Name)
		} else {
			assert.Nil(t, bodyResponse.Asset)
		}

	}

}

func TestFiberHandlersIntegrationTestDeleteAssetWithMyUserTrue(t *testing.T) {

	type body struct {
		Success bool                      `json:"success"`
		Message string                    `json:"message"`
		Error   string                    `json:"error"`
		Code    int                       `json:"code"`
		Asset   *presenter.AssetApiReturn `json:"asset"`
	}

	type test struct {
		idToken          string
		symbol           string
		myUserQuery      string
		expectedResponse body
	}

	tests := []test{
		{
			idToken:     "TestAdminID",
			symbol:      "ERRO3",
			myUserQuery: "true",
			expectedResponse: body{
				Code:    404,
				Success: false,
				Message: entity.ErrMessageApiAssetSymbolUser.Error(),
				Error:   "",
				Asset:   nil,
			},
		},
		{
			idToken:     "TestAdminID",
			symbol:      "ITUB4",
			myUserQuery: "true",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Asset was deleted successfuly",
				Error:   "",
				Asset: &presenter.AssetApiReturn{
					Symbol:     "ITUB4",
					Fullname:   "Itau Unibanco Holding S.A",
					Preference: "PN",
					AssetType: &presenter.AssetType{
						Country: "BR",
						Type:    "STOCK",
						Name:    "Ações Brasil",
					},
					Sector: &presenter.Sector{
						Name: "Finances",
					},
				},
			},
		},
	}

	DBpool := connectDatabase()

	assets, applicationLogics := configureAssetApp(DBpool)

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
	api.Delete("/asset/:symbol", assets.DeleteAsset)

	for _, testCase := range tests {
		bodyResponse := body{}

		resp, _ := fiberHandlers.MockHttpRequest(app, "DELETE", "/api/asset/"+
			testCase.symbol+"?myUser="+testCase.myUserQuery, "", testCase.idToken,
			nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &bodyResponse)
		bodyResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResponse.Code, bodyResponse.Code)
		assert.Equal(t, testCase.expectedResponse.Success, bodyResponse.Success)
		assert.Equal(t, testCase.expectedResponse.Message, bodyResponse.Message)
		assert.Equal(t, testCase.expectedResponse.Error, bodyResponse.Error)

		if testCase.expectedResponse.Asset != nil {
			assert.Equal(t, testCase.expectedResponse.Asset.Symbol,
				bodyResponse.Asset.Symbol)
			assert.Equal(t, testCase.expectedResponse.Asset.Preference,
				bodyResponse.Asset.Preference)
			assert.Equal(t, testCase.expectedResponse.Asset.Fullname,
				bodyResponse.Asset.Fullname)
			assert.Equal(t, testCase.expectedResponse.Asset.Sector.Name,
				bodyResponse.Asset.Sector.Name)
			assert.Equal(t, testCase.expectedResponse.Asset.AssetType.Type,
				bodyResponse.Asset.AssetType.Type)
			assert.Equal(t, testCase.expectedResponse.Asset.AssetType.Country,
				bodyResponse.Asset.AssetType.Country)
			assert.Equal(t, testCase.expectedResponse.Asset.AssetType.Name,
				bodyResponse.Asset.AssetType.Name)
		} else {
			assert.Nil(t, bodyResponse.Asset)
		}
	}
}
