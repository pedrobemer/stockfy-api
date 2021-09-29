package entity

import (
	"fmt"
	"strconv"
)

func StringToFloat64(stringValue string) float64 {
	var floatValue float64

	if value, err := strconv.ParseFloat(stringValue, 64); err == nil {
		floatValue = value
	}

	return floatValue
}

func InterfaceToFloat64(interfaceValue interface{}) float64 {
	var floatValue float64

	stringValue := fmt.Sprintf("%v", interfaceValue)
	floatValue = StringToFloat64(stringValue)

	return floatValue
}

func InterfaceToString(interfaceValue interface{}) string {

	stringValue := fmt.Sprintf("%v", interfaceValue)

	return stringValue
}

func IsIntegral(val float64) bool {
	return val == float64(int(val))
}
