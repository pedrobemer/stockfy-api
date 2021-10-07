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
			expectedError: entity.ErrInvalidApiMissedKeysBody,
		},
		{
			symbol:        "ITUB4",
			currency:      "",
			earningType:   "Dividendos",
			date:          "2021-07-01",
			earning:       29.12,
			expectedError: entity.ErrInvalidApiMissedKeysBody,
		},
		{
			symbol:        "ITUB4",
			currency:      "BRL",
			earningType:   "",
			date:          "2021-07-01",
			earning:       29.12,
			expectedError: entity.ErrInvalidApiMissedKeysBody,
		},
		{
			symbol:        "ITUB4",
			currency:      "BRL",
			earningType:   "Dividendos",
			date:          "",
			earning:       29.12,
			expectedError: entity.ErrInvalidApiMissedKeysBody,
		},
		{
			symbol:        "ITUB4",
			currency:      "BRL",
			earningType:   "Dividendos",
			date:          "2021-07-01",
			earning:       -29.12,
			expectedError: entity.ErrInvalidApiEarningsAmount,
		},
		{
			symbol:        "ITUB4",
			currency:      "BRL",
			earningType:   "WRONG_TYPE",
			date:          "2021-07-01",
			earning:       29.12,
			expectedError: entity.ErrInvalidApiEarningType,
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
