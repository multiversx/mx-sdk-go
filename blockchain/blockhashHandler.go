package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/multiversx/mx-chain-core-go/core/atomic"
	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-sdk-go/data"
)

const (
	minPollingInterval  = time.Second
	minBlockTtl         = time.Second
	logPath             = "blockchain/blockhashHandler"
	blockByHashEndpoint = "blocks/%s?fields=timestamp"
)

type argsBlockhashHandler struct {
	pollingInterval   time.Duration
	blockTtl          time.Duration
	httpClientWrapper httpClientWrapper
}

type blockhashHandler struct {
	httpClientWrapper
	sync.RWMutex
	log             logger.Logger
	blockhashes     map[string]int
	loopStatus      *atomic.Flag
	cancel          func()
	pollingInterval time.Duration
	blockTtl        time.Duration
	getTimeHandler  func() time.Time
}

// NewBlockhashHandler returns a new instance of blockhashHandler
func NewBlockhashHandler(args argsBlockhashHandler) (*blockhashHandler, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	bh := &blockhashHandler{
		blockhashes:       make(map[string]int),
		blockTtl:          args.blockTtl,
		loopStatus:        &atomic.Flag{},
		pollingInterval:   args.pollingInterval,
		httpClientWrapper: args.httpClientWrapper,
		getTimeHandler:    time.Now,
	}
	bh.log = logger.GetOrCreate(logPath)
	ctx, cancel := context.WithCancel(context.Background())
	bh.cancel = cancel
	go bh.processLoop(ctx)

	return bh, nil
}

func checkArgs(args argsBlockhashHandler) error {
	if args.pollingInterval < minPollingInterval {
		return fmt.Errorf("%w in checkArgs for value PollingInterval", ErrInvalidValue)
	}
	if args.blockTtl < minBlockTtl {
		return fmt.Errorf("%w in checkArgs for value BlockTtl", ErrInvalidValue)
	}
	if check.IfNil(args.httpClientWrapper) {
		return ErrNilHTTPClientWrapper
	}
	return nil
}

// GetBlockTimestampByHash returns the block by hash
func (bh *blockhashHandler) GetBlockTimestampByHash(ctx context.Context, hash string) (int, error) {
	var block data.Block
	if timestamp, ok := bh.blockhashes[hash]; ok {
		return timestamp, nil
	}

	buff, code, err := bh.httpClientWrapper.GetHTTP(ctx, fmt.Sprintf(blockByHashEndpoint, hash))
	if err != nil || code != http.StatusOK {
		return 0, createHTTPStatusErrorWithBody(code, err, buff)
	}

	err = json.Unmarshal(buff, &block)
	if err != nil {
		return 0, err
	}

	now := bh.getTimeHandler()
	if int64(block.Timestamp) > now.Add(-bh.blockTtl).Unix() {
		bh.Lock()
		defer bh.Unlock()
		bh.blockhashes[hash] = block.Timestamp
	}
	return block.Timestamp, nil
}

func (bh *blockhashHandler) processLoop(ctx context.Context) {
	bh.loopStatus.SetValue(true)
	defer bh.loopStatus.SetValue(false)

	timer := time.NewTimer(bh.pollingInterval)
	defer timer.Stop()

	for {
		bh.Lock()
		defer bh.Unlock()
		for hash, timestamp := range bh.blockhashes {
			now := bh.getTimeHandler()
			if int64(timestamp) <= now.Add(-bh.blockTtl).Unix() {
				delete(bh.blockhashes, hash)
			}
		}
		timer.Reset(bh.pollingInterval)

		select {
		case <-ctx.Done():
			bh.log.Debug("Main execute loop is closing...")
			return
		case <-timer.C:
		}
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (bh *blockhashHandler) IsInterfaceNil() bool {
	return bh == nil
}
