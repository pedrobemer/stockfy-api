package integration_tests

import (
	"encoding/json"
	"fmt"
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
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

func configureEventsApp(dbpool *pgxpool.Pool) (fiberHandlers.EventsApi,
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

	orders := fiberHandlers.EventsApi{
		ApplicationLogic:   *applicationLogics,
		ExternalInterfaces: externalInterface,
		LogicApi:           logicApiUseCases,
	}

	return orders, *applicationLogics
}

func TestFiberHandlersIntegrationTestCreateEvent(t *testing.T) {

	type body struct {
		Success bool                       `json:"success"`
		Message string                     `json:"message"`
		Error   string                     `json:"error"`
		Code    int                        `json:"code"`
		Orders  []presenter.OrderApiReturn `json:"orders"`
	}

	type test struct {
		idToken          string
		contentType      string
		symbol           string
		fullname         string
		brokerage        string
		eventRate        float64
		price            float64
		currency         string
		eventType        string
		date             string
		assetType        string
		country          string
		expectedResponse body
	}

	dateString := "2021-10-02"
	// layout := "2006-01-02"
	// dateFormatted, _ := time.Parse(layout, dateString)
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
			price:       20.00,
			eventRate:   5,
			currency:    "BRL",
			eventType:   "buy",
			date:        dateString,
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidEventType.Error(),
				Orders:  nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			symbol:      "VIVA3",
			price:       20.00,
			eventRate:   5,
			currency:    "BRL",
			eventType:   "split",
			date:        dateString,
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidAssetSymbolUserRelation.Error(),
				Orders:  nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			symbol:      "EGIE3",
			price:       20.00,
			eventRate:   5.1,
			currency:    "BRL",
			eventType:   "split",
			date:        "2021-07-02",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrEmptyQuery.Error(),
				Orders:  nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			symbol:      "EGIE3",
			price:       20.00,
			eventRate:   5,
			currency:    "USD",
			eventType:   "split",
			date:        dateString,
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidBrazilCurrency.Error(),
				Orders:  nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			symbol:      "EGIE3",
			price:       20.00,
			eventRate:   5,
			currency:    "BRL",
			eventType:   "split",
			date:        dateString,
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalideNonZeroOrderPrice.Error(),
				Orders:  nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			symbol:      "EGIE3",
			price:       -20.00,
			eventRate:   5,
			currency:    "BRL",
			eventType:   "bonification",
			date:        dateString,
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidNegativeOrderPrice.Error(),
				Orders:  nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			symbol:      "EGIE3",
			price:       20.00,
			eventRate:   5,
			currency:    "BRL",
			eventType:   "demerge",
			date:        dateString,
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidPositiveOrderPrice.Error(),
				Orders:  nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			symbol:      "EGIE3",
			price:       20.00,
			eventRate:   -5,
			currency:    "BRL",
			eventType:   "bonification",
			date:        dateString,
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidOrderBuyQuantity.Error(),
				Orders:  nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			symbol:      "EGIE3",
			price:       0,
			eventRate:   -5,
			currency:    "BRL",
			eventType:   "split",
			date:        dateString,
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidOrderBuyQuantity.Error(),
				Orders:  nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			symbol:      "TAEE3",
			price:       5.49,
			eventRate:   5,
			currency:    "BRL",
			eventType:   "bonification",
			date:        dateString,
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Event registered successfully",
				Error:   "",
				Orders: []presenter.OrderApiReturn{
					{
						Quantity:  1.6,
						Price:     5.49,
						Currency:  "BRL",
						OrderType: "bonification",
						Date:      entity.StringToTime(dateString),
						Brokerage: &presenter.Brokerage{
							Name:    "Avenue",
							Country: "US",
						},
					},
					{
						Quantity:  10.4,
						Price:     5.49,
						Currency:  "BRL",
						OrderType: "bonification",
						Date:      entity.StringToTime(dateString),
						Brokerage: &presenter.Brokerage{
							Name:    "Clear",
							Country: "BR",
						},
					},
				},
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			symbol:      "ARZZ3",
			price:       5.49,
			eventRate:   5,
			currency:    "BRL",
			eventType:   "bonification",
			date:        dateString,
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Event registered successfully",
				Error:   "",
				Orders: []presenter.OrderApiReturn{
					{
						Quantity:  4.8,
						Price:     5.49,
						Currency:  "BRL",
						OrderType: "bonification",
						Date:      entity.StringToTime(dateString),
						Brokerage: &presenter.Brokerage{
							Name:    "Avenue",
							Country: "US",
						},
					},
				},
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			symbol:      "ARZZ3",
			price:       0,
			eventRate:   2,
			currency:    "BRL",
			eventType:   "split",
			date:        "2021-10-10",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Event registered successfully",
				Error:   "",
				Orders: []presenter.OrderApiReturn{
					{
						Quantity:  16.9,
						Price:     0,
						Currency:  "BRL",
						OrderType: "split",
						Date:      entity.StringToTime("2021-10-10"),
						Brokerage: &presenter.Brokerage{
							Name:    "Avenue",
							Country: "US",
						},
					},
				},
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			symbol:      "TAEE3",
			price:       0,
			eventRate:   2,
			currency:    "BRL",
			eventType:   "split",
			date:        "2021-10-10",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Event registered successfully",
				Error:   "",
				Orders: []presenter.OrderApiReturn{
					{
						Quantity:  4.8,
						Price:     0,
						Currency:  "BRL",
						OrderType: "split",
						Date:      entity.StringToTime("2021-10-10"),
						Brokerage: &presenter.Brokerage{
							Name:    "Avenue",
							Country: "US",
						},
					},
					{
						Quantity:  31.2,
						Price:     0,
						Currency:  "BRL",
						OrderType: "split",
						Date:      entity.StringToTime("2021-10-10"),
						Brokerage: &presenter.Brokerage{
							Name:    "Clear",
							Country: "BR",
						},
					},
				},
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			symbol:      "TAEE3",
			price:       -8.23,
			eventRate:   12,
			currency:    "BRL",
			eventType:   "demerge",
			date:        "2021-10-20",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Event registered successfully",
				Error:   "",
				Orders: []presenter.OrderApiReturn{
					{
						Quantity:  0,
						Price:     -118.512,
						Currency:  "BRL",
						OrderType: "demerge",
						Date:      entity.StringToTime("2021-10-20"),
						Brokerage: &presenter.Brokerage{
							Name:    "Avenue",
							Country: "US",
						},
					},
					{
						Quantity:  0,
						Price:     -770.328,
						Currency:  "BRL",
						OrderType: "demerge",
						Date:      entity.StringToTime("2021-10-20"),
						Brokerage: &presenter.Brokerage{
							Name:    "Clear",
							Country: "BR",
						},
					},
				},
			},
		},
	}

	DBpool := connectDatabase()

	event, applicationsLogics := configureEventsApp(DBpool)

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
	api.Post("/events", event.CreateEventOrder)

	for _, testCase := range tests {
		bodyResponse := body{}
		bodyRequestStruct := presenter.EventBody{
			Symbol:    testCase.symbol,
			EventRate: testCase.eventRate,
			Price:     testCase.price,
			Currency:  testCase.currency,
			EventType: testCase.eventType,
			Date:      testCase.date,
		}

		resp, _ := fiberHandlers.MockHttpRequest(app, "POST", "/api/events",
			testCase.contentType, testCase.idToken, bodyRequestStruct)
		fmt.Println("TEST")
		fmt.Println(resp)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &bodyResponse)
		bodyResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResponse.Code, bodyResponse.Code)
		assert.Equal(t, testCase.expectedResponse.Success, bodyResponse.Success)
		assert.Equal(t, testCase.expectedResponse.Message, bodyResponse.Message)
		assert.Equal(t, testCase.expectedResponse.Error, bodyResponse.Error)

		if testCase.expectedResponse.Orders != nil {
			for i := range testCase.expectedResponse.Orders {
				for j := range bodyResponse.Orders {
					if bodyResponse.Orders[j].Quantity == testCase.expectedResponse.Orders[i].Quantity &&
						bodyResponse.Orders[j].Price == testCase.expectedResponse.Orders[i].Price {
						assert.Equal(t, testCase.expectedResponse.Orders[i].Quantity,
							bodyResponse.Orders[j].Quantity)
						assert.Equal(t, testCase.expectedResponse.Orders[i].Price,
							bodyResponse.Orders[j].Price)
						assert.Equal(t, testCase.expectedResponse.Orders[i].Currency,
							bodyResponse.Orders[j].Currency)
						assert.Equal(t, testCase.expectedResponse.Orders[i].OrderType,
							bodyResponse.Orders[j].OrderType)
						assert.Equal(t, testCase.expectedResponse.Orders[i].Date,
							bodyResponse.Orders[j].Date)
						assert.Equal(t, testCase.expectedResponse.Orders[i].Brokerage.Name,
							bodyResponse.Orders[j].Brokerage.Name)
						assert.Equal(t, testCase.expectedResponse.Orders[i].Brokerage.Country,
							bodyResponse.Orders[j].Brokerage.Country)
					}
				}

			}
		} else {
			assert.Nil(t, testCase.expectedResponse.Orders)
		}

	}

}
