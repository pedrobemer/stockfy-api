package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"stockfyApi/tables"
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
	app.Get("/asset/:symbol-orders=:orders?", func(c *fiber.Ctx) error {
		var symbolQuery []*tables.AssetQueryReturn
		var query string
		if c.Params("orders") == "" {
			query = "SELECT a.id, symbol, preference, fullname, json_build_object('id', at.id, 'type', at.type, 'name', at.name, 'country', at.country) as asset_type FROM asset as a INNER JOIN assettype as at ON a.asset_type_id = at.id INNER JOIN orders as o ON a.id = o.asset_id WHERE a.symbol=$1 GROUP BY a.symbol, a.id, preference, fullname, at.type, at.id, at.name, at.country;"
		} else if c.Params("orders") == "ALL" {
			query = "SELECT a.id, symbol, preference, a.fullname, json_build_object('id', at.id, 'type', at.type, 'name', at.name, 'country', at.country) as asset_type, json_build_object('totalQuantity', sum(o.quantity), 'weightedAdjPrice', SUM(o.quantity * price)/SUM(o.quantity), 'weightedAveragePrice', (SUM(o.quantity*o.price) FILTER(WHERE o.order_type = 'buy'))/(SUM(o.quantity) FILTER(WHERE o.order_type = 'buy'))) as orders_info ,json_agg(json_build_object('id', o.id, 'quantity', o.quantity, 'price', o.price, 'currency', o.currency, 'ordertype', o.order_type, 'date', date, 'brokerage', json_build_object('id', b.id, 'name', b.name, 'country', b.country))) as orders_list FROM asset as a INNER JOIN assettype as at ON a.asset_type_id = at.id INNER JOIN orders as o ON a.id = o.asset_id INNER JOIN brokerage as b ON o.brokerage_id = b.id WHERE a.symbol=$1 GROUP BY a.symbol, a.id, preference, a.fullname, at.type, at.id, at.name, at.country;"
		} else if c.Params("orders") == "ONLYINFO" {
			query = "SELECT a.id, symbol, preference, a.fullname, json_build_object('id', at.id, 'type', at.type, 'name', at.name, 'country', at.country) as asset_type,  json_build_object('totalQuantity', sum(o.quantity), 'weightedAdjPrice', SUM(o.quantity * price)/SUM(o.quantity), 'weightedAveragePrice', (SUM(o.quantity*o.price) FILTER(WHERE o.order_type = 'buy'))/(SUM(o.quantity) FILTER(WHERE o.order_type = 'buy'))) as orders_info FROM asset as a INNER JOIN assettype as at ON a.asset_type_id = at.id INNER JOIN orders as o ON a.id = o.asset_id INNER JOIN brokerage as b ON o.brokerage_id = b.id WHERE a.symbol=$1 GROUP BY a.symbol, a.id, preference, a.fullname, at.type, at.id, at.name, at.country;"
		} else {
			fmt.Println("Wrong API Rest")
			message := "Wrong REST API request. Please see our README.md in our Git repository to understand how to do this request."
			return c.SendString(message)
		}

		err := pgxscan.Select(context.Background(), dbpool, &symbolQuery, query,
			c.Params("symbol"))
		if err != nil {
			fmt.Println(err)
		}

		jsonQuery, err := json.Marshal(symbolQuery)
		if err != nil {
			log.Fatal(err)
		}

		return c.SendString(string(jsonQuery))

	})

	// REST API to fetch all or some asset type.
	app.Get("/assettypes/:type-:country", func(c *fiber.Ctx) error {
		var assetTypeQuery []tables.AssetTypeApiReturn

		var specificFetch bool
		if c.Params("type") == "ALL" {
			specificFetch = false
		} else {
			specificFetch = true
		}

		assetTypeQuery, err = tables.FetchAssetType(*dbpool, specificFetch,
			c.Params("type"), c.Params("country"))
		if err != nil {
			panic(err)
		}

		if assetTypeQuery == nil {
			return c.SendString("There is not any " + c.Params("type") +
				" Asset type from " + c.Params("country"))
		}

		jsonQuery, err := json.Marshal(assetTypeQuery)
		if err != nil {
			log.Fatal(err)
		}

		return c.SendString(string(jsonQuery))

	})

	// REST API to fetch a specific sector
	app.Get("/sector/:sector", func(c *fiber.Ctx) error {
		var sectorQuery []tables.SectorApiReturn

		sectorQuery = tables.FetchSector(*dbpool, c.Params("sector"))

		jsonQuery, err := json.Marshal(sectorQuery)
		if err != nil {
			log.Fatal(err)
		}

		return c.SendString(string(jsonQuery))

	})

	app.Get("/brokerage/:searchType/:name?", func(c *fiber.Ctx) error {
		var brokerageQuery []tables.BrokerageApiReturn

		if c.Params("searchType") == "SINGLE" && c.Params("name") == "" {
			return c.SendString("The " + c.Params("searchType") + " needs the name field to be populated with a name of some brokerage firm. For example: /brokerage/SINGLE/<BROKERAGE FIRM NAME>.")
		} else if c.Params("searchType") == "COUNTRY" &&
			(c.Params("name") != "BR" && c.Params("name") != "US") {
			return c.SendString("The " + c.Params("searchType") + " needs the name field to be populated with the country code BR or US. For example: /brokerage/COUNTRY/BR.")
		} else if c.Params("searchType") != "SINGLE" &&
			c.Params("searchType") != "ALL" &&
			c.Params("searchType") != "COUNTRY" {
			return c.SendString("Wrong Path. Please see the documentation from our REST API.")
		}

		brokerageQuery, _ = tables.FetchBrokerage(*dbpool, c.Params("searchType"),
			c.Params("name"))

		jsonQuery, err := json.Marshal(brokerageQuery)
		if err != nil {
			log.Fatal(err)
		}

		return c.SendString(string(jsonQuery))
	})
	// REST API to insert in the sector table a new registered sector
	app.Post("/sector", func(c *fiber.Ctx) error {

		var sectorBodyPost tables.SectorBodyPost
		if err := c.BodyParser(&sectorBodyPost); err != nil {
			fmt.Println(err)
		}
		fmt.Println(sectorBodyPost)

		var sectorInsert []tables.SectorApiReturn
		sectorInsert, _ = tables.CreateSector(*dbpool, sectorBodyPost.Sector)

		jsonQuery, err := json.Marshal(sectorInsert[0])
		if err != nil {
			log.Fatal(err)
		}

		return c.SendString(string(jsonQuery))

	})

	// REST API to insert a new asset in the asset table
	app.Post("/asset", func(c *fiber.Ctx) error {

		var assetInsert tables.AssetInsert
		if err := c.BodyParser(&assetInsert); err != nil {
			fmt.Println(err)
		}
		fmt.Println(assetInsert)

		var condAssetExist = "symbol='" + assetInsert.Symbol + "'"
		assetExist := verifyRowExistence(*dbpool, "asset", condAssetExist)
		fmt.Println(assetExist)
		if assetExist {
			return c.SendString(assetInsert.Symbol + " already exist in your database")
		}

		var assetTypeId string
		var sectorId string

		assetTypeQuery, _ := tables.FetchAssetType(*dbpool, true,
			assetInsert.AssetType, assetInsert.Country)
		assetTypeId = assetTypeQuery[0].Id

		sectorQuery, _ := tables.CreateSector(*dbpool, assetInsert.Sector)
		if sectorQuery != nil {
			sectorId = sectorQuery[0].Id
		}

		fmt.Println(assetTypeId, sectorId)

		symbolInsert := tables.CreateAsset(*dbpool, assetInsert, assetTypeId, sectorId)

		jsonQuery, err := json.Marshal(symbolInsert)
		if err != nil {
			log.Fatal(err)
		}

		return c.SendString(string(jsonQuery))

	})

	// REST API to register an order for a given asset.
	app.Post("/orders", func(c *fiber.Ctx) error {

		var orderInsert tables.OrderBodyPost
		if err := c.BodyParser(&orderInsert); err != nil {
			fmt.Println(err)
		}
		fmt.Println(orderInsert)

		var assetId string
		var assetTypeId string
		var sectorId string
		var brokerageId string

		var assetExist bool
		var condAssetExist = "symbol='" + orderInsert.Symbol + "'"
		assetExist = verifyRowExistence(*dbpool, "asset", condAssetExist)

		if !assetExist {
			var specificFetch = true
			var assetInsert tables.AssetInsert

			assetTypeQuery, err := tables.FetchAssetType(*dbpool, specificFetch, orderInsert.AssetType,
				orderInsert.Country)
			if err != nil {
				panic(err)
			}
			assetTypeId = assetTypeQuery[0].Id

			if orderInsert.Sector != "" {
				sectorInfo, _ := tables.CreateSector(*dbpool, orderInsert.Sector)
				sectorId = sectorInfo[0].Id
			}

			assetInsert.Fullname = orderInsert.Fullname
			assetInsert.Symbol = orderInsert.Symbol
			assetInsert.Country = orderInsert.Country
			assetInsert.AssetType = assetTypeQuery[0].Type
			symbolInserted := tables.CreateAsset(*dbpool, assetInsert, assetTypeId, sectorId)
			fmt.Println(symbolInserted)
		} else {
			assetId = tables.SearchAsset(*dbpool, orderInsert.Symbol)
		}

		brokerageReturn, _ := tables.FetchBrokerage(*dbpool, "SINGLE", orderInsert.Brokerage)
		brokerageId = brokerageReturn[0].Id

		fmt.Println(assetId, brokerageId)
		var orderReturn tables.OrderApiReturn

		orderReturn = tables.CreateOrder(*dbpool, orderInsert, assetId, brokerageId)

		jsonQuery, err := json.Marshal(orderReturn)
		if err != nil {
			log.Fatal(err)
		}

		return c.SendString(string(jsonQuery))

	})

	app.Listen(":3000")

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
