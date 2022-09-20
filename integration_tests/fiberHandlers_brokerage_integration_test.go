package integration_tests

import (
	"encoding/json"
	"io/ioutil"
	"stockfyApi/api/handlers/fiberHandlers"
	"stockfyApi/api/middleware"
	"stockfyApi/api/presenter"
	"stockfyApi/database/postgresql"
	"stockfyApi/entity"
	"stockfyApi/usecases"
	"stockfyApi/usecases/user"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

func configureBrokerageApp(dbpool *pgxpool.Pool) (fiberHandlers.BrokerageApi,
	usecases.Applications) {

	dbInterfaces := postgresql.NewPostgresInstance(dbpool)
	firebaseInterface := user.NewExternalApi()

	applicationLogics := usecases.NewApplications(dbInterfaces,
		firebaseInterface)

	brokerage := fiberHandlers.BrokerageApi{
		ApplicationLogic: *applicationLogics,
	}

	return brokerage, *applicationLogics
}

func TestFiberHandlersIntegrationTestGetBrokerageFirm(t *testing.T) {

	type body struct {
		Success   bool                 `json:"success"`
		Message   string               `json:"message"`
		Error     string               `json:"error"`
		Code      int                  `json:"code"`
		Brokerage *presenter.Brokerage `json:"brokerage"`
	}

	type test struct {
		idToken          string
		brokerageName    string
		expectedResponse body
	}

	tests := []test{
		{
			idToken: "INVALID_ID_TOKEN",
			expectedResponse: body{
				Code:      401,
				Success:   false,
				Message:   entity.ErrMessageApiAuthentication.Error(),
				Error:     "",
				Brokerage: nil,
			},
		},
		{
			idToken:       "TestNoAdminID",
			brokerageName: "Error",
			expectedResponse: body{
				Code:      404,
				Success:   false,
				Message:   entity.ErrInvalidBrokerageNameSearch.Error(),
				Error:     "",
				Brokerage: nil,
			},
		},
		{
			idToken:       "TestNoAdminID",
			brokerageName: "Clear",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Brokerage firm information returned successfully",
				Error:   "",
				Brokerage: &presenter.Brokerage{
					Name:    "Clear",
					Country: "BR",
				},
			},
		},
	}

	DBpool := connectDatabase()

	brokerage, applicationsLogics := configureBrokerageApp(DBpool)

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
	api.Get("/brokerage/:name", brokerage.GetBrokerageFirm)

	for _, testCase := range tests {
		bodyResponse := body{}

		resp, _ := fiberHandlers.MockHttpRequest(app, "GET", "/api/brokerage/"+
			testCase.brokerageName, "", testCase.idToken, nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &bodyResponse)
		bodyResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResponse.Code, bodyResponse.Code)
		assert.Equal(t, testCase.expectedResponse.Success, bodyResponse.Success)
		assert.Equal(t, testCase.expectedResponse.Message, bodyResponse.Message)
		assert.Equal(t, testCase.expectedResponse.Error, bodyResponse.Error)

		if testCase.expectedResponse.Brokerage != nil {
			assert.Equal(t, testCase.expectedResponse.Brokerage.Name,
				bodyResponse.Brokerage.Name)
			assert.Equal(t, testCase.expectedResponse.Brokerage.Country,
				bodyResponse.Brokerage.Country)
		} else {
			assert.Nil(t, testCase.expectedResponse.Brokerage)
		}

	}
}

func TestFiberHandlersIntegrationTestGetBrokerageFirms(t *testing.T) {

	type body struct {
		Success   bool                  `json:"success"`
		Message   string                `json:"message"`
		Error     string                `json:"error"`
		Code      int                   `json:"code"`
		Brokerage []presenter.Brokerage `json:"brokerage"`
	}

	type test struct {
		idToken          string
		countryQuery     string
		expectedResponse body
	}

	tests := []test{
		{
			idToken: "INVALID_ID_TOKEN",
			expectedResponse: body{
				Code:      401,
				Success:   false,
				Message:   entity.ErrMessageApiAuthentication.Error(),
				Error:     "",
				Brokerage: nil,
			},
		},
		{
			idToken:      "TestNoAdminID",
			countryQuery: "ERR",
			expectedResponse: body{
				Code:      400,
				Success:   false,
				Message:   entity.ErrMessageApiRequest.Error(),
				Error:     entity.ErrInvalidCountryCode.Error(),
				Brokerage: nil,
			},
		},
		{
			idToken:      "TestNoAdminID",
			countryQuery: "BR",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Returned successfully the brokerage firms information",
				Error:   "",
				Brokerage: []presenter.Brokerage{
					{
						Name:    "Clear",
						Country: "BR",
					},
					{
						Name:    "Rico",
						Country: "BR",
					},
				},
			},
		},
	}

	DBpool := connectDatabase()

	brokerage, applicationsLogics := configureBrokerageApp(DBpool)

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
	api.Get("/brokerage", brokerage.GetBrokerageFirms)

	for _, testCase := range tests {
		bodyResponse := body{}

		resp, _ := fiberHandlers.MockHttpRequest(app, "GET", "/api/brokerage?"+
			"country="+testCase.countryQuery, "", testCase.idToken, nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &bodyResponse)
		bodyResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResponse.Code, bodyResponse.Code)
		assert.Equal(t, testCase.expectedResponse.Success, bodyResponse.Success)
		assert.Equal(t, testCase.expectedResponse.Message, bodyResponse.Message)
		assert.Equal(t, testCase.expectedResponse.Error, bodyResponse.Error)

		if testCase.expectedResponse.Brokerage != nil {
			for i, brokerage := range bodyResponse.Brokerage {
				assert.Equal(t, testCase.expectedResponse.Brokerage[i].Name,
					brokerage.Name)
				assert.Equal(t, testCase.expectedResponse.Brokerage[i].Country,
					brokerage.Country)
			}
		} else {
			assert.Nil(t, testCase.expectedResponse.Brokerage)
		}

	}
}
