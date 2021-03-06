package postgresql

import (
	"context"
	"fmt"
	"log"
	"stockfyApi/entity"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

type AssetPostgres struct {
	dbpool PgxIface
}

type Asset struct {
	Id         string             `db:"id"`
	Preference *string            `db:"preference"`
	Fullname   string             `db:"fullname"`
	Symbol     string             `db:"symbol"`
	Sector     *entity.Sector     `db:"sector" json:",omitempty"`
	AssetType  *entity.AssetType  `db:"asset_type" json:",omitempty"`
	CreatedAt  time.Time          `db:"created_at" json:",omitempty"`
	UpdatedAt  time.Time          `db:"updated_at" json:",omitempty"`
	OrderInfo  *entity.OrderInfos `db:"orders_info" json:",omitempty"`
	OrdersList []Order            `db:"orders_list" json:",omitempty"`
}

type Order struct {
	Id        string            `db:"id" json:",omitempty"`
	Quantity  float64           `db:"quantity" json:",omitempty"`
	Price     float64           `db:"price" json:",omitempty"`
	Currency  string            `db:"currency" json:",omitempty"`
	OrderType string            `db:"order_type" json:",omitempty"`
	Date      string            `db:"date" json:",omitempty"`
	Brokerage *entity.Brokerage `db:"brokerage" json:",omitempty"`
	Asset     *entity.Asset     `db:"asset" json:",omitempty"`
	UserUid   string            `db:"user_uid" json:",omitempty"`
	CreatedAt time.Time         `db:"created_at" json:",omitempty"`
	UpdatedAt time.Time         `db:"updated_at" json:",omitempty"`
}

func NewAssetPostgres(db PgxIface) *AssetPostgres {
	return &AssetPostgres{
		dbpool: db,
	}
}

func (r *AssetPostgres) Create(assetInsert entity.Asset) entity.Asset {
	var ptrPreference string
	var assetReturn entity.Asset
	var row pgx.Row

	tx, err := r.dbpool.Begin(context.Background())
	if err != nil {
		log.Panic(err)
	}

	defer tx.Rollback(context.Background())

	insertRow := `
		INSERT INTO
			assets(preference, fullname, symbol, asset_type_id, sector_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, preference, fullname, symbol;
	`
	row = tx.QueryRow(context.Background(), insertRow,
		assetInsert.Preference, assetInsert.Fullname, assetInsert.Symbol,
		assetInsert.AssetType.Id, assetInsert.Sector.Id)

	err = row.Scan(&assetReturn.Id, &ptrPreference,
		&assetReturn.Fullname, &assetReturn.Symbol)
	if err != nil {
		log.Panic(err)
	}
	assetReturn.Preference = &ptrPreference

	err = tx.Commit(context.Background())
	if err != nil {
		log.Panic(err)
	}

	return assetReturn
}

func (r *AssetPostgres) Search(symbol string) ([]entity.Asset, error) {

	var symbolQuery []entity.Asset

	query := `
	SELECT
		a.id, symbol, preference, fullname,
		json_build_object(
			'id', aty.id,
			'type', aty."type",
			'name', aty."name",
			'country', aty.country
		) as asset_type,
		json_build_object(
			'id', s.id,
			'name', s."name"
		) as sector
	FROM assets as a
	INNER JOIN asset_types as aty
	ON aty.id = a.asset_type_id
	INNER JOIN sectors as s
	ON s.id = a.sector_id
	WHERE symbol=$1;
	`

	err := pgxscan.Select(context.Background(), r.dbpool, &symbolQuery, query,
		symbol)
	if err != nil {
		return nil, err
	}

	return symbolQuery, err

}

func (r *AssetPostgres) SearchByUser(symbol string, userUid string,
	orderType string) ([]entity.Asset, error) {

	var symbolQueryDb []Asset

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
			) as asset_type,
			json_build_object(
				'id', s.id,
				'name', s."name"
			) as sector
		FROM asset_users as au
		INNER JOIN assets as a
		ON a.id = au.asset_id
		INNER JOIN asset_types as aty
		ON aty.id = a.asset_type_id
		INNER JOIN sectors as s
		ON s.id = a.sector_id
		WHERE a.symbol=$1 and au.user_uid=$2
		GROUP BY a.symbol, a.id, a.preference, a.fullname, aty.id, aty."type",
		aty."name", aty.country, s.id, s."name";
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
			'id', s.id,
			'name', s."name"
		) as sector,
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
		INNER JOIN assets as a
		ON a.id = au.asset_id
		INNER JOIN asset_types as at
		ON a.asset_type_id = at.id
		INNER JOIN sectors as s
		ON s.id = a.sector_id
		INNER JOIN orders as o
		ON a.id = o.asset_id and au.user_uid = o.user_uid
		INNER JOIN brokerages as b
		ON o.brokerage_id = b.id
		WHERE a.symbol=$1 and au.user_uid =$2
		GROUP BY a.symbol, a.id, preference, a.fullname, at.type, at.id,
		at.name, at.country, s.id, s.name;
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
			'id', s.id,
			'name', s."name"
		) as sector,
		json_build_object(
			'totalQuantity', sum(o.quantity),
			'weightedAdjPrice', SUM(o.quantity * price)/SUM(o.quantity),
			'weightedAveragePrice', (
				SUM(o.quantity*o.price) FILTER(WHERE o.order_type = 'buy'))
				/(SUM(o.quantity) FILTER(WHERE o.order_type = 'buy')
			)
		) as orders_info
		FROM asset_users as au
		INNER JOIN assets as a
		ON a.id = au.asset_id
		INNER JOIN asset_types as aty
		ON a.asset_type_id = aty.id
		INNER JOIN sectors as s
		ON s.id = a.sector_id
		INNER JOIN orders as o
		ON a.id = o.asset_id and au.user_uid = o.user_uid
		INNER JOIN brokerages as b
		ON o.brokerage_id = b.id
		WHERE a.symbol=$1 and au.user_uid =$2
		GROUP BY a.symbol, a.id, preference, a.fullname, aty.type, aty.id,
		aty.name, aty.country, s.id, s.name;
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
		json_build_object(
			'id', s.id,
			'name', s."name"
		) as sector,
		json_agg(
			json_build_object(
				'id', o.id,
				'quantity', o.quantity,
				'price', o.price,
				'currency', o.currency,
				'ordertype', o.order_type,
				'date', o.date,
				'brokerage',
				json_build_object(
					'id', b.id,
					'name', b.name,
					'country', b.country
				)
			)
		) as orders_list
		FROM asset_users as au
		INNER JOIN assets as a
		ON a.id = au.asset_id
		INNER JOIN asset_types as at
		ON a.asset_type_id = at.id
		INNER JOIN sectors as s
		ON s.id = a.sector_id
		INNER JOIN orders as o
		ON a.id = o.asset_id and au.user_uid = o.user_uid
		INNER JOIN brokerages as b
		ON o.brokerage_id = b.id
		WHERE a.symbol=$1 and au.user_uid =$2
		GROUP BY a.symbol, a.id, preference, a.fullname, at.type, at.id,
		at.name, at.country, s.id, s.name;
		`
	}

	err := pgxscan.Select(context.Background(), r.dbpool, &symbolQueryDb, query,
		symbol, userUid)
	if err != nil {
		return nil, err
	}

	if symbolQueryDb == nil {
		return nil, nil
	}

	// Preliminary solution for the parsing date problem
	symbolQuery := []entity.Asset{
		{
			Id:         symbolQueryDb[0].Id,
			Symbol:     symbolQueryDb[0].Symbol,
			Preference: symbolQueryDb[0].Preference,
			Fullname:   symbolQueryDb[0].Fullname,
			Sector:     symbolQueryDb[0].Sector,
			AssetType:  symbolQueryDb[0].AssetType,
			OrderInfo:  symbolQueryDb[0].OrderInfo,
			OrdersList: func() []entity.Order {
				if symbolQueryDb[0].OrdersList == nil {
					return nil
				} else {
					var orderList []entity.Order
					layOut := "2006-01-02"
					for _, orderInfo := range symbolQueryDb[0].OrdersList {
						dateFormatted, _ := time.Parse(layOut, orderInfo.Date)
						order := entity.Order{
							Id:        orderInfo.Id,
							Quantity:  orderInfo.Quantity,
							Price:     orderInfo.Price,
							Currency:  orderInfo.Currency,
							OrderType: orderInfo.OrderType,
							Date:      dateFormatted,
							Brokerage: orderInfo.Brokerage,
						}

						orderList = append(orderList, order)
					}
					return orderList
				}
			}(),
		},
	}

	return symbolQuery, err
}

func (r *AssetPostgres) SearchPerAssetType(assetType string, country string,
	userUid string, withOrdersInfo bool) []entity.AssetType {

	var assetsPerAssetType []entity.AssetType

	var query string

	if !withOrdersInfo {
		query = `
		SELECT
			aty.id, aty.type, aty.country, aty.name,
			json_agg(
				json_build_object(
					'id', a.id,
					'symbol', a.symbol,
					'preference', a.preference,
					'fullname', a.fullname,
					'sector', json_build_object(
						'id', s.id,
						'name', s.name
					)
				)
			) as assets
		FROM asset_users as au
		INNER JOIN assets as a
		ON a.id = au.asset_id
		INNER JOIN asset_types as aty
		ON aty.id = a.asset_type_id
		INNER JOIN sectors as s
		ON s.id = a.sector_id
		WHERE au.user_uid=$1 and aty."type"=$2 and aty.country=$3
		GROUP BY aty.id, aty."type", aty."name", aty.country;
		`
	} else if withOrdersInfo {
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
				INNER JOIN assets as a
				ON a.id = au.asset_id
				INNER JOIN asset_types as aty
				ON aty.id = a.asset_type_id
				inner join sectors as s
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
	}

	err := pgxscan.Select(context.Background(), r.dbpool, &assetsPerAssetType,
		query, userUid, assetType, country)
	if err != nil {
		fmt.Println(err)
	}

	return assetsPerAssetType
}

func (r *AssetPostgres) SearchByOrderId(orderId string) []entity.Asset {
	var assetInfo []entity.Asset

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
	inner join assets as a
	on a.id = o.asset_id
	inner join asset_types as aty
	on aty.id = a.asset_type_id
	where o.id = $1;
	`

	err := pgxscan.Select(context.Background(), r.dbpool, &assetInfo,
		query, orderId)
	if err != nil {
		fmt.Println("SearchAssetByOrderId: ", err)
	}

	return assetInfo
}

func (r *AssetPostgres) Delete(assetId string) ([]entity.Asset, error) {
	var assetInfo []entity.Asset
	var err error

	queryDeleteAsset := `
	delete from assets as a
	where a.id = $1
	returning  a.id, a.symbol, a.preference, a.fullname;
	`

	err = pgxscan.Select(context.Background(), r.dbpool, &assetInfo,
		queryDeleteAsset, assetId)
	if err != nil {
		fmt.Println("entity.DeleteAsset: ", err)
	}

	return assetInfo, err
}
