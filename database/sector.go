package database

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

// Public visibility
func CreateSector(dbpool pgxpool.Pool, sector string) ([]SectorApiReturn, error) {
	var sectorInfo []SectorApiReturn
	// var sectorInfo interface{}
	var err error

	if sector == "" {
		err = errors.New("CreateSector: Impossible to create a blank sector")
		return sectorInfo, err
	}

	var sectorQuery = "WITH s as (SELECT id, name FROM sector WHERE name=$1), i as (INSERT INTO sector(name) SELECT $1 WHERE NOT EXISTS (SELECT 1 FROM s) returning id, name) SELECT id, name from i UNION ALL SELECT id, name from s;"

	err = pgxscan.Select(context.Background(), &dbpool, &sectorInfo,
		sectorQuery, sector)
	if err != nil {
		panic(err)
	}
	fmt.Println(sectorInfo)

	return sectorInfo, err
}

// Public visibility
func FetchSector(dbpool pgxpool.Pool, sector string) ([]SectorApiReturn, error) {
	var sectorQuery []SectorApiReturn
	var dbReturnError error

	query := "SELECT id, name FROM sector"
	if sector != "ALL" {
		query = query + " where name='" + sector + "'"

	}

	err := pgxscan.Select(context.Background(), &dbpool, &sectorQuery,
		query)
	if err != nil {
		log.Panic(err)
	}

	if sectorQuery == nil {
		dbReturnError = errors.New("FetchSector: Nonexistent sector in the database")
	}

	return sectorQuery, dbReturnError
}
