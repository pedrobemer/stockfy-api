package entity

import (
	"strings"
)

func ConvertAssetLookup(symbol string, fullname string,
	symbolType string) SymbolLookup {

	symbolLookup := SymbolLookup{
		Symbol:   strings.ReplaceAll(symbol, ".SAO", ""),
		Fullname: fullname,
		Type:     symbolType,
	}

	return symbolLookup
}

func ConvertAssetlPrice(symbol string, openPrice string, highPrice string,
	lowPrice string, currentPrice string, prevClosePrice string) SymbolPrice {

	symbolPrice := SymbolPrice{
		Symbol:         symbol,
		OpenPrice:      StringToFloat64(openPrice),
		HighPrice:      StringToFloat64(highPrice),
		LowPrice:       StringToFloat64(lowPrice),
		CurrentPrice:   StringToFloat64(currentPrice),
		PrevClosePrice: StringToFloat64(prevClosePrice),
	}

	return symbolPrice
}
