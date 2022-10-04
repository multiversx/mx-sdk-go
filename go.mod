module github.com/ElrondNetwork/elrond-sdk-erdgo

go 1.14

require (
	github.com/ElrondNetwork/elrond-go v1.3.42
	github.com/ElrondNetwork/elrond-go-core v1.1.19
	github.com/ElrondNetwork/elrond-go-crypto v1.0.1
	github.com/ElrondNetwork/elrond-go-logger v1.0.7
	github.com/ElrondNetwork/elrond-vm-common v1.3.15-rc1
	github.com/pborman/uuid v1.2.1
	github.com/stretchr/testify v1.7.1
	github.com/tyler-smith/go-bip39 v1.1.0
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4
)

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_2 v1.2.41 => github.com/ElrondNetwork/arwen-wasm-vm v1.2.41

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_3 v1.3.41 => github.com/ElrondNetwork/arwen-wasm-vm v1.3.41

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_4 v1.4.58 => github.com/ElrondNetwork/arwen-wasm-vm v1.4.58
