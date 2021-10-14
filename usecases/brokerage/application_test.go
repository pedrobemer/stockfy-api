package brokerage

import (
	"stockfyApi/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchBrokerage(t *testing.T) {
	type test struct {
		searchType              string
		name                    string
		country                 string
		expectedSearchBrokerage []entity.Brokerage
		expectedError           error
	}

	tests := []test{
		{
			searchType: "ALL",
			expectedSearchBrokerage: []entity.Brokerage{
				{
					Id:      "55555555-ed8b-11eb-9a03-0242ac130003",
					Name:    "Clear",
					Country: "BR",
				},
				{
					Id:      "55556666-ed8b-11eb-9a03-0242ac130003",
					Name:    "Rico",
					Country: "BR",
				},
				{
					Id:      "15151515-ed8b-11eb-9a03-0242ac130003",
					Name:    "Avenue",
					Country: "US",
				},
			},
			expectedError: nil,
		},
		{
			searchType: "SINGLE",
			name:       "Rico",
			country:    "BR",
			expectedSearchBrokerage: []entity.Brokerage{
				{
					Id:      "55556666-ed8b-11eb-9a03-0242ac130003",
					Name:    "Rico",
					Country: "BR",
				},
			},
			expectedError: nil,
		},
		{
			searchType: "COUNTRY",
			country:    "BR",
			expectedSearchBrokerage: []entity.Brokerage{
				{
					Id:      "55555555-ed8b-11eb-9a03-0242ac130003",
					Name:    "Clear",
					Country: "BR",
				},
				{
					Id:      "55556666-ed8b-11eb-9a03-0242ac130003",
					Name:    "Rico",
					Country: "BR",
				},
			},
			expectedError: nil,
		},
		{
			searchType:              "INVALID",
			country:                 "BR",
			expectedSearchBrokerage: nil,
			expectedError:           entity.ErrInvalidBrokerageSearchType,
		},
		{
			searchType:              "SINGLE",
			name:                    "",
			expectedSearchBrokerage: nil,
			expectedError:           entity.ErrInvalidBrokerageNameSearchBlank,
		},
		{
			searchType:              "SINGLE",
			name:                    "Invalid",
			expectedSearchBrokerage: nil,
			expectedError:           entity.ErrInvalidBrokerageNameSearch,
		},
	}

	mocked := NewMockRepo()
	app := NewApplication(mocked)

	for _, testCase := range tests {
		searchedBrokerage, err := app.SearchBrokerage(testCase.searchType,
			testCase.name, testCase.country)
		assert.Equal(t, testCase.expectedSearchBrokerage, searchedBrokerage)
		assert.Equal(t, testCase.expectedError, err)
	}

}
