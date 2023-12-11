package crypto

const (
	CipherAlgorithmAes128Ctr    = "aes-128-ctr"
	KeyDerivationFunctionScrypt = "scrypt"
	RandomSaltLength            = 32
	RandomIvLength              = 16
	EncryptorVersion            = 4

	ScryptN     = 4096
	ScryptR     = 8
	ScryptP     = 1
	ScryptDKLen = 32
)

type PasswordBasedEncryptorDecryptor interface {
	Decrypt(encryptedData *PasswordEncryptedData, password string) ([]byte, error)
	Encrypt(data []byte, password string) (*PasswordEncryptedData, error)
}

type PasswordEncryptedData struct {
	ID         string
	Version    int
	Cipher     string
	CipherText string
	Iv         string
	KDF        string
	KDFParams  *KeyDerivationParams
	Salt       string
	MAC        string
}

type KeyDerivationParams struct {
	N     int // numIterations
	R     int // memFactor
	P     int // pFactor
	DKLen int
}
