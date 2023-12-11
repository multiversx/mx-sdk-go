package gin

import "context"

type server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}
