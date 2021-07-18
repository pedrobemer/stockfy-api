package database

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

func FetchAssetType(dbpool pgxpool.Pool, fetchType string,
	args ...string) ([]AssetTypeApiReturn, error) {

	var assetTypeQuery []AssetTypeApiReturn
	var err error

	queryDefault := "SELECT id, type, name, country FROM assettype "

	if fetchType == "" {
		err = pgxscan.Select(context.Background(), &dbpool, &assetTypeQuery,
			queryDefault)
		if err != nil {
			panic(err)
		}
	} else if fetchType == "SPECIFIC" {
		query := queryDefault + "where type=$1 and country=$2"
		err = pgxscan.Select(context.Background(), &dbpool, &assetTypeQuery,
			query, args[0], args[1])
		if err != nil {
			panic(err)
		}
	} else if fetchType == "ONLYCOUNTRY" {
		query := queryDefault + "where country=$1"
		err = pgxscan.Select(context.Background(), &dbpool, &assetTypeQuery,
			query, args[1])
		if err != nil {
			panic(err)
		}
	} else if fetchType == "ONLYTYPE" {
		query := queryDefault + "where type=$1"
		err = pgxscan.Select(context.Background(), &dbpool, &assetTypeQuery,
			query, args[0])
		if err != nil {
			panic(err)
		}
	}

	if assetTypeQuery == nil {
		err = errors.New("FetchAssetType: There is no asset type with this specifications")
	}

	return assetTypeQuery, err
}
