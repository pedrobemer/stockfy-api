package assetusers

import "stockfyApi/entity"

type Repository interface {
	Create(assetId string, userUid string) ([]entity.AssetUsers, error)
	DeleteByAsset(assetId string) ([]entity.AssetUsers, error)
	Delete(assetId string, userUid string) ([]entity.AssetUsers, error)
	Search(assetId string, userUid string) ([]entity.AssetUsers, error)
}

type UseCases interface {
	CreateAssetUserRelation(assetId string, userUid string) (*entity.AssetUsers,
		error)
	DeleteAssetUserRelation(assetId string, userUid string) (*entity.AssetUsers,
		error)
	DeleteAssetUserRelationByAsset(assetId string) ([]entity.AssetUsers, error)
	SearchAssetUserRelation(assetId string, userUid string) (*entity.AssetUsers,
		error)
}
