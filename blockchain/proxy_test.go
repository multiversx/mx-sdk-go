package blockchain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testHttpURL = "https://test.org"

type testStruct struct {
	Nonce int
	Name  string
}

type mockHTTPClient struct {
	lastRequest *http.Request
	doCalled    func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.doCalled != nil {
		return m.doCalled(req)
	}

	m.lastRequest = req
	return &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte("account"))),
	}, nil
}

func TestElrondProxy_GetHTTPContextDone(t *testing.T) {
	t.Parallel()

	testHttpServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// simulating that the operation takes a lot of time

		time.Sleep(time.Second * 2)

		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write(nil)
	}))
	proxy := NewElrondProxy(testHttpServer.URL, &http.Client{})

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	resp, err := proxy.GetHTTP(ctx, "endpoint")
	assert.Nil(t, resp)
	require.NotNil(t, err)
	assert.Equal(t, "*url.Error", fmt.Sprintf("%T", err))
	assert.True(t, strings.Contains(err.Error(), "context deadline exceeded"))
}

func TestElrondProxy_PostHTTPContextDone(t *testing.T) {
	t.Parallel()

	testHttpServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// simulating that the operation takes a lot of time

		time.Sleep(time.Second * 2)

		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write(nil)
	}))
	proxy := NewElrondProxy(testHttpServer.URL, &http.Client{})

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	resp, err := proxy.PostHTTP(ctx, "endpoint", nil)
	assert.Nil(t, resp)
	require.NotNil(t, err)
	assert.Equal(t, "*url.Error", fmt.Sprintf("%T", err))
	assert.True(t, strings.Contains(err.Error(), "context deadline exceeded"))
}

func TestGetAccount(t *testing.T) {
	t.Parallel()

	httpClient := &mockHTTPClient{}
	proxy := NewElrondProxy(testHttpURL, httpClient)

	address, err := data.NewAddressFromBech32String("erd1qqqqqqqqqqqqqpgqfzydqmdw7m2vazsp6u5p95yxz76t2p9rd8ss0zp9ts")
	if err != nil {
		assert.Error(t, err)
	}

	_, err = proxy.GetAccount(context.Background(), address)
	if err != nil {
		assert.Error(t, err)
	}

	expected := testHttpURL + "/address/erd1qqqqqqqqqqqqqpgqfzydqmdw7m2vazsp6u5p95yxz76t2p9rd8ss0zp9ts"

	assert.Equal(t, expected, httpClient.lastRequest.URL.String())
}

func TestElrondProxy_GetNetworkEconomics(t *testing.T) {
	t.Parallel()

	responseBytes := []byte(`{"data":{"metrics":{"erd_dev_rewards":"0","erd_epoch_for_economics_data":263,"erd_inflation":"5869888769785838708144","erd_total_fees":"51189055176110000000","erd_total_staked_value":"9963775651405816710680128","erd_total_supply":"21556417261819025351089574","erd_total_top_up_value":"1146275808171377418645274"}},"code":"successful"}`)
	httpClient := &mockHTTPClient{
		doCalled: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader(responseBytes)),
			}, nil
		},
	}
	ep := NewElrondProxy("http://localhost:8079", httpClient)

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
	httpClient := &mockHTTPClient{
		doCalled: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader(responseBytes)),
			}, nil
		},
	}
	ep := NewElrondProxy("http://localhost:8079", httpClient)

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
	httpClient := &mockHTTPClient{
		doCalled: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader(responseBytes)),
			}, nil
		},
	}
	ep := NewElrondProxy("http://localhost:8079", httpClient)

	tx, err := ep.GetTransactionInfoWithResults(context.Background(), "a40e5a6af4efe221608297a73459211756ab88b96896e6e331842807a138f343")
	require.Nil(t, err)

	txBytes, _ := json.MarshalIndent(tx, "", " ")
	fmt.Println(string(txBytes))
}

func TestElrondProxy_ExecuteVmQuery(t *testing.T) {
	responseBytes := []byte(`{"data":{"data":{"returnData":["MC41LjU="],"returnCode":"ok","returnMessage":"","gasRemaining":18446744073685949187,"gasRefund":0,"outputAccounts":{"0000000000000000050033bb65a91ee17ab84c6f8a01846ef8644e15fb76696a":{"address":"erd1qqqqqqqqqqqqqpgqxwakt2g7u9atsnr03gqcgmhcv38pt7mkd94q6shuwt","nonce":0,"balance":null,"balanceDelta":0,"storageUpdates":{},"code":null,"codeMetaData":null,"outputTransfers":[],"callType":0}},"deletedAccounts":[],"touchedAccounts":[],"logs":[]}},"error":"","code":"successful"}`)
	httpClient := &mockHTTPClient{
		doCalled: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader(responseBytes)),
			}, nil
		},
	}
	_ = httpClient
	ep := NewElrondProxy("http://localhost:8079", httpClient)

	response, err := ep.ExecuteVMQuery(context.Background(), &data.VmValueRequest{
		Address:    "erd1qqqqqqqqqqqqqpgqxwakt2g7u9atsnr03gqcgmhcv38pt7mkd94q6shuwt",
		FuncName:   "version",
		CallerAddr: "erd1qqqqqqqqqqqqqpgqxwakt2g7u9atsnr03gqcgmhcv38pt7mkd94q6shuwt",
	})
	require.Nil(t, err)
	require.Equal(t, "0.5.5", string(response.Data.ReturnData[0]))
}

func TestElrondProxy_GetRawBlockByHash(t *testing.T) {
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
	rawBlockDataBytes, err := json.Marshal(rawBlockData)

	httpClient := &mockHTTPClient{
		doCalled: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader(rawBlockDataBytes)),
			}, nil
		},
	}
	ep := NewElrondProxy("http://localhost:8079", httpClient)

	response, err := ep.GetRawBlockByHash(context.Background(), 0, "aaaa")
	require.Nil(t, err)

	ts := &testStruct{}
	err = json.Unmarshal(response, ts)
	require.Nil(t, err)
	require.Equal(t, expectedTs, ts)
}

func TestElrondProxy_GetRawBlockByNonce(t *testing.T) {
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
	rawBlockDataBytes, err := json.Marshal(rawBlockData)

	httpClient := &mockHTTPClient{
		doCalled: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader(rawBlockDataBytes)),
			}, nil
		},
	}
	ep := NewElrondProxy("http://localhost:8079", httpClient)

	response, err := ep.GetRawBlockByNonce(context.Background(), 0, 10)
	require.Nil(t, err)

	ts := &testStruct{}
	err = json.Unmarshal(response, ts)
	require.Nil(t, err)
	require.Equal(t, expectedTs, ts)
}

func TestElrondProxy_GetRawMiniBlockByHash(t *testing.T) {
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
	rawBlockDataBytes, err := json.Marshal(rawBlockData)

	httpClient := &mockHTTPClient{
		doCalled: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader(rawBlockDataBytes)),
			}, nil
		},
	}
	ep := NewElrondProxy("http://localhost:8079", httpClient)

	response, err := ep.GetRawMiniBlockByHash(context.Background(), 0, "aaaa")
	require.Nil(t, err)

	ts := &testStruct{}
	err = json.Unmarshal(response, ts)
	require.Nil(t, err)
	require.Equal(t, expectedTs, ts)
}

func TestElrondProxy_GetNonceAtEpochStart(t *testing.T) {
	expectedNonce := uint64(2)
	expectedNetworkStatus := &data.NetworkStatus{
		NonceAtEpochStart: expectedNonce,
	}
	statusResponse := &data.NetworkStatusResponse{
		Data: struct {
			Status *data.NetworkStatus "json:\"status\""
		}{
			Status: expectedNetworkStatus,
		},
	}
	statusResponseBytes, err := json.Marshal(statusResponse)

	httpClient := &mockHTTPClient{
		doCalled: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader(statusResponseBytes)),
			}, nil
		},
	}
	ep := NewElrondProxy("http://localhost:8079", httpClient)

	response, err := ep.GetNonceAtEpochStart(context.Background(), core.MetachainShardId)
	require.Nil(t, err)

	require.Equal(t, expectedNonce, response)
}

func TestElrondProxy_GetRatingsConfig(t *testing.T) {
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
	ratingsResponseBytes, err := json.Marshal(ratingsResponse)

	httpClient := &mockHTTPClient{
		doCalled: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader(ratingsResponseBytes)),
			}, nil
		},
	}
	ep := NewElrondProxy("http://localhost:8079", httpClient)

	response, err := ep.GetRatingsConfig(context.Background())
	require.Nil(t, err)

	require.Equal(t, expectedRatingsConfig, response)
}

func TestElrondProxy_GetEnableEpochsConfig(t *testing.T) {
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
	enableEpochsResponseBytes, err := json.Marshal(enableEpochsResponse)

	httpClient := &mockHTTPClient{
		doCalled: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader(enableEpochsResponseBytes)),
			}, nil
		},
	}
	ep := NewElrondProxy("http://localhost:8079", httpClient)

	response, err := ep.GetEnableEpochsConfig(context.Background())
	require.Nil(t, err)

	require.Equal(t, expectedEnableEpochsConfig, response)
}
