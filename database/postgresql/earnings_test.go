package postgresql

import (
	"context"
	"regexp"
	"stockfyApi/entity"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

func TestEarningCreate(t *testing.T) {
	tr, err := time.Parse("2021-07-05", "2020-04-02")

	userUid := "eji90vl5"

	asset := entity.Asset{
		Id:     "a69a3",
		Symbol: "ITUB4",
	}

	earningOrder := entity.Earnings{
		Type:     "Dividendos",
		Earning:  5.59,
		Currency: "BRL",
		Date:     tr,
		Asset:    &asset,
		UserUid:  userUid,
	}

	expectedEarningRow := []entity.Earnings{
		{
			Id:       "akxn-1234",
			Type:     "Dividendos",
			Earning:  5.59,
			Date:     tr,
			Currency: "BRL",
			Asset:    &asset,
		},
	}

	insertRow := regexp.QuoteMeta(`
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
	`)

	columns := []string{"id", "type", "earning", "date", "currency", "asset"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(insertRow).WithArgs("Dividendos", 5.59, tr, "BRL", "a69a3",
		userUid).WillReturnRows(rows.AddRow("akxn-1234", "Dividendos", 5.59,
		tr, "BRL", &asset))

	Earnings := EarningPostgres{dbpool: mock}
	earningRow, _ := Earnings.Create(earningOrder)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, earningRow)
	assert.Equal(t, expectedEarningRow, earningRow)
}

func TestEarningSearchFromAssetUser(t *testing.T) {

	tr, err := time.Parse("2021-07-05", "2020-04-02")
	userUid := "eji90vl5"

	assetId := "ajfj49a"

	asset := entity.Asset{
		Id:     assetId,
		Symbol: "ITUB4",
	}

	expectedEarningsReturn := []entity.Earnings{
		{
			Id:       "3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
			Earning:  5.29,
			Type:     "Dividendos",
			Date:     tr,
			Currency: "BRL",
			Asset:    &asset,
		},
		{
			Id:       "4e4e4e4w-ed8b-11eb-9a03-0242ac130003",
			Earning:  10.48,
			Type:     "JCP",
			Date:     tr,
			Currency: "BRL",
			Asset:    &asset,
		},
	}

	query := regexp.QuoteMeta(`
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
	`)

	columns := []string{"id", "type", "earning", "date", "currency", "asset"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs(assetId, userUid).
		WillReturnRows(rows.AddRow("3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
			"Dividendos", 5.29, tr, "BRL", &asset).AddRow(
			"4e4e4e4w-ed8b-11eb-9a03-0242ac130003", "JCP", 10.48, tr, "BRL",
			&asset))

	Earnings := EarningPostgres{dbpool: mock}
	earningsReturn, _ := Earnings.SearchFromAssetUser(assetId, userUid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, earningsReturn)
	assert.Equal(t, expectedEarningsReturn, earningsReturn)

}

func TestEarningDeleteFromUser(t *testing.T) {

	expectedEarningId := "3e3e3e3w-ed8b-11eb-9a03-0242ac130003"
	userUid := "eji90vl5"

	query := regexp.QuoteMeta(`
	delete from earnings as e
	where e.id = $1 and e.user_uid = $2
	returning e.id
	`)

	columns := []string{"id"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs("3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
		userUid).WillReturnRows(rows.AddRow("3e3e3e3w-ed8b-11eb-9a03-0242ac130003"))

	Earnings := EarningPostgres{dbpool: mock}
	orderId := Earnings.DeleteFromUser(expectedEarningId, userUid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, orderId)
	assert.Equal(t, expectedEarningId, orderId)

}

func TestEarningDeleteFromAssetUser(t *testing.T) {

	userUid := "aji392a"
	assetId := "3e3e3e3w-ed8b-11eb-9a03-0242ac130003"

	asset := entity.Asset{
		Id:     assetId,
		Symbol: "ITUB4",
	}

	expectedOrderIds := []entity.Earnings{
		{
			Id:    "a8a8a8a8-ed8b-11eb-9a03-0242ac130003",
			Asset: &asset,
		},
		{
			Id:    "b7a8a8a8-ed8b-11eb-9a03-0242ac130003",
			Asset: &asset,
		},
	}

	queryDeleteEarnings := regexp.QuoteMeta(`
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
	`)

	columns := []string{"id", "asset"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(queryDeleteEarnings).WithArgs(assetId, userUid).
		WillReturnRows(rows.AddRow("a8a8a8a8-ed8b-11eb-9a03-0242ac130003",
			&asset).AddRow("b7a8a8a8-ed8b-11eb-9a03-0242ac130003", &asset))

	Earnings := EarningPostgres{dbpool: mock}
	orderIds, err := Earnings.DeleteFromAssetUser(assetId, userUid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, orderIds)
	assert.Equal(t, expectedOrderIds, orderIds)
}

func TestEarningUpdateFromUser(t *testing.T) {
	tr, err := time.Parse("2021-07-05", "2020-04-02")

	userUid := "eji90vl5"

	earningsUpdate := entity.Earnings{
		Id:      "3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
		Type:    "Dividendos",
		Earning: 5.29,
		Date:    tr,
		UserUid: userUid,
	}

	expectedEarningsReturn := []entity.Earnings{
		{
			Id:      "3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
			Earning: 5.29,
			Type:    "Dividendos",
			Date:    tr,
		},
	}

	query := regexp.QuoteMeta(`
	update earnings as e
	set type = $3,
		earning = $4,
		"date" = $5
	where e.id = $1 and e.user_uid = $2
	returning e.id, e.earning, e."date", e.type;
	`)

	columns := []string{"id", "earning", "date", "type"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs("3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
		userUid, "Dividendos", 5.29, tr).
		WillReturnRows(rows.AddRow("3e3e3e3w-ed8b-11eb-9a03-0242ac130003", 5.29,
			tr, "Dividendos"))

	Earnings := EarningPostgres{dbpool: mock}
	updatedOrder := Earnings.UpdateFromUser(earningsUpdate)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, updatedOrder)
	assert.Equal(t, expectedEarningsReturn, updatedOrder)
}
