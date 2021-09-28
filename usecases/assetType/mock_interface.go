package assettype

import "stockfyApi/entity"

type Mock struct {
}

func NewMockRepo() *Mock {
	return &Mock{}
}

func (m *Mock) Search(searchType string, name string,
	country string) ([]entity.AssetType,
	error) {

	if searchType == "SPECIFIC" {
		ReturnedAssetTypes := []entity.AssetType{
			{
				Id:      "1",
				Type:    "STOCK",
				Country: "US",
				Name:    "Ações EUA",
			},
		}
		return ReturnedAssetTypes, nil
	} else if searchType == "ONLYTYPE" {
		ReturnedAssetTypes := []entity.AssetType{
			{
				Id:      "1",
				Type:    "STOCK",
				Country: "US",
				Name:    "Ações EUA",
			},
			{
				Id:      "2",
				Type:    "STOCK",
				Country: "BR",
				Name:    "Ações Brasil",
			},
		}
		return ReturnedAssetTypes, nil
	} else if searchType == "ONLYCOUNTRY" {
		ReturnedAssetTypes := []entity.AssetType{
			{
				Id:      "1",
				Type:    "STOCK",
				Country: "US",
				Name:    "Ações EUA",
			},
			{
				Id:      "3",
				Type:    "ETF",
				Country: "US",
				Name:    "ETFs EUA",
			},
			{
				Id:      "5",
				Type:    "REIT",
				Country: "US",
				Name:    "REITs",
			},
		}
		return ReturnedAssetTypes, nil
	} else {
		ReturnedAssetTypes := []entity.AssetType{
			{
				Id:      "1",
				Type:    "STOCK",
				Country: "US",
				Name:    "Ações EUA",
			},
			{
				Id:      "2",
				Type:    "STOCK",
				Country: "BR",
				Name:    "Ações Brasil",
			},
			{
				Id:      "3",
				Type:    "ETF",
				Country: "US",
				Name:    "ETFs EUA",
			},
			{
				Id:      "4",
				Type:    "ETF",
				Country: "BR",
				Name:    "ETFs Brasil",
			},
			{
				Id:      "5",
				Type:    "REIT",
				Country: "US",
				Name:    "REITs",
			},
			{
				Id:      "6",
				Type:    "FII",
				Country: "BR",
				Name:    "FIIs",
			},
		}
		return ReturnedAssetTypes, nil
	}

}
