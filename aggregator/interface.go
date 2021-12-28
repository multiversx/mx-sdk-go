package aggregator

// ResponseGetter is the component able to execute a get operation on the provided URL
type ResponseGetter interface {
	Get(url string, response interface{}) error
}
