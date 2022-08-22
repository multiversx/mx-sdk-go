package workflows

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

// MoveBalanceHandlerArgs is the argument DTO for the NewMoveBalanceHandler constructor function
type MoveBalanceHandlerArgs struct {
	Proxy                      ProxyHandler
	TxInteractor               TransactionInteractor
	ReceiverAddress            string
	TrackableAddressesProvider TrackableAddressesProvider
	MinimumBalance             *big.Int
}

// moveBalanceHandler is an implementation that can create move balance transactions that will empty the balance
//of the existing accounts
type moveBalanceHandler struct {
	proxy                      ProxyHandler
	mutCachedNetworkConfigs    sync.RWMutex
	cachedNetConfigs           *data.NetworkConfig
	txInteractor               TransactionInteractor
	trackableAddressesProvider TrackableAddressesProvider
	receiverAddress            string
	minimumBalance             *big.Int
}

// NewMoveBalanceHandler creates a new instance of the moveBalanceHandler struct
func NewMoveBalanceHandler(args MoveBalanceHandlerArgs) (*moveBalanceHandler, error) {
	if check.IfNil(args.TrackableAddressesProvider) {
		return nil, ErrNilTrackableAddressesProvider
	}
	if check.IfNil(args.Proxy) {
		return nil, ErrNilProxy
	}
	if check.IfNil(args.TxInteractor) {
		return nil, ErrNilTransactionInteractor
	}
	if args.MinimumBalance == nil {
		return nil, ErrNilMinimumBalance
	}

	mbh := &moveBalanceHandler{
		proxy:                      args.Proxy,
		txInteractor:               args.TxInteractor,
		trackableAddressesProvider: args.TrackableAddressesProvider,
		receiverAddress:            args.ReceiverAddress,
		minimumBalance:             args.MinimumBalance,
	}

	return mbh, nil
}

// CacheNetworkConfigs will try to cache the network configs
func (mbh *moveBalanceHandler) CacheNetworkConfigs(ctx context.Context) error {
	cachedNetConfigs, err := mbh.proxy.GetNetworkConfig(ctx)
	if err != nil {
		return err
	}

	mbh.mutCachedNetworkConfigs.Lock()
	mbh.cachedNetConfigs = cachedNetConfigs
	mbh.mutCachedNetworkConfigs.Unlock()

	return nil
}

// GenerateMoveBalanceTransactions wil generate and add to the transaction interactor the move
// balance transactions. Will output a log error if a transaction will be failed.
func (mbh *moveBalanceHandler) GenerateMoveBalanceTransactions(ctx context.Context, addresses []string) {
	for _, address := range addresses {
		mbh.generateTransactionAndHandleErrors(ctx, address)
	}
}

func (mbh *moveBalanceHandler) generateTransactionAndHandleErrors(ctx context.Context, address string) {
	err := mbh.generateTransaction(ctx, address)
	if err != nil {
		err = fmt.Errorf("%w for provided address %s", err, address)
		log.Error(err.Error())
	}
}

func (mbh *moveBalanceHandler) generateTransaction(ctx context.Context, address string) error {
	addressHandler, err := data.NewAddressFromBech32String(address)
	if err != nil {
		return err
	}

	mbh.mutCachedNetworkConfigs.RLock()
	networkConfigs := mbh.cachedNetConfigs
	mbh.mutCachedNetworkConfigs.RUnlock()

	if networkConfigs == nil {
		return errors.New("nil cached configs")
	}

	argsCreate, err := mbh.proxy.GetDefaultTransactionArguments(ctx, addressHandler, networkConfigs)
	if err != nil {
		return err
	}

	availableBalance, ok := big.NewInt(0).SetString(argsCreate.AvailableBalance, 10)
	if !ok {
		return ErrInvalidAvailableBalanceValue
	}

	if availableBalance.Cmp(mbh.minimumBalance) < 0 {
		log.Debug("will not send move-balance transaction as it is under the set threshold",
			"address", address,
			"available", availableBalance.String(),
			"minimum allowed", mbh.minimumBalance.String(),
		)
		return nil
	}

	//add custom data bytes here if the move-balance transaction towards the hot wallet needs
	// to carry some unique information
	argsCreate.Data = nil
	argsCreate.RcvAddr = mbh.receiverAddress

	value := availableBalance.Sub(availableBalance, mbh.computeTxFee(networkConfigs, argsCreate))
	argsCreate.Value = value.String()

	skBytes := mbh.trackableAddressesProvider.PrivateKeyOfBech32Address(address)
	tx, err := mbh.txInteractor.ApplyUserSignatureAndGenerateTx(skBytes, argsCreate)
	if err != nil {
		return err
	}

	log.Debug("adding transaction", "from", address, "to", mbh.receiverAddress, "value", value.String())
	mbh.txInteractor.AddTransaction(tx)

	return nil
}

func (mbh *moveBalanceHandler) computeTxFee(networkConfigs *data.NetworkConfig, argsCreate data.ArgCreateTransaction) *big.Int {
	// this implementation should change if more complex transactions should be generated
	// if the transaction is required to do a smart contract call, wrap a transaction using the relay mechanism
	// or do an ESDT/SFT/NFT operation, then we need to query the proxy's `/transaction/cost` endpoint route
	// in order to get the correct gas limit

	argsCreate.GasLimit = networkConfigs.MinGasLimit + uint64(len(argsCreate.Data))*networkConfigs.GasPerDataByte
	result := big.NewInt(int64(argsCreate.GasPrice))
	result.Mul(result, big.NewInt(int64(argsCreate.GasLimit)))

	return result
}

// IsInterfaceNil returns true if there is no value under the interface
func (mbh *moveBalanceHandler) IsInterfaceNil() bool {
	return mbh == nil
}
