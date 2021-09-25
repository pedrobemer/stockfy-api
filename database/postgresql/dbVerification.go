package postgresql

import (
	"context"
	"fmt"

	_ "github.com/lib/pq"
)

func (r *repo) VerifyRowExistence(table string, condition string) bool {
	var rowExist bool

	var fetchRow = "SELECT exists(SELECT 1 FROM " + table + " where " +
		condition + ");"

	err := r.dbpool.QueryRow(context.Background(), fetchRow).Scan(&rowExist)
	if err != nil {
		fmt.Println(err)
	}

	return rowExist
}
