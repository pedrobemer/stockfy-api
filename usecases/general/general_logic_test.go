package general

import (
	"stockfyApi/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountryValidation(t *testing.T) {
	type test struct {
		country      string
		respExpected error
	}

	tests := []test{
		{
			country:      "BR",
			respExpected: nil,
		},
		{
			country:      "US",
			respExpected: nil,
		},
		{
			country:      "",
			respExpected: nil,
		},
		{
			country:      "DAJDAD",
			respExpected: entity.ErrInvalidCountryCode,
		},
	}

	for _, testCase := range tests {
		err := CountryValidation(testCase.country)
		assert.Equal(t, testCase.respExpected, err)
	}
}

func TestAssetTypeNameValidation(t *testing.T) {
	type test struct {
		name         string
		respExpected error
	}

	tests := []test{
		{
			name:         "ETF",
			respExpected: nil,
		},
		{
			name:         "STOCK",
			respExpected: nil,
		},
		{
			name:         "REIT",
			respExpected: nil,
		},
		{
			name:         "FII",
			respExpected: nil,
		},
		{
			name:         "",
			respExpected: nil,
		},
		{
			name:         "9a99a0dj",
			respExpected: entity.ErrInvalidAssetTypeName,
		},
	}

	for _, testCase := range tests {
		err := AssetTypeNameValidation(testCase.name)
		assert.Equal(t, testCase.respExpected, err)
	}
}
