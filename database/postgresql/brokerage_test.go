package postgresql

import (
	"context"
	"regexp"
	"stockfyApi/database"
	"testing"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

// func testFetchBrokerage(query string)

func TestFetchBrokerageWithName(t *testing.T) {

	expectedBrokerageInfo := []database.Brokerage{
		{
			Id:      "55555555-ed8b-11eb-9a03-0242ac130003",
			Name:    "Clear",
			Country: "BR",
		},
	}

	query := regexp.QuoteMeta(`
	SELECT
		id, name, country
	FROM brokerage
	where name=$1
	`)

	columns := []string{"id", "name", "country"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs("Clear").WillReturnRows(rows.AddRow(
		"55555555-ed8b-11eb-9a03-0242ac130003", "Clear", "BR"))

	Broker := repo{dbpool: mock}

	brokerageInfos, err := Broker.FetchBrokerage("SINGLE", "Clear")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, brokerageInfos)
	assert.Equal(t, expectedBrokerageInfo, brokerageInfos)

}

func TestFetchBrokerageWithCountry(t *testing.T) {
	expectedBrokerageInfo := []database.Brokerage{
		{
			Id:      "55555555-ed8b-11eb-9a03-0242ac130003",
			Name:    "Clear",
			Country: "BR",
		},
		{
			Id:      "55556666-ed8b-11eb-9a03-0242ac130003",
			Name:    "Rico",
			Country: "BR",
		},
	}

	query := regexp.QuoteMeta(`
	SELECT
		id, name, country
	FROM brokerage
	where country=$1
	`)

	columns := []string{"id", "name", "country"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs("BR").WillReturnRows(rows.AddRow(
		"55555555-ed8b-11eb-9a03-0242ac130003", "Clear", "BR").AddRow(
		"55556666-ed8b-11eb-9a03-0242ac130003", "Rico", "BR"))

	Broker := repo{dbpool: mock}
	brokerageInfos, err := Broker.FetchBrokerage("COUNTRY", "BR")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, brokerageInfos)
	assert.Equal(t, expectedBrokerageInfo, brokerageInfos)
}

func TestFetchBrokerageAll(t *testing.T) {
	expectedBrokerageInfo := []database.Brokerage{
		{
			Id:      "55555555-ed8b-11eb-9a03-0242ac130003",
			Name:    "Clear",
			Country: "BR",
		},
		{
			Id:      "55556666-ed8b-11eb-9a03-0242ac130003",
			Name:    "Rico",
			Country: "BR",
		},
		{
			Id:      "15151515-ed8b-11eb-9a03-0242ac130003",
			Name:    "Avenue",
			Country: "US",
		},
	}

	query := regexp.QuoteMeta(`
	SELECT
		id, name, country
	FROM brokerage
	`)

	columns := []string{"id", "name", "country"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WillReturnRows(rows.AddRow(
		"55555555-ed8b-11eb-9a03-0242ac130003", "Clear", "BR").AddRow(
		"55556666-ed8b-11eb-9a03-0242ac130003", "Rico", "BR").AddRow(
		"15151515-ed8b-11eb-9a03-0242ac130003", "Avenue", "US"))

	Broker := repo{dbpool: mock}
	brokerageInfos, err := Broker.FetchBrokerage("ALL")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, brokerageInfos)
	assert.Equal(t, expectedBrokerageInfo, brokerageInfos)
}
