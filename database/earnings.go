package database

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
)

func CreateEarningRow(dbpool PgxIface, earningOrder EarningsBodyPost,
	assetId string, userUid string) []EarningsApiReturn {

	var earningRow []EarningsApiReturn

	insertRow := `
	insert into
		earnings("type", earning, date, currency, asset_id, user_uid)
	values ($1, $2, $3, $4, $5, $6)
	returning id, "type", earning, "date", currency, asset_id;
	`
	err := pgxscan.Select(context.Background(), dbpool, &earningRow, insertRow,
		earningOrder.EarningType, earningOrder.Amount, earningOrder.Date,
		earningOrder.Currency, assetId, userUid)
	if err != nil {
		fmt.Println(err)
	}

	return earningRow
}

func SearchEarningFromAssetUser(dbpool PgxIface, assetId string, userUid string) (
	[]EarningsApiReturn, error) {

	var earningsReturn []EarningsApiReturn

	query := `
	SELECT
		id, type, earning, date, currency
	FROM earnings
	WHERE asset_id = $1 and user_uid = $2;
	`

	err := pgxscan.Select(context.Background(), dbpool, &earningsReturn, query,
		assetId, userUid)
	if err != nil {
		fmt.Println("database.SearchEarningFromAssetUser: ", err)
	}

	return earningsReturn, err
}

func DeleteEarningsFromAssetUser(dbpool PgxIface, assetId string, userUid string) (
	[]EarningsApiReturn, error) {
	var earningsId []EarningsApiReturn

	queryDeleteEarnings := `
	delete from earnings as e
	where e.asset_id = $1 and e.user_uid = $2
	returning e.id;
	`
	err := pgxscan.Select(context.Background(), dbpool, &earningsId,
		queryDeleteEarnings, assetId, userUid)
	if err != nil {
		fmt.Println("database.DeleteOrders: ", err)
	}

	return earningsId, err

}

func DeleteEarningFromUser(dbpool PgxIface, id string,
	userUid string) string {
	var orderId string

	query := `
	delete from earnings as e
	where e.id = $1 and e.user_uid = $2
	returning e.id
	`
	row := dbpool.QueryRow(context.Background(), query, id, userUid)
	err := row.Scan(&orderId)
	if err != nil {
		fmt.Println("database.DeleteOrder: ", err)
	}

	return orderId
}

func UpdateEarningsFromUser(dbpool PgxIface, earningsUpdate EarningsBodyPost,
	userUid string) []EarningsApiReturn {
	var earningsInfo []EarningsApiReturn

	query := `
	update earnings as e
	set type = $3,
		earning = $4,
		"date" = $5
	where e.id = $1 and e.user_uid = $2
	returning e.id, e.earning, e."date", e.type;
	`
	err := pgxscan.Select(context.Background(), dbpool, &earningsInfo,
		query, earningsUpdate.Id, userUid, earningsUpdate.EarningType,
		earningsUpdate.Amount, earningsUpdate.Date)
	if err != nil {
		fmt.Println("database.UpdateEarningsFromUser: ", err)
	}

	return earningsInfo
}
