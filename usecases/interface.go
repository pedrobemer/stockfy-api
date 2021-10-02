package usecases

import (
	"stockfyApi/usecases/asset"
	assettype "stockfyApi/usecases/assetType"
	dbverification "stockfyApi/usecases/dbVerification"
	"stockfyApi/usecases/sector"
	"stockfyApi/usecases/user"
)

type Repositories struct {
	AssetRepository          asset.Repository
	SectorRepository         sector.Repository
	AssetTypeRepository      assettype.Repository
	UserRepository           user.Repository
	DbVerificationRepository dbverification.Repository
}

type Applications struct {
	AssetApp          asset.Application
	AssetTypeApp      assettype.Application
	SectorApp         sector.Application
	UserApp           user.Application
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
		UserApp:           *user.NewApplication(repos.UserRepository, extRepo),
		DbVerificationApp: *dbverification.NewApplication(repos.DbVerificationRepository),
	}
}
