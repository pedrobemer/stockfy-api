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
