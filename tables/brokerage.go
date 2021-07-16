package tables

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

func FetchBrokerage(dbpool pgxpool.Pool, specificFetch string,
	args ...string) ([]BrokerageApiReturn, error) {

	var brokerageReturn []BrokerageApiReturn
	var err error

	queryDefault := "SELECT id, name, country FROM brokerage "

	if specificFetch == "ALL" {
		err = pgxscan.Select(context.Background(), &dbpool, &brokerageReturn,
			queryDefault)
		if err != nil {
			panic(err)
		}
	} else if specificFetch == "SINGLE" {
		query := queryDefault + "where name=$1"
		err = pgxscan.Select(context.Background(), &dbpool, &brokerageReturn,
			query, args[0])
		if err != nil {
			panic(err)
		}
	} else if specificFetch == "COUNTRY" {
		query := queryDefault + "where country=$1"
		err = pgxscan.Select(context.Background(), &dbpool, &brokerageReturn,
			query, args[0])
		if err != nil {
			panic(err)
		}
	}

	return brokerageReturn, err
}
