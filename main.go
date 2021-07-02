package main

import (
	"context"
	"encoding/json"
	"fmt"
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

func (r *Resolver) SymbolPriceUS(ctx context.Context, args FinhubArgs) (Stock, error) {
	Symbol := args.Symbol

	stock := getPrice(Symbol)
	fmt.Println("stock")
	fmt.Println(stock.C)

	return stock, nil
}

func (r *Resolver) SymbolsPriceUS(ctx context.Context, args SymbolList) ([]Stock, error) {
	var symbolsPrice []Stock
	var symbolPrice Stock

	for i, s := range args.Symbols {
		fmt.Println(i, s)
		symbolPrice = getPrice(s)
		symbolsPrice = append(symbolsPrice, symbolPrice)
	}

	return symbolsPrice, nil
}

func getPrice(symbol string) Stock {
	url := "https://finnhub.io/api/v1/quote?symbol=" + symbol + "&token=c2o3062ad3ie71thpra0"

	stock := Stock{}

	requestAndAssignToBody(url, &stock)

	return stock
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
