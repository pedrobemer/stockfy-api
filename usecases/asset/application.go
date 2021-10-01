package asset

import (
	"stockfyApi/entity"
	assettype "stockfyApi/usecases/assetType"
)

type Application struct {
	repo Repository
}

//NewApplication create new use case
func NewApplication(r Repository) *Application {
	return &Application{
		repo: r,
	}
}

func (a *Application) CreateAsset(symbol string, fullname string,
	preference *string, sectorId string, assetType assettype.AssetType) (
	entity.Asset, error) {

	assetInfo, err := entity.NewAsset(symbol, fullname, preference, sectorId,
		assetType.Id, assetType.Type, assetType.Country)

	assetCreated := a.repo.Create(*assetInfo)
	return assetCreated, err
}

func (a *Application) AssetPreferenceType(symbol string, country string,
	assetType string) string {

	var preference string

	if country == "BR" && assetType == "STOCK" {
		switch symbol[len(symbol)-1:] {
		case "3":
			preference = "ON"
			break
		case "4":
			preference = "PN"
			break
		case "1":
			preference = "UNIT"
			break
		default:
			preference = ""
			break
		}
	}

	return preference
}

func (a *Application) AssetVerificationExistence(symbol string,
	extApi ExternalApiRepository) (*entity.SymbolLookup, error) {

	symbolLookup := extApi.VerifySymbol2(symbol)
	if symbolLookup.Symbol == "" {
		return nil, entity.ErrInvalidAssetSymbol
	}

	return &symbolLookup, nil
}

func (a *Application) AssetVerificationSector(assetType string, symbol string,
	country string, extInterface ExternalApiRepository) string {

	if country == "BR" {
		symbol = symbol + ".SA"
	}

	if assetType == "STOCK" {
		companyOverview := extInterface.CompanyOverview(symbol)
		return companyOverview["finnhubIndustry"]
	} else if assetType == "ETF" {
		return "Blend"
	} else {
		return "Real Estate"
	}
}
