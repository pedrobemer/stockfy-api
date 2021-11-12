package postgresql

import (
	"context"
	"errors"
	"stockfyApi/entity"

	"github.com/georgysavva/scany/pgxscan"
	_ "github.com/lib/pq"
)

type AssetTypePostgres struct {
	dbpool PgxIface
}

func NewAssetTypePostgres(db PgxIface) *AssetTypePostgres {
	return &AssetTypePostgres{
		dbpool: db,
	}
}

func (r *AssetTypePostgres) Search(searchType string, name string,
	country string) ([]entity.AssetType, error) {

	var assetTypeQuery []entity.AssetType
	var err error

	queryDefault := "SELECT id, type, name, country FROM assettypes "

	if searchType == "" {
		err = pgxscan.Select(context.Background(), r.dbpool, &assetTypeQuery,
			queryDefault)
		if err != nil {
			panic(err)
		}
	} else if searchType == "SPECIFIC" {
		query := queryDefault + "where type=$1 and country=$2"
		err = pgxscan.Select(context.Background(), r.dbpool, &assetTypeQuery,
			query, name, country)
		if err != nil {
			panic(err)
		}
	} else if searchType == "ONLYCOUNTRY" {
		query := queryDefault + "where country=$1"
		err = pgxscan.Select(context.Background(), r.dbpool, &assetTypeQuery,
			query, country)
		if err != nil {
			panic(err)
		}
	} else if searchType == "ONLYTYPE" {
		query := queryDefault + "where type=$1"
		err = pgxscan.Select(context.Background(), r.dbpool, &assetTypeQuery,
			query, name)
		if err != nil {
			panic(err)
		}
	}

	if assetTypeQuery == nil {
		err = errors.New("FetchAssetType: There is no asset type with this specifications")
	}

	return assetTypeQuery, err
}
