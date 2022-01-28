package data

import "strconv"

func strToFloat32(s string) float32 {
	value, err := strconv.ParseFloat(s, 32)
	if err != nil {
		panic("failed to convert string to float32")
	}

	return float32(value)
}

func strToFloat64(s string) float64 {
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic("failed to convert string to float64")
	}

	return float64(value)
}

func strToBool(s string) bool {
	value, err := strconv.ParseBool(s)
	if err != nil {
		panic("failed to convert string to bool")
	}

	return value
}
