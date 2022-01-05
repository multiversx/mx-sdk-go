package aggregator

import "time"

// ComputeMedian -
func ComputeMedian(nums []float64) (float64, error) {
	return computeMedian(nums)
}

// SetLastNotifiedPrices -
func (pn *priceNotifier) SetLastNotifiedPrices(lastNotifiedPrices []float64) {
	pn.mut.Lock()
	pn.lastNotifiedPrices = lastNotifiedPrices
	pn.mut.Unlock()
}

// SetTimeSinceHandler -
func (pn *priceNotifier) SetTimeSinceHandler(handler func(time time.Time) time.Duration) {
	pn.timeSinceHandler = handler
}

// LastTimeAutoSent -
func (pn *priceNotifier) LastTimeAutoSent() time.Time {
	return pn.lastTimeAutoSent
}
