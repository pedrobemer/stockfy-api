package postgresql

import (
	"context"
	"regexp"
	"stockfyApi/entity"
	"testing"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

// func testFetchBrokerage(query string)

func TestBrokerageSearchWithName(t *testing.T) {

	expectedBrokerageInfo := []entity.Brokerage{
		{
			Id:      "55555555-ed8b-11eb-9a03-0242ac130003",
			Name:    "Clear",
			Country: "BR",
		},
	}

	query := regexp.QuoteMeta(`
	SELECT
		id, name, country
	FROM brokerages
	where name=$1
	`)

	columns := []string{"id", "name", "country"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs("Clear").WillReturnRows(rows.AddRow(
		"55555555-ed8b-11eb-9a03-0242ac130003", "Clear", "BR"))

	Broker := BrokeragePostgres{dbpool: mock}

	brokerageInfos, err := Broker.Search("SINGLE", "Clear")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, brokerageInfos)
	assert.Equal(t, expectedBrokerageInfo, brokerageInfos)

}

func TestBrokerageSearchWithCountry(t *testing.T) {
	expectedBrokerageInfo := []entity.Brokerage{
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
	FROM brokerages
	where country=$1
	`)

	columns := []string{"id", "name", "country"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs("BR").WillReturnRows(rows.AddRow(
		"55555555-ed8b-11eb-9a03-0242ac130003", "Clear", "BR").AddRow(
		"55556666-ed8b-11eb-9a03-0242ac130003", "Rico", "BR"))

	Broker := BrokeragePostgres{dbpool: mock}
	brokerageInfos, err := Broker.Search("COUNTRY", "BR")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, brokerageInfos)
	assert.Equal(t, expectedBrokerageInfo, brokerageInfos)
}

func TestBrokerageSearchAll(t *testing.T) {
	expectedBrokerageInfo := []entity.Brokerage{
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
	FROM brokerages
	`)

	columns := []string{"id", "name", "country"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WillReturnRows(rows.AddRow(
		"55555555-ed8b-11eb-9a03-0242ac130003", "Clear", "BR").AddRow(
		"55556666-ed8b-11eb-9a03-0242ac130003", "Rico", "BR").AddRow(
		"15151515-ed8b-11eb-9a03-0242ac130003", "Avenue", "US"))

	Broker := BrokeragePostgres{dbpool: mock}
	brokerageInfos, err := Broker.Search("ALL")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, brokerageInfos)
	assert.Equal(t, expectedBrokerageInfo, brokerageInfos)
}
