package database

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateEarningRow(t *testing.T) {
	tr, err := time.Parse("2021-07-05", "2020-04-02")

	userUid := "eji90vl5"

	earningOrder := EarningsBodyPost{
		EarningType: "Dividendos",
		Amount:      5.59,
		Currency:    "BRL",
		Date:        "0001-01-01 00:00:00 +0000 UTC",
	}

	expectedEarningRow := []EarningsApiReturn{
		{
			Id:       "akxn-1234",
			Type:     "Dividendos",
			Earning:  5.59,
			Date:     tr,
			Currency: "BRL",
			AssetId:  "a69a3",
		},
	}

	insertRow := regexp.QuoteMeta(`
	insert into
		earnings("type", earning, date, currency, asset_id, user_uid)
	values ($1, $2, $3, $4, $5, $6)
	returning id, "type", earning, "date", currency, asset_id;
	`)

	columns := []string{"id", "type", "earning", "date", "currency", "asset_id"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(insertRow).WithArgs("Dividendos", 5.59,
		"0001-01-01 00:00:00 +0000 UTC", "BRL", "a69a3", userUid).WillReturnRows(
		rows.AddRow("akxn-1234", "Dividendos", 5.59,
			tr, "BRL", "a69a3"))

	earningRow := CreateEarningRow(mock, earningOrder, "a69a3", userUid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, earningRow)
	assert.Equal(t, expectedEarningRow, earningRow)
}

func TestSearchEarningFromAssetUser(t *testing.T) {

	tr, err := time.Parse("2021-07-05", "2020-04-02")
	userUid := "eji90vl5"
	assetId := "ajfj49a"

	expectedEarningsReturn := []EarningsApiReturn{
		{
			Id:       "3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
			Earning:  5.29,
			Type:     "Dividendos",
			Date:     tr,
			Currency: "BRL",
		},
		{
			Id:       "4e4e4e4w-ed8b-11eb-9a03-0242ac130003",
			Earning:  10.48,
			Type:     "JCP",
			Date:     tr,
			Currency: "BRL",
		},
	}

	query := regexp.QuoteMeta(`
	SELECT
		id, type, earning, date, currency
	FROM earnings
	WHERE asset_id = $1 and user_uid = $2;
	`)

	columns := []string{"id", "type", "earning", "date", "currency"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs(assetId, userUid).
		WillReturnRows(rows.AddRow("3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
			"Dividendos", 5.29, tr, "BRL").AddRow(
			"4e4e4e4w-ed8b-11eb-9a03-0242ac130003", "JCP", 10.48, tr, "BRL"))

	earningsReturn, _ := SearchEarningFromAssetUser(mock, assetId, userUid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, earningsReturn)
	assert.Equal(t, expectedEarningsReturn, earningsReturn)

}

func TestDeleteEarningFromUser(t *testing.T) {

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
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs("3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
		userUid).WillReturnRows(rows.AddRow("3e3e3e3w-ed8b-11eb-9a03-0242ac130003"))

	orderId := DeleteEarningFromUser(mock, expectedEarningId, userUid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, orderId)
	assert.Equal(t, expectedEarningId, orderId)

}

func TestDeleteEarningsFromAssetUser(t *testing.T) {

	expectedOrderIds := []EarningsApiReturn{
		{
			Id: "a8a8a8a8-ed8b-11eb-9a03-0242ac130003",
		},
		{
			Id: "b7a8a8a8-ed8b-11eb-9a03-0242ac130003",
		},
	}

	userUid := "aji392a"
	assetId := "3e3e3e3w-ed8b-11eb-9a03-0242ac130003"

	queryDeleteEarnings := regexp.QuoteMeta(`
	delete from earnings as e
	where e.asset_id = $1 and e.user_uid = $2
	returning e.id;
	`)

	columns := []string{"id"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(queryDeleteEarnings).WithArgs(assetId, userUid).
		WillReturnRows(rows.AddRow("a8a8a8a8-ed8b-11eb-9a03-0242ac130003").
			AddRow("b7a8a8a8-ed8b-11eb-9a03-0242ac130003"))

	orderIds, err := DeleteEarningsFromAssetUser(mock, assetId, userUid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, orderIds)
	assert.Equal(t, expectedOrderIds, orderIds)
}

func TestUpdateEarningsFromUser(t *testing.T) {
	tr, err := time.Parse("2021-07-05", "2020-04-02")

	userUid := "eji90vl5"

	earningsUpdate := EarningsBodyPost{
		Id:          "3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
		EarningType: "Dividendos",
		Amount:      5.29,
		Date:        "0001-01-01 00:00:00 +0000 UTC",
	}

	expectedEarningsReturn := []EarningsApiReturn{
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
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs("3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
		userUid, "Dividendos", 5.29, "0001-01-01 00:00:00 +0000 UTC").
		WillReturnRows(rows.AddRow("3e3e3e3w-ed8b-11eb-9a03-0242ac130003", 5.29,
			tr, "Dividendos"))

	updatedOrder := UpdateEarningsFromUser(mock, earningsUpdate, userUid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, updatedOrder)
	assert.Equal(t, expectedEarningsReturn, updatedOrder)
}
