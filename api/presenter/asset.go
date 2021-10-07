package presenter

import (
	"stockfyApi/entity"
)

type AssetBody struct {
	AssetType string `json:"assetType"`
	// Sector    string `json:"sector"`
	Symbol   string `json:"symbol"`
	Fullname string `json:"fullname"`
	Country  string `json:"country"`
}

type AssetApiReturn struct {
	Id         string            `json:"id"`
	Preference string            `json:"preference,omitempty"`
	Fullname   string            `json:"fullname,omitempty"`
	Symbol     string            `json:"symbol"`
	Sector     *Sector           `json:"sector,omitempty"`
	AssetType  *AssetType        `json:"assetType,omitempty"`
	OrderInfos *OrderInfos       `json:"orderInfos,omitempty"`
	Orders     *[]OrderApiReturn `json:"orders,omitempty"`
}

func ConvertAssetToApiReturn(assetId string, preference string, fullname string,
	symbol string, sectorName string, sectorId string, assetTypeId string,
	assetType string, country string, assetTypeName string, orders []entity.Order,
	orderInfo *entity.OrderInfos) AssetApiReturn {
	var orderInfoReturn *OrderInfos

	sectorReturn := ConvertSectorToApiReturn(sectorId, sectorName)
	assetTypeReturn := ConvertAssetTypeToApiReturn(assetTypeId, assetType,
		assetTypeName, country)
	ordersReturn := ConvertOrderToApiReturn(orders)

	if orderInfo == nil {
		orderInfoReturn = nil
	} else {
		orderInfoReturn = ConvertOrderInfoToApiReturn(&orderInfo.TotalQuantity,
			&orderInfo.WeightedAdjPrice, &orderInfo.WeightedAveragePrice)

	}

	return AssetApiReturn{
		Id:         assetId,
		Preference: preference,
		Fullname:   fullname,
		Symbol:     symbol,
		Sector:     sectorReturn,
		AssetType:  assetTypeReturn,
		Orders:     ordersReturn,
		OrderInfos: orderInfoReturn,
	}
}

func ConvertArrayAssetApiReturn(assets []entity.Asset) []AssetApiReturn {
	var convertedAssets []AssetApiReturn

	for _, asset := range assets {
		convertedAsset := ConvertAssetToApiReturn(asset.Id,
			*asset.Preference, asset.Fullname, asset.Symbol,
			asset.Sector.Name, asset.Sector.Id, "", "", "", "", nil,
			asset.OrderInfo)

		convertedAssets = append(convertedAssets, convertedAsset)
	}

	return convertedAssets

}
