package asset

import (
	"errors"
	"stockfyApi/entity"
	externalapi "stockfyApi/externalApi"
	assettype "stockfyApi/usecases/assetType"
	"stockfyApi/usecases/general"
	"strings"
)

type MockApplication struct {
}

func NewMockApplication() *MockApplication {
	return &MockApplication{}
}

func (a *MockApplication) CreateAsset(symbol string, fullname string,
	preference *string, sectorId string, assetType assettype.AssetType) (
	entity.Asset, error) {

	assetInfo, err := entity.NewAsset(symbol, fullname, preference, sectorId,
		assetType.Id, assetType.Type, assetType.Country)

	return *assetInfo, err
}

func (a *MockApplication) SearchAsset(symbol string) (*entity.Asset, error) {

	preference := "TestPref"

	switch symbol {
	case "ERROR_REPOSITORY":
		return nil, errors.New("Unknown repository error")
	case "UNKNOWN_SYMBOL":
		return nil, nil
	default:
		return &entity.Asset{
			Id:         "TestID",
			Symbol:     symbol,
			Preference: &preference,
			Fullname:   "Test Name",
			AssetType: &entity.AssetType{
				Id:      "TestAssetTypeID",
				Type:    "ETF",
				Name:    "Test ETF",
				Country: "BR",
			},
			Sector: &entity.Sector{
				Id:   "TestSectorID",
				Name: "Test Sector",
			},
		}, nil
	}

}

func (a *MockApplication) DeleteAsset(assetId string) (*entity.Asset, error) {
	preference := "ON"

	switch assetId {
	case "ERROR_REPOSITORY":
		return nil, errors.New("Unknown repository error")
	case "UNKNOWN_ID":
		return nil, nil
	default:
		return &entity.Asset{
			Id:         assetId,
			Symbol:     "TEST11",
			Preference: &preference,
			Fullname:   "Test Name",
		}, nil
	}
}

func (a *MockApplication) SearchAssetByUser(symbol string, userUid string,
	withOrders bool, withOrderResume bool) (*entity.Asset, error) {

	orderType := ""
	preference := "TestPref"
	dateFormatted := entity.StringToTime("2021-10-01")

	if symbol == "" || userUid == "" || symbol == "UNKNOWN_SYMBOL" ||
		userUid == "UNKNOWN_UID" {
		return nil, entity.ErrInvalidAssetSymbolUserRelation
	} else if symbol == "ERROR_REPOSITORY" {
		return nil, errors.New("Unknown repository error")
	}

	if withOrders && !withOrderResume {
		orderType = "ONLYORDERS"
	} else if withOrders && withOrderResume {
		orderType = "ALL"
	} else if !withOrders && withOrderResume {
		orderType = "ONLYINFO"
	}

	switch orderType {
	case "ALL":
		return &entity.Asset{
			Id:         "TestID",
			Symbol:     symbol,
			Preference: &preference,
			Fullname:   "Test Name",
			AssetType: &entity.AssetType{
				Id:      "TestAssetTypeID",
				Type:    "ETF",
				Name:    "Test ETF",
				Country: "BR",
			},
			Sector: &entity.Sector{
				Id:   "TestSectorID",
				Name: "Test Sector",
			},
			OrdersList: []entity.Order{
				{
					Id:        "Order1",
					Quantity:  2,
					Price:     29.29,
					Currency:  "USD",
					OrderType: "Dividendos",
					Date:      dateFormatted,
					Brokerage: &entity.Brokerage{
						Id:      "BrokerageID",
						Name:    "Test Broker",
						Country: "US",
					},
				},
				{
					Id:        "Order2",
					Quantity:  2,
					Price:     29.29,
					Currency:  "USD",
					OrderType: "Dividendos",
					Date:      dateFormatted,
					Brokerage: &entity.Brokerage{
						Id:      "BrokerageID",
						Name:    "Test Broker",
						Country: "US",
					},
				},
			},
			OrderInfo: &entity.OrderInfos{
				TotalQuantity:        4,
				WeightedAdjPrice:     28.20,
				WeightedAveragePrice: 29.29,
			},
		}, nil
	case "ONLYORDERS":
		return &entity.Asset{
			Id:         "TestID",
			Symbol:     symbol,
			Preference: &preference,
			Fullname:   "Test Name",
			AssetType: &entity.AssetType{
				Id:      "TestAssetTypeID",
				Type:    "ETF",
				Name:    "Test ETF",
				Country: "BR",
			},
			Sector: &entity.Sector{
				Id:   "TestSectorID",
				Name: "Test Sector",
			},
			OrdersList: []entity.Order{
				{
					Id:        "Order1",
					Quantity:  2,
					Price:     29.29,
					Currency:  "USD",
					OrderType: "Dividendos",
					Date:      dateFormatted,
					Brokerage: &entity.Brokerage{
						Id:      "BrokerageID",
						Name:    "Test Broker",
						Country: "US",
					},
				},
				{
					Id:        "Order2",
					Quantity:  2,
					Price:     29.29,
					Currency:  "USD",
					OrderType: "Dividendos",
					Date:      dateFormatted,
					Brokerage: &entity.Brokerage{
						Id:      "BrokerageID",
						Name:    "Test Broker",
						Country: "US",
					},
				},
			},
		}, nil
	case "ONLYINFO":
		return &entity.Asset{
			Id:         "TestID",
			Symbol:     symbol,
			Preference: &preference,
			Fullname:   "Test Name",
			AssetType: &entity.AssetType{
				Id:      "TestAssetTypeID",
				Type:    "ETF",
				Name:    "Test ETF",
				Country: "BR",
			},
			Sector: &entity.Sector{
				Id:   "TestSectorID",
				Name: "Test Sector",
			},
			OrderInfo: &entity.OrderInfos{
				TotalQuantity:        4,
				WeightedAdjPrice:     28.20,
				WeightedAveragePrice: 29.29,
			},
		}, nil
	default:
		return &entity.Asset{
			Id:         "TestID",
			Symbol:     symbol,
			Preference: &preference,
			Fullname:   "Test Name",
			AssetType: &entity.AssetType{
				Id:      "TestAssetTypeID",
				Type:    "ETF",
				Name:    "Test ETF",
				Country: "BR",
			},
			Sector: &entity.Sector{
				Id:   "TestSectorID",
				Name: "Test Sector",
			},
		}, nil
	}
}

func (a *MockApplication) SearchAssetPerAssetType(assetType string, country string,
	userUid string, withOrdersInfo bool) (*entity.AssetType, error) {

	err := general.CountryValidation(country)
	if err != nil {
		return nil, err
	}

	err = general.AssetTypeNameValidation(assetType)
	if err != nil {
		return nil, err
	}

	if userUid == "NO_ASSET_PER_ASSETTYPE" {
		return nil, entity.ErrInvalidAssetType
	}

	switch withOrdersInfo {
	case false:
		return &entity.AssetType{
			Id:      "TestAssetTypeID",
			Type:    assetType,
			Country: "US",
			Name:    "Test " + assetType,
			Assets: []entity.Asset{
				{
					Id:         "TestID1",
					Symbol:     "Test1",
					Preference: nil,
					Fullname:   "Test Name1",
					Sector: &entity.Sector{
						Id:   "TestSectorID",
						Name: "Test Sector",
					},
				},
				{
					Id:         "TestID2",
					Symbol:     "Test2",
					Preference: nil,
					Fullname:   "Test Name2",
					Sector: &entity.Sector{
						Id:   "TestSectorID",
						Name: "Test Sector",
					},
				},
			},
		}, nil
	default:
		return &entity.AssetType{
			Id:      "TestAssetTypeID",
			Type:    assetType,
			Country: "US",
			Name:    "Test " + assetType,
			Assets: []entity.Asset{
				{
					Id:         "TestID1",
					Symbol:     "Test1",
					Preference: nil,
					Fullname:   "Test Name1",
					Sector: &entity.Sector{
						Id:   "TestSectorID",
						Name: "Test Sector",
					},
					OrderInfo: &entity.OrderInfos{
						TotalQuantity:        4,
						WeightedAdjPrice:     28.20,
						WeightedAveragePrice: 29.29,
					},
				},
				{
					Id:         "TestID2",
					Symbol:     "Test2",
					Preference: nil,
					Fullname:   "Test Name2",
					Sector: &entity.Sector{
						Id:   "TestSectorID",
						Name: "Test Sector",
					},
					OrderInfo: &entity.OrderInfos{
						TotalQuantity:        4,
						WeightedAdjPrice:     28.20,
						WeightedAveragePrice: 29.29,
					},
				},
			},
		}, nil
	}
}

func (a *MockApplication) AssetPreferenceType(symbol string, country string,
	assetType string) string {
	var preference string

	if country == "BR" && assetType == "STOCK" {
		switch symbol[len(symbol)-1:] {
		case "3":
			preference = "ON"
			break
		case "4":
			preference = "PN"
			break
		case "1":
			preference = "UNIT"
			break
		default:
			preference = ""
			break
		}
	}

	return preference
}

func (a *MockApplication) AssetVerificationExistence(symbol string, country string,
	extApi externalapi.ThirdPartyInterfaces) (*entity.SymbolLookup, error) {
	if symbol == "" {
		return nil, entity.ErrInvalidApiQuerySymbolBlank
	}

	if err := general.CountryValidation(country); err != nil {
		return nil, err
	}

	if country == "BR" {
		symbol = symbol + ".SA"
	}

	switch symbol {
	case "UNKNOWN_SYMBOL":
		return nil, entity.ErrInvalidAssetSymbol
	case "UNKNOWN_SYMBOL.SA":
		return nil, entity.ErrInvalidAssetSymbol
	default:
		return &entity.SymbolLookup{
			Fullname: "Test Name",
			Symbol:   strings.ReplaceAll(symbol, ".SA", ""),
			Type:     "ETP",
		}, nil
	}
}

func (a *MockApplication) AssetVerificationSector(assetType string, symbol string,
	country string, extInterface ExternalApiRepository) string {
	if country == "BR" {
		symbol = symbol + ".SA"
	}

	if assetType == "STOCK" {
		return "TestSector"
	} else if assetType == "ETF" {
		return "Blend"
	} else {
		return "Real Estate"
	}
}

func (a *MockApplication) AssetVerificationPrice(symbol string, country string,
	extInterface externalapi.ThirdPartyInterfaces) (*entity.SymbolPrice, error) {

	if err := general.CountryValidation(country); err != nil {
		return nil, err
	}

	if symbol == "" {
		return nil, entity.ErrInvalidApiQuerySymbolBlank
	}

	if country == "BR" {
		symbol = symbol + ".SA"
	}

	if symbol == "UNKNOWN_SYMBOL" || symbol == "UNKNOWN_SYMBOL.SA" {
		return nil, entity.ErrInvalidAssetSymbol
	} else {
		return &entity.SymbolPrice{
			Symbol:         strings.ReplaceAll(symbol, ".SA", ""),
			CurrentPrice:   29.29,
			LowPrice:       28.00,
			HighPrice:      29.89,
			OpenPrice:      29.29,
			PrevClosePrice: 29.29,
			MarketCap:      1018388,
		}, nil
	}
}
