package postgresql

import (
	"context"
	"fmt"
	"stockfyApi/entity"

	"github.com/georgysavva/scany/pgxscan"
	_ "github.com/lib/pq"
)

type AssetUserPostgres struct {
	dbpool PgxIface
}

func NewAssetUserPostgres(db PgxIface) *AssetUserPostgres {
	return &AssetUserPostgres{
		dbpool: db,
	}
}

func (r *AssetUserPostgres) Create(assetId string, userUid string) (
	[]entity.AssetUsers, error) {
	var assetUser []entity.AssetUsers

	insertRow := `
		INSERT INTO
			asset_users(asset_id, user_uid)
		VALUES ($1, $2)
		RETURNING asset_id, user_uid;
		`
	err := pgxscan.Select(context.Background(), r.dbpool, &assetUser,
		insertRow, assetId, userUid)
	if err != nil {
		fmt.Println("entity.CreateAssetUserRelation: ", err)
	}

	return assetUser, err
}

func (r *AssetUserPostgres) Delete(assetId string, userUid string) (
	[]entity.AssetUsers, error) {
	var assetUser []entity.AssetUsers

	deleteRow := `
	DELETE from asset_users as au
	WHERE au.asset_id = $1 and au.user_uid = $2
	RETURNING au.asset_id, au.user_uid;
	`
	err := pgxscan.Select(context.Background(), r.dbpool, &assetUser,
		deleteRow, assetId, userUid)
	if err != nil {
		fmt.Println("entity.DeleteAssetUserRelation: ", err)
	}

	return assetUser, err
}

func (r *AssetUserPostgres) DeleteByAsset(assetId string) (
	[]entity.AssetUsers, error) {
	var assetUser []entity.AssetUsers

	deleteRow := `
	DELETE from asset_users as au
	WHERE au.asset_id = $1
	RETURNING au.asset_id, au.user_uid;
	`
	err := pgxscan.Select(context.Background(), r.dbpool, &assetUser,
		deleteRow, assetId)
	if err != nil {
		fmt.Println("entity.DeleteAssetUserRelation: ", err)
	}

	return assetUser, err
}

func (r *AssetUserPostgres) DeleteByUser(userUid string) (
	[]entity.AssetUsers, error) {
	var assetUser []entity.AssetUsers

	deleteRow := `
	DELETE from asset_users as au
	WHERE au.user_uid = $1
	RETURNING au.asset_id, au.user_uid;
	`
	err := pgxscan.Select(context.Background(), r.dbpool, &assetUser,
		deleteRow, userUid)
	if err != nil {
		fmt.Println("entity.DeleteAssetUserRelation: ", err)
	}

	return assetUser, err
}

func (r *AssetUserPostgres) Search(assetId string, userUid string) (
	[]entity.AssetUsers, error) {
	var assetUser []entity.AssetUsers

	query := `
	SELECT
		asset_id, user_uid
	FROM asset_users
	WHERE asset_id=$1 and user_uid=$2;
	`
	err := pgxscan.Select(context.Background(), r.dbpool, &assetUser,
		query, assetId, userUid)
	if err != nil {
		fmt.Println("entity.SearchAssetUserRelation: ", err)
	}

	return assetUser, err
}
