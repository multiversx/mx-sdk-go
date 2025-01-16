package fetchers

const (
	quoteUSDFiat = "USD"
	quoteUSDT    = "USDT"

	// BinanceName defines the Binance exchange name
	BinanceName = "Binance"
	// BitfinexName defines the Bitfinex exchange name
	BitfinexName = "Bitfinex"
	// CryptocomName defines the crypto.com exchange name
	CryptocomName = "Crypto.com"
	// GeminiName defines the Gemini exchange name
	GeminiName = "Gemini"
	// HitbtcName defines the HitBTC exchange name
	HitbtcName = "HitBTC"
	// HuobiName defines the Huobi exchange name
	HuobiName = "Huobi"
	// KrakenName defines the Kraken exchange name
	KrakenName = "Kraken"
	// OkxName defines the Okx exchange name
	OkxName = "Okx"
	// XExchangeName defines the XExchange name
	XExchangeName = "XExchange"
	// EVMGasPriceStation defines an EVM gas station that will push gas prices as a full token pair price
	EVMGasPriceStation = "EVM gas price station"
)

// ImplementedFetchers is the map of all implemented exchange fetchers
var ImplementedFetchers = map[string]struct{}{
	BinanceName:   {},
	BitfinexName:  {},
	CryptocomName: {},
	GeminiName:    {},
	HitbtcName:    {},
	HuobiName:     {},
	KrakenName:    {},
	OkxName:       {},
	XExchangeName: {},
}
