package serde

// Deserializer defines the methods used to populate object fields based on the received byte array
type Deserializer interface {
	CreateStruct(obj interface{}, buff []byte) (uint64, error)
	CreatePrimitiveDataType(obj interface{}, buff []byte) error
}
