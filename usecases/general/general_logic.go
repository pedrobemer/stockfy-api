package general

import (
	"stockfyApi/entity"
)

func CountryValidation(country string) error {
	if country != "BR" && country != "US" && country != "" {
		return entity.ErrInvalidCountryCode
	}

	return nil
}

func AssetTypeNameValidation(name string) error {
	if name != "STOCK" && name != "ETF" && name != "REIT" && name != "FII" &&
		name != "" {
		return entity.ErrInvalidAssetTypeName
	}
	return nil
}
