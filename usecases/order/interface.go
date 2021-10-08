package order

import "stockfyApi/entity"

type Repository interface {
	Create(orderInsert entity.Order) entity.Order
	DeleteFromUser(id string, userUid string) (string, error)
	DeleteFromAsset(symbolId string) ([]entity.Order, error)
	SearchByOrderAndUserId(orderId string, userUid string) ([]entity.Order,
		error)
	DeleteFromAssetUser(assetId string, userUid string) ([]entity.Order, error)
	SearchFromAssetUser(assetId string, userUid string) ([]entity.Order, error)
	UpdateFromUser(orderUpdate entity.Order) []entity.Order
}
