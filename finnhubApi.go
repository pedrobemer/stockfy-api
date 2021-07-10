package main

import "context"

func (r *Resolver) SymbolPriceFinnhub(ctx context.Context, args FinhubArgs) (SymbolPrice, error) {
	Symbol := args.Symbol

	stock := getPriceFinnhub(Symbol)

	return stock, nil
}

func (r *Resolver) SymbolsPriceFinnhub(ctx context.Context, args SymbolList) ([]SymbolPrice, error) {
	var symbolsPrice []SymbolPrice
	var symbolPrice SymbolPrice

	for _, s := range args.Symbols {
		symbolPrice = getPriceFinnhub(s)

		symbolsPrice = append(symbolsPrice, symbolPrice)
	}

	return symbolsPrice, nil
}

func (r *Resolver) SymbolLookupFinnhub(ctx context.Context, args FinhubArgs) (SymbolLookupInfo, error) {

	var symbolLookupUnique SymbolLookupInfo
	var symbolTypes = map[string]string{
		"Common Stock": "STOCK",
		"ETP":          "ETF",
		"REIT":         "REIT",
	}

	var symbolLookup = verifySymbolFinnhub(args.Symbol)

	for _, s := range symbolLookup.Result {
		if s.Symbol == args.Symbol {
			symbolLookupUnique = s
			symbolLookupUnique.Type = symbolTypes[symbolLookupUnique.Type]
		}
	}

	return symbolLookupUnique, nil
}

func getPriceFinnhub(symbol string) SymbolPrice {
	url := "https://finnhub.io/api/v1/quote?symbol=" + symbol + "&token=c2o3062ad3ie71thpra0"

	symbolPriceNotFormatted := SymbolPriceNotFormatted{}
	symbolPrice := SymbolPrice{}

	requestAndAssignToBody(url, &symbolPriceNotFormatted)

	formatFinhubSymbolPrice(symbolPriceNotFormatted, &symbolPrice, symbol)

	return symbolPrice
}

func formatFinhubSymbolPrice(unformatted SymbolPriceNotFormatted, formatted *SymbolPrice, symbol string) {
	formatted.Symbol = symbol
	formatted.CurrentPrice = unformatted.C
	formatted.HighPrice = unformatted.H
	formatted.LowPrice = unformatted.L
	formatted.PrevClosePrice = unformatted.PC
	formatted.OpenPrice = unformatted.O
	formatted.MarketCap = unformatted.T
}

func verifySymbolFinnhub(symbol string) SymbolLookupFinnhub {
	url := "https://finnhub.io/api/v1/search?q=" + symbol + "&token=c2o3062ad3ie71thpra0"

	var symbolLookup SymbolLookupFinnhub

	requestAndAssignToBody(url, &symbolLookup)

	return symbolLookup
}
