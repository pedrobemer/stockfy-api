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
