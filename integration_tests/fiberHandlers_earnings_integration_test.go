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
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

func configureEarningsApp(dbpool *pgx.Conn) (fiberHandlers.EarningsApi,
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

	earnings := fiberHandlers.EarningsApi{
		ApplicationLogic: *applicationLogics,
		ApiLogic:         logicApiUseCases,
	}

	return earnings, *applicationLogics
}

func TestFiberHandlersIntegrationTestGetEarningsFromAssetUser(t *testing.T) {
	type body struct {
		Success bool                          `json:"success"`
		Message string                        `json:"message"`
		Error   string                        `json:"error"`
		Code    int                           `json:"code"`
		Earning []presenter.EarningsApiReturn `json:"earning"`
	}

	type test struct {
		idToken                   string
		symbol                    string
		orderBy                   string
		limit                     string
		offset                    string
		expectedResponse          body
		expectedOrderLengthReturn int
	}

	tests := []test{
		{
			idToken: "INVALID_ID_TOKEN",
			expectedResponse: body{
				Code:    401,
				Success: false,
				Message: entity.ErrMessageApiAuthentication.Error(),
				Error:   "",
				Earning: nil,
			},
		},
		{
			idToken: "TestNoAdminID",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidApiQuerySymbolBlank.Error(),
				Earning: nil,
			},
		},
		{
			idToken: "TestAdminID",
			symbol:  "ERRO3",
			expectedResponse: body{
				Code:    404,
				Success: false,
				Message: entity.ErrMessageApiAssetSymbolUser.Error(),
				Error:   "",
				Earning: nil,
			},
		},
		{
			idToken: "TestAdminID",
			symbol:  "EGIE3",
			limit:   "2b",
			offset:  "0",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidEarningsLimit.Error(),
				Earning: nil,
			},
			expectedOrderLengthReturn: 0,
		},
		{
			idToken: "TestAdminID",
			symbol:  "EGIE3",
			limit:   "2",
			offset:  "0B",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidEarningsOffset.Error(),
				Earning: nil,
			},
			expectedOrderLengthReturn: 0,
		},
		{
			idToken: "TestAdminID",
			symbol:  "EGIE3",
			limit:   "2",
			offset:  "0",
			orderBy: "error",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidEarningsOrderBy.Error(),
				Earning: nil,
			},
			expectedOrderLengthReturn: 0,
		},
		{
			idToken: "TestAdminID",
			symbol:  "EGIE3",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Earnings returned successfully",
				Error:   "",
				Earning: nil,
			},
			expectedOrderLengthReturn: 4,
		},
		{
			idToken: "TestAdminID",
			symbol:  "EGIE3",
			limit:   "2",
			offset:  "0",
			orderBy: "desc",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Earnings returned successfully",
				Error:   "",
				Earning: nil,
			},
			expectedOrderLengthReturn: 2,
		},
		{
			idToken: "TestAdminID",
			symbol:  "EGIE3",
			limit:   "2",
			offset:  "0",
			orderBy: "asc",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Earnings returned successfully",
				Error:   "",
				Earning: nil,
			},
			expectedOrderLengthReturn: 2,
		},
		{
			idToken: "TestAdminID",
			symbol:  "EGIE3",
			limit:   "3",
			offset:  "5",
			orderBy: "asc",
			expectedResponse: body{
				Code:    404,
				Success: false,
				Message: entity.ErrMessageApiEarningAssetUser.Error(),
				Error:   "",
				Earning: nil,
			},
			expectedOrderLengthReturn: 0,
		},
	}

	DBpool := connectDatabase()

	earnings, applicationsLogics := configureEarningsApp(DBpool)

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
	api.Get("/earnings", earnings.GetEarningsFromAssetUser)

	for _, testCase := range tests {
		bodyResponse := body{}

		resp, _ := fiberHandlers.MockHttpRequest(app, "GET", "/api/earnings?symbol="+
			testCase.symbol+"&limit="+testCase.limit+"&offset="+testCase.offset+
			"&orderBy="+testCase.orderBy, "", testCase.idToken, nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &bodyResponse)
		bodyResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResponse.Code, bodyResponse.Code)
		assert.Equal(t, testCase.expectedResponse.Success, bodyResponse.Success)
		assert.Equal(t, testCase.expectedResponse.Message, bodyResponse.Message)
		assert.Equal(t, testCase.expectedResponse.Error, bodyResponse.Error)

		if testCase.expectedResponse.Earning != nil {
			assert.Equal(t, testCase.expectedOrderLengthReturn,
				len(bodyResponse.Earning))
		} else {
			assert.Nil(t, testCase.expectedResponse.Earning)
		}

	}
}

func TestFiberHandlersIntegrationTestCreateEarnings(t *testing.T) {
	type body struct {
		Success bool                         `json:"success"`
		Message string                       `json:"message"`
		Error   string                       `json:"error"`
		Code    int                          `json:"code"`
		Earning *presenter.EarningsApiReturn `json:"earning"`
	}

	type test struct {
		idToken          string
		contentType      string
		symbol           string
		amount           float64
		currency         string
		earningType      string
		date             string
		expectedResponse body
	}

	dateString := "2021-10-05"
	layout := "2006-01-02"
	dateFormatted, _ := time.Parse(layout, dateString)
	tests := []test{
		{
			idToken:     "INVALID_ID_TOKEN",
			contentType: "application/json",
			expectedResponse: body{
				Code:    401,
				Success: false,
				Message: entity.ErrMessageApiAuthentication.Error(),
				Error:   "",
				Earning: nil,
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
				Earning: nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidEarningsCreateBlankFields.Error(),
				Earning: nil,
			},
		},
		{
			idToken:     "TestNoAdminID",
			contentType: "application/json",
			amount:      -2.19,
			symbol:      "EGIE3",
			earningType: "JCP",
			date:        "2021-10-05",
			currency:    "BRL",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidEarningsAmount.Error(),
				Earning: nil,
			},
		},
		{
			idToken:     "TestNoAdminID",
			contentType: "application/json",
			amount:      2.19,
			symbol:      "EGIE3",
			earningType: "error",
			date:        "2021-10-05",
			currency:    "BRL",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidEarningType.Error(),
				Earning: nil,
			},
		},
		{
			idToken:     "TestNoAdminID",
			contentType: "application/json",
			amount:      2.19,
			symbol:      "ERRO3",
			earningType: "JCP",
			date:        "2021-10-05",
			currency:    "BRL",
			expectedResponse: body{
				Code:    404,
				Success: false,
				Message: entity.ErrMessageApiAssetSymbolUser.Error(),
				Error:   "",
				Earning: nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			amount:      2.19,
			symbol:      "EGIE3",
			earningType: "Dividendos",
			date:        "2021-10-05",
			currency:    "BRL",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Earning registered successfully",
				Error:   "",
				Earning: &presenter.EarningsApiReturn{
					Type:     "Dividendos",
					Earning:  2.19,
					Currency: "BRL",
					Date:     &dateFormatted,
					Asset: &presenter.AssetApiReturn{
						Symbol: "EGIE3",
					},
				},
			},
		},
	}

	DBpool := connectDatabase()

	earnings, applicationsLogics := configureEarningsApp(DBpool)

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
	api.Post("/earnings", earnings.CreateEarnings)

	for _, testCase := range tests {
		bodyResponse := body{}

		bodyRequest := presenter.EarningsBody{
			Symbol:      testCase.symbol,
			Amount:      testCase.amount,
			Currency:    testCase.currency,
			EarningType: testCase.earningType,
			Date:        testCase.date,
		}

		resp, _ := fiberHandlers.MockHttpRequest(app, "POST", "/api/earnings",
			testCase.contentType, testCase.idToken, bodyRequest)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &bodyResponse)
		bodyResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResponse.Code, bodyResponse.Code)
		assert.Equal(t, testCase.expectedResponse.Success, bodyResponse.Success)
		assert.Equal(t, testCase.expectedResponse.Message, bodyResponse.Message)
		assert.Equal(t, testCase.expectedResponse.Error, bodyResponse.Error)

		if testCase.expectedResponse.Earning != nil {
			assert.Equal(t, testCase.expectedResponse.Earning.Earning,
				bodyResponse.Earning.Earning)
			assert.Equal(t, testCase.expectedResponse.Earning.Date,
				bodyResponse.Earning.Date)
			assert.Equal(t, testCase.expectedResponse.Earning.Currency,
				bodyResponse.Earning.Currency)
			assert.Equal(t, testCase.expectedResponse.Earning.Type,
				bodyResponse.Earning.Type)
			assert.Equal(t, testCase.expectedResponse.Earning.Asset.Symbol,
				bodyResponse.Earning.Asset.Symbol)
		} else {
			assert.Nil(t, testCase.expectedResponse.Earning)
		}

	}
}

func TestFiberHandlersIntegrationTestDeleteEarningFromUser(t *testing.T) {
	type body struct {
		Success bool                         `json:"success"`
		Message string                       `json:"message"`
		Error   string                       `json:"error"`
		Code    int                          `json:"code"`
		Earning *presenter.EarningsApiReturn `json:"earning"`
	}

	type bodyGetEarnings struct {
		Success    bool                          `json:"success"`
		Message    string                        `json:"message"`
		Error      string                        `json:"error"`
		Code       int                           `json:"code"`
		OrdersInfo []presenter.EarningsApiReturn `json:"earning"`
	}

	type test struct {
		idToken          string
		orderId          string
		expectedResponse body
	}

	DBpool := connectDatabase()

	earnings, applicationsLogics := configureEarningsApp(DBpool)

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
	api.Delete("/earnings/:id", earnings.DeleteEarningFromUser)
	api.Get("/earnings", earnings.GetEarningsFromAssetUser)

	resp, _ := fiberHandlers.MockHttpRequest(app, "GET", "/api/earnings?symbol="+
		"EGIE3&limit=1&offset=0&orderBy=DESC", "", "TestAdminID", nil)
	bodyResponseBytes, _ := ioutil.ReadAll(resp.Body)

	bodyResponse := bodyGetEarnings{}
	json.Unmarshal(bodyResponseBytes, &bodyResponse)
	bodyResponse.Code = resp.StatusCode
	orderTestId := bodyResponse.OrdersInfo[0].Id

	invalidId := "3erroa"
	tests := []test{
		{
			idToken: "INVALID_ID_TOKEN",
			expectedResponse: body{
				Code:    401,
				Success: false,
				Message: entity.ErrMessageApiAuthentication.Error(),
				Error:   "",
				Earning: nil,
			},
		},
		{
			idToken: "TestAdminID",
			orderId: invalidId,
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error: "ERROR: invalid input syntax for type uuid: \"" +
					invalidId + "\" (SQLSTATE 22P02)",
				Earning: nil,
			},
		},
		{
			idToken: "TestAdminID",
			orderId: orderTestId,
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Earning deleted successfully",
				Error:   "",
				Earning: &presenter.EarningsApiReturn{
					Id: orderTestId,
				},
			},
		},
	}

	for _, testCase := range tests {
		bodyResponse := body{}

		resp, _ := fiberHandlers.MockHttpRequest(app, "DELETE", "/api/earnings/"+
			testCase.orderId, "", testCase.idToken, nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &bodyResponse)
		bodyResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResponse.Code, bodyResponse.Code)
		assert.Equal(t, testCase.expectedResponse.Success, bodyResponse.Success)
		assert.Equal(t, testCase.expectedResponse.Message, bodyResponse.Message)
		assert.Equal(t, testCase.expectedResponse.Error, bodyResponse.Error)

		if testCase.expectedResponse.Earning != nil {
			assert.Equal(t, testCase.expectedResponse.Earning.Id,
				bodyResponse.Earning.Id)
		} else {
			assert.Nil(t, testCase.expectedResponse.Earning)
		}

	}
}

func TestFiberHandlersIntegrationTesUpdateEarningFromUser(t *testing.T) {
	type body struct {
		Success bool                         `json:"success"`
		Message string                       `json:"message"`
		Error   string                       `json:"error"`
		Code    int                          `json:"code"`
		Earning *presenter.EarningsApiReturn `json:"earning"`
	}

	type bodyGetEarnings struct {
		Success    bool                          `json:"success"`
		Message    string                        `json:"message"`
		Error      string                        `json:"error"`
		Code       int                           `json:"code"`
		OrdersInfo []presenter.EarningsApiReturn `json:"earning"`
	}

	type test struct {
		idToken          string
		contentType      string
		orderId          string
		symbol           string
		amount           float64
		earningType      string
		date             string
		expectedResponse body
	}

	DBpool := connectDatabase()

	earnings, applicationsLogics := configureEarningsApp(DBpool)

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
	api.Put("/earnings/:id", earnings.UpdateEarningFromUser)
	api.Get("/earnings", earnings.GetEarningsFromAssetUser)

	resp, _ := fiberHandlers.MockHttpRequest(app, "GET", "/api/earnings?symbol="+
		"EGIE3&limit=1&offset=0&orderBy=DESC", "", "TestAdminID", nil)
	bodyResponseBytes, _ := ioutil.ReadAll(resp.Body)

	bodyResponse := bodyGetEarnings{}
	json.Unmarshal(bodyResponseBytes, &bodyResponse)
	bodyResponse.Code = resp.StatusCode
	orderTestId := bodyResponse.OrdersInfo[0].Id

	invalidId := "3erroa"
	dateString := "2021-10-07"
	layout := "2006-01-02"
	dateFormatted, _ := time.Parse(layout, dateString)
	tests := []test{
		{
			idToken: "INVALID_ID_TOKEN",
			expectedResponse: body{
				Code:    401,
				Success: false,
				Message: entity.ErrMessageApiAuthentication.Error(),
				Error:   "",
				Earning: nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/pdf",
			orderId:     invalidId,
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidApiBody.Error(),
				Earning: nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			orderId:     invalidId,
			expectedResponse: body{
				Code:    404,
				Success: false,
				Message: entity.ErrMessageApiEarningId.Error(),
				Error:   "",
				Earning: nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			orderId:     orderTestId,
			amount:      -2.19,
			earningType: "JCP",
			date:        "2021-10-05",
			symbol:      "EGIE3",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidEarningsAmount.Error(),
				Earning: nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			orderId:     orderTestId,
			amount:      2.19,
			earningType: "",
			date:        "2021-10-05",
			symbol:      "EGIE3",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidEarningsCreateBlankFields.Error(),
				Earning: nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			orderId:     orderTestId,
			amount:      2.19,
			earningType: "ERROR",
			date:        "2021-10-05",
			symbol:      "EGIE3",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidEarningType.Error(),
				Earning: nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			orderId:     orderTestId,
			amount:      2.55,
			earningType: "JCP",
			date:        dateString,
			symbol:      "EGIE3",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Earning updated successfully",
				Error:   "",
				Earning: &presenter.EarningsApiReturn{
					Id:       orderTestId,
					Type:     "JCP",
					Earning:  2.55,
					Date:     &dateFormatted,
					Currency: "BRL",
					Asset: &presenter.AssetApiReturn{
						Symbol: "EGIE3",
					},
				},
			},
		},
	}

	for _, testCase := range tests {
		bodyResponse := body{}

		bodyRequest := presenter.EarningsBody{
			Symbol:      testCase.symbol,
			Amount:      testCase.amount,
			EarningType: testCase.earningType,
			Date:        testCase.date,
		}

		resp, _ := fiberHandlers.MockHttpRequest(app, "PUT", "/api/earnings/"+
			testCase.orderId, testCase.contentType, testCase.idToken, bodyRequest)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &bodyResponse)
		bodyResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResponse.Code, bodyResponse.Code)
		assert.Equal(t, testCase.expectedResponse.Success, bodyResponse.Success)
		assert.Equal(t, testCase.expectedResponse.Message, bodyResponse.Message)
		assert.Equal(t, testCase.expectedResponse.Error, bodyResponse.Error)

		if testCase.expectedResponse.Earning != nil {
			assert.Equal(t, testCase.expectedResponse.Earning.Id,
				bodyResponse.Earning.Id)
			assert.Equal(t, testCase.expectedResponse.Earning.Type,
				bodyResponse.Earning.Type)
			assert.Equal(t, testCase.expectedResponse.Earning.Earning,
				bodyResponse.Earning.Earning)
			assert.Equal(t, testCase.expectedResponse.Earning.Date,
				bodyResponse.Earning.Date)
			assert.Equal(t, testCase.expectedResponse.Earning.Currency,
				bodyResponse.Earning.Currency)
			assert.Equal(t, testCase.expectedResponse.Earning.Asset.Symbol,
				bodyResponse.Earning.Asset.Symbol)
		} else {
			assert.Nil(t, testCase.expectedResponse.Earning)
		}

	}
}
