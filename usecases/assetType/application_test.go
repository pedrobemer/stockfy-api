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
