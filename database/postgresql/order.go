package postgresql

import (
	"context"
	"fmt"
	"log"
	"stockfyApi/database"

	"github.com/georgysavva/scany/pgxscan"
	_ "github.com/lib/pq"
)

func (r *repo) CreateOrder(orderInsert database.Order) database.Order {

	var orderReturn database.Order

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

func (r *repo) SearchOrdersFromAssetUser(assetId string, userUid string) (
	[]database.Order, error) {
	var ordersReturn []database.Order

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
		fmt.Println("database.SearchOrdersFromAssetUser: ", err)
	}

	return ordersReturn, err
}

func (r *repo) DeleteOrderFromUser(id string, userUid string) string {
	var orderId string

	query := `
	delete from orders as o
	where o.id = $1 and o.user_uid = $2
	returning o.id
	`
	row := r.dbpool.QueryRow(context.Background(), query, id, userUid)
	err := row.Scan(&orderId)
	if err != nil {
		fmt.Println("database.DeleteOrder: ", err)
	}

	return orderId
}

func (r *repo) DeleteOrdersFromAsset(symbolId string) []database.Order {
	var ordersId []database.Order

	queryDeleteOrders := `
	delete from orders as o
	where o.asset_id = $1
	returning o.id;
	`

	err := pgxscan.Select(context.Background(), r.dbpool, &ordersId,
		queryDeleteOrders, symbolId)
	if err != nil {
		fmt.Println("database.DeleteOrders: ", err)
	}

	return ordersId
}

func (r *repo) DeleteOrdersFromAssetUser(assetId string, userUid string) (
	[]database.Order, error) {
	var ordersId []database.Order

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
		fmt.Println("database.DeleteOrders: ", err)
	}

	return ordersId, err

}

func (r *repo) UpdateOrderFromUser(orderUpdate database.Order) []database.Order {
	var orderInfo []database.Order

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
		fmt.Println("database.UpdateOrder: ", err)
	}

	return orderInfo
}
