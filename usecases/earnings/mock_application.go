package earnings

import (
	"errors"
	"stockfyApi/entity"
	"strings"
	"time"
)

type MockApplication struct {
}

func NewMockApplication() *MockApplication {
	return &MockApplication{}
}

func (a *MockApplication) CreateEarning(earningType string, earnings float64,
	currency string, date string, country string, assetId string,
	userUid string) (*entity.Earnings, error) {

	dateFormatted := entity.StringToTime(date)
	eargningFormatted, err := entity.NewEarnings(earningType, earnings, currency,
		dateFormatted, country, assetId, userUid)
	if err != nil {
		return nil, err
	}

	return eargningFormatted, nil
}

func (a *MockApplication) SearchEarningsFromAssetUser(assetId string,
	userUid string) ([]entity.Earnings, error) {

	dateFormatted := entity.StringToTime("2021-10-01")
	return []entity.Earnings{
		{
			Id:       "TestEarningID1",
			Type:     "Dividendos",
			Earning:  5.00,
			Date:     dateFormatted,
			Currency: "BRL",
			Asset: &entity.Asset{
				Id:     "TestAssetID",
				Symbol: "TEST3",
			},
		},
		{
			Id:       "TestEarningID2",
			Type:     "Dividendos",
			Earning:  3.00,
			Date:     dateFormatted,
			Currency: "BRL",
			Asset: &entity.Asset{
				Id:     "TestAssetID",
				Symbol: "TEST3",
			},
		},
	}, nil
}

func (a *MockApplication) SearchEarningsFromAssetUserByDate(assetId string,
	userUid string, orderBy string, limit int, offset int) ([]entity.Earnings,
	error) {

	lowerOrderBy := strings.ToLower(orderBy)
	if lowerOrderBy != "asc" && lowerOrderBy != "desc" {
		return nil, entity.ErrInvalidOrderOrderBy
	}

	if assetId == "UNKNOWN_ORDER_REPOSITORY_ERROR" {
		return nil, errors.New("Unknown error in the order repository")
	}

	layOut := "2006-01-02"
	tr, _ := time.Parse(layOut, "2021-10-01")

	asset := entity.Asset{
		Id:     "VALID_ID",
		Symbol: "ITUB4",
	}

	return []entity.Earnings{
		{
			Id:       "3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
			Earning:  5.29,
			Type:     "Dividendos",
			Date:     tr,
			Currency: "BRL",
			Asset:    &asset,
		},
		{
			Id:       "4e4e4e4w-ed8b-11eb-9a03-0242ac130003",
			Earning:  10.48,
			Type:     "JCP",
			Date:     tr,
			Currency: "BRL",
			Asset:    &asset,
		},
	}, nil
}

func (a *MockApplication) SearchEarningsFromUser(earningId string, useUid string) (
	*entity.Earnings, error) {

	dateFormatted := entity.StringToTime("2021-10-01")

	if earningId == "UNKNOWN_EARNING_ID" {
		return nil, nil
	}

	return &entity.Earnings{
		Id:       "TestEarningID1",
		Type:     "Dividendos",
		Earning:  5.00,
		Date:     dateFormatted,
		Currency: "BRL",
		Asset: &entity.Asset{
			Id:     "TestAssetID",
			Symbol: "TEST3",
		},
	}, nil
}

func (a *MockApplication) DeleteEarningsFromUser(earningId string,
	userUid string) (*string, error) {

	if earningId == "INVALID_ID" {
		return nil, errors.New("ERROR: invalid input syntax for type uuid:")
	}

	if earningId == "UNKNOWN_ID" {
		return nil, errors.New("no rows in result set")
	}

	deletedEarningId := earningId

	return &deletedEarningId, nil
}

func (a *MockApplication) DeleteEarningsFromAsset(assetId string) (
	[]entity.Earnings, error) {

	return []entity.Earnings{
		{
			Id: "TestEarningID1",
		},
		{
			Id: "TestEarningID2",
		},
		{
			Id: "TestEarningID3",
		},
	}, nil
}

func (a *MockApplication) DeleteEarningsFromAssetUser(assetId, userUid string) (
	[]entity.Earnings, error) {

	return []entity.Earnings{
		{
			Id: "TestEarningID1",
			Asset: &entity.Asset{
				Id:     "TestAssetID",
				Symbol: "TEST3",
			},
		},
		{
			Id: "TestEarningID2",
			Asset: &entity.Asset{
				Id:     "TestAssetID",
				Symbol: "TEST3",
			},
		},
	}, nil
}

func (a *MockApplication) EarningsUpdate(earningType string, earnings float64,
	currency string, date string, country string, earningId string,
	userUid string) (*entity.Earnings, error) {

	dateFormatted := entity.StringToTime(date)
	earningFormatted, err := entity.NewEarnings(earningType, earnings, currency,
		dateFormatted, country, "", userUid)
	if err != nil {
		return nil, err
	}
	earningFormatted.Id = earningId

	return &entity.Earnings{
		Id:       earningId,
		Earning:  earnings,
		Date:     dateFormatted,
		Type:     earningType,
		Currency: currency,
		Asset: &entity.Asset{
			Id:     "TestAssetID",
			Symbol: "TEST3",
		},
	}, nil
}

func (a *MockApplication) EarningsVerification(symbol string, currency string,
	earningType string, date string, earning float64) error {

	if symbol == "" || currency == "" || earningType == "" || date == "" {
		return entity.ErrInvalidEarningsCreateBlankFields
	}

	if earning <= 0 {
		return entity.ErrInvalidEarningsAmount
	}

	if !entity.ValidEarningTypes[earningType] {
		return entity.ErrInvalidEarningType
	}

	return nil
}
