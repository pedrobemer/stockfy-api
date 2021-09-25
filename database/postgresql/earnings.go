package postgresql

import (
	"context"
	"fmt"
	"stockfyApi/database"

	"github.com/georgysavva/scany/pgxscan"
)

func (r *repo) CreateEarningRow(earningOrder database.Earnings) []database.Earnings {

	var earningRow []database.Earnings

	insertRow := `
	WITH inserted as (
	INSERT INTO
		earnings("type", earning, date, currency, asset_id, user_uid)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, "type", earning, "date", currency, asset_id
	)
	SELECT
		inserted.id, inserted.type, inserted.earning, inserted.date,
		inserted.currency,
		jsonb_build_object(
			'id', ast.id,
			'symbol', ast.symbol
		) asset
	FROM inserted
	INNER JOIN asset as ast
	ON ast.id = inserted.asset_id;
	`

	err := pgxscan.Select(context.Background(), r.dbpool, &earningRow, insertRow,
		earningOrder.Type, earningOrder.Earning, earningOrder.Date,
		earningOrder.Currency, earningOrder.Asset.Id, earningOrder.UserUid)
	if err != nil {
		fmt.Println(err)
	}

	return earningRow
}

func (r *repo) SearchEarningFromAssetUser(assetId string, userUid string) (
	[]database.Earnings, error) {

	var earningsReturn []database.Earnings

	query := `
	SELECT
		eng.id, type, earning, date, currency,
		jsonb_build_object(
			'id', ast.id,
			'symbol', ast.symbol
		) as asset
	FROM earnings as eng
	INNER JOIN asset as ast
	ON ast.id = eng.asset_id
	WHERE asset_id = $1 and user_uid = $2;
	`

	err := pgxscan.Select(context.Background(), r.dbpool, &earningsReturn, query,
		assetId, userUid)
	if err != nil {
		fmt.Println("database.SearchEarningFromAssetUser: ", err)
	}

	return earningsReturn, err
}

func (r *repo) DeleteEarningsFromAssetUser(assetId string, userUid string) (
	[]database.Earnings, error) {
	var earningsId []database.Earnings

	queryDeleteEarnings := `
	WITH deleted as (
	DELETE FROM earnings
	WHERE asset_id = $1 and user_uid = $2
	RETURNIN id, asset_id
	)
	SELECT
		deleted.id,
		jsonb_build_object(
			'id', ast.id,
			'symbol', ast.symbol
		) as asset
	FROM deleted
	INNER JOIN asset as ast
	ON ast.id = deleted.asset_id;
	`
	err := pgxscan.Select(context.Background(), r.dbpool, &earningsId,
		queryDeleteEarnings, assetId, userUid)
	if err != nil {
		fmt.Println("database.DeleteOrders: ", err)
	}

	return earningsId, err

}

func (r *repo) DeleteEarningFromUser(id string, userUid string) string {
	var orderId string

	query := `
	delete from earnings as e
	where e.id = $1 and e.user_uid = $2
	returning e.id
	`
	row := r.dbpool.QueryRow(context.Background(), query, id, userUid)
	err := row.Scan(&orderId)
	if err != nil {
		fmt.Println("database.DeleteOrder: ", err)
	}

	return orderId
}

func (r *repo) UpdateEarningsFromUser(earningsUpdate database.Earnings) []database.Earnings {
	var earningsInfo []database.Earnings

	query := `
	update earnings as e
	set type = $3,
		earning = $4,
		"date" = $5
	where e.id = $1 and e.user_uid = $2
	returning e.id, e.earning, e."date", e.type;
	`
	err := pgxscan.Select(context.Background(), r.dbpool, &earningsInfo,
		query, earningsUpdate.Id, earningsUpdate.UserUid, earningsUpdate.Type,
		earningsUpdate.Earning, earningsUpdate.Date)
	if err != nil {
		fmt.Println("database.UpdateEarningsFromUser: ", err)
	}

	return earningsInfo
}
