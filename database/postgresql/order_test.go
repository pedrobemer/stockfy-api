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

func TestOrderCreate(t *testing.T) {
	tr, err := time.Parse("2021-07-05", "2020-04-02")
	userUid := "aa48fafh4"

	brokerageInfo := entity.Brokerage{
		Id:      "55555555-ed8b-11eb-9a03-0242ac130003",
		Name:    "Avenue",
		Country: "US",
	}

	assetInfo := entity.Asset{
		Id:       "1111BBBB-ed8b-11eb-9a03-0242ac130003",
		Symbol:   "VTI",
		Fullname: "Vanguard Total Stock Market US",
	}

	orderInsert := entity.Order{
		// Symbol:    "VTI",
		// Fullname:  "Vanguard Total Stock Market US",
		Id:        "3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
		Asset:     &assetInfo,
		Brokerage: &brokerageInfo,
		Quantity:  10.0,
		Price:     20.29,
		Currency:  "USD",
		OrderType: "buy",
		Date:      tr,
		UserUid:   userUid,
	}

	expectedOrderReturn := entity.Order{
		Id:        "a8a8a8a8-ed8b-11eb-9a03-0242ac130003",
		Quantity:  10.0,
		Price:     20.29,
		Currency:  "USD",
		OrderType: "buy",
		Date:      tr,
		Brokerage: &brokerageInfo,
		Asset:     &assetInfo,
	}

	insertRow := regexp.QuoteMeta(`
	WITH inserted as (
		INSERT INTO
			orders(quantity, price, currency, order_type, date, asset_id,
				brokerage_id, user_uid
			)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, quantity, price, currency, order_type, date, asset_id,
			brokerage_id
	)
	SELECT
		inserted.id, inserted.quantity, inserted.price, inserted.currency,
		inserted.order_type, inserted.date,
		json_build_object(
			'id', b.id,
			'name', b.name,
			'country', b.country
		) as brokerage,
		json_build_object(
			'id', a.id,
			'symbol', a.symbol,
			'preference', a.preference,
			'fullname', a.fullname
		) as asset
	FROM inserted
	INNER JOIN brokerages as b
	ON inserted.brokerage_id = b.id
	INNER JOIN assets as a
	ON inserted.asset_id = a.id;
	`)

	columns := []string{"id", "quantity", "price", "currency", "order_type",
		"date", "brokerage", "asset"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	mock.ExpectBegin()
	// mock.ExpectRollback()
	rows := mock.NewRows(columns)
	mock.ExpectQuery(insertRow).WithArgs(10.0, 20.29, "USD", "buy",
		tr, "1111BBBB-ed8b-11eb-9a03-0242ac130003",
		"55555555-ed8b-11eb-9a03-0242ac130003", userUid).
		WillReturnRows(rows.AddRow("a8a8a8a8-ed8b-11eb-9a03-0242ac130003", 10.0,
			20.29, "USD", "buy", tr, &brokerageInfo, &assetInfo))
	mock.ExpectCommit()

	Orders := OrderPostgres{dbpool: mock}
	orderReturn := Orders.Create(orderInsert)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, orderReturn)
	assert.Equal(t, expectedOrderReturn, orderReturn)
}

func TestOrderSearchFromAssetUser(t *testing.T) {
	tr, err := time.Parse("2021-07-05", "2020-04-02")
	userUid := "aji392a"

	brokerage := entity.Brokerage{
		Id:      "55555555-ed8b-11eb-9a03-0242ac130003",
		Name:    "Test Brokerage",
		Country: "US",
	}

	expectedOrderReturn := []entity.Order{
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
	INNER JOIN brokerages as b
	ON b.id = o.brokerage_id
	WHERE asset_id = $1 and user_uid = $2;
	`)

	columns := []string{"id", "quantity", "price", "currency", "order_type",
		"date", "brokerage"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
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

	Orders := OrderPostgres{dbpool: mock}
	ordersReturn, err := Orders.SearchFromAssetUser("aak49", userUid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, ordersReturn)
	assert.Equal(t, expectedOrderReturn, ordersReturn)

}

func TestOrderSingleDeleteFromUser(t *testing.T) {

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
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs("a8a8a8a8-ed8b-11eb-9a03-0242ac130003",
		userUid).WillReturnRows(rows.AddRow("a8a8a8a8-ed8b-11eb-9a03-0242ac130003"))

	Orders := OrderPostgres{dbpool: mock}
	orderId, _ := Orders.DeleteFromUser("a8a8a8a8-ed8b-11eb-9a03-0242ac130003",
		userUid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, orderId)
	assert.Equal(t, expectedOrderId, orderId)
}

func TestOrderDeleteFromAsset(t *testing.T) {

	expectedOrderIds := []entity.Order{
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
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(queryDeleteOrders).WithArgs(assetId).
		WillReturnRows(rows.AddRow("a8a8a8a8-ed8b-11eb-9a03-0242ac130003").
			AddRow("b7a8a8a8-ed8b-11eb-9a03-0242ac130003"))

	Orders := OrderPostgres{dbpool: mock}
	orderIds, _ := Orders.DeleteFromAsset(assetId)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, orderIds)
	assert.Equal(t, expectedOrderIds, orderIds)
}

func TestOrderDeleteFromAssetUser(t *testing.T) {

	assetInfo := entity.Asset{
		Id:     "1111BBBB-ed8b-11eb-9a03-0242ac130003",
		Symbol: "VTI",
	}

	expectedOrderIds := []entity.Order{
		{
			Id:    "a8a8a8a8-ed8b-11eb-9a03-0242ac130003",
			Asset: &assetInfo,
		},
		{
			Id:    "b7a8a8a8-ed8b-11eb-9a03-0242ac130003",
			Asset: &assetInfo,
		},
	}

	userUid := "aji392a"
	assetId := "3e3e3e3w-ed8b-11eb-9a03-0242ac130003"

	queryDeleteOrders := regexp.QuoteMeta(`
	with deleted as (
	delete from orders as o
	where o.asset_id = $1 and o.user_uid = $2
	returning o.id, o.asset_id
	)
	select
		deleted.id,
		jsonb_build_object(
			'id', ast.id,
			'symbol', ast.symbol
		) as asset
	from deleted
	inner join assets as ast
	on ast.id = deleted.asset_id;
	`)

	columns := []string{"id", "asset"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(queryDeleteOrders).WithArgs(assetId, userUid).
		WillReturnRows(rows.AddRow("a8a8a8a8-ed8b-11eb-9a03-0242ac130003",
			&assetInfo).AddRow("b7a8a8a8-ed8b-11eb-9a03-0242ac130003", &assetInfo))

	Orders := OrderPostgres{dbpool: mock}
	orderIds, err := Orders.DeleteFromAssetUser(assetId, userUid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, orderIds)
	assert.Equal(t, expectedOrderIds, orderIds)
}

func TestOrderSingleUpdateFromUser(t *testing.T) {
	tr, err := time.Parse("2021-07-05", "2020-04-02")

	userUid := "aji392a"

	brokerageInfo := entity.Brokerage{
		Id:      "55555555-ed8b-11eb-9a03-0242ac130003",
		Name:    "Avenue",
		Country: "US",
	}

	orderInsert := entity.Order{
		Id:        "3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
		Brokerage: &brokerageInfo,
		Quantity:  20.0,
		Price:     20.29,
		Currency:  "USD",
		OrderType: "buy",
		Date:      tr,
		UserUid:   userUid,
	}

	expectedUpdatedOrder := []entity.Order{
		{
			Id:        "3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
			Quantity:  20.0,
			Price:     20.29,
			Date:      tr,
			OrderType: "buy",
			Currency:  "USD",
			Brokerage: &brokerageInfo,
		},
	}

	query := regexp.QuoteMeta(`
	with updated as (
	update orders as o
	set quantity = $3,
		price = $4,
		order_type = $5,
		"date" = $6,
		brokerage_id = $7
	where o.id = $1 and o.user_uid = $2
	returning o.id, o.quantity, o.price, o."date", o.order_type, brokerage_id,
		o.currency
	)
	select
		updated.id, updated.quantity, updated.price, updated.order_type,
		updated."date", updated.currency,
		json_build_object(
			'id', updated.brokerage_id,
			'name', b."name",
			'country', b.country
		) as brokerage
	from updated
	inner join brokerages as b
	on b.id = updated.brokerage_id;
	`)

	columns := []string{"id", "quantity", "price", "date", "order_type",
		"currency", "brokerage"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs("3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
		userUid, 20.0, 20.29, "buy", tr, "55555555-ed8b-11eb-9a03-0242ac130003").
		WillReturnRows(rows.AddRow("3e3e3e3w-ed8b-11eb-9a03-0242ac130003", 20.0,
			20.29, tr, "buy", "USD", &brokerageInfo))

	Orders := OrderPostgres{dbpool: mock}
	updatedOrder := Orders.UpdateFromUser(orderInsert)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, updatedOrder)
	assert.Equal(t, expectedUpdatedOrder, updatedOrder)
}

func TestSearchByOrderAndUserId(t *testing.T) {
	tr, err := time.Parse("2021-07-05", "2020-04-02")

	userUid := "aji392a"

	brokerageInfo := entity.Brokerage{
		Id:      "55555555-ed8b-11eb-9a03-0242ac130003",
		Name:    "Avenue",
		Country: "US",
	}

	assetInfo := entity.Asset{Id: "39823-3DNC894"}
	expectedOrderInfo := []entity.Order{
		{
			Id:        "3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
			Quantity:  20.0,
			Price:     20.29,
			Date:      tr,
			OrderType: "buy",
			Currency:  "USD",
			Brokerage: &brokerageInfo,
			Asset:     &assetInfo,
		},
	}

	query := regexp.QuoteMeta(`
	SELECT
		o.id, quantity, price, currency, order_type, date,
		json_build_object(
			'id', b.id,
			'name', b."name",
			'country', b.country
		) as brokerage,
		json_build_object(
			'id', asset_id
		) as asset
	FROM orders as o
	INNER JOIN brokerages as b
	ON b.id = o.brokerage_id
	WHERE o.id = $1 and user_uid = $2;
	`)

	columns := []string{"id", "quantity", "price", "date", "order_type",
		"currency", "brokerage", "asset"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs("3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
		userUid).WillReturnRows(rows.AddRow("3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
		20.0, 20.29, tr, "buy", "USD", &brokerageInfo, &assetInfo))

	Orders := OrderPostgres{dbpool: mock}
	orderInfo, _ := Orders.SearchByOrderAndUserId("3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
		userUid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, orderInfo)
	assert.Equal(t, expectedOrderInfo, orderInfo)
}
