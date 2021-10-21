package assettype

import "stockfyApi/entity"

type Repository interface {
	Search(searchType string, name string,
		country string) ([]entity.AssetType, error)
}

type UseCases interface {
	SearchAssetType(name string, country string) ([]entity.AssetType, error)
	AssetTypeConversionToUseCaseStruct(id string, assetType string,
		country string) AssetType
	AssetTypeConversion(assetType string, country string, symbol string) string
}
