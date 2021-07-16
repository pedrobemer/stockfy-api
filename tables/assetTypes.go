package tables

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

func FetchAssetType(dbpool pgxpool.Pool, specificFetch bool,
	args ...string) ([]AssetTypeApiReturn, error) {

	var assetTypeQuery []AssetTypeApiReturn
	var err error

	queryDefault := "SELECT id, type, name, country FROM assettype "

	if !specificFetch {
		err = pgxscan.Select(context.Background(), &dbpool, &assetTypeQuery,
			queryDefault)
		if err != nil {
			panic(err)
		}
	} else {
		query := queryDefault + "where type=$1 and country=$2"
		err = pgxscan.Select(context.Background(), &dbpool, &assetTypeQuery,
			query, args[0], args[1])
		if err != nil {
			panic(err)
		}
	}

	return assetTypeQuery, err
}
