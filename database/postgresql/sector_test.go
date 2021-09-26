package postgresql

import (
	"context"
	"regexp"
	"stockfyApi/entity"
	"testing"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

func TestSectorCreate(t *testing.T) {

	var expectedSectorInfo = []entity.Sector{
		{
			Id:   "0a52d206-ed8b-11eb-9a03-0242ac130003",
			Name: "Finance",
		},
	}

	query := regexp.QuoteMeta(`
	WITH s as (
		SELECT
			id, name
		FROM sector
		WHERE name=$1
	), i as (
		INSERT INTO
			sector(name)
		SELECT $1
		WHERE NOT EXISTS (SELECT 1 FROM s)
		returning id, name
	)
	SELECT
		id, name from i
	UNION ALL
	SELECT
		id, name
	from s;
	`)

	columns := []string{"id", "name"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs("Finance").WillReturnRows(
		rows.AddRow("0a52d206-ed8b-11eb-9a03-0242ac130003", "Finance"))

	Sector := SectorPostgres{dbpool: mock}

	sectorInfo, _ := Sector.Create("Finance")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, sectorInfo)
	assert.Equal(t, expectedSectorInfo, sectorInfo)
}

func TestSectorSingleSearch(t *testing.T) {

	var expectedSectorInfo = []entity.Sector{
		{
			Id:   "0a52d206-ed8b-11eb-9a03-0242ac130003",
			Name: "Finance",
		},
	}

	query := `
	SELECT
		id, name
	FROM sector
	where name='Finance'
	`

	columns := []string{"id", "name"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WillReturnRows(
		rows.AddRow("0a52d206-ed8b-11eb-9a03-0242ac130003", "Finance"))

	Sector := SectorPostgres{dbpool: mock}

	sectorInfo, _ := Sector.SearchByName("Finance")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, sectorInfo)
	assert.Equal(t, expectedSectorInfo, sectorInfo)
}

func TestSectorSearchAll(t *testing.T) {
	var expectedSectorInfo = []entity.Sector{
		{
			Id:   "0a52d206-ed8b-11eb-9a03-0242ac130003",
			Name: "Finance",
		},
		{
			Id:   "62d4d8e2-95e5-4144-b17b-0d147c98d85c",
			Name: "Technology",
		},
	}

	query := `
	SELECT
		id, name
	FROM sector$
	`

	columns := []string{"id", "name"}

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WillReturnRows(
		rows.AddRow("0a52d206-ed8b-11eb-9a03-0242ac130003", "Finance").
			AddRow("62d4d8e2-95e5-4144-b17b-0d147c98d85c", "Technology"))

	Sector := SectorPostgres{dbpool: mock}

	sectorInfo, _ := Sector.SearchByName("ALL")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, sectorInfo)
	assert.Equal(t, expectedSectorInfo, sectorInfo)

}
