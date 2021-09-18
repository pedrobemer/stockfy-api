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
				brokerage_id, user_uid
			)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
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
	mock.ExpectQuery(insertRow).WithArgs(10.0, 20.29, "USD", "buy",
		"0001-01-01 00:00:00 +0000 UTC", "1111BBBB-ed8b-11eb-9a03-0242ac130003",
		"55555555-ed8b-11eb-9a03-0242ac130003", "aa48fafh4").
		WillReturnRows(rows.AddRow("a8a8a8a8-ed8b-11eb-9a03-0242ac130003", 10.0,
			20.29, "USD", "buy", tr, &brokerageInfo))
	mock.ExpectCommit()

	orderReturn := CreateOrder(mock, orderInsert, "1111BBBB-ed8b-11eb-9a03-0242ac130003",
		"55555555-ed8b-11eb-9a03-0242ac130003", "aa48fafh4")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, orderReturn)
	assert.Equal(t, expectedOrderReturn, orderReturn)
}

func TestSearchOrdersFromAssetUser(t *testing.T) {
	tr, err := time.Parse("2021-07-05", "2020-04-02")
	userUid := "aji392a"

	brokerage := BrokerageApiReturn{
		Id:      "55555555-ed8b-11eb-9a03-0242ac130003",
		Name:    "Test Brokerage",
		Country: "US",
	}

	expectedOrderReturn := []OrderApiReturn{
		{
			Id:        "a8a8a8a8-ed8b-11eb-9a03-0242ac130003",
			Quantity:  20,
			Price:     29.29,
			Currency:  "USD",
			OrderType: "buy",
			Date:      tr,
			Brokerage: &brokerage,
		},
		{
			Id:        "a9a999a9-ed8b-11eb-9a03-0242ac130003",
			Quantity:  198,
			Price:     20.00,
			Currency:  "USD",
			OrderType: "buy",
			Date:      tr,
			Brokerage: &brokerage,
		},
	}

	query := regexp.QuoteMeta(`
	SELECT
		o.id, quantity, price, currency, order_type, date,
		json_build_object(
			'id', b.id,
			'name', b."name",
			'country', b.country
		) as brokerage
	FROM orders as o
	INNER JOIN brokerage as b
	ON b.id = o.brokerage_id
	WHERE asset_id = $1 and user_uid = $2;
	`)

	columns := []string{"id", "quantity", "price", "currency", "order_type",
		"date", "brokerage"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs("aak49", userUid).WillReturnRows(
		rows.AddRow(expectedOrderReturn[0].Id, expectedOrderReturn[0].Quantity,
			expectedOrderReturn[0].Price, expectedOrderReturn[0].Currency,
			expectedOrderReturn[0].OrderType, expectedOrderReturn[0].Date,
			expectedOrderReturn[0].Brokerage).AddRow(expectedOrderReturn[1].Id,
			expectedOrderReturn[1].Quantity, expectedOrderReturn[1].Price,
			expectedOrderReturn[1].Currency, expectedOrderReturn[1].OrderType,
			expectedOrderReturn[1].Date, expectedOrderReturn[1].Brokerage))

	ordersReturn, err := SearchOrdersFromAssetUser(mock, "aak49", userUid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, ordersReturn)
	assert.Equal(t, expectedOrderReturn, ordersReturn)

}

func TestDeleteSingleOrderFromUser(t *testing.T) {

	expectedOrderId := "a8a8a8a8-ed8b-11eb-9a03-0242ac130003"
	userUid := "aji392a"

	query := regexp.QuoteMeta(`
	delete from orders as o
	where o.id = $1 and o.user_uid = $2
	returning o.id
	`)

	columns := []string{"id"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs("a8a8a8a8-ed8b-11eb-9a03-0242ac130003",
		userUid).WillReturnRows(rows.AddRow("a8a8a8a8-ed8b-11eb-9a03-0242ac130003"))

	orderId := DeleteOrderFromUser(mock, "a8a8a8a8-ed8b-11eb-9a03-0242ac130003",
		userUid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, orderId)
	assert.Equal(t, expectedOrderId, orderId)
}

func TestDeleteOrdersFromAsset(t *testing.T) {

	expectedOrderIds := []OrderApiReturn{
		{
			Id: "a8a8a8a8-ed8b-11eb-9a03-0242ac130003",
		},
		{
			Id: "b7a8a8a8-ed8b-11eb-9a03-0242ac130003",
		},
	}

	assetId := "3e3e3e3w-ed8b-11eb-9a03-0242ac130003"

	queryDeleteOrders := regexp.QuoteMeta(`
	delete from orders as o
	where o.asset_id = $1
	returning o.id;
	`)

	columns := []string{"id"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(queryDeleteOrders).WithArgs(assetId).
		WillReturnRows(rows.AddRow("a8a8a8a8-ed8b-11eb-9a03-0242ac130003").
			AddRow("b7a8a8a8-ed8b-11eb-9a03-0242ac130003"))

	orderIds := DeleteOrdersFromAsset(mock, assetId)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, orderIds)
	assert.Equal(t, expectedOrderIds, orderIds)
}

func TestDeleteOrdersFromAssetUser(t *testing.T) {

	expectedOrderIds := []OrderApiReturn{
		{
			Id: "a8a8a8a8-ed8b-11eb-9a03-0242ac130003",
		},
		{
			Id: "b7a8a8a8-ed8b-11eb-9a03-0242ac130003",
		},
	}

	userUid := "aji392a"
	assetId := "3e3e3e3w-ed8b-11eb-9a03-0242ac130003"

	queryDeleteOrders := regexp.QuoteMeta(`
	delete from orders as o
	where o.asset_id = $1 and o.user_uid = $2
	returning o.id;
	`)

	columns := []string{"id"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(queryDeleteOrders).WithArgs(assetId, userUid).
		WillReturnRows(rows.AddRow("a8a8a8a8-ed8b-11eb-9a03-0242ac130003").
			AddRow("b7a8a8a8-ed8b-11eb-9a03-0242ac130003"))

	orderIds, err := DeleteOrdersFromAssetUser(mock, assetId, userUid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, orderIds)
	assert.Equal(t, expectedOrderIds, orderIds)
}

func TestUpdateSingleOrderFromUser(t *testing.T) {
	tr, err := time.Parse("2021-07-05", "2020-04-02")

	orderInsert := OrderBodyPost{
		Id:        "3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
		Symbol:    "VTI",
		Fullname:  "Vanguard Total Stock Market US",
		Brokerage: "Avenue",
		Quantity:  20.0,
		Price:     20.29,
		Currency:  "USD",
		OrderType: "buy",
		Date:      "0001-01-01 00:00:00 +0000 UTC",
		Country:   "US",
		AssetType: "ETF",
	}

	expectedUpdatedOrder := []OrderApiReturn{
		{
			Id:        "3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
			Quantity:  20.0,
			Price:     20.29,
			Date:      tr,
			OrderType: "buy",
		},
	}

	userUid := "aji392a"

	query := regexp.QuoteMeta(`
	update orders as o
	set quantity = $3,
		price = $4,
		order_type = $5,
		"date" = $6
	where o.id = $1 and o.user_uid = $2
	returning o.id, o.quantity, o.price, o."date", o.order_type;
	`)

	columns := []string{"id", "quantity", "price", "date", "order_type"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs("3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
		userUid, 20.0, 20.29, "buy", "0001-01-01 00:00:00 +0000 UTC").WillReturnRows(
		rows.AddRow("3e3e3e3w-ed8b-11eb-9a03-0242ac130003", 20.0, 20.29,
			tr, "buy"))

	updatedOrder := UpdateOrderFromUser(mock, orderInsert, userUid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, updatedOrder)
	assert.Equal(t, expectedUpdatedOrder, updatedOrder)
}
