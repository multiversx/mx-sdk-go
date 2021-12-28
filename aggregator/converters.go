package aggregator

import (
	"strconv"
)

// StrToFloat64 converts the provided string to its float64 representation
func StrToFloat64(v string) (float64, error) {
	vFloat, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, err
	}

	return vFloat, nil
}
