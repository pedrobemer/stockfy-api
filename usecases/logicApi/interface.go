package logicApi

import "stockfyApi/entity"

type UseCases interface {
	ApiAssetVerification(symbol string, country string) (int, *entity.Asset,
		error)
	ApiCreateOrder(symbol string, country string, orderType string,
		quantity float64, price float64, currency string, brokerage string,
		date string, userUid string) (int, *entity.Order, error)
	ApiAssetsPerAssetType(assetType string, country string, ordersInfo bool,
		withPrice bool, userUid string) (int, *entity.AssetType, error)
	ApiDeleteAssets(myUser bool, userUid string, symbol string) (int,
		*entity.Asset, error)
	ApiGetOrdersFromAssetUser(symbol string, userUid string, orderBy string,
		limit string, offset string) (int, []entity.Order, error)
	ApiUpdateOrdersFromUser(orderId string, userUid string, orderType string,
		price float64, quantity float64, date string, brokerage string) (int,
		*entity.Order, error)
	ApiCreateEarnings(symbol string, currency string, earningType string,
		date string, earnings float64, userUid string) (int, *entity.Earnings,
		error)
	ApiGetEarningsFromAssetUser(symbol string, userUid string, orderBy string,
		limit string, offset string) (int, []entity.Earnings, error)
	ApiUpdateEarningsFromUser(earningId string, earning float64,
		earningType string, date string, userUid string) (int, *entity.Earnings,
		error)
	ApiGetAssetByUser(symbol string, userUid string, withOrders bool,
		withOrderResume bool, withPrice bool) (int, *entity.Asset, error)
	ApiCreateEvent(symbol string, symbolDemerger string,
		orderType string, eventRate float64, price float64, currency string,
		date string, userUid string) (int, []entity.Order, error)
}
