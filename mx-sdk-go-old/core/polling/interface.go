package polling

import "context"

// Executor defines the behavior of a component able to execute a certain task. This will be continuously called
// by the polling handler
type Executor interface {
	Execute(ctx context.Context) error
	IsInterfaceNil() bool
}
