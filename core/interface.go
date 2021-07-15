package core

//AddressHandler will handle different implementations of an address
type AddressHandler interface {
	AddressAsBech32String() string
	AddressBytes() []byte
	IsValid() bool
	IsInterfaceNil() bool
}
