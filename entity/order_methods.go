package entity

import "time"

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
