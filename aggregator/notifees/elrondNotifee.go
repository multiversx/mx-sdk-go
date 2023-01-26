package notifees

import (
	"context"

	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-sdk-go/aggregator"
	"github.com/multiversx/mx-sdk-go/builders"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
)

const zeroString = "0"
const txVersion = uint32(1)
const function = "submitBatch"
const minGasLimit = uint64(1)

var log = logger.GetOrCreate("mx-sdk-go/aggregator/notifees")

// ArgsMultiversXNotifee is the argument DTO for the NewMultiversXNotifee function
type ArgsMultiversXNotifee struct {
	Proxy           Proxy
	TxBuilder       TxBuilder
	TxNonceHandler  TransactionNonceHandler
	ContractAddress core.AddressHandler
	CryptoHolder    core.CryptoComponentsHolder
	BaseGasLimit    uint64
	GasLimitForEach uint64
}

type multiversXNotifee struct {
	proxy           Proxy
	txBuilder       TxBuilder
	txNonceHandler  TransactionNonceHandler
	contractAddress core.AddressHandler
	baseGasLimit    uint64
	gasLimitForEach uint64
	cryptoHolder    core.CryptoComponentsHolder
}

// NewMultiversXNotifee will create a new instance of multiversXNotifee
func NewMultiversXNotifee(args ArgsMultiversXNotifee) (*multiversXNotifee, error) {
	err := checkArgsMultiversXNotifee(args)
	if err != nil {
		return nil, err
	}

	notifee := &multiversXNotifee{
		proxy:           args.Proxy,
		txBuilder:       args.TxBuilder,
		txNonceHandler:  args.TxNonceHandler,
		contractAddress: args.ContractAddress,
		baseGasLimit:    args.BaseGasLimit,
		gasLimitForEach: args.GasLimitForEach,
		cryptoHolder:    args.CryptoHolder,
	}

	return notifee, nil
}

func checkArgsMultiversXNotifee(args ArgsMultiversXNotifee) error {
	if check.IfNil(args.Proxy) {
		return errNilProxy
	}
	if check.IfNil(args.TxBuilder) {
		return errNilTxBuilder
	}
	if check.IfNil(args.TxNonceHandler) {
		return errNilTxNonceHandler
	}
	if check.IfNil(args.ContractAddress) {
		return errNilContractAddressHandler
	}
	if !args.ContractAddress.IsValid() {
		return errInvalidContractAddress
	}
	if check.IfNil(args.CryptoHolder) {
		return builders.ErrNilCryptoComponentsHolder
	}
	if args.BaseGasLimit < minGasLimit {
		return errInvalidBaseGasLimit
	}
	if args.GasLimitForEach < minGasLimit {
		return errInvalidGasLimitForEach
	}

	return nil
}

// PriceChanged is the function that gets called by a price notifier. This function will assemble a MultiversX
// transaction, having the transaction's data field containing all the price changes information
func (en *multiversXNotifee) PriceChanged(ctx context.Context, priceChanges []*aggregator.ArgsPriceChanged) error {
	txData, err := en.prepareTxData(priceChanges)
	if err != nil {
		return err
	}

	networkConfigs, err := en.proxy.GetNetworkConfig(ctx)
	if err != nil {
		return err
	}

	gasLimit := en.baseGasLimit + uint64(len(priceChanges))*en.gasLimitForEach
	txArgs := data.ArgCreateTransaction{
		Value:    zeroString,
		RcvAddr:  en.contractAddress.AddressAsBech32String(),
		GasPrice: networkConfigs.MinGasPrice,
		GasLimit: gasLimit,
		Data:     txData,
		ChainID:  networkConfigs.ChainID,
		Version:  txVersion,
	}

	err = en.txNonceHandler.ApplyNonceAndGasPrice(ctx, en.cryptoHolder.GetAddressHandler(), &txArgs)
	if err != nil {
		return err
	}

	tx, err := en.txBuilder.ApplySignatureAndGenerateTx(en.cryptoHolder, txArgs)
	if err != nil {
		return err
	}

	txHash, err := en.txNonceHandler.SendTransaction(ctx, tx)
	if err != nil {
		return err
	}

	log.Debug("sent transaction", "hash", txHash)

	return nil
}

func (en *multiversXNotifee) prepareTxData(priceChanges []*aggregator.ArgsPriceChanged) ([]byte, error) {
	txDataBuilder := builders.NewTxDataBuilder()
	txDataBuilder.Function(function)

	for _, priceChange := range priceChanges {
		txDataBuilder.ArgBytes([]byte(priceChange.Base)).
			ArgBytes([]byte(priceChange.Quote)).
			ArgInt64(priceChange.Timestamp).
			ArgInt64(int64(priceChange.DenominatedPrice)).
			ArgInt64(int64(priceChange.Decimals))
	}

	return txDataBuilder.ToDataBytes()
}

// IsInterfaceNil returns true if there is no value under the interface
func (en *multiversXNotifee) IsInterfaceNil() bool {
	return en == nil
}
