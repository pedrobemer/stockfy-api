package earnings

import "stockfyApi/entity"

type Repository interface {
	Create(earningOrder entity.Earnings) ([]entity.Earnings, error)
	DeleteFromAsset(assetId string) ([]entity.Earnings, error)
	SearchFromAssetUser(assetId string, userUid string) ([]entity.Earnings, error)
	DeleteFromUser(id string, userUid string) (string, error)
	DeleteFromAssetUser(assetId string, userUid string) ([]entity.Earnings, error)
}
