package database

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
)

func CreateEarningRow(dbpool pgxIface, earningOrder EarningsBodyPost,
	assetId string) []EarningsApiReturn {

	var earningRow []EarningsApiReturn

	insertRow := `
	insert into
		earnings("type", earning, date, currency, asset_id)
	values ($1, $2, $3, $4, $5)
	returning id, "type", earning, "date", currency, asset_id;
	`
	err := pgxscan.Select(context.Background(), dbpool, &earningRow, insertRow,
		earningOrder.EarningType, earningOrder.Amount, earningOrder.Date,
		earningOrder.Currency, assetId)
	if err != nil {
		fmt.Println(err)
	}

	return earningRow
}
