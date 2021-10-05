package assetusers

import (
	"errors"
	"stockfyApi/entity"
)

type MockDb struct {
}

func NewMockRepo() *MockDb {
	return &MockDb{}
}

func (m *MockDb) Create(assetId string, userUid string) ([]entity.AssetUsers, error) {
	if assetId == "ERROR_DB" {
		return nil, errors.New("TRIGGERED SOME ERROR")
	}

	return []entity.AssetUsers{
		{
			AssetId: assetId,
			UserUid: userUid,
		},
	}, nil
}

func (m *MockDb) Search(assetId string, userUid string) ([]entity.AssetUsers, error) {
	if assetId == "Invalid" {
		return nil, nil
	} else if assetId == "ERROR_DB" {
		return nil, errors.New("TRIGGERED SOME ERROR")
	}

	return []entity.AssetUsers{
		{
			AssetId: assetId,
			UserUid: userUid,
		},
	}, nil
}
