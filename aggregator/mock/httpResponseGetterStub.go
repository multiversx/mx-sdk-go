package mock

// HttpResponseGetterStub -
type HttpResponseGetterStub struct {
	GetCalled func(url string, response interface{}) error
}

// Get -
func (stub *HttpResponseGetterStub) Get(url string, response interface{}) error {
	if stub.GetCalled != nil {
		return stub.GetCalled(url, response)
	}

	return nil
}
