package fetchers

const (
	quoteUSDFiat = "USD"
	quoteUSDT    = "USDT"

	// BinanceName defines the Binance exchange
	BinanceName = "Binance"
	// BitfinexName defines the Bitfinex exchange
	BitfinexName = "Bitfinex"
	// CryptocomName defines the crypto.com exchange
	CryptocomName = "Crypto.com"
	// GeminiName defines the Gemini exchange
	GeminiName = "Gemini"
	// HitbtcName defines the HitBTC exchange
	HitbtcName = "HitBTC"
	// HuobiName defines the Huobi exchange
	HuobiName = "Huobi"
	// KrakenName defines the Kraken exchange
	KrakenName = "Kraken"
	// OkexName defines the Okex exchange
	OkexName = "Okex"
)

// ImplementedFetchers is the list of all implemented exchange fetchers
var ImplementedFetchers = []string{BinanceName, BitfinexName, CryptocomName, GeminiName, HitbtcName, HuobiName, KrakenName, OkexName}
