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
	app.Get("/asset/:symbol/orders=:orders?", func(c *fiber.Ctx) error {
		var symbolQuery []tables.AssetQueryReturn

		if c.Params("orders") != "ALL" && c.Params("orders") != "ONLYINFO" &&
			c.Params("orders") != "" {
			message := "Wrong REST API request. Please see our README.md in our Git repository to understand the possible values for orders value."
			return c.SendString(message)
		}

		symbolQuery = tables.SearchAsset(*dbpool, c.Params("symbol"),
			c.Params("orders"))

		jsonQuery, err := json.Marshal(symbolQuery)
		if err != nil {
			log.Fatal(err)
		}

		return c.SendString(string(jsonQuery))

	})

	// REST API to fetch all or some asset type.
	app.Get("/assettypes/:type/:country?", func(c *fiber.Ctx) error {
		var assetTypeQuery []tables.AssetTypeApiReturn

		if c.Params("type") != "ALL" && c.Params("country") == "" {
			return c.SendString("Wrong REST API path. You need to specify the country code when you search a specific asset type. For example: /assettypes/ETF/US")
		}

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
				" asset type from " + c.Params("country"))
		}

		jsonQuery, err := json.Marshal(assetTypeQuery)
		if err != nil {
			log.Fatal(err)
		}

		return c.SendString(string(jsonQuery))

	})

	// REST API to fetch all or some asset type.
	app.Get("/assetPerAssetTypes/:type-:country/withOrders=:withOrders", func(c *fiber.Ctx) error {
		var assetTypeQuery []tables.AssetTypeApiReturn

		var withOrdersInfo bool
		if c.Params("withOrders") == "true" {
			withOrdersInfo = true
		} else if c.Params("withOrders") == "false" {
			withOrdersInfo = false
		} else {
			return c.SendString("Wrong REST API path. The withOrdes accepts only boolean variables true or false.")
		}

		assetTypeQuery = tables.SearchAssetsPerAssetType(*dbpool, c.Params("type"), c.Params("country"), withOrdersInfo)

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

			assetTypeQuery, err := tables.FetchAssetType(*dbpool, specificFetch,
				orderInsert.AssetType, orderInsert.Country)
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
			symbolInserted := tables.CreateAsset(*dbpool, assetInsert,
				assetTypeId, sectorId)
			assetId = symbolInserted.Id
			fmt.Println(symbolInserted)
		} else {
			symbolQuery := tables.SearchAsset(*dbpool, orderInsert.Symbol, "")
			assetId = symbolQuery[0].Id
		}

		brokerageReturn, _ := tables.FetchBrokerage(*dbpool, "SINGLE",
			orderInsert.Brokerage)
		brokerageId = brokerageReturn[0].Id

		fmt.Println(assetId, brokerageId)
		var orderReturn tables.OrderApiReturn

		orderReturn = tables.CreateOrder(*dbpool, orderInsert, assetId,
			brokerageId)
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
