package order

import (
	"stockfyApi/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrder(t *testing.T) {
	layOut := "2006-01-02"
	dateFormatted, _ := time.Parse(layOut, "2021-10-04")

	brokerage := entity.Brokerage{
		Id:      "BrokerageID",
		Name:    "Test Broker",
		Country: "US",
	}

	expectedOrderCreated := entity.Order{
		Quantity:  3.4,
		Price:     221.38,
		Currency:  "USD",
		OrderType: "buy",
		Date:      dateFormatted,
		Brokerage: &brokerage,
		Asset: &entity.Asset{
			Id: "AssetID",
		},
	}

	mocked := NewMockRepo()
	app := NewApplication(mocked)

	orderCreated, err := app.CreateOrder(3.4, 221.38, "USD", "buy", "2021-10-04",
		"BrokerageID", "AssetID", "userUID")

	assert.Equal(t, expectedOrderCreated, *orderCreated)
	assert.Nil(t, err)

}

func TestDeleteOrdersFromUser(t *testing.T) {
	type test struct {
		orderId         string
		userUid         string
		expectedOrderId *string
		expectedError   error
	}

	validID := "TestID"
	tests := []test{
		{
			orderId:         validID,
			userUid:         "UserUid",
			expectedOrderId: &validID,
			expectedError:   nil,
		},
		{
			orderId:         "INVALID_ID",
			userUid:         "UserUid",
			expectedOrderId: nil,
			expectedError:   nil,
		},
	}

	mocked := NewMockRepo()
	app := NewApplication(mocked)

	for _, testCase := range tests {
		deletedOrderId, err := app.DeleteOrdersFromUser(testCase.orderId,
			testCase.userUid)
		assert.Equal(t, testCase.expectedOrderId, deletedOrderId)
		assert.Equal(t, testCase.expectedError, err)
	}
}

func TestSearchOrderByIdAndUserUid(t *testing.T) {
	type test struct {
		orderId           string
		userUid           string
		expectedOrderInfo *entity.Order
		expectedError     error
	}

	layOut := "2006-01-02"
	dateFormatted, _ := time.Parse(layOut, "2021-10-04")

	tests := []test{
		{
			orderId: "ValidID",
			expectedOrderInfo: &entity.Order{
				Id:        "ValidID",
				Quantity:  20,
				Price:     2.49,
				Currency:  "BRL",
				OrderType: "buy",
				Date:      dateFormatted,
				Brokerage: &entity.Brokerage{
					Id:      "BrokerageID",
					Name:    "Broker",
					Country: "BR",
				},
				Asset: &entity.Asset{
					Id: "AssetID",
				},
			},
			expectedError: nil,
		},
		{
			orderId:           "INVALID_ORDER",
			expectedOrderInfo: nil,
			expectedError:     nil,
		},
	}

	mocked := NewMockRepo()
	app := NewApplication(mocked)

	for _, testCase := range tests {
		orderInfo, err := app.SearchOrderByIdAndUserUid(testCase.orderId,
			testCase.userUid)
		assert.Equal(t, testCase.expectedOrderInfo, orderInfo)
		assert.Equal(t, testCase.expectedError, err)
	}
}

func TestOrderVerification(t *testing.T) {

	type test struct {
		orderType     string
		country       string
		quantity      float64
		price         float64
		currency      string
		expectedError error
	}

	tests := []test{
		{
			orderType:     "sell",
			country:       "BR",
			quantity:      -20,
			price:         10.92,
			currency:      "BRL",
			expectedError: nil,
		},
		{
			orderType:     "buy",
			country:       "BR",
			quantity:      20,
			price:         10.92,
			currency:      "BRL",
			expectedError: nil,
		},
		{
			orderType:     "ai4a9",
			country:       "BR",
			quantity:      20,
			price:         10.92,
			currency:      "BRL",
			expectedError: entity.ErrInvalidOrderType,
		},
		{
			orderType:     "buy",
			country:       "AUS",
			quantity:      20.35,
			price:         10.92,
			currency:      "BRL",
			expectedError: entity.ErrInvalidCountryCode,
		},
		{
			orderType:     "buy",
			country:       "BR",
			quantity:      20.35,
			price:         10.92,
			currency:      "BRL",
			expectedError: entity.ErrInvalidOrderQuantityBrazil,
		},
		{
			orderType:     "sell",
			country:       "BR",
			quantity:      -20.35,
			price:         10.92,
			currency:      "BRL",
			expectedError: entity.ErrInvalidOrderQuantityBrazil,
		},
		{
			orderType:     "sell",
			country:       "US",
			quantity:      -20.35,
			price:         10.92,
			currency:      "BRL",
			expectedError: entity.ErrInvalidUsaCurrency,
		},
		{
			orderType:     "sell",
			country:       "BR",
			quantity:      -20,
			price:         10.92,
			currency:      "USD",
			expectedError: entity.ErrInvalidBrazilCurrency,
		},
		{
			orderType:     "buy",
			country:       "BR",
			quantity:      -1,
			price:         10.92,
			currency:      "BRL",
			expectedError: entity.ErrInvalidOrderBuyQuantity,
		},
		{
			orderType:     "sell",
			country:       "BR",
			quantity:      1,
			price:         10.92,
			currency:      "BRL",
			expectedError: entity.ErrInvalidOrderSellQuantity,
		},
		{
			orderType:     "sell",
			country:       "BR",
			quantity:      -1,
			price:         -10.92,
			currency:      "BRL",
			expectedError: entity.ErrInvalidOrderPrice,
		},
	}

	mocked := NewMockRepo()
	app := NewApplication(mocked)

	for _, testCase := range tests {
		err := app.OrderVerification(testCase.orderType, testCase.country,
			testCase.quantity, testCase.price, testCase.currency)
		assert.Equal(t, testCase.expectedError, err)

	}
}
