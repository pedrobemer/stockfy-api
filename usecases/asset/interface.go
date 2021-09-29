package asset

import (
	"stockfyApi/entity"
	"stockfyApi/externalApi/alphaVantage"
	"stockfyApi/externalApi/finnhub"
	assettype "stockfyApi/usecases/assetType"
)

// type SectorRepository interface {
// 	CreateSector(sector string) ([]entity.Sector, error)
// 	FetchSectorByName(sector string) ([]entity.Sector, error)
// 	FetchSectorByAsset(symbol string) ([]entity.Sector, error)
// }

type Repository interface {
	Create(assetInsert entity.Asset) entity.Asset
	// Search(symbol string) ([]entity.Asset, error)
	// SearchByUser(symbol string, userUid string, orderType string) (
	// 	[]entity.Asset, error)
	// SearchPerAssetType(assetType string, country string, userUid string,
	// 	withOrdersInfo bool) []entity.AssetType
	// SearchByOrderId(orderId string) []entity.Asset
	// Delete(assetId string) []entity.Asset
}

// type BrokerageRepository interface {
// 	FetchBrokerage(specificFetch string, args ...string) ([]entity.Brokerage,
// 		error)
// }

// type AssetTypeRepository interface {
// 	FetchAssetType(fetchType string, args ...string) (entity.AssetType, error)
// }

// type OrderRepository interface {
// 	CreateOrder(orderInsert entity.Order) entity.Order
// 	SearchOrdersFromAssetUser(assetId string, userUid string) ([]entity.Order,
// 		error)
// 	DeleteOrderFromUser(id string, userUid string) string
// 	DeleteOrdersFromAsset(symbolId string) []entity.Order
// 	DeleteOrdersFromAssetUser(assetId string, userUid string) ([]entity.Order,
// 		error)
// 	UpdateOrderFromUser(orderUpdate entity.Order) []entity.Order
// }

// type UsersRepository interface {
// 	CreateUser(signUp entity.Users) ([]entity.Users, error)
// 	DeleteUser(firebaseUid string) ([]entity.Users, error)
// 	UpdateUser(userInfo entity.Users) ([]entity.Users, error)
// 	SearchUser(userUid string) ([]entity.Users, error)
// }

// type EarningsRepository interface {
// 	CreateEarningRow(earningOrder entity.Earnings) []entity.Earnings
// 	SearchEarningFromAssetUser(assetId string, userUid string) (
// 		[]entity.Earnings, error)
// 	DeleteEarningsFromAssetUser(assetId string, userUid string) (
// 		[]entity.Earnings, error)
// 	DeleteEarningFromUser(id string, userUid string) string
// 	UpdateEarningsFromUser(earningsUpdate entity.Earnings) []entity.Earnings
// }

// type AssetUserRepository interface {
// 	CreateAssetUserRelation(assetId string, userUid string) ([]entity.AssetUsers,
// 		error)
// 	DeleteAssetUserRelation(assetId string, userUid string) ([]entity.AssetUsers,
// 		error)
// 	DeleteAssetUserRelationByAsset(assetId string) ([]entity.AssetUsers,
// 		error)
// 	DeleteAssetUserRelationByUser(userUid string) ([]entity.AssetUsers,
// 		error)
// 	SearchAssetUserRelation(assetId string, userUid string) ([]entity.AssetUsers,
// 		error)
// }

type UseCases interface {
	CreateAsset(symbol string, fullname string, preference *string,
		sectorId string, assetType assettype.AssetType) (entity.Asset, error)
}
