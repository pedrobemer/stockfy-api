package logicApi

import (
	"fmt"
	"stockfyApi/entity"
	externalapi "stockfyApi/externalApi"
	"stockfyApi/usecases"
)

func ApiAssetVerification(app usecases.Applications,
	externalInterfaces externalapi.ThirdPartyInterfaces, symbol string,
	country string) (int, *entity.Asset, error) {

	var symbolLookup *entity.SymbolLookup
	var err error

	// Verify if it is a valid country code. If so, this method verifies the
	// asset existence in the Alpha o Finnhub API
	switch country {
	case "BR":
		symbolLookup, err = app.AssetApp.AssetVerificationExistence(
			symbol, country, &externalInterfaces.AlphaVantageApi)
		break
	case "US":
		symbolLookup, err = app.AssetApp.AssetVerificationExistence(
			symbol, country, &externalInterfaces.FinnhubApi)
		break
	default:
		return 400, nil, entity.ErrInvalidCountryCode
	}

	if err != nil {
		return 404, nil, err
	}

	if country == "US" && symbolLookup.Type == "Equity" {
		companyOverview := externalInterfaces.AlphaVantageApi.
			CompanyOverview(symbolLookup.Symbol)
		symbolLookup.Type = companyOverview["Industry"]
	}
	fmt.Println("After CO:", symbolLookup)

	assetType := app.AssetTypeApp.AssetTypeConversion(
		symbolLookup.Type, country, symbol)
	fmt.Println("AssetType:", assetType)

	// Verify the Sector
	sectorName := app.AssetApp.AssetVerificationSector(
		assetType, symbol, country, &externalInterfaces.FinnhubApi)
	fmt.Println("SectorName:", sectorName)

	// Create Sector
	sectorInfo, err := app.SectorApp.CreateSector(sectorName)
	if err != nil {
		return 500, nil, err
	}
	fmt.Println("Sector Info:", sectorInfo)

	// Search AssetType
	assetTypeInfo, err := app.AssetTypeApp.SearchAssetType(
		assetType, country)
	if err != nil {
		return 500, nil, err
	}

	assetTypeConverted := app.AssetTypeApp.
		AssetTypeConversionToUseCaseStruct(assetTypeInfo[0].Id,
			assetTypeInfo[0].Type, assetTypeInfo[0].Country)

	// Specify the preference asset type if it is a brazilian asset
	preference := app.AssetApp.AssetPreferenceType(symbol,
		country, assetTypeInfo[0].Type)

	// Create Asset
	assetCreated, err := app.AssetApp.CreateAsset(symbol,
		symbolLookup.Fullname, &preference, sectorInfo[0].Id, assetTypeConverted)
	if err != nil {
		return 500, nil, err
	}

	return 200, &assetCreated, nil

}

func ApiCreateOrder(app usecases.Applications,
	externalInterfaces externalapi.ThirdPartyInterfaces, symbol string,
	country string, orderType string, quantity float64, price float64,
	currency string, brokerage string, date string, userUid string) (int,
	*entity.Order, error) {

	var assetInfo *entity.Asset
	httpStatusCode := 200

	err := app.OrderApp.OrderVerification(orderType, country, quantity, price,
		currency)
	if err != nil {
		return 400, nil, err
	}

	// Verify if the asset already exist in our database. If not this asset needs
	// to be created if it is a valid asset
	condAssetExist := "symbol='" + symbol + "'"
	assetExist := app.DbVerificationApp.RowValidation("asset", condAssetExist)

	if !assetExist {
		httpStatusCode, assetInfo, err = ApiAssetVerification(app,
			externalInterfaces, symbol, country)
		if err != nil {
			return httpStatusCode, nil, err
		}

	} else {
		assetInfo, err = app.AssetApp.SearchAsset(symbol)
		if err != nil {
			return 500, nil, err
		}
	}

	// Search in the AssetUser table if the user already invest in the Asset
	// based on its ID.
	assetUser, err := app.AssetUser.SearchAssetUserRelation(assetInfo.Id, userUid)
	if err != nil {
		return 500, nil, err
	}

	// If there isn't any relation between the user and the asset in the AssetUser
	// table, then, it is necessary to create such relation.
	if assetUser == nil {
		assetUser, err = app.AssetUser.CreateAssetUserRelation(assetInfo.Id,
			userUid)
		if err != nil {
			return 500, nil, err
		}
	}

	// Search if the brokerage exists
	brokerageInfo, err := app.Brokerage.SearchBrokerage("SINGLE", brokerage, "")
	if err != nil {
		return 400, nil, err
	}
	brokerageReturn := *brokerageInfo

	// Create Order
	orderReturn, err := app.OrderApp.CreateOrder(quantity, price, currency,
		orderType, date, brokerageReturn[0].Id, assetInfo.Id, userUid)
	if err != nil {
		return 500, nil, err
	}

	return httpStatusCode, orderReturn, nil
}
