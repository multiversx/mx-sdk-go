package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-sdk-go/builders"
	"github.com/multiversx/mx-sdk-go/data"
)

const (
	okCodeAfterExecution = "ok"
	internalError        = "internal error"
)

// ArgsVmQueryGetter is the arguments DTO used in the NewvmQueryGetter constructor
type ArgsVmQueryGetter struct {
	Proxy Proxy
	Log   logger.Logger
}

type vmQueryGetter struct {
	proxy Proxy
	log   logger.Logger
}

// NewVmQueryGetter creates a new instance of the vmQueryGetter type
func NewVmQueryGetter(args ArgsVmQueryGetter) (*vmQueryGetter, error) {
	if check.IfNil(args.Proxy) {
		return nil, ErrNilProxy
	}
	if check.IfNil(args.Log) {
		return nil, core.ErrNilLogger
	}

	return &vmQueryGetter{
		proxy: args.Proxy,
		log:   args.Log,
	}, nil
}

// ExecuteQueryReturningBytes will try to execute the provided query and return the result as slice of byte slices
func (dataGetter *vmQueryGetter) ExecuteQueryReturningBytes(ctx context.Context, request *data.VmValueRequest) ([][]byte, error) {
	if request == nil {
		return nil, ErrNilRequest
	}

	response, err := dataGetter.proxy.ExecuteVMQuery(ctx, request)
	if err != nil {
		return nil, err
	}
	dataGetter.log.Debug("executed VMQuery", "FuncName", request.FuncName,
		"Args", request.Args, "SC address", request.Address, "Caller", request.CallerAddr,
		"response.ReturnCode", response.Data.ReturnCode,
		"response.ReturnData", fmt.Sprintf("%+v", response.Data.ReturnData))
	if response.Data.ReturnCode != okCodeAfterExecution {
		return nil, NewQueryResponseError(
			response.Data.ReturnCode,
			response.Data.ReturnMessage,
			request.FuncName,
			request.Address,
			request.Args...,
		)
	}
	return response.Data.ReturnData, nil
}

// ExecuteQueryReturningBool will try to execute the provided query and return the result as bool
func (dataGetter *vmQueryGetter) ExecuteQueryReturningBool(ctx context.Context, request *data.VmValueRequest) (bool, error) {
	response, err := dataGetter.ExecuteQueryReturningBytes(ctx, request)
	if err != nil {
		return false, err
	}

	if len(response) == 0 {
		return false, nil
	}

	return dataGetter.parseBool(response[0], request.FuncName, request.Address, request.Args...)
}

func (dataGetter *vmQueryGetter) parseBool(buff []byte, funcName string, address string, args ...string) (bool, error) {
	if len(buff) == 0 {
		return false, nil
	}

	result, err := strconv.ParseBool(fmt.Sprintf("%d", buff[0]))
	if err != nil {
		return false, NewQueryResponseError(
			internalError,
			fmt.Sprintf("error converting the received bytes to bool, %s", err.Error()),
			funcName,
			address,
			args...,
		)
	}

	return result, nil
}

// ExecuteQueryReturningUint64 will try to execute the provided query and return the result as uint64
func (dataGetter *vmQueryGetter) ExecuteQueryReturningUint64(ctx context.Context, request *data.VmValueRequest) (uint64, error) {
	response, err := dataGetter.ExecuteQueryReturningBytes(ctx, request)
	if err != nil {
		return 0, err
	}

	if len(response) == 0 {
		return 0, nil
	}
	if len(response[0]) == 0 {
		return 0, nil
	}

	num, err := parseUInt64FromByteSlice(response[0])
	if err != nil {
		return 0, NewQueryResponseError(
			internalError,
			err.Error(),
			request.FuncName,
			request.Address,
			request.Args...,
		)
	}

	return num, nil
}

func parseUInt64FromByteSlice(bytes []byte) (uint64, error) {
	num := big.NewInt(0).SetBytes(bytes)
	if !num.IsUint64() {
		return 0, ErrNotUint64Bytes
	}

	return num.Uint64(), nil
}

// ExecuteQueryFromBuilder will try to execute the provided query and return the result as slice of byte slices
func (dataGetter *vmQueryGetter) ExecuteQueryFromBuilder(ctx context.Context, builder builders.VMQueryBuilder) ([][]byte, error) {
	vmValuesRequest, err := builder.ToVmValueRequest()
	if err != nil {
		return nil, err
	}

	return dataGetter.ExecuteQueryReturningBytes(ctx, vmValuesRequest)
}

// ExecuteQueryUint64FromBuilder will try to execute the provided query and return the result as uint64
func (dataGetter *vmQueryGetter) ExecuteQueryUint64FromBuilder(ctx context.Context, builder builders.VMQueryBuilder) (uint64, error) {
	vmValuesRequest, err := builder.ToVmValueRequest()
	if err != nil {
		return 0, err
	}

	return dataGetter.ExecuteQueryReturningUint64(ctx, vmValuesRequest)
}

// ExecuteQueryBoolFromBuilder will try to execute the provided query and return the result as bool
func (dataGetter *vmQueryGetter) ExecuteQueryBoolFromBuilder(ctx context.Context, builder builders.VMQueryBuilder) (bool, error) {
	vmValuesRequest, err := builder.ToVmValueRequest()
	if err != nil {
		return false, err
	}

	return dataGetter.ExecuteQueryReturningBool(ctx, vmValuesRequest)
}

// IsInterfaceNil returns true if there is no value under the interface
func (dataGetter *vmQueryGetter) IsInterfaceNil() bool {
	return dataGetter == nil
}
