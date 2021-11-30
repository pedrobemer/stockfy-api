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

func configureAssetTypesApp(dbpool *pgx.Conn) (fiberHandlers.AssetTypeApi,
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

	assetType := fiberHandlers.AssetTypeApi{
		ApplicationLogic: *applicationLogics,
		LogicApi:         logicApiUseCases,
	}

	return assetType, *applicationLogics
}

func TestFiberHandlersIntegrationTestGetAssetTypes(t *testing.T) {

	type body struct {
		Success   bool                 `json:"success"`
		Message   string               `json:"message"`
		Error     string               `json:"error"`
		Code      int                  `json:"code"`
		AssetType *presenter.AssetType `json:"assetType"`
	}

	type test struct {
		idToken          string
		assetType        string
		country          string
		ordersResume     string
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
				AssetType: nil,
			},
		},
		{
			idToken:      "TestNoAdminID",
			assetType:    "STOCK",
			country:      "BR",
			ordersResume: "ERROR",
			expectedResponse: body{
				Code:      400,
				Success:   false,
				Message:   entity.ErrMessageApiRequest.Error(),
				Error:     entity.ErrInvalidApiQueryWithOrderResume.Error(),
				AssetType: nil,
			},
		},
		{
			idToken:      "TestNoAdminID",
			assetType:    "",
			country:      "BR",
			ordersResume: "false",
			expectedResponse: body{
				Code:      400,
				Success:   false,
				Message:   entity.ErrMessageApiRequest.Error(),
				Error:     entity.ErrInvalidApiQueryTypeBlank.Error(),
				AssetType: nil,
			},
		},
		{
			idToken:      "TestNoAdminID",
			assetType:    "STOCK",
			country:      "",
			ordersResume: "false",
			expectedResponse: body{
				Code:      400,
				Success:   false,
				Message:   entity.ErrMessageApiRequest.Error(),
				Error:     entity.ErrInvalidApiQueryCountryBlank.Error(),
				AssetType: nil,
			},
		},
		{
			idToken:      "TestNoAdminID",
			assetType:    "STOCK",
			country:      "ERROR",
			ordersResume: "false",
			expectedResponse: body{
				Code:      400,
				Success:   false,
				Message:   entity.ErrMessageApiRequest.Error(),
				Error:     entity.ErrInvalidCountryCode.Error(),
				AssetType: nil,
			},
		},
		{
			idToken:      "TestNoAdminID",
			assetType:    "ERROR",
			country:      "US",
			ordersResume: "false",
			expectedResponse: body{
				Code:      400,
				Success:   false,
				Message:   entity.ErrMessageApiRequest.Error(),
				Error:     entity.ErrInvalidAssetTypeName.Error(),
				AssetType: nil,
			},
		},
		{
			idToken:      "TestAdminID",
			assetType:    "STOCK",
			country:      "BR",
			ordersResume: "false",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Asset type returned successfully",
				Error:   "",
				AssetType: &presenter.AssetType{
					Type:    "STOCK",
					Country: "BR",
					Name:    "Ações Brasil",
					Assets: []presenter.AssetApiReturn{
						{
							Symbol:     "ITUB4",
							Preference: "PN",
							Fullname:   "Itau Unibanco Holding S.A",
							Sector: &presenter.Sector{
								Name: "Finances",
							},
						},
						{
							Symbol:     "EGIE3",
							Preference: "ON",
							Fullname:   "Engie Brasil Energia S.A",
							Sector: &presenter.Sector{
								Name: "Utilities",
							},
						},
					},
				},
			},
		},
		{
			idToken:      "TestAdminID",
			assetType:    "STOCK",
			country:      "BR",
			ordersResume: "true",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Asset type returned successfully",
				Error:   "",
				AssetType: &presenter.AssetType{
					Type:    "STOCK",
					Country: "BR",
					Name:    "Ações Brasil",
					Assets: []presenter.AssetApiReturn{
						{
							Symbol:     "ITUB4",
							Preference: "PN",
							Fullname:   "Itau Unibanco Holding S.A",
							Sector: &presenter.Sector{
								Name: "Finances",
							},
						},
						{
							Symbol:     "EGIE3",
							Preference: "ON",
							Fullname:   "Engie Brasil Energia S.A",
							Sector: &presenter.Sector{
								Name: "Utilities",
							},
						},
					},
				},
			},
		},
		{
			idToken:      "TestAdminID",
			assetType:    "STOCK",
			country:      "US",
			ordersResume: "true",
			expectedResponse: body{
				Code:      400,
				Success:   false,
				Message:   entity.ErrMessageApiRequest.Error(),
				Error:     entity.ErrInvalidAssetType.Error(),
				AssetType: nil,
			},
		},
	}

	DBpool := connectDatabase()

	assetTypes, applicationsLogics := configureAssetTypesApp(DBpool)

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
	api.Get("/asset-types", assetTypes.GetAssetTypes)

	for _, testCase := range tests {
		bodyResponse := body{}

		resp, _ := fiberHandlers.MockHttpRequest(app, "GET", "/api/asset-types?"+
			"type="+testCase.assetType+"&country="+testCase.country+
			"&ordersResume="+testCase.ordersResume, "", testCase.idToken, nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &bodyResponse)
		bodyResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResponse.Code, bodyResponse.Code)
		assert.Equal(t, testCase.expectedResponse.Success, bodyResponse.Success)
		assert.Equal(t, testCase.expectedResponse.Message, bodyResponse.Message)
		assert.Equal(t, testCase.expectedResponse.Error, bodyResponse.Error)

		if testCase.expectedResponse.AssetType != nil {
			assert.Equal(t, testCase.expectedResponse.AssetType.Type,
				bodyResponse.AssetType.Type)
			assert.Equal(t, testCase.expectedResponse.AssetType.Country,
				bodyResponse.AssetType.Country)
			assert.Equal(t, testCase.expectedResponse.AssetType.Name,
				bodyResponse.AssetType.Name)

			for _, asset := range bodyResponse.AssetType.Assets {
				searchedAsset := presenter.AssetApiReturn{}
				for _, expAsset := range testCase.expectedResponse.AssetType.Assets {
					if asset.Symbol == expAsset.Symbol {
						searchedAsset = expAsset
					}
				}

				assert.Equal(t, searchedAsset.Symbol, asset.Symbol)
				assert.Equal(t, searchedAsset.Preference, asset.Preference)
				assert.Equal(t, searchedAsset.Fullname, asset.Fullname)
				assert.Equal(t, searchedAsset.Sector.Name, asset.Sector.Name)
				if testCase.ordersResume == "true" {
					assert.NotNil(t, asset.OrderInfos)

				} else {
					assert.Nil(t, asset.OrderInfos)
				}
			}

		} else {
			assert.Nil(t, testCase.expectedResponse.AssetType)
		}

	}

}
