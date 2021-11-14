package logicApi

import (
	"errors"
	"stockfyApi/entity"
	"stockfyApi/usecases"
	"strconv"
	"strings"
	"time"
)

type MockApplication struct {
	app usecases.Applications
}

//NewApplication create new use case
func NewMockApplication(usecases usecases.Applications) *MockApplication {
	return &MockApplication{
		app: usecases,
	}
}

func (a *MockApplication) ApiAssetVerification(symbol string, country string) (
	int, *entity.Asset, error) {
	preference := "TestPref"

	if country != "BR" && country != "US" {
		return 400, nil, entity.ErrInvalidCountryCode
	}

	if symbol == "UNKNOWN_SYMBOL" {
		return 404, nil, entity.ErrInvalidAssetSymbol
	}

	if symbol == "ERROR_SECTOR_REPO" {
		return 500, nil, errors.New("Unknown sector repository error")
	}

	if symbol == "ERROR_ASSETTYPE_REPO" {
		return 500, nil, errors.New("Unknown asset type repository error")
	}

	if symbol == "ERROR_ASSET_REPO" {
		return 500, nil, errors.New("Unknown asset repository error")
	}

	return 200, &entity.Asset{
		Id:         "TestID",
		Symbol:     symbol,
		Preference: &preference,
		Fullname:   "Test Name",
		AssetType: &entity.AssetType{
			Id:      "TestAssetTypeID",
			Type:    "ETF",
			Name:    "Test ETF",
			Country: country,
		},
		Sector: &entity.Sector{
			Id:   "TestSectorID",
			Name: "Test Sector",
		},
	}, nil
}

func (a *MockApplication) ApiCreateOrder(symbol string, country string, orderType string,
	quantity float64, price float64, currency string, brokerage string,
	date string, userUid string) (int, *entity.Order, error) {

	var assetInfo *entity.Asset
	var httpStatusCode int
	preference := "TestPref"

	layOut := "2006-01-02"
	dateFormatted, _ := time.Parse(layOut, date)

	err := a.app.OrderApp.OrderVerification(orderType, country, quantity, price,
		currency)
	if err != nil {
		return 400, nil, err
	}

	if symbol == "SYMBOL_ALREADY_EXISTS_ERROR" {
		return 500, nil, errors.New("Unknown asset repository error")
	} else if symbol == "SYMBOL_ALREADY_EXISTS" {
		assetInfo = &entity.Asset{
			Id:         "TestID",
			Symbol:     symbol,
			Preference: &preference,
			Fullname:   "Test Name",
			AssetType: &entity.AssetType{
				Id:      "TestAssetTypeID",
				Type:    "ETF",
				Name:    "Test ETF",
				Country: country,
			},
			Sector: &entity.Sector{
				Id:   "TestSectorID",
				Name: "Test Sector",
			},
		}
	} else {
		httpStatusCode, assetInfo, err = a.ApiAssetVerification(symbol, country)
		if err != nil {
			return httpStatusCode, nil, err
		}
	}

	if symbol == "ERROR_ASSETUSER_REPOSITORY" {
		return 500, nil, errors.New("Unknown asset user repository error")
	}

	if brokerage == "UNKNOWN_BROKERAGE" {
		return 400, nil, entity.ErrInvalidBrokerageNameSearch
	}

	return 200, &entity.Order{
		Id:        "TestOrderID",
		Price:     price,
		Quantity:  quantity,
		Currency:  currency,
		OrderType: orderType,
		Date:      dateFormatted,
		Brokerage: &entity.Brokerage{
			Id:      "TestBrokerageID",
			Name:    "Test Brokerage",
			Country: country,
		},
		Asset: &entity.Asset{
			Id:         assetInfo.Id,
			Symbol:     assetInfo.Symbol,
			Preference: assetInfo.Preference,
			Fullname:   assetInfo.Fullname,
		},
	}, nil
}

func (a *MockApplication) ApiAssetsPerAssetType(assetType string, country string,
	ordersInfo bool, userUid string) (int, *entity.AssetType, error) {

	if assetType == "" {
		return 400, nil, entity.ErrInvalidApiQueryTypeBlank
	}

	if country == "" {
		return 400, nil, entity.ErrInvalidApiQueryCountryBlank
	}

	if country != "BR" && country != "US" && country != "" {
		return 400, nil, entity.ErrInvalidCountryCode
	}

	if assetType == "INVALID_ASSET_TYPE" {
		return 400, nil, entity.ErrInvalidAssetTypeName
	}

	preference := "TestPref"
	if ordersInfo == true {
		return 200, &entity.AssetType{
			Id:      "TestAssetTypeID",
			Type:    assetType,
			Name:    "Test Name",
			Country: country,
			Assets: []entity.Asset{
				{
					Id:         "TestAssetID1",
					Symbol:     "TEST1",
					Preference: &preference,
					Fullname:   "Test Name 1",
					Sector: &entity.Sector{
						Id:   "TestSectorID",
						Name: "Test Sector",
					},
					OrderInfo: &entity.OrderInfos{
						WeightedAdjPrice:     20.10,
						WeightedAveragePrice: 20.5,
						TotalQuantity:        30,
					},
				},
				{
					Id:         "TestAssetID2",
					Symbol:     "TEST2",
					Preference: &preference,
					Fullname:   "Test Name 2",
					Sector: &entity.Sector{
						Id:   "TestSectorID",
						Name: "Test Sector",
					},
					OrderInfo: &entity.OrderInfos{
						WeightedAdjPrice:     20.10,
						WeightedAveragePrice: 20.5,
						TotalQuantity:        30,
					},
				},
			},
		}, nil
	} else {
		return 200, &entity.AssetType{
			Id:      "TestAssetTypeID",
			Type:    assetType,
			Name:    "Test Name",
			Country: country,
			Assets: []entity.Asset{
				{
					Id:         "TestAssetID1",
					Symbol:     "TEST1",
					Preference: &preference,
					Fullname:   "Test Name 1",
					Sector: &entity.Sector{
						Id:   "TestSectorID",
						Name: "Test Sector",
					},
				},
				{
					Id:         "TestAssetID2",
					Symbol:     "TEST2",
					Preference: &preference,
					Fullname:   "Test Name 2",
					Sector: &entity.Sector{
						Id:   "TestSectorID",
						Name: "Test Sector",
					},
				},
			},
		}, nil
	}
}

func (a *MockApplication) ApiDeleteAssets(myUser bool, userUid string, symbol string) (int,
	*entity.Asset, error) {
	preference := "TestPref"

	if !myUser && userUid != "USER_WITH_PRIVILEGE" {
		return 403, nil, entity.ErrInvalidUserAdminPrivilege
	}

	if symbol == "ERROR_ASSET_REPO" {
		return 500, nil, errors.New("Unknown asset repository error")
	}

	if symbol == "UNKNOWN_SYMBOL" {
		return 404, nil, entity.ErrInvalidAssetSymbol
	}

	if symbol == "ERROR_ASSETUSER_REPO" {
		return 500, nil, errors.New("Unknown asset user repository error")
	}

	if symbol == "ERROR_ORDERS_REPO" {
		return 500, nil, errors.New("Unknown orders repository error")
	}

	if symbol == "ERROR_EARNINGS_REPO" {
		return 500, nil, errors.New("Unknown earnings repository error")
	}

	return 200, &entity.Asset{
		Id:         "TestID",
		Symbol:     symbol,
		Preference: &preference,
		Fullname:   "Test Name",
		AssetType: &entity.AssetType{
			Id:      "TestAssetTypeID",
			Type:    "STOCK",
			Country: "US",
			Name:    "Test ASTY Name",
		},
		Sector: &entity.Sector{
			Id:   "TestSectorID",
			Name: "Test Sector",
		},
	}, nil
}

func (a *MockApplication) ApiGetOrdersFromAssetUser(symbol string,
	userUid string, orderBy string, limit string, offset string) (int,
	[]entity.Order, error) {

	var offsetInt int

	switch symbol {
	case "":
		return 400, nil, entity.ErrInvalidApiQuerySymbolBlank
	case "ERROR_REPOSITORY_ASSETBYUSER":
		return 500, nil, errors.New("Unknown error on assets repository")
	case "UNKNOWN_SYMBOL":
		return 404, nil, entity.ErrInvalidAssetSymbol

	}

	if limit == "" && offset == "" {
		if symbol == "ERROR_REPOSITORY_ORDERS" {
			return 500, nil, errors.New("Unknown error on orders repository")

		}
	} else {
		_, err := strconv.Atoi(limit)
		if err != nil {
			return 400, nil, entity.ErrInvalidOrderLimit
		}

		offsetInt, err = strconv.Atoi(offset)
		if err != nil {
			return 400, nil, entity.ErrInvalidOrderOffset
		}

		lowerOrderBy := strings.ToLower(orderBy)
		if lowerOrderBy != "asc" && lowerOrderBy != "desc" {
			return 400, nil, entity.ErrInvalidOrderOrderBy
		}

	}

	if symbol == "SYMBOL_WITHOUT_ORDERS" || offsetInt > 2 {
		return 404, nil, entity.ErrInvalidOrder
	}

	layOut := "2006-01-02"
	dateFormatted, _ := time.Parse(layOut, "2021-10-01")
	return 200, []entity.Order{
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

func (a *MockApplication) ApiUpdateOrdersFromUser(orderId string, userUid string, orderType string,
	price float64, quantity float64, date string, brokerage string) (int,
	*entity.Order, error) {

	if orderType == "" || price == 0 || quantity == 0 || date == "" ||
		brokerage == "" {
		return 400, nil, entity.ErrInvalidApiOrderUpdate
	}

	if orderId == "ERROR_ORDER_REPOSITORY" {
		return 500, nil, errors.New("Unknown error in the order repository")
	}

	if orderId == "UNKNOWN_ID" {
		return 404, nil, entity.ErrInvalidOrder
	}

	err := a.app.OrderApp.OrderVerification(orderType, "BR", quantity, price,
		"BRL")
	if err != nil {
		return 400, nil, err
	}

	if brokerage == "UNKNOWN_BROKERAGE" {
		return 400, nil, entity.ErrInvalidBrokerageNameSearch
	}

	layOut := "2006-01-02"
	dateFormatted, _ := time.Parse(layOut, date)
	return 200, &entity.Order{
		Id:        orderId,
		Price:     price,
		Quantity:  quantity,
		Currency:  "BRL",
		OrderType: orderType,
		Date:      dateFormatted,
		Brokerage: &entity.Brokerage{
			Id:      "TestBrokerageID",
			Name:    brokerage,
			Country: "BR",
		},
	}, nil
}

func (a *MockApplication) ApiCreateEarnings(symbol string, currency string,
	earningType string, date string, earnings float64, userUid string) (int,
	*entity.Earnings, error) {

	err := a.app.EarningsApp.EarningsVerification(symbol, currency, earningType,
		date, earnings)
	if err != nil {
		return 400, nil, err
	}

	if symbol == "ERROR_ASSET_REPOSITORY" {
		return 500, nil, errors.New("Unknown error in the asset repository")
	}

	if symbol == "UNKNOWN_SYMBOL" {
		return 404, nil, nil
	}

	layOut := "2006-01-02"
	dateFormatted, _ := time.Parse(layOut, date)
	return 200, &entity.Earnings{
		Id:       "TestEarningID",
		Earning:  earnings,
		Type:     earningType,
		Currency: currency,
		Date:     dateFormatted,
		Asset: &entity.Asset{
			Id:     "TestAssetID",
			Symbol: symbol,
		},
	}, nil
}

func (a *MockApplication) ApiGetEarningsFromAssetUser(symbol string,
	userUid string) (int, []entity.Earnings, error) {

	if symbol == "" {
		return 400, nil, entity.ErrInvalidApiQuerySymbolBlank
	}

	if symbol == "ERROR_ASSET_REPOSITORY" {
		return 500, nil, errors.New("Unknown error in the asset repository")
	}

	if symbol == "INVALID_SYMBOL" {
		return 404, nil, entity.ErrMessageApiAssetSymbolUser
	}

	preference := "TestPref"
	assetInfo := &entity.Asset{
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
	}

	if symbol == "ERRO_EARNINGS_REPOSITORY" {
		return 500, nil, errors.New("Unknown error in the earnings repository")
	}

	if symbol == "SYMBOL_WITHOUT_EARNINGS" {
		return 404, nil, entity.ErrMessageApiEarningAssetUser
	}

	layOut := "2006-01-02"
	dateFormatted, _ := time.Parse(layOut, "2021-10-01")
	return 200, []entity.Earnings{
		{
			Id:       "Earnings1",
			Type:     "Dividendos",
			Earning:  5.29,
			Date:     dateFormatted,
			Currency: "BRL",
			Asset: &entity.Asset{
				Id:     assetInfo.Id,
				Symbol: assetInfo.Symbol,
			},
		},
	}, nil
}

func (a *MockApplication) ApiUpdateEarningsFromUser(earningId string, earning float64,
	earningType string, date string, userUid string) (int, *entity.Earnings,
	error) {

	if earningId == "ERROR_EARNING_REPOSITORY" {
		return 500, nil, errors.New("Unknown error in the earning repository")
	}

	if earningId == "UNKNOWN_EARNING_ID" {
		return 404, nil, entity.ErrMessageApiEarningId
	}

	if earningId == "ERROR_ASSET_REPOSITORY" {
		return 500, nil, errors.New("Unknown error in the asset repository")
	}

	// Verification if the information received in the body attends the
	// requirements of the Earning table
	err := a.app.EarningsApp.EarningsVerification("TEST3",
		"BRL", earningType, date, earning)
	if err != nil {
		return 400, nil, err
	}

	if earningId == "ERROR_UPDATE_EARNING_REPOSITORY" {
		return 500, nil, errors.New("Unknown in the update earning function")
	}
	layOut := "2006-01-02"
	dateFormatted, _ := time.Parse(layOut, date)
	return 200, &entity.Earnings{
		Id:       earningId,
		Earning:  earning,
		Type:     earningType,
		Date:     dateFormatted,
		Currency: "BRL",
		Asset: &entity.Asset{
			Id:     "TestAssetID",
			Symbol: "TEST3",
		},
	}, nil
}
