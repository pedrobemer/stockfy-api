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

func TestSearchSectorByName(t *testing.T) {
	type test struct {
		name               string
		expectedSectorInfo *entity.Sector
		expectedError      error
	}

	tests := []test{
		{
			name: "Finance",
			expectedSectorInfo: &entity.Sector{
				Id:   "TestID",
				Name: "Finance",
			},
			expectedError: nil,
		},
		{
			name:               "INVALID",
			expectedSectorInfo: nil,
			expectedError:      nil,
		},
	}

	mocked := NewMockRepo()
	app := NewApplication(mocked)

	for _, testCase := range tests {
		sectorInfo, err := app.SearchSectorByName(testCase.name)
		assert.Equal(t, testCase.expectedError, err)
		assert.Equal(t, testCase.expectedSectorInfo, sectorInfo)
	}
}
