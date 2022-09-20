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

func (a *MockApplication) MeasureAssetTotalQuantityForSpecificDate(
	assetId string, userUid string, date string) (map[string]float64, error) {

	var orders []entity.Order

	ordersQuantityByBrokerage := make(map[string]float64)

	// dateFormatted := entity.StringToTime(date)
	// orders, err := a.repo.SearchFromAssetUserSpecificDate(assetId, userUid,
	// 	dateFormatted)
	// if err != nil {
	// 	return nil, err
	// }

	if assetId != "EmptyOrders" {
		orders = []entity.Order{
			{
				Id:        "Order1",
				Quantity:  2,
				Price:     29.29,
				Currency:  "BRL",
				OrderType: "buy",
				Date:      entity.StringToTime("2021-10-01"),
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
				Date:      entity.StringToTime("2021-11-01"),
				Brokerage: &entity.Brokerage{
					Id:      "TestBrokerageID",
					Name:    "Test",
					Country: "BR",
				},
			},
			{
				Id:        "Order3",
				Quantity:  45,
				Price:     30.90,
				Currency:  "BRL",
				OrderType: "buy",
				Date:      entity.StringToTime("2021-12-01"),
				Brokerage: &entity.Brokerage{
					Id:      "TestBrokerageID",
					Name:    "Test",
					Country: "BR",
				},
			},
			{
				Id:        "Order4",
				Quantity:  60,
				Price:     27.90,
				Currency:  "BRL",
				OrderType: "buy",
				Date:      entity.StringToTime("2021-08-01"),
				Brokerage: &entity.Brokerage{
					Id:      "TestBrokerageID2",
					Name:    "Test2",
					Country: "BR",
				},
			},
			{
				Id:        "Order5",
				Quantity:  14,
				Price:     27.90,
				Currency:  "BRL",
				OrderType: "buy",
				Date:      entity.StringToTime("2021-08-06"),
				Brokerage: &entity.Brokerage{
					Id:      "TestBrokerageID2",
					Name:    "Test2",
					Country: "BR",
				},
			},
		}
	} else if assetId == "ErrorQuery" {
		return nil, errors.New("Unknown error in the order repository")
	}

	if orders == nil {
		return nil, entity.ErrEmptyQuery
	}

	for _, order := range orders {
		ordersQuantityByBrokerage[order.Brokerage.Name] += order.Quantity
	}

	return ordersQuantityByBrokerage, nil
}

func (a *MockApplication) OrderVerification(orderType string, country string,
	quantity float64, price float64, currency string) error {

	if !entity.ValidOrderType[orderType] {
		return entity.ErrInvalidOrderType
	}

	if country != "BR" && country != "US" {
		return entity.ErrInvalidCountryCode
	}

	if country == "BR" && orderType == "buy" {
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

	// Split Price Verification
	if orderType == "split" && price != 0 {
		return entity.ErrInvalideNonZeroOrderPrice
	}

	// Bonification Price Verification
	if (orderType == "bonification" || orderType == "sell" ||
		orderType == "buy") && price <= 0 {
		return entity.ErrInvalidNegativeOrderPrice
	}

	// Demerger Price Verification
	if orderType == "demerge" && price >= 0 {
		return entity.ErrInvalidPositiveOrderPrice
	}

	// Quantity Verification
	if (orderType == "buy" || orderType == "bonification" ||
		orderType == "split") && quantity <= 0 {
		return entity.ErrInvalidOrderBuyQuantity
	} else if orderType == "sell" && quantity > 0 {
		return entity.ErrInvalidOrderSellQuantity
	} else if orderType == "demerge" && quantity != 0 {
		return entity.ErrInvalidOrderDemergeQuantity
	}

	return nil
}

func (a *MockApplication) EventTypeValueVerification(eventType string) error {
	if eventType == "bonification" || eventType == "split" {
		return nil
	} else {
		return entity.ErrInvalidEventType
	}
}
