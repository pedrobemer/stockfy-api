package assetusers

import (
	"errors"
	"stockfyApi/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAssetUserRelation(t *testing.T) {
	type test struct {
		assetId           string
		userUid           string
		expectedAssetUser *entity.AssetUsers
		expectedError     error
	}

	tests := []test{
		{
			assetId: "Asset TestValid",
			userUid: "User TestValid",
			expectedAssetUser: &entity.AssetUsers{
				AssetId: "Asset TestValid",
				UserUid: "User TestValid",
			},
			expectedError: nil,
		},
		{
			assetId:           "ERROR_DB",
			userUid:           "User TestValid",
			expectedAssetUser: nil,
			expectedError:     errors.New("TRIGGERED SOME ERROR"),
		},
	}

	mocked := NewMockRepo()
	app := NewApplication(mocked)

	for _, testCase := range tests {
		assetUserRel, err := app.CreateAssetUserRelation(testCase.assetId,
			testCase.userUid)

		assert.Equal(t, testCase.expectedAssetUser, assetUserRel)
		// assert.Equal(t, testCase.userUid, assetUserRel.UserUid)
		assert.Equal(t, testCase.expectedError, err)

	}
}

func TestSearchAssetUserRelation(t *testing.T) {
	type test struct {
		assetId           string
		userUid           string
		expectedAssetUser *entity.AssetUsers
		expectedError     error
	}

	tests := []test{
		{
			assetId: "ValidAsset",
			userUid: "ValidUser",
			expectedAssetUser: &entity.AssetUsers{
				AssetId: "ValidAsset",
				UserUid: "ValidUser",
			},
			expectedError: nil,
		},
		{
			assetId:           "Invalid",
			userUid:           "ValidUser",
			expectedAssetUser: nil,
			expectedError:     nil,
		},
		{
			assetId:           "ERROR_DB",
			userUid:           "ValidUser",
			expectedAssetUser: nil,
			expectedError:     errors.New("TRIGGERED SOME ERROR"),
		},
	}

	mocked := NewMockRepo()
	app := NewApplication(mocked)

	for _, testCase := range tests {
		assetUserRelation, err := app.SearchAssetUserRelation(testCase.assetId,
			testCase.userUid)
		assert.Equal(t, testCase.expectedAssetUser, assetUserRelation)
		assert.Equal(t, testCase.expectedError, err)
	}
}
