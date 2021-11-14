package presenter

import (
	"stockfyApi/entity"
	"time"
)

type EarningsBody struct {
	Id          string  `json:"id,omitempty"`
	Symbol      string  `json:"symbol"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	EarningType string  `json:"earningType"`
	Date        string  `json:"date"`
}

type EarningsApiReturn struct {
	Id       string          `json:"id"`
	Type     string          `json:"type,omitempty"`
	Earning  float64         `json:"earning,omitempty"`
	Currency string          `json:"currency,omitempty"`
	Date     *time.Time      `json:"date,omitempty"`
	Asset    *AssetApiReturn `json:"asset_id,omitempty"`
}

func ConvertEarningToApiReturn(earningId string, earningType string,
	earning float64, currency string, date time.Time, assetId string,
	assetSymbol string) EarningsApiReturn {
	return EarningsApiReturn{
		Id:       earningId,
		Type:     earningType,
		Earning:  earning,
		Currency: currency,
		Date:     &date,
		Asset: &AssetApiReturn{
			Id:     assetId,
			Symbol: assetSymbol,
		},
	}
}

func ConvertArrayEarningToApiReturn(earnings []entity.Earnings) []EarningsApiReturn {
	var earningsApi []EarningsApiReturn
	for _, earning := range earnings {
		earningApi := ConvertEarningToApiReturn(earning.Id, earning.Type,
			earning.Earning, earning.Currency, earning.Date, earning.Asset.Id,
			earning.Asset.Symbol)
		earningsApi = append(earningsApi, earningApi)
	}

	return earningsApi
}
