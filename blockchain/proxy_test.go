package blockchain

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/multiversx/mx-chain-storage-go/lrucache"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/api"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-go/config"
	"github.com/multiversx/mx-chain-go/state"
	sdkCore "github.com/multiversx/mx-sdk-go/core"
	sdkHttp "github.com/multiversx/mx-sdk-go/core/http"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/testsCommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testHttpURL = "https://test.org"
const networkConfigEndpoint = "network/config"
const getNetworkStatusEndpoint = "network/status/%d"
const getNodeStatusEndpoint = "node/status"

// not a real-world valid test query option but rather a test one to check all fields are properly set
var testQueryOptions = api.AccountQueryOptions{
	OnFinalBlock: true,
	OnStartOfEpoch: core.OptionalUint32{
		Value:    3737,
		HasValue: true,
	},
	BlockNonce: core.OptionalUint64{
		Value:    3838,
		HasValue: true,
	},
	BlockHash:     []byte("block hash"),
	BlockRootHash: []byte("block root hash"),
	HintEpoch: core.OptionalUint32{
		Value:    3939,
		HasValue: true,
	},
}

type testStruct struct {
	Nonce int
	Name  string
}

type mockHTTPClient struct {
	doCalled func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.doCalled != nil {
		return m.doCalled(req)
	}

	return nil, errors.New("not implemented")
}

func createMockClientRespondingBytes(responseBytes []byte) *mockHTTPClient {
	return &mockHTTPClient{
		doCalled: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       io.NopCloser(bytes.NewReader(responseBytes)),
				StatusCode: http.StatusOK,
			}, nil
		},
	}
}

func createMockClientMultiResponse(responseMap map[string][]byte) *mockHTTPClient {
	return &mockHTTPClient{
		doCalled: func(req *http.Request) (*http.Response, error) {
			responseBytes, exists := responseMap[req.URL.String()]
			if !exists {
				return nil, fmt.Errorf("no response for URL: %s", req.URL.String())
			}
			return &http.Response{
				Body:       io.NopCloser(bytes.NewReader(responseBytes)),
				StatusCode: http.StatusOK,
			}, nil
		},
	}
}

func createMockClientRespondingBytesWithStatus(responseBytes []byte, status int) *mockHTTPClient {
	return &mockHTTPClient{
		doCalled: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       io.NopCloser(bytes.NewReader(responseBytes)),
				StatusCode: status,
			}, nil
		},
	}
}

func createMockClientRespondingError(err error) *mockHTTPClient {
	return &mockHTTPClient{
		doCalled: func(req *http.Request) (*http.Response, error) {
			return nil, err
		},
	}
}

func createMockArgsProxy(httpClient sdkHttp.Client) ArgsProxy {
	return ArgsProxy{
		ProxyURL:            testHttpURL,
		Client:              httpClient,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		AllowedDeltaToFinal: 1,
		CacheExpirationTime: time.Second,
		EntityType:          sdkCore.ObserverNode,
	}
}

func createMockArgsProxyWithCache(httpClient sdkHttp.Client) ArgsProxy {
	// Create an LRU cache instance with the desired size
	cacheSize := 1000
	lruCacheInstance, err := lrucache.NewCache(cacheSize)
	if err != nil {
		panic("Failed to create LRU cache: " + err.Error())
	}

	return ArgsProxy{
		ProxyURL:               testHttpURL,
		Client:                 httpClient,
		SameScState:            false,
		ShouldBeSynced:         false,
		FinalityCheck:          false,
		AllowedDeltaToFinal:    1,
		CacheExpirationTime:    time.Second,
		EntityType:             sdkCore.ObserverNode,
		FilterQueryBlockCacher: lruCacheInstance,
	}
}

func handleRequestNetworkConfigAndStatus(
	req *http.Request,
	numShards uint32,
	currentNonce uint64,
	highestNonce uint64,
) (*http.Response, bool, error) {

	handled := false
	url := req.URL.String()
	var response interface{}
	switch url {
	case fmt.Sprintf("%s/%s", testHttpURL, networkConfigEndpoint):
		handled = true
		response = data.NetworkConfigResponse{
			Data: struct {
				Config *data.NetworkConfig `json:"config"`
			}{
				Config: &data.NetworkConfig{
					NumShardsWithoutMeta: numShards,
				},
			},
		}

	case fmt.Sprintf("%s/%s", testHttpURL, getNodeStatusEndpoint):
		handled = true
		response = data.NodeStatusResponse{
			Data: struct {
				Status *data.NetworkStatus `json:"metrics"`
			}{
				Status: &data.NetworkStatus{
					Nonce:                currentNonce,
					HighestNonce:         highestNonce,
					ProbableHighestNonce: currentNonce,
					ShardID:              2,
				},
			},
			Error: "",
			Code:  "",
		}
	case fmt.Sprintf("%s/%s", testHttpURL, fmt.Sprintf(getNetworkStatusEndpoint, 2)):
		handled = true
		response = data.NetworkStatusResponse{
			Data: struct {
				Status *data.NetworkStatus `json:"status"`
			}{
				Status: &data.NetworkStatus{
					Nonce:                currentNonce,
					HighestNonce:         highestNonce,
					ProbableHighestNonce: currentNonce,
					ShardID:              2,
				},
			},
			Error: "",
			Code:  "",
		}
	}

	if !handled {
		return nil, handled, nil
	}

	buff, _ := json.Marshal(response)
	return &http.Response{
		Body:       io.NopCloser(bytes.NewReader(buff)),
		StatusCode: http.StatusOK,
	}, handled, nil
}

func TestNewProxy(t *testing.T) {
	t.Parallel()

	t.Run("invalid time cache should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsProxy(nil)
		args.CacheExpirationTime = time.Second - time.Nanosecond
		proxyInstance, err := NewProxy(args)

		assert.True(t, check.IfNil(proxyInstance))
		assert.True(t, errors.Is(err, ErrInvalidCacherDuration))
	})
	t.Run("invalid nonce delta should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsProxy(nil)
		args.FinalityCheck = true
		args.AllowedDeltaToFinal = 0
		proxyInstance, err := NewProxy(args)

		assert.True(t, check.IfNil(proxyInstance))
		assert.True(t, errors.Is(err, ErrInvalidAllowedDeltaToFinal))
	})
	t.Run("should work with finality check", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsProxy(nil)
		args.FinalityCheck = true
		proxyInstance, err := NewProxy(args)

		assert.False(t, check.IfNil(proxyInstance))
		assert.Nil(t, err)
	})
	t.Run("should work without finality check", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsProxy(nil)
		proxyInstance, err := NewProxy(args)

		assert.False(t, check.IfNil(proxyInstance))
		assert.Nil(t, err)
	})
}

func TestGetAccount(t *testing.T) {
	t.Parallel()

	numAccountQueries := uint32(0)
	httpClient := &mockHTTPClient{
		doCalled: func(req *http.Request) (*http.Response, error) {
			response, handled, err := handleRequestNetworkConfigAndStatus(req, 3, 9170526, 9170526)
			if handled {
				return response, err
			}

			account := data.AccountResponse{
				Data: struct {
					Account *data.Account `json:"account"`
				}{
					Account: &data.Account{
						Nonce:   37,
						Balance: "38",
					},
				},
			}
			accountBytes, _ := json.Marshal(account)
			atomic.AddUint32(&numAccountQueries, 1)
			return &http.Response{
				Body:       io.NopCloser(bytes.NewReader(accountBytes)),
				StatusCode: http.StatusOK,
			}, nil
		},
	}
	args := createMockArgsProxy(httpClient)
	args.FinalityCheck = true
	proxyInstance, _ := NewProxy(args)

	address, err := data.NewAddressFromBech32String("erd1qqqqqqqqqqqqqpgqfzydqmdw7m2vazsp6u5p95yxz76t2p9rd8ss0zp9ts")
	if err != nil {
		assert.Error(t, err)
	}
	expectedErr := errors.New("expected error")

	t.Run("nil address should error", func(t *testing.T) {
		t.Parallel()

		response, errGet := proxyInstance.GetAccount(context.Background(), nil)
		require.Equal(t, ErrNilAddress, errGet)
		require.Nil(t, response)
	})
	t.Run("invalid address should error", func(t *testing.T) {
		t.Parallel()

		invalidAddress := data.NewAddressFromBytes([]byte("invalid address"))
		response, errGet := proxyInstance.GetAccount(context.Background(), invalidAddress)
		require.Equal(t, ErrInvalidAddress, errGet)
		require.Nil(t, response)
	})
	t.Run("finality checker errors should not query", func(t *testing.T) {
		proxyInstance.finalityProvider = &testsCommon.FinalityProviderStub{
			CheckShardFinalizationCalled: func(ctx context.Context, targetShardID uint32, maxNoncesDelta uint64) error {
				return expectedErr
			},
		}

		account, errGet := proxyInstance.GetAccount(context.Background(), address)
		assert.Nil(t, account)
		assert.True(t, errors.Is(errGet, expectedErr))
		assert.Equal(t, uint32(0), atomic.LoadUint32(&numAccountQueries))
	})
	t.Run("finality checker returns nil should return account", func(t *testing.T) {
		finalityCheckCalled := uint32(0)
		proxyInstance.finalityProvider = &testsCommon.FinalityProviderStub{
			CheckShardFinalizationCalled: func(ctx context.Context, targetShardID uint32, maxNoncesDelta uint64) error {
				atomic.AddUint32(&finalityCheckCalled, 1)
				return nil
			},
		}

		account, errGet := proxyInstance.GetAccount(context.Background(), address)
		assert.NotNil(t, account)
		assert.Equal(t, uint64(37), account.Nonce)
		assert.Nil(t, errGet)
		assert.Equal(t, uint32(1), atomic.LoadUint32(&numAccountQueries))
		assert.Equal(t, uint32(1), atomic.LoadUint32(&finalityCheckCalled))
	})
}

func TestProxy_GetNetworkEconomics(t *testing.T) {
	t.Parallel()

	responseBytes := []byte(`{"data":{"metrics":{"erd_dev_rewards":"0","erd_epoch_for_economics_data":263,"erd_inflation":"5869888769785838708144","erd_total_fees":"51189055176110000000","erd_total_staked_value":"9963775651405816710680128","erd_total_supply":"21556417261819025351089574","erd_total_top_up_value":"1146275808171377418645274"}},"code":"successful"}`)
	httpClient := createMockClientRespondingBytes(responseBytes)
	args := createMockArgsProxy(httpClient)
	ep, _ := NewProxy(args)

	networkEconomics, err := ep.GetNetworkEconomics(context.Background())
	require.Nil(t, err)
	require.Equal(t, &data.NetworkEconomics{
		DevRewards:            "0",
		EpochForEconomicsData: 263,
		Inflation:             "5869888769785838708144",
		TotalFees:             "51189055176110000000",
		TotalStakedValue:      "9963775651405816710680128",
		TotalSupply:           "21556417261819025351089574",
		TotalTopUpValue:       "1146275808171377418645274",
	}, networkEconomics)
}

func TestProxy_RequestTransactionCost(t *testing.T) {
	t.Parallel()

	responseBytes := []byte(`{"data":{"txGasUnits":24273810,"returnMessage":""},"error":"","code":"successful"}`)
	httpClient := createMockClientRespondingBytes(responseBytes)
	args := createMockArgsProxy(httpClient)
	ep, _ := NewProxy(args)

	tx := &transaction.FrontendTransaction{
		Nonce:    1,
		Value:    "50",
		Receiver: "erd1rh5ws22jxm9pe7dtvhfy6j3uttuupkepferdwtmslms5fydtrh5sx3xr8r",
		Sender:   "erd1rh5ws22jxm9pe7dtvhfy6j3uttuupkepferdwtmslms5fydtrh5sx3xr8r",
		Data:     []byte("hello"),
		ChainID:  "1",
		Version:  1,
		Options:  0,
	}
	txCost, err := ep.RequestTransactionCost(context.Background(), tx)
	require.Nil(t, err)
	require.Equal(t, &data.TxCostResponseData{
		TxCost:     24273810,
		RetMessage: "",
	}, txCost)
}

func TestProxy_GetTransactionInfoWithResults(t *testing.T) {
	t.Parallel()

	responseBytes := []byte(`{"data":{"transaction":{"type":"normal","nonce":22100,"round":53057,"epoch":263,"value":"0","receiver":"erd1e6c9vcga5lyhwdu9nr9lya4ujz2r5w2egsfjp0lslrgv5dsccpdsmre6va","sender":"erd1e6c9vcga5lyhwdu9nr9lya4ujz2r5w2egsfjp0lslrgv5dsccpdsmre6va","gasPrice":1000000001,"gasLimit":19210000,"data":"RVNEVE5GVFRyYW5zZmVyQDU5NTk1OTQ1MzY0MzM5MkQzMDM5MzQzOTMwMzBAMDFAMDFAMDAwMDAwMDAwMDAwMDAwMDA1MDA1QzgzRTBDNDJFRENFMzk0RjQwQjI0RDI5RDI5OEIwMjQ5QzQxRjAyODk3NEA2Njc1NkU2NEA1N0M1MDZBMTlEOEVBRTE4Mjk0MDNBOEYzRjU0RTFGMDM3OUYzODE1N0ZDRjUzRDVCQ0E2RjIzN0U0QTRDRjYxQDFjMjA=","signature":"37922ccb13d46857d819cd618f1fd8e76777e14f1c16a54eeedd55c67d5883db0e3a91cd0ff0f541e2c3d7347488b4020c0ba4765c9bc02ea0ae1c3f6db1ec05","sourceShard":1,"destinationShard":1,"blockNonce":53052,"blockHash":"4a63312d1bfe48aa516185d12abff5daf6343fce1f298db51291168cf97a790c","notarizedAtSourceInMetaNonce":53053,"NotarizedAtSourceInMetaHash":"342d189e36ef5cbf9f8b3f2cb5bf2cb8e2260062c6acdf89ce5faacd99f4dbcc","notarizedAtDestinationInMetaNonce":53053,"notarizedAtDestinationInMetaHash":"342d189e36ef5cbf9f8b3f2cb5bf2cb8e2260062c6acdf89ce5faacd99f4dbcc","miniblockType":"TxBlock","miniblockHash":"0c659cce5e2653522cc0e3cf35571264522035a7aef4ffa5244d1ed3d8bc01a8","status":"success","hyperblockNonce":53053,"hyperblockHash":"342d189e36ef5cbf9f8b3f2cb5bf2cb8e2260062c6acdf89ce5faacd99f4dbcc","smartContractResults":[{"hash":"5ab14959aaeb3a20d95ec6bbefc03f251732d9368711c55c63d3811e70903f4e","nonce":0,"value":0,"receiver":"erd1qqqqqqqqqqqqqpgqtjp7p3pwmn3efaqtynff62vtqfyug8cz396qs5vnsy","sender":"erd1e6c9vcga5lyhwdu9nr9lya4ujz2r5w2egsfjp0lslrgv5dsccpdsmre6va","data":"ESDTNFTTransfer@595959453643392d303934393030@01@01@080112020001226f08011204746573741a20ceb056611da7c977378598cbf276bc90943a3959441320bff0f8d0ca3618c05b20e8072a206161616161616161616161616161616161616161616161616161616161616161320461626261321268747470733a2f2f656c726f6e642e636f6d3a0474657374@66756e64@57c506a19d8eae1829403a8f3f54e1f0379f38157fcf53d5bca6f237e4a4cf61@1c20","prevTxHash":"c7fadeaccce0673bd6ce3a1f472f7cc1beef20c0b3131cfa9866cd5075816639","originalTxHash":"c7fadeaccce0673bd6ce3a1f472f7cc1beef20c0b3131cfa9866cd5075816639","gasLimit":18250000,"gasPrice":1000000001,"callType":0},{"hash":"f337d2705d2b644f9d8f75cf270e879b6ada51c4c54009e94305adf368d1adbc","nonce":1,"value":102793170000000,"receiver":"erd1e6c9vcga5lyhwdu9nr9lya4ujz2r5w2egsfjp0lslrgv5dsccpdsmre6va","sender":"erd1qqqqqqqqqqqqqpgqtjp7p3pwmn3efaqtynff62vtqfyug8cz396qs5vnsy","data":"@6f6b","prevTxHash":"5ab14959aaeb3a20d95ec6bbefc03f251732d9368711c55c63d3811e70903f4e","originalTxHash":"c7fadeaccce0673bd6ce3a1f472f7cc1beef20c0b3131cfa9866cd5075816639","gasLimit":0,"gasPrice":1000000001,"callType":0}]}},"error":"","code":"successful"}`)
	httpClient := createMockClientRespondingBytes(responseBytes)
	args := createMockArgsProxy(httpClient)
	ep, _ := NewProxy(args)

	tx, err := ep.GetTransactionInfoWithResults(context.Background(), "a40e5a6af4efe221608297a73459211756ab88b96896e6e331842807a138f343")
	require.Nil(t, err)

	txBytes, _ := json.MarshalIndent(tx, "", " ")
	fmt.Println(string(txBytes))
}

func TestProxy_ExecuteVmQuery(t *testing.T) {
	t.Parallel()

	responseBytes := []byte(`{"data":{"data":{"returnData":["MC41LjU="],"returnCode":"ok","returnMessage":"","gasRemaining":18446744073685949187,"gasRefund":0,"outputAccounts":{"0000000000000000050033bb65a91ee17ab84c6f8a01846ef8644e15fb76696a":{"address":"erd1qqqqqqqqqqqqqpgqxwakt2g7u9atsnr03gqcgmhcv38pt7mkd94q6shuwt","nonce":0,"balance":null,"balanceDelta":0,"storageUpdates":{},"code":null,"codeMetaData":null,"outputTransfers":[],"callType":0}},"deletedAccounts":[],"touchedAccounts":[],"logs":[]}},"error":"","code":"successful"}`)
	t.Run("no finality check", func(t *testing.T) {
		httpClient := createMockClientRespondingBytes(responseBytes)
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		response, err := ep.ExecuteVMQuery(context.Background(), &data.VmValueRequest{
			Address:    "erd1qqqqqqqqqqqqqpgqxwakt2g7u9atsnr03gqcgmhcv38pt7mkd94q6shuwt",
			FuncName:   "version",
			CallerAddr: "erd1qqqqqqqqqqqqqpgqxwakt2g7u9atsnr03gqcgmhcv38pt7mkd94q6shuwt",
		})
		require.Nil(t, err)
		require.Equal(t, "0.5.5", string(response.Data.ReturnData[0]))
	})
	t.Run("with finality check, chain is stuck", func(t *testing.T) {
		httpClient := &mockHTTPClient{
			doCalled: func(req *http.Request) (*http.Response, error) {
				response, handled, err := handleRequestNetworkConfigAndStatus(req, 3, 9170528, 9170526)
				if handled {
					return response, err
				}

				assert.Fail(t, "should have not reached this point in which the VM query is actually requested")
				return nil, nil
			},
		}
		args := createMockArgsProxy(httpClient)
		args.FinalityCheck = true
		ep, _ := NewProxy(args)

		response, err := ep.ExecuteVMQuery(context.Background(), &data.VmValueRequest{
			Address:    "erd1qqqqqqqqqqqqqpgqxwakt2g7u9atsnr03gqcgmhcv38pt7mkd94q6shuwt",
			FuncName:   "version",
			CallerAddr: "erd1qqqqqqqqqqqqqpgqxwakt2g7u9atsnr03gqcgmhcv38pt7mkd94q6shuwt",
		})

		assert.NotNil(t, err)
		assert.True(t, strings.Contains(err.Error(), "shardID 2 is stuck"))
		assert.Nil(t, response)
	})
	t.Run("with finality check, invalid address", func(t *testing.T) {
		httpClient := &mockHTTPClient{
			doCalled: func(req *http.Request) (*http.Response, error) {
				response, handled, err := handleRequestNetworkConfigAndStatus(req, 3, 9170526, 9170526)
				if handled {
					return response, err
				}

				assert.Fail(t, "should have not reached this point in which the VM query is actually requested")
				return nil, nil
			},
		}
		args := createMockArgsProxy(httpClient)
		args.FinalityCheck = true
		ep, _ := NewProxy(args)

		response, err := ep.ExecuteVMQuery(context.Background(), &data.VmValueRequest{
			Address:    "invalid",
			FuncName:   "version",
			CallerAddr: "erd1qqqqqqqqqqqqqpgqxwakt2g7u9atsnr03gqcgmhcv38pt7mkd94q6shuwt",
		})

		assert.NotNil(t, err)
		assert.True(t, strings.Contains(err.Error(), "invalid bech32 string length 7"))
		assert.Nil(t, response)
	})
	t.Run("with finality check, should work", func(t *testing.T) {
		wasHandled := false
		httpClient := &mockHTTPClient{
			doCalled: func(req *http.Request) (*http.Response, error) {
				response, handled, err := handleRequestNetworkConfigAndStatus(req, 3, 9170526, 9170525)
				if handled {
					wasHandled = true
					return response, err
				}

				return &http.Response{
					Body:       io.NopCloser(bytes.NewReader(responseBytes)),
					StatusCode: http.StatusOK,
				}, nil
			},
		}
		args := createMockArgsProxy(httpClient)
		args.FinalityCheck = true
		ep, _ := NewProxy(args)

		response, err := ep.ExecuteVMQuery(context.Background(), &data.VmValueRequest{
			Address:    "erd1qqqqqqqqqqqqqpgqxwakt2g7u9atsnr03gqcgmhcv38pt7mkd94q6shuwt",
			FuncName:   "version",
			CallerAddr: "erd1qqqqqqqqqqqqqpgqxwakt2g7u9atsnr03gqcgmhcv38pt7mkd94q6shuwt",
		})

		assert.True(t, wasHandled)
		require.Nil(t, err)
		require.Equal(t, "0.5.5", string(response.Data.ReturnData[0]))
	})
}

func TestProxy_GetRawBlockByHash(t *testing.T) {
	t.Parallel()

	expectedTs := &testStruct{
		Nonce: 1,
		Name:  "a test struct to be sent and received",
	}
	responseBytes, _ := json.Marshal(expectedTs)
	rawBlockData := &data.RawBlockRespone{
		Data: struct {
			Block []byte "json:\"block\""
		}{
			Block: responseBytes,
		},
		Error: "",
		Code:  "200",
	}
	rawBlockDataBytes, _ := json.Marshal(rawBlockData)

	httpClient := createMockClientRespondingBytes(rawBlockDataBytes)
	args := createMockArgsProxy(httpClient)
	ep, _ := NewProxy(args)

	response, err := ep.GetRawBlockByHash(context.Background(), 0, "aaaa")
	require.Nil(t, err)

	ts := &testStruct{}
	err = json.Unmarshal(response, ts)
	require.Nil(t, err)
	require.Equal(t, expectedTs, ts)
}

func TestProxy_GetRawBlockByNonce(t *testing.T) {
	t.Parallel()

	expectedTs := &testStruct{
		Nonce: 10,
		Name:  "a test struct to be sent and received",
	}
	responseBytes, _ := json.Marshal(expectedTs)
	rawBlockData := &data.RawBlockRespone{
		Data: struct {
			Block []byte "json:\"block\""
		}{
			Block: responseBytes,
		},
	}
	rawBlockDataBytes, _ := json.Marshal(rawBlockData)

	httpClient := createMockClientRespondingBytes(rawBlockDataBytes)
	args := createMockArgsProxy(httpClient)
	ep, _ := NewProxy(args)

	response, err := ep.GetRawBlockByNonce(context.Background(), 0, 10)
	require.Nil(t, err)

	ts := &testStruct{}
	err = json.Unmarshal(response, ts)
	require.Nil(t, err)
	require.Equal(t, expectedTs, ts)
}

func TestProxy_GetRawMiniBlockByHash(t *testing.T) {
	t.Parallel()

	expectedTs := &testStruct{
		Nonce: 10,
		Name:  "a test struct to be sent and received",
	}
	responseBytes, _ := json.Marshal(expectedTs)
	rawBlockData := &data.RawMiniBlockRespone{
		Data: struct {
			MiniBlock []byte "json:\"miniblock\""
		}{
			MiniBlock: responseBytes,
		},
	}
	rawBlockDataBytes, _ := json.Marshal(rawBlockData)

	httpClient := createMockClientRespondingBytes(rawBlockDataBytes)
	args := createMockArgsProxy(httpClient)
	ep, _ := NewProxy(args)

	response, err := ep.GetRawMiniBlockByHash(context.Background(), 0, "aaaa", 1)
	require.Nil(t, err)

	ts := &testStruct{}
	err = json.Unmarshal(response, ts)
	require.Nil(t, err)
	require.Equal(t, expectedTs, ts)
}

func TestProxy_GetNonceAtEpochStart(t *testing.T) {
	t.Parallel()

	expectedNonce := uint64(2)
	expectedNetworkStatus := &data.NetworkStatus{
		NonceAtEpochStart: expectedNonce,
		ShardID:           core.MetachainShardId,
	}
	statusResponse := &data.NodeStatusResponse{
		Data: struct {
			Status *data.NetworkStatus "json:\"metrics\""
		}{
			Status: expectedNetworkStatus,
		},
	}
	statusResponseBytes, _ := json.Marshal(statusResponse)

	httpClient := createMockClientRespondingBytes(statusResponseBytes)
	args := createMockArgsProxy(httpClient)
	ep, _ := NewProxy(args)

	response, err := ep.GetNonceAtEpochStart(context.Background(), core.MetachainShardId)
	require.Nil(t, err)

	require.Equal(t, expectedNonce, response)
}

func TestProxy_GetRatingsConfig(t *testing.T) {
	t.Parallel()

	expectedRatingsConfig := &data.RatingsConfig{
		GeneralMaxRating: 0,
		GeneralMinRating: 0,
	}
	ratingsResponse := &data.RatingsConfigResponse{
		Data: struct {
			Config *data.RatingsConfig "json:\"config\""
		}{
			Config: expectedRatingsConfig,
		},
	}
	ratingsResponseBytes, _ := json.Marshal(ratingsResponse)

	httpClient := createMockClientRespondingBytes(ratingsResponseBytes)
	args := createMockArgsProxy(httpClient)
	ep, _ := NewProxy(args)

	response, err := ep.GetRatingsConfig(context.Background())
	require.Nil(t, err)

	require.Equal(t, expectedRatingsConfig, response)
}

func TestProxy_GetEnableEpochsConfig(t *testing.T) {
	t.Parallel()

	enableEpochs := config.EnableEpochs{
		BalanceWaitingListsEnableEpoch: 1,
	}

	expectedEnableEpochsConfig := &data.EnableEpochsConfig{
		EnableEpochs: enableEpochs,
	}
	enableEpochsResponse := &data.EnableEpochsConfigResponse{
		Data: struct {
			Config *data.EnableEpochsConfig "json:\"enableEpochs\""
		}{
			Config: expectedEnableEpochsConfig},
	}
	enableEpochsResponseBytes, _ := json.Marshal(enableEpochsResponse)

	httpClient := createMockClientRespondingBytes(enableEpochsResponseBytes)
	args := createMockArgsProxy(httpClient)
	ep, _ := NewProxy(args)

	response, err := ep.GetEnableEpochsConfig(context.Background())
	require.Nil(t, err)

	require.Equal(t, expectedEnableEpochsConfig, response)
}

func TestProxy_GetGenesisNodesPubKeys(t *testing.T) {
	t.Parallel()

	expectedGenesisNodes := &data.GenesisNodes{
		Eligible: map[uint32][]string{
			0: {"pubkey1"},
		},
	}
	genesisNodesResponse := &data.GenesisNodesResponse{
		Data: struct {
			Nodes *data.GenesisNodes "json:\"nodes\""
		}{
			Nodes: expectedGenesisNodes},
	}
	genesisNodesResponseBytes, _ := json.Marshal(genesisNodesResponse)

	httpClient := createMockClientRespondingBytes(genesisNodesResponseBytes)
	args := createMockArgsProxy(httpClient)
	ep, _ := NewProxy(args)

	response, err := ep.GetGenesisNodesPubKeys(context.Background())
	require.Nil(t, err)

	require.Equal(t, expectedGenesisNodes, response)
}

func TestProxy_GetValidatorsInfoByEpoch(t *testing.T) {
	t.Parallel()

	expectedValidatorsInfo := []*state.ShardValidatorInfo{
		{
			PublicKey: []byte("pubkey1"),
			ShardId:   1,
		},
	}
	validatorsInfoResponse := &data.ValidatorsInfoResponse{
		Data: struct {
			ValidatorsInfo []*state.ShardValidatorInfo "json:\"validators\""
		}{
			ValidatorsInfo: expectedValidatorsInfo,
		},
	}
	validatorsInfoResponseBytes, _ := json.Marshal(validatorsInfoResponse)

	httpClient := createMockClientRespondingBytes(validatorsInfoResponseBytes)
	args := createMockArgsProxy(httpClient)
	ep, _ := NewProxy(args)

	response, err := ep.GetValidatorsInfoByEpoch(context.Background(), 1)
	require.Nil(t, err)

	require.Equal(t, expectedValidatorsInfo, response)
}

func TestElrondProxy_GetESDTTokenData(t *testing.T) {
	t.Parallel()

	token := "TKN-001122"
	expectedErr := errors.New("expected error")
	validAddress := data.NewAddressFromBytes(bytes.Repeat([]byte("1"), 32))
	emptyQueryOptions := api.AccountQueryOptions{}
	t.Run("nil address, should error", func(t *testing.T) {
		t.Parallel()

		httpClient := createMockClientRespondingBytes(make([]byte, 0))
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		tokenData, err := ep.GetESDTTokenData(context.Background(), nil, token, emptyQueryOptions)
		assert.Nil(t, tokenData)
		assert.Equal(t, ErrNilAddress, err)
	})
	t.Run("invalid address, should error", func(t *testing.T) {
		t.Parallel()

		httpClient := createMockClientRespondingBytes(make([]byte, 0))
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		address := data.NewAddressFromBytes([]byte("invalid"))
		tokenData, err := ep.GetESDTTokenData(context.Background(), address, token, emptyQueryOptions)
		assert.Nil(t, tokenData)
		assert.Equal(t, ErrInvalidAddress, err)
	})
	t.Run("http client errors, should error", func(t *testing.T) {
		t.Parallel()

		httpClient := createMockClientRespondingError(expectedErr)
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		tokenData, err := ep.GetESDTTokenData(context.Background(), validAddress, token, emptyQueryOptions)
		assert.Nil(t, tokenData)
		assert.ErrorIs(t, err, expectedErr)
	})
	t.Run("invalid status, should error", func(t *testing.T) {
		t.Parallel()

		httpClient := createMockClientRespondingBytesWithStatus(make([]byte, 0), http.StatusNotFound)
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		tokenData, err := ep.GetESDTTokenData(context.Background(), validAddress, token, emptyQueryOptions)
		assert.Nil(t, tokenData)
		assert.ErrorIs(t, err, ErrHTTPStatusCodeIsNotOK)
	})
	t.Run("invalid response bytes, should error", func(t *testing.T) {
		t.Parallel()

		httpClient := createMockClientRespondingBytes([]byte("invalid json"))
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		tokenData, err := ep.GetESDTTokenData(context.Background(), validAddress, token, emptyQueryOptions)
		assert.Nil(t, tokenData)
		assert.NotNil(t, err)
	})
	t.Run("response returned error, should error", func(t *testing.T) {
		t.Parallel()

		response := &data.ESDTFungibleResponse{
			Error: expectedErr.Error(),
		}
		responseBytes, _ := json.Marshal(response)

		httpClient := createMockClientRespondingBytes(responseBytes)
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		tokenData, err := ep.GetESDTTokenData(context.Background(), validAddress, token, emptyQueryOptions)
		assert.Nil(t, tokenData)
		assert.NotNil(t, err)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		responseTokenData := &data.ESDTFungibleTokenData{
			TokenIdentifier: "identifier",
			Balance:         "balance",
			Properties:      "properties",
		}
		response := &data.ESDTFungibleResponse{
			Data: struct {
				TokenData *data.ESDTFungibleTokenData `json:"tokenData"`
			}{
				TokenData: responseTokenData,
			},
		}
		responseBytes, _ := json.Marshal(response)

		httpClient := createMockClientRespondingBytes(responseBytes)
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		tokenData, err := ep.GetESDTTokenData(context.Background(), validAddress, token, emptyQueryOptions)
		assert.NotNil(t, tokenData)
		assert.Nil(t, err)
		assert.Equal(t, responseTokenData, tokenData)
		assert.False(t, responseTokenData == tokenData) // pointer testing
	})
	t.Run("should work with query options", func(t *testing.T) {
		t.Parallel()

		responseTokenData := &data.ESDTFungibleTokenData{
			TokenIdentifier: "identifier",
			Balance:         "balance",
			Properties:      "properties",
		}
		response := &data.ESDTFungibleResponse{
			Data: struct {
				TokenData *data.ESDTFungibleTokenData `json:"tokenData"`
			}{
				TokenData: responseTokenData,
			},
		}
		responseBytes, _ := json.Marshal(response)
		expectedSuffix := "?blockHash=626c6f636b2068617368&blockNonce=3838&blockRootHash=626c6f636b20726f6f742068617368&hintEpoch=3939&onFinalBlock=true&onStartOfEpoch=3737"

		httpClient := &mockHTTPClient{
			doCalled: func(req *http.Request) (*http.Response, error) {
				assert.True(t, strings.HasSuffix(req.URL.String(), expectedSuffix))

				return &http.Response{
					Body:       io.NopCloser(bytes.NewReader(responseBytes)),
					StatusCode: http.StatusOK,
				}, nil
			},
		}
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		tokenData, err := ep.GetESDTTokenData(context.Background(), validAddress, token, testQueryOptions)
		assert.NotNil(t, tokenData)
		assert.Nil(t, err)
		assert.Equal(t, responseTokenData, tokenData)
		assert.False(t, responseTokenData == tokenData) // pointer testing
	})
}

func TestElrondProxy_GetNFTTokenData(t *testing.T) {
	t.Parallel()

	token := "TKN-001122"
	nonce := uint64(37)
	expectedErr := errors.New("expected error")
	validAddress := data.NewAddressFromBytes(bytes.Repeat([]byte("1"), 32))
	emptyQueryOptions := api.AccountQueryOptions{}
	t.Run("nil address, should error", func(t *testing.T) {
		t.Parallel()

		httpClient := createMockClientRespondingBytes(make([]byte, 0))
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		tokenData, err := ep.GetNFTTokenData(context.Background(), nil, token, nonce, emptyQueryOptions)
		assert.Nil(t, tokenData)
		assert.Equal(t, ErrNilAddress, err)
	})
	t.Run("invalid address, should error", func(t *testing.T) {
		t.Parallel()

		httpClient := createMockClientRespondingBytes(make([]byte, 0))
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		address := data.NewAddressFromBytes([]byte("invalid"))
		tokenData, err := ep.GetNFTTokenData(context.Background(), address, token, nonce, emptyQueryOptions)
		assert.Nil(t, tokenData)
		assert.Equal(t, ErrInvalidAddress, err)
	})
	t.Run("http client errors, should error", func(t *testing.T) {
		t.Parallel()

		httpClient := createMockClientRespondingError(expectedErr)
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		tokenData, err := ep.GetNFTTokenData(context.Background(), validAddress, token, nonce, emptyQueryOptions)
		assert.Nil(t, tokenData)
		assert.ErrorIs(t, err, expectedErr)
	})
	t.Run("invalid status, should error", func(t *testing.T) {
		t.Parallel()

		httpClient := createMockClientRespondingBytesWithStatus(make([]byte, 0), http.StatusNotFound)
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		tokenData, err := ep.GetNFTTokenData(context.Background(), validAddress, token, nonce, emptyQueryOptions)
		assert.Nil(t, tokenData)
		assert.ErrorIs(t, err, ErrHTTPStatusCodeIsNotOK)
	})
	t.Run("invalid response bytes, should error", func(t *testing.T) {
		t.Parallel()

		httpClient := createMockClientRespondingBytes([]byte("invalid json"))
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		tokenData, err := ep.GetNFTTokenData(context.Background(), validAddress, token, nonce, emptyQueryOptions)
		assert.Nil(t, tokenData)
		assert.NotNil(t, err)
	})
	t.Run("response returned error, should error", func(t *testing.T) {
		t.Parallel()

		response := &data.ESDTNFTResponse{
			Error: expectedErr.Error(),
		}
		responseBytes, _ := json.Marshal(response)

		httpClient := createMockClientRespondingBytes(responseBytes)
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		tokenData, err := ep.GetNFTTokenData(context.Background(), validAddress, token, nonce, emptyQueryOptions)
		assert.Nil(t, tokenData)
		assert.NotNil(t, err)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		responseTokenData := &data.ESDTNFTTokenData{
			TokenIdentifier: "identifier",
			Balance:         "balance",
			Properties:      "properties",
			Name:            "name",
			Nonce:           nonce,
			Creator:         "creator",
			Royalties:       "royalties",
			Hash:            []byte("hash"),
			URIs:            [][]byte{[]byte("uri1"), []byte("uri2")},
			Attributes:      []byte("attributes"),
		}
		response := &data.ESDTNFTResponse{
			Data: struct {
				TokenData *data.ESDTNFTTokenData `json:"tokenData"`
			}{
				TokenData: responseTokenData,
			},
		}
		responseBytes, _ := json.Marshal(response)

		httpClient := createMockClientRespondingBytes(responseBytes)
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		tokenData, err := ep.GetNFTTokenData(context.Background(), validAddress, token, nonce, emptyQueryOptions)
		assert.NotNil(t, tokenData)
		assert.Nil(t, err)
		assert.Equal(t, responseTokenData, tokenData)
		assert.False(t, responseTokenData == tokenData) // pointer testing
	})
	t.Run("should work with query options", func(t *testing.T) {
		t.Parallel()

		responseTokenData := &data.ESDTNFTTokenData{
			TokenIdentifier: "identifier",
			Balance:         "balance",
			Properties:      "properties",
			Name:            "name",
			Nonce:           nonce,
			Creator:         "creator",
			Royalties:       "royalties",
			Hash:            []byte("hash"),
			URIs:            [][]byte{[]byte("uri1"), []byte("uri2")},
			Attributes:      []byte("attributes"),
		}
		response := &data.ESDTNFTResponse{
			Data: struct {
				TokenData *data.ESDTNFTTokenData `json:"tokenData"`
			}{
				TokenData: responseTokenData,
			},
		}
		responseBytes, _ := json.Marshal(response)
		expectedSuffix := "?blockHash=626c6f636b2068617368&blockNonce=3838&blockRootHash=626c6f636b20726f6f742068617368&hintEpoch=3939&onFinalBlock=true&onStartOfEpoch=3737"

		httpClient := &mockHTTPClient{
			doCalled: func(req *http.Request) (*http.Response, error) {
				assert.True(t, strings.HasSuffix(req.URL.String(), expectedSuffix))

				return &http.Response{
					Body:       io.NopCloser(bytes.NewReader(responseBytes)),
					StatusCode: http.StatusOK,
				}, nil
			},
		}
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		tokenData, err := ep.GetNFTTokenData(context.Background(), validAddress, token, nonce, testQueryOptions)
		assert.NotNil(t, tokenData)
		assert.Nil(t, err)
		assert.Equal(t, responseTokenData, tokenData)
		assert.False(t, responseTokenData == tokenData) // pointer testing
	})
}

func TestProxy_GetGuardianData(t *testing.T) {
	t.Parallel()

	t.Run("nil address should error", func(t *testing.T) {
		t.Parallel()

		httpClient := createMockClientRespondingBytes([]byte("dummy response"))
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		response, err := ep.GetGuardianData(context.Background(), nil)
		require.Equal(t, err, ErrNilAddress)
		require.Nil(t, response)
	})
	t.Run("invalid address should error", func(t *testing.T) {
		t.Parallel()

		httpClient := createMockClientRespondingBytes([]byte("dummy response"))
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		address := data.NewAddressFromBytes([]byte("invalid address"))
		response, err := ep.GetGuardianData(context.Background(), address)
		require.Equal(t, err, ErrInvalidAddress)
		require.Nil(t, response)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		expectedGuardianData := &api.GuardianData{
			ActiveGuardian: &api.Guardian{
				Address:         "active guardian",
				ActivationEpoch: 100,
			},
			PendingGuardian: &api.Guardian{
				Address:         "pending guardian",
				ActivationEpoch: 200,
			},
			Guarded: false,
		}
		guardianDataResponse := &data.GuardianDataResponse{
			Data: struct {
				GuardianData *api.GuardianData `json:"guardianData"`
			}{
				GuardianData: expectedGuardianData,
			},
		}
		guardianDataResponseBytes, _ := json.Marshal(guardianDataResponse)

		httpClient := createMockClientRespondingBytes(guardianDataResponseBytes)
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		address, _ := data.NewAddressFromBech32String("erd1qqqqqqqqqqqqqpgqfzydqmdw7m2vazsp6u5p95yxz76t2p9rd8ss0zp9ts")
		response, err := ep.GetGuardianData(context.Background(), address)
		require.Nil(t, err)

		require.Equal(t, expectedGuardianData, response)
	})
}

func TestProxy_IsDataTrieMigrated(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("expected error")
	validAddress := data.NewAddressFromBytes(bytes.Repeat([]byte("1"), 32))
	t.Run("nil address", func(t *testing.T) {
		t.Parallel()

		httpClient := createMockClientRespondingBytes(make([]byte, 0))
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		isMigrated, err := ep.IsDataTrieMigrated(context.Background(), nil)
		assert.False(t, isMigrated)
		assert.Equal(t, ErrNilAddress, err)
	})
	t.Run("invalid bech32 address", func(t *testing.T) {
		t.Parallel()

		httpClient := createMockClientRespondingBytes(make([]byte, 0))
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		invalidAddress := data.NewAddressFromBytes([]byte("invalid"))

		isMigrated, err := ep.IsDataTrieMigrated(context.Background(), invalidAddress)
		assert.False(t, isMigrated)
		assert.True(t, strings.Contains(err.Error(), "wrong size when encoding address"))
	})
	t.Run("http client errors", func(t *testing.T) {
		t.Parallel()

		httpClient := createMockClientRespondingError(expectedErr)
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		isMigrated, err := ep.IsDataTrieMigrated(context.Background(), validAddress)
		assert.False(t, isMigrated)
		assert.ErrorIs(t, err, expectedErr)
	})
	t.Run("invalid status", func(t *testing.T) {
		t.Parallel()

		httpClient := createMockClientRespondingBytesWithStatus(make([]byte, 0), http.StatusNotFound)
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		isMigrated, err := ep.IsDataTrieMigrated(context.Background(), validAddress)
		assert.False(t, isMigrated)
		assert.ErrorIs(t, err, ErrHTTPStatusCodeIsNotOK)
	})
	t.Run("invalid response bytes", func(t *testing.T) {
		t.Parallel()

		httpClient := createMockClientRespondingBytes([]byte("invalid json"))
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		isMigrated, err := ep.IsDataTrieMigrated(context.Background(), validAddress)
		assert.False(t, isMigrated)
		assert.NotNil(t, err)
	})
	t.Run("response returned error", func(t *testing.T) {
		t.Parallel()

		response := &data.IsDataTrieMigratedResponse{
			Error: expectedErr.Error(),
		}
		responseBytes, _ := json.Marshal(response)

		httpClient := createMockClientRespondingBytes(responseBytes)
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		isMigrated, err := ep.IsDataTrieMigrated(context.Background(), validAddress)
		assert.False(t, isMigrated)
		assert.NotNil(t, err)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})
	t.Run("isMigrated key not found in map", func(t *testing.T) {
		t.Parallel()

		responseMap := make(map[string]bool)
		responseMap["random key"] = true

		response := &data.IsDataTrieMigratedResponse{
			Data: responseMap,
		}
		responseBytes, _ := json.Marshal(response)

		httpClient := createMockClientRespondingBytes(responseBytes)
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		isMigrated, err := ep.IsDataTrieMigrated(context.Background(), validAddress)
		assert.False(t, isMigrated)
		assert.Contains(t, err.Error(), "isMigrated key not found in response map")
	})
	t.Run("migrated trie", func(t *testing.T) {
		t.Parallel()

		responseMap := make(map[string]bool)
		responseMap["isMigrated"] = true

		response := &data.IsDataTrieMigratedResponse{
			Data: responseMap,
		}
		responseBytes, _ := json.Marshal(response)

		httpClient := createMockClientRespondingBytes(responseBytes)
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		isMigrated, err := ep.IsDataTrieMigrated(context.Background(), validAddress)
		assert.True(t, isMigrated)
		assert.Nil(t, err)
	})
	t.Run("not migrated trie", func(t *testing.T) {
		t.Parallel()

		responseMap := make(map[string]bool)
		responseMap["isMigrated"] = false

		response := &data.IsDataTrieMigratedResponse{
			Data: responseMap,
		}
		responseBytes, _ := json.Marshal(response)

		httpClient := createMockClientRespondingBytes(responseBytes)
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		isMigrated, err := ep.IsDataTrieMigrated(context.Background(), validAddress)
		assert.False(t, isMigrated)
		assert.Nil(t, err)
	})
}

func TestProxy_FilterLogs(t *testing.T) {
	t.Parallel()

	httpDataBlock21000000 := `{"data":{"block":{"nonce":21000000,"round":21008275,"epoch":1458,"shard":0,"numTxs":2,"hash":"94f7d0e806c5ff3768236f78c67007c82d36a2d917657ae265409a07ff835d8e","prevBlockHash":"7d74f44411c7ceacc0a78ac10e1da07b4b5e5bf54dd2570d8771df0666906e7b","stateRootHash":"6991a3c0d0f491ab466e5a7a6ec7c82a82fce07d62868de40f2aa7cce91b2f14","accumulatedFees":"0","developerFees":"0","status":"on-chain","randSeed":"d45aa08f56aa025f7e6f1d6b6239f0a18a95624ba6378c37dd609cb2b92da4564e33f87ab706d8f0f8a06ba3494e7891","prevRandSeed":"c90e3d77f157df90ea9237e76d9e3f09d923eed11321b72b6648b86ea07241d6b4e6fdc98ace4a31dcf9d9d95097f318","pubKeyBitmap":"ffffffffffffff7f","signature":"35099a6af1ef34d11f3e348dba9634072b8ef2ed90c9a04bb8dc028e2198684c43a9c8922af772bcea80c3419a6c9016","leaderSignature":"ebc7eccc88d50b8ab56ab6a503b75a2b7ccde2bcd02ca8f2d474cc5769b8d48e5f619d39c9ba95fbf80f9cda760d3c99","chainID":"1","softwareVersion":"32","receiptsHash":"0e5751c026e543b2e8ab2eb06099daa1d1e5df47778f7787faab45cdf12fe3a8","timestamp":1722167250,"miniBlocks":[{"hash":"13b5b686e27dca2c3cbd131837ba8050639970661717a4fef8bde968f738f766","type":"SmartContractResultBlock","processingType":"Normal","constructionState":"Final","sourceShard":1,"destinationShard":0,"transactions":[{"type":"unsigned","processingTypeOnSource":"BuiltInFunctionCall","processingTypeOnDestination":"BuiltInFunctionCall","hash":"564da8da1fcb85509f9b1fe73a5c0e2d18362265b2ea1f322baa7c255163bcd3","nonce":0,"round":21008275,"epoch":1458,"value":"0","receiver":"erd1mmh2j5y8esv2tmmyeau3hr4xa2u2te3zc3j9wumn3v5vm8uvsczqnucj5l","sender":"erd1qqqqqqqqqqqqqpgqta0tv8d5pjzmwzshrtw62n4nww9kxtl278ssspxpxu","gasPrice":1000000000,"data":"RVNEVFRyYW5zZmVyQDQ4NTU1NDRiMmQzNDY2NjEzNDYyMzJANThhODZmYzg0MzQw","previousTransactionHash":"be72a85d881e83a22196715a54f5897f40210c1e62162014f47ab1be89e7702c","originalTransactionHash":"be72a85d881e83a22196715a54f5897f40210c1e62162014f47ab1be89e7702c","originalSender":"erd1mmh2j5y8esv2tmmyeau3hr4xa2u2te3zc3j9wumn3v5vm8uvsczqnucj5l","sourceShard":1,"destinationShard":0,"miniblockType":"SmartContractResultBlock","miniblockHash":"13b5b686e27dca2c3cbd131837ba8050639970661717a4fef8bde968f738f766","logs":{"address":"erd1mmh2j5y8esv2tmmyeau3hr4xa2u2te3zc3j9wumn3v5vm8uvsczqnucj5l","events":[{"address":"erd1qqqqqqqqqqqqqpgqta0tv8d5pjzmwzshrtw62n4nww9kxtl278ssspxpxu","identifier":"ESDTTransfer","topics":["SFVUSy00ZmE0YjI=","","WKhvyENA","3u6pUIfMGKXvZM95G46m6ril5iLEZFdzc4sozZ+MhgQ="],"data":null,"additionalData":["","RVNEVFRyYW5zZmVy","SFVUSy00ZmE0YjI=","WKhvyENA"]},{"address":"erd1mmh2j5y8esv2tmmyeau3hr4xa2u2te3zc3j9wumn3v5vm8uvsczqnucj5l","identifier":"writeLog","topics":["AAAAAAAAAAAFAF9eth20DIW3Chca3aVOs3OLYy/q8eE="],"data":"QDZmNmI=","additionalData":["QDZmNmI="]},{"address":"erd1mmh2j5y8esv2tmmyeau3hr4xa2u2te3zc3j9wumn3v5vm8uvsczqnucj5l","identifier":"completedTxEvent","topics":["vnKoXYgeg6IhlnFaVPWJf0AhDB5iFiAU9HqxvonncCw="],"data":null,"additionalData":null}]},"status":"success","tokens":["HUTK-4fa4b2"],"esdtValues":["97480453145408"],"operation":"ESDTTransfer","callType":"directCall","options":0},{"type":"unsigned","processingTypeOnSource":"MoveBalance","processingTypeOnDestination":"MoveBalance","hash":"43a29c450a306a8e82c40f3a453892a8827d986367b95d5e7768f35e1b582c26","nonce":30335,"round":21008275,"epoch":1458,"value":"77706790000000","receiver":"erd1mmh2j5y8esv2tmmyeau3hr4xa2u2te3zc3j9wumn3v5vm8uvsczqnucj5l","sender":"erd1qqqqqqqqqqqqqpgqta0tv8d5pjzmwzshrtw62n4nww9kxtl278ssspxpxu","gasPrice":1000000000,"data":"QDZmNmJAMDAwMDAwMGI0ODU1NTQ0YjJkMzQ2NjYxMzQ2MjMyMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDA2NThhODZmYzg0MzQw","previousTransactionHash":"be72a85d881e83a22196715a54f5897f40210c1e62162014f47ab1be89e7702c","originalTransactionHash":"be72a85d881e83a22196715a54f5897f40210c1e62162014f47ab1be89e7702c","sourceShard":1,"destinationShard":0,"miniblockType":"SmartContractResultBlock","miniblockHash":"13b5b686e27dca2c3cbd131837ba8050639970661717a4fef8bde968f738f766","logs":{"address":"erd1mmh2j5y8esv2tmmyeau3hr4xa2u2te3zc3j9wumn3v5vm8uvsczqnucj5l","events":[{"address":"erd1mmh2j5y8esv2tmmyeau3hr4xa2u2te3zc3j9wumn3v5vm8uvsczqnucj5l","identifier":"completedTxEvent","topics":["vnKoXYgeg6IhlnFaVPWJf0AhDB5iFiAU9HqxvonncCw="],"data":null,"additionalData":null}]},"status":"success","operation":"transfer","isRefund":true,"callType":"directCall","options":0}],"indexOfFirstTxProcessed":0,"indexOfLastTxProcessed":1}],"scheduledData":{"rootHash":"a0f86f73d739d8490d07d3a47acdd19a14c80ecfe8247aca2877d4e06944cc8f","accumulatedFees":"0","developerFees":"0"}}},"error":"","code":"successful"}`
	httpDataBlock21000001 := `{"data":{"block":{"nonce":21000001,"round":21008276,"epoch":1458,"shard":0,"numTxs":1,"hash":"6151d0be315f423218bc3d5552607208101d872cb65be19b2e34b1840e3de0fd","prevBlockHash":"94f7d0e806c5ff3768236f78c67007c82d36a2d917657ae265409a07ff835d8e","stateRootHash":"27e7a29cf7ef5b4b6bda589c2b1d6d95e640a5d0f6218d50ee45227965ea21c2","accumulatedFees":"633302170000000","developerFees":"9639651000000","status":"on-chain","randSeed":"f631b3d7bb1793cec31b1094b483aed5e3a97f07decdd14fc7b67e942c6953e42bb0a81d7bb068d4fdfdc1d80aae1c13","prevRandSeed":"d45aa08f56aa025f7e6f1d6b6239f0a18a95624ba6378c37dd609cb2b92da4564e33f87ab706d8f0f8a06ba3494e7891","pubKeyBitmap":"ffffffffffffff7f","signature":"a4bf779974fc7bd1a5545702db0137de7968fe9b0f938f6e6cf1389b0e852273a1bfaf40aeef9972a96251164cb36a13","leaderSignature":"7bca4323d89918f8ba63d56ca146ebb2956cef9e9495c63a367a905b496f965641422afd29874ba259206905c097068d","chainID":"1","softwareVersion":"32","receiptsHash":"ecc180171f7fa8b8be10a13946f8087e215b4abd2e4acef67354e4cd077f711d","timestamp":1722167256,"miniBlocks":[{"hash":"10ed7dbead935b88975cd545ad8a69b57c4c96f4e2eb44e3f90965e9a964164a","type":"TxBlock","processingType":"Normal","constructionState":"Final","sourceShard":0,"destinationShard":0,"transactions":[{"type":"normal","processingTypeOnSource":"RelayedTxV2","processingTypeOnDestination":"RelayedTxV2","hash":"7918af2e7eac08999938e164896b6eb41b0255668c7bbcc79274c769c6870b7c","nonce":787136,"round":21008276,"epoch":1458,"value":"0","receiver":"erd15aq4rug5rxjnu88723f2y5fx2w9kzzw2rha89jzfpfhfy9huuxyq4zrm0t","sender":"erd1y9upyp2e4udn9ehxfam5l7gd5upenan9skg3mfm85f36qsvrpr5q4sgqfv","gasPrice":1000000000,"gasLimit":5597500,"data":"cmVsYXllZFR4VjJAMDAwMDAwMDAwMDAwMDAwMDA1MDA4MmY2ODE2ZGM5ZGQwMTEyNDQ2ZjQ1OThmNDg2MzI3ZmIxNzIzMzQ1YmQ1OEAxY0A0NTUzNDQ1NDU0NzI2MTZlNzM2NjY1NzI0MDM0MzMzNDM3MzQ2NjMyNjQzMzM1MzYzNTMzMzkzMzM1MzMzMjMzMzg0MDMwMzIzMDMyNjY2NTY2NjI2NjMyNjQzNzYzMzI2NjMwMzAzMDMwMzA0MDM2MzIzNzM1MzczOTM0NjYzNjM2MzYzNjM2MzUzNzMyNDAzMDM0NjEzODM5NjFAN2IxNGIyNzM4ZTQ2ODQzMzAwOTZlMjA0YjM4NmU4MmZkODk3YWFlNTQwODkyZTk0ZDMwZDA5YTIwZmJiZjM2M2U0NzM2ZDczM2E0NTg2MmIyMjllODg3NGUyODVmYjhhNjFhNTRjMjZkYmExN2ZlMTcwMjZmNzRlZDdhYjIxMGY=","signature":"39d1798e2191d8191dc3256b4efc244c1fa822870915ba80ad9394bceab0e4291efe798b95ab939b9693d0aebef9a1df3b2cdead89c9d2055dc063745857aa00","sourceShard":0,"destinationShard":0,"miniblockType":"TxBlock","miniblockHash":"10ed7dbead935b88975cd545ad8a69b57c4c96f4e2eb44e3f90965e9a964164a","status":"success","tokens":["CGO-5e9528"],"esdtValues":["9500000000000000000000"],"receivers":["erd1qqqqqqqqqqqqqpgqstmgzmwfm5q3y3r0gkv0fp3j07chyv69h4vq7md7fd"],"receiversShardIDs":[0],"operation":"ESDTTransfer","function":"buyOffer","initiallyPaidFee":"647500000000000","isRelayed":true,"chainID":"1","version":1,"options":0}],"indexOfFirstTxProcessed":0,"indexOfLastTxProcessed":0},{"hash":"59afb0fe188a917202bb0d47d9c3a3fd83650a102fee1ff85c896343ecf1d39b","type":"SmartContractResultBlock","processingType":"Normal","isFromReceiptsStorage":true,"sourceShard":0,"destinationShard":0,"transactions":[{"type":"unsigned","processingTypeOnSource":"SCInvoking","processingTypeOnDestination":"SCInvoking","hash":"1f18fd88045146a6596f5bec22f7611a4d902011328cb5ea462d79e501e66ec1","nonce":0,"round":21008276,"epoch":1458,"value":"0","receiver":"erd1qqqqqqqqqqqqqpgqstmgzmwfm5q3y3r0gkv0fp3j07chyv69h4vq7md7fd","sender":"erd15aq4rug5rxjnu88723f2y5fx2w9kzzw2rha89jzfpfhfy9huuxyq4zrm0t","gasPrice":1000000000,"gasLimit":4633000,"data":"YnV5T2ZmZXJAMDRhODlh","previousTransactionHash":"13976575c92dbc52fbe512ea57d6f7bec42cfd9c07c069518cb7b8d64bc1b5b4","originalTransactionHash":"7918af2e7eac08999938e164896b6eb41b0255668c7bbcc79274c769c6870b7c","sourceShard":0,"destinationShard":0,"miniblockType":"SmartContractResultBlock","miniblockHash":"59afb0fe188a917202bb0d47d9c3a3fd83650a102fee1ff85c896343ecf1d39b","status":"success","operation":"transfer","function":"buyOffer","callType":"directCall","relayerAddress":"erd1y9upyp2e4udn9ehxfam5l7gd5upenan9skg3mfm85f36qsvrpr5q4sgqfv","relayedValue":"0","options":0},{"type":"unsigned","processingTypeOnSource":"BuiltInFunctionCall","processingTypeOnDestination":"BuiltInFunctionCall","hash":"6ef90ec8d57e3caeac23b3830fac9ce3b88caf6accbc8e876c02913f41959da4","nonce":0,"round":21008276,"epoch":1458,"value":"0","receiver":"erd1tlatx0xnpvn3smsskh04xaqcfwqtttqs3ryn862f9c73w5epr4uqg6jwgq","sender":"erd1qqqqqqqqqqqqqpgqstmgzmwfm5q3y3r0gkv0fp3j07chyv69h4vq7md7fd","gasPrice":1000000000,"data":"RVNEVFRyYW5zZmVyQDQzNDc0ZjJkMzU2NTM5MzUzMjM4QDAxZTkzZjA4ZjM4MDJjNjQwMDAw","previousTransactionHash":"13976575c92dbc52fbe512ea57d6f7bec42cfd9c07c069518cb7b8d64bc1b5b4","originalTransactionHash":"7918af2e7eac08999938e164896b6eb41b0255668c7bbcc79274c769c6870b7c","sourceShard":0,"destinationShard":0,"miniblockType":"SmartContractResultBlock","miniblockHash":"59afb0fe188a917202bb0d47d9c3a3fd83650a102fee1ff85c896343ecf1d39b","status":"success","tokens":["CGO-5e9528"],"esdtValues":["9025000000000000000000"],"operation":"ESDTTransfer","callType":"directCall","relayerAddress":"erd1y9upyp2e4udn9ehxfam5l7gd5upenan9skg3mfm85f36qsvrpr5q4sgqfv","relayedValue":"0","options":0},{"type":"unsigned","processingTypeOnSource":"BuiltInFunctionCall","processingTypeOnDestination":"BuiltInFunctionCall","hash":"1a69f8683e4fe1f2de99febd05c88f41f6f13bacc63cd8852b63e61e92d59d72","nonce":0,"round":21008276,"epoch":1458,"value":"0","receiver":"erd15aq4rug5rxjnu88723f2y5fx2w9kzzw2rha89jzfpfhfy9huuxyq4zrm0t","sender":"erd1qqqqqqqqqqqqqpgqstmgzmwfm5q3y3r0gkv0fp3j07chyv69h4vq7md7fd","gasPrice":1000000000,"data":"RVNEVE5GVFRyYW5zZmVyQDQzNDk1NDQ1NGQyZDYyNjQ2NjM1NjYzMUAwODA3N2JAMDFAYTc0MTUxZjExNDE5YTUzZTFjZmU1NDUyYTI1MTI2NTM4YjYxMDljYTFkZmE3MmM4NDkwYTZlOTIxNmZjZTE4OA==","previousTransactionHash":"13976575c92dbc52fbe512ea57d6f7bec42cfd9c07c069518cb7b8d64bc1b5b4","originalTransactionHash":"7918af2e7eac08999938e164896b6eb41b0255668c7bbcc79274c769c6870b7c","sourceShard":0,"destinationShard":0,"miniblockType":"SmartContractResultBlock","miniblockHash":"59afb0fe188a917202bb0d47d9c3a3fd83650a102fee1ff85c896343ecf1d39b","status":"success","tokens":["CITEM-bdf5f1-08077b"],"esdtValues":["1"],"receivers":["erd15aq4rug5rxjnu88723f2y5fx2w9kzzw2rha89jzfpfhfy9huuxyq4zrm0t"],"receiversShardIDs":[0],"operation":"ESDTNFTTransfer","callType":"directCall","relayerAddress":"erd1y9upyp2e4udn9ehxfam5l7gd5upenan9skg3mfm85f36qsvrpr5q4sgqfv","relayedValue":"0","options":0},{"type":"unsigned","processingTypeOnSource":"MoveBalance","processingTypeOnDestination":"MoveBalance","hash":"5a3e158e23bca6f245f28ba13af23d746e5d21b0de140cbb03dc09879e28981a","nonce":29,"round":21008276,"epoch":1458,"value":"14197830000000","receiver":"erd1y9upyp2e4udn9ehxfam5l7gd5upenan9skg3mfm85f36qsvrpr5q4sgqfv","sender":"erd1qqqqqqqqqqqqqpgqstmgzmwfm5q3y3r0gkv0fp3j07chyv69h4vq7md7fd","gasPrice":1000000000,"previousTransactionHash":"13976575c92dbc52fbe512ea57d6f7bec42cfd9c07c069518cb7b8d64bc1b5b4","originalTransactionHash":"7918af2e7eac08999938e164896b6eb41b0255668c7bbcc79274c769c6870b7c","returnMessage":"gas refund for relayer","sourceShard":0,"destinationShard":0,"miniblockType":"SmartContractResultBlock","miniblockHash":"59afb0fe188a917202bb0d47d9c3a3fd83650a102fee1ff85c896343ecf1d39b","status":"success","operation":"transfer","isRefund":true,"callType":"directCall","options":0},{"type":"unsigned","processingTypeOnSource":"BuiltInFunctionCall","processingTypeOnDestination":"SCInvoking","hash":"13976575c92dbc52fbe512ea57d6f7bec42cfd9c07c069518cb7b8d64bc1b5b4","nonce":28,"round":21008276,"epoch":1458,"value":"0","receiver":"erd1qqqqqqqqqqqqqpgqstmgzmwfm5q3y3r0gkv0fp3j07chyv69h4vq7md7fd","sender":"erd15aq4rug5rxjnu88723f2y5fx2w9kzzw2rha89jzfpfhfy9huuxyq4zrm0t","gasPrice":1000000000,"gasLimit":4833000,"data":"RVNEVFRyYW5zZmVyQDQzNDc0ZjJkMzU2NTM5MzUzMjM4QDAyMDJmZWZiZjJkN2MyZjAwMDAwQDYyNzU3OTRmNjY2NjY1NzJAMDRhODlh","previousTransactionHash":"7918af2e7eac08999938e164896b6eb41b0255668c7bbcc79274c769c6870b7c","originalTransactionHash":"7918af2e7eac08999938e164896b6eb41b0255668c7bbcc79274c769c6870b7c","sourceShard":0,"destinationShard":0,"miniblockType":"SmartContractResultBlock","miniblockHash":"59afb0fe188a917202bb0d47d9c3a3fd83650a102fee1ff85c896343ecf1d39b","logs":{"address":"erd1qqqqqqqqqqqqqpgqstmgzmwfm5q3y3r0gkv0fp3j07chyv69h4vq7md7fd","events":[{"address":"erd15aq4rug5rxjnu88723f2y5fx2w9kzzw2rha89jzfpfhfy9huuxyq4zrm0t","identifier":"ESDTTransfer","topics":["Q0dPLTVlOTUyOA==","","AgL++/LXwvAAAA==","AAAAAAAAAAAFAIL2gW3J3QESRG9FmPSGMn+xcjNFvVg="],"data":null,"additionalData":["","RVNEVFRyYW5zZmVy","Q0dPLTVlOTUyOA==","AgL++/LXwvAAAA==","YnV5T2ZmZXI=","BKia"]},{"address":"erd1qqqqqqqqqqqqqpgqstmgzmwfm5q3y3r0gkv0fp3j07chyv69h4vq7md7fd","identifier":"ESDTLocalBurn","topics":["Q0dPLTVlOTUyOA==","","CkzHmVY8OAAA"],"data":null,"additionalData":null},{"address":"erd1qqqqqqqqqqqqqpgqstmgzmwfm5q3y3r0gkv0fp3j07chyv69h4vq7md7fd","identifier":"ESDTTransfer","topics":["Q0dPLTVlOTUyOA==","","Aek/CPOALGQAAA==","X/qzPNMLJxhuELXfU3QYS4C1rBCIyTPpSS49F1MhHXg="],"data":"RGlyZWN0Q2FsbA==","additionalData":["RGlyZWN0Q2FsbA==","RVNEVFRyYW5zZmVy","Q0dPLTVlOTUyOA==","Aek/CPOALGQAAA=="]},{"address":"erd1qqqqqqqqqqqqqpgqstmgzmwfm5q3y3r0gkv0fp3j07chyv69h4vq7md7fd","identifier":"ESDTNFTTransfer","topics":["Q0lURU0tYmRmNWYx","CAd7","AQ==","p0FR8RQZpT4c/lRSolEmU4thCcod+nLISQpukhb84Yg="],"data":"RGlyZWN0Q2FsbA==","additionalData":["RGlyZWN0Q2FsbA==","RVNEVE5GVFRyYW5zZmVy","Q0lURU0tYmRmNWYx","CAd7","AQ==","p0FR8RQZpT4c/lRSolEmU4thCcod+nLISQpukhb84Yg="]},{"address":"erd1qqqqqqqqqqqqqpgqstmgzmwfm5q3y3r0gkv0fp3j07chyv69h4vq7md7fd","identifier":"buyOffer","topics":["YnV5X29mZmVy","BKia","Q0lURU0tYmRmNWYx","CAd7","Q0dPLTVlOTUyOA==","AgL++/LXwvAAAA==","X/qzPNMLJxhuELXfU3QYS4C1rBCIyTPpSS49F1MhHXg=","p0FR8RQZpT4c/lRSolEmU4thCcod+nLISQpukhb84Yg="],"data":null,"additionalData":[""]},{"address":"erd1qqqqqqqqqqqqqpgqstmgzmwfm5q3y3r0gkv0fp3j07chyv69h4vq7md7fd","identifier":"writeLog","topics":["p0FR8RQZpT4c/lRSolEmU4thCcod+nLISQpukhb84Yg="],"data":"QDZmNmJAMDRhODlh","additionalData":["QDZmNmJAMDRhODlh"]},{"address":"erd1qqqqqqqqqqqqqpgqstmgzmwfm5q3y3r0gkv0fp3j07chyv69h4vq7md7fd","identifier":"completedTxEvent","topics":["eRivLn6sCJmZOOFkiWtutBsCVWaMe7zHknTHacaHC3w="],"data":null,"additionalData":null}]},"status":"success","tokens":["CGO-5e9528"],"esdtValues":["9500000000000000000000"],"operation":"ESDTTransfer","function":"buyOffer","callType":"directCall","relayerAddress":"erd1y9upyp2e4udn9ehxfam5l7gd5upenan9skg3mfm85f36qsvrpr5q4sgqfv","relayedValue":"0","options":0}],"indexOfFirstTxProcessed":0,"indexOfLastTxProcessed":0}],"scheduledData":{"rootHash":"6991a3c0d0f491ab466e5a7a6ec7c82a82fce07d62868de40f2aa7cce91b2f14","accumulatedFees":"0","developerFees":"0"}}},"error":"","code":"successful"}`
	httpDataBlock21000005 := `{"data":{"block":{"nonce":21000005,"round":21008280,"epoch":1458,"shard":0,"numTxs":4,"hash":"98a546ea97ff90f36b3d18f1ef1cc26a60ccfd384f90a83f9dec145081970522","prevBlockHash":"85114eb317ac8efc9c7d9bee6ac5a33657faa34b1171b81dee0609633de655fd","stateRootHash":"7de197cd2a287ed2f34fc6a2558070dc7880ee2186ad7abe2aeac1742d88aff1","accumulatedFees":"1283137500000000","developerFees":"0","status":"on-chain","randSeed":"cf0c4f0172892c50d58eb03c47fae450c19e128c99b832e70a88916e19af76c1873bd16fe047823252280eb116b93c17","prevRandSeed":"ddaeedc6822b86a3a69d89238b7fb21514fd437ec670e39c088c73c7d7ded8c3dcd41ba86adee10c56e281112fdc678e","pubKeyBitmap":"ffffffffffffff7f","signature":"11dff3f68fb4789930875586ac34696410dd9de1734f72e8144fe99533c26e3be6d940b0d88b3700e1e8f805dd9eef10","leaderSignature":"00bcefd1b724a2194885dde6ec0334c5c9bc527dd6e64eb6394989345d2b66467de3ea6cb57f7f3cb2b27ab07548a385","chainID":"1","softwareVersion":"32","receiptsHash":"d1bff2f1e0a148164c619fb48ad8d215b6a1c8c06c3345cfc85662eb9ddf6d05","timestamp":1722167280,"miniBlocks":[{"hash":"9f9e6aca41da76e7f3c7b2078c30280b45b4496ad9c2390f2af3517e3eb8d051","type":"SmartContractResultBlock","processingType":"Normal","constructionState":"Final","sourceShard":1,"destinationShard":0,"transactions":[{"type":"unsigned","processingTypeOnSource":"BuiltInFunctionCall","processingTypeOnDestination":"BuiltInFunctionCall","hash":"abfe1307dd36ce8d884e7037f5e9d2a63d10618c962ce08822e755ecb102cfd5","nonce":0,"round":21008280,"epoch":1458,"value":"0","receiver":"erd1d7y4a8wtykxnxxjhywzk0q5tkey4g9z6rhalefw6syr779kh77yqd0fj5y","sender":"erd1qqqqqqqqqqqqqpgq50dge6rrpcra4tp9hl57jl0893a4r2r72jpsk39rjj","gasPrice":1000000000,"data":"TXVsdGlFU0RUTkZUVHJhbnNmZXJAMDNANTE1NzU0MmQzNDM2NjE2MzMwMzFAMDBAZTE4OWMyQDUzNDY0OTU0MmQ2MTY1NjI2MzM5MzBAMDBAMDFlNDViYmI2NGQxOGEwYjY4QDQzNTk0MjQ1NTIyZDM0MzgzOTYzMzE2M0AwMEAzMDZmOTJiZDQ4Mjc2Nzg4","previousTransactionHash":"ecf50772daf1def3ca02a7ca24925628597aafc68d2bf987857f430832e36cde","originalTransactionHash":"ecf50772daf1def3ca02a7ca24925628597aafc68d2bf987857f430832e36cde","originalSender":"erd1d7y4a8wtykxnxxjhywzk0q5tkey4g9z6rhalefw6syr779kh77yqd0fj5y","sourceShard":1,"destinationShard":0,"miniblockType":"SmartContractResultBlock","miniblockHash":"9f9e6aca41da76e7f3c7b2078c30280b45b4496ad9c2390f2af3517e3eb8d051","logs":{"address":"erd1d7y4a8wtykxnxxjhywzk0q5tkey4g9z6rhalefw6syr779kh77yqd0fj5y","events":[{"address":"erd1qqqqqqqqqqqqqpgq50dge6rrpcra4tp9hl57jl0893a4r2r72jpsk39rjj","identifier":"MultiESDTNFTTransfer","topics":["UVdULTQ2YWMwMQ==","","4YnC","U0ZJVC1hZWJjOTA=","","AeRbu2TRigto","Q1lCRVItNDg5YzFj","","MG+SvUgnZ4g=","b4lencsljTMaVyOFZ4KLtklUFFod+/yl2oEH7xbX94g="],"data":null,"additionalData":["","TXVsdGlFU0RUTkZUVHJhbnNmZXI=","Aw==","UVdULTQ2YWMwMQ==","AA==","4YnC","U0ZJVC1hZWJjOTA=","AA==","AeRbu2TRigto","Q1lCRVItNDg5YzFj","AA==","MG+SvUgnZ4g="]},{"address":"erd1d7y4a8wtykxnxxjhywzk0q5tkey4g9z6rhalefw6syr779kh77yqd0fj5y","identifier":"writeLog","topics":["AAAAAAAAAAAFAKPajOhjDgfarCW/6el95yx7Uah+VIM="],"data":"QDZmNmI=","additionalData":["QDZmNmI="]},{"address":"erd1d7y4a8wtykxnxxjhywzk0q5tkey4g9z6rhalefw6syr779kh77yqd0fj5y","identifier":"completedTxEvent","topics":["7PUHctrx3vPKAqfKJJJWKFl6r8aNK/mHhX9DCDLjbN4="],"data":null,"additionalData":null}]},"status":"success","tokens":["QWT-46ac01","SFIT-aebc90","CYBER-489c1c"],"esdtValues":["14780866","34901695778924399464","3490169577892439944"],"receivers":["erd1d7y4a8wtykxnxxjhywzk0q5tkey4g9z6rhalefw6syr779kh77yqd0fj5y","erd1d7y4a8wtykxnxxjhywzk0q5tkey4g9z6rhalefw6syr779kh77yqd0fj5y","erd1d7y4a8wtykxnxxjhywzk0q5tkey4g9z6rhalefw6syr779kh77yqd0fj5y"],"receiversShardIDs":[0,0,0],"operation":"MultiESDTNFTTransfer","callType":"directCall","options":0},{"type":"unsigned","processingTypeOnSource":"MoveBalance","processingTypeOnDestination":"MoveBalance","hash":"b1531a179dde42a824049e011a31e8a88c3410697ffa50500c168bf67c55dec5","nonce":4170,"round":21008280,"epoch":1458,"value":"329471860000000","receiver":"erd1d7y4a8wtykxnxxjhywzk0q5tkey4g9z6rhalefw6syr779kh77yqd0fj5y","sender":"erd1qqqqqqqqqqqqqpgq50dge6rrpcra4tp9hl57jl0893a4r2r72jpsk39rjj","gasPrice":1000000000,"data":"QDZmNmI=","previousTransactionHash":"ecf50772daf1def3ca02a7ca24925628597aafc68d2bf987857f430832e36cde","originalTransactionHash":"ecf50772daf1def3ca02a7ca24925628597aafc68d2bf987857f430832e36cde","sourceShard":1,"destinationShard":0,"miniblockType":"SmartContractResultBlock","miniblockHash":"9f9e6aca41da76e7f3c7b2078c30280b45b4496ad9c2390f2af3517e3eb8d051","logs":{"address":"erd1d7y4a8wtykxnxxjhywzk0q5tkey4g9z6rhalefw6syr779kh77yqd0fj5y","events":[{"address":"erd1d7y4a8wtykxnxxjhywzk0q5tkey4g9z6rhalefw6syr779kh77yqd0fj5y","identifier":"completedTxEvent","topics":["7PUHctrx3vPKAqfKJJJWKFl6r8aNK/mHhX9DCDLjbN4="],"data":null,"additionalData":null}]},"status":"success","operation":"transfer","isRefund":true,"callType":"directCall","options":0}],"indexOfFirstTxProcessed":0,"indexOfLastTxProcessed":1},{"hash":"13982e2094e459870f6551e1e95f3451d89822e820a4f78dd6a3e06fd938496b","type":"TxBlock","processingType":"Normal","constructionState":"Final","sourceShard":0,"destinationShard":0,"transactions":[{"type":"normal","processingTypeOnSource":"BuiltInFunctionCall","processingTypeOnDestination":"BuiltInFunctionCall","hash":"a9d17c85a5fb9abec85564201ca821f26d90c624c5545b98fb0a808e4ec63c15","nonce":747,"round":21008280,"epoch":1458,"value":"0","receiver":"erd1dfxlnuwp5clncpzwy7m3y2e0kzk6spkyf7fh9kn4gleccvwlvsgqrtr40r","sender":"erd1dfxlnuwp5clncpzwy7m3y2e0kzk6spkyf7fh9kn4gleccvwlvsgqrtr40r","gasPrice":1000000000,"gasLimit":1000000,"data":"RVNEVE5GVFRyYW5zZmVyQDQyNTg0OTRlNTYyZDY2MzMzNjMzNjEzMUAwMUAwMUBiNjNlZTk5YmNiZWJiYTZiMTc1MmMxNTc4MzYwMDgwMDdlNmFlNTMyZTM2MTBmN2U0YzJjZDc1MDYzNzc4MzI0","signature":"4518a4f37424999baf3294b7105fcd89fdea77d39652345799413d0ad62f4b9816d32537b7d9e0d4e235bf091ff775406bd24c08ced0a87da67e7a0e7d15b60f","sourceShard":0,"destinationShard":0,"miniblockType":"TxBlock","miniblockHash":"13982e2094e459870f6551e1e95f3451d89822e820a4f78dd6a3e06fd938496b","logs":{"address":"erd1dfxlnuwp5clncpzwy7m3y2e0kzk6spkyf7fh9kn4gleccvwlvsgqrtr40r","events":[{"address":"erd1dfxlnuwp5clncpzwy7m3y2e0kzk6spkyf7fh9kn4gleccvwlvsgqrtr40r","identifier":"ESDTNFTTransfer","topics":["QlhJTlYtZjM2M2Ex","AQ==","AQ==","tj7pm8vrumsXUsFXg2AIAH5q5TLjYQ9+TCzXUGN3gyQ="],"data":null,"additionalData":["","RVNEVE5GVFRyYW5zZmVy","QlhJTlYtZjM2M2Ex","AQ==","AQ==","tj7pm8vrumsXUsFXg2AIAH5q5TLjYQ9+TCzXUGN3gyQ="]},{"address":"erd1dfxlnuwp5clncpzwy7m3y2e0kzk6spkyf7fh9kn4gleccvwlvsgqrtr40r","identifier":"completedTxEvent","topics":["qdF8haX7mr7IVWQgHKgh8m2QxiTFVFuY+wqAjk7GPBU="],"data":null,"additionalData":null}]},"status":"success","tokens":["BXINV-f363a1-01"],"esdtValues":["1"],"receivers":["erd1kclwnx7tawaxk96jc9tcxcqgqplx4efjudss7ljv9nt4qcmhsvjq5l0wx8"],"receiversShardIDs":[0],"operation":"ESDTNFTTransfer","initiallyPaidFee":"224335000000000","chainID":"1","version":1,"options":0}],"indexOfFirstTxProcessed":0,"indexOfLastTxProcessed":0},{"hash":"21ef17f13edf62c320c6b149dab91c27fb99a52fef9dc04a789a2acf2cc03d45","type":"TxBlock","processingType":"Normal","constructionState":"Final","sourceShard":0,"destinationShard":1,"transactions":[{"type":"normal","processingTypeOnSource":"MoveBalance","processingTypeOnDestination":"SCInvoking","hash":"8e5d4a7b26a010f45e4f46bd23a69da773f09e920e34d7c1c76aa5381e9ca79c","nonce":30335,"round":21008280,"epoch":1458,"value":"2000000000000000000","receiver":"erd1qqqqqqqqqqqqqpgqcc69ts8409p3h77q5chsaqz57y6hugvc4fvs64k74v","sender":"erd1mmh2j5y8esv2tmmyeau3hr4xa2u2te3zc3j9wumn3v5vm8uvsczqnucj5l","gasPrice":1000000000,"gasLimit":70050000,"data":"YWdncmVnYXRlRWdsZEAwMDAwMDAwYzU3NDU0NzRjNDQyZDYyNjQzNDY0MzczOTAwMDAwMDBhNDI0NjU5MmQzODMzMzQzNDY2NjYwMDAwMDAwODE3ZTFiOWQ0YmEwOGMxMjkwMDAwMDAwMDAwMDAwMDAwMDUwMDJlZGYwOGU3ZmQyODc3NGZhZjI3MThjNGUwODk2Mjg5NGFmMTA5NTE1NDgzMDAwMDAwMTQ3Mzc3NjE3MDU0NmY2YjY1NmU3MzQ2Njk3ODY1NjQ0OTZlNzA3NTc0MDAwMDAwMDIwMDAwMDAwYTQyNDY1OTJkMzgzMzM0MzQ2NjY2MDAwMDAwMDEwMTAwMDAwMDBjNTc0NTQ3NGM0NDJkNjI2NDM0NjQzNzM5MDAwMDAwMGE0MjQ2NTkyZDM4MzMzNDM0NjY2NjAwMDAwMDA4MDNkZmIzOTI5NGJmM2VkNzAwMDAwMDAwMDAwMDAwMDAwNTAwMDBiNGMwOTQ5NDdlNDI3ZDc5OTMxYThiYWQ4MTMxNmI3OTdkMjM4Y2RiM2YwMDAwMDAxOTczNzc2MTcwNGQ3NTZjNzQ2OTU0NmY2YjY1NmU3MzQ2Njk3ODY1NjQ0OTZlNzA3NTc0MDAwMDAwMDQwMDAwMDAwMTAxMDAwMDAwMDAwMDAwMDAwYzU3NDU0NzRjNDQyZDYyNjQzNDY0MzczOTAwMDAwMDBhNDI0NjU5MmQzODMzMzQzNDY2NjZAMDAwMDAwMGM1NzQ1NDc0YzQ0MmQ2MjY0MzQ2NDM3MzkwMDAwMDAwMDAwMDAwMDBhNDI0NjU5MmQzODMzMzQzNDY2NjYwMDAwMDAwOTE0Y2I0NWM3YzVmMjFhZWJiZQ==","signature":"e305df6b6b6f56391e3416e325abfa32a295402671be3c28a20d89e956c8b9c3fb38b5d7e05bfd20d35b933233fa7d2a6c209d17c410e83b691ec7116b41ef04","sourceShard":0,"destinationShard":1,"miniblockType":"TxBlock","miniblockHash":"21ef17f13edf62c320c6b149dab91c27fb99a52fef9dc04a789a2acf2cc03d45","status":"pending","operation":"transfer","function":"aggregateEgld","initiallyPaidFee":"1754355000000000","chainID":"1","version":2,"options":2,"guardian":"erd1y4nexjs97zjtnzqqtxjgkn2dsjgt85dc9kesx87vx6ker3nfac5sr6333t","guardianSignature":"bf7bd02cd997a39a8798f1bc51fd03d1f3068fd052b3c354a9652dbf98532a6be83e07d4db3c0c2f4f40ccb6947e74a2e617876c5a82830c04064e8745bc5a0f"}],"indexOfFirstTxProcessed":0,"indexOfLastTxProcessed":0},{"hash":"882e7558709a8813578583b2470cb19ff0540a1059a79b37751d6492fe82decd","type":"SmartContractResultBlock","processingType":"Normal","isFromReceiptsStorage":true,"sourceShard":0,"destinationShard":0,"transactions":[{"type":"unsigned","processingTypeOnSource":"MoveBalance","processingTypeOnDestination":"MoveBalance","hash":"0f27d677690fd9b71eb3e04771687d0f7635c6acd9b75bbfc8f2e45bfce7ad32","nonce":748,"round":21008280,"epoch":1458,"value":"5697500000000","receiver":"erd1dfxlnuwp5clncpzwy7m3y2e0kzk6spkyf7fh9kn4gleccvwlvsgqrtr40r","sender":"erd1dfxlnuwp5clncpzwy7m3y2e0kzk6spkyf7fh9kn4gleccvwlvsgqrtr40r","gasPrice":1000000000,"data":"QDZmNmI=","previousTransactionHash":"a9d17c85a5fb9abec85564201ca821f26d90c624c5545b98fb0a808e4ec63c15","originalTransactionHash":"a9d17c85a5fb9abec85564201ca821f26d90c624c5545b98fb0a808e4ec63c15","sourceShard":0,"destinationShard":0,"miniblockType":"SmartContractResultBlock","miniblockHash":"882e7558709a8813578583b2470cb19ff0540a1059a79b37751d6492fe82decd","status":"success","operation":"transfer","isRefund":true,"callType":"directCall","options":0}],"indexOfFirstTxProcessed":0,"indexOfLastTxProcessed":0}],"scheduledData":{"rootHash":"8b42e3aa3975148a5ac6757b37182e9034f91a68e26a106138ce90844c8b7f79","accumulatedFees":"0","developerFees":"0"}}},"error":"","code":"successful"}`
	httpNodeStatus := `{"data":{"metrics":{"erd_current_round":187068,"erd_epoch_number":12,"erd_highest_final_nonce":21980319,"erd_nonce":21980327,"erd_nonce_at_epoch_start":172770,"erd_nonces_passed_in_current_epoch":14253,"erd_round_at_epoch_start":172814,"erd_rounds_passed_in_current_epoch":14254,"erd_rounds_per_epoch":14400}},"error":"","code":"successful"}`

	t.Run("invalid block range", func(t *testing.T) {
		invalidFilter := &sdkCore.FilterQuery{
			FromBlock: core.OptionalUint64{21000000, true},
			ToBlock:   core.OptionalUint64{21001000, true},
			ShardID:   0,
			Topics:    nil,
			BlockHash: nil,
		}

		statusResponseBytes := []byte(httpNodeStatus)

		responseMap := map[string][]byte{
			"https://test.org/node/status": statusResponseBytes,
		}

		httpClient := createMockClientMultiResponse(responseMap)
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		res, err := ep.FilterLogs(context.Background(), invalidFilter)
		assert.Equal(t, ErrInvalidBlockRange, err)
		assert.Nil(t, res)
	})

	t.Run("toBlock greater than last block", func(t *testing.T) {
		invalidFilter := &sdkCore.FilterQuery{
			FromBlock: core.OptionalUint64{21000000, false},
			ToBlock:   core.OptionalUint64{1000000000, false},
			ShardID:   0,
			Topics:    nil,
			BlockHash: nil,
		}

		statusResponseBytes := []byte(httpNodeStatus)

		responseMap := map[string][]byte{
			"https://test.org/node/status": statusResponseBytes,
		}

		httpClient := createMockClientMultiResponse(responseMap)
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		res, err := ep.FilterLogs(context.Background(), invalidFilter)
		assert.Equal(t, ErrNoBlockRangeProvided, err)
		assert.Nil(t, res)
	})

	t.Run("no block range provided", func(t *testing.T) {
		invalidFilter := &sdkCore.FilterQuery{
			FromBlock: core.OptionalUint64{0, false},
			ToBlock:   core.OptionalUint64{0, false},
			ShardID:   0,
			Topics:    nil,
			BlockHash: nil,
		}

		statusResponseBytes := []byte(httpNodeStatus)

		responseMap := map[string][]byte{
			"https://test.org/node/status": statusResponseBytes,
		}

		httpClient := createMockClientMultiResponse(responseMap)
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		res, err := ep.FilterLogs(context.Background(), invalidFilter)
		assert.Equal(t, ErrNoBlockRangeProvided, err)
		assert.Nil(t, res)
	})

	t.Run("restricts matches for specific addresses", func(t *testing.T) {

		validFilter := &sdkCore.FilterQuery{
			FromBlock: core.OptionalUint64{21000005, true},
			ToBlock:   core.OptionalUint64{21000005, true},
			Addresses: []string{"erd1qqqqqqqqqqqqqpgq50dge6rrpcra4tp9hl57jl0893a4r2r72jpsk39rjj", "erd1d7y4a8wtykxnxxjhywzk0q5tkey4g9z6rhalefw6syr779kh77yqd0fj5y"},
			ShardID:   0,
			Topics:    nil,
			BlockHash: nil,
		}

		statusResponseBytes := []byte(httpNodeStatus)
		blockResponseBytes := []byte(httpDataBlock21000005)

		responseMap := map[string][]byte{
			"https://test.org/node/status":             statusResponseBytes,
			"https://test.org/block/by-nonce/21000005": blockResponseBytes,
		}

		httpClient := createMockClientMultiResponse(responseMap)
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		res, err := ep.FilterLogs(context.Background(), validFilter)

		assert.Nil(t, err)
		assert.Equal(t, len(res), 4)
		assert.Equal(t, res[0].Identifier, "MultiESDTNFTTransfer")
		assert.Equal(t, res[1].Identifier, "writeLog")
		assert.Equal(t, res[2].Identifier, "completedTxEvent")
		assert.Equal(t, res[3].Identifier, "completedTxEvent")
	})

	t.Run("with topics", func(t *testing.T) {
		validFilter := &sdkCore.FilterQuery{
			BlockHash: nil,
			FromBlock: core.OptionalUint64{21000000, true},
			ToBlock:   core.OptionalUint64{21000000, true},
			ShardID:   0,
			Topics: [][]byte{
				{72, 85, 84, 75, 45, 52, 102, 97, 52, 98, 50}, // HUTK-4fa4b2
			},
		}

		statusResponseBytes := []byte(httpNodeStatus)
		blockResponseBytes := []byte(httpDataBlock21000000)

		responseMap := map[string][]byte{
			"https://test.org/node/status":             statusResponseBytes,
			"https://test.org/block/by-nonce/21000000": blockResponseBytes,
			"https://test.org/block/by-hash/00f7d0e806c5ff3700236f78c67007c8000000000000000000409a07ff835d8e": blockResponseBytes,
		}

		httpClient := createMockClientMultiResponse(responseMap)
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		res, err := ep.FilterLogs(context.Background(), validFilter)

		assert.Nil(t, err)
		assert.Equal(t, len(res), 1)
		assert.Equal(t, res[0].Identifier, "ESDTTransfer")
	})

	t.Run("should work without caching", func(t *testing.T) {
		validFilter := &sdkCore.FilterQuery{
			BlockHash: &[32]byte{
				0x00, 0xf7, 0xd0, 0xe8, 0x06, 0xc5, 0xff, 0x37,
				0x00, 0x23, 0x6f, 0x78, 0xc6, 0x70, 0x07, 0xc8,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x40, 0x9a, 0x07, 0xff, 0x83, 0x5d, 0x8e,
			},
			FromBlock: core.OptionalUint64{21000000, false},
			ToBlock:   core.OptionalUint64{21000000, false},
			ShardID:   0,
			Topics:    nil,
		}

		blockResponseBytes := []byte(httpDataBlock21000000)
		statusResponseBytes := []byte(httpNodeStatus)

		responseMap := map[string][]byte{
			"https://test.org/node/status":             statusResponseBytes,
			"https://test.org/block/by-nonce/21000000": blockResponseBytes,
			"https://test.org/block/by-hash/00f7d0e806c5ff3700236f78c67007c8000000000000000000409a07ff835d8e": blockResponseBytes,
		}

		httpClient := createMockClientMultiResponse(responseMap)
		args := createMockArgsProxy(httpClient)
		ep, _ := NewProxy(args)

		res, err := ep.FilterLogs(context.Background(), validFilter)

		assert.Nil(t, err)
		assert.Equal(t, len(res), 4)
		assert.Equal(t, res[0].Identifier, "ESDTTransfer")
		assert.Equal(t, res[1].Identifier, "writeLog")
		assert.Equal(t, res[2].Identifier, "completedTxEvent")
		assert.Equal(t, res[3].Identifier, "completedTxEvent")
	})

	t.Run("should work with caching", func(t *testing.T) {

		validFilter := &sdkCore.FilterQuery{
			FromBlock: core.OptionalUint64{21000001, true},
			ToBlock:   core.OptionalUint64{21000001, true},
			ShardID:   0,
			Topics:    nil,
			BlockHash: nil,
		}

		blockResponseBytes := []byte(httpDataBlock21000001)
		statusResponseBytes := []byte(httpNodeStatus)

		responseMap := map[string][]byte{
			"https://test.org/node/status":             statusResponseBytes,
			"https://test.org/block/by-nonce/21000001": blockResponseBytes,
			"https://test.org/block/by-hash/00f7d0e806c5ff3700236f78c67007c8000000000000000000409a07ff835d8e": blockResponseBytes,
		}

		httpClient := createMockClientMultiResponse(responseMap)
		args := createMockArgsProxyWithCache(httpClient)
		ep, _ := NewProxy(args)

		_, err := ep.FilterLogs(context.Background(), validFilter)
		assert.Nil(t, err)

		// check if the block is cached
		cacheKey := make([]byte, 8)
		binary.BigEndian.PutUint64(cacheKey, 21000001)
		value, _ := ep.filterQueryBlockCacher.Get(cacheKey)
		cachedBlockBytes, _ := value.([]byte)
		assert.Equal(t, cachedBlockBytes, blockResponseBytes)

		res2, err := ep.FilterLogs(context.Background(), validFilter)

		assert.Nil(t, err)
		assert.Equal(t, len(res2), 7)
		assert.Equal(t, res2[0].Identifier, "ESDTTransfer")
		assert.Equal(t, res2[1].Identifier, "ESDTLocalBurn")
		assert.Equal(t, res2[2].Identifier, "ESDTTransfer")
		assert.Equal(t, res2[3].Identifier, "ESDTNFTTransfer")
		assert.Equal(t, res2[4].Identifier, "buyOffer")
		assert.Equal(t, res2[5].Identifier, "writeLog")
		assert.Equal(t, res2[6].Identifier, "completedTxEvent")
	})
}
