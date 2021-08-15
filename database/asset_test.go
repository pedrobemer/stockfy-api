package database

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

func TestSingleSearchAsset(t *testing.T) {

	assetType := AssetTypeApiReturn{
		Id:      "28ccf27a-ed8b-11eb-9a03-0242ac130003",
		Type:    "STOCK",
		Name:    "Ações Brasil",
		Country: "BR",
	}
	preference := "ON"

	sectorInfo := SectorApiReturn{
		Id:   "83ae92f8-ed8b-11eb-9a03-0242ac130003",
		Name: "Finance",
	}

	var expectedAsset = []AssetQueryReturn{
		{
			Id:         "0a52d206-ed8b-11eb-9a03-0242ac130003",
			Symbol:     "ITUB4",
			Preference: &preference,
			Fullname:   "Itau Unibanco Holding SA",
			AssetType:  &assetType,
			Sector:     &sectorInfo,
		},
	}

	query := `
	SELECT
		(.+)
	FROM asset as a
	INNER JOIN assettype as at
	ON a.asset_type_id = at.id
	(.+)
	GROUP BY a.symbol, a.id, preference, fullname, at.type, at.id, at.name,
	at.country;
	`

	queryTest := `
	select
		(.+)
	from sector as s
	inner join asset as a
	on a.sector_id = s.id
	(.+)
	`
	columns := []string{"id", "symbol", "preference", "fullname", "asset_type"}
	columns_test := []string{"id", "name"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs("ITUB4").WillReturnRows(
		rows.AddRow("0a52d206-ed8b-11eb-9a03-0242ac130003", "ITUB4", &preference,
			"Itau Unibanco Holding SA", &assetType))

	rows_test := mock.NewRows(columns_test)
	mock.ExpectQuery(queryTest).WithArgs("ITUB4").WillReturnRows(
		rows_test.AddRow("83ae92f8-ed8b-11eb-9a03-0242ac130003", "Finance"))

	asset, err := SearchAsset(mock, "ITUB4", "")
	if err != nil {
		fmt.Println(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, asset)
	assert.Equal(t, expectedAsset, asset)
}

func TestSingleSearchAssetWithOrders(t *testing.T) {
	tr, err := time.Parse("2021-07-05", "2021-07-21")
	tr2, err := time.Parse("2021-07-05", "2020-04-02")

	assetType := AssetTypeApiReturn{
		Id:      "28ccf27a-ed8b-11eb-9a03-0242ac130003",
		Type:    "STOCK",
		Name:    "Ações Brasil",
		Country: "BR",
	}

	preference := "ON"

	brokerageInfo := BrokerageApiReturn{
		Id:      "55555555-ed8b-11eb-9a03-0242ac130003",
		Name:    "Clear",
		Country: "BR",
	}

	orderList := []OrderApiReturn{
		{
			Id:        "44444444-ed8b-11eb-9a03-0242ac130003",
			Quantity:  20,
			Price:     39.93,
			Currency:  "BRL",
			OrderType: "buy",
			Date:      tr,
			Brokerage: &brokerageInfo,
		},
		{
			Id:        "yeid847e-ed8b-11eb-9a03-0242ac130003",
			Quantity:  5,
			Price:     27.13,
			Currency:  "BRL",
			OrderType: "buy",
			Date:      tr2,
			Brokerage: &brokerageInfo,
		},
	}

	sectorInfo := SectorApiReturn{
		Id:   "83ae92f8-ed8b-11eb-9a03-0242ac130003",
		Name: "Finance",
	}

	var expectedAsset = []AssetQueryReturn{
		{
			Id:         "0a52d206-ed8b-11eb-9a03-0242ac130003",
			Symbol:     "ITUB4",
			Preference: &preference,
			Fullname:   "Itau Unibanco Holding SA",
			AssetType:  &assetType,
			Sector:     &sectorInfo,
			OrdersList: orderList,
		},
	}

	query := regexp.QuoteMeta(`
	SELECT
		a.id, symbol, preference, a.fullname,
		json_build_object(
			'id', at.id,
			'type', at.type,
			'name', at.name,
			'country', at.country
		) as asset_type,
		json_agg(
			json_build_object(
				'id', o.id,
				'quantity', o.quantity,
				'price', o.price,
				'currency', o.currency,
				'ordertype', o.order_type,
				'date', date,
				'brokerage',
				json_build_object(
					'id', b.id,
					'name', b.name,
					'country', b.country
				)
			)
		) as orders_list
		FROM asset as a
		INNER JOIN assettype as at
		ON a.asset_type_id = at.id
		INNER JOIN orders as o
		ON a.id = o.asset_id
		INNER JOIN brokerage as b
		ON o.brokerage_id = b.id
		WHERE a.symbol=$1
		GROUP BY a.symbol, a.id, preference, a.fullname, at.type, at.id,
			at.name, at.country;
	`)

	querySector := `
	select
		(.+)
	from sector as s
	inner join asset as a
	on a.sector_id = s.id
	(.+)
	`

	columnsAsset := []string{"id", "symbol", "preference", "fullname",
		"asset_type", "orders_list"}
	columnsSector := []string{"id", "name"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columnsAsset)
	mock.ExpectQuery(query).WithArgs("ITUB4").WillReturnRows(
		rows.AddRow("0a52d206-ed8b-11eb-9a03-0242ac130003", "ITUB4", &preference,
			"Itau Unibanco Holding SA", &assetType, orderList))

	rows_sector := mock.NewRows(columnsSector)
	mock.ExpectQuery(querySector).WithArgs("ITUB4").WillReturnRows(
		rows_sector.AddRow("83ae92f8-ed8b-11eb-9a03-0242ac130003", "Finance"))

	asset, err := SearchAsset(mock, "ITUB4", "ONLYORDERS")
	if err != nil {
		fmt.Println(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, asset)
	assert.Equal(t, expectedAsset, asset)

}

func TestSingleSearchAssetWithOrderInfo(t *testing.T) {

	assetType := AssetTypeApiReturn{
		Id:      "28ccf27a-ed8b-11eb-9a03-0242ac130003",
		Type:    "STOCK",
		Name:    "Ações Brasil",
		Country: "BR",
	}

	preference := "ON"

	ordersInfo := OrderGeneralInfos{
		TotalQuantity:        25,
		WeightedAdjPrice:     37.37,
		WeightedAveragePrice: 37.37,
	}

	sectorInfo := SectorApiReturn{
		Id:   "83ae92f8-ed8b-11eb-9a03-0242ac130003",
		Name: "Finance",
	}

	var expectedAsset = []AssetQueryReturn{
		{
			Id:         "0a52d206-ed8b-11eb-9a03-0242ac130003",
			Symbol:     "ITUB4",
			Preference: &preference,
			Fullname:   "Itau Unibanco Holding SA",
			AssetType:  &assetType,
			Sector:     &sectorInfo,
			OrderInfo:  &ordersInfo,
		},
	}

	query := regexp.QuoteMeta(`
	SELECT
		a.id, symbol, preference, a.fullname,
		json_build_object(
			'id', at.id,
			'type', at.type,
			'name', at.name,
			'country', at.country
		) as asset_type,
		json_build_object(
			'totalQuantity', sum(o.quantity),
			'weightedAdjPrice', SUM(o.quantity * price)/SUM(o.quantity),
			'weightedAveragePrice', (
				SUM(o.quantity*o.price) FILTER(WHERE o.order_type = 'buy'))/
				(SUM(o.quantity) FILTER(WHERE o.order_type = 'buy')
			)
		) as orders_info
	FROM asset as a
	INNER JOIN assettype as at
	ON a.asset_type_id = at.id
	INNER JOIN orders as o
	ON a.id = o.asset_id
	INNER JOIN brokerage as b
	ON o.brokerage_id = b.id
	WHERE a.symbol=$1
	GROUP BY a.symbol, a.id, preference, a.fullname, at.type, at.id,
		at.name, at.country;
	`)

	querySector := `
	select
		(.+)
	from sector as s
	inner join asset as a
	on a.sector_id = s.id
	(.+)
	`

	columnsAsset := []string{"id", "symbol", "preference", "fullname",
		"asset_type", "orders_info"}
	columnsSector := []string{"id", "name"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columnsAsset)
	mock.ExpectQuery(query).WithArgs("ITUB4").WillReturnRows(
		rows.AddRow("0a52d206-ed8b-11eb-9a03-0242ac130003", "ITUB4", &preference,
			"Itau Unibanco Holding SA", &assetType, &ordersInfo))

	rows_sector := mock.NewRows(columnsSector)
	mock.ExpectQuery(querySector).WithArgs("ITUB4").WillReturnRows(
		rows_sector.AddRow("83ae92f8-ed8b-11eb-9a03-0242ac130003", "Finance"))

	asset, err := SearchAsset(mock, "ITUB4", "ONLYINFO")
	if err != nil {
		fmt.Println(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, asset)
	assert.Equal(t, expectedAsset, asset)

}

func TestSingleSearchAssetAllInfo(t *testing.T) {
	tr, err := time.Parse("2021-07-05", "2021-07-21")
	tr2, err := time.Parse("2021-07-05", "2020-04-02")

	assetType := AssetTypeApiReturn{
		Id:      "28ccf27a-ed8b-11eb-9a03-0242ac130003",
		Type:    "STOCK",
		Name:    "Ações Brasil",
		Country: "BR",
	}

	preference := "ON"

	brokerageInfo := BrokerageApiReturn{
		Id:      "55555555-ed8b-11eb-9a03-0242ac130003",
		Name:    "Clear",
		Country: "BR",
	}

	orderList := []OrderApiReturn{
		{
			Id:        "44444444-ed8b-11eb-9a03-0242ac130003",
			Quantity:  20,
			Price:     39.93,
			Currency:  "BRL",
			OrderType: "buy",
			Date:      tr,
			Brokerage: &brokerageInfo,
		},
		{
			Id:        "yeid847e-ed8b-11eb-9a03-0242ac130003",
			Quantity:  5,
			Price:     27.13,
			Currency:  "BRL",
			OrderType: "buy",
			Date:      tr2,
			Brokerage: &brokerageInfo,
		},
	}

	ordersInfo := OrderGeneralInfos{
		TotalQuantity:        25,
		WeightedAdjPrice:     37.37,
		WeightedAveragePrice: 37.37,
	}

	sectorInfo := SectorApiReturn{
		Id:   "83ae92f8-ed8b-11eb-9a03-0242ac130003",
		Name: "Finance",
	}

	var expectedAsset = []AssetQueryReturn{
		{
			Id:         "0a52d206-ed8b-11eb-9a03-0242ac130003",
			Symbol:     "ITUB4",
			Preference: &preference,
			Fullname:   "Itau Unibanco Holding SA",
			AssetType:  &assetType,
			Sector:     &sectorInfo,
			OrdersList: orderList,
			OrderInfo:  &ordersInfo,
		},
	}

	query := regexp.QuoteMeta(`
	SELECT
		a.id, symbol, preference, a.fullname,
		json_build_object(
			'id', at.id,
			'type', at.type,
			'name', at.name,
			'country', at.country
		) as asset_type,
		json_build_object(
			'totalQuantity', sum(o.quantity),
			'weightedAdjPrice', SUM(o.quantity * price)/SUM(o.quantity),
			'weightedAveragePrice', (
				SUM(o.quantity*o.price) FILTER(WHERE o.order_type = 'buy'))
				/(SUM(o.quantity) FILTER(WHERE o.order_type = 'buy')
			)
		) as orders_info,
		json_agg(
			json_build_object(
				'id', o.id,
				'quantity', o.quantity,
				'price', o.price,
				'currency', o.currency,
				'ordertype', o.order_type,
				'date', date,
				'brokerage',
				json_build_object(
					'id', b.id,
					'name', b.name,
					'country', b.country
				)
			)
		) as orders_list
	FROM asset as a
	INNER JOIN assettype as at
	ON a.asset_type_id = at.id
	INNER JOIN orders as o
	ON a.id = o.asset_id
	INNER JOIN brokerage as b
	ON o.brokerage_id = b.id
	WHERE a.symbol=$1
	GROUP BY a.symbol, a.id, preference, a.fullname, at.type, at.id,
		at.name, at.country;
	`)

	querySector := `
	select
		(.+)
	from sector as s
	inner join asset as a
	on a.sector_id = s.id
	(.+)
	`

	columnsAsset := []string{"id", "symbol", "preference", "fullname",
		"asset_type", "orders_info", "orders_list"}
	columnsSector := []string{"id", "name"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columnsAsset)
	mock.ExpectQuery(query).WithArgs("ITUB4").WillReturnRows(
		rows.AddRow("0a52d206-ed8b-11eb-9a03-0242ac130003", "ITUB4", &preference,
			"Itau Unibanco Holding SA", &assetType, &ordersInfo, orderList))

	rows_sector := mock.NewRows(columnsSector)
	mock.ExpectQuery(querySector).WithArgs("ITUB4").WillReturnRows(
		rows_sector.AddRow("83ae92f8-ed8b-11eb-9a03-0242ac130003", "Finance"))

	asset, err := SearchAsset(mock, "ITUB4", "ALL")
	if err != nil {
		fmt.Println(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, asset)
	assert.Equal(t, expectedAsset, asset)

}

func TestCreateAsset(t *testing.T) {

	asset := AssetInsert{
		AssetType: "STOCK",
		Sector:    "Finance",
		Symbol:    "ITUB4",
		Country:   "BR",
		Fullname:  "Itau Unibanco Holding SA",
	}

	expectedAssetReturn := AssetApiReturn{
		Id:         "000aaaa6-ed8b-11eb-9a03-0242ac130003",
		Preference: "PN",
		Fullname:   "Itau Unibanco Holding SA",
		Symbol:     "ITUB4",
	}

	insertRow := regexp.QuoteMeta(`
	INSERT INTO
		asset(preference, fullname, symbol, asset_type_id, sector_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, preference, fullname, symbol;
	`)

	columns := []string{"id", "preference", "fullname", "symbol"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectBegin()
	mock.ExpectQuery(insertRow).WithArgs("PN", "Itau Unibanco Holding SA",
		"ITUB4", "0a52d206-ed8b-11eb-9a03-0242ac130003",
		"83ae92f8-ed8b-11eb-9a03-0242ac130003").WillReturnRows(
		rows.AddRow("000aaaa6-ed8b-11eb-9a03-0242ac130003", "PN",
			"Itau Unibanco Holding SA", "ITUB4"))
	mock.ExpectCommit()

	assetRtr := CreateAsset(mock, asset, "0a52d206-ed8b-11eb-9a03-0242ac130003",
		"83ae92f8-ed8b-11eb-9a03-0242ac130003")

	assert.NotNil(t, assetRtr)
	assert.Equal(t, expectedAssetReturn, assetRtr)
}
