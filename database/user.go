package database

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
)

func CreateUser(dbpool PgxIface, signUp UserDatabase) ([]UserDatabase, error) {

	var userRow []UserDatabase

	insertRow := `
	INSERT INTO
		users(username, email, uid, type)
	VALUES ($1, $2, $3, $4)
	RETURNING id, username, email, uid, type;
	`
	err := pgxscan.Select(context.Background(), dbpool, &userRow, insertRow,
		signUp.Username, signUp.Email, signUp.Uid, "normal")
	if err != nil {
		fmt.Println("database.CreateUser: ", err)
	}

	return userRow, err
}

func DeleteUser(dbpool PgxIface, firebaseUid string) ([]UserDatabase, error) {
	var userRow []UserDatabase

	deleteRow := `
	DELETE from users as u
	WHERE u.uid = $1
	RETURNING u.id, u.uid, u.username, u.email, u.type;
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
	UPDATE users as u
	SET email = $2,
		username = $3
	WHERE u.uid = $1
	RETURNING u.id, u.uid, u.username, u.email, u.type;
	`

	err := pgxscan.Select(context.Background(), dbpool, &userRow, query,
		userInfo.Uid, userInfo.Email, userInfo.Username)
	if err != nil {
		fmt.Println("database.UpdateUser: ", err)
	}

	return userRow, err
}

func SearchUser(dbpool PgxIface, userUid string) ([]UserDatabase, error) {
	var userRow []UserDatabase

	query := `
	SELECT
		uid, email, username, "type"
	FROM users
	WHERE uid=$1;
	`
	err := pgxscan.Select(context.Background(), dbpool, &userRow, query, userUid)
	if err != nil {
		fmt.Println("database.UpdateUser: ", err)
	}

	return userRow, err
}
