package alphaVantage

import (
	"stockfyApi/client"
	"stockfyApi/commonTypes"
	"strings"
)

func VerifySymbolAlpha(symbol string) SymbolLookupAlpha {
	url := "https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=" +
		symbol + "&apikey=KIUG1ZKFZ13BI08F"

	var symbolLookup SymbolLookupAlpha

	client.RequestAndAssignToBody(url, &symbolLookup)

	return symbolLookup
}

func ConvertSymbolLookup(queryResult SymbolLookupInfo) commonTypes.SymbolLookup {
	var symbolLookup commonTypes.SymbolLookup

	symbolLookup.Symbol = strings.ReplaceAll(queryResult["1. symbol"],
		".SA", "")
	symbolLookup.Fullname = queryResult["2. name"]
	symbolLookup.Type = queryResult["3. type"]

	return symbolLookup

}
