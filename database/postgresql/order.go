package postgresql

import (
	"context"
	"fmt"
	"log"
	"stockfyApi/entity"

	"github.com/georgysavva/scany/pgxscan"
	_ "github.com/lib/pq"
)

type OrderPostgres struct {
	dbpool PgxIface
}

func NewOrderPostgres(db *PgxIface) *OrderPostgres {
	return &OrderPostgres{
		dbpool: *db,
	}
}

func (r *OrderPostgres) Create(orderInsert entity.Order) entity.Order {

	var orderReturn entity.Order

	tx, err := r.dbpool.Begin(context.Background())
	if err != nil {
		log.Panic(err)
	}

	defer tx.Rollback(context.Background())

	insertRow := `
	WITH inserted as (
		INSERT INTO
			orders(quantity, price, currency, order_type, date, asset_id,
				brokerage_id, user_uid
			)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, quantity, price, currency, order_type, date, brokerage_id
	)
	SELECT
		inserted.id, inserted.quantity, inserted.price, inserted.currency,
		inserted.order_type, inserted.date,
		json_build_object(
			'id', b.id,
			'name', b.name,
			'country', b.country
		) as brokerage
	FROM inserted
	INNER JOIN brokerage as b
	ON inserted.brokerage_id = b.id;
	`

	row := tx.QueryRow(context.Background(), insertRow,
		orderInsert.Quantity, orderInsert.Price, orderInsert.Currency,
		orderInsert.OrderType, orderInsert.Date, orderInsert.Asset.Id,
		orderInsert.Brokerage.Id, orderInsert.UserUid)
	err = row.Scan(&orderReturn.Id, &orderReturn.Quantity,
		&orderReturn.Price, &orderReturn.Currency,
		&orderReturn.OrderType, &orderReturn.Date, &orderReturn.Brokerage)
	if err != nil {
		log.Panic(err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		log.Panic(err)
	}

	return orderReturn
}

func (r *OrderPostgres) SearchFromAssetUser(assetId string, userUid string) (
	[]entity.Order, error) {
	var ordersReturn []entity.Order

	query := `
	SELECT
		o.id, quantity, price, currency, order_type, date,
		json_build_object(
			'id', b.id,
			'name', b."name",
			'country', b.country
		) as brokerage
	FROM orders as o
	INNER JOIN brokerage as b
	ON b.id = o.brokerage_id
	WHERE asset_id = $1 and user_uid = $2;
	`

	err := pgxscan.Select(context.Background(), r.dbpool, &ordersReturn, query,
		assetId, userUid)
	if err != nil {
		fmt.Println("entity.SearchOrdersFromAssetUser: ", err)
	}

	return ordersReturn, err
}

func (r *OrderPostgres) DeleteFromUser(id string, userUid string) string {
	var orderId string

	query := `
	delete from orders as o
	where o.id = $1 and o.user_uid = $2
	returning o.id
	`
	row := r.dbpool.QueryRow(context.Background(), query, id, userUid)
	err := row.Scan(&orderId)
	if err != nil {
		fmt.Println("entity.DeleteOrder: ", err)
	}

	return orderId
}

func (r *OrderPostgres) DeleteFromAsset(symbolId string) []entity.Order {
	var ordersId []entity.Order

	queryDeleteOrders := `
	delete from orders as o
	where o.asset_id = $1
	returning o.id;
	`

	err := pgxscan.Select(context.Background(), r.dbpool, &ordersId,
		queryDeleteOrders, symbolId)
	if err != nil {
		fmt.Println("entity.DeleteOrders: ", err)
	}

	return ordersId
}

func (r *OrderPostgres) DeleteFromAssetUser(assetId string, userUid string) (
	[]entity.Order, error) {
	var ordersId []entity.Order

	queryDeleteOrders := `
	with deleted as (
	delete from orders as o
	where o.asset_id = $1 and o.user_uid = $2
	returning o.id
	)
	select
		deleted.id,
		jsonb_build_object(
			'id', ast.id,
			'symbol', ast.symbol
		) as asset
	from deleted
	inner join asset as ast
	on ast.id = deleted.asset_id;
	`
	err := pgxscan.Select(context.Background(), r.dbpool, &ordersId,
		queryDeleteOrders, assetId, userUid)
	if err != nil {
		fmt.Println("entity.DeleteOrders: ", err)
	}

	return ordersId, err

}

func (r *OrderPostgres) UpdateFromUser(orderUpdate entity.Order) []entity.Order {
	var orderInfo []entity.Order

	query := `
	update orders as o
	set quantity = $3,
		price = $4,
		order_type = $5,
		"date" = $6
	where o.id = $1 and o.user_uid = $2
	returning o.id, o.quantity, o.price, o."date", o.order_type;
	`
	err := pgxscan.Select(context.Background(), r.dbpool, &orderInfo,
		query, orderUpdate.Id, orderUpdate.UserUid, orderUpdate.Quantity,
		orderUpdate.Price, orderUpdate.OrderType, orderUpdate.Date)
	if err != nil {
		fmt.Println("entity.UpdateOrder: ", err)
	}

	return orderInfo
}
