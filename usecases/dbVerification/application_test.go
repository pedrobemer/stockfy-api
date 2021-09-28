package dbverification

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRowValidation(t *testing.T) {

	mockedRepo := NewMockRepo()

	dbVerApp := NewApplication(mockedRepo)

	type test struct {
		table        string
		condition    string
		respExpected bool
	}

	tests := []test{
		{
			table:        "Test",
			condition:    "True",
			respExpected: true,
		},
		{
			table:        "Test",
			condition:    "False",
			respExpected: false,
		},
	}

	for _, testCase := range tests {
		validation := dbVerApp.RowValidation(testCase.table, testCase.condition)
		assert.Equal(t, testCase.respExpected, validation)
	}

}
