package core

import (
	"encoding/hex"
	"net/url"
	"strconv"

	"github.com/multiversx/mx-chain-core-go/data/api"
)

const (
	// UrlParameterOnFinalBlock represents the name of a URL parameter
	UrlParameterOnFinalBlock = "onFinalBlock"
	// UrlParameterOnStartOfEpoch represents the name of a URL parameter
	UrlParameterOnStartOfEpoch = "onStartOfEpoch"
	// UrlParameterBlockNonce represents the name of a URL parameter
	UrlParameterBlockNonce = "blockNonce"
	// UrlParameterBlockHash represents the name of an URL parameter
	UrlParameterBlockHash = "blockHash"
	// UrlParameterBlockRootHash represents the name of an URL parameter
	UrlParameterBlockRootHash = "blockRootHash"
	// UrlParameterHintEpoch represents the name of an URL parameter
	UrlParameterHintEpoch = "hintEpoch"
)

// BuildUrlWithAccountQueryOptions builds an URL with block query parameters
// TODO: move this to mx-chain-core-go & remove also from mx-chain-proxy-go common/options.go
func BuildUrlWithAccountQueryOptions(path string, options api.AccountQueryOptions) string {
	u := url.URL{Path: path}
	query := u.Query()

	if options.OnFinalBlock {
		query.Set(UrlParameterOnFinalBlock, "true")
	}
	if options.OnStartOfEpoch.HasValue {
		query.Set(UrlParameterOnStartOfEpoch, strconv.Itoa(int(options.OnStartOfEpoch.Value)))
	}
	if options.BlockNonce.HasValue {
		query.Set(UrlParameterBlockNonce, strconv.FormatUint(options.BlockNonce.Value, 10))
	}
	if len(options.BlockHash) > 0 {
		query.Set(UrlParameterBlockHash, hex.EncodeToString(options.BlockHash))
	}
	if len(options.BlockRootHash) > 0 {
		query.Set(UrlParameterBlockRootHash, hex.EncodeToString(options.BlockRootHash))
	}
	if options.HintEpoch.HasValue {
		query.Set(UrlParameterHintEpoch, strconv.Itoa(int(options.HintEpoch.Value)))
	}

	u.RawQuery = query.Encode()
	return u.String()
}
