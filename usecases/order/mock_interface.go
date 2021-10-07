package order

import (
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
		return "", nil
	}

	return id, nil
}

func (m *MockDb) SearchFromAssetUser(assetId string, userUid string) (
	[]entity.Order, error) {
	return []entity.Order{}, nil
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
