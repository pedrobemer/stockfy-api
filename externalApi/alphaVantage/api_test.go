package alphaVantage

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"stockfyApi/api/handlers/fiberHandlers"
	"stockfyApi/entity"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifySymbol2(t *testing.T) {
	MockDoFunc := func(req *http.Request) (*http.Response, error) {

		var symbol string
		bodyResp := SymbolLookupAlpha{}

		// Treat body from the request to get the symbol value from the URL query
		urlQuery := strings.Split(req.URL.RawQuery, "&")
		for _, query := range urlQuery {
			queryParams := strings.Split(string(query), "=")

			if queryParams[0] == "keywords" {
				symbol = queryParams[1]
			}
		}

		// If the symbol is invalid then return error, else returns the information
		// with the asset information based on the Alpha Vantage JSON template
		switch symbol {
		case "ITUB4.SA":
			bodyResp = SymbolLookupAlpha{
				BestMatches: []SymbolLookupInfo{
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
			break
		case "KNRI11.SA":
			bodyResp = SymbolLookupAlpha{
				BestMatches: []SymbolLookupInfo{
					{
						Symbol:      symbol + "O",
						Name:        "Kinea Renda Imobiliária Fundo de Investimento Imobiliário",
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
			bodyResp = SymbolLookupAlpha{
				BestMatches: []SymbolLookupInfo{
					{
						Symbol:      symbol + "O",
						Name:        "iShares S&P 500 Fundo de Investimento - Investimento No Exterior",
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
		case "AAPL":
			bodyResp = SymbolLookupAlpha{
				BestMatches: []SymbolLookupInfo{
					{
						Symbol:      symbol,
						Name:        "Apple Inc",
						Type:        "Equity",
						Region:      "United States",
						MarketOpen:  "09:30",
						MarketClose: "16:00",
						Timezone:    "UTC-04",
						Currency:    "USD",
						MatchScore:  "1.0000",
					},
					{
						Symbol:      symbol + "34.SAO",
						Name:        "Apple Inc",
						Type:        "Equity",
						Region:      "Brazil/Sao Paolo",
						MarketOpen:  "10:00",
						MarketClose: "17:30",
						Timezone:    "UTC-03",
						Currency:    "BRL",
						MatchScore:  "0.7500",
					},
				},
			}
			break
		case "VTI":
			bodyResp = SymbolLookupAlpha{
				BestMatches: []SymbolLookupInfo{
					{
						Symbol:      symbol,
						Name:        "Vanguard Total Stock Market ETF",
						Type:        "ETF",
						Region:      "United States",
						MarketOpen:  "09:30",
						MarketClose: "16:00",
						Timezone:    "UTC-04",
						Currency:    "USD",
						MatchScore:  "1.0000",
					},
					{
						Symbol:      symbol + "AX",
						Name:        "VANGUARD TOTAL INTERNATIONAL STOCK INDEX FUND ADMIRAL SHARES",
						Type:        "Mutual Fund",
						Region:      "United States",
						MarketOpen:  "09:30",
						MarketClose: "16:00",
						Timezone:    "UTC-04",
						Currency:    "USD",
						MatchScore:  "0.7500",
					},
				},
			}
			break
		case "AMT":
			bodyResp = SymbolLookupAlpha{
				BestMatches: []SymbolLookupInfo{
					{
						Symbol:      symbol,
						Name:        "American Tower Corp",
						Type:        "Equity",
						Region:      "United States",
						MarketOpen:  "09:30",
						MarketClose: "16:00",
						Timezone:    "UTC-04",
						Currency:    "USD",
						MatchScore:  "1.0000",
					},
					{
						Symbol:      symbol + "A",
						Name:        "Amistar Corp",
						Type:        "Equity",
						Region:      "United States",
						MarketOpen:  "09:30",
						MarketClose: "16:00",
						Timezone:    "UTC-04",
						Currency:    "USD",
						MatchScore:  "0.8571",
					},
				},
			}
			break
		default:
			bodyResp = SymbolLookupAlpha{}
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

	type test struct {
		symbol               string
		expectedSymbolLookup entity.SymbolLookup
	}

	tests := []test{
		{
			symbol: "ITUB4.SA",
			expectedSymbolLookup: entity.SymbolLookup{
				Symbol:   "ITUB4",
				Fullname: "Itaú Unibanco Holding S.A",
				Type:     "Equity",
			},
		},
		{
			symbol: "KNRI11.SA",
			expectedSymbolLookup: entity.SymbolLookup{
				Symbol:   "KNRI11",
				Fullname: "Kinea Renda Imobiliária Fundo De Investimento Imobiliário",
				Type:     "ETF",
			},
		},
		{
			symbol: "AAPL",
			expectedSymbolLookup: entity.SymbolLookup{
				Symbol:   "AAPL",
				Fullname: "Apple Inc",
				Type:     "Equity",
			},
		},
		{
			symbol: "AMT",
			expectedSymbolLookup: entity.SymbolLookup{
				Symbol:   "AMT",
				Fullname: "American Tower Corp",
				Type:     "Equity",
			},
		},
		{
			symbol: "VTI",
			expectedSymbolLookup: entity.SymbolLookup{
				Symbol:   "VTI",
				Fullname: "Vanguard Total Stock Market ETF",
				Type:     "ETF",
			},
		},
	}

	mockAlphaClient := MockClient{
		Client: fiberHandlers.MockClient{
			DoFunc: MockDoFunc,
		},
	}

	alpha := AlphaApi{
		Token:              "Test",
		HttpOutsideRequest: mockAlphaClient.HttpOutsideClientRequest,
	}

	for _, testCase := range tests {
		symbolLookup := alpha.VerifySymbol2(testCase.symbol)

		assert.Equal(t, testCase.expectedSymbolLookup, symbolLookup)
	}

}

func TestGetPrice(t *testing.T) {
	MockDoFunc := func(req *http.Request) (*http.Response, error) {

		var symbol string
		bodyResp := SymbolPriceAlpha{}

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
		case "ITUB4.SA":
			bodyResp = SymbolPriceAlpha{
				GlobalQuote: SymbolPriceInfo{
					Symbol:        symbol,
					Open:          "22.5000",
					High:          "22.4400",
					Low:           "21.9100",
					Price:         "22.4400",
					Volume:        "95434007",
					LatestDay:     "2021-11-23",
					PrevClose:     "22.0700",
					Change:        "0.3700",
					ChangePercent: "1.6765%",
				},
			}
			break
		case "KNRI11.SA":
			bodyResp = SymbolPriceAlpha{
				GlobalQuote: SymbolPriceInfo{
					Symbol:        symbol,
					Open:          "133.6000",
					High:          "134.0000",
					Low:           "132.7001",
					Price:         "132.9500",
					Volume:        "23140",
					LatestDay:     "2021-11-23",
					PrevClose:     "134.0000",
					Change:        "-1.0500",
					ChangePercent: "-0.7836%",
				},
			}
			break
		case "IVVB11.SA":
			bodyResp = SymbolPriceAlpha{
				GlobalQuote: SymbolPriceInfo{
					Symbol:        symbol,
					Open:          "285.8200",
					High:          "288.4500",
					Low:           "284.5700",
					Price:         "285.4500",
					Volume:        "344217",
					LatestDay:     "2021-11-23",
					PrevClose:     "285.6900",
					Change:        "-0.2400",
					ChangePercent: "-0.0840%",
				},
			}
			break
		case "AAPL":
			bodyResp = SymbolPriceAlpha{
				GlobalQuote: SymbolPriceInfo{
					Symbol:        symbol,
					Open:          "161.2000",
					High:          "161.8000",
					Low:           "159.0601",
					Price:         "161.4100",
					Volume:        "95434007",
					LatestDay:     "2021-11-23",
					PrevClose:     "161.2000",
					Change:        "0.3900",
					ChangePercent: "0.2422%",
				},
			}
			break
		case "VTI":
			bodyResp = SymbolPriceAlpha{
				GlobalQuote: SymbolPriceInfo{
					Symbol:        symbol,
					Open:          "240.1000",
					High:          "241.1000",
					Low:           "238.2364",
					Price:         "240.4000",
					Volume:        "3779100",
					LatestDay:     "2021-11-23",
					PrevClose:     "240.3200",
					Change:        "0.0800",
					ChangePercent: "0.0333%",
				},
			}
			break
		case "AMT":
			bodyResp = SymbolPriceAlpha{
				GlobalQuote: SymbolPriceInfo{
					Symbol:        symbol,
					Open:          "258.2600",
					High:          "262.2100",
					Low:           "256.7500",
					Price:         "262.0000",
					Volume:        "3779100",
					LatestDay:     "2021-11-23",
					PrevClose:     "257.6300",
					Change:        "4.3700",
					ChangePercent: "1.6962%",
				},
			}
			break
		default:
			bodyResp = SymbolPriceAlpha{}
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

	type test struct {
		symbol               string
		expectedSymbolLookup entity.SymbolPrice
	}

	tests := []test{
		{
			symbol: "ITUB4.SA",
			expectedSymbolLookup: entity.SymbolPrice{
				Symbol:         "ITUB4",
				OpenPrice:      22.5000,
				HighPrice:      22.4400,
				LowPrice:       21.9100,
				CurrentPrice:   22.4400,
				PrevClosePrice: 22.0700,
			},
		},
		{
			symbol: "KNRI11.SA",
			expectedSymbolLookup: entity.SymbolPrice{
				Symbol:         "KNRI11",
				OpenPrice:      133.6000,
				HighPrice:      134.0000,
				LowPrice:       132.7001,
				CurrentPrice:   132.9500,
				PrevClosePrice: 134.0000,
			},
		},
		{
			symbol: "AAPL",
			expectedSymbolLookup: entity.SymbolPrice{
				Symbol:         "AAPL",
				OpenPrice:      161.2000,
				HighPrice:      161.8000,
				LowPrice:       159.0601,
				CurrentPrice:   161.4100,
				PrevClosePrice: 161.2000,
			},
		},
		{
			symbol: "AMT",
			expectedSymbolLookup: entity.SymbolPrice{
				Symbol:         "AMT",
				OpenPrice:      258.2600,
				HighPrice:      262.2100,
				LowPrice:       256.7500,
				CurrentPrice:   262.0000,
				PrevClosePrice: 257.6300,
			},
		},
		{
			symbol: "VTI",
			expectedSymbolLookup: entity.SymbolPrice{
				Symbol:         "VTI",
				OpenPrice:      240.100,
				HighPrice:      241.1000,
				LowPrice:       238.2364,
				CurrentPrice:   240.4000,
				PrevClosePrice: 240.3200,
			},
		},
	}

	mockAlphaClient := MockClient{
		Client: fiberHandlers.MockClient{
			DoFunc: MockDoFunc,
		},
	}

	alpha := AlphaApi{
		Token:              "Test",
		HttpOutsideRequest: mockAlphaClient.HttpOutsideClientRequest,
	}

	for _, testCase := range tests {
		symbolLookup := alpha.GetPrice(testCase.symbol)

		assert.Equal(t, testCase.expectedSymbolLookup, symbolLookup)
	}

}
