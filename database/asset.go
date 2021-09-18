package database

import (
	"context"
	"fmt"
	"log"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

func CreateAsset(dbpool PgxIface, assetInsert AssetInsert,
	assetTypeId string, sectorId string) AssetApiReturn {

	var preference string
	if assetInsert.Country == "BR" && assetInsert.AssetType == "STOCK" {
		switch assetInsert.Symbol[len(assetInsert.Symbol)-1:] {
		case "3":
			preference = "ON"
			break
		case "4":
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

	tx, err := dbpool.Begin(context.Background())
	if err != nil {
		log.Panic(err)
	}

	defer tx.Rollback(context.Background())

	var symbolInsert AssetApiReturn
	var insertRow string
	var row pgx.Row
	if sectorId != "" {
		insertRow = `
		INSERT INTO
			asset(preference, fullname, symbol, asset_type_id, sector_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, preference, fullname, symbol;
		`
		row = tx.QueryRow(context.Background(), insertRow,
			preference, assetInsert.Fullname, assetInsert.Symbol,
			assetTypeId, sectorId)
	} else {
		insertRow = `
		INSERT INTO
			asset(preference, fullname, symbol, asset_type_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, preference, fullname, symbol;
		`
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

func SearchAsset(dbpool PgxIface, symbol string) ([]AssetQueryReturn, error) {

	var symbolQuery []AssetQueryReturn

	query := `
	SELECT
		a.id, symbol, preference, fullname,
		json_build_object(
			'id', aty.id,
			'type', aty."type",
			'name', aty."name",
			'country', aty.country
		) as asset_type
	FROM asset as a
	INNER JOIN assettype as aty
	ON aty.id = a.asset_type_id
	WHERE symbol=$1;
	`

	err := pgxscan.Select(context.Background(), dbpool, &symbolQuery, query,
		symbol)
	if err != nil {
		return symbolQuery, err
	}
	if symbolQuery == nil {
		return symbolQuery, err
	}

	if symbolQuery[0].AssetType.Type != "ETF" &&
		symbolQuery[0].AssetType.Type != "FII" {
		var err error
		var sector []SectorApiReturn
		sector, err = FetchSectorByAsset(dbpool, symbol)
		if err != nil {
			return symbolQuery, err
		}
		symbolQuery[0].Sector = &sector[0]
	}

	return symbolQuery, err

}

func SearchAssetByUser(dbpool PgxIface, symbol string, userUid string,
	orderType string) ([]AssetQueryReturn, error) {

	var symbolQuery []AssetQueryReturn

	var query string
	if orderType == "" {
		query = `
		SELECT
			a.id, symbol, preference, fullname,
			json_build_object(
				'id', aty.id,
				'type', aty."type",
				'name', aty."name",
				'country', aty.country
			) as asset_type
		FROM asset_users as au
		INNER JOIN asset as a
		ON a.id = au.asset_id
		INNER JOIN assettype as aty
		ON aty.id = a.asset_type_id
		WHERE au.user_uid=$2 and a.symbol=$1
		GROUP BY a.symbol, a.id, a.preference, a.fullname, aty.id, aty."type",
		aty."name", aty.country;
		`
	} else if orderType == "ALL" {
		query = `
		SELECT
			a.id, symbol, preference, a.fullname,
		json_build_object(
			'id', at.id,
			'type', at.type,
			'name', at.name,
			'country', at.country
		) as asset_type,
		json_build_object(
			'totalQuantity', sum(o.quantity),
			'weightedAdjPrice', SUM(o.quantity * price)/SUM(o.quantity),
			'weightedAveragePrice', (
				SUM(o.quantity*o.price) FILTER(WHERE o.order_type = 'buy'))
				/(SUM(o.quantity) FILTER(WHERE o.order_type = 'buy')
			)
		) as orders_info,
		json_agg(
			json_build_object(
				'id', o.id,
				'quantity', o.quantity,
				'price', o.price,
				'currency', o.currency,
				'ordertype', o.order_type,
				'date', date,
				'brokerage',
				json_build_object(
					'id', b.id,
					'name', b.name,
					'country', b.country
				)
			)
		) as orders_list
		FROM asset_users as au
		INNER JOIN asset as a
		ON a.id = au.asset_id
		INNER JOIN assettype as at
		ON a.asset_type_id = at.id
		INNER JOIN orders as o
		ON a.id = o.asset_id and au.user_uid = o.user_uid
		INNER JOIN brokerage as b
		ON o.brokerage_id = b.id
		WHERE a.symbol=$1 and au.user_uid =$2
		GROUP BY a.symbol, a.id, preference, a.fullname, at.type, at.id,
		at.name, at.country;
		`
	} else if orderType == "ONLYINFO" {
		query = `
		SELECT
			a.id, symbol, preference, a.fullname,
		json_build_object(
			'id', aty.id,
			'type', aty.type,
			'name', aty.name,
			'country', aty.country
		) as asset_type,
		json_build_object(
			'totalQuantity', sum(o.quantity),
			'weightedAdjPrice', SUM(o.quantity * price)/SUM(o.quantity),
			'weightedAveragePrice', (
				SUM(o.quantity*o.price) FILTER(WHERE o.order_type = 'buy'))
				/(SUM(o.quantity) FILTER(WHERE o.order_type = 'buy')
			)
		) as orders_info
		FROM asset_users as au
		INNER JOIN asset as a
		ON a.id = au.asset_id
		INNER JOIN assettype as aty
		ON a.asset_type_id = aty.id
		INNER JOIN orders as o
		ON a.id = o.asset_id and au.user_uid = o.user_uid
		INNER JOIN brokerage as b
		ON o.brokerage_id = b.id
		WHERE a.symbol=$1 and au.user_uid =$2
		GROUP BY a.symbol, a.id, preference, a.fullname, aty.type, aty.id,
		aty.name, aty.country;
		`
	} else if orderType == "ONLYORDERS" {
		query = `
		SELECT
			a.id, symbol, preference, a.fullname,
		json_build_object(
			'id', at.id,
			'type', at.type,
			'name', at.name,
			'country', at.country
		) as asset_type,
		json_agg(
			json_build_object(
				'id', o.id,
				'quantity', o.quantity,
				'price', o.price,
				'currency', o.currency,
				'ordertype', o.order_type,
				'date', date,
				'brokerage',
				json_build_object(
					'id', b.id,
					'name', b.name,
					'country', b.country
				)
			)
		) as orders_list
		FROM asset_users as au
		INNER JOIN asset as a
		ON a.id = au.asset_id
		INNER JOIN assettype as at
		ON a.asset_type_id = at.id
		INNER JOIN orders as o
		ON a.id = o.asset_id and au.user_uid = o.user_uid
		INNER JOIN brokerage as b
		ON o.brokerage_id = b.id
		WHERE a.symbol=$1 and au.user_uid =$2
		GROUP BY a.symbol, a.id, preference, a.fullname, at.type, at.id,
		at.name, at.country;
		`
	}

	err := pgxscan.Select(context.Background(), dbpool, &symbolQuery, query,
		symbol, userUid)
	if err != nil {
		return symbolQuery, err
	}

	if symbolQuery[0].AssetType.Type != "ETF" &&
		symbolQuery[0].AssetType.Type != "FII" {
		var err error
		var sector []SectorApiReturn
		sector, err = FetchSectorByAsset(dbpool, symbol)
		if err != nil {
			return symbolQuery, err
		}
		symbolQuery[0].Sector = &sector[0]
	}

	return symbolQuery, err
}

func SearchAssetsPerAssetType(dbpool PgxIface, assetType string,
	country string, userUid string, withOrdersInfo bool) []AssetTypeApiReturn {

	var assetsPerAssetType []AssetTypeApiReturn

	var query string

	if !withOrdersInfo && assetType != "ETF" && assetType != "FII" {
		query = `
		SELECT
			aty.id, aty.type, aty.country, aty.name,
			json_agg(
				json_build_object(
					'id', a.id,
					'symbol', a.symbol,
					'preference', a.preference,
					'fullname', a.fullname, 'sector',
					json_build_object(
						'id', s.id,
						'name', s.name
					)
				)
			) as assets
		FROM asset_users as au
		INNER JOIN asset as a
		ON a.id = au.asset_id
		INNER JOIN assettype as aty
		ON aty.id = a.asset_type_id
		INNER JOIN sector as s
		ON s.id = a.sector_id
		WHERE au.user_uid=$1 and aty."type"=$2 and aty.country=$3
		GROUP BY aty.id, aty."type", aty."name", aty.country;
		`
	} else if !withOrdersInfo && (assetType == "ETF" || assetType == "FII") {
		query = `
		SELECT
			aty.id, aty.type, aty.country, aty.name,
			json_agg(
				json_build_object(
					'id', a.id,
					'symbol', a.symbol,
					'preference', a.preference,
					'fullname', a.fullname
				)
			) as assets
		FROM asset_users as au
		INNER JOIN asset as a
		ON a.id = au.asset_id
		INNER JOIN assettype as aty
		ON aty.id = a.asset_type_id
		WHERE au.user_uid=$1 and aty."type"=$2 and aty.country=$3
		GROUP BY aty.id, aty."type", aty."name", aty.country;
		`
	} else if withOrdersInfo && assetType != "ETF" && assetType != "FII" {
		query = `
		SELECT
			f_query.at_id as id, f_query.at_type as type, f_query.at_name as name,
			f_query.at_country as country,
			json_agg(
				json_build_object(
					'id', f_query.id,
					'symbol', f_query.symbol,
					'preference', f_query.preference,
					'fullname', f_query.fullname,
					'sector', f_query.sector,
					'orderInfo', f_query.order_info
				)
			) as assets
		FROM (
			SELECT
				valid_assets.id, valid_assets.symbol, valid_assets.preference,
				valid_assets.fullname, valid_assets.at_id, valid_assets.at_type,
				valid_assets.at_name, valid_assets.at_country,
				json_build_object(
					'id', valid_assets.s_id,
					'name', valid_assets.s_name
				) as sector,
				json_build_object(
					'totalQuantity', sum(o.quantity),
					'weightedAdjPrice', SUM(o.quantity * price)/SUM(o.quantity),
					'weightedAveragePrice', (
						SUM(o.quantity*o.price)
						FILTER(WHERE o.order_type = 'buy'))
						/(SUM(o.quantity) FILTER(WHERE o.order_type = 'buy'))
				) as order_info
			FROM (
				select
					a.id, a.symbol, a.preference, a.fullname, s.id as s_id,
					s."name" as s_name, aty.id as at_id, aty."type" as at_type,
					aty."name" as at_name, aty.country as at_country
				FROM asset_users as au
				INNER JOIN asset as a
				ON a.id = au.asset_id
				INNER JOIN assettype as aty
				ON aty.id = a.asset_type_id
				inner join sector as s
				on s.id = a.sector_id
				WHERE au.user_uid=$1 and aty."type"=$2 and aty.country=$3
				GROUP BY a.symbol, a.id, a.preference, a.fullname, aty.id, aty."type",
				aty."name", aty.country, s.id, s."name"
			) valid_assets
			INNER JOIN orders as o
			ON o.asset_id = valid_assets.id
			WHERE o.user_uid = $1
			GROUP BY valid_assets.id, valid_assets.symbol,
			valid_assets.preference, valid_assets.fullname, valid_assets.s_id,
			valid_assets.s_name, valid_assets.at_id, valid_assets.at_type,
			valid_assets.at_name, valid_assets.at_country
		) as f_query
		GROUP BY f_query.at_id, f_query.at_type, f_query.at_country,
		f_query.at_name;
		`
	} else if withOrdersInfo && (assetType == "ETF" || assetType == "FII") {
		query = `
		SELECT
			f_query.at_id as id, f_query.at_type as type, f_query.at_name as name,
			f_query.at_country as country,
			json_agg(
				json_build_object(
					'id', f_query.id,
					'symbol', f_query.symbol,
					'preference', f_query.preference,
					'fullname', f_query.fullname,
					'orderInfo', f_query.order_info
				)
			) as assets
		FROM (
			SELECT
				valid_assets.id, valid_assets.symbol, valid_assets.preference,
				valid_assets.fullname, valid_assets.at_id, valid_assets.at_type,
				valid_assets.at_name, valid_assets.at_country,
				json_build_object(
					'totalQuantity', sum(o.quantity),
					'weightedAdjPrice', SUM(o.quantity * price)/SUM(o.quantity),
					'weightedAveragePrice', (
						SUM(o.quantity*o.price)
						FILTER(WHERE o.order_type = 'buy'))
						/(SUM(o.quantity) FILTER(WHERE o.order_type = 'buy'))
				) as order_info
			FROM (
				SELECT
					a.id, a.symbol, a.preference, a.fullname, aty.id as at_id,
					aty."type" as at_type, aty."name" as at_name,
					aty.country as at_country
				FROM asset_users as au
				INNER JOIN asset as a
				ON a.id = au.asset_id
				INNER JOIN assettype as aty
				ON aty.id = a.asset_type_id
				WHERE au.user_uid=$1 and aty."type"=$2 and aty.country=$3
				GROUP BY a.symbol, a.id, a.preference, a.fullname, aty.id, aty."type",
				aty."name", aty.country
			) valid_assets
			INNER JOIN orders as o
			ON o.asset_id = valid_assets.id
			WHERE o.user_uid =$1
			GROUP BY valid_assets.id, valid_assets.symbol,
			valid_assets.preference, valid_assets.fullname, valid_assets.at_id,
			valid_assets.at_type, valid_assets.at_name, valid_assets.at_country
		) as f_query
		GROUP BY f_query.at_id, f_query.at_type, f_query.at_country,
		f_query.at_name;
		`
	}

	err := pgxscan.Select(context.Background(), dbpool, &assetsPerAssetType,
		query, userUid, assetType, country)
	if err != nil {
		fmt.Println(err)
	}

	return assetsPerAssetType
}

func SearchAssetByOrderId(dpbool PgxIface, orderId string) []AssetQueryReturn {
	var assetInfo []AssetQueryReturn

	query := `
	select
		a.id, a.preference , a.symbol,
		json_build_object(
			'id', aty.id,
			'type', aty."type",
			'name', aty."name",
			'country', aty.country
		) as asset_type
	from orders as o
	inner join asset as a
	on a.id = o.asset_id
	inner join assettype as aty
	on aty.id = a.asset_type_id
	where o.id = $1;
	`

	err := pgxscan.Select(context.Background(), dpbool, &assetInfo,
		query, orderId)
	if err != nil {
		fmt.Println("SearchAssetByOrderId: ", err)
	}

	return assetInfo
}

func DeleteAsset(dbpool PgxIface, assetId string) []AssetQueryReturn {
	var assetInfo []AssetQueryReturn
	var err error

	queryDeleteAsset := `
	delete from asset as a
	where a.id = $1
	returning  a.id, a.symbol, a.preference, a.fullname;
	`

	err = pgxscan.Select(context.Background(), dbpool, &assetInfo,
		queryDeleteAsset, assetId)
	if err != nil {
		fmt.Println("database.DeleteAsset: ", err)
	}

	return assetInfo
}
