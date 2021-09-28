package asset

import (
	"stockfyApi/entity"
)

type MockRepo interface {
	Create(assetInsert entity.Asset) entity.Asset
	Search(symbol string) ([]entity.Asset, error)
	SearchByUser(symbol string, userUid string, orderType string) (
		[]entity.Asset, error)
	SearchPerAssetType(assetType string, country string, userUid string,
		withOrdersInfo bool) []entity.AssetType
	SearchByOrderId(orderId string) []entity.Asset
	Delete(assetId string) []entity.Asset
}

type Mock struct {
	// Create(assetInsert entity.Asset) entity.Asset
	// Search(symbol string) ([]entity.Asset, error)
	// SearchByUser(symbol string, userUid string, orderType string) (
	// 	[]entity.Asset, error)
	// SearchPerAssetType(assetType string, country string, userUid string,
	// 	withOrdersInfo bool) []entity.AssetType
	// SearchByOrderId(orderId string) []entity.Asset
	// Delete(assetId string) []entity.Asset
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
