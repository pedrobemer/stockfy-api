package order

import (
	"errors"
	"stockfyApi/entity"
	"time"
)

type MockDb struct {
}

func NewMockRepo() *MockDb {
	return &MockDb{}
}

func (m *MockDb) Create(orderInsert entity.Order) entity.Order {

	layOut := "2006-01-02"
	dateFormatted, _ := time.Parse(layOut, "2021-10-04")

	brokerage := entity.Brokerage{
		Id:      "BrokerageID",
		Name:    "Test Broker",
		Country: "US",
	}

	return entity.Order{
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
}

func (m *MockDb) DeleteFromAsset(symbolId string) ([]entity.Order, error) {
	return []entity.Order{}, nil
}

func (m *MockDb) DeleteFromAssetUser(assetId string, userUid string) (
	[]entity.Order, error) {
	return []entity.Order{}, nil
}

func (m *MockDb) DeleteFromUser(id string, userUid string) (string, error) {
	if id == "INVALID_ID" {
		return "", errors.New("no rows in result set")
	}

	return id, nil
}

func (m *MockDb) SearchFromAssetUser(assetId string, userUid string) (
	[]entity.Order, error) {
	return []entity.Order{}, nil
}

func (m *MockDb) SearchFromAssetUserOrderByDate(assetId string,
	userUid string, orderBy string, limit int, offset int) ([]entity.Order,
	error) {

	if assetId == "UNKNOWN_ID" || offset > 2 {
		return []entity.Order{}, nil
	}

	if assetId == "INVALID_ID" {
		return nil, errors.New("UUID SQL ERROR")
	}

	layOut := "2006-01-02"
	tr, _ := time.Parse(layOut, "2021-10-01")

	brokerage := entity.Brokerage{
		Id:      "55555555-ed8b-11eb-9a03-0242ac130003",
		Name:    "Test Brokerage",
		Country: "US",
	}

	return []entity.Order{
		{
			Id:        "a8a8a8a8-ed8b-11eb-9a03-0242ac130003",
			Quantity:  20,
			Price:     29.29,
			Currency:  "USD",
			OrderType: "buy",
			Date:      tr,
			Brokerage: &brokerage,
		},
		{
			Id:        "a9a999a9-ed8b-11eb-9a03-0242ac130003",
			Quantity:  198,
			Price:     20.00,
			Currency:  "USD",
			OrderType: "buy",
			Date:      tr,
			Brokerage: &brokerage,
		},
	}, nil
}

func (m *MockDb) SearchByOrderAndUserId(orderId string, userUid string) (
	[]entity.Order, error) {
	layOut := "2006-01-02"
	dateFormatted, _ := time.Parse(layOut, "2021-10-04")

	if orderId == "INVALID_ORDER" {
		return nil, nil
	}

	return []entity.Order{
		{
			Id:        orderId,
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
	}, nil
}

func (m *MockDb) UpdateFromUser(orderUpdate entity.Order) []entity.Order {
	return []entity.Order{}
}
