package workflows

import (
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/elrond-go/core/check"
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

	var err error
	mbh.cachedNetConfigs, err = args.Proxy.GetNetworkConfig()
	if err != nil {
		return nil, err
	}

	return mbh, nil
}

// GenerateMoveBalanceTransactions wil generate and add to the transaction interactor the move
// balance transactions. Will output a log error if a transaction will be failed.
func (mbh *moveBalanceHandler) GenerateMoveBalanceTransactions(addresses []string) {
	for _, address := range addresses {
		mbh.generateTransactionHandled(address)
	}
}

func (mbh *moveBalanceHandler) generateTransactionHandled(address string) {
	err := mbh.generateTransaction(address)
	if err != nil {
		err = fmt.Errorf("%w for provided address string %s", err, address)
		log.Error(err.Error())
	}
}

func (mbh *moveBalanceHandler) generateTransaction(address string) error {
	addressHandler, err := data.NewAddressFromBech32String(address)
	if err != nil {
		return err
	}

	argsCreate, err := mbh.proxy.GetDefaultTransactionArguments(addressHandler, mbh.cachedNetConfigs)
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

	value := availableBalance.Sub(availableBalance, mbh.computeTxFee(argsCreate))
	argsCreate.Value = value.String()
	argsCreate.RcvAddr = mbh.receiverAddress

	skBytes := mbh.trackableAddressesProvider.PrivateKeyOfBech32Address(address)
	tx, err := mbh.txInteractor.ApplySignatureAndGenerateTransaction(skBytes, argsCreate)
	if err != nil {
		return err
	}

	log.Debug("adding transaction", "from", address, "to", mbh.receiverAddress, "value", value.String())
	mbh.txInteractor.AddTransaction(tx)

	return nil
}

func (mbh *moveBalanceHandler) computeTxFee(argsCreate data.ArgCreateTransaction) *big.Int {
	// this implementation should change if more complex transactions should be generated
	result := big.NewInt(int64(argsCreate.GasPrice))
	result.Mul(result, big.NewInt(int64(argsCreate.GasLimit)))

	return result
}

// IsInterfaceNil returns true if there is no value under the interface
func (mbh *moveBalanceHandler) IsInterfaceNil() bool {
	return mbh == nil
}
