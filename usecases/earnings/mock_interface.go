package earnings

import (
	"errors"
	"stockfyApi/entity"
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
