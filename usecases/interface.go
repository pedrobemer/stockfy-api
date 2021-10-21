package usecases

import (
	"stockfyApi/usecases/asset"
	assettype "stockfyApi/usecases/assetType"
	assetusers "stockfyApi/usecases/assetUser"
	"stockfyApi/usecases/brokerage"
	dbverification "stockfyApi/usecases/dbVerification"
	"stockfyApi/usecases/earnings"
	"stockfyApi/usecases/order"
	"stockfyApi/usecases/sector"
	"stockfyApi/usecases/user"
)

type Repositories struct {
	AssetRepository          asset.Repository
	SectorRepository         sector.Repository
	AssetTypeRepository      assettype.Repository
	UserRepository           user.Repository
	OrderRepository          order.Repository
	AssetUserRepository      assetusers.Repository
	BrokerageRepository      brokerage.Repository
	EarningsRepository       earnings.Repository
	DbVerificationRepository dbverification.Repository
}

type Applications struct {
	AssetApp          asset.UseCases
	AssetTypeApp      assettype.UseCases
	AssetUserApp      assetusers.UseCases
	SectorApp         sector.UseCases
	UserApp           user.UseCases
	OrderApp          order.UseCases
	BrokerageApp      brokerage.UseCases
	EarningsApp       earnings.UseCases
	DbVerificationApp dbverification.UseCases
}

func NewApplications(repos Repositories, extRepo user.ExternalUserDatabase) *Applications {
	return &Applications{
		SectorApp:         sector.NewApplication(repos.SectorRepository),
		AssetTypeApp:      assettype.NewApplication(repos.AssetTypeRepository),
		AssetApp:          asset.NewApplication(repos.AssetRepository),
		AssetUserApp:      assetusers.NewApplication(repos.AssetUserRepository),
		UserApp:           user.NewApplication(repos.UserRepository, extRepo),
		OrderApp:          order.NewApplication(repos.OrderRepository),
		BrokerageApp:      brokerage.NewApplication(repos.BrokerageRepository),
		EarningsApp:       earnings.NewApplication(repos.EarningsRepository),
		DbVerificationApp: dbverification.NewApplication(repos.DbVerificationRepository),
	}
}
