package order

import (
	"errors"
	"stockfyApi/entity"
	"strings"
)

type MockApplication struct {
}

func NewMockApplication() *MockApplication {
	return &MockApplication{}
}

func (a *MockApplication) CreateOrder(quantity float64, price float64,
	currency string, orderType string, date string, brokerageId string,
	assetId string, userUid string) (*entity.Order, error) {

	dateFormatted := entity.StringToTime(date)
	orderFormatted, err := entity.NewOrder(quantity, price, currency, orderType,
		dateFormatted, brokerageId, assetId, userUid)
	if err != nil {
		return nil, err
	}

	return orderFormatted, nil
}

func (a *MockApplication) DeleteOrdersFromAsset(assetId string) ([]entity.Order,
	error) {

	if assetId == "ERROR_REPOSITORY" {
		return nil, errors.New("Unknown orders repository error")
	}

	return []entity.Order{
		{
			Id: "TestOrderID1",
		},
		{
			Id: "TestOrderID1",
		},
	}, nil
}

func (a *MockApplication) DeleteOrdersFromAssetUser(assetId string, userUid string) (
	*[]entity.Order, error) {
	if assetId == "ERROR_REPOSITORY" || userUid == "ERROR_REPOSITORY" {
		return nil, errors.New("Unknown orders repository error")
	}

	return &[]entity.Order{
		{
			Id: "TestOrderID1",
			Asset: &entity.Asset{
				Id:     assetId,
				Symbol: "TEST3",
			},
		},
		{
			Id: "TestOrderID2",
			Asset: &entity.Asset{
				Id:     assetId,
				Symbol: "TEST3",
			},
		},
	}, nil
}

func (a *MockApplication) DeleteOrdersFromUser(orderId string, userUid string) (
	*string, error) {
	if orderId == "ERROR_REPOSITORY" || userUid == "ERROR_REPOSITORY" {
		return nil, errors.New("Unknown orders repository error")
	}

	if orderId == "INVALID_ORDER_ID" {
		return nil, nil
	}

	return &orderId, nil
}

func (a *MockApplication) SearchOrderByIdAndUserUid(orderId string, userUid string) (
	*entity.Order, error) {

	dateFormatted := entity.StringToTime("2021-10-01")

	if orderId == "ERROR_REPOSITORY" || userUid == "ERROR_REPOSITORY" {
		return nil, errors.New("Unknown orders repository error")
	}

	if orderId == "INVALID_ORDER_ID" {
		return nil, nil
	}

	return &entity.Order{
		Id:        orderId,
		Price:     29.29,
		Quantity:  10,
		Currency:  "BRL",
		OrderType: "Dividendos",
		Date:      dateFormatted,
		Brokerage: &entity.Brokerage{
			Id:      "TestBrokerageID",
			Name:    "Test Brokerage",
			Country: "BR",
		},
		Asset: &entity.Asset{
			Id: "TestAssetID",
		},
	}, nil
}

func (a *MockApplication) SearchOrdersFromAssetUser(assetId string, userUid string) (
	[]entity.Order, error) {

	dateFormatted := entity.StringToTime("2021-10-01")

	if assetId == "ERROR_REPOSITORY" || userUid == "ERROR_REPOSITORY" {
		return nil, errors.New("Unknown orders repository error")
	}

	return []entity.Order{
		{
			Id:        "TestOrderID",
			Price:     29.29,
			Quantity:  10,
			Currency:  "BRL",
			OrderType: "Dividendos",
			Date:      dateFormatted,
			Brokerage: &entity.Brokerage{
				Id:      "TestBrokerageID",
				Name:    "Test Brokerage",
				Country: "BR",
			},
		},
	}, nil
}

func (a *MockApplication) SearchOrdersSearchFromAssetUserByDate(assetId string,
	userUid string, orderBy string, limit int, offset int) ([]entity.Order,
	error) {

	lowerOrderBy := strings.ToLower(orderBy)
	if lowerOrderBy != "asc" && lowerOrderBy != "desc" {
		return nil, entity.ErrInvalidOrderOrderBy
	}

	if assetId == "UNKNOWN_ORDER_REPOSITORY_ERROR" {
		return nil, errors.New("Unknown error in the order repository")
	}

	dateFormatted := entity.StringToTime("2021-10-01")
	return []entity.Order{
		{
			Id:        "Order1",
			Quantity:  2,
			Price:     29.29,
			Currency:  "BRL",
			OrderType: "buy",
			Date:      dateFormatted,
			Brokerage: &entity.Brokerage{
				Id:      "TestBrokerageID",
				Name:    "Test",
				Country: "BR",
			},
		},
		{
			Id:        "Order2",
			Quantity:  3,
			Price:     31.90,
			Currency:  "BRL",
			OrderType: "buy",
			Date:      dateFormatted,
			Brokerage: &entity.Brokerage{
				Id:      "TestBrokerageID",
				Name:    "Test",
				Country: "BR",
			},
		},
	}, nil
}

func (a *MockApplication) UpdateOrder(orderId string, userUid string, price float64,
	quantity float64, orderType, date string, brokerageId string,
	currency string) (*entity.Order, error) {

	dateFormatted := entity.StringToTime(date)

	orderFormatted, err := entity.NewOrder(quantity, price, currency,
		orderType, dateFormatted, brokerageId, "", userUid)
	if err != nil {
		return nil, err
	}
	orderFormatted.Id = orderId

	return &entity.Order{
		Id:        "TestOrderID",
		Price:     29.29,
		Quantity:  10,
		Currency:  "BRL",
		OrderType: "Dividendos",
		Date:      dateFormatted,
		Brokerage: &entity.Brokerage{
			Id:      "TestBrokerageID",
			Name:    "Test Brokerage",
			Country: "BR",
		},
	}, nil
}

func (a *MockApplication) OrderVerification(orderType string, country string,
	quantity float64, price float64, currency string) error {

	if orderType != "sell" && orderType != "buy" {
		return entity.ErrInvalidOrderType
	}

	if country != "BR" && country != "US" {
		return entity.ErrInvalidCountryCode
	}

	if country == "BR" && (orderType == "sell" || orderType == "buy") {
		if !entity.IsIntegral(quantity) {
			return entity.ErrInvalidOrderQuantityBrazil
		}
	}

	if country == "BR" && currency != "BRL" {
		return entity.ErrInvalidBrazilCurrency
	}

	if country == "US" && currency != "USD" {
		return entity.ErrInvalidUsaCurrency
	}

	if orderType == "buy" && quantity < 0 {
		return entity.ErrInvalidOrderBuyQuantity
	} else if orderType == "sell" && quantity > 0 {
		return entity.ErrInvalidOrderSellQuantity
	} else if price < 0 {
		return entity.ErrInvalidOrderPrice
	}

	return nil
}
