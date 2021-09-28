package asset

import (
	"stockfyApi/entity"
	assettype "stockfyApi/usecases/assetType"
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

func (a *Application) CreateAsset(symbol string, fullname string,
	preference *string, sectorId string, assetType assettype.AssetType) (
	entity.Asset, error) {

	assetInfo, err := entity.NewAsset(symbol, fullname, preference, sectorId,
		assetType.Id, assetType.Type, assetType.Country)

	// var preference, ptrPreference string
	// if assetInsert.AssetType.Country == "BR" &&
	// 	assetInsert.AssetType.Type == "STOCK" {
	// 	switch assetInsert.Symbol[len(assetInsert.Symbol)-1:] {
	// 	case "3":
	// 		preference = "ON"
	// 		break
	// 	case "4":
	// 		preference = "PN"
	// 		break
	// 	case "11":
	// 		preference = "UNIT"
	// 		break
	// 	default:
	// 		preference = ""
	// 		break
	// 	}
	// }

	assetCreated := a.repo.Create(*assetInfo)
	return assetCreated, err
}
