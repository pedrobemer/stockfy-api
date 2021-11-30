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

func configureOrdersApp(dbpool *pgx.Conn) (fiberHandlers.OrderApi,
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

	orders := fiberHandlers.OrderApi{
		ApplicationLogic:   *applicationLogics,
		ExternalInterfaces: externalInterface,
		LogicApi:           logicApiUseCases,
	}

	return orders, *applicationLogics
}

func TestFiberHandlersIntegrationTestCreateOrder(t *testing.T) {

	type body struct {
		Success bool                      `json:"success"`
		Message string                    `json:"message"`
		Error   string                    `json:"error"`
		Code    int                       `json:"code"`
		Orders  *presenter.OrderApiReturn `json:"orders"`
	}

	type test struct {
		idToken          string
		contentType      string
		symbol           string
		fullname         string
		brokerage        string
		quantity         float64
		price            float64
		currency         string
		orderType        string
		date             string
		assetType        string
		country          string
		expectedResponse body
	}

	dateString := "2021-10-02"
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
				Orders:  nil,
			},
		},
		{
			idToken:     "TestNoAdminID",
			contentType: "application/pdf",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidApiBody.Error(),
				Orders:  nil,
			},
		},
		{
			idToken:     "TestNoAdminID",
			contentType: "application/json",
			symbol:      "ITUB4",
			fullname:    "Itau Unibanco Holding S.A",
			brokerage:   "Clear",
			price:       20.00,
			quantity:    5,
			currency:    "BRL",
			orderType:   "error",
			date:        dateString,
			assetType:   "STOCK",
			country:     "BR",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidOrderType.Error(),
				Orders:  nil,
			},
		},
		{
			idToken:     "TestNoAdminID",
			contentType: "application/json",
			symbol:      "ITUB4",
			fullname:    "Itau Unibanco Holding S.A",
			brokerage:   "Clear",
			price:       20.00,
			quantity:    5,
			currency:    "BRL",
			orderType:   "buy",
			date:        dateString,
			assetType:   "STOCK",
			country:     "ERR",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidCountryCode.Error(),
				Orders:  nil,
			},
		},
		{
			idToken:     "TestNoAdminID",
			contentType: "application/json",
			symbol:      "ITUB4",
			fullname:    "Itau Unibanco Holding S.A",
			brokerage:   "Clear",
			price:       20.00,
			quantity:    5.1,
			currency:    "BRL",
			orderType:   "buy",
			date:        dateString,
			assetType:   "STOCK",
			country:     "BR",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidOrderQuantityBrazil.Error(),
				Orders:  nil,
			},
		},
		{
			idToken:     "TestNoAdminID",
			contentType: "application/json",
			symbol:      "ITUB4",
			fullname:    "Itau Unibanco Holding S.A",
			brokerage:   "Clear",
			price:       20.00,
			quantity:    5,
			currency:    "USD",
			orderType:   "buy",
			date:        dateString,
			assetType:   "STOCK",
			country:     "BR",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidBrazilCurrency.Error(),
				Orders:  nil,
			},
		},
		{
			idToken:     "TestNoAdminID",
			contentType: "application/json",
			symbol:      "ITUB4",
			fullname:    "Itau Unibanco Holding S.A",
			brokerage:   "Clear",
			price:       20.00,
			quantity:    5,
			currency:    "BRL",
			orderType:   "buy",
			date:        dateString,
			assetType:   "STOCK",
			country:     "US",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidUsaCurrency.Error(),
				Orders:  nil,
			},
		},
		{
			idToken:     "TestNoAdminID",
			contentType: "application/json",
			symbol:      "ITUB4",
			fullname:    "Itau Unibanco Holding S.A",
			brokerage:   "Clear",
			price:       20.00,
			quantity:    -5,
			currency:    "BRL",
			orderType:   "buy",
			date:        dateString,
			assetType:   "STOCK",
			country:     "BR",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidOrderBuyQuantity.Error(),
				Orders:  nil,
			},
		},
		{
			idToken:     "TestNoAdminID",
			contentType: "application/json",
			symbol:      "ITUB4",
			fullname:    "Itau Unibanco Holding S.A",
			brokerage:   "Clear",
			price:       20.00,
			quantity:    5,
			currency:    "BRL",
			orderType:   "sell",
			date:        dateString,
			assetType:   "STOCK",
			country:     "BR",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidOrderSellQuantity.Error(),
				Orders:  nil,
			},
		},
		{
			idToken:     "TestNoAdminID",
			contentType: "application/json",
			symbol:      "ITUB4",
			fullname:    "Itau Unibanco Holding S.A",
			brokerage:   "Error",
			price:       20.00,
			quantity:    5,
			currency:    "BRL",
			orderType:   "buy",
			date:        dateString,
			assetType:   "STOCK",
			country:     "BR",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidBrokerageNameSearch.Error(),
				Orders:  nil,
			},
		},
		{
			idToken:     "TestNoAdminID",
			contentType: "application/json",
			symbol:      "ERRO11",
			fullname:    "Itau Unibanco Holding S.A",
			brokerage:   "Error",
			price:       20.00,
			quantity:    5,
			currency:    "BRL",
			orderType:   "buy",
			date:        dateString,
			assetType:   "STOCK",
			country:     "BR",
			expectedResponse: body{
				Code:    404,
				Success: false,
				Message: entity.ErrMessageApiAssetSymbolUser.Error(),
				Error:   "",
				Orders:  nil,
			},
		},
		{
			idToken:     "TestNoAdminID",
			contentType: "application/json",
			symbol:      "ITUB4",
			fullname:    "Itau Unibanco Holding S.A",
			brokerage:   "Clear",
			price:       20.00,
			quantity:    5,
			currency:    "BRL",
			orderType:   "buy",
			date:        dateString,
			assetType:   "STOCK",
			country:     "BR",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Order registered successfully",
				Error:   "",
				Orders: &presenter.OrderApiReturn{
					Quantity:  5,
					Price:     20.00,
					Currency:  "BRL",
					OrderType: "buy",
					Date:      dateFormatted,
					Brokerage: &presenter.Brokerage{
						Name:    "Clear",
						Country: "BR",
					},
					Asset: &presenter.AssetApiReturn{
						Symbol:   "ITUB4",
						Fullname: "Itau Unibanco Holding S.A",
					},
				},
			},
		},
		{
			idToken:     "TestNoAdminID",
			contentType: "application/json",
			symbol:      "KNRI11",
			fullname:    "Kinea Renda Imobiliária Fundo de Investimento Imobiliário",
			brokerage:   "Clear",
			price:       129.98,
			quantity:    4,
			currency:    "BRL",
			orderType:   "buy",
			date:        dateString,
			assetType:   "STOCK",
			country:     "BR",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Order registered successfully",
				Error:   "",
				Orders: &presenter.OrderApiReturn{
					Quantity:  4,
					Price:     129.98,
					Currency:  "BRL",
					OrderType: "buy",
					Date:      dateFormatted,
					Brokerage: &presenter.Brokerage{
						Name:    "Clear",
						Country: "BR",
					},
					Asset: &presenter.AssetApiReturn{
						Symbol:   "KNRI11",
						Fullname: "Kinea Renda Imobiliária Fundo De Investimento Imobiliário",
					},
				},
			},
		},
	}

	DBpool := connectDatabase()

	order, applicationsLogics := configureOrdersApp(DBpool)

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
	api.Post("/orders", order.CreateUserOrder)

	for _, testCase := range tests {
		bodyResponse := body{}
		bodyRequestStruct := presenter.OrderBody{
			Symbol:    testCase.symbol,
			Fullname:  testCase.fullname,
			Brokerage: testCase.brokerage,
			Quantity:  testCase.quantity,
			Price:     testCase.price,
			Currency:  testCase.currency,
			OrderType: testCase.orderType,
			Date:      testCase.date,
			AssetType: testCase.assetType,
			Country:   testCase.country,
		}

		resp, _ := fiberHandlers.MockHttpRequest(app, "POST", "/api/orders",
			testCase.contentType, testCase.idToken, bodyRequestStruct)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &bodyResponse)
		bodyResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResponse.Code, bodyResponse.Code)
		assert.Equal(t, testCase.expectedResponse.Success, bodyResponse.Success)
		assert.Equal(t, testCase.expectedResponse.Message, bodyResponse.Message)
		assert.Equal(t, testCase.expectedResponse.Error, bodyResponse.Error)

		if testCase.expectedResponse.Orders != nil {
			assert.Equal(t, testCase.expectedResponse.Orders.Quantity,
				bodyResponse.Orders.Quantity)
			assert.Equal(t, testCase.expectedResponse.Orders.Price,
				bodyResponse.Orders.Price)
			assert.Equal(t, testCase.expectedResponse.Orders.Currency,
				bodyResponse.Orders.Currency)
			assert.Equal(t, testCase.expectedResponse.Orders.OrderType,
				bodyResponse.Orders.OrderType)
			assert.Equal(t, testCase.expectedResponse.Orders.Date,
				bodyResponse.Orders.Date)
			assert.Equal(t, testCase.expectedResponse.Orders.Brokerage.Name,
				bodyResponse.Orders.Brokerage.Name)
			assert.Equal(t, testCase.expectedResponse.Orders.Brokerage.Country,
				bodyResponse.Orders.Brokerage.Country)
			assert.Equal(t, testCase.expectedResponse.Orders.Asset.Symbol,
				bodyResponse.Orders.Asset.Symbol)
			assert.Equal(t, testCase.expectedResponse.Orders.Asset.Fullname,
				bodyResponse.Orders.Asset.Fullname)
		} else {
			assert.Nil(t, testCase.expectedResponse.Orders)
		}

	}

}

func TestFiberHandlersIntegrationTestGetOrdersFromAssetUser(t *testing.T) {
	type body struct {
		Success    bool                       `json:"success"`
		Message    string                     `json:"message"`
		Error      string                     `json:"error"`
		Code       int                        `json:"code"`
		OrdersInfo []presenter.OrderApiReturn `json:"ordersInfo"`
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
				Code:       401,
				Success:    false,
				Message:    entity.ErrMessageApiAuthentication.Error(),
				Error:      "",
				OrdersInfo: nil,
			},
		},
		{
			idToken: "TestNoAdminID",
			expectedResponse: body{
				Code:       400,
				Success:    false,
				Message:    entity.ErrMessageApiRequest.Error(),
				Error:      entity.ErrInvalidApiQuerySymbolBlank.Error(),
				OrdersInfo: nil,
			},
		},
		{
			idToken: "TestNoAdminID",
			symbol:  "ERRO3",
			expectedResponse: body{
				Code:       404,
				Success:    false,
				Message:    entity.ErrMessageApiAssetSymbolUser.Error(),
				Error:      "",
				OrdersInfo: nil,
			},
		},
		{
			idToken: "TestAdminID",
			symbol:  "EGIE3",
			limit:   "2b",
			offset:  "0",
			expectedResponse: body{
				Code:       400,
				Success:    false,
				Message:    entity.ErrMessageApiRequest.Error(),
				Error:      entity.ErrInvalidOrderLimit.Error(),
				OrdersInfo: nil,
			},
			expectedOrderLengthReturn: 0,
		},
		{
			idToken: "TestAdminID",
			symbol:  "EGIE3",
			limit:   "2",
			offset:  "0a",
			expectedResponse: body{
				Code:       400,
				Success:    false,
				Message:    entity.ErrMessageApiRequest.Error(),
				Error:      entity.ErrInvalidOrderOffset.Error(),
				OrdersInfo: nil,
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
				Code:       500,
				Success:    false,
				Message:    entity.ErrMessageApiInternalError.Error(),
				Error:      entity.ErrInvalidOrderOrderBy.Error(),
				OrdersInfo: nil,
			},
			expectedOrderLengthReturn: 0,
		},
		{
			idToken: "TestAdminID",
			symbol:  "EGIE3",
			expectedResponse: body{
				Code:       200,
				Success:    true,
				Message:    "Orders returned successfully",
				Error:      "",
				OrdersInfo: nil,
			},
			expectedOrderLengthReturn: 6,
		},
		{
			idToken: "TestAdminID",
			symbol:  "EGIE3",
			limit:   "2",
			offset:  "0",
			orderBy: "desc",
			expectedResponse: body{
				Code:       200,
				Success:    true,
				Message:    "Orders returned successfully",
				Error:      "",
				OrdersInfo: nil,
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
				Code:       200,
				Success:    true,
				Message:    "Orders returned successfully",
				Error:      "",
				OrdersInfo: nil,
			},
			expectedOrderLengthReturn: 2,
		},
		{
			idToken: "TestAdminID",
			symbol:  "EGIE3",
			limit:   "3",
			offset:  "6",
			orderBy: "asc",
			expectedResponse: body{
				Code:       404,
				Success:    false,
				Message:    entity.ErrMessageApiAssetSymbolUser.Error(),
				Error:      "",
				OrdersInfo: nil,
			},
			expectedOrderLengthReturn: 0,
		},
	}

	DBpool := connectDatabase()

	order, applicationsLogics := configureOrdersApp(DBpool)

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
	api.Get("/orders", order.GetOrdersFromAssetUser)

	for _, testCase := range tests {
		bodyResponse := body{}

		resp, _ := fiberHandlers.MockHttpRequest(app, "GET", "/api/orders?symbol="+
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

		if testCase.expectedResponse.OrdersInfo != nil {
			assert.Equal(t, testCase.expectedOrderLengthReturn,
				len(bodyResponse.OrdersInfo))
		} else {
			assert.Nil(t, testCase.expectedResponse.OrdersInfo)
		}

	}
}

func TestFiberHandlersIntegrationTestDeleteOrderFromUser(t *testing.T) {
	type body struct {
		Success bool    `json:"success"`
		Message string  `json:"message"`
		Error   string  `json:"error"`
		Code    int     `json:"code"`
		Order   *string `json:"order"`
	}

	type bodyGetOrders struct {
		Success    bool                       `json:"success"`
		Message    string                     `json:"message"`
		Error      string                     `json:"error"`
		Code       int                        `json:"code"`
		OrdersInfo []presenter.OrderApiReturn `json:"ordersInfo"`
	}

	type test struct {
		idToken                   string
		orderId                   string
		expectedResponse          body
		expectedOrderLengthReturn int
	}

	DBpool := connectDatabase()

	order, applicationsLogics := configureOrdersApp(DBpool)

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
	api.Delete("orders/:id", order.DeleteOrderFromUser)
	api.Get("/orders", order.GetOrdersFromAssetUser) //

	resp, _ := fiberHandlers.MockHttpRequest(app, "GET", "/api/orders?symbol="+
		"EGIE3&limit=1&offset=0&orderBy=DESC", "", "TestAdminID", nil)
	bodyResponseBytes, _ := ioutil.ReadAll(resp.Body)

	bodyResponse := bodyGetOrders{}
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
				Order:   nil,
			},
		},
		{
			idToken: "TestNoAdminID",
			orderId: invalidId,
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   "ERROR: invalid input syntax for type uuid: \"" + invalidId + "\" (SQLSTATE 22P02)",
				Order:   nil,
			},
		},
		{
			idToken: "TestAdminID",
			orderId: "30ba4b45-4a70-11eb-99fd-ad786a821574",
			expectedResponse: body{
				Code:    404,
				Success: false,
				Message: entity.ErrMessageApiOrderId.Error(),
				Error:   "",
				Order:   nil,
			},
		},
		{
			idToken: "TestAdminID",
			orderId: orderTestId,
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Order deleted successfully",
				Error:   "",
				Order:   &orderTestId,
			},
		},
	}

	for _, testCase := range tests {
		bodyResponse := body{}

		resp, _ := fiberHandlers.MockHttpRequest(app, "DELETE", "/api/orders/"+
			testCase.orderId, "", testCase.idToken, nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &bodyResponse)
		bodyResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResponse, bodyResponse)
	}
}

func TestFiberHandlersIntegrationTestUpdateOrderFromUser(t *testing.T) {
	type body struct {
		Success bool                      `json:"success"`
		Message string                    `json:"message"`
		Error   string                    `json:"error"`
		Code    int                       `json:"code"`
		Order   *presenter.OrderApiReturn `json:"order"`
	}

	type bodyGetOrders struct {
		Success    bool                       `json:"success"`
		Message    string                     `json:"message"`
		Error      string                     `json:"error"`
		Code       int                        `json:"code"`
		OrdersInfo []presenter.OrderApiReturn `json:"ordersInfo"`
	}

	type test struct {
		idToken                   string
		contentType               string
		orderId                   string
		price                     float64
		quantity                  float64
		orderType                 string
		date                      string
		brokerage                 string
		expectedResponse          body
		expectedOrderLengthReturn int
	}

	DBpool := connectDatabase()

	order, applicationsLogics := configureOrdersApp(DBpool)

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
	api.Put("orders/:id", order.UpdateOrderFromUser)
	api.Get("/orders", order.GetOrdersFromAssetUser) //

	resp, _ := fiberHandlers.MockHttpRequest(app, "GET", "/api/orders?symbol="+
		"EGIE3&limit=1&offset=0&orderBy=DESC", "", "TestAdminID", nil)
	bodyResponseBytes, _ := ioutil.ReadAll(resp.Body)

	bodyResponse := bodyGetOrders{}
	json.Unmarshal(bodyResponseBytes, &bodyResponse)
	bodyResponse.Code = resp.StatusCode
	orderTestId := bodyResponse.OrdersInfo[0].Id

	invalidId := "3erroa"
	layout := "2006-01-02"
	dateString := "2021-10-05"
	dateFormatted, _ := time.Parse(layout, dateString)
	tests := []test{
		{
			idToken: "INVALID_ID_TOKEN",
			expectedResponse: body{
				Code:    401,
				Success: false,
				Message: entity.ErrMessageApiAuthentication.Error(),
				Error:   "",
				Order:   nil,
			},
		},
		{
			idToken: "TestNoAdminID",

			orderId:     invalidId,
			contentType: "application/pdf",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidApiBody.Error(),
				Order:   nil,
			},
		},
		{
			idToken:     "TestNoAdminID",
			orderId:     invalidId,
			contentType: "application/json",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidApiOrderUpdate.Error(),
				Order:   nil,
			},
		},
		{
			idToken:     "TestNoAdminID",
			orderId:     invalidId,
			contentType: "application/json",
			price:       20.20,
			quantity:    5,
			orderType:   "buy",
			date:        dateString,
			brokerage:   "Clear",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   "scany: rows final error: ERROR: invalid input syntax for type uuid: \"" + invalidId + "\" (SQLSTATE 22P02)",
				Order:   nil,
			},
		},
		{
			idToken:     "TestNoAdminID",
			orderId:     "30ba4b45-4a70-11eb-99fd-ad786a821574",
			contentType: "application/json",
			price:       20.20,
			quantity:    5,
			orderType:   "buy",
			date:        dateString,
			brokerage:   "Clear",
			expectedResponse: body{
				Code:    404,
				Success: false,
				Message: entity.ErrMessageApiOrderId.Error(),
				Error:   "",
				Order:   nil,
			},
		},
		{
			idToken:     "TestAdminID",
			orderId:     orderTestId,
			contentType: "application/json",
			price:       20.20,
			quantity:    5.2,
			orderType:   "buy",
			date:        dateString,
			brokerage:   "Clear",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidOrderQuantityBrazil.Error(),
				Order:   nil,
			},
		},
		{
			idToken:     "TestAdminID",
			orderId:     orderTestId,
			contentType: "application/json",
			price:       20.20,
			quantity:    5,
			orderType:   "err",
			date:        dateString,
			brokerage:   "Clear",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidOrderType.Error(),
				Order:   nil,
			},
		},
		{
			idToken:     "TestAdminID",
			orderId:     orderTestId,
			contentType: "application/json",
			price:       20.20,
			quantity:    5,
			orderType:   "buy",
			date:        dateString,
			brokerage:   "Error",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidBrokerageNameSearch.Error(),
				Order:   nil,
			},
		},
		{
			idToken:     "TestAdminID",
			orderId:     orderTestId,
			contentType: "application/json",
			price:       20.20,
			quantity:    5,
			orderType:   "buy",
			date:        dateString,
			brokerage:   "Clear",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Order updated successfully",
				Error:   "",
				Order: &presenter.OrderApiReturn{
					Id:        orderTestId,
					Quantity:  5,
					Price:     20.20,
					OrderType: "buy",
					Date:      dateFormatted,
					Currency:  "BRL",
					Brokerage: &presenter.Brokerage{
						Id:      bodyResponse.OrdersInfo[0].Brokerage.Id,
						Name:    bodyResponse.OrdersInfo[0].Brokerage.Name,
						Country: bodyResponse.OrdersInfo[0].Brokerage.Country,
					},
				},
			},
		},
	}

	for _, testCase := range tests {
		bodyResponse := body{}

		bodyRequest := presenter.OrderBody{
			Price:     testCase.price,
			Quantity:  testCase.quantity,
			OrderType: testCase.orderType,
			Date:      testCase.date,
			Brokerage: testCase.brokerage,
		}

		resp, _ := fiberHandlers.MockHttpRequest(app, "PUT", "/api/orders/"+
			testCase.orderId, testCase.contentType, testCase.idToken, bodyRequest)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &bodyResponse)
		bodyResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)

		assert.Equal(t, testCase.expectedResponse, bodyResponse)
	}
}
