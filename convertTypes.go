package main

import (
	"fmt"
	"strconv"
)

func stringToFloat64(stringValue string) float64 {
	var floatValue float64

	if value, err := strconv.ParseFloat(stringValue, 64); err == nil {
		floatValue = value
	}

	return floatValue
}

func interfaceToFloat64(interfaceValue interface{}) float64 {
	var floatValue float64

	stringValue := fmt.Sprintf("%v", interfaceValue)
	floatValue = stringToFloat64(stringValue)

	return floatValue
}

func interfaceToString(interfaceValue interface{}) string {

	stringValue := fmt.Sprintf("%v", interfaceValue)

	return stringValue
}
