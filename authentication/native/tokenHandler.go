package native

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/authentication"
)

// authTokenHandler will handle encoding and decoding native authentication tokens
type authTokenHandler struct {
	decodeHandler    func(s string) ([]byte, error)
	hexDecodeHandler func(s string) ([]byte, error)
	encodeHandler    func(src []byte) string
}

// NewAuthTokenHandler returns a new instance of a native authentication token handler
func NewAuthTokenHandler() *authTokenHandler {
	return &authTokenHandler{
		decodeHandler:    decodeHandler,
		hexDecodeHandler: hex.DecodeString,
		encodeHandler:    encodeHandler,
	}
}

// Decode decodes the given access token
func (th *authTokenHandler) Decode(accessToken string) (authentication.AuthToken, error) {
	token := &AuthToken{}
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
	token.signature, err = th.hexDecodeHandler(strs[2])
	if err != nil {
		return nil, err
	}
	strs = strings.Split(string(body), ".")
	token.host, err = th.decodeHandler(strs[0])
	if err != nil {
		return nil, err
	}
	token.blockHash = strs[1]
	token.ttl, err = strconv.ParseInt(strs[2], 10, 64)
	if err != nil {
		return nil, err
	}
	token.extraInfo, err = th.decodeHandler(strs[3])
	if err != nil {
		return nil, err
	}

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

	encodedToken := th.encodeHandler(th.GetUnsignedToken(authToken))
	if len(encodedToken) == 0 {
		return "", authentication.ErrNilBody
	}

	return fmt.Sprintf("%s.%s.%x", encodedAddress, encodedToken, signature), nil
}

// GetUnsignedToken returns an unsigned authentication token
func (th *authTokenHandler) GetUnsignedToken(token authentication.AuthToken) []byte {
	encodedExtraInfo := th.encodeHandler(token.GetExtraInfo())
	encodedHost := th.encodeHandler(token.GetHost())
	return []byte(fmt.Sprintf("%s.%s.%d.%s", encodedHost, token.GetBlockHash(), token.GetTtl(), encodedExtraInfo))
}

// GetSignableMessage returns a message to be signed
func (th *authTokenHandler) GetSignableMessage(address, unsignedToken []byte) []byte {
	return []byte(fmt.Sprintf("%s%s", address, unsignedToken))
}

func decodeHandler(source string) ([]byte, error) {
	source = strings.ReplaceAll(source, "-", "+")
	source = strings.ReplaceAll(source, "_", "/")
	source = strings.TrimRight(source, "=")
	return base64.RawStdEncoding.DecodeString(source)
}

func encodeHandler(source []byte) string {
	encoded := base64.StdEncoding.EncodeToString(source)
	encoded = strings.ReplaceAll(encoded, "+", "-")
	encoded = strings.ReplaceAll(encoded, "/", "_")
	encoded = strings.TrimRight(encoded, "=")
	return encoded
}

// IsInterfaceNil returns true if there is no value under the interface
func (th *authTokenHandler) IsInterfaceNil() bool {
	return th == nil
}
