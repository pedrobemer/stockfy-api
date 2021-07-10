package main

import (
	"context"
	"fmt"
)

func (r *Resolver) SymbolsPriceAlpha(ctx context.Context, args SymbolList) ([]SymbolPrice, error) {
	var symbolsPrice []SymbolPrice
	var symbolPrice SymbolPrice

	for _, s := range args.Symbols {
		symbolPrice = getPriceAlphaVantage(s)

		symbolsPrice = append(symbolsPrice, symbolPrice)
	}

	return symbolsPrice, nil
}

func (r *Resolver) SymbolLookupAlpha(ctx context.Context, args FinhubArgs) (SymbolLookupInfo, error) {

	var symbolLookupUnique SymbolLookupInfo
	// var symbolTypes = map[string]string{
	// 	"Equity": "STOCK",
	// 	"ETF":    "ETF",
	// 	"REIT":   "REIT",
	// 	"FII":    "FII",
	// }

	var symbolLookup = verifySymbolAlpha(args.Symbol)
	fmt.Println(symbolLookup)
	// for _, s := range symbolLookup.Result {
	// 	if s.Symbol == args.Symbol {
	// 		symbolLookupUnique = s
	// 		symbolLookupUnique.Type = symbolTypes[symbolLookupUnique.Type]
	// 	}
	// }

	return symbolLookupUnique, nil
}

func getPriceAlphaVantage(symbol string) SymbolPrice {
	url := "https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=" + symbol + "&apikey=KIUG1ZKFZ13BI08F"

	var symbolPriceNotFormatted map[string]interface{}
	symbolPrice := SymbolPrice{}

	requestAndAssignToBody(url, &symbolPriceNotFormatted)

	formatAlphaVantageSymbolPrice(symbolPriceNotFormatted, &symbolPrice, symbol)

	return symbolPrice
}

func formatAlphaVantageSymbolPrice(unformatted map[string]interface{}, formatted *SymbolPrice, symbol string) {

	for k, v := range unformatted["Global Quote"].(map[string]interface{}) {
		switch k {
		case "01. symbol":
			formatted.Symbol = interfaceToString(v)
		case "02. open":
			formatted.OpenPrice = interfaceToFloat64(v)
		case "03. high":
			formatted.HighPrice = interfaceToFloat64(v)
		case "04. low":
			formatted.LowPrice = interfaceToFloat64(v)
		case "05. price":
			formatted.CurrentPrice = interfaceToFloat64(v)
		case "08. previous close":
			formatted.PrevClosePrice = interfaceToFloat64(v)
		default:
		}
	}
}

func verifySymbolAlpha(symbol string) SymbolLookupAlpha {
	url := "https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=" + symbol + "&apikey=KIUG1ZKFZ13BI08F"

	var symbolLookup SymbolLookupAlpha
	var test interface{}

	requestAndAssignToBody(url, &test)
	fmt.Println(test)

	return symbolLookup
}
