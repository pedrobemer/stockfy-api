package order

import (
	"stockfyApi/entity"
	"time"
)

type Application struct {
	repo Repository
}

//NewApplication create new use case
func NewApplication(r Repository) *Application {
	return &Application{
		repo: r,
	}
}

func (a *Application) CreateOrder(quantity float64, price float64,
	currency string, orderType string, date string, brokerageId string,
	assetId string, userUid string) (*entity.Order, error) {

	layOut := "2006-01-02"
	dateFormatted, _ := time.Parse(layOut, date)
	orderFormatted, err := entity.NewOrder(quantity, price, currency, orderType,
		dateFormatted, brokerageId, assetId, userUid)
	if err != nil {
		return nil, err
	}

	orderCreated := a.repo.Create(*orderFormatted)

	return &orderCreated, nil
}

func (a *Application) OrderVerification(orderType string, country string,
	quantity float64, price float64, currency string) error {

	if orderType != "sell" && orderType != "buy" {
		return entity.ErrInvalidApiOrderType
	}

	if country != "BR" && country != "US" {
		return entity.ErrInvalidCountryCode
	}

	if country == "BR" && (orderType == "sell" || orderType == "buy") {
		if !entity.IsIntegral(quantity) {
			return entity.ErrInvalidApiBrazilOrderQuantity
		}
	}

	if country == "BR" && currency != "BRL" {
		return entity.ErrInvalidApiBrazilOrderCurrency
	}

	if country == "US" && currency != "USD" {
		return entity.ErrInvalidApiUsaOrderCurrency
	}

	if orderType == "buy" && quantity < 0 {
		return entity.ErrInvalidApiOrderBuyQuantity
	} else if orderType == "sell" && quantity > 0 {
		return entity.ErrInvalidApiOrderSellQuantity
	} else if price < 0 {
		return entity.ErrInvalidApiOrderPrice
	}

	return nil
}
