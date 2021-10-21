package serde

type Deserializer interface {
	CreateStruct(obj interface{}, buff []byte) (uint64, error)
}
