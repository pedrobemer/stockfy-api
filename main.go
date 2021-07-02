package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

type Resolver struct{}

// FinhubArgs is the args to finhub
type FinhubArgs struct {
	Symbol string
}

type SymbolList struct {
	Symbols []string
}

func formatSymbolPrice(unformatted SymbolPriceNotFormatted, formatted *SymbolPrice) {
	formatted.CurrentPrice = unformatted.C
	formatted.HighPrice = unformatted.H
	formatted.LowPrice = unformatted.L
	formatted.PrevClosePrice = unformatted.PC
	formatted.OpenPrice = unformatted.O
	formatted.MarketCap = unformatted.T
}

func (r *Resolver) SymbolPriceUS(ctx context.Context, args FinhubArgs) (SymbolPrice, error) {
	Symbol := args.Symbol

	stock := getPrice(Symbol)

	return stock, nil
}

func (r *Resolver) SymbolsPriceUS(ctx context.Context, args SymbolList) ([]SymbolPrice, error) {
	var symbolsPrice []SymbolPrice
	var symbolPrice SymbolPrice

	for _, s := range args.Symbols {
		symbolPrice = getPrice(s)

		symbolsPrice = append(symbolsPrice, symbolPrice)
	}

	return symbolsPrice, nil
}

func getPrice(symbol string) SymbolPrice {
	url := "https://finnhub.io/api/v1/quote?symbol=" + symbol + "&token=c2o3062ad3ie71thpra0"

	symbolPriceNotFormatted := SymbolPriceNotFormatted{}
	symbolPrice := SymbolPrice{}

	requestAndAssignToBody(url, &symbolPriceNotFormatted)

	formatSymbolPrice(symbolPriceNotFormatted, &symbolPrice)

	return symbolPrice
}

func main() {
	s, err := getSchema("./schema.graphql")
	if err != nil {
		panic(err)
	}

	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers()}

	schema := graphql.MustParseSchema(s, &Resolver{}, opts...)

	http.Handle("/", &relay.Handler{Schema: schema})
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func requestAndAssignToBody(url string, anyThing interface{}) {
	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		log.Fatal(err)
	}

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	jsonErr := json.Unmarshal(body, &anyThing)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}

func getSchema(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
