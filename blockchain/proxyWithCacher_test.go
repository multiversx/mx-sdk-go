package blockchain

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func Test1(t *testing.T) {
	responseBytes := []byte(`{"data":{"metrics":{"erd_dev_rewards":"0","erd_epoch_for_economics_data":263,"erd_inflation":"5869888769785838708144","erd_total_fees":"51189055176110000000","erd_total_staked_value":"9963775651405816710680128","erd_total_supply":"21556417261819025351089574","erd_total_top_up_value":"1146275808171377418645274"}},"code":"successful"}`)
	httpClient := &mockHTTPClient{
		doCalled: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader(responseBytes)),
			}, nil
		},
	}
	ep := NewElrondProxy("http://localhost:8079", httpClient)

	networkEconomics, _ := ep.GetNetworkEconomics(context.Background())
	fmt.Println(reflect.TypeOf(networkEconomics))
}
