package integration_tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"stockfyApi/externalApi/alphaVantage"
	"stockfyApi/externalApi/finnhub"
	"strings"
)

func mockDoFuncAlphaVerifySymbol(req *http.Request) (*http.Response, error) {
	var symbol, funcType string

	// Treat body from the request to get the symbol value from the URL query
	urlQuery := strings.Split(req.URL.RawQuery, "&")
	for _, query := range urlQuery {
		queryParams := strings.Split(string(query), "=")

		switch queryParams[0] {
		case "function":
			funcType = queryParams[1]
			break
		case "keywords":
			symbol = queryParams[1]
		case "symbol":
			symbol = queryParams[1]
		}

	}

	if funcType == "SYMBOL_SEARCH" {
		bodyResp := alphaVantage.SymbolLookupAlpha{}

		switch symbol {
		case "ITUB4.SA":
			bodyResp = alphaVantage.SymbolLookupAlpha{
				BestMatches: []alphaVantage.SymbolLookupInfo{
					{
						Symbol:      symbol + "O",
						Name:        "Itaú Unibanco Holding S.A",
						Type:        "Equity",
						Region:      "Brazil/Sao Paolo",
						MarketOpen:  "10:00",
						MarketClose: "17:30",
						Timezone:    "UTC-03",
						Currency:    "BRL",
						MatchScore:  "1.0000",
					},
				},
			}
		case "FLRY3.SA":
			bodyResp = alphaVantage.SymbolLookupAlpha{
				BestMatches: []alphaVantage.SymbolLookupInfo{
					{
						Symbol:      symbol + "O",
						Name:        "Fleury S.A",
						Type:        "Equity",
						Region:      "Brazil/Sao Paolo",
						MarketOpen:  "10:00",
						MarketClose: "17:30",
						Timezone:    "UTC-03",
						Currency:    "BRL",
						MatchScore:  "1.0000",
					},
				},
			}
			break
		case "KNRI11.SA":
			bodyResp = alphaVantage.SymbolLookupAlpha{
				BestMatches: []alphaVantage.SymbolLookupInfo{
					{
						Symbol: symbol + "O",
						Name: "Kinea Renda Imobiliária Fundo de " +
							"Investimento Imobiliário",
						Type:        "ETF",
						Region:      "Brazil/Sao Paolo",
						MarketOpen:  "10:00",
						MarketClose: "17:30",
						Timezone:    "UTC-03",
						Currency:    "BRL",
						MatchScore:  "1.0000",
					},
				},
			}
			break
		case "IVVB11.SA":
			bodyResp = alphaVantage.SymbolLookupAlpha{
				BestMatches: []alphaVantage.SymbolLookupInfo{
					{
						Symbol: symbol + "O",
						Name: "iShares S&P 500 Fundo de Investimento - " +
							"Investimento No Exterior",
						Type:        "ETF",
						Region:      "Brazil/Sao Paolo",
						MarketOpen:  "10:00",
						MarketClose: "17:30",
						Timezone:    "UTC-03",
						Currency:    "BRL",
						MatchScore:  "1.0000",
					},
				},
			}
			break
		default:
			bodyResp = alphaVantage.SymbolLookupAlpha{}
		}

		bodyByte, _ := json.Marshal(bodyResp)

		respHeader := http.Header{
			"Content-Type": {"application/json"},
		}

		return &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     respHeader,
			Body:       ioutil.NopCloser(bytes.NewReader(bodyByte)),
			Request:    req,
		}, nil
	} else if funcType == "GLOBAL_QUOTE" {
		var bodyResp alphaVantage.SymbolPriceAlpha

		switch symbol {
		case "ITUB4.SA":
			bodyResp = alphaVantage.SymbolPriceAlpha{
				GlobalQuote: alphaVantage.SymbolPriceInfo{
					Symbol:        "ITUB4.SA",
					Open:          "22.44",
					High:          "23.15",
					Low:           "22.44",
					Price:         "23.06",
					Volume:        "10345",
					LatestDay:     "22.45",
					PrevClose:     "22.44",
					Change:        "1.39",
					ChangePercent: "2.09%",
				},
			}
			break
		case "FLRY3.SA":
			bodyResp = alphaVantage.SymbolPriceAlpha{
				GlobalQuote: alphaVantage.SymbolPriceInfo{
					Symbol:        "FLRY3.SA",
					Open:          "18.00",
					High:          "18.03",
					Low:           "17.16",
					Price:         "17.93",
					Volume:        "10345",
					LatestDay:     "18.01",
					PrevClose:     "18.01",
					Change:        "-0.01",
					ChangePercent: "-0.01%",
				},
			}
			break
		case "KNRI11.SA":
			bodyResp = alphaVantage.SymbolPriceAlpha{
				GlobalQuote: alphaVantage.SymbolPriceInfo{
					Symbol:        "KNRI11.SA",
					Open:          "128.20",
					High:          "129.1",
					Low:           "127.95",
					Price:         "128.5",
					Volume:        "10345",
					LatestDay:     "128.18",
					PrevClose:     "128.18",
					Change:        "0.32",
					ChangePercent: "0.25%",
				},
			}
			break
		case "IVVB11.SA":

			bodyResp = alphaVantage.SymbolPriceAlpha{
				GlobalQuote: alphaVantage.SymbolPriceInfo{
					Symbol:        "IVVB11.SA",
					Open:          "281.00",
					High:          "282.5",
					Low:           "277.75",
					Price:         "281.00",
					Volume:        "10345",
					LatestDay:     "281.00",
					PrevClose:     "281.00",
					Change:        "0.00",
					ChangePercent: "0.00%",
				},
			}
			break
		default:
			bodyResp = alphaVantage.SymbolPriceAlpha{}
		}

		bodyByte, _ := json.Marshal(bodyResp)

		respHeader := http.Header{
			"Content-Type": {"application/json"},
		}
		// If the symbol is invalid then return error, else returns the information
		// with the asset information based on the Alpha Vantage JSON template
		return &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     respHeader,
			Body:       ioutil.NopCloser(bytes.NewReader(bodyByte)),
			Request:    req,
		}, nil
	}

	return &http.Response{}, nil
}

func mockDoFuncFinnhubVerifySymbol(req *http.Request) (*http.Response, error) {

	var symbol string

	urlPath := strings.Split(req.URL.Path, "/")
	reqType := urlPath[len(urlPath)-1]

	if reqType == "search" {
		bodyResp := finnhub.SymbolLookupFinnhub{}

		// Treat body from the request to get the symbol value from the URL query
		urlQuery := strings.Split(req.URL.RawQuery, "&")
		for _, query := range urlQuery {
			queryParams := strings.Split(string(query), "=")

			if queryParams[0] == "q" {
				symbol = queryParams[1]
			}
		}

		// If the symbol is invalid then return error, else returns the information
		// with the asset information based on the Alpha Vantage JSON template
		switch symbol {
		case "AAPL":
			bodyResp = finnhub.SymbolLookupFinnhub{
				Count: 2,
				Result: []finnhub.SymbolLookupInfo{
					{
						Symbol:        symbol,
						DisplaySymbol: symbol,
						Type:          "Common Stock",
						Description:   "APPLE INC",
					},
					{
						Symbol:        symbol + ".MX",
						DisplaySymbol: symbol + ".MX",
						Type:          "Common Stock",
						Description:   "APPLE INC",
					},
				},
			}
			break
		case "VTI":
			bodyResp = finnhub.SymbolLookupFinnhub{
				Count: 2,
				Result: []finnhub.SymbolLookupInfo{
					{
						Symbol:        symbol,
						DisplaySymbol: symbol,
						Type:          "ETP",
						Description:   "VANGUARD TOTAL STOCK MKT ETF",
					},
					{
						Symbol:        symbol + ".MX",
						DisplaySymbol: symbol + ".MX",
						Type:          "ETP",
						Description:   "VANGUARD TOTAL STOCK MKT ETF",
					},
				},
			}
			break
		case "AMT":
			bodyResp = finnhub.SymbolLookupFinnhub{
				Count: 2,
				Result: []finnhub.SymbolLookupInfo{
					{
						Symbol:        symbol,
						DisplaySymbol: symbol,
						Type:          "REIT",
						Description:   "AMERICAN TOWER CORP",
					},
					{
						Symbol:        symbol + ".MX",
						DisplaySymbol: symbol + ".MX",
						Type:          "ETP",
						Description:   "AMERICAN TOWER CORP",
					},
				},
			}
			break
		default:
			bodyResp = finnhub.SymbolLookupFinnhub{}
		}

		bodyByte, _ := json.Marshal(bodyResp)

		respHeader := http.Header{
			"Content-Type": {"application/json"},
		}
		return &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     respHeader,
			Body:       ioutil.NopCloser(bytes.NewReader(bodyByte)),
			Request:    req,
		}, nil
	} else if reqType == "profile2" {
		bodyResp := finnhub.CompanyProfile2{}

		// Treat body from the request to get the symbol value from the URL query
		urlQuery := strings.Split(req.URL.RawQuery, "&")
		for _, query := range urlQuery {
			queryParams := strings.Split(string(query), "=")

			if queryParams[0] == "symbol" {
				symbol = queryParams[1]
			}
		}

		// If the symbol is invalid then return error, else returns the information
		// with the asset information based on the Alpha Vantage JSON template
		switch symbol {
		case "FLRY3.SA":
			bodyResp = finnhub.CompanyProfile2{
				Country:              "BR",
				Currency:             "BRL",
				Exchange:             "Sao Paolo",
				FinnhubIndustry:      "Health Care",
				Ipo:                  "2009-12-16",
				Logo:                 "https://finnhub.io/api/logo?symbol=FLRY3.SA",
				MarketCapitalization: 6194.829,
				Name:                 "Fleury SA",
				Phone:                "551150351986.0",
				ShareOutstanding:     316.968763,
				Ticker:               "FLRY3.SA",
				Weburl:               "http://www.fleury.com.br/",
			}
		case "AAPL":
			bodyResp = finnhub.CompanyProfile2{
				Country:              "US",
				Currency:             "USD",
				Exchange:             "NASDAQ NMS - GLOBAL MARKET",
				FinnhubIndustry:      "Technology",
				Ipo:                  "1980-12-12",
				Logo:                 "https://finnhub.io/api/logo?symbol=AAPL",
				MarketCapitalization: 2634047,
				Name:                 "Apple Inc",
				Phone:                "14089961010.0",
				ShareOutstanding:     16426.79,
				Ticker:               "AAPL",
				Weburl:               "https://www.apple.com/",
			}
			break
		case "AMT":
			bodyResp = finnhub.CompanyProfile2{
				Country:              "US",
				Currency:             "USD",
				Exchange:             "NEW YORK STOCK EXCHANGE, INC.",
				FinnhubIndustry:      "Real Estate",
				Ipo:                  "1998-06-05",
				Logo:                 "https://finnhub.io/api/logo?symbol=AMT",
				MarketCapitalization: 118853.9,
				Name:                 "American Tower Corp",
				Phone:                "16173757500.0",
				ShareOutstanding:     444.33,
				Ticker:               "AMT",
				Weburl:               "http://www.americantower.com/",
			}
			break
		default:
			bodyResp = finnhub.CompanyProfile2{}
		}

		bodyByte, _ := json.Marshal(bodyResp)

		respHeader := http.Header{
			"Content-Type": {"application/json"},
		}
		return &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     respHeader,
			Body:       ioutil.NopCloser(bytes.NewReader(bodyByte)),
			Request:    req,
		}, nil
	} else if reqType == "quote" {
		bodyResp := finnhub.SymbolPriceFinnhub{}

		// Treat body from the request to get the symbol value from the URL query
		urlQuery := strings.Split(req.URL.RawQuery, "&")
		for _, query := range urlQuery {
			queryParams := strings.Split(string(query), "=")

			if queryParams[0] == "symbol" {
				symbol = queryParams[1]
			}
		}

		// If the symbol is invalid then return error, else returns the information
		// with the asset information based on the Alpha Vantage JSON template
		switch symbol {
		case "AAPL":
			bodyResp = finnhub.SymbolPriceFinnhub{
				C:  163.76,
				O:  158.735,
				H:  164.2,
				PC: 164.77,
				L:  157.8,
				T:  1638478802,
			}
			break
		case "VTI":
			bodyResp = finnhub.SymbolPriceFinnhub{
				C:  233.9,
				H:  234.68,
				L:  230.1,
				O:  230.21,
				PC: 230.16,
				T:  1638478800,
			}
			break
		case "AMT":
			bodyResp = finnhub.SymbolPriceFinnhub{
				C:  266.24,
				H:  267.74,
				L:  259.27,
				O:  259.73,
				PC: 258.95,
				T:  1638478002,
			}
			break
		default:
			bodyResp = finnhub.SymbolPriceFinnhub{}
		}

		bodyByte, _ := json.Marshal(bodyResp)

		respHeader := http.Header{
			"Content-Type": {"application/json"},
		}
		return &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     respHeader,
			Body:       ioutil.NopCloser(bytes.NewReader(bodyByte)),
			Request:    req,
		}, nil
	}

	return &http.Response{}, nil
}
