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

type AssetPrice struct {
	ActualPrice float64 `json:"actualPrice"`
	OpenPrice   float64 `json:"openPrice"`
}

type AssetApiReturn struct {
	Id         string           `json:"id,omitempty"`
	Preference string           `json:"preference,omitempty"`
	Fullname   string           `json:"fullname,omitempty"`
	Symbol     string           `json:"symbol,omitempty"`
	Sector     *Sector          `json:"sector,omitempty"`
	AssetType  *AssetType       `json:"assetType,omitempty"`
	OrderInfos *OrderInfos      `json:"orderResume,omitempty"`
	Orders     []OrderApiReturn `json:"orders,omitempty"`
	Price      *AssetPrice      `json:"price,omitempty"`
}

func ConvertAssetToApiReturn(assetId string, preference string, fullname string,
	symbol string, sectorName string, sectorId string, assetTypeId string,
	assetType string, country string, assetTypeName string, orders []entity.Order,
	orderInfo *entity.OrderInfos, price *entity.SymbolPrice) AssetApiReturn {
	var orderInfoReturn *OrderInfos
	var priceInfo *AssetPrice

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

	if price == nil {
		priceInfo = nil
	} else {
		priceInfo = &AssetPrice{
			ActualPrice: price.CurrentPrice,
			OpenPrice:   price.OpenPrice,
		}
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
		Price:      priceInfo,
	}
}

func ConvertArrayAssetApiReturn(assets []entity.Asset) []AssetApiReturn {
	var convertedAssets []AssetApiReturn

	for _, asset := range assets {
		convertedAsset := ConvertAssetToApiReturn(asset.Id,
			*asset.Preference, asset.Fullname, asset.Symbol,
			asset.Sector.Name, asset.Sector.Id, "", "", "", "", nil,
			asset.OrderInfo, asset.Price)

		convertedAssets = append(convertedAssets, convertedAsset)
	}

	return convertedAssets

}
