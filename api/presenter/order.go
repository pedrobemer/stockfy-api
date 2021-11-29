package presenter

import (
	"stockfyApi/entity"
	"time"
)

type OrderBody struct {
	Symbol    string  `json:"symbol"`
	Fullname  string  `json:"fullname"`
	Brokerage string  `json:"brokerage"`
	Quantity  float64 `json:"quantity"`
	Price     float64 `json:"price"`
	Currency  string  `json:"currency"`
	OrderType string  `json:"orderType"`
	Date      string  `json:"date"`
	Country   string  `json:"country"`
	AssetType string  `json:"assetType"`
}

type OrderApiReturn struct {
	Id        string          `json:"id,omitempty"`
	Quantity  float64         `json:"quantity,omitempty"`
	Price     float64         `json:"price,omitempty"`
	Currency  string          `json:"currency,omitempty"`
	OrderType string          `json:"orderType,omitempty"`
	Date      time.Time       `json:"date,omitempty"`
	Brokerage *Brokerage      `json:"brokerage,omitempty"`
	Asset     *AssetApiReturn `json:"asset,omitempty"`
}

type OrderInfos struct {
	TotalQuantity        float64 `json:"totalQuantity,omitempty"`
	WeightedAdjPrice     float64 `json:"weightedAdjPrice,omitempty"`
	WeightedAveragePrice float64 `json:"weightedAveragePrice,omitempty"`
}

func ConvertOrderToApiReturn(orders []entity.Order) []OrderApiReturn {
	convertedOrders := []OrderApiReturn{}
	if orders == nil {
		return nil
	}

	for _, o := range orders {
		convertedOrder := OrderApiReturn{
			Id:        o.Id,
			Quantity:  o.Quantity,
			Price:     o.Price,
			Currency:  o.Currency,
			OrderType: o.OrderType,
			Date:      o.Date,
			Brokerage: ConvertBrokerageToApiReturn(o.Brokerage.Id,
				o.Brokerage.Name, o.Brokerage.Country),
		}

		convertedOrders = append(convertedOrders, convertedOrder)
	}

	return convertedOrders
}

func ConvertSingleOrderToApiReturn(order entity.Order) OrderApiReturn {

	if order.Asset == nil {
		return OrderApiReturn{
			Id:        order.Id,
			Quantity:  order.Quantity,
			Price:     order.Price,
			Currency:  order.Currency,
			OrderType: order.OrderType,
			Date:      order.Date,
			Brokerage: ConvertBrokerageToApiReturn(order.Brokerage.Id,
				order.Brokerage.Name, order.Brokerage.Country),
		}
	}

	return OrderApiReturn{
		Id:        order.Id,
		Quantity:  order.Quantity,
		Price:     order.Price,
		Currency:  order.Currency,
		OrderType: order.OrderType,
		Date:      order.Date,
		Brokerage: ConvertBrokerageToApiReturn(order.Brokerage.Id,
			order.Brokerage.Name, order.Brokerage.Country),
		Asset: &AssetApiReturn{
			Id:       order.Asset.Id,
			Symbol:   order.Asset.Symbol,
			Fullname: order.Asset.Fullname,
		},
	}
}

func ConvertOrderInfoToApiReturn(totalQuantity *float64, weightedAdjPrice *float64,
	weightedAveragePrice *float64) *OrderInfos {

	if totalQuantity == nil && weightedAdjPrice == nil &&
		weightedAveragePrice == nil {
		return nil
	}

	return &OrderInfos{
		TotalQuantity:        *totalQuantity,
		WeightedAdjPrice:     *weightedAdjPrice,
		WeightedAveragePrice: *weightedAveragePrice,
	}
}
