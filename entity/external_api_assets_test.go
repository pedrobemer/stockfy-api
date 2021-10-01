package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertAssetLookup(t *testing.T) {

	type test struct {
		symbol               string
		fullname             string
		symbolType           string
		expectedSymbolLookup SymbolLookup
	}

	tests := []test{
		{
			symbol:     "VEA",
			fullname:   "VANGUARD FTSE DEVELOPED ETF",
			symbolType: "ETF",
			expectedSymbolLookup: SymbolLookup{
				Symbol:   "VEA",
				Fullname: "Vanguard FTSE Developed ETF",
				Type:     "ETF",
			},
		},
		{
			symbol:     "ITUB4.SA",
			fullname:   "Itaú Unibanco Holding S.A",
			symbolType: "STOCK",
			expectedSymbolLookup: SymbolLookup{
				Symbol:   "ITUB4",
				Fullname: "Itaú Unibanco Holding S.A",
				Type:     "STOCK",
			},
		},
	}

	for _, testCase := range tests {
		symbolLookup := ConvertAssetLookup(testCase.symbol, testCase.fullname,
			testCase.symbolType)
		assert.Equal(t, testCase.expectedSymbolLookup, symbolLookup)
	}

}

func TestConvertAssetPrice(t *testing.T) {
	type test struct {
		symbol              string
		openPrice           string
		highPrice           string
		lowPrice            string
		currentPrice        string
		prevClosePrice      string
		expectedSymbolPrice SymbolPrice
	}

	tests := []test{
		{
			symbol:         "ITUB4",
			openPrice:      "39.29",
			highPrice:      "39.35",
			lowPrice:       "38.3",
			currentPrice:   "38.5839",
			prevClosePrice: "39.29",
			expectedSymbolPrice: SymbolPrice{
				Symbol:         "ITUB4",
				OpenPrice:      39.29,
				HighPrice:      39.35,
				LowPrice:       38.3,
				CurrentPrice:   38.5839,
				PrevClosePrice: 39.29,
			},
		},
	}

	for _, testCase := range tests {
		symbolPrice := ConvertAssetPrice(testCase.symbol, testCase.openPrice,
			testCase.highPrice, testCase.lowPrice, testCase.currentPrice,
			testCase.prevClosePrice)
		assert.Equal(t, testCase.expectedSymbolPrice, symbolPrice)
	}
}
