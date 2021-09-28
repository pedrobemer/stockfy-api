package assettype

import "stockfyApi/entity"

type AssetTypeRepository interface {
	Search(searchType string, name string,
		country string) ([]entity.AssetType, error)
}
