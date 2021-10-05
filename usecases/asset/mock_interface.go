package asset

import (
	"stockfyApi/entity"
	"time"
)

type MockDb struct {
}

type MockExternal struct {
}

func NewMockRepo() *MockDb {
	return &MockDb{}
}

func NewExternalApi() *MockExternal {
	return &MockExternal{}
}

func (m *MockDb) Create(assetInsert entity.Asset) entity.Asset {

	assetCreated := entity.Asset{
		Id:         "a38a9jkrh40a",
		Symbol:     assetInsert.Symbol,
		Preference: assetInsert.Preference,
		Fullname:   assetInsert.Fullname,
	}

	return assetCreated
}

func (m *MockDb) Search(symbol string) ([]entity.Asset, error) {

	if symbol == "Invalid" {
		return nil, entity.ErrInvalidSearchAssetName
	}

	assetType := entity.AssetType{
		Id:      "28ccf27a-ed8b-11eb-9a03-0242ac130003",
		Type:    "STOCK",
		Name:    "Ações Brasil",
		Country: "BR",
	}
	preference := "ON"

	sectorInfo := entity.Sector{
		Id:   "83ae92f8-ed8b-11eb-9a03-0242ac130003",
		Name: "Finance",
	}

	return []entity.Asset{
		{
			Id:         "0a52d206-ed8b-11eb-9a03-0242ac130003",
			Symbol:     symbol,
			Preference: &preference,
			Fullname:   "Itau Unibanco Holding SA",
			AssetType:  &assetType,
			Sector:     &sectorInfo,
		},
	}, nil
}

func (m *MockDb) SearchByUser(symbol string, userUid string, orderType string) (
	[]entity.Asset, error) {

	if symbol == "Invalid" {
		return nil, entity.ErrInvalidSearchAssetName
	}

	tr, _ := time.Parse("2021-07-05", "2021-07-21")
	tr2, _ := time.Parse("2021-07-05", "2020-04-02")

	assetType := entity.AssetType{
		Id:      "28ccf27a-ed8b-11eb-9a03-0242ac130003",
		Type:    "STOCK",
		Name:    "Ações Brasil",
		Country: "BR",
	}

	preference := "ON"

	brokerageInfo := entity.Brokerage{
		Id:      "55555555-ed8b-11eb-9a03-0242ac130003",
		Name:    "Clear",
		Country: "BR",
	}

	orderList := []entity.Order{
		{
			Id:        "44444444-ed8b-11eb-9a03-0242ac130003",
			Quantity:  20,
			Price:     39.93,
			Currency:  "BRL",
			OrderType: "buy",
			Date:      tr,
			Brokerage: &brokerageInfo,
		},
		{
			Id:        "yeid847e-ed8b-11eb-9a03-0242ac130003",
			Quantity:  5,
			Price:     27.13,
			Currency:  "BRL",
			OrderType: "buy",
			Date:      tr2,
			Brokerage: &brokerageInfo,
		},
	}

	sectorInfo := entity.Sector{
		Id:   "83ae92f8-ed8b-11eb-9a03-0242ac130003",
		Name: "Finance",
	}

	ordersInfo := entity.OrderInfos{
		TotalQuantity:        25,
		WeightedAdjPrice:     37.37,
		WeightedAveragePrice: 37.37,
	}

	if orderType == "ONLYORDERS" {
		return []entity.Asset{
			{
				Id:         "0a52d206-ed8b-11eb-9a03-0242ac130003",
				Symbol:     symbol,
				Preference: &preference,
				Fullname:   "Itau Unibanco Holding SA",
				AssetType:  &assetType,
				Sector:     &sectorInfo,
				OrdersList: orderList,
			},
		}, nil
	}

	if orderType == "ONLYINFO" {
		return []entity.Asset{
			{
				Id:         "0a52d206-ed8b-11eb-9a03-0242ac130003",
				Symbol:     symbol,
				Preference: &preference,
				Fullname:   "Itau Unibanco Holding SA",
				AssetType:  &assetType,
				Sector:     &sectorInfo,
				OrderInfo:  &ordersInfo,
			},
		}, nil
	}

	if orderType == "ALL" {
		return []entity.Asset{
			{
				Id:         "0a52d206-ed8b-11eb-9a03-0242ac130003",
				Symbol:     symbol,
				Preference: &preference,
				Fullname:   "Itau Unibanco Holding SA",
				AssetType:  &assetType,
				Sector:     &sectorInfo,
				OrdersList: orderList,
				OrderInfo:  &ordersInfo,
			},
		}, nil
	}

	return []entity.Asset{}, nil
}

func (m *MockExternal) CompanyOverview(symbol string) map[string]string {
	return map[string]string{
		"country":         "BR",
		"currency":        "BRL",
		"exchange":        "Sao Paolo",
		"finnhubIndustry": "Banking",
		"ipo":             "",
		"logo":            "https://test.com",
		"name":            "Bradesco S.A",
		"phone":           "918480347",
		"ticker":          "BBDC3",
		"weburl":          "https://test.com",
	}
}

func (m *MockExternal) GetPrice(symbol string) entity.SymbolPrice {
	if symbol != "ITUB3.SA" {
		return entity.SymbolPrice{}
	}

	return entity.SymbolPrice{
		Symbol:         "ITUB3",
		CurrentPrice:   29.93,
		HighPrice:      31.00,
		LowPrice:       29.56,
		OpenPrice:      30.99,
		PrevClosePrice: 30.99,
		MarketCap:      1478481948,
	}
}

func (m *MockExternal) VerifySymbol2(symbol string) entity.SymbolLookup {
	if symbol != "ITUB4.SA" {
		return entity.SymbolLookup{}
	} else {
		return entity.SymbolLookup{
			Fullname: "Itau Unibanco Holding SA",
			Symbol:   "ITUB4",
			Type:     "STOCK",
		}
	}

}
