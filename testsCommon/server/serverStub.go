package server

import "context"

// ServerStub -
type ServerStub struct {
	ListenAndServeCalled func() error
	ShutdownCalled       func(ctx context.Context) error
	CloseCalled          func() error
}

// ListenAndServe -
func (s *ServerStub) ListenAndServe() error {
	if s.ListenAndServeCalled != nil {
		return s.ListenAndServeCalled()
	}
	return nil
}

// Shutdown -
func (s *ServerStub) Shutdown(ctx context.Context) error {
	if s.ShutdownCalled != nil {
		return s.ShutdownCalled(ctx)
	}
	return nil
}

// Close -
func (s *ServerStub) Close() error {
	if s.CloseCalled != nil {
		return s.CloseCalled()
	}
	return nil
}
