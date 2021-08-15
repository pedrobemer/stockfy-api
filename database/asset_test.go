package database

import (
	"context"
	"errors"
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

func TestSearchAssetByOrderId(t *testing.T) {

	assetType := AssetTypeApiReturn{
		Id:      "28ccf27a-ed8b-11eb-9a03-0242ac130003",
		Type:    "STOCK",
		Name:    "Ações Brasil",
		Country: "BR",
	}

	preference := "ON"

	var expectedAssetInfo = []AssetQueryReturn{
		{
			Id:         "0a52d206-ed8b-11eb-9a03-0242ac130003",
			Symbol:     "ITUB4",
			Preference: &preference,
			AssetType:  &assetType,
		},
	}

	query := regexp.QuoteMeta(`
	select
		a.id, a.preference , a.symbol,
		json_build_object(
			'id', aty.id,
			'type', aty."type",
			'name', aty."name",
			'country', aty.country
		) as asset_type
	from orders as o
	inner join asset as a
	on a.id = o.asset_id
	inner join assettype as aty
	on aty.id = a.asset_type_id
	where o.id = $1;
	`)

	columns := []string{"id", "preference", "symbol", "asset_type"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs("6669aaaa-ed8b-11eb-9a03-0242ac130003").
		WillReturnRows(rows.AddRow(
			"0a52d206-ed8b-11eb-9a03-0242ac130003", &preference, "ITUB4",
			&assetType))

	assetInfo := SearchAssetByOrderId(mock,
		"6669aaaa-ed8b-11eb-9a03-0242ac130003")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, assetInfo)
	assert.Equal(t, expectedAssetInfo, assetInfo)
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

func TestDeleteAsset(t *testing.T) {

	assetType := AssetTypeApiReturn{
		Id:      "28ccf27a-ed8b-11eb-9a03-0242ac130003",
		Type:    "STOCK",
		Name:    "Ações Brasil",
		Country: "BR",
	}
	preference := "PN"

	orderList := []OrderApiReturn{
		{
			Id: "6669aaaa-ed8b-11eb-9a03-0242ac130003",
		},
	}

	var expectedDelAsset = []AssetQueryReturn{
		{
			Id:         "0a52d206-ed8b-11eb-9a03-0242ac130003",
			Symbol:     "ITUB4",
			Preference: &preference,
			Fullname:   "Itau Unibanco Holding SA",
			OrdersList: orderList,
		},
	}

	queryAsset := `
	SELECT
		(.+)
	FROM asset as a
	INNER JOIN assettype as at
	ON a.asset_type_id = at.id
	(.+)
	GROUP BY a.symbol, a.id, preference, fullname, at.type, at.id, at.name,
	at.country;
	`
	querySector := `
	select
		(.+)
	from sector as s
	inner join asset as a
	on a.sector_id = s.id
	(.+)
	`

	queryDeleteOrders := regexp.QuoteMeta(`
	delete from orders as o
	where o.asset_id = $1
	returning o.id;
	`)

	queryDeleteAsset := regexp.QuoteMeta(`
	delete from asset as a
	where a.id = $1
	returning  a.id, a.symbol, a.preference, a.fullname;
	`)

	columnsAsset := []string{"id", "symbol", "preference", "fullname",
		"asset_type"}
	columnsSector := []string{"id", "name"}
	columnsDelOrder := []string{"id"}
	columnsDelAsset := []string{"id", "symbol", "preference", "fullname"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	rowsAsset := mock.NewRows(columnsAsset)
	rowsSector := mock.NewRows(columnsSector)
	rowsDelOrders := mock.NewRows(columnsDelOrder)
	rowsDelAsset := mock.NewRows(columnsDelAsset)

	mock.ExpectQuery(queryAsset).WithArgs("ITUB4").WillReturnRows(
		rowsAsset.AddRow("0a52d206-ed8b-11eb-9a03-0242ac130003", "ITUB4",
			&preference, "Itau Unibanco Holding SA", &assetType))
	mock.ExpectQuery(querySector).WithArgs("ITUB4").WillReturnRows(
		rowsSector.AddRow("83ae92f8-ed8b-11eb-9a03-0242ac130003", "Finance"))

	mock.ExpectQuery(queryDeleteOrders).WithArgs(
		"0a52d206-ed8b-11eb-9a03-0242ac130003").WillReturnRows(
		rowsDelOrders.AddRow("6669aaaa-ed8b-11eb-9a03-0242ac130003"))

	mock.ExpectQuery(queryDeleteAsset).WithArgs(
		"0a52d206-ed8b-11eb-9a03-0242ac130003").WillReturnRows(
		rowsDelAsset.AddRow("0a52d206-ed8b-11eb-9a03-0242ac130003", "ITUB4",
			&preference, "Itau Unibanco Holding SA"))

	assetInfo := DeleteAsset(mock, "ITUB4")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, assetInfo)
	assert.Equal(t, expectedDelAsset, assetInfo)
}

func TestSearchAssetPerAssetTypeWithoutOrderAndNotETFandFII(t *testing.T) {

	preference := "PN"
	preference2 := "ON"

	sectorInfo := SectorApiReturn{
		Id:   "83ae92f8-ed8b-11eb-9a03-0242ac130003",
		Name: "Finance",
	}

	sectorInfo2 := SectorApiReturn{
		Id:   "83838383-ed8b-11eb-9a03-0242ac130003",
		Name: "Health",
	}

	var assets = []AssetQueryReturn{
		{
			Id:         "0a52d206-ed8b-11eb-9a03-0242ac130003",
			Symbol:     "ITUB4",
			Preference: &preference,
			Fullname:   "Itau Unibanco Holding SA",
			Sector:     &sectorInfo,
		},
		{
			Id:         "11111111-ed8b-11eb-9a03-0242ac130003",
			Symbol:     "FLRY3",
			Preference: &preference2,
			Fullname:   "Fleury SA",
			Sector:     &sectorInfo2,
		},
	}

	expectedAssetTypeInfo := []AssetTypeApiReturn{
		{
			Id:      "00000000-ed8b-11eb-9a03-0242ac130003",
			Type:    "STOCK",
			Country: "BR",
			Name:    "Ações Brasil",
			Assets:  assets,
		},
	}

	query := regexp.QuoteMeta(`
	SELECT
		at.id, at.type, at.country, at.name,
		json_agg(
			json_build_object(
				'id', a.id,
				'symbol', a.symbol,
				'preference', a.preference,
				'fullname', a.fullname, 'sector',
				json_build_object(
					'id', s.id,
					'name', s.name
				)
			)
		) as assets
	FROM assettype as at
	INNER JOIN asset as a
	ON at.id = a.asset_type_id
	INNER JOIN sector as s
	ON a.sector_id = s.id
	WHERE at.type = $1 and at.country = $2
	GROUP BY at.id, at.type, at.country, at.name;
	`)

	assetTypeInfo, err := testSearchAssetPerAssetType("STOCK", "BR", "Ações Brasil", false, assets,
		query)
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.NotNil(t, assetTypeInfo)
	assert.Equal(t, expectedAssetTypeInfo, assetTypeInfo)
}

func TestSearchAssetPerAssetTypeWithoutOrderAndIsETForFII(t *testing.T) {

	var assets = []AssetQueryReturn{
		{
			Id:       "0a52d206-ed8b-11eb-9a03-0242ac130003",
			Symbol:   "VTI",
			Fullname: "Vanguard Total Market US",
		},
		{
			Id:       "11111111-ed8b-11eb-9a03-0242ac130003",
			Symbol:   "IJR",
			Fullname: "iShares S&P 600 Small-Caps",
		},
	}

	expectedAssetTypeInfo := []AssetTypeApiReturn{
		{
			Id:      "00000000-ed8b-11eb-9a03-0242ac130003",
			Type:    "ETF",
			Country: "US",
			Name:    "ETFs EUA",
			Assets:  assets,
		},
	}

	query := regexp.QuoteMeta(`
	SELECT
		at.id, at.type, at.country, at.name,
		json_agg(
			json_build_object(
				'id', a.id,
				'symbol', a.symbol,
				'preference', a.preference,
				'fullname', a.fullname
			)
		) as assets
	FROM assettype as at
	INNER JOIN asset as a
	ON at.id = a.asset_type_id
	WHERE at.type = $1 and at.country = $2
	GROUP BY at.id, at.type, at.country, at.name;
	`)

	assetTypeInfo, errorMock := testSearchAssetPerAssetType("ETF", "US",
		"ETFs EUA", false, assets, query)
	if errorMock != nil {
		t.Fatalf(errorMock.Error())
	}

	assert.NotNil(t, assetTypeInfo)
	assert.Equal(t, expectedAssetTypeInfo, assetTypeInfo)

}

func TestSearchAssetPerAssetTypeWithOrderAndNotETFandFII(t *testing.T) {
	preference := "PN"
	preference2 := "ON"

	ordersInfo := OrderGeneralInfos{
		TotalQuantity:        25,
		WeightedAdjPrice:     37.37,
		WeightedAveragePrice: 37.37,
	}

	sectorInfo := SectorApiReturn{
		Id:   "83ae92f8-ed8b-11eb-9a03-0242ac130003",
		Name: "Finance",
	}

	sectorInfo2 := SectorApiReturn{
		Id:   "83838383-ed8b-11eb-9a03-0242ac130003",
		Name: "Health",
	}

	var assets = []AssetQueryReturn{
		{
			Id:         "0a52d206-ed8b-11eb-9a03-0242ac130003",
			Symbol:     "ITUB4",
			Preference: &preference,
			Fullname:   "Itau Unibanco Holding SA",
			Sector:     &sectorInfo,
			OrderInfo:  &ordersInfo,
		},
		{
			Id:         "11111111-ed8b-11eb-9a03-0242ac130003",
			Symbol:     "FLRY3",
			Preference: &preference2,
			Fullname:   "Fleury SA",
			Sector:     &sectorInfo2,
			OrderInfo:  &ordersInfo,
		},
	}

	expectedAssetTypeInfo := []AssetTypeApiReturn{
		{
			Id:      "00000000-ed8b-11eb-9a03-0242ac130003",
			Type:    "STOCK",
			Country: "BR",
			Name:    "Ações Brasil",
			Assets:  assets,
		},
	}

	query := regexp.QuoteMeta(`
		select
			f_query.atid as id, f_query.attype as type, f_query.atname as name,
			f_query.atcountry as country,
			json_agg(
				json_build_object(
					'id', f_query.id,
					'symbol', f_query.symbol,
					'preference', f_query.preference,
					'fullname', f_query.fullname,
					'sector', f_query.sector,
					'orderInfo', f_query.order_info
				)
			) as assets
		from (
			select
				a2.id, a2.symbol, a2.preference, a2.fullname,
				a2.atid, a2.attype, a2.atname, a2.atcountry,
				json_build_object('id', a2.sid, 'name', a2.sname) as sector,
				json_build_object(
					'totalQuantity', sum(o.quantity),
					'weightedAdjPrice', SUM(o.quantity * price)/SUM(o.quantity),
					'weightedAveragePrice', (
						SUM(o.quantity*o.price)
						FILTER(WHERE o.order_type = 'buy'))
						/(SUM(o.quantity) FILTER(WHERE o.order_type = 'buy'))
			) as order_info
			from (
				select
					a.id, a.symbol, a.preference, a.fullname,
					s.id as sid, s.name as sname,
					at.id as atid, at.type as attype, at.name as atname,
					at.country as atcountry
				from asset as a
				inner join assettype as at
				on at.id = a.asset_type_id
				inner join sector as s
				on s.id = a.sector_id
				where at.type = $1 and at.country = $2
			) a2
			inner join orders as o
			on o.asset_id = a2.id
			group by a2.id, a2.symbol, a2.preference, a2.fullname, a2.sid,
				a2.sname, a2.atid, a2.attype, a2.atname, a2.atcountry
		) as f_query
		group by f_query.atid, f_query.attype, f_query.atcountry, f_query.atname
		`)

	assetTypeInfo, errorMock := testSearchAssetPerAssetType("STOCK", "BR",
		"Ações Brasil", true, assets, query)
	if errorMock != nil {
		t.Fatalf(errorMock.Error())
	}

	assert.NotNil(t, assetTypeInfo)
	assert.Equal(t, expectedAssetTypeInfo, assetTypeInfo)

}

func TestSearchAssetPerAssetTypeWithOrderAndIsETForFII(t *testing.T) {

	ordersInfo := OrderGeneralInfos{
		TotalQuantity:        25,
		WeightedAdjPrice:     37.37,
		WeightedAveragePrice: 37.37,
	}

	var assets = []AssetQueryReturn{
		{
			Id:        "0a52d206-ed8b-11eb-9a03-0242ac130003",
			Symbol:    "VTI",
			Fullname:  "Vanguard Total Market US",
			OrderInfo: &ordersInfo,
		},
		{
			Id:        "11111111-ed8b-11eb-9a03-0242ac130003",
			Symbol:    "IJR",
			Fullname:  "iShares S&P 600 Small-Caps",
			OrderInfo: &ordersInfo,
		},
	}

	expectedAssetTypeInfo := []AssetTypeApiReturn{
		{
			Id:      "00000000-ed8b-11eb-9a03-0242ac130003",
			Type:    "ETF",
			Country: "US",
			Name:    "ETFs EUA",
			Assets:  assets,
		},
	}

	query := regexp.QuoteMeta(`
	select
		f_query.atid as id,
		f_query.attype as type,
		f_query.atname as name,
		f_query.atcountry as country,
		json_agg(
			json_build_object(
				'id', f_query.id,
				'symbol', f_query.symbol,
				'preference', f_query.preference,
				'fullname', f_query.fullname,
				'orderInfo', f_query.order_info
			)
		) as assets
	from (
		select
			a2.id, a2.symbol, a2.preference, a2.fullname,
			a2.atid, a2.attype, a2.atname, a2.atcountry,
			json_build_object(
				'totalQuantity', sum(o.quantity),
				'weightedAdjPrice', SUM(o.quantity * price)/SUM(o.quantity),
				'weightedAveragePrice', (
					SUM(o.quantity*o.price) FILTER(WHERE o.order_type = 'buy'))
					/(SUM(o.quantity) FILTER(WHERE o.order_type = 'buy'))
		) as order_info
		from (
			select
				a.id, a.symbol, a.preference, a.fullname,
				at.id as atid, at.type as attype, at.name as atname,
				at.country as atcountry
			from asset as a
			inner join assettype as at
			on at.id = a.asset_type_id
			where at.type = $1 and at.country = $2
		) a2
		inner join orders as o
		on o.asset_id = a2.id
		group by a2.id, a2.symbol, a2.preference, a2.fullname,
			a2.atid, a2.attype, a2.atname, a2.atcountry
	) as f_query
	group by f_query.atid, f_query.attype, f_query.atcountry, f_query.atname
	`)

	assetTypeInfo, errorMock := testSearchAssetPerAssetType("ETF", "US",
		"ETFs EUA", true, assets, query)
	if errorMock != nil {
		t.Fatalf(errorMock.Error())
	}

	assert.NotNil(t, assetTypeInfo)
	assert.Equal(t, expectedAssetTypeInfo, assetTypeInfo)

}

func testSearchAssetPerAssetType(assetType string, country string,
	assetTypeName string, withOrders bool, assets []AssetQueryReturn, query string) (
	[]AssetTypeApiReturn, error) {
	columns := []string{"id", "type", "country", "name", "assets"}

	var err2 error
	var assetTypeInfo []AssetTypeApiReturn

	mock, err := pgxmock.NewConn()
	if err != nil {
		return assetTypeInfo, err
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs(assetType, country).WillReturnRows(
		rows.AddRow("00000000-ed8b-11eb-9a03-0242ac130003", assetType, country,
			assetTypeName, assets))

	assetTypeInfo = SearchAssetsPerAssetType(mock, assetType, country,
		withOrders)
	if assetTypeInfo[0].Id == "" {
		return assetTypeInfo, errors.New("Wrong Query")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		return assetTypeInfo, err
	}

	return assetTypeInfo, err2
}
