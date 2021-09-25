package postgresql

import (
	"context"
	"fmt"
	"os"
	"stockfyApi/database"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(pgx.Tx) error) (err error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)
	Ping(context.Context) error
	Prepare(context.Context, string, string) (*pgconn.StatementDescription, error)
	// Deallocate(ctx context.Context, name string) error
	Close(context.Context) error
}

type repo struct {
	dbpool PgxIface
}

func NewPostgresqlRepository(dbname string, user string,
	password string) database.Repository {

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		user, password, dbname)

	DBpool, err := pgx.Connect(context.Background(), dbinfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer DBpool.Close(context.Background())

	return &repo{
		dbpool: DBpool,
	}
}
