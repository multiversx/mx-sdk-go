package aggregator_test

import (
	"testing"

	"github.com/multiversx/mx-sdk-go/mx-sdk-go-old/aggregator"
	"github.com/stretchr/testify/assert"
)

func TestComputeMedian(t *testing.T) {
	t.Parallel()

	t.Run("nil slice should err", func(t *testing.T) {
		t.Parallel()

		median, err := aggregator.ComputeMedian(nil)
		assert.Equal(t, 0.0, median)
		assert.Equal(t, aggregator.ErrInvalidNumOfElementsToComputeMedian, err)
	})
	t.Run("one value should return that value", func(t *testing.T) {
		t.Parallel()

		prices := []float64{1.0045}
		median, err := aggregator.ComputeMedian(prices)
		assert.Equal(t, 1.0045, median)
		assert.Nil(t, err)
	})
	t.Run("two values should return the median", func(t *testing.T) {
		t.Parallel()

		prices := []float64{1.0045, 1.0047}
		median, err := aggregator.ComputeMedian(prices)
		assert.Equal(t, 1.0046, median)
		assert.Nil(t, err)
	})
	t.Run("three values should return the median", func(t *testing.T) {
		t.Parallel()

		prices := []float64{1.0045, 1.0047, 1.0049}
		median, err := aggregator.ComputeMedian(prices)
		assert.Equal(t, 1.0047, median)
		assert.Nil(t, err)
	})
	t.Run("extreme values should be eliminated", func(t *testing.T) {
		t.Parallel()

		prices := []float64{0.0001, 1.0045, 1.0047, 892789.0}
		median, err := aggregator.ComputeMedian(prices)
		assert.Equal(t, 1.0046, median)
		assert.Nil(t, err)
	})
}
