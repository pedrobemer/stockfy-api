package postgresql

import (
	"context"
	"fmt"
	"stockfyApi/database"

	"github.com/georgysavva/scany/pgxscan"
	_ "github.com/lib/pq"
)

func (r *repo) CreateAssetUserRelation(assetId string, userUid string) (
	[]database.AssetUsers, error) {
	var assetUser []database.AssetUsers

	insertRow := `
		INSERT INTO
			asset_users(asset_id, user_uid)
		VALUES ($1, $2)
		RETURNING asset_id, user_uid;
		`
	err := pgxscan.Select(context.Background(), r.dbpool, &assetUser,
		insertRow, assetId, userUid)
	if err != nil {
		fmt.Println("database.CreateAssetUserRelation: ", err)
	}

	return assetUser, err
}

func (r *repo) DeleteAssetUserRelation(assetId string, userUid string) (
	[]database.AssetUsers, error) {
	var assetUser []database.AssetUsers

	deleteRow := `
	DELETE from asset_users as au
	WHERE au.asset_id = $1 and au.user_uid = $2
	RETURNING au.asset_id, au.user_uid;
	`
	err := pgxscan.Select(context.Background(), r.dbpool, &assetUser,
		deleteRow, assetId, userUid)
	if err != nil {
		fmt.Println("database.DeleteAssetUserRelation: ", err)
	}

	return assetUser, err
}

func (r *repo) DeleteAssetUserRelationByAsset(assetId string) (
	[]database.AssetUsers, error) {
	var assetUser []database.AssetUsers

	deleteRow := `
	DELETE from asset_users as au
	WHERE au.asset_id = $1
	RETURNING au.asset_id, au.user_uid;
	`
	err := pgxscan.Select(context.Background(), r.dbpool, &assetUser,
		deleteRow, assetId)
	if err != nil {
		fmt.Println("database.DeleteAssetUserRelation: ", err)
	}

	return assetUser, err
}

func (r *repo) DeleteAssetUserRelationByUser(userUid string) (
	[]database.AssetUsers, error) {
	var assetUser []database.AssetUsers

	deleteRow := `
	DELETE from asset_users as au
	WHERE au.user_uid = $1
	RETURNING au.asset_id, au.user_uid;
	`
	err := pgxscan.Select(context.Background(), r.dbpool, &assetUser,
		deleteRow, userUid)
	if err != nil {
		fmt.Println("database.DeleteAssetUserRelation: ", err)
	}

	return assetUser, err
}

func (r *repo) SearchAssetUserRelation(assetId string, userUid string) (
	[]database.AssetUsers, error) {
	var assetUser []database.AssetUsers

	query := `
	SELECT
		asset_id, user_uid
	FROM asset_users
	WHERE asset_id=$1 and user_uid=$2;
	`
	err := pgxscan.Select(context.Background(), r.dbpool, &assetUser,
		query, assetId, userUid)
	if err != nil {
		fmt.Println("database.SearchAssetUserRelation: ", err)
	}

	return assetUser, err
}
