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

type UseCases interface {
	CreateOrder(quantity float64, price float64, currency string,
		orderType string, date string, brokerageId string, assetId string,
		userUid string) (*entity.Order, error)
	DeleteOrdersFromAsset(assetId string) ([]entity.Order, error)
	DeleteOrdersFromAssetUser(assetId string, userUid string) (*[]entity.Order,
		error)
	DeleteOrdersFromUser(orderId string, userUid string) (*string, error)
	SearchOrderByIdAndUserUid(orderId string, userUid string) (*entity.Order,
		error)
	SearchOrdersFromAssetUser(assetId string, userUid string) ([]entity.Order,
		error)
	UpdateOrder(orderId string, userUid string, price float64, quantity float64,
		orderType, date string, brokerageId string, currency string) (
		*entity.Order, error)
	OrderVerification(orderType string, country string, quantity float64,
		price float64, currency string) error
}
