package order

import (
	"errors"
	"stockfyApi/entity"
	"strings"
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

func (a *Application) CreateOrder(quantity float64, price float64,
	currency string, orderType string, date string, brokerageId string,
	assetId string, userUid string) (*entity.Order, error) {

	dateFormatted := entity.StringToTime(date)
	orderFormatted, err := entity.NewOrder(quantity, price, currency, orderType,
		dateFormatted, brokerageId, assetId, userUid)
	if err != nil {
		return nil, err
	}

	orderCreated := a.repo.Create(*orderFormatted)

	return &orderCreated, nil
}

func (a *Application) DeleteOrdersFromAsset(assetId string) ([]entity.Order,
	error) {
	ordersDeleted, err := a.repo.DeleteFromAsset(assetId)
	if err != nil {
		return nil, err
	}

	return ordersDeleted, nil
}

func (a *Application) DeleteOrdersFromAssetUser(assetId string, userUid string) (
	*[]entity.Order, error) {
	ordersDeleted, err := a.repo.DeleteFromAssetUser(assetId, userUid)
	if err != nil {
		return nil, err
	}

	return &ordersDeleted, nil
}

func (a *Application) DeleteOrdersFromUser(orderId string, userUid string) (
	*string, error) {

	deletedOrderId, err := a.repo.DeleteFromUser(orderId, userUid)
	if err != nil {
		if err.Error() == errors.New("no rows in result set").Error() {
			return nil, nil
		}
		return nil, err
	}

	return &deletedOrderId, nil
}

func (a *Application) SearchOrderByIdAndUserUid(orderId string, userUid string) (
	*entity.Order, error) {

	orderInfo, err := a.repo.SearchByOrderAndUserId(orderId, userUid)
	if err != nil {
		return nil, err
	}

	if orderInfo == nil {
		return nil, nil
	}

	return &orderInfo[0], nil
}

func (a *Application) SearchOrdersFromAssetUser(assetId string, userUid string) (
	[]entity.Order, error) {
	assetInfo, err := a.repo.SearchFromAssetUser(assetId, userUid)
	if err != nil {
		return nil, err
	}

	return assetInfo, nil
}

func (a *Application) SearchOrdersSearchFromAssetUserByDate(assetId string,
	userUid string, orderBy string, limit int, offset int) ([]entity.Order,
	error) {

	lowerOrderBy := strings.ToLower(orderBy)
	if lowerOrderBy != "asc" && lowerOrderBy != "desc" {
		return nil, entity.ErrInvalidOrderOrderBy
	}

	assetOrderInfo, err := a.repo.SearchFromAssetUserOrderByDate(assetId, userUid,
		orderBy, limit, offset)
	if err != nil {
		return nil, err
	}

	return assetOrderInfo, nil
}

func (a *Application) UpdateOrder(orderId string, userUid string, price float64,
	quantity float64, orderType, date string, brokerageId string,
	currency string) (*entity.Order, error) {

	dateFormatted := entity.StringToTime(date)

	orderFormatted, err := entity.NewOrder(quantity, price, currency,
		orderType, dateFormatted, brokerageId, "", userUid)
	if err != nil {
		return nil, err
	}
	orderFormatted.Id = orderId

	updatedOrder := a.repo.UpdateFromUser(*orderFormatted)

	return &updatedOrder[0], nil
}

func (a *Application) MeasureAssetTotalQuantityForSpecificDate(assetId string,
	userUid string, date string) (map[string]float64, error) {

	ordersQuantityByBrokerage := make(map[string]float64)

	dateFormatted := entity.StringToTime(date)

	orders, err := a.repo.SearchFromAssetUserSpecificDate(assetId, userUid,
		dateFormatted)
	if err != nil {
		return nil, err
	}

	if orders == nil {
		return nil, entity.ErrEmptyQuery
	}

	for _, order := range orders {
		ordersQuantityByBrokerage[order.Brokerage.Name] += order.Quantity
	}

	return ordersQuantityByBrokerage, nil
}

func orderBuyVerification(quantity float64, price float64) error {

	if !entity.IsIntegral(quantity) {
		return entity.ErrInvalidOrderQuantityBrazil
	}

	if price <= 0 {
		return entity.ErrInvalidNegativeOrderPrice
	}

	if quantity <= 0 {
		return entity.ErrInvalidOrderBuyQuantity
	}

	return nil
}

func orderSellVerification(quantity float64, price float64) error {

	if price <= 0 {
		return entity.ErrInvalidNegativeOrderPrice
	}

	if quantity > 0 {
		return entity.ErrInvalidOrderSellQuantity
	}

	return nil
}

func orderBonificationVerification(quantity float64, price float64) error {

	if price <= 0 {
		return entity.ErrInvalidNegativeOrderPrice
	}

	if quantity <= 0 {
		return entity.ErrInvalidOrderBuyQuantity
	}

	return nil
}

func orderSplitVerification(quantity float64, price float64) error {

	if price != 0 {
		return entity.ErrInvalideNonZeroOrderPrice
	}

	if quantity <= 0 {
		return entity.ErrInvalidOrderBuyQuantity
	}

	return nil
}

func orderDemergeVerification(quantity float64, price float64) error {

	if price >= 0 {
		return entity.ErrInvalidPositiveOrderPrice
	}

	if quantity != 0 {
		return entity.ErrInvalidOrderDemergeQuantity
	}

	return nil
}

func (a *Application) OrderVerification(orderType string, country string,
	quantity float64, price float64, currency string) error {

	if !entity.ValidOrderType[orderType] {
		return entity.ErrInvalidOrderType
	}

	if country != "BR" && country != "US" {
		return entity.ErrInvalidCountryCode
	}

	if country == "BR" && currency != "BRL" {
		return entity.ErrInvalidBrazilCurrency
	}

	if country == "US" && currency != "USD" {
		return entity.ErrInvalidUsaCurrency
	}

	verificationError := map[string]func(float64, float64) error{
		"sell":         orderSellVerification,
		"buy":          orderBuyVerification,
		"bonification": orderBonificationVerification,
		"split":        orderSplitVerification,
		"demerge":      orderDemergeVerification,
	}[orderType](quantity, price)

	return verificationError
}

func (a *Application) EventTypeValueVerification(eventType string) error {

	if !entity.ValidEventsType[eventType] {
		return entity.ErrInvalidEventType
	}

	return nil
}
