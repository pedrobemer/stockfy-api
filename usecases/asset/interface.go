package asset

import (
	"stockfyApi/entity"
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
	// Delete(assetId string) []entity.Asset
}

type ExternalApiRepository interface {
	VerifySymbol2(symbol string) entity.SymbolLookup
	GetPrice(symbol string) entity.SymbolPrice
	CompanyOverview(symbol string) map[string]string
}

type UseCases interface {
	CreateAsset(symbol string, fullname string, preference *string,
		sectorId string, assetType assettype.AssetType) (entity.Asset, error)
}
