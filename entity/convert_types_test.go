package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringToFloat64(t *testing.T) {

	type test struct {
		input          string
		expectedOutput float64
	}

	tests := []test{
		{
			input:          "64.6593",
			expectedOutput: 64.6593,
		},
		{
			input:          "ada90rffj",
			expectedOutput: 0,
		},
		{
			input:          "65,24981",
			expectedOutput: 0,
		},
	}

	for _, testCase := range tests {
		value := StringToFloat64(testCase.input)
		assert.Equal(t, testCase.expectedOutput, value)
	}

}

func TestInterfaceToFloat64(t *testing.T) {

	type test struct {
		input          interface{}
		expectedOutput float64
	}

	tests := []test{
		{
			input:          "64.6593",
			expectedOutput: 64.6593,
		},
		{
			input:          "ada90rffj",
			expectedOutput: 0,
		},
		{
			input:          "65,24981",
			expectedOutput: 0,
		},
	}

	for _, testCase := range tests {
		value := InterfaceToFloat64(testCase.input)
		assert.Equal(t, testCase.expectedOutput, value)
	}

}

func TestInterfaceToString(t *testing.T) {

	type test struct {
		input          interface{}
		expectedOutput string
	}

	tests := []test{
		{
			input:          5938,
			expectedOutput: "5938",
		},
		{
			input:          41029.495,
			expectedOutput: "41029.495",
		},
		{
			input:          "aakdaçde",
			expectedOutput: "aakdaçde",
		},
	}

	for _, testCase := range tests {
		value := InterfaceToString(testCase.input)
		assert.Equal(t, testCase.expectedOutput, value)
	}

}
