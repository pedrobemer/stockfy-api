package tables

import "time"

type OrderGeneralInfos struct {
	TotalQuantity        float64 `json:"totalQuantity,omitempty"`
	WeightedAdjPrice     float64 `json:"weightedAdjPrice,omitempty"`
	WeightedAveragePrice float64 `json:"weightedAveragePrice,omitempty"`
}

type AssetQueryReturn struct {
	Id         string `db:"id"`
	Preference *string
	Fullname   string              `db:"fullname"`
	Symbol     string              `db:"symbol"`
	Sector     *SectorApiReturn    `db:"sector" json:",omitempty"`
	AssetType  *AssetTypeApiReturn `db:"asset_type" json:",omitempty"`
	OrderInfo  *OrderGeneralInfos  `db:"orders_info" json:",omitempty"`
	OrdersList []OrderApiReturn    `db:"orders_list" json:",omitempty"`
}

type SectorBodyPost struct {
	Sector string `json:"sector"`
}

type OrderBodyPost struct {
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
	Sector    string  `json:"sector"`
}

type AssetBodyPost struct {
	AssetType string `json:"assetType"`
	Sector    string `json:"sector"`
	Symbol    string `json:"symbol"`
	Fullname  string `json:"fullname"`
	Country   string `json:"country"`
}

type AssetInsert struct {
	AssetType string `json:"assetType"`
	Sector    string `json:"sector"`
	Symbol    string `json:"symbol"`
	Fullname  string `json:"fullname"`
	Country   string `json:"country"`
}

type AssetTypeApiReturn struct {
	Id      string             `db:"id" json:",omitempty"`
	Type    string             `db:"type" json:",omitempty"`
	Name    string             `db:"name" json:",omitempty"`
	Country string             `db:"country" json:",omitempty"`
	Assets  []AssetQueryReturn `db:"assets" json:",omitempty"`
}

type SectorApiReturn struct {
	Id   string `db:"id" json:",omitempty"`
	Name string `db:"name" json:",omitempty"`
}

type BrokerageApiReturn struct {
	Id      string `db:"id"`
	Name    string `db:"name"`
	Country string `db:"country"`
}

type OrderApiReturn struct {
	Id        string             `db:"id"`
	Quantity  float64            `db:"quantity"`
	Price     float64            `db:"price"`
	Currency  string             `db:"currency"`
	OrderType string             `db:"order_type"`
	Date      time.Time          `db:"date"`
	Brokerage BrokerageApiReturn `db:"brokerage"`
}

type AssetApiReturn struct {
	Id         string `db:"id"`
	Preference string `db:"preference"`
	Fullname   string `db:"fullname"`
	Symbol     string `db:"symbol"`
}