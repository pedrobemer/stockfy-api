package alphaVantage

import (
	"fmt"
	"stockfyApi/client"
)

type AlphaApi struct {
	Token string
}

func (a *AlphaApi) VerifySymbol(symbol string) SymbolLookupInfo {
	url := "https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=" +
		symbol + "&apikey=" + a.Token

	var symbolLookupAlpha SymbolLookupAlpha
	var symbolLookup SymbolLookupInfo

	client.RequestAndAssignToBody("GET", url, nil, &symbolLookupAlpha)

	for _, s := range symbolLookupAlpha.BestMatches {
		if s.MatchScore == "1.0000" {
			symbolLookup = s
		}
	}

	return symbolLookup
}

func (a *AlphaApi) GetPrice(symbol string) SymbolPriceAlpha {
	url := "https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=" +
		symbol + "&apikey=" + a.Token

	var symbolPrice SymbolPriceAlpha
	// var symbolPrice commonTypes.SymbolPrice

	client.RequestAndAssignToBody("GET", url, nil, &symbolPrice)
	fmt.Println(symbolPrice)

	// formatAlphaVantageSymbolPrice(symbolPriceNotFormatted, &symbolPrice, symbol)

	return symbolPrice
}

func (a *AlphaApi) CompanyOverview(symbol string) CompanyOverview {
	url := "https://www.alphavantage.co/query?function=OVERVIEW&symbol=" +
		symbol + "&apikey=" + a.Token

	var companyOverview CompanyOverview

	client.RequestAndAssignToBody("GET", url, nil, &companyOverview)

	return companyOverview
}
