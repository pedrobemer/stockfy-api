package entity

func NewAsset(symbol string, fullname string, preference *string,
	sectorId string, assetTypeId string, assetType string, country string) (
	*Asset, error) {

	asset := &Asset{
		Symbol:     symbol,
		Fullname:   fullname,
		Preference: preference,
		Sector:     &Sector{Id: sectorId},
		AssetType:  &AssetType{Id: assetTypeId},
	}

	err := asset.Validate(assetType, country)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func (a *Asset) Validate(assetType string, country string) error {
	if a.Symbol == "" || a.Fullname == "" || a.AssetType.Id == "" ||
		a.Sector.Id == "" {
		return ErrInvalidAssetEntity

	}

	if a.Preference == nil && assetType == "STOCK" && country == "BR" {
		return ErrInvalidAssetEntity
	}

	if (assetType != "STOCK" && assetType != "ETF" && assetType != "REIT" &&
		assetType != "FII") || (country != "BR" && country != "US") {
		return ErrInvalidAssetEntity
	}

	return nil
}
