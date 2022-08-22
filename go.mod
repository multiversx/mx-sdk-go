module github.com/ElrondNetwork/elrond-sdk-erdgo

go 1.14

require (
	github.com/ElrondNetwork/elrond-go v1.3.36-0.20220711104849-8a728880e1e1
	github.com/ElrondNetwork/elrond-go-core v1.1.16-0.20220711092037-f35a3a0faf0f
	github.com/ElrondNetwork/elrond-go-crypto v1.0.1
	github.com/ElrondNetwork/elrond-go-logger v1.0.7
	github.com/ElrondNetwork/elrond-vm-common v1.3.13-0.20220708125052-5343b3b65f3e
	github.com/pborman/uuid v1.2.1
	github.com/stretchr/testify v1.7.1
	github.com/tyler-smith/go-bip39 v1.1.0
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4
)

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_2 v1.2.40 => github.com/ElrondNetwork/arwen-wasm-vm v1.2.41-0.20220708132302-7590ce4497ec

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_3 v1.3.40 => github.com/ElrondNetwork/arwen-wasm-vm v1.3.41-0.20220708143408-ab05f75aa3a6

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_4 v1.4.54-rc3 => github.com/ElrondNetwork/arwen-wasm-vm v1.4.57-0.20220708144802-8e7c159a12c6
