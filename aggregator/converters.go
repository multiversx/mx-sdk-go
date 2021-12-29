package aggregator

import (
	"math"
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

func trim(v float64, precision float64) float64 {
	return math.Round(v/precision) * precision
}
