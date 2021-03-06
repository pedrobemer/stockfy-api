package earnings

import (
	"errors"
	"stockfyApi/entity"
	"time"
)

type MockDb struct {
}

func NewMockRepo() *MockDb {
	return &MockDb{}
}

func (m *MockDb) Create(earningOrder entity.Earnings) ([]entity.Earnings, error) {
	if earningOrder.Asset.Id == "WRONG_ID" {
		return nil, errors.New("Some Database Error")
	}

	return []entity.Earnings{
		{
			Id:       "ORDER_ID",
			Type:     earningOrder.Type,
			Earning:  earningOrder.Earning,
			Currency: earningOrder.Currency,
			Date:     earningOrder.Date,
			UserUid:  earningOrder.UserUid,
			Asset: &entity.Asset{
				Id: earningOrder.Asset.Id,
			},
		},
	}, nil
}

func (m *MockDb) SearchFromAssetUser(assetId string, userUid string) (
	[]entity.Earnings, error) {
	return []entity.Earnings{}, nil
}

func (r *MockDb) SearchFromAssetUserEarningsByDate(assetId string,
	userUid string, orderBy string, limit int, offset int) ([]entity.Earnings,
	error) {

	if assetId == "UNKNOWN_ID" || offset > 2 {
		return []entity.Earnings{}, nil
	}

	if assetId == "INVALID_ID" {
		return nil, errors.New("UUID SQL ERROR")
	}

	layOut := "2006-01-02"
	tr, _ := time.Parse(layOut, "2021-10-01")

	asset := entity.Asset{
		Id:     assetId,
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

func (m *MockDb) SearchFromUser(earningsId string, userUid string) (
	[]entity.Earnings, error) {
	layout := "2006-01-02"
	dateFormatted, _ := time.Parse(layout, "2021-10-07")

	if earningsId == "INVALID" {
		return nil, errors.New("Some Error")
	}

	return []entity.Earnings{
		{
			Id:       earningsId,
			Type:     "Dividendos",
			Earning:  29.29,
			Date:     dateFormatted,
			Currency: "BRL",
			Asset: &entity.Asset{
				Id:     "AssetID",
				Symbol: "ITUB4",
			},
		},
	}, nil
}

func (m *MockDb) DeleteFromAssetUser(assetId string, userUid string) (
	[]entity.Earnings, error) {
	return []entity.Earnings{}, nil
}

func (m *MockDb) DeleteFromUser(id string, userUid string) (string, error) {
	return "", nil
}

func (m *MockDb) DeleteFromAsset(assetId string) ([]entity.Earnings, error) {
	return []entity.Earnings{}, nil
}

func (m *MockDb) UpdateFromUser(earningsUpdate entity.Earnings) (
	[]entity.Earnings, error) {
	return []entity.Earnings{
		{
			Id:       earningsUpdate.Id,
			Earning:  earningsUpdate.Earning,
			Date:     earningsUpdate.Date,
			Type:     earningsUpdate.Type,
			Currency: earningsUpdate.Currency,
			Asset: &entity.Asset{
				Id:     "AssetID",
				Symbol: "ASSET",
			},
		},
	}, nil
}
