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

func TestApiGetEarningsFromAssetUser(t *testing.T) {
	type body struct {
		Success bool                          `json:"success"`
		Message string                        `json:"message"`
		Error   string                        `json:"error"`
		Code    int                           `json:"code"`
		Earning []presenter.EarningsApiReturn `json:"earning"`
	}

	type test struct {
		idToken      string
		contentType  string
		pathQuery    string
		expectedResp body
	}

	layOut := "2006-01-02"
	dateFormatted, _ := time.Parse(layOut, "2021-10-01")
	tests := []test{
		{
			idToken:     "ValidIdTokenWithoutEmailVerification",
			contentType: "application/json",
			pathQuery:   "?symbol=valid",
			expectedResp: body{
				Code:    401,
				Success: false,
				Message: entity.ErrMessageApiAuthentication.Error(),
				Error:   "",
				Earning: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			pathQuery:   "",
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidApiQuerySymbolBlank.Error(),
				Earning: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			pathQuery:   "?symbol=ERROR_ASSET_REPOSITORY",
			expectedResp: body{
				Code:    500,
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Error:   errors.New("Unknown error in the asset repository").Error(),
				Earning: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			pathQuery:   "?symbol=ERROR_EARNINGS_REPOSITORY",
			expectedResp: body{
				Code:    500,
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Error:   errors.New("Unknown error in the earnings repository").Error(),
				Earning: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			pathQuery:   "?symbol=INVALID_SYMBOL",
			expectedResp: body{
				Code:    404,
				Success: false,
				Message: entity.ErrMessageApiAssetSymbolUser.Error(),
				Error:   "",
				Earning: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			pathQuery:   "?symbol=SYMBOL_WITHOUT_EARNINGS",
			expectedResp: body{
				Code:    404,
				Success: false,
				Message: entity.ErrMessageApiEarningAssetUser.Error(),
				Error:   "",
				Earning: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			pathQuery:   "?symbol=VALID_SYMBOL&orderBy=DESC&limit=ab&offset=2",
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidEarningsLimit.Error(),
				Earning: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			pathQuery:   "?symbol=VALID_SYMBOL&orderBy=DESC&limit=2&offset=2a",
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidEarningsOffset.Error(),
				Earning: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			pathQuery:   "?symbol=VALID_SYMBOL&orderBy=ERROR&limit=2&offset=2",
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidEarningsOrderBy.Error(),
				Earning: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			pathQuery:   "?symbol=VALID_SYMBOL&orderBy=DESC&limit=2&offset=0",
			expectedResp: body{
				Code:    200,
				Success: true,
				Message: "Earnings returned successfully",
				Error:   "",
				Earning: []presenter.EarningsApiReturn{
					{
						Id:       "Earnings1",
						Type:     "Dividendos",
						Earning:  5.29,
						Date:     &dateFormatted,
						Currency: "BRL",
						Asset: &presenter.AssetApiReturn{
							Id:     "TestID",
							Symbol: "VALID_SYMBOL",
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
	earnings := EarningsApi{
		ApplicationLogic: *usecases,
		ApiLogic:         logicApi,
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
	api.Get("/earnings", earnings.GetEarningsFromAssetUser)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "GET", "/api/earnings"+testCase.pathQuery,
			testCase.contentType, testCase.idToken, nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}

}

func TestApiCreateEarnings(t *testing.T) {
	type body struct {
		Success bool                         `json:"success"`
		Message string                       `json:"message"`
		Error   string                       `json:"error"`
		Code    int                          `json:"code"`
		Earning *presenter.EarningsApiReturn `json:"earning"`
	}

	type test struct {
		idToken      string
		contentType  string
		bodyReq      presenter.EarningsBody
		expectedResp body
	}

	layOut := "2006-01-02"
	dateFormatted, _ := time.Parse(layOut, "2021-10-01")
	tests := []test{
		{
			idToken:     "ValidIdTokenWithoutEmailVerification",
			contentType: "application/json",
			expectedResp: body{
				Code:    401,
				Success: false,
				Message: entity.ErrMessageApiAuthentication.Error(),
				Error:   "",
				Earning: nil,
			},
		},
		{
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/pdf",
			bodyReq: presenter.EarningsBody{
				Symbol:      "TEST3",
				Amount:      6.00,
				Currency:    "BRL",
				EarningType: "Dividendos",
				Date:        "2021-10-01",
			},
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Earning: nil,
				Error:   entity.ErrInvalidApiBody.Error(),
			},
		},
		{
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/json",
			bodyReq: presenter.EarningsBody{
				Symbol:      "",
				Amount:      6.00,
				Currency:    "BRL",
				EarningType: "Dividendos",
				Date:        "2021-10-01",
			},
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Earning: nil,
				Error:   entity.ErrInvalidEarningsCreateBlankFields.Error(),
			},
		},
		{
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/json",
			bodyReq: presenter.EarningsBody{
				Symbol:      "TEST3",
				Amount:      -6.00,
				Currency:    "BRL",
				EarningType: "Dividendos",
				Date:        "2021-10-01",
			},
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Earning: nil,
				Error:   entity.ErrInvalidEarningsAmount.Error(),
			},
		},
		{
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/json",
			bodyReq: presenter.EarningsBody{
				Symbol:      "TEST3",
				Amount:      6.00,
				Currency:    "BRL",
				EarningType: "ERROR",
				Date:        "2021-10-01",
			},
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Earning: nil,
				Error:   entity.ErrInvalidEarningType.Error(),
			},
		},
		{
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/json",
			bodyReq: presenter.EarningsBody{
				Symbol:      "ERROR_ASSET_REPOSITORY",
				Amount:      6.00,
				Currency:    "BRL",
				EarningType: "JCP",
				Date:        "2021-10-01",
			},
			expectedResp: body{
				Code:    500,
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Earning: nil,
				Error:   errors.New("Unknown error in the asset repository").Error(),
			},
		},
		{
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/json",
			bodyReq: presenter.EarningsBody{
				Symbol:      "UNKNOWN_SYMBOL",
				Amount:      6.00,
				Currency:    "BRL",
				EarningType: "JCP",
				Date:        "2021-10-01",
			},
			expectedResp: body{
				Code:    404,
				Success: false,
				Message: entity.ErrMessageApiAssetSymbolUser.Error(),
				Earning: nil,
				Error:   "",
			},
		},
		{
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/json",
			bodyReq: presenter.EarningsBody{
				Symbol:      "TEST3",
				Amount:      6.00,
				Currency:    "BRL",
				EarningType: "JCP",
				Date:        "2021-10-01",
			},
			expectedResp: body{
				Code:    200,
				Success: true,
				Message: "Earning registered successfully",
				Earning: &presenter.EarningsApiReturn{
					Id:       "TestEarningID",
					Earning:  6.00,
					Type:     "JCP",
					Currency: "BRL",
					Date:     &dateFormatted,
					Asset: &presenter.AssetApiReturn{
						Id:     "TestAssetID",
						Symbol: "TEST3",
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
	earnings := EarningsApi{
		ApplicationLogic: *usecases,
		ApiLogic:         logicApi,
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
	api.Post("/earnings", earnings.CreateEarnings)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "POST", "/api/earnings",
			testCase.contentType, testCase.idToken, testCase.bodyReq)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}
}

func TestApiUpdateEarnings(t *testing.T) {
	type body struct {
		Success bool                         `json:"success"`
		Message string                       `json:"message"`
		Error   string                       `json:"error"`
		Code    int                          `json:"code"`
		Earning *presenter.EarningsApiReturn `json:"earning"`
	}

	type test struct {
		idToken      string
		contentType  string
		earningId    string
		bodyReq      presenter.EarningsBody
		expectedResp body
	}

	date := "2021-10-01"
	layOut := "2006-01-02"
	dateFormatted, _ := time.Parse(layOut, date)
	tests := []test{
		{
			idToken:     "ValidIdTokenWithoutEmailVerification",
			contentType: "application/json",
			earningId:   "TestEarningId",
			expectedResp: body{
				Code:    401,
				Success: false,
				Message: entity.ErrMessageApiAuthentication.Error(),
				Error:   "",
				Earning: nil,
			},
		},
		{
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/pdf",
			earningId:   "TestEarningId",
			bodyReq: presenter.EarningsBody{
				Amount:      6.00,
				EarningType: "Dividendos",
				Date:        "2021-10-01",
			},
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Earning: nil,
				Error:   entity.ErrInvalidApiBody.Error(),
			},
		},
		{
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/json",
			earningId:   "ERROR_EARNING_REPOSITORY",
			bodyReq: presenter.EarningsBody{
				Amount:      6.00,
				EarningType: "Dividendos",
				Date:        "2021-10-01",
			},
			expectedResp: body{
				Code:    500,
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Earning: nil,
				Error:   errors.New("Unknown error in the earning repository").Error(),
			},
		},
		{
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/json",
			earningId:   "ERROR_ASSET_REPOSITORY",
			bodyReq: presenter.EarningsBody{
				Amount:      6.00,
				EarningType: "Dividendos",
				Date:        "2021-10-01",
			},
			expectedResp: body{
				Code:    500,
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Earning: nil,
				Error:   errors.New("Unknown error in the asset repository").Error(),
			},
		},
		{
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/json",
			earningId:   "ERROR_UPDATE_EARNING_REPOSITORY",
			bodyReq: presenter.EarningsBody{
				Amount:      6.00,
				EarningType: "Dividendos",
				Date:        "2021-10-01",
			},
			expectedResp: body{
				Code:    500,
				Success: false,
				Message: entity.ErrMessageApiInternalError.Error(),
				Earning: nil,
				Error:   errors.New("Unknown in the update earning function").Error(),
			},
		},
		{
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/json",
			earningId:   "UNKNOWN_EARNING_ID",
			bodyReq: presenter.EarningsBody{
				Amount:      6.00,
				EarningType: "Dividendos",
				Date:        "2021-10-01",
			},
			expectedResp: body{
				Code:    404,
				Success: false,
				Message: entity.ErrMessageApiEarningId.Error(),
				Earning: nil,
				Error:   "",
			},
		},
		{
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/json",
			earningId:   "TestEarningId",
			bodyReq: presenter.EarningsBody{
				Amount:      -6.00,
				EarningType: "Dividendos",
				Date:        "2021-10-01",
			},
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Earning: nil,
				Error:   entity.ErrInvalidEarningsAmount.Error(),
			},
		},
		{
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/json",
			earningId:   "TestEarningId",
			bodyReq: presenter.EarningsBody{
				Amount:      6.00,
				EarningType: "Dividendos",
				Date:        date,
			},
			expectedResp: body{
				Code:    200,
				Success: true,
				Message: "Earning updated successfully",
				Earning: &presenter.EarningsApiReturn{
					Id:       "TestEarningId",
					Earning:  6.00,
					Type:     "Dividendos",
					Date:     &dateFormatted,
					Currency: "BRL",
					Asset: &presenter.AssetApiReturn{
						Id:     "TestAssetID",
						Symbol: "TEST3",
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
	earnings := EarningsApi{
		ApplicationLogic: *usecases,
		ApiLogic:         logicApi,
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
	api.Put("/earnings/:id", earnings.UpdateEarningFromUser)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "PUT", "/api/earnings/"+testCase.earningId,
			testCase.contentType, testCase.idToken, testCase.bodyReq)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}
}

func TestApiDeleteEarning(t *testing.T) {
	type body struct {
		Success bool                         `json:"success"`
		Message string                       `json:"message"`
		Error   string                       `json:"error"`
		Code    int                          `json:"code"`
		Earning *presenter.EarningsApiReturn `json:"earning"`
	}

	type test struct {
		idToken      string
		contentType  string
		earningId    string
		expectedResp body
	}

	tests := []test{
		{
			idToken:     "ValidIdTokenWithoutEmailVerification",
			contentType: "application/json",
			earningId:   "TestEarningId",
			expectedResp: body{
				Code:    401,
				Success: false,
				Message: entity.ErrMessageApiAuthentication.Error(),
				Error:   "",
				Earning: nil,
			},
		},
		{
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/json",
			earningId:   "INVALID_ID",
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   errors.New("ERROR: invalid input syntax for type uuid:").Error(),
				Earning: nil,
			},
		},
		{
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/json",
			earningId:   "UNKNOWN_ID",
			expectedResp: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   errors.New("no rows in result set").Error(),
				Earning: nil,
			},
		},
		{
			idToken:     "ValidIdTokenPrivilegeUser",
			contentType: "application/json",
			earningId:   "TestEarningId",
			expectedResp: body{
				Code:    200,
				Success: true,
				Message: "Earning deleted successfully",
				Error:   "",
				Earning: &presenter.EarningsApiReturn{
					Id: "TestEarningId",
					Asset: &presenter.AssetApiReturn{
						Id:     "",
						Symbol: "",
					},
				},
			},
		},
	}

	// Mock UseCases function (Sector Application Logic)
	usecases := usecases.NewMockApplications()
	logicApi := logicApi.NewMockApplication(*usecases)

	// Declare Sector Application Logic
	earnings := EarningsApi{
		ApplicationLogic: *usecases,
		ApiLogic:         logicApi,
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
	api.Delete("/earnings/:id", earnings.DeleteEarningFromUser)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "DELETE", "/api/earnings/"+testCase.earningId,
			testCase.contentType, testCase.idToken, nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
		// if jsonResponse.Earning != nil {
		// 	assert.Equal(t, testCase.expectedResp.Earning.Id,
		// 		jsonResponse.Earning.Id)
		// }
	}
}
