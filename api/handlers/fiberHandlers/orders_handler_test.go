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
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestApiGetOrdersFromAssetUser(t *testing.T) {
	type body struct {
		Success    bool                        `json:"success"`
		Message    string                      `json:"message"`
		Error      string                      `json:"error"`
		Code       int                         `json:"code"`
		OrdersInfo *[]presenter.OrderApiReturn `json:"ordersInfo"`
	}

	type test struct {
		idToken      string
		contentType  string
		symbol       string
		orderBy      string
		limit        string
		offset       string
		expectedResp body
	}

	layOut := "2006-01-02"
	dateFormatted, _ := time.Parse(layOut, "2021-10-01")
	tests := []test{
		{
			idToken:     "ValidIdTokenWithoutEmailVerification",
			contentType: "application/json",
			symbol:      "",
			expectedResp: body{
				Success:    false,
				Message:    entity.ErrMessageApiAuthentication.Error(),
				Error:      "",
				Code:       401,
				OrdersInfo: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			symbol:      "",
			expectedResp: body{
				Success:    false,
				Message:    entity.ErrMessageApiRequest.Error(),
				Error:      entity.ErrInvalidApiQuerySymbolBlank.Error(),
				Code:       400,
				OrdersInfo: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			symbol:      "ERROR_REPOSITORY_ASSETBYUSER",
			expectedResp: body{
				Success:    false,
				Message:    entity.ErrMessageApiInternalError.Error(),
				Error:      errors.New("Unknown error on assets repository").Error(),
				Code:       500,
				OrdersInfo: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			symbol:      "UNKNOWN_SYMBOL",
			expectedResp: body{
				Success:    false,
				Message:    entity.ErrMessageApiAssetSymbolUser.Error(),
				Error:      "",
				Code:       404,
				OrdersInfo: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			symbol:      "ERROR_REPOSITORY_ORDERS",
			expectedResp: body{
				Success:    false,
				Message:    entity.ErrMessageApiInternalError.Error(),
				Error:      errors.New("Unknown error on orders repository").Error(),
				Code:       500,
				OrdersInfo: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			symbol:      "SYMBOL_WITHOUT_ORDERS",
			expectedResp: body{
				Success:    false,
				Message:    entity.ErrMessageApiAssetSymbolUser.Error(),
				Error:      "",
				Code:       404,
				OrdersInfo: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			symbol:      "VALID_SYMBOL",
			limit:       "4",
			offset:      "2a",
			expectedResp: body{
				Success:    false,
				Message:    entity.ErrMessageApiRequest.Error(),
				Error:      entity.ErrInvalidOrderOffset.Error(),
				Code:       400,
				OrdersInfo: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			symbol:      "VALID_SYMBOL",
			limit:       "4a",
			offset:      "2",
			expectedResp: body{
				Success:    false,
				Message:    entity.ErrMessageApiRequest.Error(),
				Error:      entity.ErrInvalidOrderLimit.Error(),
				Code:       400,
				OrdersInfo: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			symbol:      "VALID_SYMBOL",
			orderBy:     "error",
			limit:       "4",
			offset:      "0",
			expectedResp: body{
				Success:    false,
				Message:    entity.ErrMessageApiRequest.Error(),
				Error:      entity.ErrInvalidOrderOrderBy.Error(),
				Code:       400,
				OrdersInfo: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			symbol:      "VALID_SYMBOL",
			orderBy:     "desc",
			limit:       "4",
			offset:      "3",
			expectedResp: body{
				Success:    false,
				Message:    entity.ErrMessageApiAssetSymbolUser.Error(),
				Error:      "",
				Code:       404,
				OrdersInfo: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			symbol:      "TEST3",
			expectedResp: body{
				Success: true,
				Message: "Orders returned successfully",
				Error:   "",
				Code:    200,
				OrdersInfo: &[]presenter.OrderApiReturn{
					{
						Id:        "Order1",
						Quantity:  2,
						Price:     29.29,
						Currency:  "BRL",
						OrderType: "buy",
						Date:      dateFormatted,
						Brokerage: &presenter.Brokerage{
							Id:      "TestBrokerageID",
							Name:    "Test",
							Country: "BR",
						},
					},
					{
						Id:        "Order2",
						Quantity:  3,
						Price:     31.90,
						Currency:  "BRL",
						OrderType: "buy",
						Date:      dateFormatted,
						Brokerage: &presenter.Brokerage{
							Id:      "TestBrokerageID",
							Name:    "Test",
							Country: "BR",
						},
					},
				},
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			symbol:      "TEST3",
			orderBy:     "desc",
			limit:       "2",
			offset:      "0",
			expectedResp: body{
				Success: true,
				Message: "Orders returned successfully",
				Error:   "",
				Code:    200,
				OrdersInfo: &[]presenter.OrderApiReturn{
					{
						Id:        "Order1",
						Quantity:  2,
						Price:     29.29,
						Currency:  "BRL",
						OrderType: "buy",
						Date:      dateFormatted,
						Brokerage: &presenter.Brokerage{
							Id:      "TestBrokerageID",
							Name:    "Test",
							Country: "BR",
						},
					},
					{
						Id:        "Order2",
						Quantity:  3,
						Price:     31.90,
						Currency:  "BRL",
						OrderType: "buy",
						Date:      dateFormatted,
						Brokerage: &presenter.Brokerage{
							Id:      "TestBrokerageID",
							Name:    "Test",
							Country: "BR",
						},
					},
				},
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()
	logicApi := logicApi.NewMockApplication(*usecases)

	// Declare Sector Application Logic
	orders := OrderApi{
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
	api.Get("/orders", orders.GetOrdersFromAssetUser)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "GET", "/api/orders?symbol="+
			testCase.symbol+"&orderBy="+testCase.orderBy+"&limit="+testCase.limit+
			"&offset="+testCase.offset, testCase.contentType, testCase.idToken,
			nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}
}

func TestApiCreateOrder(t *testing.T) {
	type body struct {
		Success bool                      `json:"success"`
		Message string                    `json:"message"`
		Error   string                    `json:"error"`
		Code    int                       `json:"code"`
		Orders  *presenter.OrderApiReturn `json:"orders"`
	}

	type test struct {
		idToken      string
		contentType  string
		bodyReq      presenter.OrderBody
		expectedResp body
	}

	layOut := "2006-01-02"
	dateFormatted, _ := time.Parse(layOut, "2021-10-01")
	tests := []test{
		{
			idToken:     "INVALID_TOKEN",
			contentType: "application/json",
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiAuthentication.Error(),
				Error:   "",
				Code:    401,
				Orders:  nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/pdf",
			bodyReq: presenter.OrderBody{
				Symbol:    "TEST3",
				Fullname:  "Test Name",
				Brokerage: "Test Brokerage",
				Quantity:  -2,
				Price:     29.10,
				OrderType: "buy",
				Currency:  "BRL",
				Date:      "2021-10-01",
				Country:   "BR",
				AssetType: "ETF",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidApiBody.Error(),
				Code:    400,
				Orders:  nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			bodyReq: presenter.OrderBody{
				Symbol:    "TEST3",
				Fullname:  "Test Name",
				Brokerage: "Test Brokerage",
				Quantity:  -2,
				Price:     29.10,
				OrderType: "buy",
				Currency:  "BRL",
				Date:      "2021-10-01",
				Country:   "BR",
				AssetType: "ETF",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidOrderBuyQuantity.Error(),
				Code:    400,
				Orders:  nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			bodyReq: presenter.OrderBody{
				Symbol:    "TEST3",
				Fullname:  "Test Name",
				Brokerage: "Test Brokerage",
				Quantity:  2,
				Price:     29.10,
				OrderType: "sell",
				Currency:  "BRL",
				Date:      "2021-10-01",
				Country:   "BR",
				AssetType: "ETF",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidOrderSellQuantity.Error(),
				Code:    400,
				Orders:  nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			bodyReq: presenter.OrderBody{
				Symbol:    "TEST3",
				Fullname:  "Test Name",
				Brokerage: "Test Brokerage",
				Quantity:  2,
				Price:     29.10,
				OrderType: "error",
				Currency:  "BRL",
				Date:      "2021-10-01",
				Country:   "BR",
				AssetType: "ETF",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidOrderType.Error(),
				Code:    400,
				Orders:  nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			bodyReq: presenter.OrderBody{
				Symbol:    "TEST3",
				Fullname:  "Test Name",
				Brokerage: "Test Brokerage",
				Quantity:  2,
				Price:     29.10,
				OrderType: "buy",
				Currency:  "BRL",
				Date:      "2021-10-01",
				Country:   "ERROR",
				AssetType: "ETF",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidCountryCode.Error(),
				Code:    400,
				Orders:  nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			bodyReq: presenter.OrderBody{
				Symbol:    "TEST3",
				Fullname:  "Test Name",
				Brokerage: "Test Brokerage",
				Quantity:  2,
				Price:     29.10,
				OrderType: "buy",
				Currency:  "USD",
				Date:      "2021-10-01",
				Country:   "BR",
				AssetType: "ETF",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidBrazilCurrency.Error(),
				Code:    400,
				Orders:  nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			bodyReq: presenter.OrderBody{
				Symbol:    "TEST3",
				Fullname:  "Test Name",
				Brokerage: "Test Brokerage",
				Quantity:  2,
				Price:     29.10,
				OrderType: "buy",
				Currency:  "BRL",
				Date:      "2021-10-01",
				Country:   "US",
				AssetType: "ETF",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidUsaCurrency.Error(),
				Code:    400,
				Orders:  nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			bodyReq: presenter.OrderBody{
				Symbol:    "TEST3",
				Fullname:  "Test Name",
				Brokerage: "Test Brokerage",
				Quantity:  2.19,
				Price:     29.10,
				OrderType: "buy",
				Currency:  "BRL",
				Date:      "2021-10-01",
				Country:   "BR",
				AssetType: "ETF",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidOrderQuantityBrazil.Error(),
				Code:    400,
				Orders:  nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			bodyReq: presenter.OrderBody{
				Symbol:    "SYMBOL_ALREADY_EXISTS_ERROR",
				Fullname:  "Test Name",
				Brokerage: "Test Brokerage",
				Quantity:  2,
				Price:     29.10,
				OrderType: "buy",
				Currency:  "BRL",
				Date:      "2021-10-01",
				Country:   "BR",
				AssetType: "ETF",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Error:   errors.New("Unknown asset repository error").Error(),
				Code:    500,
				Orders:  nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			bodyReq: presenter.OrderBody{
				Symbol:    "ERROR_ASSETUSER_REPOSITORY",
				Fullname:  "Test Name",
				Brokerage: "Test Brokerage",
				Quantity:  2,
				Price:     29.10,
				OrderType: "buy",
				Currency:  "BRL",
				Date:      "2021-10-01",
				Country:   "BR",
				AssetType: "ETF",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Error:   errors.New("Unknown asset user repository error").Error(),
				Code:    500,
				Orders:  nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			bodyReq: presenter.OrderBody{
				Symbol:    "TEST3",
				Fullname:  "Test Name",
				Brokerage: "UNKNOWN_BROKERAGE",
				Quantity:  2,
				Price:     29.10,
				OrderType: "buy",
				Currency:  "BRL",
				Date:      "2021-10-01",
				Country:   "BR",
				AssetType: "ETF",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidBrokerageNameSearch.Error(),
				Code:    400,
				Orders:  nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			bodyReq: presenter.OrderBody{
				Symbol:    "SYMBOL_ALREADY_EXISTS",
				Fullname:  "Test Name",
				Brokerage: "Test Brokerage",
				Quantity:  2,
				Price:     29.10,
				OrderType: "buy",
				Currency:  "BRL",
				Date:      "2021-10-01",
				Country:   "BR",
				AssetType: "ETF",
			},
			expectedResp: body{
				Success: true,
				Message: "Order registered successfully",
				Error:   "",
				Code:    200,
				Orders: &presenter.OrderApiReturn{
					Id:        "TestOrderID",
					Quantity:  2,
					Price:     29.1,
					Currency:  "BRL",
					OrderType: "buy",
					Date:      dateFormatted,
					Brokerage: &presenter.Brokerage{
						Id:      "TestBrokerageID",
						Name:    "Test Brokerage",
						Country: "BR",
					},
					Asset: &presenter.AssetApiReturn{
						Id:       "TestID",
						Symbol:   "SYMBOL_ALREADY_EXISTS",
						Fullname: "Test Name",
					},
				},
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			bodyReq: presenter.OrderBody{
				Symbol:    "TESTE3",
				Fullname:  "Test Name",
				Brokerage: "Test Brokerage",
				Quantity:  2,
				Price:     29.10,
				OrderType: "buy",
				Currency:  "BRL",
				Date:      "2021-10-01",
				Country:   "BR",
				AssetType: "ETF",
			},
			expectedResp: body{
				Success: true,
				Message: "Order registered successfully",
				Error:   "",
				Code:    200,
				Orders: &presenter.OrderApiReturn{
					Id:        "TestOrderID",
					Quantity:  2,
					Price:     29.1,
					Currency:  "BRL",
					OrderType: "buy",
					Date:      dateFormatted,
					Brokerage: &presenter.Brokerage{
						Id:      "TestBrokerageID",
						Name:    "Test Brokerage",
						Country: "BR",
					},
					Asset: &presenter.AssetApiReturn{
						Id:       "TestID",
						Symbol:   "TESTE3",
						Fullname: "Test Name",
					},
				},
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()
	logicApi := logicApi.NewMockApplication(*usecases)

	// Declare Sector Application Logic
	orders := OrderApi{
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
	api.Post("/orders", orders.CreateUserOrder)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "POST", "/api/orders", testCase.contentType,
			testCase.idToken, testCase.bodyReq)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}
}

func TestApiDeleteOrderFromUser(t *testing.T) {

	type body struct {
		Success bool    `json:"success"`
		Message string  `json:"message"`
		Error   string  `json:"error"`
		Code    int     `json:"code"`
		Order   *string `json:"order"`
	}

	type test struct {
		idToken      string
		orderId      string
		expectedResp body
	}

	orderId := "TestOrderID"
	tests := []test{
		{
			idToken: "INVALID_TOKEN",
			orderId: "ERROR_REPOSITORY",
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiAuthentication.Error(),
				Error:   "",
				Code:    401,
				Order:   nil,
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			orderId: "ERROR_REPOSITORY",
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   errors.New("Unknown orders repository error").Error(),
				Code:    400,
				Order:   nil,
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			orderId: "INVALID_ORDER_ID",
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiOrderId.Error(),
				Error:   "",
				Code:    404,
				Order:   nil,
			},
		},
		{
			idToken: "ValidIdTokenWithoutPrivilegedUser",
			orderId: "TestOrderID",
			expectedResp: body{
				Success: true,
				Message: "Order deleted successfully",
				Error:   "",
				Code:    200,
				Order:   &orderId,
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()
	logicApi := logicApi.NewMockApplication(*usecases)

	// Declare Sector Application Logic
	orders := OrderApi{
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
	api.Delete("/orders/:id", orders.DeleteOrderFromUser)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "DELETE", "/api/orders/"+testCase.orderId,
			"application/json", testCase.idToken, nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}
}

func TestApiUpdateOrder(t *testing.T) {

	type body struct {
		Success bool                      `json:"success"`
		Message string                    `json:"message"`
		Error   string                    `json:"error"`
		Code    int                       `json:"code"`
		Order   *presenter.OrderApiReturn `json:"order"`
	}

	type test struct {
		idToken      string
		contentType  string
		orderId      string
		bodyReq      presenter.OrderBody
		expectedResp body
	}

	layOut := "2006-01-02"
	date := "2021-10-01"
	dateFormatted, _ := time.Parse(layOut, date)
	tests := []test{
		{
			contentType: "application/json",
			idToken:     "INVALID_TOKEN",
			orderId:     "ERROR_REPOSITORY",
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiAuthentication.Error(),
				Error:   "",
				Code:    401,
				Order:   nil,
			},
		},
		{
			contentType: "application/pdf",
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			orderId:     "VALID_ORDER_ID",
			bodyReq: presenter.OrderBody{
				OrderType: "Dividendos",
				Price:     30.29,
				Quantity:  2,
				Date:      "2021-10-01",
				Brokerage: "",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidApiBody.Error(),
				Code:    400,
				Order:   nil,
			},
		},
		{
			contentType: "application/json",
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			orderId:     "VALID_ORDER_ID",
			bodyReq: presenter.OrderBody{
				OrderType: "Dividendos",
				Price:     30.29,
				Quantity:  2,
				Date:      "2021-10-01",
				Brokerage: "",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidApiOrderUpdate.Error(),
				Code:    400,
				Order:   nil,
			},
		},
		{
			contentType: "application/json",
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			orderId:     "ERROR_ORDER_REPOSITORY",
			bodyReq: presenter.OrderBody{
				OrderType: "Dividendos",
				Price:     30.29,
				Quantity:  2,
				Date:      "2021-10-01",
				Brokerage: "Test Name",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Error:   errors.New("Unknown error in the order repository").Error(),
				Code:    500,
				Order:   nil,
			},
		},
		{
			contentType: "application/json",
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			orderId:     "UNKNOWN_ID",
			bodyReq: presenter.OrderBody{
				OrderType: "Dividendos",
				Price:     30.29,
				Quantity:  2,
				Date:      "2021-10-01",
				Brokerage: "Test Name",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiOrderId.Error(),
				Error:   "",
				Code:    404,
				Order:   nil,
			},
		},
		{
			contentType: "application/json",
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			orderId:     "VALID_ID",
			bodyReq: presenter.OrderBody{
				OrderType: "Dividendos",
				Price:     30.29,
				Quantity:  2,
				Date:      "2021-10-01",
				Brokerage: "Test Name",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidOrderType.Error(),
				Code:    400,
				Order:   nil,
			},
		},
		{
			contentType: "application/json",
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			orderId:     "VALID_ID",
			bodyReq: presenter.OrderBody{
				OrderType: "buy",
				Price:     30.29,
				Quantity:  -2,
				Date:      "2021-10-01",
				Brokerage: "Test Name",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidOrderBuyQuantity.Error(),
				Code:    400,
				Order:   nil,
			},
		},
		{
			contentType: "application/json",
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			orderId:     "VALID_ID",
			bodyReq: presenter.OrderBody{
				OrderType: "sell",
				Price:     30.29,
				Quantity:  2,
				Date:      "2021-10-01",
				Brokerage: "Test Name",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidOrderSellQuantity.Error(),
				Code:    400,
				Order:   nil,
			},
		},
		{
			contentType: "application/json",
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			orderId:     "VALID_ID",
			bodyReq: presenter.OrderBody{
				OrderType: "buy",
				Price:     -30.29,
				Quantity:  2,
				Date:      "2021-10-01",
				Brokerage: "Test Name",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidOrderPrice.Error(),
				Code:    400,
				Order:   nil,
			},
		},
		{
			contentType: "application/json",
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			orderId:     "VALID_ID",
			bodyReq: presenter.OrderBody{
				OrderType: "buy",
				Price:     30.29,
				Quantity:  2.19,
				Date:      "2021-10-01",
				Brokerage: "Test Name",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidOrderQuantityBrazil.Error(),
				Code:    400,
				Order:   nil,
			},
		},
		{
			contentType: "application/json",
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			orderId:     "VALID_ID",
			bodyReq: presenter.OrderBody{
				OrderType: "buy",
				Price:     30.29,
				Quantity:  2,
				Date:      "2021-10-01",
				Brokerage: "UNKNOWN_BROKERAGE",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidBrokerageNameSearch.Error(),
				Code:    400,
				Order:   nil,
			},
		},
		{
			contentType: "application/json",
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			orderId:     "ORDER_VALID_ID",
			bodyReq: presenter.OrderBody{
				OrderType: "buy",
				Price:     30.29,
				Quantity:  2,
				Date:      date,
				Brokerage: "Test Brokerage",
			},
			expectedResp: body{
				Success: true,
				Message: "Order updated successfully",
				Error:   "",
				Code:    200,
				Order: &presenter.OrderApiReturn{
					Id:        "ORDER_VALID_ID",
					Price:     30.29,
					Quantity:  2,
					Date:      dateFormatted,
					Currency:  "BRL",
					OrderType: "buy",
					Brokerage: &presenter.Brokerage{
						Id:      "TestBrokerageID",
						Name:    "Test Brokerage",
						Country: "BR",
					},
				},
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()
	logicApi := logicApi.NewMockApplication(*usecases)

	// Declare Sector Application Logic
	orders := OrderApi{
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
	api.Put("/orders/:id", orders.UpdateOrderFromUser)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "PUT", "/api/orders/"+testCase.orderId,
			testCase.contentType, testCase.idToken, testCase.bodyReq)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}
}
