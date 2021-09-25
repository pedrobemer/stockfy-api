package postgresql

import (
	"context"
	"stockfyApi/database"

	"github.com/georgysavva/scany/pgxscan"
	_ "github.com/lib/pq"
)

func (r *repo) FetchBrokerage(specificFetch string, args ...string) (
	[]database.Brokerage, error) {

	var brokerageReturn []database.Brokerage
	var err error

	queryDefault := "SELECT id, name, country FROM brokerage "

	if specificFetch == "ALL" {
		err = pgxscan.Select(context.Background(), r.dbpool, &brokerageReturn,
			queryDefault)
		if err != nil {
			panic(err)
		}
	} else if specificFetch == "SINGLE" {
		query := queryDefault + "where name=$1"
		err = pgxscan.Select(context.Background(), r.dbpool, &brokerageReturn,
			query, args[0])
		if err != nil {
			panic(err)
		}
	} else if specificFetch == "COUNTRY" {
		query := queryDefault + "where country=$1"
		err = pgxscan.Select(context.Background(), r.dbpool, &brokerageReturn,
			query, args[0])
		if err != nil {
			panic(err)
		}
	}

	return brokerageReturn, err
}
