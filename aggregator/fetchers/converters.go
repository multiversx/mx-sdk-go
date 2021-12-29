package fetchers

import (
	"strconv"
)

// StrToPositiveFloat64 converts the provided string to its float64 representation
func StrToPositiveFloat64(v string) (float64, error) {
	vFloat, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, err
	}
	if vFloat <= 0 {
		return 0, errInvalidResponseData
	}

	return vFloat, nil
}
