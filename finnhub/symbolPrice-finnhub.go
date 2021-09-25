package finnhub

import (
	"stockfyApi/client"
	"stockfyApi/commonTypes"
)

func (f *FinnhubApi) GetPriceFinnhub(symbol string) commonTypes.SymbolPrice {
	url := "https://finnhub.io/api/v1/quote?symbol=" + symbol + "&token=" +
		f.token

	symbolPriceNotFormatted := SymbolPriceFinnhub{}
	symbolPrice := commonTypes.SymbolPrice{}

	client.RequestAndAssignToBody("GET", url, nil, &symbolPriceNotFormatted)

	formatFinhubSymbolPrice(symbolPriceNotFormatted, &symbolPrice, symbol)

	return symbolPrice
}

func formatFinhubSymbolPrice(unformatted SymbolPriceFinnhub,
	formatted *commonTypes.SymbolPrice, symbol string) {
	formatted.Symbol = symbol
	formatted.CurrentPrice = unformatted.C
	formatted.HighPrice = unformatted.H
	formatted.LowPrice = unformatted.L
	formatted.PrevClosePrice = unformatted.PC
	formatted.OpenPrice = unformatted.O
	formatted.MarketCap = unformatted.T
}
