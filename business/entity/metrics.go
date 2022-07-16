package entity

import (
	"strconv"
)

var stringConvertMap = map[string]float64{
	"on":  1,
	"off": 0,
}

func ConvertMetricValue(value string) (float64, error) {
	if v, ok := stringConvertMap[value]; ok {
		return v, nil
	}

	v, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return 0, err
	}

	return v, nil
}
