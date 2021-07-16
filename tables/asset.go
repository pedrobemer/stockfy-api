package tables

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

func CreateAsset(dbpool pgxpool.Pool, assetInsert AssetInsert,
	assetTypeId string, sectorId string) AssetApiReturn {

	var preference string
	if assetInsert.Country == "BR" && assetInsert.AssetType == "STOCK" {
		fmt.Println(len(assetInsert.Symbol[len(assetInsert.Symbol)-1:]))
		switch assetInsert.Symbol[len(assetInsert.Symbol)-1:] {
		case "3":
			preference = "ON"
			break
		case "4":
			// fmt.Println("entrou")
			preference = "PN"
			break
		case "11":
			preference = "UNIT"
			break
		default:
			preference = ""
			break
		}
	}
	fmt.Println(preference)

	tx, err := dbpool.Begin(context.Background())
	if err != nil {
		log.Panic(err)
	}

	defer tx.Rollback(context.Background())

	fmt.Println(assetTypeId, sectorId)
	var symbolInsert AssetApiReturn
	var insertRow string
	var row pgx.Row
	if sectorId != "" {
		insertRow = "INSERT INTO asset(preference, fullname, symbol, asset_type_id, sector_id) VALUES ($1, $2, $3, $4, $5) RETURNING id, preference, fullname, symbol;"
		row = tx.QueryRow(context.Background(), insertRow,
			preference, assetInsert.Fullname, assetInsert.Symbol,
			assetTypeId, sectorId)
	} else {
		insertRow = "INSERT INTO asset(preference, fullname, symbol, asset_type_id) VALUES ($1, $2, $3, $4) RETURNING id, preference, fullname, symbol;"
		row = tx.QueryRow(context.Background(), insertRow,
			preference, assetInsert.Fullname, assetInsert.Symbol,
			assetTypeId)
	}

	err = row.Scan(&symbolInsert.Id, &symbolInsert.Preference,
		&symbolInsert.Fullname, &symbolInsert.Symbol)
	if err != nil {
		log.Panic(err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		log.Panic(err)
	}

	return symbolInsert
}

func SearchAsset(dbpool pgxpool.Pool, symbol string) string {
	var assetId string

	fetchAssetId := "SELECT id FROM asset WHERE symbol=$1"
	row := dbpool.QueryRow(context.Background(), fetchAssetId, symbol)
	row.Scan(&assetId)

	return assetId
}
