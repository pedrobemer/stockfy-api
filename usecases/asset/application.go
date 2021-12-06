package asset

import (
	"stockfyApi/entity"
	externalapi "stockfyApi/externalApi"
	assettype "stockfyApi/usecases/assetType"
	"stockfyApi/usecases/general"
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

func (a *Application) SearchAsset(symbol string) (*entity.Asset, error) {
	asset, err := a.repo.Search(symbol)
	if err != nil {
		return nil, err
	}

	if asset == nil {
		return nil, nil
	}

	return &asset[0], nil
}

func (a *Application) DeleteAsset(assetId string) (*entity.Asset, error) {
	deletedAsset, err := a.repo.Delete(assetId)
	if err != nil {
		return nil, err
	}

	if deletedAsset == nil {
		return nil, nil
	}

	return &deletedAsset[0], nil
}

func (a *Application) SearchAssetByUser(symbol string, userUid string,
	withOrders bool, withOrderResume bool) (*entity.Asset, error) {
	orderType := ""

	if withOrders && !withOrderResume {
		orderType = "ONLYORDERS"
	} else if withOrders && withOrderResume {
		orderType = "ALL"
	} else if !withOrders && withOrderResume {
		orderType = "ONLYINFO"
	}

	asset, err := a.repo.SearchByUser(symbol, userUid, orderType)
	if err != nil {
		return nil, err
	}

	if asset == nil {
		return nil, nil
	}

	return &asset[0], nil
}

func (a *Application) SearchAssetPerAssetType(assetType string, country string,
	userUid string, withOrdersInfo bool) (*entity.AssetType, error) {

	err := general.CountryValidation(country)
	if err != nil {
		return nil, err
	}

	err = general.AssetTypeNameValidation(assetType)
	if err != nil {
		return nil, err
	}

	assetsPerAssetType := a.repo.SearchPerAssetType(assetType, country, userUid,
		withOrdersInfo)
	if assetsPerAssetType == nil {
		return nil, entity.ErrInvalidAssetType
	}

	return &assetsPerAssetType[0], nil
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

// func (a *Application) AssetVerificationExistence(symbol string, country string,
// 	extApi ExternalApiRepository) (*entity.SymbolLookup, error) {
func (a *Application) AssetVerificationExistence(symbol string, country string,
	extApi externalapi.ThirdPartyInterfaces) (*entity.SymbolLookup, error) {

	var symbolLookup entity.SymbolLookup

	if symbol == "" {
		return nil, entity.ErrInvalidApiQuerySymbolBlank
	}

	if err := general.CountryValidation(country); err != nil {
		return nil, err
	}

	if country == "BR" {
		symbol = symbol + ".SA"
		symbolLookup = extApi.AlphaVantageApi.VerifySymbol2(symbol)
	} else {
		symbolLookup = extApi.FinnhubApi.VerifySymbol2(symbol)
	}

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

func (a *Application) AssetVerificationPrice(symbol string, country string,
	extInterface externalapi.ThirdPartyInterfaces) (*entity.SymbolPrice, error) {

	var symbolPrice entity.SymbolPrice

	if err := general.CountryValidation(country); err != nil {
		return nil, err
	}

	if symbol == "" {
		return nil, entity.ErrInvalidApiQuerySymbolBlank
	}

	if country == "BR" {
		symbol = symbol + ".SA"
		symbolPrice = extInterface.AlphaVantageApi.GetPrice(symbol)
	} else {
		symbolPrice = extInterface.FinnhubApi.GetPrice(symbol)
	}

	if symbolPrice.CurrentPrice == 0 {
		return nil, entity.ErrInvalidAssetSymbol
	}

	return &symbolPrice, nil
}
