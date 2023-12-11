package aggregator

import "sort"

const minNumberOfElementsToComputeMedian = 1

func computeMedian(nums []float64) (float64, error) {
	if len(nums) < minNumberOfElementsToComputeMedian {
		return 0, ErrInvalidNumOfElementsToComputeMedian
	}

	sort.Float64s(nums)

	numsLen := len(nums)
	mid := numsLen / 2

	if numsLen&1 != 0 {
		return nums[mid], nil
	}

	return (nums[mid-1] + nums[mid]) / 2, nil
}
