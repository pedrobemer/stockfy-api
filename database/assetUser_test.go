package database

import (
	"context"
	"regexp"
	"testing"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

var assetUserInput = AssetUsersApiReturn{
	AssetId: "a48a93kdjfaj4a",
	UserUid: "b948aliru78",
}

var expectedAssetUser = []AssetUsersApiReturn{
	{
		AssetId: "a48a93kdjfaj4a",
		UserUid: "b948aliru78",
	},
}

func assetUserMockDatabase(query string) (pgxmock.PgxConnIface, error) {
	columns := []string{"asset_id", "user_uid"}

	mock, err := pgxmock.NewConn()
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	mock.ExpectQuery(query).WithArgs(assetUserInput.AssetId,
		assetUserInput.UserUid).WillReturnRows(rows.AddRow("a48a93kdjfaj4a",
		"b948aliru78"))

	return mock, err
}

func TestCreateAssetUserRelation(t *testing.T) {

	query := regexp.QuoteMeta(`
	INSERT INTO
		asset_users(asset_id, user_uid)
	VALUES ($1, $2)
	RETURNING asset_id, user_uid;
	`)

	mock, err := assetUserMockDatabase(query)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	assetUserRow, _ := CreateAssetUserRelation(mock, assetUserInput.AssetId,
		assetUserInput.UserUid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, assetUserRow)
	assert.Equal(t, expectedAssetUser, assetUserRow)
}

func TestDeleteAssetUserRelation(t *testing.T) {

	query := regexp.QuoteMeta(`
	DELETE from asset_users as au
	WHERE au.asset_id = $1 and au.user_uid = $2
	RETURNING au.asset_id, au.user_uid;
	`)

	mock, err := assetUserMockDatabase(query)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	assetUserRow, _ := DeleteAssetUserRelation(mock, assetUserInput.AssetId,
		assetUserInput.UserUid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, assetUserRow)
	assert.Equal(t, expectedAssetUser, assetUserRow)
}

func TestSearchAssetUserRelation(t *testing.T) {

	query := regexp.QuoteMeta(`
	SELECT
		asset_id, user_uid
	FROM asset_users
	WHERE asset_id=$1 and user_uid=$2;
	`)

	mock, err := assetUserMockDatabase(query)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	assetUserRow, _ := SearchAssetUserRelation(mock, assetUserInput.AssetId,
		assetUserInput.UserUid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, assetUserRow)
	assert.Equal(t, expectedAssetUser, assetUserRow)
}
