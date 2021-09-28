package sector

import "stockfyApi/entity"

// type Sector struct {
// 	Id string
// 	Name
// }

type Application struct {
	repo SectorRepository
}

//NewApplication create new use case
func NewApplication(r SectorRepository) *Application {
	return &Application{
		repo: r,
	}
}

func (a *Application) CreateSector(name string) ([]entity.Sector, error) {
	return a.repo.Create(name)
}
