package finnhub

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
		bodyResp := SymbolLookupFinnhub{}

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
		case "ITUB3.SA":
			bodyResp = SymbolLookupFinnhub{
				Count: 1,
				Result: []SymbolLookupInfo{
					{
						Symbol:        symbol,
						DisplaySymbol: symbol,
						Type:          "Common Stock",
						Description:   "ITAU UNIBANCO HOLDING SA",
					},
				},
			}
			break
		case "IVVB11.SA":
			bodyResp = SymbolLookupFinnhub{
				Count: 1,
				Result: []SymbolLookupInfo{
					{
						Symbol:        symbol,
						DisplaySymbol: symbol,
						Type:          "ETP",
						Description:   "ISHARES S&P 500 FIC FI IE",
					},
				},
			}
			break
		case "AAPL":
			bodyResp = SymbolLookupFinnhub{
				Count: 2,
				Result: []SymbolLookupInfo{
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
			bodyResp = SymbolLookupFinnhub{
				Count: 2,
				Result: []SymbolLookupInfo{
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
			bodyResp = SymbolLookupFinnhub{
				Count: 2,
				Result: []SymbolLookupInfo{
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
			bodyResp = SymbolLookupFinnhub{}
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
			symbol: "ITUB3.SA",
			expectedSymbolLookup: entity.SymbolLookup{
				Symbol:   "ITUB3",
				Fullname: "Itau Unibanco Holding SA",
				Type:     "Common Stock",
			},
		},
		{
			symbol: "AAPL",
			expectedSymbolLookup: entity.SymbolLookup{
				Symbol:   "AAPL",
				Fullname: "Apple Inc",
				Type:     "Common Stock",
			},
		},
		{
			symbol: "AMT",
			expectedSymbolLookup: entity.SymbolLookup{
				Symbol:   "AMT",
				Fullname: "American Tower Corp",
				Type:     "REIT",
			},
		},
		{
			symbol: "VTI",
			expectedSymbolLookup: entity.SymbolLookup{
				Symbol:   "VTI",
				Fullname: "Vanguard Total Stock Mkt ETF",
				Type:     "ETP",
			},
		},
	}

	mockFinnhubClient := MockClient{
		Client: fiberHandlers.MockClient{
			DoFunc: MockDoFunc,
		},
	}

	finnhubApi := FinnhubApi{
		Token:              "Test",
		HttpOutsideRequest: mockFinnhubClient.HttpOutsideClientRequest,
	}

	for _, testCase := range tests {
		symbolLookup := finnhubApi.VerifySymbol2(testCase.symbol)

		assert.Equal(t, testCase.expectedSymbolLookup, symbolLookup)
	}

}

func TestGetPrice(t *testing.T) {
	MockDoFunc := func(req *http.Request) (*http.Response, error) {

		var symbol string
		bodyResp := SymbolPriceFinnhub{}

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
			bodyResp = SymbolPriceFinnhub{
				O:  161.2000,
				H:  161.8000,
				L:  159.0601,
				C:  161.4100,
				PC: 161.2000,
			}
			break
		case "VTI":
			bodyResp = SymbolPriceFinnhub{
				O:  240.100,
				H:  241.1000,
				L:  238.2364,
				C:  240.4000,
				PC: 240.3200,
			}
			break
		case "AMT":
			bodyResp = SymbolPriceFinnhub{
				O:  258.2600,
				H:  262.2100,
				L:  256.7500,
				C:  262.0000,
				PC: 257.6300,
			}
			break
		default:
			bodyResp = SymbolPriceFinnhub{}
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

	mockFinnhubClient := MockClient{
		Client: fiberHandlers.MockClient{
			DoFunc: MockDoFunc,
		},
	}

	finnhubAPi := FinnhubApi{
		Token:              "Test",
		HttpOutsideRequest: mockFinnhubClient.HttpOutsideClientRequest,
	}

	for _, testCase := range tests {
		symbolLookup := finnhubAPi.GetPrice(testCase.symbol)

		assert.Equal(t, testCase.expectedSymbolLookup, symbolLookup)
	}

}
