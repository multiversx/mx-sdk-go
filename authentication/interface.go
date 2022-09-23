package authentication

// AuthClient defines the behavior of an authentication client
type AuthClient interface {
	GetAccessToken() (string, error)
	IsInterfaceNil() bool
}
