package blockchain

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	"github.com/multiversx/mx-sdk-go/builders"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/testsCommon"
	"github.com/stretchr/testify/assert"
)

const (
	returnCode     = "return code"
	returnMessage  = "return message"
	calledFunction = "called function"
)

var calledArgs = []string{"6172677331", "6172677332"}
var testSCAddressBech32 = "erd1zptg3eu7uw0qvzhnu009lwxupcn6ntjxptj5gaxt8curhxjqr9tsqpsnht"

func createMockProxy(returningBytes [][]byte) *testsCommon.ProxyStub {
	return &testsCommon.ProxyStub{
		ExecuteVMQueryCalled: func(ctx context.Context, vmRequest *data.VmValueRequest) (*data.VmValuesResponseData, error) {
			return &data.VmValuesResponseData{
				Data: &vm.VMOutputApi{
					ReturnCode: okCodeAfterExecution,
					ReturnData: returningBytes,
				},
			}, nil
		},
	}
}

func TestNewVmQueryGetter(t *testing.T) {
	t.Parallel()

	t.Run("nil proxy", func(t *testing.T) {
		t.Parallel()

		dg, err := NewVmQueryGetter(nil)
		assert.Equal(t, ErrNilProxy, err)
		assert.True(t, check.IfNil(dg))
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		dg, err := NewVmQueryGetter(&testsCommon.ProxyStub{})
		assert.Nil(t, err)
		assert.False(t, check.IfNil(dg))
	})
}

func TestNewVmQueryGetter_ExecuteQueryReturningBytes(t *testing.T) {
	t.Parallel()

	testSCAddress, _ := data.NewAddressFromBech32String(testSCAddressBech32)
	t.Run("nil vm ", func(t *testing.T) {
		t.Parallel()

		dg, _ := NewVmQueryGetter(&testsCommon.ProxyStub{})

		result, err := dg.ExecuteQueryReturningBytes(context.Background(), nil)
		assert.Nil(t, result)
		assert.Equal(t, ErrNilRequest, err)
	})
	t.Run("proxy errors", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("expected error")
		proxyInstance := &testsCommon.ProxyStub{
			ExecuteVMQueryCalled: func(ctx context.Context, vmRequest *data.VmValueRequest) (*data.VmValuesResponseData, error) {
				return nil, expectedErr
			},
		}
		dg, _ := NewVmQueryGetter(proxyInstance)

		result, err := dg.ExecuteQueryReturningBytes(context.Background(), &data.VmValueRequest{})
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})

	builder := builders.NewVMQueryBuilder()
	builder.
		Address(testSCAddress).
		Function(calledFunction).
		ArgHexString(calledArgs[0]).
		ArgHexString(calledArgs[1])
	t.Run("return code not ok", func(t *testing.T) {
		t.Parallel()

		expectedErr := NewQueryResponseError(returnCode, returnMessage, calledFunction, testSCAddressBech32, calledArgs...)
		proxyInstance := &testsCommon.ProxyStub{
			ExecuteVMQueryCalled: func(ctx context.Context, vmRequest *data.VmValueRequest) (*data.VmValuesResponseData, error) {
				return &data.VmValuesResponseData{
					Data: &vm.VMOutputApi{
						ReturnData:      nil,
						ReturnCode:      returnCode,
						ReturnMessage:   returnMessage,
						GasRemaining:    0,
						GasRefund:       nil,
						OutputAccounts:  nil,
						DeletedAccounts: nil,
						TouchedAccounts: nil,
						Logs:            nil,
					},
				}, nil
			},
		}

		dg, _ := NewVmQueryGetter(proxyInstance)
		result, err := dg.ExecuteQueryFromBuilder(context.Background(), builder)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, result)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		retData := [][]byte{[]byte("response 1"), []byte("response 2")}
		proxyInstance := &testsCommon.ProxyStub{
			ExecuteVMQueryCalled: func(ctx context.Context, vmRequest *data.VmValueRequest) (*data.VmValuesResponseData, error) {
				return &data.VmValuesResponseData{
					Data: &vm.VMOutputApi{
						ReturnData:      retData,
						ReturnCode:      okCodeAfterExecution,
						ReturnMessage:   returnMessage,
						GasRemaining:    0,
						GasRefund:       nil,
						OutputAccounts:  nil,
						DeletedAccounts: nil,
						TouchedAccounts: nil,
						Logs:            nil,
					},
				}, nil
			},
		}
		dg, _ := NewVmQueryGetter(proxyInstance)

		result, err := dg.ExecuteQueryFromBuilder(context.Background(), builder)
		assert.Nil(t, err)
		assert.Equal(t, retData, result)
	})
}

func TestNewVmQueryGetter_ExecuteQueryReturningBool(t *testing.T) {
	t.Parallel()

	t.Run("nil request", func(t *testing.T) {
		t.Parallel()

		proxyInstance := createMockProxy(make([][]byte, 0))
		dg, _ := NewVmQueryGetter(proxyInstance)

		result, err := dg.ExecuteQueryReturningBool(context.Background(), nil)
		assert.False(t, result)
		assert.Equal(t, ErrNilRequest, err)
	})
	t.Run("empty response", func(t *testing.T) {
		t.Parallel()

		proxyInstance := createMockProxy(make([][]byte, 0))
		dg, _ := NewVmQueryGetter(proxyInstance)

		result, err := dg.ExecuteQueryReturningBool(context.Background(), &data.VmValueRequest{})
		assert.False(t, result)
		assert.Nil(t, err)
	})
	t.Run("empty byte slice on first element", func(t *testing.T) {
		t.Parallel()

		proxyInstance := createMockProxy(make([][]byte, 0))
		dg, _ := NewVmQueryGetter(proxyInstance)

		result, err := dg.ExecuteQueryReturningBool(context.Background(), &data.VmValueRequest{})
		assert.False(t, result)
		assert.Nil(t, err)
	})
	t.Run("not a bool result", func(t *testing.T) {
		t.Parallel()

		proxyInstance := createMockProxy([][]byte{[]byte("random bytes")})
		dg, _ := NewVmQueryGetter(proxyInstance)

		expectedError := NewQueryResponseError(
			internalError,
			`error converting the received bytes to bool, strconv.ParseBool: parsing "114": invalid syntax`,
			"",
			"",
		)

		result, err := dg.ExecuteQueryReturningBool(context.Background(), &data.VmValueRequest{})
		assert.False(t, result)
		assert.Equal(t, expectedError, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		proxyInstance := createMockProxy([][]byte{{1}})
		dg, _ := NewVmQueryGetter(proxyInstance)

		result, err := dg.ExecuteQueryBoolFromBuilder(context.Background(), builders.NewVMQueryBuilder())
		assert.True(t, result)
		assert.Nil(t, err)

		dg.proxy = createMockProxy([][]byte{{0}})

		result, err = dg.ExecuteQueryReturningBool(context.Background(), &data.VmValueRequest{})
		assert.False(t, result)
		assert.Nil(t, err)
	})
}

func TestNewVmQueryGetter_ExecuteQueryReturningUint64(t *testing.T) {
	t.Parallel()

	t.Run("nil request", func(t *testing.T) {
		t.Parallel()

		proxyInstance := createMockProxy(make([][]byte, 0))
		dg, _ := NewVmQueryGetter(proxyInstance)

		result, err := dg.ExecuteQueryReturningUint64(context.Background(), nil)
		assert.Zero(t, result)
		assert.Equal(t, ErrNilRequest, err)
	})
	t.Run("empty response", func(t *testing.T) {
		t.Parallel()

		proxyInstance := createMockProxy(make([][]byte, 0))
		dg, _ := NewVmQueryGetter(proxyInstance)

		result, err := dg.ExecuteQueryReturningUint64(context.Background(), &data.VmValueRequest{})
		assert.Zero(t, result)
		assert.Nil(t, err)
	})
	t.Run("empty byte slice on first element", func(t *testing.T) {
		t.Parallel()

		dg, _ := NewVmQueryGetter(createMockProxy([][]byte{make([]byte, 0)}))

		result, err := dg.ExecuteQueryUint64FromBuilder(context.Background(), builders.NewVMQueryBuilder())
		assert.Zero(t, result)
		assert.Nil(t, err)
	})
	t.Run("large buffer", func(t *testing.T) {
		t.Parallel()

		proxyInstance := createMockProxy([][]byte{[]byte("random bytes")})
		dg, _ := NewVmQueryGetter(proxyInstance)

		expectedError := NewQueryResponseError(
			internalError,
			ErrNotUint64Bytes.Error(),
			"",
			"",
		)

		result, err := dg.ExecuteQueryReturningUint64(context.Background(), &data.VmValueRequest{})
		assert.Zero(t, result)
		assert.Equal(t, expectedError, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		proxyInstance := createMockProxy([][]byte{{1}})
		dg, _ := NewVmQueryGetter(proxyInstance)

		result, err := dg.ExecuteQueryReturningUint64(context.Background(), &data.VmValueRequest{})
		assert.Equal(t, uint64(1), result)
		assert.Nil(t, err)

		dg.proxy = createMockProxy([][]byte{{0xFF, 0xFF}})

		result, err = dg.ExecuteQueryReturningUint64(context.Background(), &data.VmValueRequest{})
		assert.Equal(t, uint64(65535), result)
		assert.Nil(t, err)
	})
}

func TestNewVmQueryGetter_executeQueryWithErroredBuilder(t *testing.T) {
	t.Parallel()

	builder := builders.NewVMQueryBuilder().ArgBytes(nil)

	dg, _ := NewVmQueryGetter(&testsCommon.ProxyStub{})

	resultBytes, err := dg.ExecuteQueryFromBuilder(context.Background(), builder)
	assert.Nil(t, resultBytes)
	assert.True(t, errors.Is(err, builders.ErrInvalidValue))
	assert.True(t, strings.Contains(err.Error(), "builder.ArgBytes"))

	resultUint64, err := dg.ExecuteQueryUint64FromBuilder(context.Background(), builder)
	assert.Zero(t, resultUint64)
	assert.True(t, errors.Is(err, builders.ErrInvalidValue))
	assert.True(t, strings.Contains(err.Error(), "builder.ArgBytes"))

	resultBool, err := dg.ExecuteQueryBoolFromBuilder(context.Background(), builder)
	assert.False(t, resultBool)
	assert.True(t, errors.Is(err, builders.ErrInvalidValue))
	assert.True(t, strings.Contains(err.Error(), "builder.ArgBytes"))
}
