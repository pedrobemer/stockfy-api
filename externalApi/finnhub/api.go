package finnhub

import (
	"stockfyApi/client"
)

type FinnhubApi struct {
	Token string
}

func (f *FinnhubApi) VerifySymbol(symbol string) SymbolLookupInfo {
	url := "https://finnhub.io/api/v1/search?q=" + symbol + "&token=" +
		f.Token

	var symbolLookupFinnhub SymbolLookupFinnhub
	var symbolLookupInfo SymbolLookupInfo

	client.RequestAndAssignToBody("GET", url, nil, &symbolLookupFinnhub)

	for _, s := range symbolLookupFinnhub.Result {
		if s.Symbol == symbol {
			symbolLookupInfo = s
		}
	}

	return symbolLookupInfo
}

func (f *FinnhubApi) CompanyProfile2(symbol string) CompanyProfile2 {
	url := "https://finnhub.io/api/v1/stock/profile2?symbol=" + symbol +
		"&token=" + f.Token

	var companyProfile2 CompanyProfile2

	client.RequestAndAssignToBody("GET", url, nil, &companyProfile2)

	return companyProfile2
}

func (f *FinnhubApi) GetPrice(symbol string) SymbolPriceFinnhub {
	url := "https://finnhub.io/api/v1/quote?symbol=" + symbol + "&token=" +
		f.Token

	symbolPrice := SymbolPriceFinnhub{}
	// symbolPrice := commonTypes.SymbolPrice{}

	client.RequestAndAssignToBody("GET", url, nil, &symbolPrice)

	// formatFinhubSymbolPrice(symbolPriceNotFormatted, &symbolPrice, symbol)

	return symbolPrice
}
