package postgresql

import (
	"context"
	"errors"
	"stockfyApi/database"

	"github.com/georgysavva/scany/pgxscan"
	_ "github.com/lib/pq"
)

func (r *repo) FetchAssetType(fetchType string, args ...string) (
	[]database.AssetType, error) {

	var assetTypeQuery []database.AssetType
	var err error

	queryDefault := "SELECT id, type, name, country FROM assettype "

	if fetchType == "" {
		err = pgxscan.Select(context.Background(), r.dbpool, &assetTypeQuery,
			queryDefault)
		if err != nil {
			panic(err)
		}
	} else if fetchType == "SPECIFIC" {
		query := queryDefault + "where type=$1 and country=$2"
		err = pgxscan.Select(context.Background(), r.dbpool, &assetTypeQuery,
			query, args[0], args[1])
		if err != nil {
			panic(err)
		}
	} else if fetchType == "ONLYCOUNTRY" {
		query := queryDefault + "where country=$1"
		err = pgxscan.Select(context.Background(), r.dbpool, &assetTypeQuery,
			query, args[1])
		if err != nil {
			panic(err)
		}
	} else if fetchType == "ONLYTYPE" {
		query := queryDefault + "where type=$1"
		err = pgxscan.Select(context.Background(), r.dbpool, &assetTypeQuery,
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
