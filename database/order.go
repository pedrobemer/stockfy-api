package database

import (
	"context"
	"fmt"
	"log"

	"github.com/georgysavva/scany/pgxscan"
	_ "github.com/lib/pq"
)

func CreateOrder(dbpool PgxIface, orderInsert OrderBodyPost, assetId string,
	brokerageId string) OrderApiReturn {

	var orderReturn OrderApiReturn

	tx, err := dbpool.Begin(context.Background())
	if err != nil {
		log.Panic(err)
	}

	defer tx.Rollback(context.Background())

	insertRow := `
	WITH inserted as (
		INSERT INTO
			orders(quantity, price, currency, order_type, date, asset_id,
				brokerage_id
			)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
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
		orderInsert.OrderType, orderInsert.Date, assetId,
		brokerageId)
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

func DeleteOrder(dbpool PgxIface, id string) string {
	var orderId string

	query := `
	delete from orders as o
	where o.id = $1
	returning o.id
	`
	row := dbpool.QueryRow(context.Background(), query, id)
	err := row.Scan(&orderId)
	if err != nil {
		fmt.Println("database.DeleteOrder: ", err)
	}

	return orderId
}

func DeleteOrders(dbpool PgxIface, symbolId string) []OrderApiReturn {
	var ordersId []OrderApiReturn

	queryDeleteOrders := `
	delete from orders as o
	where o.asset_id = $1
	returning o.id;
	`

	err := pgxscan.Select(context.Background(), dbpool, &ordersId,
		queryDeleteOrders, symbolId)
	if err != nil {
		fmt.Println("database.DeleteOrders: ", err)
	}

	return ordersId
}

func UpdateOrder(dbpool PgxIface, orderUpdate OrderBodyPost) []OrderApiReturn {
	var orderInfo []OrderApiReturn

	query := `
	update orders as o
	set quantity = $2,
		price = $3,
		order_type = $4,
		"date" = $5
	where o.id = $1
	returning o.id, o.quantity, o.price, o."date", o.order_type;
	`
	err := pgxscan.Select(context.Background(), dbpool, &orderInfo,
		query, orderUpdate.Id, orderUpdate.Quantity, orderUpdate.Price,
		orderUpdate.OrderType, orderUpdate.Date)
	if err != nil {
		fmt.Println("database.UpdateOrder: ", err)
	}

	return orderInfo
}
