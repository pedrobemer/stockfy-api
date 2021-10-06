package asset

import (
	"stockfyApi/entity"
	assettype "stockfyApi/usecases/assetType"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {

	preference := "ON"
	assetType := assettype.AssetType{
		Id:      "50vjfnsa",
		Type:    "STOCK",
		Country: "BR",
	}

	expectedAssetCreated := entity.Asset{
		Id:         "a38a9jkrh40a",
		Symbol:     "ITUB4",
		Preference: &preference,
		Fullname:   "Itau Unibanco Holding SA",
	}

	mockedRepo := NewMockRepo()

	assetApp := NewApplication(mockedRepo)

	assetCreated, err := assetApp.CreateAsset("ITUB4", "Itau Unibanco Holding SA",
		&preference, "a40vn4", assetType)

	assert.Nil(t, err)
	assert.Equal(t, expectedAssetCreated, assetCreated)

}

func TestSearchAsset(t *testing.T) {
	type test struct {
		symbol        string
		expectedAsset *entity.Asset
		expectedError error
	}

	assetType := entity.AssetType{
		Id:      "28ccf27a-ed8b-11eb-9a03-0242ac130003",
		Type:    "STOCK",
		Name:    "Ações Brasil",
		Country: "BR",
	}
	preference := "ON"

	sectorInfo := entity.Sector{
		Id:   "83ae92f8-ed8b-11eb-9a03-0242ac130003",
		Name: "Finance",
	}

	tests := []test{
		{
			symbol: "ITUB4",
			expectedAsset: &entity.Asset{
				Id:         "0a52d206-ed8b-11eb-9a03-0242ac130003",
				Symbol:     "ITUB4",
				Preference: &preference,
				Fullname:   "Itau Unibanco Holding SA",
				AssetType:  &assetType,
				Sector:     &sectorInfo,
			},
			expectedError: nil,
		},
		{
			symbol:        "Invalid",
			expectedAsset: nil,
			expectedError: entity.ErrInvalidSearchAssetName,
		},
	}

	mockedRepo := NewMockRepo()
	assetApp := NewApplication(mockedRepo)

	for _, testCase := range tests {
		searchedAsset, err := assetApp.SearchAsset(testCase.symbol)
		assert.Equal(t, testCase.expectedAsset, searchedAsset)
		assert.Equal(t, testCase.expectedError, err)
	}
}

func TestSearchAssetByUser(t *testing.T) {
	type test struct {
		symbol        string
		userUid       string
		withInfo      bool
		onlyInfo      bool
		bypassInfo    bool
		expectedAsset *entity.Asset
		expectedError error
	}

	tr, _ := time.Parse("2021-07-05", "2021-07-21")
	tr2, _ := time.Parse("2021-07-05", "2020-04-02")

	assetType := entity.AssetType{
		Id:      "28ccf27a-ed8b-11eb-9a03-0242ac130003",
		Type:    "STOCK",
		Name:    "Ações Brasil",
		Country: "BR",
	}

	preference := "ON"

	brokerageInfo := entity.Brokerage{
		Id:      "55555555-ed8b-11eb-9a03-0242ac130003",
		Name:    "Clear",
		Country: "BR",
	}

	orderList := []entity.Order{
		{
			Id:        "44444444-ed8b-11eb-9a03-0242ac130003",
			Quantity:  20,
			Price:     39.93,
			Currency:  "BRL",
			OrderType: "buy",
			Date:      tr,
			Brokerage: &brokerageInfo,
		},
		{
			Id:        "yeid847e-ed8b-11eb-9a03-0242ac130003",
			Quantity:  5,
			Price:     27.13,
			Currency:  "BRL",
			OrderType: "buy",
			Date:      tr2,
			Brokerage: &brokerageInfo,
		},
	}

	sectorInfo := entity.Sector{
		Id:   "83ae92f8-ed8b-11eb-9a03-0242ac130003",
		Name: "Finance",
	}

	ordersInfo := entity.OrderInfos{
		TotalQuantity:        25,
		WeightedAdjPrice:     37.37,
		WeightedAveragePrice: 37.37,
	}

	tests := []test{
		{
			symbol:     "ITUB4",
			userUid:    "TestID",
			withInfo:   false,
			onlyInfo:   false,
			bypassInfo: false,
			expectedAsset: &entity.Asset{
				Id:         "0a52d206-ed8b-11eb-9a03-0242ac130003",
				Symbol:     "ITUB4",
				Preference: &preference,
				Fullname:   "Itau Unibanco Holding SA",
				AssetType:  &assetType,
				Sector:     &sectorInfo,
				OrdersList: orderList,
			},
			expectedError: nil,
		},
		{
			symbol:     "ITUB4",
			userUid:    "TestID",
			withInfo:   true,
			onlyInfo:   false,
			bypassInfo: false,
			expectedAsset: &entity.Asset{
				Id:         "0a52d206-ed8b-11eb-9a03-0242ac130003",
				Symbol:     "ITUB4",
				Preference: &preference,
				Fullname:   "Itau Unibanco Holding SA",
				AssetType:  &assetType,
				Sector:     &sectorInfo,
				OrdersList: orderList,
				OrderInfo:  &ordersInfo,
			},
			expectedError: nil,
		},
		{
			symbol:     "ITUB4",
			userUid:    "TestID",
			withInfo:   false,
			onlyInfo:   true,
			bypassInfo: false,
			expectedAsset: &entity.Asset{
				Id:         "0a52d206-ed8b-11eb-9a03-0242ac130003",
				Symbol:     "ITUB4",
				Preference: &preference,
				Fullname:   "Itau Unibanco Holding SA",
				AssetType:  &assetType,
				Sector:     &sectorInfo,
				OrderInfo:  &ordersInfo,
			},
			expectedError: nil,
		},
		{
			symbol:        "Invalid",
			userUid:       "TestID",
			withInfo:      false,
			onlyInfo:      true,
			bypassInfo:    false,
			expectedAsset: nil,
			expectedError: entity.ErrInvalidSearchAssetName,
		},
	}

	mockedRepo := NewMockRepo()
	assetApp := NewApplication(mockedRepo)

	for _, testCase := range tests {
		searchedAsset, err := assetApp.SearchAssetByUser(testCase.symbol,
			testCase.userUid, testCase.withInfo, testCase.onlyInfo,
			testCase.bypassInfo)
		assert.Equal(t, testCase.expectedAsset, searchedAsset)
		assert.Equal(t, testCase.expectedError, err)
	}

}

func TestSearchAssetPerAssetType(t *testing.T) {
	type test struct {
		assetType      string
		country        string
		userUid        string
		withOrdersInfo bool
		expectedReturn *entity.AssetType
		expectedError  error
	}

	searchedAssetType := []entity.AssetType{
		{
			Id:      "6582b653-eb19-465b-892a-d6d74d61932c",
			Type:    "ETF",
			Name:    "Asset" + "ETF",
			Country: "US",
			Assets: []entity.Asset{
				{
					Id:         "1f1c2ad6-f16c-4659-826c-e5fd328461f7",
					Preference: nil,
					Fullname:   "ISHARES S&P SMALL-CAP 600 VA",
					Symbol:     "IJS",
					Sector: &entity.Sector{
						Id:   "17d1937c-4e47-4c66-994b-1409c8526cea",
						Name: "Blend",
					},
				},
				{
					Id:         "9ed5cfdc-e4d2-4780-a962-29bb2d11716c",
					Preference: nil,
					Fullname:   "Vanguard FTSE Emerging Market ETF",
					Symbol:     "VWO",
					Sector: &entity.Sector{
						Id:   "17d1937c-4e47-4c66-994b-1409c8526cea",
						Name: "Blend",
					},
				},
			},
		},
	}

	orderInfos := entity.OrderInfos{
		TotalQuantity:        20.09,
		WeightedAdjPrice:     81.56562966650074,
		WeightedAveragePrice: 81.56562966650074,
	}

	searchedAssetTypeWithInfo := []entity.AssetType{
		{
			Id:      searchedAssetType[0].Id,
			Type:    searchedAssetType[0].Type,
			Name:    searchedAssetType[0].Name,
			Country: searchedAssetType[0].Country,
			Assets: []entity.Asset{
				{
					Id:         searchedAssetType[0].Assets[0].Id,
					Preference: searchedAssetType[0].Assets[0].Preference,
					Fullname:   searchedAssetType[0].Assets[0].Fullname,
					Symbol:     searchedAssetType[0].Assets[0].Symbol,
					Sector: &entity.Sector{
						Id:   "17d1937c-4e47-4c66-994b-1409c8526cea",
						Name: "Blend",
					},
					OrderInfo: &orderInfos,
				},
				{
					Id:         searchedAssetType[0].Assets[1].Id,
					Preference: searchedAssetType[0].Assets[1].Preference,
					Fullname:   searchedAssetType[0].Assets[1].Fullname,
					Symbol:     searchedAssetType[0].Assets[1].Symbol,
					Sector: &entity.Sector{
						Id:   "17d1937c-4e47-4c66-994b-1409c8526cea",
						Name: "Blend",
					},
					OrderInfo: &orderInfos,
				},
			},
		},
	}

	tests := []test{
		{
			assetType:      "ETF",
			country:        "US",
			userUid:        "TestID",
			withOrdersInfo: false,
			expectedReturn: &searchedAssetType[0],
			expectedError:  nil,
		},
		{
			assetType:      "ETF",
			country:        "US",
			userUid:        "TestID",
			withOrdersInfo: true,
			expectedReturn: &searchedAssetTypeWithInfo[0],
			expectedError:  nil,
		},
		{
			assetType:      "ETF",
			country:        "BR",
			userUid:        "No Asset",
			withOrdersInfo: true,
			expectedReturn: nil,
			expectedError:  entity.ErrInvalidAssetType,
		},
	}

	mockedRepo := NewMockRepo()
	assetApp := NewApplication(mockedRepo)

	for _, testCase := range tests {
		assetTypeReturn, err := assetApp.SearchAssetPerAssetType(testCase.assetType,
			testCase.country, testCase.userUid, testCase.withOrdersInfo)

		if assetTypeReturn == nil {
			assert.Equal(t, testCase.expectedReturn, assetTypeReturn)
			assert.Equal(t, testCase.expectedError, err)
		} else {
			assert.Equal(t, testCase.expectedReturn.Type, assetTypeReturn.Type)
			assert.Equal(t, testCase.expectedReturn.Country, assetTypeReturn.Country)
			assert.Equal(t, testCase.expectedReturn.Name, assetTypeReturn.Name)
			assert.Equal(t, &testCase.expectedReturn.Assets[0].Sector,
				&assetTypeReturn.Assets[0].Sector)
			assert.Equal(t, testCase.expectedReturn.Assets[0].OrderInfo,
				assetTypeReturn.Assets[0].OrderInfo)
			assert.Equal(t, testCase.expectedError, err)
		}

	}
}

func TestAssetPreferenceType(t *testing.T) {
	type test struct {
		symbol             string
		country            string
		assetType          string
		expectedPreference string
	}

	tests := []test{
		{
			symbol:             "ITUB3",
			country:            "BR",
			assetType:          "STOCK",
			expectedPreference: "ON",
		},
		{
			symbol:             "ITUB4",
			country:            "BR",
			assetType:          "STOCK",
			expectedPreference: "PN",
		},
		{
			symbol:             "TAEE11",
			country:            "BR",
			assetType:          "STOCK",
			expectedPreference: "UNIT",
		},
		{
			symbol:             "AAPL",
			country:            "US",
			assetType:          "STOCK",
			expectedPreference: "",
		},
	}

	mockedRepo := NewMockRepo()

	assetApp := NewApplication(mockedRepo)

	for _, testCase := range tests {
		preference := assetApp.AssetPreferenceType(testCase.symbol,
			testCase.country, testCase.assetType)
		assert.Equal(t, testCase.expectedPreference, preference)
	}
}

func TestAssetVerificationExistence(t *testing.T) {
	type test struct {
		symbol               string
		country              string
		expectedSymbolLookup *entity.SymbolLookup
		expectedError        error
	}

	tests := []test{
		{
			symbol:  "ITUB4",
			country: "BR",
			expectedSymbolLookup: &entity.SymbolLookup{
				Fullname: "Itau Unibanco Holding SA",
				Symbol:   "ITUB4",
				Type:     "STOCK",
			},
			expectedError: nil,
		},
		{
			symbol:               "AAJRI",
			country:              "US",
			expectedSymbolLookup: nil,
			expectedError:        entity.ErrInvalidAssetSymbol,
		},
		{
			symbol:               "AAPL",
			country:              "BAU",
			expectedSymbolLookup: nil,
			expectedError:        entity.ErrInvalidCountryCode,
		},
	}

	mockedDb := NewMockRepo()
	extApiMocked := NewExternalApi()
	assetApp := NewApplication(mockedDb)

	for _, testCase := range tests {
		symbolLookup, err := assetApp.AssetVerificationExistence(testCase.symbol,
			testCase.country, extApiMocked)
		assert.Equal(t, testCase.expectedSymbolLookup, symbolLookup)
		assert.Equal(t, testCase.expectedError, err)
	}

}

func TestAssetVerificationSector(t *testing.T) {
	type test struct {
		assetType      string
		symbol         string
		country        string
		expectedSector string
	}

	tests := []test{
		{
			assetType:      "STOCK",
			symbol:         "BBDC3",
			country:        "BR",
			expectedSector: "Banking",
		},
		{
			assetType:      "ETF",
			symbol:         "IVVB11",
			country:        "BR",
			expectedSector: "Blend",
		},
		{
			assetType:      "STOCK",
			symbol:         "AAPL",
			country:        "US",
			expectedSector: "Banking",
		},
		{
			assetType:      "FII",
			symbol:         "KNRI11",
			country:        "BR",
			expectedSector: "Real Estate",
		},
		{
			assetType:      "REIT",
			symbol:         "AMT",
			country:        "US",
			expectedSector: "Real Estate",
		},
	}

	mockedDb := NewMockRepo()
	extApiMocked := NewExternalApi()
	assetApp := NewApplication(mockedDb)

	for _, testCase := range tests {
		sectorName := assetApp.AssetVerificationSector(testCase.assetType,
			testCase.symbol, testCase.country, extApiMocked)
		assert.Equal(t, testCase.expectedSector, sectorName)
	}

}

func TestAssetVerificationPrice(t *testing.T) {
	type test struct {
		symbol              string
		country             string
		expectedSymbolPrice *entity.SymbolPrice
		expectedError       error
	}

	tests := []test{
		{
			symbol:  "ITUB3",
			country: "BR",
			expectedSymbolPrice: &entity.SymbolPrice{
				Symbol:         "ITUB3",
				CurrentPrice:   29.93,
				HighPrice:      31.00,
				LowPrice:       29.56,
				OpenPrice:      30.99,
				PrevClosePrice: 30.99,
				MarketCap:      1478481948,
			},
			expectedError: nil,
		},
		{
			symbol:              "AAAPDK",
			country:             "US",
			expectedSymbolPrice: nil,
			expectedError:       entity.ErrInvalidAssetSymbol,
		},
		{
			symbol:              "",
			country:             "BR",
			expectedSymbolPrice: nil,
			expectedError:       entity.ErrInvalidAssetSymbol,
		},
		{
			symbol:              "ITUB4",
			country:             "AOS",
			expectedSymbolPrice: nil,
			expectedError:       entity.ErrInvalidCountryCode,
		},
	}

	mockedDb := NewMockRepo()
	extApiMocked := NewExternalApi()
	assetApp := NewApplication(mockedDb)

	for _, testCase := range tests {
		symbolPrice, err := assetApp.AssetVerificationPrice(testCase.symbol,
			testCase.country, extApiMocked)
		assert.Equal(t, testCase.expectedSymbolPrice, symbolPrice)
		assert.Equal(t, testCase.expectedError, err)
	}
}
