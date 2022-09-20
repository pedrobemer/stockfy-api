package postgresql

import (
	"context"
	"stockfyApi/usecases"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type PgxIface interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row
}

func NewPostgresInstance(dbpool PgxIface) usecases.Repositories {
	return usecases.Repositories{
		AssetRepository:          NewAssetPostgres(dbpool),
		SectorRepository:         NewSectorPostgres(dbpool),
		AssetTypeRepository:      NewAssetTypePostgres(dbpool),
		UserRepository:           NewUserPostgres(dbpool),
		OrderRepository:          NewOrderPostgres(dbpool),
		AssetUserRepository:      NewAssetUserPostgres(dbpool),
		BrokerageRepository:      NewBrokeragePostgres(dbpool),
		EarningsRepository:       NewEarningPostgres(dbpool),
		DbVerificationRepository: NewDbVerificationPostgres(dbpool),
	}
}
