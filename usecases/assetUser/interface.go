package assetusers

import "stockfyApi/entity"

type Repository interface {
	Create(assetId string, userUid string) ([]entity.AssetUsers, error)
	// DeleteByAsset(assetId string) ([]entity.AssetUsers, error)
	Search(assetId string, userUid string) ([]entity.AssetUsers, error)
}
