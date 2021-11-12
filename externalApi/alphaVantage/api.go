package alphaVantage

import (
	"fmt"
	"stockfyApi/client"
	"stockfyApi/entity"
)

type AlphaApi struct {
	Token string
}

func NewAlphaVantageApi(token string) *AlphaApi {
	return &AlphaApi{
		Token: token,
	}
}

func (a *AlphaApi) VerifySymbol2(symbol string) entity.SymbolLookup {
	url := "https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=" +
		symbol + "&apikey=" + a.Token

	var symbolLookupAlpha SymbolLookupAlpha
	var symbolLookupBest SymbolLookupInfo

	client.RequestAndAssignToBody("GET", url, "", nil, &symbolLookupAlpha)

	for _, s := range symbolLookupAlpha.BestMatches {
		if s.MatchScore == "1.0000" {
			symbolLookupBest = s
		}
	}

	symbolLookup := entity.ConvertAssetLookup(symbolLookupBest.Symbol,
		symbolLookupBest.Name, symbolLookupBest.Type)

	return symbolLookup
}

func (a *AlphaApi) GetPrice(symbol string) entity.SymbolPrice {
	url := "https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=" +
		symbol + "&apikey=" + a.Token

	var symbolPriceNotFormatted SymbolPriceAlpha

	client.RequestAndAssignToBody("GET", url, "", nil, &symbolPriceNotFormatted)
	fmt.Println(symbolPriceNotFormatted)

	symbolPrice := entity.ConvertAssetPrice(symbol,
		symbolPriceNotFormatted.GlobalQuote.Open,
		symbolPriceNotFormatted.GlobalQuote.High,
		symbolPriceNotFormatted.GlobalQuote.Low,
		symbolPriceNotFormatted.GlobalQuote.Price,
		symbolPriceNotFormatted.GlobalQuote.LatestDay)

	return symbolPrice
}

func (a *AlphaApi) CompanyOverview(symbol string) map[string]string {
	url := "https://www.alphavantage.co/query?function=OVERVIEW&symbol=" +
		symbol + "&apikey=" + a.Token

	var companyOverview map[string]string

	client.RequestAndAssignToBody("GET", url, "", nil, &companyOverview)

	return companyOverview
}
