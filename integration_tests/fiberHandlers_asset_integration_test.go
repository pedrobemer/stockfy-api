package integration_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
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
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

func mockDoFuncAlphaVerifySymbol(req *http.Request) (*http.Response, error) {
	var symbol string
	bodyResp := alphaVantage.SymbolLookupAlpha{}

	// Treat body from the request to get the symbol value from the URL query
	urlQuery := strings.Split(req.URL.RawQuery, "&")
	for _, query := range urlQuery {
		queryParams := strings.Split(string(query), "=")

		if queryParams[0] == "keywords" {
			symbol = queryParams[1]
		}
	}

	// If the symbol is invalid then return error, else returns the information
	// with the asset information based on the Alpha Vantage JSON template
	switch symbol {
	case "ITUB4.SA":
		bodyResp = alphaVantage.SymbolLookupAlpha{
			BestMatches: []alphaVantage.SymbolLookupInfo{
				{
					Symbol:      symbol + "O",
					Name:        "Itaú Unibanco Holding S.A",
					Type:        "Equity",
					Region:      "Brazil/Sao Paolo",
					MarketOpen:  "10:00",
					MarketClose: "17:30",
					Timezone:    "UTC-03",
					Currency:    "BRL",
					MatchScore:  "1.0000",
				},
			},
		}
	case "FLRY3.SA":
		bodyResp = alphaVantage.SymbolLookupAlpha{
			BestMatches: []alphaVantage.SymbolLookupInfo{
				{
					Symbol:      symbol + "O",
					Name:        "Fleury S.A",
					Type:        "Equity",
					Region:      "Brazil/Sao Paolo",
					MarketOpen:  "10:00",
					MarketClose: "17:30",
					Timezone:    "UTC-03",
					Currency:    "BRL",
					MatchScore:  "1.0000",
				},
			},
		}
		break
	case "KNRI11.SA":
		bodyResp = alphaVantage.SymbolLookupAlpha{
			BestMatches: []alphaVantage.SymbolLookupInfo{
				{
					Symbol: symbol + "O",
					Name: "Kinea Renda Imobiliária Fundo de " +
						"Investimento Imobiliário",
					Type:        "ETF",
					Region:      "Brazil/Sao Paolo",
					MarketOpen:  "10:00",
					MarketClose: "17:30",
					Timezone:    "UTC-03",
					Currency:    "BRL",
					MatchScore:  "1.0000",
				},
			},
		}
		break
	case "IVVB11.SA":
		bodyResp = alphaVantage.SymbolLookupAlpha{
			BestMatches: []alphaVantage.SymbolLookupInfo{
				{
					Symbol: symbol + "O",
					Name: "iShares S&P 500 Fundo de Investimento - " +
						"Investimento No Exterior",
					Type:        "ETF",
					Region:      "Brazil/Sao Paolo",
					MarketOpen:  "10:00",
					MarketClose: "17:30",
					Timezone:    "UTC-03",
					Currency:    "BRL",
					MatchScore:  "1.0000",
				},
			},
		}
		break
	default:
		bodyResp = alphaVantage.SymbolLookupAlpha{}
	}

	bodyByte, _ := json.Marshal(bodyResp)

	respHeader := http.Header{
		"Content-Type": {"application/json"},
	}
	return &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     respHeader,
		Body:       ioutil.NopCloser(bytes.NewReader(bodyByte)),
		Request:    req,
	}, nil
}

func mockDoFuncFinnhubVerifySymbol(req *http.Request) (*http.Response, error) {

	var symbol string

	urlPath := strings.Split(req.URL.Path, "/")
	reqType := urlPath[len(urlPath)-1]

	if reqType == "search" {
		bodyResp := finnhub.SymbolLookupFinnhub{}

		// Treat body from the request to get the symbol value from the URL query
		urlQuery := strings.Split(req.URL.RawQuery, "&")
		for _, query := range urlQuery {
			queryParams := strings.Split(string(query), "=")

			if queryParams[0] == "q" {
				symbol = queryParams[1]
			}
		}

		// If the symbol is invalid then return error, else returns the information
		// with the asset information based on the Alpha Vantage JSON template
		switch symbol {
		case "AAPL":
			bodyResp = finnhub.SymbolLookupFinnhub{
				Count: 2,
				Result: []finnhub.SymbolLookupInfo{
					{
						Symbol:        symbol,
						DisplaySymbol: symbol,
						Type:          "Common Stock",
						Description:   "APPLE INC",
					},
					{
						Symbol:        symbol + ".MX",
						DisplaySymbol: symbol + ".MX",
						Type:          "Common Stock",
						Description:   "APPLE INC",
					},
				},
			}
			break
		case "VTI":
			bodyResp = finnhub.SymbolLookupFinnhub{
				Count: 2,
				Result: []finnhub.SymbolLookupInfo{
					{
						Symbol:        symbol,
						DisplaySymbol: symbol,
						Type:          "ETP",
						Description:   "VANGUARD TOTAL STOCK MKT ETF",
					},
					{
						Symbol:        symbol + ".MX",
						DisplaySymbol: symbol + ".MX",
						Type:          "ETP",
						Description:   "VANGUARD TOTAL STOCK MKT ETF",
					},
				},
			}
			break
		case "AMT":
			bodyResp = finnhub.SymbolLookupFinnhub{
				Count: 2,
				Result: []finnhub.SymbolLookupInfo{
					{
						Symbol:        symbol,
						DisplaySymbol: symbol,
						Type:          "REIT",
						Description:   "AMERICAN TOWER CORP",
					},
					{
						Symbol:        symbol + ".MX",
						DisplaySymbol: symbol + ".MX",
						Type:          "ETP",
						Description:   "AMERICAN TOWER CORP",
					},
				},
			}
			break
		default:
			bodyResp = finnhub.SymbolLookupFinnhub{}
		}

		bodyByte, _ := json.Marshal(bodyResp)

		respHeader := http.Header{
			"Content-Type": {"application/json"},
		}
		return &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     respHeader,
			Body:       ioutil.NopCloser(bytes.NewReader(bodyByte)),
			Request:    req,
		}, nil
	} else if reqType == "profile2" {
		bodyResp := finnhub.CompanyProfile2{}

		// Treat body from the request to get the symbol value from the URL query
		urlQuery := strings.Split(req.URL.RawQuery, "&")
		for _, query := range urlQuery {
			queryParams := strings.Split(string(query), "=")

			if queryParams[0] == "symbol" {
				symbol = queryParams[1]
			}
		}

		// If the symbol is invalid then return error, else returns the information
		// with the asset information based on the Alpha Vantage JSON template
		switch symbol {
		case "FLRY3.SA":
			bodyResp = finnhub.CompanyProfile2{
				Country:              "BR",
				Currency:             "BRL",
				Exchange:             "Sao Paolo",
				FinnhubIndustry:      "Health Care",
				Ipo:                  "2009-12-16",
				Logo:                 "https://finnhub.io/api/logo?symbol=FLRY3.SA",
				MarketCapitalization: 6194.829,
				Name:                 "Fleury SA",
				Phone:                "551150351986.0",
				ShareOutstanding:     316.968763,
				Ticker:               "FLRY3.SA",
				Weburl:               "http://www.fleury.com.br/",
			}
		case "AAPL":
			bodyResp = finnhub.CompanyProfile2{
				Country:              "US",
				Currency:             "USD",
				Exchange:             "NASDAQ NMS - GLOBAL MARKET",
				FinnhubIndustry:      "Technology",
				Ipo:                  "1980-12-12",
				Logo:                 "https://finnhub.io/api/logo?symbol=AAPL",
				MarketCapitalization: 2634047,
				Name:                 "Apple Inc",
				Phone:                "14089961010.0",
				ShareOutstanding:     16426.79,
				Ticker:               "AAPL",
				Weburl:               "https://www.apple.com/",
			}
			break
		case "AMT":
			bodyResp = finnhub.CompanyProfile2{
				Country:              "US",
				Currency:             "USD",
				Exchange:             "NEW YORK STOCK EXCHANGE, INC.",
				FinnhubIndustry:      "Real Estate",
				Ipo:                  "1998-06-05",
				Logo:                 "https://finnhub.io/api/logo?symbol=AMT",
				MarketCapitalization: 118853.9,
				Name:                 "American Tower Corp",
				Phone:                "16173757500.0",
				ShareOutstanding:     444.33,
				Ticker:               "AMT",
				Weburl:               "http://www.americantower.com/",
			}
			break
		default:
			bodyResp = finnhub.CompanyProfile2{}
		}

		bodyByte, _ := json.Marshal(bodyResp)

		respHeader := http.Header{
			"Content-Type": {"application/json"},
		}
		return &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     respHeader,
			Body:       ioutil.NopCloser(bytes.NewReader(bodyByte)),
			Request:    req,
		}, nil
	}

	return &http.Response{}, nil
}

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
