package alphaVantage

import (
	"stockfyApi/client"
)

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

func ConvertSymbolLookup(queryResult SymbolLookupInfo) commonTypes.SymbolLookup {
	var symbolLookup commonTypes.SymbolLookup

	symbolLookup.Symbol = strings.ReplaceAll(queryResult.Symbol, ".SAO", "")
	symbolLookup.Fullname = queryResult.Name
	symbolLookup.Type = queryResult.Type

	return symbolLookup

}
