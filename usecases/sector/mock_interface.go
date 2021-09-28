package sector

import "stockfyApi/entity"

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
