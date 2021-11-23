package postgresql

import (
	"context"
	"errors"
	"fmt"
	"stockfyApi/entity"

	"github.com/georgysavva/scany/pgxscan"
)

type UserPostgres struct {
	dbpool PgxIface
}

func NewUserPostgres(db PgxIface) *UserPostgres {
	return &UserPostgres{
		dbpool: db,
	}
}

func (r *UserPostgres) Create(signUp entity.Users) ([]entity.Users, error) {

	var userRow []entity.Users

	insertRow := `
	INSERT INTO
		users(username, email, uid, type)
	VALUES ($1, $2, $3, $4)
	RETURNING uid, username, email, type;
	`
	err := pgxscan.Select(context.Background(), r.dbpool, &userRow, insertRow,
		signUp.Username, signUp.Email, signUp.Uid, "normal")
	if err != nil {
		// fmt.Println("entity.CreateUser: ", err)
		return nil, errors.New("entity.CreateUser: " + err.Error())
	}

	return userRow, err
}

func (r *UserPostgres) Delete(firebaseUid string) ([]entity.Users, error) {
	var userRow []entity.Users

	deleteRow := `
	DELETE from users as u
	WHERE u.uid = $1
	RETURNING u.uid, u.username, u.email, u.type;
	`

	err := pgxscan.Select(context.Background(), r.dbpool, &userRow, deleteRow,
		firebaseUid)
	if err != nil {
		// fmt.Println("entity.DeleteUser: ", err)
		return nil, errors.New("entity.DeleteUser: " + err.Error())
	}

	return userRow, nil
}

func (r *UserPostgres) Update(userInfo entity.Users) ([]entity.Users, error) {
	var userRow []entity.Users

	query := `
	UPDATE users as u
	SET email = $2,
		username = $3
	WHERE u.uid = $1
	RETURNING u.uid, u.username, u.email, u.type;
	`

	err := pgxscan.Select(context.Background(), r.dbpool, &userRow, query,
		userInfo.Uid, userInfo.Email, userInfo.Username)
	if err != nil {
		fmt.Println("entity.UpdateUser: ", err)
	}

	return userRow, err
}

func (r *UserPostgres) Search(userUid string) ([]entity.Users, error) {
	var userRow []entity.Users

	query := `
	SELECT
		uid, email, username, "type"
	FROM users
	WHERE uid=$1;
	`
	err := pgxscan.Select(context.Background(), r.dbpool, &userRow, query, userUid)
	if err != nil {
		fmt.Println("entity.UpdateUser: ", err)
	}

	if userRow == nil {
		return nil, entity.ErrInvalidUserSearch
	}

	return userRow, err
}
