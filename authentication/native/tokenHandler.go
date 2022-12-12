package native

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/authentication"
)

// authTokenHandler will handle encoding and decoding native authentication tokens
type authTokenHandler struct {
	decodeHandler func(s string) ([]byte, error)
	encodeHandler func(src []byte) string
}

// NewAuthTokenHandler returns a new instance of a native authentication token handler
func NewAuthTokenHandler() *authTokenHandler {
	return &authTokenHandler{
		decodeHandler: decodeHandler,
		encodeHandler: base64.StdEncoding.EncodeToString,
	}
}

// Decode decodes the given access token
func (th *authTokenHandler) Decode(accessToken string) (authentication.AuthToken, error) {
	token := AuthToken{}
	var err error
	strs := strings.Split(accessToken, ".")
	token.address, err = th.decodeHandler(strs[0])
	if err != nil {
		return nil, err
	}
	body, err := th.decodeHandler(strs[1])
	if err != nil {
		return nil, err
	}
	token.signature = []byte(strs[2])
	strs = strings.Split(string(body), ".")
	token.blockHash = strs[0]
	token.ttl, err = strconv.ParseInt(strs[1], 10, 64)
	if err != nil {
		return nil, err
	}
	token.extraInfo = strs[2]

	return token, nil
}

// Encode encodes the given authentication token
func (th *authTokenHandler) Encode(authToken authentication.AuthToken) (string, error) {
	signature := authToken.GetSignature()
	if len(signature) == 0 {
		return "", authentication.ErrNilSignature
	}

	encodedAddress := th.encodeHandler(authToken.GetAddress())
	if len(encodedAddress) == 0 {
		return "", authentication.ErrNilAddress
	}

	encodedToken := th.encodeHandler(authToken.GetBody())
	if len(encodedToken) == 0 {
		return "", authentication.ErrNilBody
	}

	return fmt.Sprintf("%s.%s.%s", encodedAddress, encodedToken, signature), nil
}

func decodeHandler(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(strings.TrimRight(s, "="))
}

// IsInterfaceNil returns true if there is no value under the interface
func (th *authTokenHandler) IsInterfaceNil() bool {
	return th == nil
}
