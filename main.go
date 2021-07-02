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

func (r *Resolver) Finhub(ctx context.Context, args FinhubArgs) (Stock, error) {
	fmt.Println("symbol")
	fmt.Println(args.Symbol)

	Symbol := args.Symbol
	url := "https://finnhub.io/api/v1/quote?symbol=" + Symbol + "&token=c2o3062ad3ie71thpra0"

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

	apple := Stock{}
	jsonErr := json.Unmarshal(body, &apple)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return apple, nil
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

func getSchema(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
