package earnings

import (
	"stockfyApi/entity"
	"strings"
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

	dateFormatted := entity.StringToTime(date)
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

func (a *Application) SearchEarningsFromAssetUser(assetId string, userUid string) (
	[]entity.Earnings, error) {
	earnings, err := a.repo.SearchFromAssetUser(assetId, userUid)
	if err != nil {
		return nil, err
	}

	return earnings, nil
}

func (a *Application) SearchEarningsFromAssetUserByDate(assetId string,
	userUid string, orderBy string, limit int, offset int) ([]entity.Earnings,
	error) {

	lowerOrderBy := strings.ToLower(orderBy)
	if lowerOrderBy != "asc" && lowerOrderBy != "desc" {
		return nil, entity.ErrInvalidEarningsOrderBy
	}

	earnings, err := a.repo.SearchFromAssetUserEarningsByDate(assetId, userUid,
		lowerOrderBy, limit, offset)
	if err != nil {
		return nil, err
	}

	return earnings, nil
}

func (a *Application) SearchEarningsFromUser(earningId string, useUid string) (
	*entity.Earnings, error) {
	earningReturn, err := a.repo.SearchFromUser(earningId, useUid)
	if err != nil {
		if strings.Contains(err.Error(), "invalid input syntax for type uuid") {
			return nil, nil
		}

		return nil, err
	}

	if earningReturn == nil {
		return nil, nil
	}

	return &earningReturn[0], err
}

func (a *Application) DeleteEarningsFromUser(earningId string,
	userUid string) (*string, error) {
	deletedEarningId, err := a.repo.DeleteFromUser(earningId, userUid)
	if err != nil {
		return nil, err
	}

	return &deletedEarningId, nil
}

func (a *Application) DeleteEarningsFromAsset(assetId string) ([]entity.Earnings,
	error) {

	deletedEarnings, err := a.repo.DeleteFromAsset(assetId)
	if err != nil {
		return nil, err
	}

	return deletedEarnings, nil
}

func (a *Application) DeleteEarningsFromAssetUser(assetId, userUid string) (
	[]entity.Earnings, error) {
	deletedEarnings, err := a.repo.DeleteFromAssetUser(assetId, userUid)
	if err != nil {
		return nil, err
	}

	return deletedEarnings, nil
}

func (a *Application) EarningsUpdate(earningType string, earnings float64,
	currency string, date string, country string, earningId string,
	userUid string) (*entity.Earnings, error) {

	dateFormatted := entity.StringToTime(date)
	earningFormatted, err := entity.NewEarnings(earningType, earnings, currency,
		dateFormatted, country, "", userUid)
	if err != nil {
		return nil, err
	}
	earningFormatted.Id = earningId

	updatedEarning, err := a.repo.UpdateFromUser(*earningFormatted)
	if err != nil {
		return nil, err
	}

	return &updatedEarning[0], nil
}

func (a *Application) EarningsVerification(symbol string, currency string,
	earningType string, date string, earning float64) error {

	if symbol == "" || currency == "" || earningType == "" || date == "" {
		return entity.ErrInvalidEarningsCreateBlankFields
	}

	if earning <= 0 {
		return entity.ErrInvalidEarningsAmount
	}

	if !entity.ValidEarningTypes[earningType] {
		return entity.ErrInvalidEarningType
	}

	return nil
}
