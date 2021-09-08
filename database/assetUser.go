package database

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
	_ "github.com/lib/pq"
)

func CreateAssetUserRelation(dbpool PgxIface, assetId string, userUid string) (
	[]AssetUsersApiReturn, error) {
	var assetUser []AssetUsersApiReturn

	insertRow := `
		INSERT INTO
			asset_users(asset_id, user_uid)
		VALUES ($1, $2)
		RETURNING asset_id, user_uid;
		`
	err := pgxscan.Select(context.Background(), dbpool, &assetUser,
		insertRow, assetId, userUid)
	if err != nil {
		fmt.Println("database.CreateAssetUserRelation: ", err)
	}

	return assetUser, err
}

func DeleteAssetUserRelation(dbpool PgxIface, assetId string, userUid string) (
	[]AssetUsersApiReturn, error) {
	var assetUser []AssetUsersApiReturn

	deleteRow := `
	DELETE from asset_users as au
	WHERE au.asset_id = $1 and au.user_uid = $2
	RETURNING au.asset_id, au.user_uid;
	`
	err := pgxscan.Select(context.Background(), dbpool, &assetUser,
		deleteRow, assetId, userUid)
	if err != nil {
		fmt.Println("database.DeleteAssetUserRelation: ", err)
	}

	return assetUser, err
}

func DeleteAssetUserRelationByAsset(dbpool PgxIface, assetId string) (
	[]AssetUsersApiReturn, error) {
	var assetUser []AssetUsersApiReturn

	deleteRow := `
	DELETE from asset_users as au
	WHERE au.asset_id = $1
	RETURNING au.asset_id, au.user_uid;
	`
	err := pgxscan.Select(context.Background(), dbpool, &assetUser,
		deleteRow, assetId)
	if err != nil {
		fmt.Println("database.DeleteAssetUserRelation: ", err)
	}

	return assetUser, err
}

func DeleteAssetUserRelationByUser(dbpool PgxIface, userUid string) (
	[]AssetUsersApiReturn, error) {
	var assetUser []AssetUsersApiReturn

	deleteRow := `
	DELETE from asset_users as au
	WHERE au.user_uid = $1
	RETURNING au.asset_id, au.user_uid;
	`
	err := pgxscan.Select(context.Background(), dbpool, &assetUser,
		deleteRow, userUid)
	if err != nil {
		fmt.Println("database.DeleteAssetUserRelation: ", err)
	}

	return assetUser, err
}

func SearchAssetUserRelation(dbpool PgxIface, assetId string, userUid string) (
	[]AssetUsersApiReturn, error) {
	var assetUser []AssetUsersApiReturn

	query := `
	SELECT
		asset_id, user_uid
	FROM asset_users
	WHERE asset_id=$1 and user_uid=$2;
	`
	err := pgxscan.Select(context.Background(), dbpool, &assetUser,
		query, assetId, userUid)
	if err != nil {
		fmt.Println("database.SearchAssetUserRelation: ", err)
	}

	return assetUser, err
}
