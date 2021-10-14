package brokerage

import (
	"stockfyApi/entity"
)

type MockDb struct {
}

func NewMockRepo() *MockDb {
	return &MockDb{}
}

func (m *MockDb) Search(specificFetch string, args ...string) (
	[]entity.Brokerage, error) {

	switch specificFetch {
	case "ALL":
		return []entity.Brokerage{
			{
				Id:      "55555555-ed8b-11eb-9a03-0242ac130003",
				Name:    "Clear",
				Country: "BR",
			},
			{
				Id:      "55556666-ed8b-11eb-9a03-0242ac130003",
				Name:    "Rico",
				Country: "BR",
			},
			{
				Id:      "15151515-ed8b-11eb-9a03-0242ac130003",
				Name:    "Avenue",
				Country: "US",
			},
		}, nil
	case "SINGLE":
		if args[0] == "Invalid" {
			return nil, nil
		} else {
			return []entity.Brokerage{
				{
					Id:      "55556666-ed8b-11eb-9a03-0242ac130003",
					Name:    args[0],
					Country: "BR",
				},
			}, nil
		}
	case "COUNTRY":
		return []entity.Brokerage{
			{
				Id:      "55555555-ed8b-11eb-9a03-0242ac130003",
				Name:    "Clear",
				Country: args[0],
			},
			{
				Id:      "55556666-ed8b-11eb-9a03-0242ac130003",
				Name:    "Rico",
				Country: args[0],
			},
		}, nil
	default:
		return nil, entity.ErrInvalidBrokerageSearchType
	}
}
