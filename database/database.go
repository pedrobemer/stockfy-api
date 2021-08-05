package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

var DBpool *pgx.Conn

const (
	DB_USER     = "pedrobemer"
	DB_PASSWORD = "pirulito"
	DB_NAME     = "stockfy"
)

func Connect() error {
	var err error
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)

	DBpool, err = pgx.Connect(context.Background(), dbinfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer DBpool.Close(context.Background())

	return err
}
