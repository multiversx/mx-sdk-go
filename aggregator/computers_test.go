package aggregator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeMedian(t *testing.T) {
	t.Parallel()

	t.Run("one value should return that value", func(t *testing.T) {
		t.Parallel()

		prices := []float64{1.0045}
		assert.Equal(t, 1.0045, computeMedian(prices))
	})
	t.Run("two values should return the median", func(t *testing.T) {
		t.Parallel()

		prices := []float64{1.0045, 1.0047}
		assert.Equal(t, 1.0046, computeMedian(prices))
	})
	t.Run("three values should return the median", func(t *testing.T) {
		t.Parallel()

		prices := []float64{1.0045, 1.0047, 1.0049}
		assert.Equal(t, 1.0047, computeMedian(prices))
	})
	t.Run("extreme values should be eliminated", func(t *testing.T) {
		t.Parallel()

		prices := []float64{0.0001, 1.0045, 1.0047, 892789.0}
		assert.Equal(t, 1.0046, computeMedian(prices))
	})
}
