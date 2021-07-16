package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

func verifyRowExistence(dbpool pgxpool.Pool, table string, condition string) bool {
	var rowExist bool

	var fetchRow = "SELECT exists(SELECT 1 FROM " + table + " where " +
		condition + ");"

	err := dbpool.QueryRow(context.Background(), fetchRow).Scan(&rowExist)
	if err != nil {
		fmt.Println(err)
	}

	return rowExist
}
