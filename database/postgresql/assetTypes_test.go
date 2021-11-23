package postgresql

import (
	"context"
	"fmt"
	"regexp"
	"stockfyApi/entity"
	"testing"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

func TestAssetTypeSearchAll(t *testing.T) {

	expectedAssetTypes := []entity.AssetType{
		{
			Id:      "1",
			Type:    "STOCK",
			Country: "US",
			Name:    "Ações EUA",
		},
		{
			Id:      "2",
			Type:    "STOCK",
			Country: "BR",
			Name:    "Ações Brasil",
		},
		{
			Id:      "3",
			Type:    "ETF",
			Country: "US",
			Name:    "ETFs EUA",
		},
		{
			Id:      "4",
			Type:    "ETF",
			Country: "BR",
			Name:    "ETFs Brasil",
		},
		{
			Id:      "5",
			Type:    "REIT",
			Country: "US",
			Name:    "REITs",
		},
		{
			Id:      "6",
			Type:    "FII",
			Country: "BR",
			Name:    "FIIs",
		},
	}

	query := "SELECT id, type, name, country FROM asset_types"

	columns := []string{"id", "type", "name", "country"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WillReturnRows(rows.AddRow(
		"1", "STOCK", "Ações EUA", "US").AddRow("2", "STOCK", "Ações Brasil",
		"BR").AddRow("3", "ETF", "ETFs EUA", "US").AddRow("4", "ETF", "ETFs Brasil", "BR").
		AddRow("5", "REIT", "REITs", "US").AddRow("6", "FII", "FIIs", "BR"))

	At := AssetTypePostgres{dbpool: mock}
	assetType, err := At.Search("", "", "")
	if err != nil {
		fmt.Println(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, assetType)
	assert.Equal(t, expectedAssetTypes, assetType)
}

func TestAssetTypeSingleSearch(t *testing.T) {
	expectedAssetTypes := []entity.AssetType{
		{
			Id:      "1",
			Type:    "STOCK",
			Country: "US",
			Name:    "Ações EUA",
		},
	}

	query := regexp.QuoteMeta(`
				SELECT
					id, type, name, country
				FROM asset_types
				where type=$1 and country=$2
			`)

	columns := []string{"id", "type", "name", "country"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs("STOCK", "US").WillReturnRows(
		rows.AddRow("1", "STOCK", "Ações EUA", "US"))

	At := AssetTypePostgres{dbpool: mock}

	assetType, err := At.Search("SPECIFIC", expectedAssetTypes[0].Type,
		expectedAssetTypes[0].Country)
	if err != nil {
		fmt.Println(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, assetType)
	assert.Equal(t, expectedAssetTypes, assetType)
}

func TestAssetTypePerCountry(t *testing.T) {
	country := "US"
	expectedAssetTypes := []entity.AssetType{
		{
			Id:      "1",
			Type:    "STOCK",
			Country: "US",
			Name:    "Ações EUA",
		},
		{
			Id:      "3",
			Type:    "ETF",
			Country: "US",
			Name:    "ETFs EUA",
		},
		{
			Id:      "5",
			Type:    "REIT",
			Country: "US",
			Name:    "REITs",
		},
	}

	query := regexp.QuoteMeta(`
				SELECT
					id, type, name, country
				FROM asset_types
				where country=$1
			`)

	columns := []string{"id", "type", "name", "country"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs("US").WillReturnRows(rows.AddRow(
		"1", "STOCK", "Ações EUA", "US").AddRow("3", "ETF", "ETFs EUA", "US").
		AddRow("5", "REIT", "REITs", "US"))

	At := AssetTypePostgres{dbpool: mock}
	assetType, err := At.Search("ONLYCOUNTRY", "", country)
	if err != nil {
		fmt.Println(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, assetType)
	assert.Equal(t, expectedAssetTypes, assetType)
}

func TestAssetTypeSearchPerType(t *testing.T) {
	name := "STOCK"
	expectedAssetTypes := []entity.AssetType{
		{
			Id:      "1",
			Type:    "STOCK",
			Country: "US",
			Name:    "Ações EUA",
		},
		{
			Id:      "2",
			Type:    "STOCK",
			Country: "BR",
			Name:    "Ações Brasil",
		},
	}

	query := regexp.QuoteMeta(`
				SELECT
					id, type, name, country
				FROM asset_types
				where type=$1
			`)

	columns := []string{"id", "type", "name", "country"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs("STOCK").WillReturnRows(rows.AddRow(
		"1", "STOCK", "Ações EUA", "US").AddRow("2", "STOCK", "Ações Brasil",
		"BR"))

	At := AssetTypePostgres{dbpool: mock}

	assetType, err := At.Search("ONLYTYPE", name, "")
	if err != nil {
		fmt.Println(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, assetType)
	assert.Equal(t, expectedAssetTypes, assetType)
}
