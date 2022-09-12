package authentication

// AuthClient defines the behavior of an authentication client
type AuthClient interface {
	GetAccessToken(address string, token string) (string, error)
}
