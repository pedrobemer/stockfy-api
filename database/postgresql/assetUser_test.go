package postgresql

import (
	"context"
	"regexp"
	"stockfyApi/entity"
	"testing"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

var assetUserInput = entity.AssetUsers{
	AssetId: "a48a93kdjfaj4a",
	UserUid: "b948aliru78",
}

var expectedAssetUser = []entity.AssetUsers{
	{
		AssetId: "a48a93kdjfaj4a",
		UserUid: "b948aliru78",
	},
}

func assetUserMockDatabase(query string, arguments int, args ...string) (
	pgxmock.PgxConnIface, error) {
	columns := []string{"asset_id", "user_uid"}

	mock, err := pgxmock.NewConn()
	defer mock.Close(context.Background())

	rows := mock.NewRows(columns)
	if arguments == 2 {
		mock.ExpectQuery(query).WithArgs(assetUserInput.AssetId,
			assetUserInput.UserUid).WillReturnRows(rows.AddRow("a48a93kdjfaj4a",
			"b948aliru78"))
	} else {
		mock.ExpectQuery(query).WithArgs(args[0]).WillReturnRows(rows.AddRow(
			"a48a93kdjfaj4a", "b948aliru78"))
	}

	return mock, err
}

func TestAssetUserCreate(t *testing.T) {

	query := regexp.QuoteMeta(`
	INSERT INTO
		asset_users(asset_id, user_uid)
	VALUES ($1, $2)
	RETURNING asset_id, user_uid;
	`)

	mock, err := assetUserMockDatabase(query, 2)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}

	Au := AssetUserPostgres{dbpool: mock}
	assetUserRow, _ := Au.Create(assetUserInput.AssetId,
		assetUserInput.UserUid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, assetUserRow)
	assert.Equal(t, expectedAssetUser, assetUserRow)
}

func TestAssetUserDelete(t *testing.T) {

	query := regexp.QuoteMeta(`
	DELETE from asset_users as au
	WHERE au.asset_id = $1 and au.user_uid = $2
	RETURNING au.asset_id, au.user_uid;
	`)

	mock, err := assetUserMockDatabase(query, 2)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}

	Au := AssetUserPostgres{dbpool: mock}
	assetUserRow, _ := Au.Delete(assetUserInput.AssetId,
		assetUserInput.UserUid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, assetUserRow)
	assert.Equal(t, expectedAssetUser, assetUserRow)
}

func TestAssetUserDeleteByAsset(t *testing.T) {

	query := regexp.QuoteMeta(`
	DELETE from asset_users as au
	WHERE au.asset_id = $1
	RETURNING au.asset_id, au.user_uid;
	`)

	mock, err := assetUserMockDatabase(query, 1, assetUserInput.AssetId)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}

	Au := AssetUserPostgres{dbpool: mock}
	assetUserRow, _ := Au.DeleteByAsset(assetUserInput.AssetId)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, assetUserRow)
	assert.Equal(t, expectedAssetUser, assetUserRow)
}

func TestAssetUserDeleteByUser(t *testing.T) {

	query := regexp.QuoteMeta(`
	DELETE from asset_users as au
	WHERE au.user_uid = $1
	RETURNING au.asset_id, au.user_uid;
	`)

	mock, err := assetUserMockDatabase(query, 1, assetUserInput.UserUid)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}

	Au := AssetUserPostgres{dbpool: mock}
	assetUserRow, _ := Au.DeleteByUser(assetUserInput.UserUid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, assetUserRow)
	assert.Equal(t, expectedAssetUser, assetUserRow)
}

func TestAssetUserSearch(t *testing.T) {

	query := regexp.QuoteMeta(`
	SELECT
		asset_id, user_uid
	FROM asset_users
	WHERE asset_id=$1 and user_uid=$2;
	`)

	mock, err := assetUserMockDatabase(query, 2)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub entity connection", err)
	}

	Au := AssetUserPostgres{dbpool: mock}

	assetUserRow, _ := Au.Search(assetUserInput.AssetId,
		assetUserInput.UserUid)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, assetUserRow)
	assert.Equal(t, expectedAssetUser, assetUserRow)
}
