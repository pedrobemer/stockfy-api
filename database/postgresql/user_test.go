package postgresql

import (
	"context"
	"regexp"
	"stockfyApi/database"
	"testing"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

var userCreate = database.Users{
	Uid:      "a48a93kdjfaj4a",
	Username: "Pedro Soares",
	Email:    "test@gmail.com",
	Type:     "normal",
}

var expectedSectorInfo = []database.Users{
	{
		Id:       "0a52d206-ed8b-11eb-9a03-0242ac130003",
		Uid:      "a48a93kdjfaj4a",
		Username: "Pedro Soares",
		Email:    "test@gmail.com",
		Type:     "normal",
	},
}

func userMockDatabase() (pgxmock.PgxConnIface, *pgxmock.Rows, error) {
	columns := []string{"id", "uid", "username", "email", "type"}

	mock, err := pgxmock.NewConn()
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)

	return mock, rows, err
}

func TestCreateUser(t *testing.T) {

	query := regexp.QuoteMeta(`
	INSERT INTO
		users(username, email, uid, type)
	VALUES ($1, $2, $3, $4)
	RETURNING id, username, email, uid, type;
	`)

	mock, rows, err := userMockDatabase()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(query).WithArgs(userCreate.Username, userCreate.Email,
		userCreate.Uid, userCreate.Type).WillReturnRows(rows.AddRow(
		"0a52d206-ed8b-11eb-9a03-0242ac130003", "a48a93kdjfaj4a", "Pedro Soares",
		"test@gmail.com", "normal"))

	Users := repo{dbpool: mock}
	userRow, _ := Users.CreateUser(userCreate)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, userRow)
	assert.Equal(t, expectedSectorInfo, userRow)
}

func TestDeleteUser(t *testing.T) {

	query := regexp.QuoteMeta(`
	DELETE from users as u
	WHERE u.uid = $1
	RETURNING u.id, u.uid, u.username, u.email, u.type;
	`)

	mock, rows, err := userMockDatabase()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(query).WithArgs(userCreate.Uid).WillReturnRows(rows.AddRow(
		"0a52d206-ed8b-11eb-9a03-0242ac130003", "a48a93kdjfaj4a", "Pedro Soares",
		"test@gmail.com", "normal"))

	Users := repo{dbpool: mock}
	userRow, _ := Users.DeleteUser(userCreate.Uid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, userRow)
	assert.Equal(t, expectedSectorInfo, userRow)
}

func TestUpdateUser(t *testing.T) {

	query := regexp.QuoteMeta(`
	UPDATE users as u
	SET email = $2,
		username = $3
	WHERE u.uid = $1
	RETURNING u.id, u.uid, u.username, u.email, u.type;
	`)

	mock, rows, err := userMockDatabase()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(query).WithArgs(userCreate.Uid, userCreate.Email,
		userCreate.Username).WillReturnRows(rows.AddRow(
		"0a52d206-ed8b-11eb-9a03-0242ac130003", "a48a93kdjfaj4a", "Pedro Soares",
		"test@gmail.com", "normal"))

	Users := repo{dbpool: mock}
	userRow, _ := Users.UpdateUser(userCreate)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, userRow)
	assert.Equal(t, expectedSectorInfo, userRow)
}

func TestSearchUser(t *testing.T) {

	query := regexp.QuoteMeta(`
	SELECT
		uid, email, username, "type"
	FROM users
	WHERE uid=$1;
	`)

	mock, rows, err := userMockDatabase()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(query).WithArgs(userCreate.Uid).WillReturnRows(rows.AddRow(
		"0a52d206-ed8b-11eb-9a03-0242ac130003", "a48a93kdjfaj4a", "Pedro Soares",
		"test@gmail.com", "normal"))

	Users := repo{dbpool: mock}
	userRow, _ := Users.SearchUser(userCreate.Uid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, userRow)
	assert.Equal(t, expectedSectorInfo, userRow)
}
