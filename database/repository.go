package database

type SectorRepository interface {
	CreateSector(sector string) ([]Sector, error)
	FetchSectorByName(sector string) ([]Sector, error)
	FetchSectorByAsset(symbol string) ([]Sector, error)
}

type AssetRepository interface {
	CreateAsset(assetInsert Asset, assetTypeId string, sectorId string) Asset
	SearchAsset(symbol string) ([]Asset, error)
	SearchAssetByUser(symbol string, userUid string, orderType string) (
		[]Asset, error)
	SearchAssetsPerAssetType(assetType string, country string, userUid string,
		withOrdersInfo bool) []AssetType
	SearchAssetByOrderId(orderId string) []Asset
	DeleteAsset(assetId string) []Asset
}

type BrokerageRepository interface {
	FetchBrokerage(specificFetch string, args ...string) ([]Brokerage, error)
}

type AssetTypeRepository interface {
	FetchAssetType(fetchType string, args ...string) (AssetType, error)
}

type OrderRepository interface {
	CreateOrder(orderInsert Order) Order
	SearchOrdersFromAssetUser(assetId string, userUid string) ([]Order, error)
	DeleteOrderFromUser(id string, userUid string) string
	DeleteOrdersFromAsset(symbolId string) []Order
	DeleteOrdersFromAssetUser(assetId string, userUid string) ([]Order, error)
	UpdateOrderFromUser(orderUpdate Order) []Order
}

type UsersRepository interface {
	CreateUser(signUp Users) ([]Users, error)
	DeleteUser(firebaseUid string) ([]Users, error)
	UpdateUser(userInfo Users) ([]Users, error)
	SearchUser(userUid string) ([]Users, error)
}

type EarningsRepository interface {
	CreateEarningRow(earningOrder Earnings) []Earnings
	SearchEarningFromAssetUser(assetId string, userUid string) ([]Earnings,
		error)
	DeleteEarningsFromAssetUser(assetId string, userUid string) ([]Earnings,
		error)
	DeleteEarningFromUser(id string, userUid string) string
	UpdateEarningsFromUser(earningsUpdate Earnings) []Earnings
}

type AssetUserRepository interface {
	CreateAssetUserRelation(assetId string, userUid string) ([]AssetUsers, error)
	DeleteAssetUserRelation(assetId string, userUid string) ([]AssetUsers, error)
	DeleteAssetUserRelationByAsset(assetId string) ([]AssetUsers, error)
	DeleteAssetUserRelationByUser(userUid string) ([]AssetUsers, error)
	SearchAssetUserRelation(assetId string, userUid string) ([]AssetUsers, error)
}

type Repository interface {
	SectorRepository
	AssetRepository
	BrokerageRepository
	OrderRepository
	UsersRepository
	EarningsRepository
	AssetUserRepository
}
