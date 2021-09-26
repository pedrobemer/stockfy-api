package postgresql

import (
	"context"
	"errors"
	"fmt"
	"stockfyApi/entity"

	"github.com/georgysavva/scany/pgxscan"
	_ "github.com/lib/pq"
)

type SectorPostgres struct {
	dbpool PgxIface
}

func NewSectorPostgres(db *PgxIface) *SectorPostgres {
	return &SectorPostgres{
		dbpool: *db,
	}
}

func (r *SectorPostgres) Create(sector string) ([]entity.Sector, error) {

	var sectorInfo []entity.Sector
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

	err = pgxscan.Select(context.Background(), r.dbpool, &sectorInfo,
		sectorQuery, sector)
	if err != nil {
		panic(err)
	}

	return sectorInfo, err
}

func (r *SectorPostgres) SearchByName(sector string) ([]entity.Sector, error) {

	var sectorQuery []entity.Sector
	var dbReturnError error

	query := `
	SELECT
		id, name
	FROM sector
	`
	if sector != "ALL" {
		query = query + "where name='" + sector + "'"

	}

	err := pgxscan.Select(context.Background(), r.dbpool, &sectorQuery,
		query)
	if err != nil {
		fmt.Println(err)
	}

	if sectorQuery == nil {
		dbReturnError = errors.New("FetchSector: Nonexistent sector in the entity")
	}

	return sectorQuery, dbReturnError
}

func (r *SectorPostgres) SearchByAsset(symbol string) ([]entity.Sector, error) {
	var sectorQuery []entity.Sector
	var dbReturnError error

	query := `
	select
		s.id,
		s.name
	from sector as s
	inner join asset as a
	on a.sector_id = s.id
	where a.symbol = $1;
	`

	err := pgxscan.Select(context.Background(), r.dbpool, &sectorQuery,
		query, symbol)
	if err != nil {
		fmt.Println(err)
	}

	if sectorQuery == nil {
		dbReturnError = errors.New("FetchSector: Nonexistent sector in the entity")
	}

	return sectorQuery, dbReturnError

}
