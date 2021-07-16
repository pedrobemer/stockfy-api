package tables

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

func CreateOrder(dbpool pgxpool.Pool, orderInsert OrderBodyPost, assetId string,
	brokerageId string) OrderApiReturn {

	var orderReturn OrderApiReturn

	tx, err := dbpool.Begin(context.Background())
	if err != nil {
		log.Panic(err)
	}

	defer tx.Rollback(context.Background())

	insertRow := "WITH inserted as (INSERT INTO orders(quantity, price, currency, order_type, date, asset_id, brokerage_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, quantity, price, currency, order_type, date, brokerage_id) SELECT inserted.id, inserted.quantity, inserted.price, inserted.currency, inserted.order_type, inserted.date, json_build_object('id', b.id, 'name', b.name, 'country', b.country) as brokerage FROM inserted INNER JOIN brokerage as b ON inserted.brokerage_id = b.id;"

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
