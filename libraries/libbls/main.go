package main

import (
	"C"
)

func main() {
}

//export generatePublicKey
func generatePublicKey(privateKey *C.char) *C.char {
	privateKeyHex := C.GoString(privateKey)
	publicKeyHex := doGeneratePublicKey(privateKeyHex)
	return C.CString(publicKeyHex)
}

//export signMessage
func signMessage(message *C.char, privateKey *C.char) *C.char {
	messageHex := C.GoString(message)
	privateKeyHex := C.GoString(privateKey)
	signatureHex := doSignMessage(messageHex, privateKeyHex)
	return C.CString(signatureHex)
}

//export verifyMessage
func verifyMessage(publicKey *C.char, message *C.char, signature *C.char) int {
	publicKeyHex := C.GoString(publicKey)
	messageHex := C.GoString(message)
	signatureHex := C.GoString(signature)
	ok := doVerifyMessage(publicKeyHex, messageHex, signatureHex)
	if ok {
		return 1
	}

	return 0
}

//export generatePrivateKey
func generatePrivateKey() *C.char {
	privateKeyHex := doGeneratePrivateKey()
	return C.CString(privateKeyHex)
}
