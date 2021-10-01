package asset

import (
	"stockfyApi/entity"
	assettype "stockfyApi/usecases/assetType"
	"testing"

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
