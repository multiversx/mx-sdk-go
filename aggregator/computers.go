package aggregator

import "sort"

func computeMedian(nums []float64) float64 {
	sort.Float64s(nums)

	numsLen := len(nums)
	mid := numsLen / 2

	if numsLen&1 != 0 {
		return nums[mid]
	}

	return (nums[mid-1] + nums[mid]) / 2
}
