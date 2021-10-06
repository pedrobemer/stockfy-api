package order

import "stockfyApi/entity"

type Repository interface {
	Create(orderInsert entity.Order) entity.Order
	DeleteFromAsset(symbolId string) ([]entity.Order, error)
	DeleteFromAssetUser(assetId string, userUid string) ([]entity.Order, error)
}
