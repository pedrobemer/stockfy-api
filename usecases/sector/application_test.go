package sector

import (
	"stockfyApi/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {

	expectedSectorReturn := []entity.Sector{
		{
			Id:   "a38a9jkrh40a",
			Name: "Finance",
		},
	}

	mockedRepo := NewMockRepo()

	sectorApp := NewApplication(mockedRepo)

	sectorReturned, err := sectorApp.CreateSector("Finance")

	assert.Nil(t, err)
	assert.Equal(t, expectedSectorReturn, sectorReturned)

}
