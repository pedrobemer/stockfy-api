package logicApi

import (
	"stockfyApi/entity"
	externalapi "stockfyApi/externalApi"
	"stockfyApi/usecases"
	"strconv"
	"strings"
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

	symbolLookup, err := a.app.AssetApp.AssetVerificationExistence(symbol,
		country, a.externalInterfaces)

	if err != nil {
		if err.Error() == entity.ErrInvalidAssetSymbol.Error() {
			return 404, nil, err
		}

		return 400, nil, err
	}

	if country == "US" && symbolLookup.Type == "Equity" {
		companyOverview := a.externalInterfaces.AlphaVantageApi.
			CompanyOverview(symbolLookup.Symbol)
		symbolLookup.Type = companyOverview["Industry"]
	}

	assetType := a.app.AssetTypeApp.AssetTypeConversion(
		symbolLookup.Type, country, symbol)

	// Verify the Sector
	sectorName := a.app.AssetApp.AssetVerificationSector(
		assetType, symbol, country, a.externalInterfaces.FinnhubApi)

	// Create Sector
	sectorInfo, err := a.app.SectorApp.CreateSector(sectorName)
	if err != nil {
		return 500, nil, err
	}

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
	assetExist := a.app.DbVerificationApp.RowValidation("assets", condAssetExist)

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

	// Create a new AssetUser relation, if the relation does not exist already.
	// If the relation exists this method will return an error equal to:
	// entity.ErrinvalidAssetUserAlreadyExists. In this case, our API will simple
	// ignore the error because the relation already exist and hence we can
	// create the order.
	_, err = a.app.AssetUserApp.CreateAssetUserRelation(assetInfo.Id, userUid)
	if err != nil {
		if err.Error() != entity.ErrinvalidAssetUserAlreadyExists.Error() {
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
	ordersInfo bool, withPrice bool, userUid string) (int, *entity.AssetType,
	error) {

	chPrice := make(chan *entity.SymbolPrice)
	defer close(chPrice)
	var assetPrice *entity.SymbolPrice

	if assetType == "" {
		return 400, nil, entity.ErrInvalidApiQueryTypeBlank
	}

	if country == "" {
		return 400, nil, entity.ErrInvalidApiQueryCountryBlank
	}

	searchedAssetType, err := a.app.AssetApp.SearchAssetPerAssetType(assetType,
		country, userUid, ordersInfo)
	if err != nil {
		return 400, nil, err
	}

	if withPrice == true {
		for _, assetInfo := range searchedAssetType.Assets {
			go func(assetSymbol string) {
				var assetPrice *entity.SymbolPrice

				assetPrice, err = a.app.AssetApp.AssetVerificationPrice(
					assetSymbol, searchedAssetType.Country, a.externalInterfaces)

				chPrice <- assetPrice
			}(assetInfo.Symbol)
		}

		for i := 0; i < len(searchedAssetType.Assets); i++ {
			assetPrice = <-chPrice
			if assetPrice != nil {
				for i, assetInfo := range searchedAssetType.Assets {
					if assetPrice.Symbol == assetInfo.Symbol {
						searchedAssetType.Assets[i].Price = assetPrice
					}
				}
			}

		}

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
			return 403, nil, entity.ErrInvalidUserAdminPrivilege
		}

		// Search the Asset information
		assetInfo, err := a.app.AssetApp.SearchAsset(symbol)
		if err != nil {
			return 500, nil, err
		}

		if assetInfo == nil {
			return 404, nil, entity.ErrInvalidAssetSymbol
		}

		// Delete Asset from the database
		deletedAsset, err := a.app.AssetApp.DeleteAsset(assetInfo.Id)
		if err != nil {
			return 500, nil, err
		}

		if deletedAsset == nil {
			return 404, nil, entity.ErrInvalidDeleteAsset
		}

		deletedAssetInfo = assetInfo

	} else {

		// Search if the user has this Asset
		assetInfo, err := a.app.AssetApp.SearchAssetByUser(symbol, userUid,
			false, false)
		if err != nil {
			return 500, nil, err
		}

		if assetInfo == nil {
			return 404, nil, entity.ErrInvalidAssetSymbol
		}

		// Delete Orders from this Asset for a specific user
		_, err = a.app.OrderApp.DeleteOrdersFromAssetUser(
			assetInfo.Id, userUid)
		if err != nil {
			return 500, nil, err
		}

		_, err = a.app.EarningsApp.DeleteEarningsFromAssetUser(assetInfo.Id,
			userUid)
		if err != nil {
			return 500, nil, err
		}

		// Delete Asset for this user
		assetUserDeleted, err = a.app.AssetUserApp.DeleteAssetUserRelation(
			assetInfo.Id, userUid)
		if err != nil {
			return 500, nil, err
		}

		if assetUserDeleted == nil {
			return 404, nil, entity.ErrInvalidAssetUser
		}

		deletedAssetInfo = assetInfo
	}

	return 200, deletedAssetInfo, nil
}

func (a *Application) ApiGetOrdersFromAssetUser(symbol string, userUid string,
	orderBy string, limit string, offset string) (int, []entity.Order, error) {
	var ordersInfo []entity.Order

	if symbol == "" {
		return 400, nil, entity.ErrInvalidApiQuerySymbolBlank
	}

	assetInfo, err := a.app.AssetApp.SearchAssetByUser(symbol, userUid, false,
		false)
	if err != nil {
		return 500, nil, err
	}

	if assetInfo == nil {
		return 404, nil, entity.ErrInvalidAssetSymbol
	}

	if limit == "" && offset == "" {
		ordersInfo, err = a.app.OrderApp.SearchOrdersFromAssetUser(assetInfo.Id,
			userUid)
		if err != nil {
			return 500, nil, err
		}
	} else {

		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return 400, nil, entity.ErrInvalidOrderLimit
		}

		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			return 400, nil, entity.ErrInvalidOrderOffset
		}

		ordersInfo, err = a.app.OrderApp.SearchOrdersSearchFromAssetUserByDate(
			assetInfo.Id, userUid, orderBy, limitInt, offsetInt)
		if err != nil {
			return 500, nil, err
		}
	}

	if ordersInfo == nil {
		return 404, nil, entity.ErrInvalidOrder
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
		if strings.Contains(err.Error(), "ERROR: invalid input syntax for type uuid:") {
			return 400, nil, err
		}

		return 500, nil, err
	}

	if orderInfo == nil {
		return 404, nil, entity.ErrInvalidOrder
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
		false)
	if err != nil {
		return 500, nil, err
	}

	if assetInfo == nil {
		return 404, nil, nil
	}

	earningCreated, err := a.app.EarningsApp.CreateEarning(earningType, earnings,
		currency, date, assetInfo.AssetType.Country, assetInfo.Id, userUid)
	if err != nil {
		return 400, nil, err
	}

	return 200, earningCreated, nil
}

func (a *Application) ApiGetEarningsFromAssetUser(symbol string, userUid string,
	orderBy string, limit string, offset string) (
	int, []entity.Earnings, error) {

	var earningsReturn []entity.Earnings

	if symbol == "" {
		return 400, nil, entity.ErrInvalidApiQuerySymbolBlank
	}

	assetInfo, err := a.app.AssetApp.SearchAssetByUser(symbol, userUid, false,
		false)
	if err != nil {
		return 500, nil, err
	}

	if assetInfo == nil {
		return 404, nil, entity.ErrMessageApiAssetSymbolUser
	}

	if limit == "" && offset == "" {
		earningsReturn, err = a.app.EarningsApp.SearchEarningsFromAssetUser(
			assetInfo.Id, userUid)
		if err != nil {
			return 500, nil, err
		}
	} else {

		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return 400, nil, entity.ErrInvalidEarningsLimit
		}

		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			return 400, nil, entity.ErrInvalidEarningsOffset
		}

		earningsReturn, err = a.app.EarningsApp.SearchEarningsFromAssetUserByDate(
			assetInfo.Id, userUid, orderBy, limitInt, offsetInt)
		if err != nil {
			return 400, nil, err
		}
	}

	if earningsReturn == nil {
		return 404, nil, entity.ErrMessageApiEarningAssetUser
	}

	return 200, earningsReturn, nil
}

func (a *Application) ApiUpdateEarningsFromUser(earningId string, earning float64,
	earningType string, date string, userUid string) (int, *entity.Earnings,
	error) {

	// Get actual information about the requested earning for update
	searchedEarning, err := a.app.EarningsApp.SearchEarningsFromUser(earningId,
		userUid)
	if err != nil {
		return 500, nil, err
	}

	if searchedEarning == nil {
		return 404, nil, entity.ErrMessageApiEarningId
	}

	// Get information about the asset associated to this earning
	assetInfo, err := a.app.AssetApp.SearchAsset(searchedEarning.Asset.Symbol)
	if err != nil {
		return 500, nil, err
	}

	// Verification if the information received in the body attends the
	// requirements of the Earning table
	err = a.app.EarningsApp.EarningsVerification(searchedEarning.Asset.Symbol,
		searchedEarning.Currency, earningType, date, earning)
	if err != nil {
		return 400, nil, err
	}

	// Update the earning information of the earning with specific ID
	earningsUpdate, err := a.app.EarningsApp.EarningsUpdate(earningType,
		earning, searchedEarning.Currency, date, assetInfo.AssetType.Country,
		earningId, userUid)
	if err != nil {
		return 500, nil, err
	}

	return 200, earningsUpdate, nil
}

func (a *Application) ApiGetAssetByUser(symbol string, userUid string,
	withOrders bool, withOrderResume bool, withPrice bool) (int, *entity.Asset,
	error) {

	var assetPrice *entity.SymbolPrice
	var err error
	chAssetInfo := make(chan *entity.Asset)

	chAssetInfoErr := make(chan error)
	chPrice := make(chan *entity.SymbolPrice)
	chPriceErr := make(chan error)

	go func() {
		assetInfo, err := a.app.AssetApp.SearchAsset(symbol)
		chAssetInfoErr <- err
		chAssetInfo <- assetInfo
		close(chAssetInfoErr)
		close(chAssetInfo)
	}()

	if err := <-chAssetInfoErr; err != nil {
		return 500, nil, entity.ErrInvalidAssetSymbol
	}

	assetInfo := <-chAssetInfo
	if assetInfo == nil {
		return 400, nil, entity.ErrInvalidAssetSymbol
	}

	if withPrice == true {
		go func() {
			var assetPrice *entity.SymbolPrice
			var err error

			assetPrice, err = a.app.AssetApp.AssetVerificationPrice(
				assetInfo.Symbol, assetInfo.AssetType.Country,
				a.externalInterfaces)

			chPrice <- assetPrice
			chPriceErr <- err
			close(chPrice)
			close(chPriceErr)
		}()
	}

	searchedAsset, err := a.app.AssetApp.SearchAssetByUser(
		assetInfo.Symbol, userUid, withOrders, withOrderResume)
	if err != nil {
		return 500, nil, err
	}

	if searchedAsset == nil {
		return 404, nil, nil
	}

	if withPrice == true {
		assetPrice = <-chPrice
		err = <-chPriceErr
		if err != nil {
			return 400, nil, err
		}
	}

	searchedAsset.Price = assetPrice

	return 200, searchedAsset, nil
}

func (a *Application) ApiCreateEvent(symbol string, symbolDemerger string,
	orderType string, eventRate float64, price float64, currency string,
	date string, userUid string) (int, []entity.Order, error) {

	type createOrderGoRoutine struct {
		assetInfo      entity.Asset
		orderType      string
		price          float64
		gainedQuantity float64
		currency       string
		brokerageName  string
		date           string
		userUid        string
	}

	type createOrderResponse struct {
		StatusCode int
		Order      *entity.Order
		Error      error
	}

	var orderInfo createOrderGoRoutine
	var eventsCreated []entity.Order

	err := a.app.OrderApp.EventTypeValueVerification(orderType)
	if err != nil {
		return 400, nil, err
	}

	assetInfo, err := a.app.AssetApp.SearchAssetByUser(symbol, userUid, false,
		false)
	if err != nil {
		return 500, nil, err
	}

	if assetInfo == nil {
		return 400, nil, entity.ErrInvalidAssetSymbolUserRelation
	}

	quantityPerBrokerage, err :=
		a.app.OrderApp.MeasureAssetTotalQuantityForSpecificDate(assetInfo.Id,
			userUid, date)
	if err != nil {
		return 400, nil, err
	}

	apiCreateOrderResponses := []chan createOrderResponse{}
	i := 0
	gainedQuantity := 0.0
	for brokerageName, quantity := range quantityPerBrokerage {
		if orderType != "demerge" {
			gainedQuantity = quantity / eventRate
		}

		apiCreateOrderResponses = append(apiCreateOrderResponses,
			make(chan createOrderResponse))

		orderInfo = createOrderGoRoutine{
			assetInfo: *assetInfo,
			orderType: orderType,
			price: func() float64 {
				if orderType == "demerge" {
					return price * quantity
				} else {
					return price
				}
			}(),
			gainedQuantity: gainedQuantity,
			currency:       currency,
			brokerageName:  brokerageName,
			date:           date,
			userUid:        userUid,
		}

		go func(orderInfo createOrderGoRoutine,
			apiOrderResponse chan createOrderResponse) {

			httpStatusCode, orderCreated, err := a.ApiCreateOrder(
				orderInfo.assetInfo.Symbol, orderInfo.assetInfo.AssetType.Country,
				orderInfo.orderType, orderInfo.gainedQuantity, orderInfo.price,
				orderInfo.currency, orderInfo.brokerageName, orderInfo.date,
				orderInfo.userUid)

			apiOrderResponse <- createOrderResponse{
				StatusCode: httpStatusCode,
				Order:      orderCreated,
				Error:      err,
			}

			close(apiOrderResponse)
		}(orderInfo, apiCreateOrderResponses[i])

		i++
	}

	for index := range apiCreateOrderResponses {
		channelResponse := apiCreateOrderResponses[index]
		orderResponse := <-channelResponse
		if orderResponse.Error != nil {
			return 400, nil, orderResponse.Error
		}

		eventsCreated = append(eventsCreated, *orderResponse.Order)
	}

	// TODO: Implement rollback after failure
	// if createOrderResp.Error != nil {
	// 	return 500, nil, createOrderResp.Error
	// }

	return 200, eventsCreated, nil
}
