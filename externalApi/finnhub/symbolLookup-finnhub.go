package finnhub

import (
	"stockfyApi/client"
	"stockfyApi/commonTypes"
	"strings"
)

func (f *FinnhubApi) VerifySymbolFinnhub(symbol string) SymbolLookupInfo {
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

func ConvertSymbolLookup(queryResult SymbolLookupInfo) commonTypes.SymbolLookup {

	var symbolLookup commonTypes.SymbolLookup
	fullname_title := strings.Title(strings.ToLower(queryResult.Description))

	for _, s := range strings.Fields(fullname_title) {
		if s == "Sa" || s == "Edp" || s == "Etf" || s == "Ftse" ||
			s == "Msci" || s == "Usa" {
			fullname_title = strings.ReplaceAll(fullname_title, s,
				strings.ToUpper(s))
		}
	}

	symbolLookup.Symbol = strings.ReplaceAll(queryResult.Symbol, ".SA", "")
	symbolLookup.Fullname = fullname_title
	symbolLookup.Type = SymbolTypesFinnhub[queryResult.Type]

	return symbolLookup

}
