package sector

import "stockfyApi/entity"

type Repository interface {
	Create(sector string) ([]entity.Sector, error)
	SearchByName(sector string) ([]entity.Sector, error)
	// SearchByAsset(symbol string) ([]entity.Sector, error)
}

type UseCases interface {
	CreateSector(name string) ([]entity.Sector, error)
}
