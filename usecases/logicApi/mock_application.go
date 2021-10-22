package logicApi

import (
	"errors"
	"stockfyApi/entity"
)

type MockApplication struct {
	//app usecases.Applications
}

//NewApplication create new use case
func NewMockApplication() *MockApplication {
	return &MockApplication{}
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
	return 200, nil, nil
}

func (a *MockApplication) ApiAssetsPerAssetType(assetType string, country string, ordersInfo bool,
	userUid string) (int, *entity.AssetType, error) {
	return 200, nil, nil
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

func (a *MockApplication) ApiGetOrdersFromAssetUser(symbol string, userUid string) (int,
	[]entity.Order, error) {
	return 200, nil, nil
}

func (a *MockApplication) ApiUpdateOrdersFromUser(orderId string, userUid string, orderType string,
	price float64, quantity float64, date string, brokerage string) (int,
	*entity.Order, error) {
	return 200, nil, nil
}

func (a *MockApplication) ApiCreateEarnings(symbol string, currency string, earningType string,
	date string, earnings float64, userUid string) (int, *entity.Earnings,
	error) {
	return 200, nil, nil
}

func (a *MockApplication) ApiGetEarningsFromAssetUser(symbol string, userUid string) (int,
	[]entity.Earnings, error) {
	return 200, nil, nil
}

func (a *MockApplication) ApiUpdateEarningsFromUser(earningId string, earning float64,
	earningType string, date string, userUid string) (int, *entity.Earnings,
	error) {
	return 200, nil, nil
}
