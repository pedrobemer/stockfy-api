package postgresql

import (
	"context"
	"fmt"
	"stockfyApi/entity"

	"github.com/georgysavva/scany/pgxscan"
)

type EarningPostgres struct {
	dbpool PgxIface
}

func NewEarningPostgres(db *PgxIface) *EarningPostgres {
	return &EarningPostgres{
		dbpool: *db,
	}
}

func (r *EarningPostgres) Create(earningOrder entity.Earnings) []entity.Earnings {

	var earningRow []entity.Earnings

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

func (r *EarningPostgres) SearchFromAssetUser(assetId string, userUid string) (
	[]entity.Earnings, error) {

	var earningsReturn []entity.Earnings

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
		fmt.Println("entity.SearchEarningFromAssetUser: ", err)
	}

	return earningsReturn, err
}

func (r *EarningPostgres) DeleteFromAssetUser(assetId string, userUid string) (
	[]entity.Earnings, error) {
	var earningsId []entity.Earnings

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
		fmt.Println("entity.DeleteOrders: ", err)
	}

	return earningsId, err

}

func (r *EarningPostgres) DeleteFromUser(id string, userUid string) string {
	var orderId string

	query := `
	delete from earnings as e
	where e.id = $1 and e.user_uid = $2
	returning e.id
	`
	row := r.dbpool.QueryRow(context.Background(), query, id, userUid)
	err := row.Scan(&orderId)
	if err != nil {
		fmt.Println("entity.DeleteOrder: ", err)
	}

	return orderId
}

func (r *EarningPostgres) UpdateFromUser(earningsUpdate entity.Earnings) []entity.Earnings {
	var earningsInfo []entity.Earnings

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
		fmt.Println("entity.UpdateEarningsFromUser: ", err)
	}

	return earningsInfo
}
