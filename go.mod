module github.com/ElrondNetwork/elrond-sdk-erdgo

go 1.14

require (
	github.com/ElrondNetwork/elrond-go v1.3.7-0.20220310094258-ead8cd541713
	github.com/ElrondNetwork/elrond-go-core v1.1.14
	github.com/ElrondNetwork/elrond-go-crypto v1.0.1
	github.com/ElrondNetwork/elrond-go-logger v1.0.6
	github.com/ElrondNetwork/elrond-vm-common v1.3.2
	github.com/btcsuite/websocket v0.0.0-20150119174127-31079b680792
	github.com/gin-contrib/cors v0.0.0-20190301062745-f9e10995c85a
	github.com/gin-gonic/gin v1.7.7
	github.com/pborman/uuid v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/tyler-smith/go-bip39 v1.1.0
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	golang.org/x/oauth2 v0.0.0-20190226205417-e64efc72b421
)

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_2 v1.2.35 => github.com/ElrondNetwork/arwen-wasm-vm v1.2.35

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_3 v1.3.35 => github.com/ElrondNetwork/arwen-wasm-vm v1.3.35

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_4 v1.4.40 => github.com/ElrondNetwork/arwen-wasm-vm v1.4.40
