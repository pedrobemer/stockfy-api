package entity

import "time"

func NewEarnings(earningType string, earnings float64, currency string,
	date time.Time, country string, assetId string, userUid string) (*Earnings,
	error) {

	earning := &Earnings{
		Type:     earningType,
		Earning:  earnings,
		Currency: currency,
		Date:     date,
		Asset: &Asset{
			Id: assetId,
		},
		UserUid: userUid,
	}

	err := earning.Validate(country)
	if err != nil {
		return nil, err
	}

	return earning, nil
}

func (a *Earnings) Validate(country string) error {
	if a.Currency != "BRL" && a.Currency != "USD" {
		return ErrInvalidCurrency
	}

	if a.Currency == "BRL" && country != "BR" {
		return ErrInvalidBrazilCurrency
	}

	if a.Currency == "USD" && country != "US" {
		return ErrInvalidUsaCurrency
	}

	return nil
}
