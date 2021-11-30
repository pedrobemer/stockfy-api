package fiberHandlers

import (
	"encoding/json"
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

func TestApiGetAssetType(t *testing.T) {
	type body struct {
		Success   bool                 `json:"success"`
		Message   string               `json:"message"`
		Error     string               `json:"error"`
		Code      int                  `json:"code"`
		AssetType *presenter.AssetType `json:"assetType"`
	}

	type test struct {
		idToken      string
		contentType  string
		path         string
		expectedResp body
	}

	tests := []test{
		{
			idToken:     "ValidIdTokenWithoutEmailVerification",
			contentType: "application/json",
			path:        "type=STOCK&country=BR&ordersResume=true",
			expectedResp: body{
				Code:      401,
				Success:   false,
				Message:   entity.ErrMessageApiAuthentication.Error(),
				Error:     "",
				AssetType: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			path:        "type=STOCK&country=BR&ordersResume=ERROR",
			expectedResp: body{
				Code:      400,
				Success:   false,
				Message:   entity.ErrMessageApiRequest.Error(),
				Error:     entity.ErrInvalidApiQueryWithOrderResume.Error(),
				AssetType: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			path:        "country=BR&ordersResume=false",
			expectedResp: body{
				Code:      400,
				Success:   false,
				Message:   entity.ErrMessageApiRequest.Error(),
				Error:     entity.ErrInvalidApiQueryTypeBlank.Error(),
				AssetType: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			path:        "type=STOCK&ordersResume=false",
			expectedResp: body{
				Code:      400,
				Success:   false,
				Message:   entity.ErrMessageApiRequest.Error(),
				Error:     entity.ErrInvalidApiQueryCountryBlank.Error(),
				AssetType: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			path:        "type=INVALID_ASSET_TYPE&country=BR&ordersResume=true",
			expectedResp: body{
				Code:      400,
				Success:   false,
				Message:   entity.ErrMessageApiRequest.Error(),
				Error:     entity.ErrInvalidAssetTypeName.Error(),
				AssetType: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			path:        "type=STOCK&country=ERROR&ordersResume=false",
			expectedResp: body{
				Code:      400,
				Success:   false,
				Message:   entity.ErrMessageApiRequest.Error(),
				Error:     entity.ErrInvalidCountryCode.Error(),
				AssetType: nil,
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			path:        "type=STOCK&country=US&ordersResume=true",
			expectedResp: body{
				Code:    200,
				Success: true,
				Message: "Asset type returned successfully",
				Error:   "",
				AssetType: &presenter.AssetType{
					Id:      "TestAssetTypeID",
					Type:    "STOCK",
					Name:    "Test Name",
					Country: "US",
					Assets: []presenter.AssetApiReturn{
						{
							Id:         "TestAssetID1",
							Symbol:     "TEST1",
							Preference: "TestPref",
							Fullname:   "Test Name 1",
							Sector: &presenter.Sector{
								Id:   "TestSectorID",
								Name: "Test Sector",
							},
							OrderInfos: &presenter.OrderInfos{
								WeightedAdjPrice:     20.10,
								WeightedAveragePrice: 20.5,
								TotalQuantity:        30,
							},
						},
						{
							Id:         "TestAssetID2",
							Symbol:     "TEST2",
							Preference: "TestPref",
							Fullname:   "Test Name 2",
							Sector: &presenter.Sector{
								Id:   "TestSectorID",
								Name: "Test Sector",
							},
							OrderInfos: &presenter.OrderInfos{
								WeightedAdjPrice:     20.10,
								WeightedAveragePrice: 20.5,
								TotalQuantity:        30,
							},
						},
					},
				},
			},
		},
		{
			idToken:     "ValidIdTokenWithoutPrivilegedUser",
			contentType: "application/json",
			path:        "type=STOCK&country=US&ordersResume=false",
			expectedResp: body{
				Code:    200,
				Success: true,
				Message: "Asset type returned successfully",
				Error:   "",
				AssetType: &presenter.AssetType{
					Id:      "TestAssetTypeID",
					Type:    "STOCK",
					Name:    "Test Name",
					Country: "US",
					Assets: []presenter.AssetApiReturn{
						{
							Id:         "TestAssetID1",
							Symbol:     "TEST1",
							Preference: "TestPref",
							Fullname:   "Test Name 1",
							Sector: &presenter.Sector{
								Id:   "TestSectorID",
								Name: "Test Sector",
							},
						},
						{
							Id:         "TestAssetID2",
							Symbol:     "TEST2",
							Preference: "TestPref",
							Fullname:   "Test Name 2",
							Sector: &presenter.Sector{
								Id:   "TestSectorID",
								Name: "Test Sector",
							},
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
	assetTypes := AssetTypeApi{
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
	api.Get("/asset-types", assetTypes.GetAssetTypes)

	for _, testCase := range tests {
		jsonResponse := body{}
		resp, _ := MockHttpRequest(app, "GET", "/api/asset-types?"+testCase.path,
			testCase.contentType, testCase.idToken, nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &jsonResponse)
		jsonResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResp, jsonResponse)
	}
}
