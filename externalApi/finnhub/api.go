package finnhub

import (
	"io"
	"stockfyApi/entity"
)

type FinnhubApi struct {
	Token              string
	HttpOutsideRequest func(method string, url string, contentType string,
		bodyReq io.Reader, bodyResp interface{})
}

func NewFinnhubApi(token string, httpClient func(method string,
	url string, contentType string, bodyReq io.Reader,
	bodyResp interface{})) *FinnhubApi {
	return &FinnhubApi{
		Token:              token,
		HttpOutsideRequest: httpClient,
	}
}

func (f *FinnhubApi) VerifySymbol2(symbol string) entity.SymbolLookup {
	url := "https://finnhub.io/api/v1/search?q=" + symbol + "&token=" +
		f.Token

	var symbolLookupFinnhub SymbolLookupFinnhub
	var symbolLookupInfo SymbolLookupInfo

	f.HttpOutsideRequest("GET", url, "", nil, &symbolLookupFinnhub)

	for _, s := range symbolLookupFinnhub.Result {
		if s.Symbol == symbol {
			symbolLookupInfo = s
		}
	}

	symbolLookup := entity.ConvertAssetLookup(symbolLookupInfo.Symbol,
		symbolLookupInfo.Description, symbolLookupInfo.Type)

	return symbolLookup
}

func (f *FinnhubApi) CompanyOverview(symbol string) map[string]string {
	url := "https://finnhub.io/api/v1/stock/profile2?symbol=" + symbol +
		"&token=" + f.Token

	var companyProfile2 CompanyProfile2

	f.HttpOutsideRequest("GET", url, "", nil, &companyProfile2)

	return map[string]string{
		"country":         companyProfile2.Country,
		"currency":        companyProfile2.Currency,
		"exchange":        companyProfile2.Exchange,
		"finnhubIndustry": companyProfile2.FinnhubIndustry,
		"ipo":             "",
		"logo":            companyProfile2.Logo,
		"name":            companyProfile2.Name,
		"phone":           companyProfile2.Phone,
		"ticker":          companyProfile2.Ticker,
		"weburl":          companyProfile2.Weburl,
	}
}

func (f *FinnhubApi) GetPrice(symbol string) entity.SymbolPrice {
	url := "https://finnhub.io/api/v1/quote?symbol=" + symbol + "&token=" +
		f.Token

	symbolPrice := SymbolPriceFinnhub{}

	f.HttpOutsideRequest("GET", url, "", nil, &symbolPrice)

	return entity.SymbolPrice{
		Symbol:         symbol,
		OpenPrice:      symbolPrice.O,
		HighPrice:      symbolPrice.H,
		LowPrice:       symbolPrice.L,
		CurrentPrice:   symbolPrice.C,
		PrevClosePrice: symbolPrice.PC,
	}
}
