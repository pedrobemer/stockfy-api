package asset

import (
	"stockfyApi/entity"
	externalapi "stockfyApi/externalApi"
	assettype "stockfyApi/usecases/assetType"
)

type Repository interface {
	Create(assetInsert entity.Asset) entity.Asset
	Search(symbol string) ([]entity.Asset, error)
	SearchByUser(symbol string, userUid string, orderType string) (
		[]entity.Asset, error)
	SearchPerAssetType(assetType string, country string, userUid string,
		withOrdersInfo bool) []entity.AssetType
	// SearchByOrderId(orderId string) []entity.Asset
	Delete(assetId string) ([]entity.Asset, error)
}

type ExternalApiRepository interface {
	VerifySymbol2(symbol string) entity.SymbolLookup
	GetPrice(symbol string) entity.SymbolPrice
	CompanyOverview(symbol string) map[string]string
}

type UseCases interface {
	CreateAsset(symbol string, fullname string, preference *string,
		sectorId string, assetType assettype.AssetType) (entity.Asset, error)
	SearchAsset(symbol string) (*entity.Asset, error)
	DeleteAsset(assetId string) (*entity.Asset, error)
	SearchAssetByUser(symbol string, userUid string, withOrders bool,
		withOrderResume bool) (*entity.Asset, error)
	SearchAssetPerAssetType(assetType string, country string, userUid string,
		withOrdersInfo bool) (*entity.AssetType, error)
	AssetPreferenceType(symbol string, country string, assetType string) string
	AssetVerificationExistence(symbol string, country string,
		extApi externalapi.ThirdPartyInterfaces) (*entity.SymbolLookup, error)
	AssetVerificationSector(assetType string, symbol string, country string,
		extInterface ExternalApiRepository) string
	AssetVerificationPrice(symbol string, country string,
		extInterface externalapi.ThirdPartyInterfaces) (*entity.SymbolPrice, error)
}
