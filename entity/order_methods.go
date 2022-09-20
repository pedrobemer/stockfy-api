package entity

import "time"

var ValidOrderType = map[string]bool{
	"sell":         true,
	"buy":          true,
	"bonification": true,
	"split":        true,
	"demerge":      true,
}

func NewOrder(quantity float64, price float64, currency string, orderType string,
	date time.Time, brokerageId, assetId, userUid string) (*Order, error) {

	order := &Order{
		Quantity:  quantity,
		Price:     price,
		Currency:  currency,
		OrderType: orderType,
		Date:      date,
		Brokerage: &Brokerage{Id: brokerageId},
		Asset:     &Asset{Id: assetId},
		UserUid:   userUid,
	}

	return order, nil
}
