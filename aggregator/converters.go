package aggregator

import "math"

func trim(v float64, precision float64) float64 {
	return math.Round(v/precision) * precision
}
