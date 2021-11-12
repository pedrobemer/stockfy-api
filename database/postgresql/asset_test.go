package postgresql

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"stockfyApi/entity"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

func TestAssetSearch(t *testing.T) {

	symbol := "ITUB4"

	assetType := entity.AssetType{
		Id:      "28ccf27a-ed8b-11eb-9a03-0242ac130003",
		Type:    "STOCK",
		Name:    "Ações Brasil",
		Country: "BR",
	}
	preference := "ON"

	sectorInfo := entity.Sector{
		Id:   "83ae92f8-ed8b-11eb-9a03-0242ac130003",
		Name: "Finance",
	}

	var expectedAsset = []entity.Asset{
		{
			Id:         "0a52d206-ed8b-11eb-9a03-0242ac130003",
			Symbol:     symbol,
			Preference: &preference,
			Fullname:   "Itau Unibanco Holding SA",
			AssetType:  &assetType,
			Sector:     &sectorInfo,
		},
	}

	query := regexp.QuoteMeta(`
	SELECT
		a.id, symbol, preference, fullname,
		json_build_object(
			'id', aty.id,
			'type', aty."type",
			'name', aty."name",
			'country', aty.country
		) as asset_type,
		json_build_object(
			'id', s.id,
			'name', s."name"
		) as sector
	FROM assets as a
	INNER JOIN assettypes as aty
	ON aty.id = a.asset_type_id
	INNER JOIN sectors as s
	ON s.id = a.sector_id
	WHERE symbol=$1;
	`)

	columns := []string{"id", "symbol", "preference", "fullname", "asset_type",
		"sector"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs(symbol).WillReturnRows(
		rows.AddRow("0a52d206-ed8b-11eb-9a03-0242ac130003", "ITUB4", &preference,
			"Itau Unibanco Holding SA", &assetType, &sectorInfo))

	Asset := AssetPostgres{dbpool: mock}

	asset, err := Asset.Search(symbol)
	if err != nil {
		fmt.Println(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, asset)
	assert.Equal(t, expectedAsset, asset)
}

func TestAssetSingleSearchByUser(t *testing.T) {

	symbol := "ITUB4"
	userUid := "afauaf4s29f"

	assetType := entity.AssetType{
		Id:      "28ccf27a-ed8b-11eb-9a03-0242ac130003",
		Type:    "STOCK",
		Name:    "Ações Brasil",
		Country: "BR",
	}
	preference := "ON"

	sectorInfo := entity.Sector{
		Id:   "83ae92f8-ed8b-11eb-9a03-0242ac130003",
		Name: "Finance",
	}

	var expectedAsset = []entity.Asset{
		{
			Id:         "0a52d206-ed8b-11eb-9a03-0242ac130003",
			Symbol:     symbol,
			Preference: &preference,
			Fullname:   "Itau Unibanco Holding SA",
			AssetType:  &assetType,
			Sector:     &sectorInfo,
		},
	}

	query := regexp.QuoteMeta(`
	SELECT
		a.id, symbol, preference, fullname,
		json_build_object(
			'id', aty.id,
			'type', aty."type",
			'name', aty."name",
			'country', aty.country
		) as asset_type,
		json_build_object(
			'id', s.id,
			'name', s."name"
		) as sector
	FROM asset_users as au
	INNER JOIN assets as a
	ON a.id = au.asset_id
	INNER JOIN assettypes as aty
	ON aty.id = a.asset_type_id
	INNER JOIN sectors as s
	ON s.id = a.sector_id
	WHERE a.symbol=$1 and au.user_uid=$2
	GROUP BY a.symbol, a.id, a.preference, a.fullname, aty.id, aty."type",
	aty."name", aty.country, s.id, s."name";
	`)

	columns := []string{"id", "symbol", "preference", "fullname", "asset_type",
		"sector"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs(symbol, userUid).WillReturnRows(
		rows.AddRow("0a52d206-ed8b-11eb-9a03-0242ac130003", "ITUB4", &preference,
			"Itau Unibanco Holding SA", &assetType, &sectorInfo))

	Asset := AssetPostgres{dbpool: mock}

	asset, err := Asset.SearchByUser(symbol, userUid, "")
	if err != nil {
		fmt.Println(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, asset)
	assert.Equal(t, expectedAsset, asset)
}

func TestAssetSingleSearchByUserWithOrders(t *testing.T) {

	symbol := "ITUB4"
	userUid := "afauaf4s29f"

	tr, err := time.Parse("2021-07-05", "2021-07-21")
	tr2, err := time.Parse("2021-07-05", "2020-04-02")

	assetType := entity.AssetType{
		Id:      "28ccf27a-ed8b-11eb-9a03-0242ac130003",
		Type:    "STOCK",
		Name:    "Ações Brasil",
		Country: "BR",
	}

	preference := "ON"

	brokerageInfo := entity.Brokerage{
		Id:      "55555555-ed8b-11eb-9a03-0242ac130003",
		Name:    "Clear",
		Country: "BR",
	}

	orderList := []entity.Order{
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

	sectorInfo := entity.Sector{
		Id:   "83ae92f8-ed8b-11eb-9a03-0242ac130003",
		Name: "Finance",
	}

	var expectedAsset = []entity.Asset{
		{
			Id:         "0a52d206-ed8b-11eb-9a03-0242ac130003",
			Symbol:     symbol,
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
	json_build_object(
		'id', s.id,
		'name', s."name"
	) as sector,
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
	FROM asset_users as au
	INNER JOIN assets as a
	ON a.id = au.asset_id
	INNER JOIN assettypes as at
	ON a.asset_type_id = at.id
	INNER JOIN sectors as s
	ON s.id = a.sector_id
	INNER JOIN orders as o
	ON a.id = o.asset_id and au.user_uid = o.user_uid
	INNER JOIN brokerages as b
	ON o.brokerage_id = b.id
	WHERE a.symbol=$1 and au.user_uid =$2
	GROUP BY a.symbol, a.id, preference, a.fullname, at.type, at.id,
	at.name, at.country, s.id, s.name;
	`)

	columnsAsset := []string{"id", "symbol", "preference", "fullname",
		"asset_type", "orders_list", "sector"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columnsAsset)
	mock.ExpectQuery(query).WithArgs(symbol, userUid).WillReturnRows(
		rows.AddRow("0a52d206-ed8b-11eb-9a03-0242ac130003", "ITUB4", &preference,
			"Itau Unibanco Holding SA", &assetType, orderList, &sectorInfo))

	Asset := AssetPostgres{dbpool: mock}

	asset, err := Asset.SearchByUser("ITUB4", userUid, "ONLYORDERS")
	if err != nil {
		fmt.Println(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, asset)
	assert.Equal(t, expectedAsset, asset)

}

func TestAssetSingleSearchByUserWithOrderInfo(t *testing.T) {

	symbol := "ITUB4"
	userUid := "afauaf4s29f"

	assetType := entity.AssetType{
		Id:      "28ccf27a-ed8b-11eb-9a03-0242ac130003",
		Type:    "STOCK",
		Name:    "Ações Brasil",
		Country: "BR",
	}

	preference := "ON"

	ordersInfo := entity.OrderInfos{
		TotalQuantity:        25,
		WeightedAdjPrice:     37.37,
		WeightedAveragePrice: 37.37,
	}

	sectorInfo := entity.Sector{
		Id:   "83ae92f8-ed8b-11eb-9a03-0242ac130003",
		Name: "Finance",
	}

	var expectedAsset = []entity.Asset{
		{
			Id:         "0a52d206-ed8b-11eb-9a03-0242ac130003",
			Symbol:     symbol,
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
		'id', aty.id,
		'type', aty.type,
		'name', aty.name,
		'country', aty.country
	) as asset_type,
	json_build_object(
		'id', s.id,
		'name', s."name"
	) as sector,
	json_build_object(
		'totalQuantity', sum(o.quantity),
		'weightedAdjPrice', SUM(o.quantity * price)/SUM(o.quantity),
		'weightedAveragePrice', (
			SUM(o.quantity*o.price) FILTER(WHERE o.order_type = 'buy'))
			/(SUM(o.quantity) FILTER(WHERE o.order_type = 'buy')
		)
	) as orders_info
	FROM asset_users as au
	INNER JOIN assets as a
	ON a.id = au.asset_id
	INNER JOIN assettypes as aty
	ON a.asset_type_id = aty.id
	INNER JOIN sectors as s
	ON s.id = a.sector_id
	INNER JOIN orders as o
	ON a.id = o.asset_id and au.user_uid = o.user_uid
	INNER JOIN brokerages as b
	ON o.brokerage_id = b.id
	WHERE a.symbol=$1 and au.user_uid =$2
	GROUP BY a.symbol, a.id, preference, a.fullname, aty.type, aty.id,
	aty.name, aty.country, s.id, s.name;
	`)

	columnsAsset := []string{"id", "symbol", "preference", "fullname",
		"asset_type", "orders_info", "sector"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columnsAsset)
	mock.ExpectQuery(query).WithArgs(symbol, userUid).WillReturnRows(
		rows.AddRow("0a52d206-ed8b-11eb-9a03-0242ac130003", "ITUB4", &preference,
			"Itau Unibanco Holding SA", &assetType, &ordersInfo, &sectorInfo))

	Asset := AssetPostgres{dbpool: mock}

	asset, err := Asset.SearchByUser(symbol, userUid, "ONLYINFO")
	if err != nil {
		fmt.Println(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, asset)
	assert.Equal(t, expectedAsset, asset)

}

func TestAssetSingleSearchAllInfo(t *testing.T) {
	tr, err := time.Parse("2021-07-05", "2021-07-21")
	tr2, err := time.Parse("2021-07-05", "2020-04-02")

	symbol := "ITUB4"
	userUid := "afauaf4s29f"

	assetType := entity.AssetType{
		Id:      "28ccf27a-ed8b-11eb-9a03-0242ac130003",
		Type:    "STOCK",
		Name:    "Ações Brasil",
		Country: "BR",
	}

	preference := "ON"

	brokerageInfo := entity.Brokerage{
		Id:      "55555555-ed8b-11eb-9a03-0242ac130003",
		Name:    "Clear",
		Country: "BR",
	}

	orderList := []entity.Order{
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

	ordersInfo := entity.OrderInfos{
		TotalQuantity:        25,
		WeightedAdjPrice:     37.37,
		WeightedAveragePrice: 37.37,
	}

	sectorInfo := entity.Sector{
		Id:   "83ae92f8-ed8b-11eb-9a03-0242ac130003",
		Name: "Finance",
	}

	var expectedAsset = []entity.Asset{
		{
			Id:         "0a52d206-ed8b-11eb-9a03-0242ac130003",
			Symbol:     symbol,
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
			'id', s.id,
			'name', s."name"
		) as sector,
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
	FROM asset_users as au
	INNER JOIN assets as a
	ON a.id = au.asset_id
	INNER JOIN assettypes as at
	ON a.asset_type_id = at.id
	INNER JOIN sectors as s
	ON s.id = a.sector_id
	INNER JOIN orders as o
	ON a.id = o.asset_id and au.user_uid = o.user_uid
	INNER JOIN brokerages as b
	ON o.brokerage_id = b.id
	WHERE a.symbol=$1 and au.user_uid =$2
	GROUP BY a.symbol, a.id, preference, a.fullname, at.type, at.id,
	at.name, at.country, s.id, s.name;
	`)

	columnsAsset := []string{"id", "symbol", "preference", "fullname",
		"asset_type", "sector", "orders_info", "orders_list"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columnsAsset)
	mock.ExpectQuery(query).WithArgs(symbol, userUid).WillReturnRows(
		rows.AddRow("0a52d206-ed8b-11eb-9a03-0242ac130003", "ITUB4", &preference,
			"Itau Unibanco Holding SA", &assetType, &sectorInfo, &ordersInfo,
			orderList))

	Asset := AssetPostgres{dbpool: mock}

	asset, err := Asset.SearchByUser("ITUB4", userUid, "ALL")
	if err != nil {
		fmt.Println(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, asset)
	assert.Equal(t, expectedAsset, asset)

}

func TestAssetSearchByOrderId(t *testing.T) {
	assetType := entity.AssetType{
		Id:      "28ccf27a-ed8b-11eb-9a03-0242ac130003",
		Type:    "STOCK",
		Name:    "Ações Brasil",
		Country: "BR",
	}

	preference := "ON"

	var expectedAssetInfo = []entity.Asset{
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
	inner join assets as a
	on a.id = o.asset_id
	inner join assettypes as aty
	on aty.id = a.asset_type_id
	where o.id = $1;
	`)

	columns := []string{"id", "preference", "symbol", "asset_type"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs("6669aaaa-ed8b-11eb-9a03-0242ac130003").
		WillReturnRows(rows.AddRow(
			"0a52d206-ed8b-11eb-9a03-0242ac130003", &preference, "ITUB4",
			&assetType))

	Asset := AssetPostgres{dbpool: mock}
	assetInfo := Asset.SearchByOrderId(
		"6669aaaa-ed8b-11eb-9a03-0242ac130003")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, assetInfo)
	assert.Equal(t, expectedAssetInfo, assetInfo)
}

func TestAssetCreate(t *testing.T) {

	assetType := entity.AssetType{
		Id: "0a52d206-ed8b-11eb-9a03-0242ac130003",
	}

	sectorInfo := entity.Sector{
		Id: "83ae92f8-ed8b-11eb-9a03-0242ac130003",
	}
	preference := "PN"

	asset := entity.Asset{
		AssetType:  &assetType,
		Sector:     &sectorInfo,
		Preference: &preference,
		Symbol:     "ITUB4",
		Fullname:   "Itau Unibanco Holding SA",
	}

	expectedAssetReturn := entity.Asset{
		Id:         "000aaaa6-ed8b-11eb-9a03-0242ac130003",
		Preference: &preference,
		Fullname:   "Itau Unibanco Holding SA",
		Symbol:     "ITUB4",
	}

	insertRow := regexp.QuoteMeta(`
	INSERT INTO
		assets(preference, fullname, symbol, asset_type_id, sector_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, preference, fullname, symbol;
	`)

	columns := []string{"id", "preference", "fullname", "symbol"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectBegin()
	mock.ExpectQuery(insertRow).WithArgs(&preference, "Itau Unibanco Holding SA",
		"ITUB4", "0a52d206-ed8b-11eb-9a03-0242ac130003",
		"83ae92f8-ed8b-11eb-9a03-0242ac130003").WillReturnRows(
		rows.AddRow("000aaaa6-ed8b-11eb-9a03-0242ac130003", "PN",
			"Itau Unibanco Holding SA", "ITUB4"))
	mock.ExpectCommit()

	Asset := AssetPostgres{dbpool: mock}

	assetRtr := Asset.Create(asset)

	assert.NotNil(t, assetRtr)
	assert.Equal(t, expectedAssetReturn, assetRtr)
}

func TestAssetDelete(t *testing.T) {

	preference := "PN"

	var expectedDelAsset = []entity.Asset{
		{
			Id:         "0a52d206-ed8b-11eb-9a03-0242ac130003",
			Symbol:     "ITUB4",
			Preference: &preference,
			Fullname:   "Itau Unibanco Holding SA",
		},
	}

	queryDeleteAsset := regexp.QuoteMeta(`
	delete from assets as a
	where a.id = $1
	returning  a.id, a.symbol, a.preference, a.fullname;
	`)

	columnsDelAsset := []string{"id", "symbol", "preference", "fullname"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rowsDelAsset := mock.NewRows(columnsDelAsset)

	mock.ExpectQuery(queryDeleteAsset).WithArgs(
		"0a52d206-ed8b-11eb-9a03-0242ac130003").WillReturnRows(
		rowsDelAsset.AddRow("0a52d206-ed8b-11eb-9a03-0242ac130003", "ITUB4",
			&preference, "Itau Unibanco Holding SA"))

	Asset := AssetPostgres{dbpool: mock}
	assetInfo, err := Asset.Delete("0a52d206-ed8b-11eb-9a03-0242ac130003")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, assetInfo)
	assert.Nil(t, err)
	assert.Equal(t, expectedDelAsset, assetInfo)
}

func TestAssetSearchPerAssetTypeWithoutOrderInfo(t *testing.T) {

	preference := "PN"
	preference2 := "ON"

	assetType := "STOCK"
	country := "BR"
	name := "Ações Brasil"
	userUid := "afauaf4s29f"

	sectorInfo := entity.Sector{
		Id:   "83ae92f8-ed8b-11eb-9a03-0242ac130003",
		Name: "Finance",
	}

	sectorInfo2 := entity.Sector{
		Id:   "83838383-ed8b-11eb-9a03-0242ac130003",
		Name: "Health",
	}

	assets := []entity.Asset{
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

	expectedAssetTypeInfo := []entity.AssetType{
		{
			Id:      "00000000-ed8b-11eb-9a03-0242ac130003",
			Type:    assetType,
			Country: country,
			Name:    name,
			Assets:  assets,
		},
	}

	query := regexp.QuoteMeta(`
	SELECT
		aty.id, aty.type, aty.country, aty.name,
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
	FROM asset_users as au
	INNER JOIN assets as a
	ON a.id = au.asset_id
	INNER JOIN assettypes as aty
	ON aty.id = a.asset_type_id
	INNER JOIN sectors as s
	ON s.id = a.sector_id
	WHERE au.user_uid=$1 and aty."type"=$2 and aty.country=$3
	GROUP BY aty.id, aty."type", aty."name", aty.country;
	`)

	assetTypeInfo, err := testSearchAssetPerAssetType(userUid, assetType,
		country, name, false, assets, query)
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.NotNil(t, assetTypeInfo)
	assert.Equal(t, expectedAssetTypeInfo, assetTypeInfo)
}

func TestAssetSearchPerAssetTypeWithOrderInfo(t *testing.T) {
	preference := "PN"
	preference2 := "ON"

	assetType := "STOCK"
	country := "BR"
	name := "Ações Brasil"
	userUid := "afauaf4s29f"

	ordersInfo := entity.OrderInfos{
		TotalQuantity:        25,
		WeightedAdjPrice:     37.37,
		WeightedAveragePrice: 37.37,
	}

	sectorInfo := entity.Sector{
		Id:   "83ae92f8-ed8b-11eb-9a03-0242ac130003",
		Name: "Finance",
	}

	sectorInfo2 := entity.Sector{
		Id:   "83838383-ed8b-11eb-9a03-0242ac130003",
		Name: "Health",
	}

	var assets = []entity.Asset{
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

	expectedAssetTypeInfo := []entity.AssetType{
		{
			Id:      "00000000-ed8b-11eb-9a03-0242ac130003",
			Type:    assetType,
			Country: country,
			Name:    name,
			Assets:  assets,
		},
	}

	query := regexp.QuoteMeta(`
		SELECT
			f_query.at_id as id, f_query.at_type as type, f_query.at_name as name,
			f_query.at_country as country,
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
		FROM (
			SELECT
				valid_assets.id, valid_assets.symbol, valid_assets.preference,
				valid_assets.fullname, valid_assets.at_id, valid_assets.at_type,
				valid_assets.at_name, valid_assets.at_country,
				json_build_object(
					'id', valid_assets.s_id,
					'name', valid_assets.s_name
				) as sector,
				json_build_object(
					'totalQuantity', sum(o.quantity),
					'weightedAdjPrice', SUM(o.quantity * price)/SUM(o.quantity),
					'weightedAveragePrice', (
						SUM(o.quantity*o.price)
						FILTER(WHERE o.order_type = 'buy'))
						/(SUM(o.quantity) FILTER(WHERE o.order_type = 'buy'))
				) as order_info
			FROM (
				select
					a.id, a.symbol, a.preference, a.fullname, s.id as s_id,
					s."name" as s_name, aty.id as at_id, aty."type" as at_type,
					aty."name" as at_name, aty.country as at_country
				FROM asset_users as au
				INNER JOIN assets as a
				ON a.id = au.asset_id
				INNER JOIN assettypes as aty
				ON aty.id = a.asset_type_id
				inner join sectors as s
				on s.id = a.sector_id
				WHERE au.user_uid=$1 and aty."type"=$2 and aty.country=$3
				GROUP BY a.symbol, a.id, a.preference, a.fullname, aty.id, aty."type",
				aty."name", aty.country, s.id, s."name"
			) valid_assets
			INNER JOIN orders as o
			ON o.asset_id = valid_assets.id
			WHERE o.user_uid = $1
			GROUP BY valid_assets.id, valid_assets.symbol,
			valid_assets.preference, valid_assets.fullname, valid_assets.s_id,
			valid_assets.s_name, valid_assets.at_id, valid_assets.at_type,
			valid_assets.at_name, valid_assets.at_country
		) as f_query
		GROUP BY f_query.at_id, f_query.at_type, f_query.at_country,
		f_query.at_name;
		`)

	assetTypeInfo, errorMock := testSearchAssetPerAssetType(userUid, assetType,
		country, name, true, assets, query)
	if errorMock != nil {
		t.Fatalf(errorMock.Error())
	}

	assert.NotNil(t, assetTypeInfo)
	assert.Equal(t, expectedAssetTypeInfo, assetTypeInfo)

}

func testSearchAssetPerAssetType(userUid string, assetType string,
	country string, assetTypeName string, withOrders bool,
	assets []entity.Asset, query string) ([]entity.AssetType, error) {
	columns := []string{"id", "type", "country", "name", "assets"}

	var err2 error
	var assetTypeInfo []entity.AssetType

	mock, err := pgxmock.NewConn()
	if err != nil {
		return assetTypeInfo, err
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs(userUid, assetType, country).WillReturnRows(
		rows.AddRow("00000000-ed8b-11eb-9a03-0242ac130003", assetType, country,
			assetTypeName, assets))

	Asset := AssetPostgres{dbpool: mock}
	assetTypeInfo = Asset.SearchPerAssetType(assetType, country, userUid,
		withOrders)
	if assetTypeInfo[0].Id == "" {
		return assetTypeInfo, errors.New("Wrong Query")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		return assetTypeInfo, err
	}

	return assetTypeInfo, err2
}
