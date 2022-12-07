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

func NewAuthTokenHandler() *authTokenHandler {
	return &authTokenHandler{
		decodeHandler: base64.StdEncoding.DecodeString,
		encodeHandler: base64.StdEncoding.EncodeToString,
	}
}

// Decode decodes the given access token
func (th *authTokenHandler) Decode(accessToken string) (authentication.AuthToken, error) {
	token := NativeAuthToken{}
	var err error
	strs := strings.Split(accessToken, ".")
	token.Address, err = th.decodeHandler(strs[0])
	if err != nil {
		return nil, err
	}
	body, err := th.decodeHandler(strs[1])
	if err != nil {
		return nil, err
	}
	token.Signature = []byte(strs[2])
	strs = strings.Split(string(body), ".")
	token.Host = strs[0]
	token.BlockHash = strs[1]
	token.Ttl, err = strconv.ParseInt(strs[2], 10, 64)
	if err != nil {
		return nil, err
	}
	token.ExtraInfo = strs[3]

	return token, nil
}

// Encode encodes the given authentication token
func (th *authTokenHandler) Encode(authToken authentication.AuthToken) (string, error) {
	token, ok := authToken.(*NativeAuthToken)
	if !ok {
		return "", authentication.ErrCannotConvertToken
	}
	signature := token.Signature
	if len(signature) == 0 {
		return "", authentication.ErrNilSignature
	}

	encodedAddress := th.encodeHandler(token.Address)
	if len(encodedAddress) == 0 {
		return "", authentication.ErrNilAddress
	}

	encodedToken := th.encodeHandler(token.Body())
	if len(encodedToken) == 0 {
		return "", authentication.ErrNilBody
	}
	encodedSignature := th.encodeHandler(signature)
	if len(encodedSignature) == 0 {
		return "", authentication.ErrNilSignature
	}
	return fmt.Sprintf("%s.%s.%s", encodedAddress, encodedToken, encodedSignature), nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (th *authTokenHandler) IsInterfaceNil() bool {
	return th == nil
}
