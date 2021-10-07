package earnings

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

func (a *Application) CreateEarning(earningType string, earnings float64,
	currency string, date string, country string, assetId string,
	userUid string) (*entity.Earnings, error) {

	layout := "2006-01-02"
	dateFormatted, _ := time.Parse(layout, date)
	eargningFormatted, err := entity.NewEarnings(earningType, earnings, currency,
		dateFormatted, country, assetId, userUid)
	if err != nil {
		return nil, err
	}

	earningCreated, err := a.repo.Create(*eargningFormatted)
	if err != nil {
		return nil, err
	}

	return &earningCreated[0], nil
}

func (a *Application) EarningsVerification(symbol string, currency string,
	earningType string, date string, earning float64) error {

	if symbol == "" || currency == "" || earningType == "" || date == "" {
		return entity.ErrInvalidApiMissedKeysBody
	}

	if earning <= 0 {
		return entity.ErrInvalidApiEarningsAmount
	}

	if !entity.ValidEarningTypes[earningType] {
		return entity.ErrInvalidApiEarningType
	}

	return nil
}
