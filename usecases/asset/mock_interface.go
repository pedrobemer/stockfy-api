package asset

import (
	"stockfyApi/entity"
)

type Mock struct {
}

func NewMockRepo() *Mock {
	return &Mock{}
}

func (m *Mock) Create(assetInsert entity.Asset) entity.Asset {

	assetCreated := entity.Asset{
		Id:         "a38a9jkrh40a",
		Symbol:     assetInsert.Symbol,
		Preference: assetInsert.Preference,
		Fullname:   assetInsert.Fullname,
	}

	return assetCreated
}
