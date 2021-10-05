package usecases

import (
	"stockfyApi/usecases/asset"
	assettype "stockfyApi/usecases/assetType"
	assetusers "stockfyApi/usecases/assetUser"
	"stockfyApi/usecases/brokerage"
	dbverification "stockfyApi/usecases/dbVerification"
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
	DbVerificationRepository dbverification.Repository
}

type Applications struct {
	AssetApp          asset.Application
	AssetTypeApp      assettype.Application
	AssetUser         assetusers.Application
	SectorApp         sector.Application
	UserApp           user.Application
	OrderApp          order.Application
	Brokerage         brokerage.Application
	DbVerificationApp dbverification.Application
}

type UseCases struct {
	sector         sector.UseCases
	dbVerification dbverification.UseCases
	assetType      assettype.UseCases
	asset          asset.UseCases
}

func NewApplications(repos Repositories, extRepo user.ExternalUserDatabase) *Applications {
	return &Applications{
		SectorApp:         *sector.NewApplication(repos.SectorRepository),
		AssetTypeApp:      *assettype.NewApplication(repos.AssetTypeRepository),
		AssetApp:          *asset.NewApplication(repos.AssetRepository),
		AssetUser:         *assetusers.NewApplication(repos.AssetUserRepository),
		UserApp:           *user.NewApplication(repos.UserRepository, extRepo),
		OrderApp:          *order.NewApplication(repos.OrderRepository),
		Brokerage:         *brokerage.NewApplication(repos.BrokerageRepository),
		DbVerificationApp: *dbverification.NewApplication(repos.DbVerificationRepository),
	}
}
