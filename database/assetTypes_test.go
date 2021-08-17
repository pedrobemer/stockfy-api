package database

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

func TestFetchAllAssetType(t *testing.T) {

	expectedAssetTypes := []AssetTypeApiReturn{
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

	query := "SELECT id, type, name, country FROM assettype"

	columns := []string{"id", "type", "name", "country"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WillReturnRows(rows.AddRow(
		"1", "STOCK", "Ações EUA", "US").AddRow("2", "STOCK", "Ações Brasil",
		"BR").AddRow("3", "ETF", "ETFs EUA", "US").AddRow("4", "ETF", "ETFs Brasil", "BR").
		AddRow("5", "REIT", "REITs", "US").AddRow("6", "FII", "FIIs", "BR"))

	assetType, err := FetchAssetType(mock, "")
	if err != nil {
		fmt.Println(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, assetType)
	assert.Equal(t, expectedAssetTypes, assetType)
}

func TestFetchSpecificAssetType(t *testing.T) {
	expectedAssetTypes := []AssetTypeApiReturn{
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
				FROM assettype
				where type=$1 and country=$2
			`)

	columns := []string{"id", "type", "name", "country"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs("ETF", "US").WillReturnRows(
		rows.AddRow("1", "STOCK", "Ações EUA", "US"))

	assetType, err := FetchAssetType(mock, "SPECIFIC", "ETF", "US")
	if err != nil {
		fmt.Println(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, assetType)
	assert.Equal(t, expectedAssetTypes, assetType)
}

func TestFetchAssetTypePerCountry(t *testing.T) {
	expectedAssetTypes := []AssetTypeApiReturn{
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
				FROM assettype
				where country=$1
			`)

	columns := []string{"id", "type", "name", "country"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs("US").WillReturnRows(rows.AddRow(
		"1", "STOCK", "Ações EUA", "US").AddRow("3", "ETF", "ETFs EUA", "US").
		AddRow("5", "REIT", "REITs", "US"))

	assetType, err := FetchAssetType(mock, "ONLYCOUNTRY", "AAAA", "US")
	if err != nil {
		fmt.Println(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, assetType)
	assert.Equal(t, expectedAssetTypes, assetType)
}

func TestFetchAssetTypePerType(t *testing.T) {
	expectedAssetTypes := []AssetTypeApiReturn{
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
				FROM assettype
				where type=$1
			`)

	columns := []string{"id", "type", "name", "country"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs("STOCK").WillReturnRows(rows.AddRow(
		"1", "STOCK", "Ações EUA", "US").AddRow("2", "STOCK", "Ações Brasil",
		"BR"))

	assetType, err := FetchAssetType(mock, "ONLYTYPE", "STOCK", "")
	if err != nil {
		fmt.Println(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, assetType)
	assert.Equal(t, expectedAssetTypes, assetType)
}
