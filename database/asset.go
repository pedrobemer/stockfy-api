package database

import (
	"context"
	"fmt"
	"log"

	"github.com/georgysavva/scany/pgxscan"
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

func SearchAsset(dbpool *pgxpool.Pool, symbol string, orderType string) []AssetQueryReturn {

	var symbolQuery []AssetQueryReturn

	var query string
	if orderType == "" {
		query = `
		SELECT
			a.id, symbol, preference, fullname,
			json_build_object(
				'id', at.id,
				'type', at.type,
				'name', at.name,
				'country', at.country
			) as asset_type
		FROM asset as a
		INNER JOIN assettype as at
		ON a.asset_type_id = at.id
		WHERE a.symbol=$1
		GROUP BY a.symbol, a.id, preference, fullname, at.type, at.id, at.name,
		at.country;
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
			FROM asset as a
			INNER JOIN assettype as at
			ON a.asset_type_id = at.id
			INNER JOIN orders as o
			ON a.id = o.asset_id
			INNER JOIN brokerage as b
			ON o.brokerage_id = b.id
			WHERE a.symbol=$1
			GROUP BY a.symbol, a.id, preference, a.fullname, at.type, at.id,
				at.name, at.country;
			`
	} else if orderType == "ONLYINFO" {
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
					SUM(o.quantity*o.price) FILTER(WHERE o.order_type = 'buy'))/
					(SUM(o.quantity) FILTER(WHERE o.order_type = 'buy')
				)
			) as orders_info
			FROM asset as a
			INNER JOIN assettype as at
			ON a.asset_type_id = at.id
			INNER JOIN orders as o
			ON a.id = o.asset_id
			INNER JOIN brokerage as b
			ON o.brokerage_id = b.id
			WHERE a.symbol=$1
			GROUP BY a.symbol, a.id, preference, a.fullname, at.type, at.id,
				at.name, at.country;
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
			FROM asset as a
			INNER JOIN assettype as at
			ON a.asset_type_id = at.id
			INNER JOIN orders as o
			ON a.id = o.asset_id
			INNER JOIN brokerage as b
			ON o.brokerage_id = b.id
			WHERE a.symbol=$1
			GROUP BY a.symbol, a.id, preference, a.fullname, at.type, at.id,
				at.name, at.country;
			`

	}

	err := pgxscan.Select(context.Background(), dbpool, &symbolQuery, query,
		symbol)
	if err != nil {
		fmt.Println(err)
	}

	if symbolQuery[0].AssetType.Type != "ETF" &&
		symbolQuery[0].AssetType.Type != "FII" {
		var sector SectorApiReturn
		query = `
		select
			s.id,
			s.name
		from sector as s
		inner join asset as a
		on a.sector_id = s.id
		where a.symbol = $1;
		`
		row := dbpool.QueryRow(context.Background(), query, symbol)
		fmt.Println(row)
		err = row.Scan(&sector.Id, &sector.Name)
		if err != nil {
			fmt.Println(err)
		}
		symbolQuery[0].Sector = &sector
		// fmt.Println(symbolQuery)
	}

	return symbolQuery
}

func SearchAssetsPerAssetType(dbpool pgxpool.Pool, assetType string,
	country string, withOrdersInfo bool) []AssetTypeApiReturn {

	var assetsPerAssetType []AssetTypeApiReturn

	var query string

	if !withOrdersInfo && assetType != "ETF" && assetType != "FII" {
		query = `
		SELECT
			at.id, at.type, at.country, at.name,
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
		FROM assettype as at
		INNER JOIN asset as a
		ON at.id = a.asset_type_id
		INNER JOIN sector as s
		ON a.sector_id = s.id
		WHERE at.type = $1 and at.country = $2
		GROUP BY at.id, at.type, at.country, at.name;
		`
	} else if !withOrdersInfo && (assetType == "ETF" || assetType == "FII") {
		query = `
		SELECT
			at.id, at.type, at.country, at.name,
			json_agg(
				json_build_object(
					'id', a.id,
					'symbol', a.symbol,
					'preference', a.preference,
					'fullname', a.fullname
				)
			) as assets
		FROM assettype as at
		INNER JOIN asset as a
		ON at.id = a.asset_type_id
		WHERE at.type = $1 and at.country = $2
		GROUP BY at.id, at.type, at.country, at.name;
		`
	} else if withOrdersInfo && assetType != "ETF" && assetType != "FII" {
		query = `
		select
			f_query.atid as id, f_query.attype as type, f_query.atname as name,
			f_query.atcountry as country,
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
		from (
			select
				a2.id, a2.symbol, a2.preference, a2.fullname,
				a2.atid, a2.attype, a2.atname, a2.atcountry,
				json_build_object('id', a2.sid, 'name', a2.sname) as sector,
				json_build_object(
					'totalQuantity', sum(o.quantity),
					'weightedAdjPrice', SUM(o.quantity * price)/SUM(o.quantity),
					'weightedAveragePrice', (
						SUM(o.quantity*o.price)
						FILTER(WHERE o.order_type = 'buy'))
						/(SUM(o.quantity) FILTER(WHERE o.order_type = 'buy'))
			) as order_info
			from (
				select
					a.id, a.symbol, a.preference, a.fullname,
					s.id as sid, s.name as sname,
					at.id as atid, at.type as attype, at.name as atname,
					at.country as atcountry
				from asset as a
				inner join assettype as at
				on at.id = a.asset_type_id
				inner join sector as s
				on s.id = a.sector_id
				where at.type = $1 and at.country = $2
			) a2
			inner join orders as o
			on o.asset_id = a2.id
			group by a2.id, a2.symbol, a2.preference, a2.fullname, a2.sid,
				a2.sname, a2.atid, a2.attype, a2.atname, a2.atcountry
		) as f_query
		group by f_query.atid, f_query.attype, f_query.atcountry, f_query.atname
		`
	} else if withOrdersInfo && (assetType == "ETF" || assetType == "FII") {
		query = `
		select
			f_query.atid as id,
			f_query.attype as type,
			f_query.atname as name,
			f_query.atcountry as country,
			json_agg(
				json_build_object(
					'id', f_query.id,
					'symbol', f_query.symbol,
					'preference', f_query.preference,
					'fullname', f_query.fullname,
					'orderInfo', f_query.order_info
				)
			) as assets
		from (
			select
				a2.id, a2.symbol, a2.preference, a2.fullname,
				a2.atid, a2.attype, a2.atname, a2.atcountry,
				json_build_object(
					'totalQuantity', sum(o.quantity),
					'weightedAdjPrice', SUM(o.quantity * price)/SUM(o.quantity),
					'weightedAveragePrice', (
						SUM(o.quantity*o.price) FILTER(WHERE o.order_type = 'buy'))
						/(SUM(o.quantity) FILTER(WHERE o.order_type = 'buy'))
			) as order_info
			from (
				select
					a.id, a.symbol, a.preference, a.fullname,
					at.id as atid, at.type as attype, at.name as atname,
					at.country as atcountry
				from asset as a
				inner join assettype as at
				on at.id = a.asset_type_id
				where at.type = $1 and at.country = $2
			) a2
			inner join orders as o
			on o.asset_id = a2.id
			group by a2.id, a2.symbol, a2.preference, a2.fullname,
				a2.atid, a2.attype, a2.atname, a2.atcountry
		) as f_query
		group by f_query.atid, f_query.attype, f_query.atcountry, f_query.atname
		`
	}
	err := pgxscan.Select(context.Background(), &dbpool, &assetsPerAssetType,
		query, assetType, country)
	if err != nil {
		fmt.Println(err)
	}

	return assetsPerAssetType
}
