package alphaVantage

import (
	"io"
	"stockfyApi/entity"
	"strings"
)

type AlphaApi struct {
	Token              string
	HttpOutsideRequest func(method string, url string, contentType string,
		bodyReq io.Reader, bodyResp interface{})
}

func NewAlphaVantageApi(token string, httpClient func(method string,
	url string, contentType string, bodyReq io.Reader,
	bodyResp interface{})) *AlphaApi {
	return &AlphaApi{
		Token:              token,
		HttpOutsideRequest: httpClient,
	}
}

func (a *AlphaApi) VerifySymbol2(symbol string) entity.SymbolLookup {
	url := "https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=" +
		symbol + "&apikey=" + a.Token

	var symbolLookupAlpha SymbolLookupAlpha
	var symbolLookupBest SymbolLookupInfo

	a.HttpOutsideRequest("GET", url, "", nil, &symbolLookupAlpha)

	for _, s := range symbolLookupAlpha.BestMatches {
		if s.MatchScore == "1.0000" {
			symbolLookupBest = s
		}
	}

	symbolLookupBest.Symbol = strings.ReplaceAll(symbolLookupBest.Symbol,
		".SAO", ".SA")
	symbolLookup := entity.ConvertAssetLookup(symbolLookupBest.Symbol,
		symbolLookupBest.Name, symbolLookupBest.Type)

	return symbolLookup
}

func (a *AlphaApi) GetPrice(symbol string) entity.SymbolPrice {
	url := "https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=" +
		symbol + "&apikey=" + a.Token

	var symbolPriceNotFormatted SymbolPriceAlpha

	a.HttpOutsideRequest("GET", url, "", nil, &symbolPriceNotFormatted)

	symbolPrice := entity.ConvertAssetPrice(symbol,
		symbolPriceNotFormatted.GlobalQuote.Open,
		symbolPriceNotFormatted.GlobalQuote.High,
		symbolPriceNotFormatted.GlobalQuote.Low,
		symbolPriceNotFormatted.GlobalQuote.Price,
		symbolPriceNotFormatted.GlobalQuote.PrevClose)

	return symbolPrice
}

func (a *AlphaApi) CompanyOverview(symbol string) map[string]string {
	url := "https://www.alphavantage.co/query?function=OVERVIEW&symbol=" +
		symbol + "&apikey=" + a.Token

	var companyOverview map[string]string

	a.HttpOutsideRequest("GET", url, "", nil, &companyOverview)

	return companyOverview
}
