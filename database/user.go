package database

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
)

func CreateUser(dbpool PgxIface, signUp UserDatabase) ([]UserDatabase, error) {

	var userRow []UserDatabase

	insertRow := `
	insert into
		users(username, email, uid)
	values ($1, $2, $3)
	returning id, username, email, uid;
	`
	err := pgxscan.Select(context.Background(), dbpool, &userRow, insertRow,
		signUp.Username, signUp.Email, signUp.Uid)
	if err != nil {
		fmt.Println("database.CreateUser: ", err)
	}

	return userRow, err
}

func DeleteUser(dbpool PgxIface, firebaseUid string) ([]UserDatabase, error) {
	var userRow []UserDatabase

	deleteRow := `
	delete from users as u
	where u.uid = $1
	returning u.id, u.uid, u.username, u.email;
	`

	err := pgxscan.Select(context.Background(), dbpool, &userRow, deleteRow,
		firebaseUid)
	if err != nil {
		fmt.Println("database.DeleteUser: ", err)
	}

	return userRow, err
}

func UpdateUser(dbpool PgxIface, userInfo UserDatabase) ([]UserDatabase, error) {
	var userRow []UserDatabase

	query := `
	update users as u
	set email = $2,
		username = $3
	where u.uid = $1
	returning u.id, u.uid, u.username, u.email;
	`

	err := pgxscan.Select(context.Background(), dbpool, &userRow, query,
		userInfo.Uid, userInfo.Email, userInfo.Username)
	if err != nil {
		fmt.Println("database.DeleteOrders: ", err)
	}

	return userRow, err
}
