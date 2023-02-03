package main

import "C"

func main() {
}

//export generatePublicKey
func generatePublicKey(privateKey *C.char) *C.char {
	privateKeyHex := C.GoString(privateKey)
	publicKeyHex := doGeneratePublicKeyAsHex(privateKeyHex)
	return C.CString(publicKeyHex)
}

//export computeMessageSignature
func computeMessageSignature(message *C.char, privateKey *C.char) *C.char {
	messageHex := C.GoString(message)
	privateKeyHex := C.GoString(privateKey)
	signatureHex := doComputeMessageSignatureAsHex(messageHex, privateKeyHex)
	return C.CString(signatureHex)
}

//export verifyMessageSignature
func verifyMessageSignature(publicKey *C.char, message *C.char, signature *C.char) int {
	publicKeyHex := C.GoString(publicKey)
	messageHex := C.GoString(message)
	signatureHex := C.GoString(signature)
	ok := doVerifyMessageSignature(publicKeyHex, messageHex, signatureHex)
	if ok {
		return 1
	}

	return 0
}

//export generatePrivateKey
func generatePrivateKey() *C.char {
	privateKeyHex := doGeneratePrivateKeyAsHex()
	return C.CString(privateKeyHex)
}
