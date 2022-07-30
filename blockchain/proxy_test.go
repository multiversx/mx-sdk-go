package blockchain

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	erdgoCore "github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	erdgoHttp "github.com/ElrondNetwork/elrond-sdk-erdgo/core/http"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/testsCommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testHttpURL = "https://test.org"
const networkConfigEndpoint = "network/config"
const getNetworkStatusEndpoint = "network/status/%d"
const getNodeStatusEndpoint = "node/status"

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
	return createMockClientRespondingBytesWithStatusCode(responseBytes, http.StatusOK)
}

func createMockClientRespondingBytesWithStatusCode(responseBytes []byte, status int) *mockHTTPClient {
	return &mockHTTPClient{
		doCalled: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       ioutil.NopCloser(bytes.NewReader(responseBytes)),
				StatusCode: status,
			}, nil
		},
	}
}

func createMockArgsElrondProxy(httpClient erdgoHttp.Client) ArgsElrondProxy {
	return ArgsElrondProxy{
		ProxyURL:            testHttpURL,
		Client:              httpClient,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		AllowedDeltaToFinal: 1,
		CacheExpirationTime: time.Second,
		EntityType:          erdgoCore.ObserverNode,
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
		Body:       ioutil.NopCloser(bytes.NewReader(buff)),
		StatusCode: http.StatusOK,
	}, handled, nil
}

func TestNewElrondProxy(t *testing.T) {
	t.Parallel()

	t.Run("invalid time cache should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsElrondProxy(nil)
		args.CacheExpirationTime = time.Second - time.Nanosecond
		proxy, err := NewElrondProxy(args)

		assert.True(t, check.IfNil(proxy))
		assert.True(t, errors.Is(err, ErrInvalidCacherDuration))
	})
	t.Run("invalid nonce delta should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsElrondProxy(nil)
		args.FinalityCheck = true
		args.AllowedDeltaToFinal = 0
		proxy, err := NewElrondProxy(args)

		assert.True(t, check.IfNil(proxy))
		assert.True(t, errors.Is(err, ErrInvalidAllowedDeltaToFinal))
	})
	t.Run("should work with finality check", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsElrondProxy(nil)
		args.FinalityCheck = true
		proxy, err := NewElrondProxy(args)

		assert.False(t, check.IfNil(proxy))
		assert.Nil(t, err)
	})
	t.Run("should work without finality check", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsElrondProxy(nil)
		proxy, err := NewElrondProxy(args)

		assert.False(t, check.IfNil(proxy))
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
				Body:       ioutil.NopCloser(bytes.NewReader(accountBytes)),
				StatusCode: http.StatusOK,
			}, nil
		},
	}
	args := createMockArgsElrondProxy(httpClient)
	args.FinalityCheck = true
	proxy, _ := NewElrondProxy(args)

	address, err := data.NewAddressFromBech32String("erd1qqqqqqqqqqqqqpgqfzydqmdw7m2vazsp6u5p95yxz76t2p9rd8ss0zp9ts")
	if err != nil {
		assert.Error(t, err)
	}
	expectedErr := errors.New("expected error")

	t.Run("finality checker errors should not query", func(t *testing.T) {
		proxy.finalityProvider = &testsCommon.FinalityProviderStub{
			CheckShardFinalizationCalled: func(ctx context.Context, targetShardID uint32, maxNoncesDelta uint64) error {
				return expectedErr
			},
		}

		account, errGet := proxy.GetAccount(context.Background(), address)
		assert.Nil(t, account)
		assert.True(t, errors.Is(errGet, expectedErr))
		assert.Equal(t, uint32(0), atomic.LoadUint32(&numAccountQueries))
	})
	t.Run("finality checker returns nil should return account", func(t *testing.T) {
		finalityCheckCalled := uint32(0)
		proxy.finalityProvider = &testsCommon.FinalityProviderStub{
			CheckShardFinalizationCalled: func(ctx context.Context, targetShardID uint32, maxNoncesDelta uint64) error {
				atomic.AddUint32(&finalityCheckCalled, 1)
				return nil
			},
		}

		account, errGet := proxy.GetAccount(context.Background(), address)
		assert.NotNil(t, account)
		assert.Equal(t, uint64(37), account.Nonce)
		assert.Nil(t, errGet)
		assert.Equal(t, uint32(1), atomic.LoadUint32(&numAccountQueries))
		assert.Equal(t, uint32(1), atomic.LoadUint32(&finalityCheckCalled))
	})
}

func TestElrondProxy_GetNetworkEconomics(t *testing.T) {
	t.Parallel()

	responseBytes := []byte(`{"data":{"metrics":{"erd_dev_rewards":"0","erd_epoch_for_economics_data":263,"erd_inflation":"5869888769785838708144","erd_total_fees":"51189055176110000000","erd_total_staked_value":"9963775651405816710680128","erd_total_supply":"21556417261819025351089574","erd_total_top_up_value":"1146275808171377418645274"}},"code":"successful"}`)
	httpClient := createMockClientRespondingBytes(responseBytes)
	args := createMockArgsElrondProxy(httpClient)
	ep, _ := NewElrondProxy(args)

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

func TestElrondProxy_RequestTransactionCost(t *testing.T) {
	t.Parallel()

	responseBytes := []byte(`{"data":{"txGasUnits":24273810,"returnMessage":""},"error":"","code":"successful"}`)
	httpClient := createMockClientRespondingBytes(responseBytes)
	args := createMockArgsElrondProxy(httpClient)
	ep, _ := NewElrondProxy(args)

	tx := &data.Transaction{
		Nonce:   1,
		Value:   "50",
		RcvAddr: "erd1rh5ws22jxm9pe7dtvhfy6j3uttuupkepferdwtmslms5fydtrh5sx3xr8r",
		SndAddr: "erd1rh5ws22jxm9pe7dtvhfy6j3uttuupkepferdwtmslms5fydtrh5sx3xr8r",
		Data:    []byte("hello"),
		ChainID: "1",
		Version: 1,
		Options: 0,
	}
	txCost, err := ep.RequestTransactionCost(context.Background(), tx)
	require.Nil(t, err)
	require.Equal(t, &data.TxCostResponseData{
		TxCost:     24273810,
		RetMessage: "",
	}, txCost)
}

func TestElrondProxy_GetTransactionInfoWithResults(t *testing.T) {
	t.Parallel()

	responseBytes := []byte(`{"data":{"transaction":{"type":"normal","nonce":22100,"round":53057,"epoch":263,"value":"0","receiver":"erd1e6c9vcga5lyhwdu9nr9lya4ujz2r5w2egsfjp0lslrgv5dsccpdsmre6va","sender":"erd1e6c9vcga5lyhwdu9nr9lya4ujz2r5w2egsfjp0lslrgv5dsccpdsmre6va","gasPrice":1000000001,"gasLimit":19210000,"data":"RVNEVE5GVFRyYW5zZmVyQDU5NTk1OTQ1MzY0MzM5MkQzMDM5MzQzOTMwMzBAMDFAMDFAMDAwMDAwMDAwMDAwMDAwMDA1MDA1QzgzRTBDNDJFRENFMzk0RjQwQjI0RDI5RDI5OEIwMjQ5QzQxRjAyODk3NEA2Njc1NkU2NEA1N0M1MDZBMTlEOEVBRTE4Mjk0MDNBOEYzRjU0RTFGMDM3OUYzODE1N0ZDRjUzRDVCQ0E2RjIzN0U0QTRDRjYxQDFjMjA=","signature":"37922ccb13d46857d819cd618f1fd8e76777e14f1c16a54eeedd55c67d5883db0e3a91cd0ff0f541e2c3d7347488b4020c0ba4765c9bc02ea0ae1c3f6db1ec05","sourceShard":1,"destinationShard":1,"blockNonce":53052,"blockHash":"4a63312d1bfe48aa516185d12abff5daf6343fce1f298db51291168cf97a790c","notarizedAtSourceInMetaNonce":53053,"NotarizedAtSourceInMetaHash":"342d189e36ef5cbf9f8b3f2cb5bf2cb8e2260062c6acdf89ce5faacd99f4dbcc","notarizedAtDestinationInMetaNonce":53053,"notarizedAtDestinationInMetaHash":"342d189e36ef5cbf9f8b3f2cb5bf2cb8e2260062c6acdf89ce5faacd99f4dbcc","miniblockType":"TxBlock","miniblockHash":"0c659cce5e2653522cc0e3cf35571264522035a7aef4ffa5244d1ed3d8bc01a8","status":"success","hyperblockNonce":53053,"hyperblockHash":"342d189e36ef5cbf9f8b3f2cb5bf2cb8e2260062c6acdf89ce5faacd99f4dbcc","smartContractResults":[{"hash":"5ab14959aaeb3a20d95ec6bbefc03f251732d9368711c55c63d3811e70903f4e","nonce":0,"value":0,"receiver":"erd1qqqqqqqqqqqqqpgqtjp7p3pwmn3efaqtynff62vtqfyug8cz396qs5vnsy","sender":"erd1e6c9vcga5lyhwdu9nr9lya4ujz2r5w2egsfjp0lslrgv5dsccpdsmre6va","data":"ESDTNFTTransfer@595959453643392d303934393030@01@01@080112020001226f08011204746573741a20ceb056611da7c977378598cbf276bc90943a3959441320bff0f8d0ca3618c05b20e8072a206161616161616161616161616161616161616161616161616161616161616161320461626261321268747470733a2f2f656c726f6e642e636f6d3a0474657374@66756e64@57c506a19d8eae1829403a8f3f54e1f0379f38157fcf53d5bca6f237e4a4cf61@1c20","prevTxHash":"c7fadeaccce0673bd6ce3a1f472f7cc1beef20c0b3131cfa9866cd5075816639","originalTxHash":"c7fadeaccce0673bd6ce3a1f472f7cc1beef20c0b3131cfa9866cd5075816639","gasLimit":18250000,"gasPrice":1000000001,"callType":0},{"hash":"f337d2705d2b644f9d8f75cf270e879b6ada51c4c54009e94305adf368d1adbc","nonce":1,"value":102793170000000,"receiver":"erd1e6c9vcga5lyhwdu9nr9lya4ujz2r5w2egsfjp0lslrgv5dsccpdsmre6va","sender":"erd1qqqqqqqqqqqqqpgqtjp7p3pwmn3efaqtynff62vtqfyug8cz396qs5vnsy","data":"@6f6b","prevTxHash":"5ab14959aaeb3a20d95ec6bbefc03f251732d9368711c55c63d3811e70903f4e","originalTxHash":"c7fadeaccce0673bd6ce3a1f472f7cc1beef20c0b3131cfa9866cd5075816639","gasLimit":0,"gasPrice":1000000001,"callType":0}]}},"error":"","code":"successful"}`)
	httpClient := createMockClientRespondingBytes(responseBytes)
	args := createMockArgsElrondProxy(httpClient)
	ep, _ := NewElrondProxy(args)

	tx, err := ep.GetTransactionInfoWithResults(context.Background(), "a40e5a6af4efe221608297a73459211756ab88b96896e6e331842807a138f343")
	require.Nil(t, err)

	txBytes, _ := json.MarshalIndent(tx, "", " ")
	fmt.Println(string(txBytes))
}

func TestElrondProxy_ExecuteVmQuery(t *testing.T) {
	t.Parallel()

	responseBytes := []byte(`{"data":{"data":{"returnData":["MC41LjU="],"returnCode":"ok","returnMessage":"","gasRemaining":18446744073685949187,"gasRefund":0,"outputAccounts":{"0000000000000000050033bb65a91ee17ab84c6f8a01846ef8644e15fb76696a":{"address":"erd1qqqqqqqqqqqqqpgqxwakt2g7u9atsnr03gqcgmhcv38pt7mkd94q6shuwt","nonce":0,"balance":null,"balanceDelta":0,"storageUpdates":{},"code":null,"codeMetaData":null,"outputTransfers":[],"callType":0}},"deletedAccounts":[],"touchedAccounts":[],"logs":[]}},"error":"","code":"successful"}`)
	t.Run("no finality check", func(t *testing.T) {
		httpClient := createMockClientRespondingBytes(responseBytes)
		args := createMockArgsElrondProxy(httpClient)
		ep, _ := NewElrondProxy(args)

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
		args := createMockArgsElrondProxy(httpClient)
		args.FinalityCheck = true
		ep, _ := NewElrondProxy(args)

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
		args := createMockArgsElrondProxy(httpClient)
		args.FinalityCheck = true
		ep, _ := NewElrondProxy(args)

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
					Body:       ioutil.NopCloser(bytes.NewReader(responseBytes)),
					StatusCode: http.StatusOK,
				}, nil
			},
		}
		args := createMockArgsElrondProxy(httpClient)
		args.FinalityCheck = true
		ep, _ := NewElrondProxy(args)

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

func TestElrondProxy_GetRawBlockByHash(t *testing.T) {
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
	args := createMockArgsElrondProxy(httpClient)
	ep, _ := NewElrondProxy(args)

	response, err := ep.GetRawBlockByHash(context.Background(), 0, "aaaa")
	require.Nil(t, err)

	ts := &testStruct{}
	err = json.Unmarshal(response, ts)
	require.Nil(t, err)
	require.Equal(t, expectedTs, ts)
}

func TestElrondProxy_GetRawBlockByNonce(t *testing.T) {
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
	args := createMockArgsElrondProxy(httpClient)
	ep, _ := NewElrondProxy(args)

	response, err := ep.GetRawBlockByNonce(context.Background(), 0, 10)
	require.Nil(t, err)

	ts := &testStruct{}
	err = json.Unmarshal(response, ts)
	require.Nil(t, err)
	require.Equal(t, expectedTs, ts)
}

func TestElrondProxy_GetRawMiniBlockByHash(t *testing.T) {
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
	args := createMockArgsElrondProxy(httpClient)
	ep, _ := NewElrondProxy(args)

	response, err := ep.GetRawMiniBlockByHash(context.Background(), 0, "aaaa", 1)
	require.Nil(t, err)

	ts := &testStruct{}
	err = json.Unmarshal(response, ts)
	require.Nil(t, err)
	require.Equal(t, expectedTs, ts)
}

func TestElrondProxy_GetNonceAtEpochStart(t *testing.T) {
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
	args := createMockArgsElrondProxy(httpClient)
	ep, _ := NewElrondProxy(args)

	response, err := ep.GetNonceAtEpochStart(context.Background(), core.MetachainShardId)
	require.Nil(t, err)

	require.Equal(t, expectedNonce, response)
}

func TestElrondProxy_GetRatingsConfig(t *testing.T) {
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
	args := createMockArgsElrondProxy(httpClient)
	ep, _ := NewElrondProxy(args)

	response, err := ep.GetRatingsConfig(context.Background())
	require.Nil(t, err)

	require.Equal(t, expectedRatingsConfig, response)
}

func TestElrondProxy_GetEnableEpochsConfig(t *testing.T) {
	t.Parallel()

	expectedEnableEpochsConfig := &data.EnableEpochsConfig{
		BalanceWaitingListsEnableEpoch: 1,
		WaitingListFixEnableEpoch:      1,
	}
	enableEpochsResponse := &data.EnableEpochsConfigResponse{
		Data: struct {
			Config *data.EnableEpochsConfig "json:\"enableEpochs\""
		}{
			Config: expectedEnableEpochsConfig},
	}
	enableEpochsResponseBytes, _ := json.Marshal(enableEpochsResponse)

	httpClient := createMockClientRespondingBytes(enableEpochsResponseBytes)
	args := createMockArgsElrondProxy(httpClient)
	ep, _ := NewElrondProxy(args)

	response, err := ep.GetEnableEpochsConfig(context.Background())
	require.Nil(t, err)

	require.Equal(t, expectedEnableEpochsConfig, response)
}

func TestElrondProxy_GetGenesisNodesPubKeys(t *testing.T) {
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
	args := createMockArgsElrondProxy(httpClient)
	ep, _ := NewElrondProxy(args)

	response, err := ep.GetGenesisNodesPubKeys(context.Background())
	require.Nil(t, err)

	require.Equal(t, expectedGenesisNodes, response)
}

func TestElrondProxy_GetAccountKeys(t *testing.T) {
	t.Parallel()

	defaultResponseBytes := []byte(`{"data":{"pairs":{"666f6f":"626172"}},"error":"","code":""}`) // "key": "value"
	address, err := data.NewAddressFromBech32String("erd1qqqqqqqqqqqqqpgqfzydqmdw7m2vazsp6u5p95yxz76t2p9rd8ss0zp9ts")
	assert.NotNil(t, address)
	assert.Nil(t, err)

	t.Run("retrieve all account storage data", func(t *testing.T) {
		httpClient := createMockClientRespondingBytes(defaultResponseBytes)
		args := createMockArgsElrondProxy(httpClient)
		proxy, _ := NewElrondProxy(args)

		expectedAccountKeys := &data.AccountKeys{
			"666f6f": "626172", // "key": "value"
		}
		response, err := proxy.GetAccountKeys(context.Background(), address)
		assert.Nil(t, err)
		assert.Equal(t, expectedAccountKeys, response)
	})
	t.Run("fail on empty address", func(t *testing.T) {
		httpClient := createMockClientRespondingBytes(defaultResponseBytes)
		args := createMockArgsElrondProxy(httpClient)
		proxy, _ := NewElrondProxy(args)

		response, err := proxy.GetAccountKeys(context.Background(), nil)
		assert.Error(t, err)
		assert.Nil(t, response)
	})
	t.Run("fail on invalid address", func(t *testing.T) {
		httpClient := createMockClientRespondingBytes(defaultResponseBytes)
		args := createMockArgsElrondProxy(httpClient)
		proxy, _ := NewElrondProxy(args)

		response, err := proxy.GetAccountKeys(context.Background(), data.NewAddressFromBytes([]byte{}))
		assert.Error(t, err)
		assert.Nil(t, response)
	})
	t.Run("fail on empty response", func(t *testing.T) {
		httpClient := createMockClientRespondingBytes([]byte{})
		args := createMockArgsElrondProxy(httpClient)
		proxy, _ := NewElrondProxy(args)

		response, err := proxy.GetAccountKeys(context.Background(), address)
		assert.NotNil(t, err)
		assert.IsType(t, &json.SyntaxError{}, err)
		assert.Nil(t, response)
	})
	t.Run("fail on invalid response", func(t *testing.T) {
		httpClient := createMockClientRespondingBytes([]byte("invalid"))
		args := createMockArgsElrondProxy(httpClient)
		proxy, _ := NewElrondProxy(args)

		response, err := proxy.GetAccountKeys(context.Background(), address)
		assert.NotNil(t, err)
		assert.IsType(t, &json.SyntaxError{}, err)
		assert.Nil(t, response)
	})
	t.Run("fail on response error status", func(t *testing.T) {
		httpClient := createMockClientRespondingBytesWithStatusCode(defaultResponseBytes, http.StatusBadRequest)
		args := createMockArgsElrondProxy(httpClient)
		proxy, _ := NewElrondProxy(args)

		response, err := proxy.GetAccountKeys(context.Background(), address)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Bad Request")
		assert.Nil(t, response)
	})
	t.Run("fail on response error message", func(t *testing.T) {
		httpClient := createMockClientRespondingBytes([]byte(`{"data":{"pairs":{"666f6f":"626172"}},"error":"fail","code":""}`))
		args := createMockArgsElrondProxy(httpClient)
		proxy, _ := NewElrondProxy(args)

		response, err := proxy.GetAccountKeys(context.Background(), address)
		assert.NotNil(t, err)
		assert.IsType(t, errors.New(""), err)
		assert.Nil(t, response)
	})
}
