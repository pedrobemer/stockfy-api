package logicApi

import (
	"fmt"
	"stockfyApi/entity"
	externalapi "stockfyApi/externalApi"
	"stockfyApi/usecases"
)

type Application struct {
	app                usecases.Applications
	externalInterfaces externalapi.ThirdPartyInterfaces
}

//NewApplication create new use case
func NewApplication(a usecases.Applications,
	ext externalapi.ThirdPartyInterfaces) *Application {
	return &Application{
		app:                a,
		externalInterfaces: ext,
	}
}

func (a *Application) ApiAssetVerification(symbol string, country string) (
	int, *entity.Asset, error) {

	var symbolLookup *entity.SymbolLookup
	var err error

	// Verify if it is a valid country code. If so, this method verifies the
	// asset existence in the Alpha o Finnhub API
	switch country {
	case "BR":
		symbolLookup, err = a.app.AssetApp.AssetVerificationExistence(
			symbol, country, &a.externalInterfaces.AlphaVantageApi)
		break
	case "US":
		symbolLookup, err = a.app.AssetApp.AssetVerificationExistence(
			symbol, country, &a.externalInterfaces.FinnhubApi)
		break
	default:
		return 400, nil, entity.ErrInvalidCountryCode
	}

	if err != nil {
		return 404, nil, err
	}

	if country == "US" && symbolLookup.Type == "Equity" {
		companyOverview := a.externalInterfaces.AlphaVantageApi.
			CompanyOverview(symbolLookup.Symbol)
		symbolLookup.Type = companyOverview["Industry"]
	}
	fmt.Println("After CO:", symbolLookup)

	assetType := a.app.AssetTypeApp.AssetTypeConversion(
		symbolLookup.Type, country, symbol)
	fmt.Println("AssetType:", assetType)

	// Verify the Sector
	sectorName := a.app.AssetApp.AssetVerificationSector(
		assetType, symbol, country, &a.externalInterfaces.FinnhubApi)
	fmt.Println("SectorName:", sectorName)

	// Create Sector
	sectorInfo, err := a.app.SectorApp.CreateSector(sectorName)
	if err != nil {
		return 500, nil, err
	}
	fmt.Println("Sector Info:", sectorInfo)

	// Search AssetType
	assetTypeInfo, err := a.app.AssetTypeApp.SearchAssetType(
		assetType, country)
	if err != nil {
		return 500, nil, err
	}

	assetTypeConverted := a.app.AssetTypeApp.
		AssetTypeConversionToUseCaseStruct(assetTypeInfo[0].Id,
			assetTypeInfo[0].Type, assetTypeInfo[0].Country)

	// Specify the preference asset type if it is a brazilian asset
	preference := a.app.AssetApp.AssetPreferenceType(symbol,
		country, assetTypeInfo[0].Type)

	// Create Asset
	assetCreated, err := a.app.AssetApp.CreateAsset(symbol,
		symbolLookup.Fullname, &preference, sectorInfo[0].Id, assetTypeConverted)
	if err != nil {
		return 500, nil, err
	}

	return 200, &assetCreated, nil

}

func (a *Application) ApiCreateOrder(symbol string, country string,
	orderType string, quantity float64, price float64, currency string,
	brokerage string, date string, userUid string) (int, *entity.Order, error) {

	var assetInfo *entity.Asset
	httpStatusCode := 200

	err := a.app.OrderApp.OrderVerification(orderType, country, quantity, price,
		currency)
	if err != nil {
		return 400, nil, err
	}

	// Verify if the asset already exist in our database. If not this asset needs
	// to be created if it is a valid asset
	condAssetExist := "symbol='" + symbol + "'"
	assetExist := a.app.DbVerificationApp.RowValidation("asset", condAssetExist)

	if !assetExist {
		httpStatusCode, assetInfo, err = a.ApiAssetVerification(symbol, country)
		if err != nil {
			return httpStatusCode, nil, err
		}

	} else {
		assetInfo, err = a.app.AssetApp.SearchAsset(symbol)
		if err != nil {
			return 500, nil, err
		}
	}

	// Search in the AssetUser table if the user already invest in the Asset
	// based on its ID.
	assetUser, err := a.app.AssetUserApp.SearchAssetUserRelation(assetInfo.Id,
		userUid)
	if err != nil {
		return 500, nil, err
	}

	// If there isn't any relation between the user and the asset in the AssetUser
	// table, then, it is necessary to create such relation.
	if assetUser == nil {
		assetUser, err = a.app.AssetUserApp.CreateAssetUserRelation(assetInfo.Id,
			userUid)
		if err != nil {
			return 500, nil, err
		}
	}

	// Search if the brokerage exists
	brokerageInfo, err := a.app.BrokerageApp.SearchBrokerage("SINGLE",
		brokerage, "")
	if err != nil {
		return 400, nil, err
	}
	brokerageReturn := *brokerageInfo

	// Create Order
	orderReturn, err := a.app.OrderApp.CreateOrder(quantity, price, currency,
		orderType, date, brokerageReturn[0].Id, assetInfo.Id, userUid)
	if err != nil {
		return 500, nil, err
	}

	return httpStatusCode, orderReturn, nil
}

func (a *Application) ApiAssetsPerAssetType(assetType string, country string,
	ordersInfo bool, userUid string) (int, *entity.AssetType,
	error) {

	if assetType == "" || country == "" {
		return 400, nil, entity.ErrInvalidApiRequest
	}

	searchedAssetType, err := a.app.AssetApp.SearchAssetPerAssetType(assetType,
		country, userUid, ordersInfo)
	if err != nil {
		return 400, nil, err
	}

	return 200, searchedAssetType, nil

}
