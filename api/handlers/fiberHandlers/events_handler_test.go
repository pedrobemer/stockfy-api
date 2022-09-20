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

func TestApiCreateEventOrder(t *testing.T) {
	type body struct {
		Success bool                       `json:"success"`
		Message string                     `json:"message"`
		Error   string                     `json:"error"`
		Code    int                        `json:"code"`
		Orders  []presenter.OrderApiReturn `json:"orders"`
	}

	type test struct {
		idToken      string
		contentType  string
		bodyReq      presenter.EventBody
		expectedResp body
	}

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
			bodyReq: presenter.EventBody{
				Symbol:    "TEST3",
				EventRate: 4,
				Price:     29.10,
				EventType: "buy",
				Currency:  "BRL",
				Date:      "2021-10-01",
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
			bodyReq: presenter.EventBody{
				Symbol:    "TEST3",
				EventRate: 4,
				Price:     29.10,
				EventType: "buy",
				Currency:  "BRL",
				Date:      "2021-10-01",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidEventType.Error(),
				Code:    400,
				Orders:  nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			bodyReq: presenter.EventBody{
				Symbol:    "TEST3",
				EventRate: 4,
				Price:     29.10,
				EventType: "buy",
				Currency:  "BRL",
				Date:      "2021-10-01",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidEventType.Error(),
				Code:    400,
				Orders:  nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			bodyReq: presenter.EventBody{
				Symbol:    "EMPTYQUERY",
				EventRate: 4,
				Price:     29.10,
				EventType: "split",
				Currency:  "BRL",
				Date:      "2021-10-01",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrEmptyQuery.Error(),
				Code:    400,
				Orders:  nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			bodyReq: presenter.EventBody{
				Symbol:    "CREATEORDERERROR",
				EventRate: 4,
				Price:     29.10,
				EventType: "split",
				Currency:  "BRL",
				Date:      "2021-10-01",
			},
			expectedResp: body{
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Error:   errors.New("Create Order Error").Error(),
				Code:    500,
				Orders:  nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			bodyReq: presenter.EventBody{
				Symbol:    "TEST3",
				EventRate: 4,
				Price:     5.19,
				EventType: "split",
				Currency:  "BRL",
				Date:      "2021-10-01",
			},
			expectedResp: body{
				Success: true,
				Message: "Event registered successfully",
				Error:   "",
				Code:    200,
				Orders: []presenter.OrderApiReturn{
					{
						Id:        "TestOrderID",
						Price:     5.19,
						Quantity:  5,
						Currency:  "BRL",
						OrderType: "split",
						Date:      entity.StringToTime("2021-10-01"),
						Brokerage: &presenter.Brokerage{
							Id:      "TestBrokerageID",
							Name:    "Test Brokerage",
							Country: "BR",
						},
					},
					{
						Id:        "TestOrderID2",
						Price:     5.19,
						Quantity:  5.4,
						Currency:  "BRL",
						OrderType: "split",
						Date:      entity.StringToTime("2021-10-01"),
						Brokerage: &presenter.Brokerage{
							Id:      "TestBrokerageID2",
							Name:    "Test Brokerage 2",
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

	// Declare Events Application Logic
	events := EventsApi{
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
	api.Post("/events", events.CreateEventOrder)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "POST", "/api/events", testCase.contentType,
			testCase.idToken, testCase.bodyReq)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}
}
