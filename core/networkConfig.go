package core

type NetworkConfig struct {
	MinGasLimit      int
	GasPerDataByte   int
	GasPriceModifier float32
	ChainID          string
}
