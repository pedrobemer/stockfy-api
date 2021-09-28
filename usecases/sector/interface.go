package sector

import "stockfyApi/entity"

type SectorRepository interface {
	Create(sector string) ([]entity.Sector, error)
	// SearchByName(sector string) ([]entity.Sector, error)
	// SearchByAsset(symbol string) ([]entity.Sector, error)
}
