package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "pedrobemer"
	DB_PASSWORD = "pirulito"
	DB_NAME     = "stockfy"
)

type AssetType struct {
	Id      string `db:"id"`
	Type    string `db:"type"`
	Name    string `db:"name"`
	Country string `db:"country"`
}

type SymbolQuery struct {
	Id               string `db:"id"`
	Preference       *string
	Fullname         string `db:"fullname"`
	Symbol           string `db:"symbol"`
	AssetTypeId      string `db:"assettype_id"`
	AssetTypeType    string `db:"assettype_type"`
	AssetTypeName    string `db:"assettype_name"`
	AssetTypeCountry string `db:"assettype_country"`
	// AssetType  AssetTypeStr `db:"asset_type"`
}

func main() {

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)

	dbpool, err := pgxpool.Connect(context.Background(), dbinfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	app := fiber.New()

	// REST API to fetch some asset symbol.
	app.Get("/query/asset/Symbol=:symbol", func(c *fiber.Ctx) error {
		var symbolQuery []*SymbolQuery
		columns := " asset.id, fullname, symbol, "
		fk_columns := "assettype.id as assettype_id" +
			", assettype.type as assettype_type" +
			", assettype.name as assettype_name" +
			", assettype.country as assettype_country "
		condition := " WHERE symbol = $1"
		query := "SELECT" + columns + fk_columns +
			"FROM asset JOIN assettype ON asset.asset_type_id = assettype.id" +
			condition

		err := pgxscan.Select(context.Background(), dbpool, &symbolQuery, query,
			c.Params("symbol"))
		if err != nil {
			fmt.Println("ERRROU")
		}

		jsonQuery, err := json.Marshal(symbolQuery)
		if err != nil {
			log.Fatal(err)
		}

		return c.SendString(string(jsonQuery))

	})

	// REST API to fetch all or some asset type.
	app.Get("/query/assettypes/type=:type", func(c *fiber.Ctx) error {
		var assetTypeQuery []*AssetType

		queryDefault := "SELECT id, type, name, country FROM assettype "
		if c.Params("type") == "ALL" {
			err := pgxscan.Select(context.Background(), dbpool, &assetTypeQuery,
				queryDefault)
			if err != nil {
				fmt.Println("ERRROU")
			}
		} else {
			query := queryDefault + "where type=$1"
			err := pgxscan.Select(context.Background(), dbpool, &assetTypeQuery,
				query, c.Params("type"))
			if err != nil {
				fmt.Println("ERRROU")
			}
		}

		jsonQuery, err := json.Marshal(assetTypeQuery)
		if err != nil {
			log.Fatal(err)
		}

		return c.SendString(string(jsonQuery))

	})

	app.Listen(":3000")

	// var fullname
	// rows, err := dbpool.Query(context.Background(), "SELECT asset.id, fullname, symbol FROM asset JOIN assettype ON asset.asset_type_id = assettype.id")
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	// 	os.Exit(1)
	// }

	// s, err := getSchema("./schema.graphql")
	// if err != nil {
	// 	panic(err)
	// }

	// opts := []graphql.SchemaOpt{graphql.UseFieldResolvers()}

	// schema := graphql.MustParseSchema(s, &Resolver{}, opts...)

	// http.Handle("/", &relay.Handler{Schema: schema})
	// log.Fatal(http.ListenAndServe(":3000", nil))
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
