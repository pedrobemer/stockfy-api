package sector

import (
	"errors"
	"stockfyApi/entity"
)

type MockApplication struct {
}

func NewMockApplication() *MockApplication {
	return &MockApplication{}
}

func (a *MockApplication) CreateSector(name string) ([]entity.Sector, error) {

	if name == "ERROR_SECTOR" {
		return nil, errors.New("Some Error")
	}

	return []entity.Sector{
		{
			Id:   "TestID",
			Name: name,
		},
	}, nil
}

func (a *MockApplication) SearchSectorByName(name string) (*entity.Sector, error) {

	if name == "" {
		return nil, entity.ErrInvalidSectorSearchName
	}

	switch name {
	case "INVALID_NAME":
		return nil, nil
	case "ERROR_NAME":
		return nil, errors.New("Some Error")
	default:
		return &entity.Sector{
			Id:   "TestID",
			Name: name,
		}, nil
	}

}
