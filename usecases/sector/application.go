package sector

import (
	"stockfyApi/entity"
)

// type Sector struct {
// 	Id string
// 	Name
// }

type Application struct {
	repo Repository
}

//NewApplication create new use case
func NewApplication(r Repository) *Application {
	return &Application{
		repo: r,
	}
}

func (a *Application) CreateSector(name string) ([]entity.Sector, error) {
	return a.repo.Create(name)
}

func (a *Application) SearchSectorByName(name string) (*entity.Sector, error) {

	if name == "" {
		return nil, entity.ErrInvalidSectorSearchName
	}

	sectorInfo, err := a.repo.SearchByName(name)
	if err != nil {
		return nil, err
	}
	if sectorInfo == nil {
		return nil, nil
	}

	return &sectorInfo[0], nil

}
