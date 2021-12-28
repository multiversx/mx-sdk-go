package polling

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
)

const minimumPollingInterval = time.Millisecond

// ArgsPollingHandler is the DTO used in the polling handler constructor
type ArgsPollingHandler struct {
	Log              logger.Logger
	Name             string
	PollingInterval  time.Duration
	PollingWhenError time.Duration
	Executor         Executor
}

// pollingHandler represents the component that is able to coordinate the process of calling the
// callBackFunction function continuously until the call to Close is done.
type pollingHandler struct {
	log              logger.Logger
	name             string
	pollingInterval  time.Duration
	pollingWhenError time.Duration
	executor         Executor
	mutState         sync.RWMutex
	cancel           func()
}

// NewPollingHandler will create a new polling handler instance
func NewPollingHandler(args ArgsPollingHandler) (*pollingHandler, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	return &pollingHandler{
		log:              args.Log,
		name:             args.Name,
		pollingInterval:  args.PollingInterval,
		pollingWhenError: args.PollingWhenError,
		executor:         args.Executor,
	}, nil
}

func checkArgs(args ArgsPollingHandler) error {
	if check.IfNil(args.Log) {
		return ErrNilLogger
	}
	if args.PollingInterval < minimumPollingInterval {
		return fmt.Errorf("%w for PollingInterval", ErrInvalidValue)
	}
	if args.PollingWhenError < minimumPollingInterval {
		return fmt.Errorf("%w for PollingWhenError", ErrInvalidValue)
	}
	if check.IfNil(args.Executor) {
		return ErrNilExecutor
	}

	return nil
}

// StartProcessingLoop will start the processing loop
func (ph *pollingHandler) StartProcessingLoop() error {
	ph.mutState.Lock()
	defer ph.mutState.Unlock()

	if ph.cancel != nil {
		return ErrLoopAlreadyStarted
	}

	ctx, cancel := context.WithCancel(context.Background())
	ph.cancel = cancel

	go ph.processLoop(ctx)

	return nil
}

func (ph *pollingHandler) processLoop(ctx context.Context) {
	defer ph.cleanup()

	for {
		pollingChan := time.After(ph.pollingInterval)

		err := ph.executor.Execute(ctx)
		if err != nil {
			ph.log.Error("error in pollingHandler.processLoop",
				"name", ph.name, "error", err,
				"retrying after", ph.pollingWhenError)
			pollingChan = time.After(ph.pollingWhenError)
		}

		select {
		case <-pollingChan:
		case <-ctx.Done():
			ph.log.Debug("pollingHandler's processing loop is closing...",
				"name", ph.name)
			return
		}
	}
}

func (ph *pollingHandler) cleanup() {
	ph.mutState.Lock()
	defer ph.mutState.Unlock()

	ph.cancel = nil
}

// IsRunning returns true if the processing loop is running
func (ph *pollingHandler) IsRunning() bool {
	ph.mutState.RLock()
	defer ph.mutState.RUnlock()

	return ph.cancel != nil
}

// Close will close any containing members and clean any go routines associated
func (ph *pollingHandler) Close() error {
	ph.mutState.RLock()
	defer ph.mutState.RUnlock()

	if ph.cancel != nil {
		ph.cancel()
	}

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (ph *pollingHandler) IsInterfaceNil() bool {
	return ph == nil
}
