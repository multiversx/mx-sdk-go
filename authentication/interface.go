package authentication

// AuthClient defines the behavior of an authentication client
type AuthClient interface {
	GetAccessToken() (string, error)
	IsInterfaceNil() bool
}

// AuthServer defines the behavior of an authentication server
type AuthServer interface {
	Validate(accessToken string) error
	IsInterfaceNil() bool
}

// AuthTokenHandler defines the behavior of an authentication token handler
type AuthTokenHandler interface {
	Decode(accessToken string) (AuthToken, error)
	Encode(authToken AuthToken) (string, error)
	IsInterfaceNil() bool
}

// AuthToken defines the behavior of an authentication token
type AuthToken interface {
	Body() []byte
}
