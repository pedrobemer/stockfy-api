package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAsset(t *testing.T) {
	preference := "ON"

	expectedAssetReturn := Asset{
		Symbol:     "ITUB4",
		Fullname:   "Itau Unibanco Holding SA",
		Preference: &preference,
		Sector:     &Sector{Id: "13a48ajp4"},
		AssetType:  &AssetType{Id: "avn48ak"},
	}

	assetInfo, err := NewAsset("ITUB4", "Itau Unibanco Holding SA", &preference,
		"13a48ajp4", "avn48ak", "STOCK", "BR")

	assert.Nil(t, err)
	assert.Equal(t, expectedAssetReturn.Symbol, assetInfo.Symbol)
	assert.Equal(t, expectedAssetReturn.Fullname, assetInfo.Fullname)
	assert.Equal(t, expectedAssetReturn.Preference, assetInfo.Preference)
	assert.Equal(t, expectedAssetReturn.Sector.Id, assetInfo.Sector.Id)
	assert.Equal(t, expectedAssetReturn.AssetType.Id, assetInfo.AssetType.Id)

}

func TestNewAssetValidation(t *testing.T) {
	type test struct {
		symbol       string
		fullname     string
		preference   *string
		assetTypeId  string
		sectorId     string
		assetType    string
		country      string
		respExpected error
	}

	var preferenceNil *string
	preference := "ON"

	tests := []test{
		{
			symbol:       "ITUB4",
			fullname:     "Itau Unibanco Holding SA",
			preference:   &preference,
			sectorId:     "13a48ajp4",
			assetTypeId:  "avn48ak",
			assetType:    "STOCK",
			country:      "BR",
			respExpected: nil,
		},
		{
			symbol:       "ITUB4",
			fullname:     "Itau Unibanco Holding SA",
			preference:   preferenceNil,
			sectorId:     "13a48ajp4",
			assetTypeId:  "avn48ak",
			assetType:    "STOCK",
			country:      "BR",
			respExpected: ErrInvalidAssetPreferenceUndefined,
		},
		{
			symbol:       "",
			fullname:     "Itau Unibanco Holding SA",
			preference:   &preference,
			sectorId:     "13a48ajp4",
			assetTypeId:  "avn48ak",
			assetType:    "STOCK",
			country:      "BR",
			respExpected: ErrInvalidAssetEntityBlank,
		},
		{
			symbol:       "ITUB4",
			fullname:     "",
			preference:   &preference,
			sectorId:     "13a48ajp4",
			assetTypeId:  "avn48ak",
			assetType:    "STOCK",
			country:      "BR",
			respExpected: ErrInvalidAssetEntityBlank,
		},
		{
			symbol:       "ITUB4",
			fullname:     "Itau Unibanco Holding SA",
			preference:   &preference,
			sectorId:     "",
			assetTypeId:  "avn48ak",
			assetType:    "STOCK",
			country:      "BR",
			respExpected: ErrInvalidAssetEntityBlank,
		},
		{
			symbol:       "ITUB4",
			fullname:     "Itau Unibanco Holding SA",
			preference:   &preference,
			sectorId:     "13a48ajp4",
			assetTypeId:  "",
			assetType:    "STOCK",
			country:      "BR",
			respExpected: ErrInvalidAssetEntityBlank,
		},
		{
			symbol:       "ITUB4",
			fullname:     "Itau Unibanco Holding SA",
			preference:   &preference,
			sectorId:     "13a48ajp4",
			assetTypeId:  "avn48ak",
			assetType:    "AAAAA",
			country:      "BR",
			respExpected: ErrInvalidAssetEntityValues,
		},
		{
			symbol:       "ITUB4",
			fullname:     "Itau Unibanco Holding SA",
			preference:   &preference,
			sectorId:     "13a48ajp4",
			assetTypeId:  "avn48ak",
			assetType:    "ETF",
			country:      "AAAA",
			respExpected: ErrInvalidAssetEntityValues,
		},
	}

	for _, testCase := range tests {
		_, err := NewAsset(testCase.symbol, testCase.fullname,
			testCase.preference, testCase.sectorId, testCase.assetTypeId,
			testCase.assetType, testCase.country)
		assert.Equal(t, testCase.respExpected, err)
	}
}
