package postgresql

import (
	"context"
	"fmt"
	"stockfyApi/database"

	"github.com/georgysavva/scany/pgxscan"
)

func (r *repo) CreateUser(signUp database.Users) ([]database.Users, error) {

	var userRow []database.Users

	insertRow := `
	INSERT INTO
		users(username, email, uid, type)
	VALUES ($1, $2, $3, $4)
	RETURNING id, username, email, uid, type;
	`
	err := pgxscan.Select(context.Background(), r.dbpool, &userRow, insertRow,
		signUp.Username, signUp.Email, signUp.Uid, "normal")
	if err != nil {
		fmt.Println("database.CreateUser: ", err)
	}

	return userRow, err
}

func (r *repo) DeleteUser(firebaseUid string) ([]database.Users, error) {
	var userRow []database.Users

	deleteRow := `
	DELETE from users as u
	WHERE u.uid = $1
	RETURNING u.id, u.uid, u.username, u.email, u.type;
	`

	err := pgxscan.Select(context.Background(), r.dbpool, &userRow, deleteRow,
		firebaseUid)
	if err != nil {
		fmt.Println("database.DeleteUser: ", err)
	}

	return userRow, err
}

func (r *repo) UpdateUser(userInfo database.Users) ([]database.Users, error) {
	var userRow []database.Users

	query := `
	UPDATE users as u
	SET email = $2,
		username = $3
	WHERE u.uid = $1
	RETURNING u.id, u.uid, u.username, u.email, u.type;
	`

	err := pgxscan.Select(context.Background(), r.dbpool, &userRow, query,
		userInfo.Uid, userInfo.Email, userInfo.Username)
	if err != nil {
		fmt.Println("database.UpdateUser: ", err)
	}

	return userRow, err
}

func (r *repo) SearchUser(userUid string) ([]database.Users, error) {
	var userRow []database.Users

	query := `
	SELECT
		uid, email, username, "type"
	FROM users
	WHERE uid=$1;
	`
	err := pgxscan.Select(context.Background(), r.dbpool, &userRow, query, userUid)
	if err != nil {
		fmt.Println("database.UpdateUser: ", err)
	}

	return userRow, err
}
