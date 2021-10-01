package asset

import (
	"stockfyApi/entity"
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
