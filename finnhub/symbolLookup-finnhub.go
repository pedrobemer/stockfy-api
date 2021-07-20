package finnhub

import (
	"stockfyApi/client"
	"stockfyApi/commonTypes"
	"strings"
)

func VerifySymbolFinnhub(symbol string) SymbolLookupFinnhub {
	url := "https://finnhub.io/api/v1/search?q=" + symbol + "&token=c2o3062ad3ie71thpra0"

	var symbolLookup SymbolLookupFinnhub

	client.RequestAndAssignToBody(url, &symbolLookup)

	return symbolLookup
}

func ConvertSymbolLookup(queryResult SymbolLookupInfo) commonTypes.SymbolLookup {

	var symbolLookup commonTypes.SymbolLookup

	symbolLookup.Symbol = strings.ReplaceAll(queryResult.Symbol, ".SA", "")
	symbolLookup.Fullname = queryResult.Description
	symbolLookup.Type = SymbolTypesFinnhub[queryResult.Type]

	return symbolLookup

}
