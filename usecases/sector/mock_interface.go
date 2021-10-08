package sector

import (
	"stockfyApi/entity"
)

type Mock struct {
}

func NewMockRepo() *Mock {
	return &Mock{}
}

func (m *Mock) Create(name string) ([]entity.Sector, error) {

	sectorCreated := []entity.Sector{
		{
			Id:   "a38a9jkrh40a",
			Name: name,
		},
	}

	return sectorCreated, nil
}

func (m *Mock) SearchByName(sector string) ([]entity.Sector, error) {
	if sector == "INVALID" {
		return nil, nil
	}

	return []entity.Sector{
		{
			Id:   "TestID",
			Name: sector,
		},
	}, nil
}
