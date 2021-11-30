package assettype

import (
	"stockfyApi/entity"
	"stockfyApi/usecases/general"
)

type AssetType struct {
	Id      string
	Type    string
	Country string
}

type Application struct {
	repo Repository
}

//NewApplication create new use case
func NewApplication(r Repository) *Application {
	return &Application{
		repo: r,
	}
}

func (a *Application) SearchAssetType(name string, country string) (
	[]entity.AssetType, error) {
	var searchType string

	if country != "" {
		err := general.CountryValidation(country)
		if err != nil {
			return nil, err
		}
	}

	if name != "" {
		err := general.AssetTypeNameValidation(name)
		if err != nil {
			return nil, err
		}
	}

	if name == "" && country == "" {
		searchType = ""
	} else if name != "" && country == "" {
		searchType = "ONLYTYPE"
	} else if name == "" && country != "" {
		searchType = "ONLYCOUNTRY"
	} else if name != "" && country != "" {
		searchType = "SPECIFIC"
	}

	assetTypeReturn, err := a.repo.Search(searchType, name, country)
	if err != nil {
		return nil, err
	}

	return assetTypeReturn, nil
}

func (a *Application) AssetTypeConversionToUseCaseStruct(id string,
	assetType string, country string) AssetType {
	return AssetType{
		Id:      id,
		Type:    assetType,
		Country: country,
	}
}

func (a *Application) AssetTypeConversion(assetType string, country string,
	symbol string) string {
	if assetType == "ETP" {
		return "ETF"
	} else if assetType == "REIT" {
		return "REIT"
	} else if assetType == "Common Stock" {
		return "STOCK"
	} else if country == "BR" && assetType == "ETF" {
		for _, validEtf := range entity.ListValidBrETF {
			if symbol == validEtf {
				return "ETF"
			}
		}
		return "FII"
	} else if country == "BR" && assetType == "Equity" {
		return "STOCK"
	} else if country != "BR" && assetType == "REAL ESTATE INVESTMENT TRUSTS" {
		return "REIT"
	} else if country != "BR" && assetType != "REAL ESTATE INVESTMENT TRUSTS" {
		return "STOCK"
	}

	return assetType
}
