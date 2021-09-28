package asset

import (
	"stockfyApi/entity"
	assettype "stockfyApi/usecases/assetType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {

	preference := "ON"
	assetType := assettype.AssetType{
		Id:      "50vjfnsa",
		Type:    "STOCK",
		Country: "BR",
	}

	expectedAssetCreated := entity.Asset{
		Id:         "a38a9jkrh40a",
		Symbol:     "ITUB4",
		Preference: &preference,
		Fullname:   "Itau Unibanco Holding SA",
	}

	mockedRepo := NewMockRepo()

	assetApp := NewApplication(mockedRepo)

	assetCreated, err := assetApp.CreateAsset("ITUB4", "Itau Unibanco Holding SA",
		&preference, "a40vn4", assetType)

	assert.Nil(t, err)
	assert.Equal(t, expectedAssetCreated, assetCreated)

}
