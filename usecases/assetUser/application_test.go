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
		{
			assetId:           "DUPLICATED_RELATION",
			userUid:           "User TestValid",
			expectedAssetUser: nil,
			expectedError:     entity.ErrinvalidAssetUserAlreadyExists,
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

func TestDeleteAssetUserRelation(t *testing.T) {
	type test struct {
		assetId                   string
		userUid                   string
		expectedAssetUserRelation *entity.AssetUsers
		expectedError             error
	}

	tests := []test{
		{
			assetId: "TestID",
			userUid: "UserValid",
			expectedAssetUserRelation: &entity.AssetUsers{
				AssetId: "TestID",
				UserUid: "UserValid",
			},
			expectedError: nil,
		},
		{
			assetId:                   "DO_NOT_EXIST",
			userUid:                   "UserValid",
			expectedAssetUserRelation: nil,
			expectedError:             nil,
		},
	}

	mocked := NewMockRepo()
	app := NewApplication(mocked)

	for _, testCase := range tests {
		assestUserRel, err := app.DeleteAssetUserRelation(testCase.assetId,
			testCase.userUid)
		assert.Equal(t, testCase.expectedError, err)
		assert.Equal(t, testCase.expectedAssetUserRelation, assestUserRel)
	}
}
