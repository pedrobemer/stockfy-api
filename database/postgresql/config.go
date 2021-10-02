package postgresql

import (
	"context"
	"stockfyApi/usecases"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions,
		f func(pgx.Tx) error) (err error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	QueryFunc(ctx context.Context, sql string, args []interface{},
		scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag,
		error)
	Ping(context.Context) error
	Prepare(context.Context, string, string) (*pgconn.StatementDescription,
		error)
	// Deallocate(ctx context.Context, name string) error
	Close(context.Context) error
}

func NewPostgresInstance(dbpool PgxIface) usecases.Repositories {
	return usecases.Repositories{
		AssetRepository:          NewAssetPostgres(dbpool),
		SectorRepository:         NewSectorPostgres(dbpool),
		AssetTypeRepository:      NewAssetTypePostgres(dbpool),
		UserRepository:           NewUserPostgres(dbpool),
		DbVerificationRepository: NewDbVerificationPostgres(dbpool),
	}
}
