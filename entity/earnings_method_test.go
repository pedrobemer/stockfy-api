package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEarnings(t *testing.T) {
	tr := time.Now()

	expectedEarnings := &Earnings{
		Type:     "Dividendos",
		Earning:  39.49,
		Currency: "BRL",
		Date:     tr,
		Asset: &Asset{
			Id: "TestID",
		},
		UserUid: "TestUserUID",
	}

	earningsCreated, err := NewEarnings("Dividendos", 39.49, "BRL", tr, "BR",
		"TestID", "TestUserUID")

	assert.Nil(t, err)
	assert.Equal(t, expectedEarnings, earningsCreated)
}

func TestNewEarningsValidation(t *testing.T) {
	tr := time.Now()

	type test struct {
		earningType   string
		earnings      float64
		currency      string
		date          time.Time
		country       string
		assetId       string
		userUid       string
		expectedError error
	}

	tests := []test{
		{
			earningType:   "Dividendos",
			earnings:      39.49,
			currency:      "BRL",
			date:          tr,
			country:       "BR",
			assetId:       "TestID",
			userUid:       "TestUserUID",
			expectedError: nil,
		},
		{
			earningType:   "Dividendos",
			earnings:      39.49,
			currency:      "",
			date:          tr,
			country:       "BR",
			assetId:       "TestID",
			userUid:       "TestUserUID",
			expectedError: ErrInvalidCurrency,
		},
		{
			earningType:   "Dividendos",
			earnings:      39.49,
			currency:      "AAKA",
			date:          tr,
			country:       "BR",
			assetId:       "TestID",
			userUid:       "TestUserUID",
			expectedError: ErrInvalidCurrency,
		},
		{
			earningType:   "Dividendos",
			earnings:      39.49,
			currency:      "BRL",
			date:          tr,
			country:       "US",
			assetId:       "TestID",
			userUid:       "TestUserUID",
			expectedError: ErrInvalidBrazilCurrency,
		},
		{
			earningType:   "Dividendos",
			earnings:      39.49,
			currency:      "USD",
			date:          tr,
			country:       "BR",
			assetId:       "TestID",
			userUid:       "TestUserUID",
			expectedError: ErrInvalidUsaCurrency,
		},
	}

	for _, testCase := range tests {
		_, err := NewEarnings(testCase.earningType, testCase.earnings,
			testCase.currency, testCase.date, testCase.country, testCase.assetId,
			testCase.userUid)
		assert.Equal(t, testCase.expectedError, err)
	}
}
