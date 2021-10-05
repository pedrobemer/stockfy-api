package order

import (
	"stockfyApi/entity"
	"time"
)

type MockDb struct {
}

// type MockExternal struct {
// }

func NewMockRepo() *MockDb {
	return &MockDb{}
}

// func NewExternalApi() *MockExternal {
// 	return &MockExternal{}
// }

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
