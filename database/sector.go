package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
	_ "github.com/lib/pq"
)

func CreateSector(dbpool pgxIface, sector string) ([]SectorApiReturn, error) {

	var sectorInfo []SectorApiReturn
	var err error

	if sector == "" {
		err = errors.New("CreateSector: Impossible to create a blank sector")
		return sectorInfo, err
	}

	var sectorQuery = `
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
	`

	err = pgxscan.Select(context.Background(), dbpool, &sectorInfo,
		sectorQuery, sector)
	if err != nil {
		panic(err)
	}
	fmt.Println(sectorInfo)

	return sectorInfo, err
}

func FetchSector(dbpool pgxIface, sector string) ([]SectorApiReturn, error) {

	var sectorQuery []SectorApiReturn
	var dbReturnError error

	query := `
	SELECT
		id, name
	FROM sector
	`
	if sector != "ALL" {
		query = query + "where name='" + sector + "'"

	}

	err := pgxscan.Select(context.Background(), dbpool, &sectorQuery,
		query)
	if err != nil {
		fmt.Println(err)
	}

	if sectorQuery == nil {
		dbReturnError = errors.New("FetchSector: Nonexistent sector in the database")
	}

	return sectorQuery, dbReturnError
}
