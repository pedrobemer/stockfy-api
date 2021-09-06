package database

import (
	"context"
	"regexp"
	"testing"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

var userCreate = UserDatabase{
	Uid:      "a48a93kdjfaj4a",
	Username: "Pedro Soares",
	Email:    "test@gmail.com",
}

var expectedSectorInfo = []UserDatabase{
	{
		Id:       "0a52d206-ed8b-11eb-9a03-0242ac130003",
		Uid:      "a48a93kdjfaj4a",
		Username: "Pedro Soares",
		Email:    "test@gmail.com",
	},
}

func userMockDatabase() (pgxmock.PgxConnIface, *pgxmock.Rows, error) {
	columns := []string{"id", "uid", "username", "email"}

	mock, err := pgxmock.NewConn()
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)

	return mock, rows, err
}

func TestCreateUser(t *testing.T) {

	query := regexp.QuoteMeta(`
	insert into
		users(username, email, uid)
	values ($1, $2, $3)
	returning id, username, email, uid;
	`)

	mock, rows, err := userMockDatabase()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(query).WithArgs(userCreate.Username, userCreate.Email,
		userCreate.Uid).WillReturnRows(rows.AddRow(
		"0a52d206-ed8b-11eb-9a03-0242ac130003", "a48a93kdjfaj4a", "Pedro Soares",
		"test@gmail.com"))

	userRow, _ := CreateUser(mock, userCreate)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, userRow)
	assert.Equal(t, expectedSectorInfo, userRow)
}

func TestDeleteUser(t *testing.T) {

	query := regexp.QuoteMeta(`
	delete from users as u
	where u.uid = $1
	returning u.id, u.uid, u.username, u.email;
	`)

	mock, rows, err := userMockDatabase()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(query).WithArgs(userCreate.Uid).WillReturnRows(rows.AddRow(
		"0a52d206-ed8b-11eb-9a03-0242ac130003", "a48a93kdjfaj4a", "Pedro Soares",
		"test@gmail.com"))

	userRow, _ := DeleteUser(mock, userCreate.Uid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, userRow)
	assert.Equal(t, expectedSectorInfo, userRow)
}

func TestUpdateUser(t *testing.T) {

	query := regexp.QuoteMeta(`
	update users as u
	set email = $2,
		username = $3
	where u.uid = $1
	returning u.id, u.uid, u.username, u.email;
	`)

	mock, rows, err := userMockDatabase()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(query).WithArgs(userCreate.Uid, userCreate.Email,
		userCreate.Username).WillReturnRows(rows.AddRow(
		"0a52d206-ed8b-11eb-9a03-0242ac130003", "a48a93kdjfaj4a", "Pedro Soares",
		"test@gmail.com"))

	userRow, _ := UpdateUser(mock, userCreate)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, userRow)
	assert.Equal(t, expectedSectorInfo, userRow)
}
