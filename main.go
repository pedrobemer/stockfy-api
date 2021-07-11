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
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "pedrobemer"
	DB_PASSWORD = "pirulito"
	DB_NAME     = "stockfy"
)

type AssetQueryReturn struct {
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

type SectorBodyPost struct {
	Sector string `json:"sector"`
}

type OrderBodyPost struct {
	Symbol    string  `json:"symbol"`
	Brokerage string  `json:"brokerage"`
	Quantity  float64 `json:"quantity"`
	Price     float64 `json:"price"`
	Currency  string  `json:"currency"`
	OrderType string  `json:"orderType"`
	Date      string  `json:"date"`
}

type AssetBodyPost struct {
	AssetType  string `json:"assetType"`
	Sector     string `json:"sector"`
	Symbol     string `json:"symbol"`
	Fullname   string `json:"fullname"`
	Preference string `json:"preference"`
}

type AssetTypeApiReturn struct {
	Id      string `db:"id"`
	Type    string `db:"type"`
	Name    string `db:"name"`
	Country string `db:"country"`
}

type SectorApiReturn struct {
	Id   string `db:"id"`
	Name string `db:"name"`
}

type OrderApiReturn struct {
	Id        string  `db:"id"`
	Quantity  float64 `db:"quantity"`
	Price     float64 `db:"price"`
	Currency  string  `db:"currency"`
	OrderType string  `db:"order_type"`
	Date      string  `db:"date"`
}

type AssetApiReturn struct {
	Id         string `db:"id"`
	Preference string `db:"preference"`
	Fullname   string `db:"fullname"`
	Symbol     string `db:"symbol"`
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
	app.Get("/asset/:symbol", func(c *fiber.Ctx) error {
		var symbolQuery []*AssetQueryReturn
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
	app.Get("/assettypes/:type", func(c *fiber.Ctx) error {
		var assetTypeQuery []*AssetTypeApiReturn

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

	// REST API to fetch a specific sector
	app.Get("/sector/:sector", func(c *fiber.Ctx) error {
		var sectorQuery []*SectorApiReturn

		queryDefault := "SELECT id, name FROM sector "
		if c.Params("type") == "ALL" {
			err := pgxscan.Select(context.Background(), dbpool, &sectorQuery,
				queryDefault)
			if err != nil {
				fmt.Println("ERRROU")
			}
		} else {
			query := queryDefault + "where name=$1"
			err := pgxscan.Select(context.Background(), dbpool, &sectorQuery,
				query, c.Params("sector"))
			if err != nil {
				fmt.Println("ERRROU")
			}
		}

		jsonQuery, err := json.Marshal(sectorQuery)
		if err != nil {
			log.Fatal(err)
		}

		return c.SendString(string(jsonQuery))

	})

	// REST API to insert in the sector table a new registered sector
	app.Post("/sector", func(c *fiber.Ctx) error {

		var sectorBodyPost SectorBodyPost
		if err := c.BodyParser(&sectorBodyPost); err != nil {
			fmt.Println(err)
		}
		fmt.Println(sectorBodyPost)

		tx, err := dbpool.Begin(context.Background())
		if err != nil {
			log.Panic(err)
		}

		defer tx.Rollback(context.Background())

		var sectorInsert SectorApiReturn
		insertRow := "INSERT INTO sector(name) VALUES ($1) RETURNING id, name;"

		row := tx.QueryRow(context.Background(), insertRow,
			sectorBodyPost.Sector)
		err = row.Scan(&sectorInsert.Id, &sectorInsert.Name)
		if err != nil {
			log.Panic(err)
		}

		err = tx.Commit(context.Background())
		if err != nil {
			log.Panic(err)
		}

		jsonQuery, err := json.Marshal(sectorInsert)
		if err != nil {
			log.Fatal(err)
		}

		return c.SendString(string(jsonQuery))

	})

	// REST API to insert a new asset in the asset table
	app.Post("/asset", func(c *fiber.Ctx) error {

		var assetInsert AssetBodyPost
		if err := c.BodyParser(&assetInsert); err != nil {
			fmt.Println(err)
		}
		fmt.Println(assetInsert)

		tx, err := dbpool.Begin(context.Background())
		if err != nil {
			log.Panic(err)
		}

		defer tx.Rollback(context.Background())

		var assetTypeId string
		var sectorId string
		queryAssetTypeId := "SELECT id FROM assettype WHERE type=$1"
		querySectorId := "SELECT id FROM sector WHERE name=$1"
		tx.QueryRow(context.Background(), queryAssetTypeId,
			assetInsert.AssetType).Scan(&assetTypeId)
		tx.QueryRow(context.Background(), querySectorId,
			assetInsert.Sector).Scan(&sectorId)

		fmt.Println(assetTypeId, sectorId)
		var symbolInsert AssetApiReturn
		var insertRow string
		var row pgx.Row
		if sectorId != "" {
			insertRow = "INSERT INTO asset(preference, fullname, symbol, asset_type_id, sector_id) VALUES ($1, $2, $3, $4, $5) RETURNING id, preference, fullname, symbol;"
			row = tx.QueryRow(context.Background(), insertRow,
				assetInsert.Preference, assetInsert.Fullname, assetInsert.Symbol,
				assetTypeId, sectorId)
		} else {
			insertRow = "INSERT INTO asset(preference, fullname, symbol, asset_type_id) VALUES ($1, $2, $3, $4) RETURNING id, preference, fullname, symbol;"
			row = tx.QueryRow(context.Background(), insertRow,
				assetInsert.Preference, assetInsert.Fullname, assetInsert.Symbol,
				assetTypeId)
		}

		fmt.Println(row)
		err = row.Scan(&symbolInsert.Id, &symbolInsert.Preference,
			&symbolInsert.Fullname, &symbolInsert.Symbol)
		if err != nil {
			log.Panic(err)
		}

		err = tx.Commit(context.Background())
		if err != nil {
			log.Panic(err)
		}

		jsonQuery, err := json.Marshal(symbolInsert)
		if err != nil {
			log.Fatal(err)
		}

		return c.SendString(string(jsonQuery))

	})

	// REST API to register an order for a given asset.
	app.Post("/orders", func(c *fiber.Ctx) error {

		var orderInsert OrderBodyPost
		if err := c.BodyParser(&orderInsert); err != nil {
			fmt.Println(err)
		}
		fmt.Println(orderInsert)

		tx, err := dbpool.Begin(context.Background())
		if err != nil {
			log.Panic(err)
		}

		defer tx.Rollback(context.Background())

		var assetId string
		var brokerageId string
		queryAssetId := "SELECT id FROM asset WHERE symbol=$1"
		queryBrokerageId := "SELECT id FROM brokerage WHERE name=$1"
		tx.QueryRow(context.Background(), queryAssetId,
			orderInsert.Symbol).Scan(&assetId)
		tx.QueryRow(context.Background(), queryBrokerageId,
			orderInsert.Brokerage).Scan(&brokerageId)

		fmt.Println(brokerageId, assetId)
		var ordersInsert OrderApiReturn
		insertRow := "INSERT INTO orders(quantity, price, currency, order_type, date, asset_id, brokerage_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, quantity, price, currency, order_type, date;"

		row := tx.QueryRow(context.Background(), insertRow,
			orderInsert.Quantity, orderInsert.Price, orderInsert.Currency,
			orderInsert.OrderType, orderInsert.Date, assetId,
			brokerageId)
		err = row.Scan(&ordersInsert.Id, &orderInsert.Quantity,
			&orderInsert.Price, &orderInsert.Currency, &orderInsert.OrderType,
			&orderInsert.Date)
		if err != nil {
			log.Panic(err)
		}

		err = tx.Commit(context.Background())
		if err != nil {
			log.Panic(err)
		}

		jsonQuery, err := json.Marshal(ordersInsert)
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
