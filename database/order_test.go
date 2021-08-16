package database

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrder(t *testing.T) {
	tr, err := time.Parse("2021-07-05", "2020-04-02")

	brokerageInfo := BrokerageApiReturn{
		Id:      "55555555-ed8b-11eb-9a03-0242ac130003",
		Name:    "Avenue",
		Country: "US",
	}

	orderInsert := OrderBodyPost{
		Id:        "3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
		Symbol:    "VTI",
		Fullname:  "Vanguard Total Stock Market US",
		Brokerage: "Avenue",
		Quantity:  10.0,
		Price:     20.29,
		Currency:  "USD",
		OrderType: "buy",
		Date:      "0001-01-01 00:00:00 +0000 UTC",
		Country:   "US",
		AssetType: "ETF",
	}

	expectedOrderReturn := OrderApiReturn{
		Id:        "a8a8a8a8-ed8b-11eb-9a03-0242ac130003",
		Quantity:  10.0,
		Price:     20.29,
		Currency:  "USD",
		OrderType: "buy",
		Date:      tr,
		Brokerage: &brokerageInfo,
	}

	insertRow := regexp.QuoteMeta(`
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
	`)

	columns := []string{"id", "quantity", "price", "currency", "order_type",
		"date", "brokerage"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	mock.ExpectBegin()
	// mock.ExpectRollback()
	rows := mock.NewRows(columns)
	mock.ExpectQuery(insertRow).WithArgs(10.0, 20.29, "USD", "buy", "0001-01-01 00:00:00 +0000 UTC",
		"1111BBBB-ed8b-11eb-9a03-0242ac130003",
		"55555555-ed8b-11eb-9a03-0242ac130003").
		WillReturnRows(rows.AddRow("a8a8a8a8-ed8b-11eb-9a03-0242ac130003", 10.0,
			20.29, "USD", "buy", tr, &brokerageInfo))
	mock.ExpectCommit()

	orderReturn := CreateOrder(mock, orderInsert, "1111BBBB-ed8b-11eb-9a03-0242ac130003",
		"55555555-ed8b-11eb-9a03-0242ac130003")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, orderReturn)
	assert.Equal(t, expectedOrderReturn, orderReturn)
}
