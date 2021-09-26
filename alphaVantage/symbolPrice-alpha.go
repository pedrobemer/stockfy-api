package alphaVantage

import (
	"fmt"
	"stockfyApi/client"
	"stockfyApi/commonTypes"
	"stockfyApi/convertVariables"
	"strings"
)

func (a *AlphaApi) GetPriceAlphaVantage(symbol string) commonTypes.SymbolPrice {
	url := "https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=" +
		symbol + "&apikey=" + a.Token

	var symbolPriceNotFormatted SymbolPriceAlpha
	var symbolPrice commonTypes.SymbolPrice

	client.RequestAndAssignToBody("GET", url, nil, &symbolPriceNotFormatted)
	fmt.Println(symbolPriceNotFormatted)

	formatAlphaVantageSymbolPrice(symbolPriceNotFormatted, &symbolPrice, symbol)

	return symbolPrice
}

func formatAlphaVantageSymbolPrice(unformatted SymbolPriceAlpha,
	formatted *commonTypes.SymbolPrice, symbol string) {

	formatted.Symbol = strings.ReplaceAll(unformatted.GlobalQuote.Symbol,
		".SAO", "")
	formatted.OpenPrice = convertVariables.StringToFloat64(
		unformatted.GlobalQuote.Open)
	formatted.HighPrice = convertVariables.StringToFloat64(
		unformatted.GlobalQuote.High)
	formatted.LowPrice = convertVariables.StringToFloat64(
		unformatted.GlobalQuote.Low)
	formatted.CurrentPrice = convertVariables.StringToFloat64(
		unformatted.GlobalQuote.Price)
	formatted.PrevClosePrice = convertVariables.StringToFloat64(
		unformatted.GlobalQuote.PrevClose)
}
