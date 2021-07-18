package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

var DBpool *pgxpool.Pool

const (
	DB_USER     = "pedrobemer"
	DB_PASSWORD = "pirulito"
	DB_NAME     = "stockfy"
)

func Connect() error {
	var err error
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)

	DBpool, err = pgxpool.Connect(context.Background(), dbinfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer DBpool.Close()

	return err
}
