package earnings

import (
	"errors"
	"stockfyApi/entity"
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

	layout := "2006-01-02"
	dateFormatted, _ := time.Parse(layout, date)
	eargningFormatted, err := entity.NewEarnings(earningType, earnings, currency,
		dateFormatted, country, assetId, userUid)
	if err != nil {
		return nil, err
	}

	return eargningFormatted, nil
}

func (a *MockApplication) SearchEarningsFromAssetUser(assetId string,
	userUid string) ([]entity.Earnings, error) {
	layout := "2006-01-02"
	dateFormatted, _ := time.Parse(layout, "2021-10-01")
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

func (a *MockApplication) SearchEarningsFromUser(earningId string, useUid string) (
	*entity.Earnings, error) {

	layout := "2006-01-02"
	dateFormatted, _ := time.Parse(layout, "2021-10-01")

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

	layout := "2006-01-02"
	dateFormatted, _ := time.Parse(layout, date)
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
