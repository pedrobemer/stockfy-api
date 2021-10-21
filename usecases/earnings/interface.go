package earnings

import "stockfyApi/entity"

type Repository interface {
	Create(earningOrder entity.Earnings) ([]entity.Earnings, error)
	DeleteFromAsset(assetId string) ([]entity.Earnings, error)
	SearchFromUser(earningsId string, userUid string) ([]entity.Earnings, error)
	SearchFromAssetUser(assetId string, userUid string) ([]entity.Earnings, error)
	DeleteFromUser(id string, userUid string) (string, error)
	DeleteFromAssetUser(assetId string, userUid string) ([]entity.Earnings, error)
	UpdateFromUser(earningsUpdate entity.Earnings) ([]entity.Earnings, error)
}

type UseCases interface {
	CreateEarning(earningType string, earnings float64, currency string,
		date string, country string, assetId string, userUid string) (
		*entity.Earnings, error)
	SearchEarningsFromAssetUser(assetId string, userUid string) (
		[]entity.Earnings, error)
	SearchEarningsFromUser(earningId string, useUid string) (*entity.Earnings,
		error)
	DeleteEarningsFromUser(earningId string, userUid string) (*string, error)
	DeleteEarningsFromAsset(assetId string) ([]entity.Earnings, error)
	DeleteEarningsFromAssetUser(assetId, userUid string) (*[]entity.Earnings,
		error)
	EarningsUpdate(earningType string, earnings float64, currency string,
		date string, country string, earningId string, userUid string) (
		*entity.Earnings, error)
	EarningsVerification(symbol string, currency string, earningType string,
		date string, earning float64) error
}
