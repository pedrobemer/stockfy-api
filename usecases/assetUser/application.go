package assetusers

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

func (a *Application) CreateAssetUserRelation(assetId string, userUid string) (
	*entity.AssetUsers, error) {

	assetUserRelation, err := a.repo.Create(assetId, userUid)
	if err != nil {
		if strings.Contains(err.Error(), "ERROR: duplicate key value violates"+
			" unique constraint \"asset_users_pk\"") {
			return nil, entity.ErrinvalidAssetUserAlreadyExists
		}
		return nil, err
	}

	return &assetUserRelation[0], nil
}

func (a *Application) DeleteAssetUserRelation(assetId string, userUid string) (
	*entity.AssetUsers, error) {
	assetUserRelation, err := a.repo.Delete(assetId, userUid)
	if err != nil {
		return nil, err
	}

	if assetUserRelation == nil {
		return nil, nil
	}

	return &assetUserRelation[0], nil
}

func (a *Application) DeleteAssetUserRelationByAsset(assetId string) (
	[]entity.AssetUsers, error) {
	assetUserRelation, err := a.repo.DeleteByAsset(assetId)
	if err != nil {
		return nil, err
	}

	return assetUserRelation, nil
}

func (a *Application) SearchAssetUserRelation(assetId string, userUid string) (
	*entity.AssetUsers, error) {

	assetUserRelation, err := a.repo.Search(assetId, userUid)
	if err != nil {
		return nil, err
	}

	if assetUserRelation == nil {
		return nil, nil
	}

	return &assetUserRelation[0], nil
}
