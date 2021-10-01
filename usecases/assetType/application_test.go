package assettype

import (
	"stockfyApi/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createApp() *Application {
	mockedRepo := NewMockRepo()

	assetApp := NewApplication(mockedRepo)

	return assetApp
}

func TestSearch(t *testing.T) {

	astApp := createApp()

	type test struct {
		name                  string
		country               string
		respExpectedAssetType []entity.AssetType
		respExpectedError     error
	}

	expectedAllAssetTypes, _ := NewMockRepo().Search("", "", "")
	expectedSpecificAssetTypes, _ := NewMockRepo().Search("SPECIFIC", "STOCK",
		"US")
	expectedOnlyAssetTypes, _ := NewMockRepo().Search("ONLYTYPE", "STOCK", "")
	expectedOnlyCountryAssetTypes, _ := NewMockRepo().Search("ONLYCOUNTRY", "",
		"US")

	tests := []test{
		{
			name:                  "",
			country:               "",
			respExpectedAssetType: expectedAllAssetTypes,
			respExpectedError:     nil,
		},
		{
			name:                  "STOCK",
			country:               "US",
			respExpectedAssetType: expectedSpecificAssetTypes,
			respExpectedError:     nil,
		},
		{
			name:                  "STOCK",
			country:               "",
			respExpectedAssetType: expectedOnlyAssetTypes,
			respExpectedError:     nil,
		},
		{
			name:                  "",
			country:               "US",
			respExpectedAssetType: expectedOnlyCountryAssetTypes,
			respExpectedError:     nil,
		},
		{
			name:                  "",
			country:               "AAODIASIDJSAO",
			respExpectedAssetType: nil,
			respExpectedError:     entity.ErrInvalidCountryCode,
		},
		{
			name:                  "AODADA",
			country:               "",
			respExpectedAssetType: nil,
			respExpectedError:     entity.ErrInvalidAssetTypeName,
		},
	}

	for _, testCase := range tests {
		assetTypeReturned, err := astApp.SearchAssetType(testCase.name,
			testCase.country)
		assert.Equal(t, testCase.respExpectedError, err)
		assert.Equal(t, testCase.respExpectedAssetType, assetTypeReturned)
	}

}

func TestAssetTypeConversion(t *testing.T) {
	type test struct {
		assetType                  string
		country                    string
		symbol                     string
		expectedAssetTypeConverted string
	}

	tests := []test{
		{
			assetType:                  "ETP",
			country:                    "US",
			symbol:                     "VTI",
			expectedAssetTypeConverted: "ETF",
		},
		{
			assetType:                  "Common Stock",
			country:                    "US",
			symbol:                     "AAPL",
			expectedAssetTypeConverted: "STOCK",
		},
		{
			assetType:                  "ETF",
			country:                    "BR",
			symbol:                     "IVVB11",
			expectedAssetTypeConverted: "ETF",
		},
		{
			assetType:                  "ETF",
			country:                    "BR",
			symbol:                     "KNRI11",
			expectedAssetTypeConverted: "FII",
		},
		{
			assetType:                  "Equity",
			country:                    "BR",
			symbol:                     "FLRY3",
			expectedAssetTypeConverted: "STOCK",
		},
		{
			assetType:                  "REAL ESTATE INVESTMENT TRUSTS",
			country:                    "US",
			symbol:                     "AMT",
			expectedAssetTypeConverted: "REIT",
		},
		{
			assetType:                  "Equity",
			country:                    "US",
			symbol:                     "AAPL",
			expectedAssetTypeConverted: "STOCK",
		},
	}

	assetApp := createApp()

	for _, testCase := range tests {
		assetType := assetApp.AssetTypeConversion(testCase.assetType,
			testCase.country, testCase.symbol)
		assert.Equal(t, testCase.expectedAssetTypeConverted, assetType)
	}
}
