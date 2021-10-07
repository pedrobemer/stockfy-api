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

	// Create Order
	orderReturn, err := a.app.OrderApp.CreateOrder(quantity, price, currency,
		orderType, date, brokerageInfo[0].Id, assetInfo.Id, userUid)
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

func (a *Application) ApiDeleteAssets(myUser bool, userUid string,
	symbol string) (int, *entity.Asset, error) {

	var deletedAssetInfo *entity.Asset
	var assetUserDeleted *entity.AssetUsers

	// If the myUser flag is not true, our application assumes that the Asset
	// will be deleted for every user. Only user with admin privileges are enabled
	// for such action. When myUser is false, then, our application assumes that
	// the asset will be deleted only for the user requested for it.
	if !myUser {

		// Search to find if the user has admin privileges
		searchedUser, _ := a.app.UserApp.SearchUser(userUid)
		if searchedUser.Type != "admin" {
			return 405, nil, entity.ErrInvalidApiAuthorization
		}

		// Search the Asset information
		assetInfo, err := a.app.AssetApp.SearchAsset(symbol)
		if err != nil {
			return 400, nil, err
		}

		// Delete the Asset for all the users
		_, err = a.app.AssetUserApp.DeleteAssetUserRelationByAsset(assetInfo.Id)
		if err != nil {
			return 500, nil, err
		}

		// Delete the Orders for this Asset for all the users
		_, err = a.app.OrderApp.DeleteOrdersFromAsset(assetInfo.Id)
		if err != nil {
			return 500, nil, err
		}

		// Delete Earnings for this asset for all the users
		_, err = a.app.EarningsApp.DeleteEarningsFromAsset(assetInfo.Id)

		// Delete Asset from the database
		deletedAsset, err := a.app.AssetApp.DeleteAsset(assetInfo.Id)
		if err != nil {
			return 500, nil, err
		}

		if deletedAsset == nil {
			return 404, nil, entity.ErrInvalidDeleteAsset
		}

		deletedAssetInfo = assetInfo

	} else if myUser {

		// Search if the user has this Asset
		assetInfo, err := a.app.AssetApp.SearchAssetByUser(symbol, userUid,
			false, false, true)
		if err != nil {
			return 404, nil, err
		}

		// Delete Orders from this Asset for a specific user
		_, err = a.app.OrderApp.DeleteOrdersFromAssetUser(
			assetInfo.Id, userUid)
		if err != nil {
			return 500, nil, err
		}

		_, err = a.app.EarningsApp.DeleteEarningsFromAssetUser(assetInfo.Id,
			userUid)

		// Delete Asset for this user
		assetUserDeleted, err = a.app.AssetUserApp.DeleteAssetUserRelation(
			assetInfo.Id, userUid)
		if err != nil {
			return 500, nil, err
		}

		if assetUserDeleted == nil {
			return 400, nil, entity.ErrInvalidAssetUser
		}

		deletedAssetInfo = assetInfo

	} else {
		return 400, nil, entity.ErrInvalidApiRequest
	}

	return 200, deletedAssetInfo, nil
}

func (a *Application) ApiGetOrdersFromAssetUser(symbol string, userUid string) (
	int, []entity.Order, error) {
	if symbol == "" {
		return 400, nil, entity.ErrInvalidApiAssetSymbol
	}

	assetInfo, err := a.app.AssetApp.SearchAssetByUser(symbol, userUid, false,
		false, true)
	if err != nil {
		return 400, nil, err
	}

	ordersInfo, err := a.app.OrderApp.SearchOrdersFromAssetUser(assetInfo.Id,
		userUid)
	if err != nil {
		return 500, nil, err
	}

	if ordersInfo == nil {
		return 404, nil, entity.ErrInvalidOrdersFromAssetUser
	}

	return 200, ordersInfo, nil
}

func (a *Application) ApiUpdateOrdersFromUser(orderId string, userUid string,
	orderType string, price float64, quantity float64, date string,
	brokerage string) (int, *entity.Order, error) {

	if orderType == "" || price == 0 || quantity == 0 || date == "" ||
		brokerage == "" {
		return 400, nil, entity.ErrInvalidApiOrderUpdate
	}

	orderInfo, err := a.app.OrderApp.SearchOrderByIdAndUserUid(orderId, userUid)
	if err != nil {
		return 500, nil, err
	}

	if orderInfo == nil {
		return 400, nil, entity.ErrInvalidOrderId
	}

	err = a.app.OrderApp.OrderVerification(orderType, orderInfo.Brokerage.Country,
		quantity, price, orderInfo.Currency)
	if err != nil {
		return 400, nil, err
	}

	brokerageInfo, err := a.app.BrokerageApp.SearchBrokerage("SINGLE",
		brokerage, "")
	if err != nil {
		return 400, nil, err
	}

	updatedOrder, err := a.app.OrderApp.UpdateOrder(orderId, userUid, price,
		quantity, orderType, date, brokerageInfo[0].Id, orderInfo.Currency)
	if err != nil {
		return 500, nil, err
	}

	return 200, updatedOrder, nil
}

func (a *Application) ApiCreateEarnings(symbol string, currency string,
	earningType string, date string, earnings float64, userUid string) (int,
	*entity.Earnings, error) {

	err := a.app.EarningsApp.EarningsVerification(symbol, currency, earningType,
		date, earnings)
	if err != nil {
		return 400, nil, err
	}

	assetInfo, err := a.app.AssetApp.SearchAssetByUser(symbol, userUid, false,
		false, true)
	if err != nil {
		return 400, nil, err
	}

	if assetInfo == nil {
		return 400, nil, entity.ErrInvalidApiEarningSymbol
	}

	earningCreated, err := a.app.EarningsApp.CreateEarning(earningType, earnings,
		currency, date, assetInfo.AssetType.Country, assetInfo.Id, userUid)
	if err != nil {
		return 400, nil, err
	}

	return 200, earningCreated, nil
}

func (a *Application) ApiGetEarningsFromAssetUser(symbol string, userUid string) (
	int, []entity.Earnings, error) {

	if symbol == "" {
		return 400, nil, entity.ErrInvalidApiRequest
	}

	assetInfo, err := a.app.AssetApp.SearchAssetByUser(symbol, userUid, false,
		false, true)
	if err != nil {
		return 400, nil, err
	}

	if assetInfo == nil {
		return 404, nil, entity.ErrInvalidAssetSymbol
	}

	earningsReturn, err := a.app.EarningsApp.SearchEarningsFromAssetUser(
		assetInfo.Id, userUid)
	if err != nil {
		return 500, nil, err
	}

	if earningsReturn == nil {
		return 404, nil, entity.ErrInvalidApiEarningAssetUser
	}

	return 200, earningsReturn, nil
}
