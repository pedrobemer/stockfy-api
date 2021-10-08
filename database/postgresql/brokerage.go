package postgresql

import (
	"context"
	"stockfyApi/entity"

	"github.com/georgysavva/scany/pgxscan"
	_ "github.com/lib/pq"
)

type BrokeragePostgres struct {
	dbpool PgxIface
}

func NewBrokeragePostgres(db PgxIface) *BrokeragePostgres {
	return &BrokeragePostgres{
		dbpool: db,
	}
}

func (r *BrokeragePostgres) Search(specificFetch string, args ...string) (
	[]entity.Brokerage, error) {

	var brokerageReturn []entity.Brokerage
	var err error

	if specificFetch != "ALL" && specificFetch != "SINGLE" &&
		specificFetch != "COUNTRY" {
		return brokerageReturn, entity.ErrInvalidBrokerageSearchType
	}

	queryDefault := "SELECT id, name, country FROM brokerage "

	if specificFetch == "ALL" {
		err = pgxscan.Select(context.Background(), r.dbpool, &brokerageReturn,
			queryDefault)
		if err != nil {
			return nil, err
		}
	} else if specificFetch == "SINGLE" {
		query := queryDefault + "where name=$1"
		err = pgxscan.Select(context.Background(), r.dbpool, &brokerageReturn,
			query, args[0])
		if err != nil {
			return nil, err
		}
	} else if specificFetch == "COUNTRY" {
		query := queryDefault + "where country=$1"
		err = pgxscan.Select(context.Background(), r.dbpool, &brokerageReturn,
			query, args[0])
		if err != nil {
			return nil, err
		}
	}

	return brokerageReturn, err
}
