package postgresql

import (
	"context"
	"regexp"
	"stockfyApi/entity"
	"testing"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

var userCreate = entity.Users{
	Uid:      "a48a93kdjfaj4a",
	Username: "Pedro Soares",
	Email:    "test@gmail.com",
	Type:     "normal",
}

var expectedSectorInfo = []entity.Users{
	{
		Uid:      "a48a93kdjfaj4a",
		Username: "Pedro Soares",
		Email:    "test@gmail.com",
		Type:     "normal",
	},
}

func userMockDatabase() (pgxmock.PgxConnIface, *pgxmock.Rows, error) {
	columns := []string{"uid", "username", "email", "type"}

	mock, err := pgxmock.NewConn()
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)

	return mock, rows, err
}

func TestUserCreate(t *testing.T) {

	query := regexp.QuoteMeta(`
	INSERT INTO
		users(username, email, uid, type)
	VALUES ($1, $2, $3, $4)
	RETURNING uid, username, email, type;
	`)

	mock, rows, err := userMockDatabase()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}

	mock.ExpectQuery(query).WithArgs(userCreate.Username, userCreate.Email,
		userCreate.Uid, userCreate.Type).WillReturnRows(rows.AddRow(
		"a48a93kdjfaj4a", "Pedro Soares", "test@gmail.com", "normal"))

	Users := UserPostgres{dbpool: mock}
	userRow, _ := Users.Create(userCreate)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, userRow)
	assert.Equal(t, expectedSectorInfo, userRow)
}

func TestUserDelete(t *testing.T) {

	query := regexp.QuoteMeta(`
	DELETE from users as u
	WHERE u.uid = $1
	RETURNING u.uid, u.username, u.email, u.type;
	`)

	mock, rows, err := userMockDatabase()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}

	mock.ExpectQuery(query).WithArgs(userCreate.Uid).WillReturnRows(rows.AddRow(
		"a48a93kdjfaj4a", "Pedro Soares", "test@gmail.com", "normal"))

	Users := UserPostgres{dbpool: mock}
	userRow, _ := Users.Delete(userCreate.Uid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, userRow)
	assert.Equal(t, expectedSectorInfo, userRow)
}

func TestUserUpdate(t *testing.T) {

	query := regexp.QuoteMeta(`
	UPDATE users as u
	SET email = $2,
		username = $3
	WHERE u.uid = $1
	RETURNING u.uid, u.username, u.email, u.type;
	`)

	mock, rows, err := userMockDatabase()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}

	mock.ExpectQuery(query).WithArgs(userCreate.Uid, userCreate.Email,
		userCreate.Username).WillReturnRows(rows.AddRow(
		"a48a93kdjfaj4a", "Pedro Soares", "test@gmail.com", "normal"))

	Users := UserPostgres{dbpool: mock}
	userRow, _ := Users.Update(userCreate)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, userRow)
	assert.Equal(t, expectedSectorInfo, userRow)
}

func TestUserSearch(t *testing.T) {

	query := regexp.QuoteMeta(`
	SELECT
		uid, email, username, "type"
	FROM users
	WHERE uid=$1;
	`)

	mock, rows, err := userMockDatabase()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}

	mock.ExpectQuery(query).WithArgs(userCreate.Uid).WillReturnRows(rows.AddRow(
		"a48a93kdjfaj4a", "Pedro Soares", "test@gmail.com", "normal"))

	Users := UserPostgres{dbpool: mock}
	userRow, _ := Users.Search(userCreate.Uid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, userRow)
	assert.Equal(t, expectedSectorInfo, userRow)
}
