package earnings

import (
	"errors"
	"stockfyApi/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateEarning(t *testing.T) {
	type test struct {
		earningType     string
		earnings        float64
		currency        string
		date            string
		country         string
		assetId         string
		userUid         string
		expectedEarning *entity.Earnings
		expectedError   error
	}

	layout := "2006-01-02"
	dateFormatted, _ := time.Parse(layout, "2021-07-01")

	tests := []test{
		{
			earningType: "Dividendos",
			earnings:    39.19,
			currency:    "BRL",
			date:        "2021-07-01",
			country:     "BR",
			assetId:     "TestID",
			userUid:     "TestUserUID",
			expectedEarning: &entity.Earnings{
				Id:       "ORDER_ID",
				Type:     "Dividendos",
				Earning:  39.19,
				Currency: "BRL",
				Date:     dateFormatted,
				UserUid:  "TestUserUID",
				Asset: &entity.Asset{
					Id: "TestID",
				},
			},
			expectedError: nil,
		},
		{
			earningType:     "Dividendos",
			earnings:        39.19,
			currency:        "BRL",
			date:            "2021-07-01",
			country:         "US",
			assetId:         "TestID",
			userUid:         "TestUserUID",
			expectedEarning: nil,
			expectedError:   entity.ErrInvalidBrazilCurrency,
		},
		{
			earningType:     "Dividendos",
			earnings:        39.19,
			currency:        "BRL",
			date:            "2021-07-01",
			country:         "BR",
			assetId:         "WRONG_ID",
			userUid:         "TestUserUID",
			expectedEarning: nil,
			expectedError:   errors.New("Some Database Error"),
		},
	}

	mocked := NewMockRepo()
	app := NewApplication(mocked)

	for _, testCase := range tests {
		earningCreated, err := app.CreateEarning(testCase.earningType, testCase.earnings,
			testCase.currency, testCase.date, testCase.country, testCase.assetId,
			testCase.userUid)

		assert.Equal(t, testCase.expectedEarning, earningCreated)
		assert.Equal(t, testCase.expectedError, err)
	}

}

func TestSearchEarningsFromAssetUserByDate(t *testing.T) {
	type test struct {
		assetId          string
		userUid          string
		orderBy          string
		limit            int
		offset           int
		expectedEarnings []entity.Earnings
		expectedError    error
	}

	layOut := "2006-01-02"
	tr, _ := time.Parse(layOut, "2021-10-01")

	asset := entity.Asset{
		Id:     "VALID_ID",
		Symbol: "ITUB4",
	}

	tests := []test{
		{
			assetId:          "VALID_ID",
			userUid:          "VALID_UID",
			orderBy:          "error",
			limit:            2,
			offset:           0,
			expectedEarnings: nil,
			expectedError:    entity.ErrInvalidEarningsOrderBy,
		},
		{
			assetId:          "UNKNOWN_ID",
			userUid:          "VALID_UID",
			orderBy:          "desc",
			limit:            2,
			offset:           0,
			expectedEarnings: []entity.Earnings{},
			expectedError:    nil,
		},
		{
			assetId:          "VALID_ID",
			userUid:          "VALID_UID",
			orderBy:          "desc",
			limit:            2,
			offset:           3,
			expectedEarnings: []entity.Earnings{},
			expectedError:    nil,
		},
		{
			assetId:          "INVALID_ID",
			userUid:          "VALID_UID",
			orderBy:          "desc",
			limit:            2,
			offset:           0,
			expectedEarnings: nil,
			expectedError:    errors.New("UUID SQL ERROR"),
		},
		{
			assetId: "VALID_ID",
			userUid: "VALID_UID",
			orderBy: "desc",
			limit:   2,
			offset:  0,
			expectedEarnings: []entity.Earnings{
				{
					Id:       "3e3e3e3w-ed8b-11eb-9a03-0242ac130003",
					Earning:  5.29,
					Type:     "Dividendos",
					Date:     tr,
					Currency: "BRL",
					Asset:    &asset,
				},
				{
					Id:       "4e4e4e4w-ed8b-11eb-9a03-0242ac130003",
					Earning:  10.48,
					Type:     "JCP",
					Date:     tr,
					Currency: "BRL",
					Asset:    &asset,
				},
			},
			expectedError: nil,
		},
	}

	mocked := NewMockRepo()
	app := NewApplication(mocked)

	for _, testCase := range tests {
		earningsSearch, err := app.SearchEarningsFromAssetUserByDate(
			testCase.assetId, testCase.userUid, testCase.orderBy, testCase.limit,
			testCase.offset)
		assert.Equal(t, testCase.expectedEarnings, earningsSearch)
		assert.Equal(t, testCase.expectedError, err)
	}
}

func TestSearchEarningsFromUser(t *testing.T) {
	type test struct {
		earningId        string
		userUid          string
		expectedEarnings *entity.Earnings
		expectedError    error
	}

	layout := "2006-01-02"
	dateFormatted, _ := time.Parse(layout, "2021-10-07")
	tests := []test{
		{
			earningId: "TestID",
			userUid:   "UserUID",
			expectedEarnings: &entity.Earnings{
				Id:       "TestID",
				Type:     "Dividendos",
				Earning:  29.29,
				Date:     dateFormatted,
				Currency: "BRL",
				Asset: &entity.Asset{
					Id:     "AssetID",
					Symbol: "ITUB4",
				},
			},
			expectedError: nil,
		},
		{
			earningId:        "INVALID",
			userUid:          "UserUID",
			expectedEarnings: nil,
			expectedError:    errors.New("Some Error"),
		},
	}

	mocked := NewMockRepo()
	app := NewApplication(mocked)

	for _, testCase := range tests {
		earningsSearch, err := app.SearchEarningsFromUser(testCase.earningId,
			testCase.userUid)
		assert.Equal(t, testCase.expectedEarnings, earningsSearch)
		assert.Equal(t, testCase.expectedError, err)
	}
}

func TestEarningsUpdate(t *testing.T) {
	type test struct {
		earningType      string
		earnings         float64
		currency         string
		date             string
		country          string
		earningId        string
		userUid          string
		expectedEarnings *entity.Earnings
		expectedError    error
	}

	layout := "2006-01-02"
	dateFormatted, _ := time.Parse(layout, "2021-10-07")
	tests := []test{
		{
			earningType: "Dividendos",
			earnings:    10.49,
			currency:    "USD",
			date:        "2021-10-07",
			country:     "US",
			earningId:   "TestID",
			userUid:     "UserUID",
			expectedEarnings: &entity.Earnings{
				Id:       "TestID",
				Earning:  10.49,
				Date:     dateFormatted,
				Type:     "Dividendos",
				Currency: "USD",
				Asset: &entity.Asset{
					Id:     "AssetID",
					Symbol: "ASSET",
				},
			},
			expectedError: nil,
		},
		{
			earningType:      "Dividendos",
			earnings:         10.49,
			currency:         "USD",
			date:             "2021-10-07",
			country:          "BR",
			earningId:        "TestID",
			userUid:          "UserUID",
			expectedEarnings: nil,
			expectedError:    entity.ErrInvalidUsaCurrency,
		},
	}

	mocked := NewMockRepo()
	app := NewApplication(mocked)

	for _, testCase := range tests {
		earningsUpdated, err := app.EarningsUpdate(testCase.earningType,
			testCase.earnings, testCase.currency, testCase.date, testCase.country,
			testCase.earningId, testCase.userUid)
		assert.Equal(t, testCase.expectedEarnings, earningsUpdated)
		assert.Equal(t, testCase.expectedError, err)
	}
}

func TestEarningsVerification(t *testing.T) {
	type test struct {
		symbol        string
		currency      string
		earningType   string
		date          string
		earning       float64
		expectedError error
	}

	tests := []test{
		{
			symbol:        "ITUB4",
			currency:      "BRL",
			earningType:   "Dividendos",
			date:          "2021-07-01",
			earning:       29.12,
			expectedError: nil,
		},
		{
			symbol:        "",
			currency:      "BRL",
			earningType:   "Dividendos",
			date:          "2021-07-01",
			earning:       29.12,
			expectedError: entity.ErrInvalidEarningsCreateBlankFields,
		},
		{
			symbol:        "ITUB4",
			currency:      "",
			earningType:   "Dividendos",
			date:          "2021-07-01",
			earning:       29.12,
			expectedError: entity.ErrInvalidEarningsCreateBlankFields,
		},
		{
			symbol:        "ITUB4",
			currency:      "BRL",
			earningType:   "",
			date:          "2021-07-01",
			earning:       29.12,
			expectedError: entity.ErrInvalidEarningsCreateBlankFields,
		},
		{
			symbol:        "ITUB4",
			currency:      "BRL",
			earningType:   "Dividendos",
			date:          "",
			earning:       29.12,
			expectedError: entity.ErrInvalidEarningsCreateBlankFields,
		},
		{
			symbol:        "ITUB4",
			currency:      "BRL",
			earningType:   "Dividendos",
			date:          "2021-07-01",
			earning:       -29.12,
			expectedError: entity.ErrInvalidEarningsAmount,
		},
		{
			symbol:        "ITUB4",
			currency:      "BRL",
			earningType:   "WRONG_TYPE",
			date:          "2021-07-01",
			earning:       29.12,
			expectedError: entity.ErrInvalidEarningType,
		},
	}

	mocked := NewMockRepo()
	app := NewApplication(mocked)

	for _, testCase := range tests {
		err := app.EarningsVerification(testCase.symbol, testCase.currency,
			testCase.earningType, testCase.date, testCase.earning)
		assert.Equal(t, testCase.expectedError, err)
	}
}
