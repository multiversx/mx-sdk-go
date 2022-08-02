package interactors

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"io"
	"io/ioutil"
	"os"

	"github.com/ElrondNetwork/elrond-go-crypto/signing"
	"github.com/ElrondNetwork/elrond-go-crypto/signing/ed25519"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/pborman/uuid"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/scrypt"
)

const (
	mnemonicBitSize = 256
	egldCoinType    = uint32(508)
	hardened        = uint32(0x80000000)
	keystoreVersion = 4
	keyHeaderKDF    = "scrypt"
	scryptN         = 4096
	scryptR         = 8
	scryptP         = 1
	scryptDKLen     = 32
	addressLen      = 32
)

type bip32Path []uint32

type bip32 struct {
	Key       []byte
	ChainCode []byte
}

var suite = ed25519.NewEd25519()
var keyGenerator = signing.NewKeyGenerator(suite)

type encryptedKeyJSONV4 struct {
	Address string `json:"address"`
	Bech32  string `json:"bech32"`
	Crypto  struct {
		Cipher       string `json:"cipher"`
		CipherText   string `json:"ciphertext"`
		CipherParams struct {
			IV string `json:"iv"`
		} `json:"cipherparams"`
		KDF       string `json:"kdf"`
		KDFParams struct {
			DkLen int    `json:"dklen"`
			Salt  string `json:"salt"`
			N     int    `json:"n"`
			R     int    `json:"r"`
			P     int    `json:"p"`
		} `json:"kdfparams"`
		MAC string `json:"mac"`
	} `json:"crypto"`
	Id      string `json:"id"`
	Version int    `json:"version"`
}

type wallet struct {
}

// NewWallet creates a new wallet instance
func NewWallet() *wallet {
	return &wallet{}
}

// GenerateMnemonic will generate a new mnemonic value using the bip39 implementation
func (w *wallet) GenerateMnemonic() (data.Mnemonic, error) {
	entropy, err := bip39.NewEntropy(mnemonicBitSize)
	if err != nil {
		return "", err
	}

	mnemonicString, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return data.Mnemonic(mnemonicString), nil
}

// GetPrivateKeyFromMnemonic generates a private key based on mnemonic, account and address index
func (w *wallet) GetPrivateKeyFromMnemonic(mnemonic data.Mnemonic, account, addressIndex uint32) []byte {
	seed := w.CreateSeedFromMnemonic(mnemonic)
	privateKey := w.GetPrivateKeyFromSeed(seed, account, addressIndex)
	return privateKey
}

// GetPrivateKeyFromSeed generates a private key based on seed, account and address index
func (w *wallet) GetPrivateKeyFromSeed(seed []byte, account, addressIndex uint32) []byte {
	var egldPath = bip32Path{
		44 | hardened,
		egldCoinType | hardened,
		hardened, // account
		hardened,
		hardened, // addressIndex
	}

	egldPath[2] = account | hardened
	egldPath[4] = addressIndex | hardened
	keyData := derivePrivateKey(seed, egldPath)

	return keyData.Key
}

// GetSeedFromMnemonic creates a seed for a given mnemonic
func (w *wallet) CreateSeedFromMnemonic(mnemonic data.Mnemonic) []byte {
	seed := bip39.NewSeed(string(mnemonic), "")
	return seed
}

func derivePrivateKey(seed []byte, path bip32Path) *bip32 {
	b := &bip32{}
	digest := hmac.New(sha512.New, []byte("ed25519 seed"))
	digest.Write(seed)
	intermediary := digest.Sum(nil)
	serializedKeyLen := 32
	serializedChildIndexLen := 4
	hardenedChildPadding := byte(0x00)
	b.Key = intermediary[:serializedKeyLen]
	b.ChainCode = intermediary[serializedKeyLen:]
	for _, childIdx := range path {
		buff := make([]byte, 1+serializedKeyLen+4)
		buff[0] = hardenedChildPadding
		copy(buff[1:1+serializedKeyLen], b.Key)
		binary.BigEndian.PutUint32(buff[1+serializedKeyLen:1+serializedKeyLen+serializedChildIndexLen], childIdx)
		digest = hmac.New(sha512.New, b.ChainCode)
		digest.Write(buff)
		intermediary = digest.Sum(nil)
		b.Key = intermediary[:serializedKeyLen]
		b.ChainCode = intermediary[serializedKeyLen:]
	}

	return b
}

// GetAddressFromPrivateKey generates the bech32 address from a private key
func (w *wallet) GetAddressFromPrivateKey(privateKeyBytes []byte) (core.AddressHandler, error) {
	privateKey, err := keyGenerator.PrivateKeyFromByteArray(privateKeyBytes)
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.GeneratePublic()
	publicKeyBytes, err := publicKey.ToByteArray()
	if err != nil {
		return nil, err
	}

	return data.NewAddressFromBytes(publicKeyBytes), nil
}

// LoadPrivateKeyFromJsonFile loads a password encrypted private key from a .json file
func (w *wallet) LoadPrivateKeyFromJsonFile(filename string, password string) ([]byte, error) {
	buff, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	key := &encryptedKeyJSONV4{}
	err = json.Unmarshal(buff, key)
	if err != nil {
		return nil, err
	}

	mac, err := hex.DecodeString(key.Crypto.MAC)
	if err != nil {
		return nil, err
	}

	iv, err := hex.DecodeString(key.Crypto.CipherParams.IV)
	if err != nil {
		return nil, err
	}

	cipherText, err := hex.DecodeString(key.Crypto.CipherText)
	if err != nil {
		return nil, err
	}

	salt, err := hex.DecodeString(key.Crypto.KDFParams.Salt)
	if err != nil {
		return nil, err
	}

	derivedKey, err := scrypt.Key([]byte(password), salt,
		key.Crypto.KDFParams.N,
		key.Crypto.KDFParams.R,
		key.Crypto.KDFParams.P,
		key.Crypto.KDFParams.DkLen)
	if err != nil {
		return nil, err
	}

	hash := hmac.New(sha256.New, derivedKey[16:32])
	_, err = hash.Write(cipherText)
	if err != nil {
		return nil, err
	}

	sha := hash.Sum(nil)
	if !bytes.Equal(sha, mac) {
		return nil, ErrWrongPassword
	}

	aesBlock, err := aes.NewCipher(derivedKey[:16])
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(aesBlock, iv)
	privateKey := make([]byte, len(cipherText))
	stream.XORKeyStream(privateKey, cipherText)

	if len(privateKey) > 32 {
		privateKey = privateKey[:32]
	}

	address, err := w.GetAddressFromPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	isSameAccount := hex.EncodeToString(address.AddressBytes()) == key.Address &&
		address.AddressAsBech32String() == key.Bech32
	if !isSameAccount {
		return nil, ErrDifferentAccountRecovered
	}
	return privateKey, nil
}

// SavePrivateKeyToJsonFile saves a password encrypted private key to a .json file
func (w *wallet) SavePrivateKeyToJsonFile(privateKey []byte, password string, filename string) error {
	salt := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return err
	}

	derivedKey, err := scrypt.Key([]byte(password), salt, scryptN, scryptR, scryptP, scryptDKLen)
	if err != nil {
		return err
	}

	encryptKey := derivedKey[:16]
	iv := make([]byte, aes.BlockSize) // 16
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return err
	}

	aesBlock, err := aes.NewCipher(encryptKey)
	if err != nil {
		return err
	}

	stream := cipher.NewCTR(aesBlock, iv)
	cipherText := make([]byte, len(privateKey))
	stream.XORKeyStream(cipherText, privateKey)

	hash := hmac.New(sha256.New, derivedKey[16:32])
	_, err = hash.Write(cipherText)
	if err != nil {
		return err
	}

	mac := hash.Sum(nil)

	address, err := w.GetAddressFromPrivateKey(privateKey)
	if err != nil {
		return err
	}

	keystoreJson := &encryptedKeyJSONV4{
		Bech32:  address.AddressAsBech32String(),
		Address: hex.EncodeToString(address.AddressBytes()),
		Version: keystoreVersion,
		Id:      uuid.New(),
	}
	keystoreJson.Crypto.CipherParams.IV = hex.EncodeToString(iv)
	keystoreJson.Crypto.Cipher = "aes-128-ctr"
	keystoreJson.Crypto.CipherText = hex.EncodeToString(cipherText)
	keystoreJson.Crypto.KDF = keyHeaderKDF
	keystoreJson.Crypto.MAC = hex.EncodeToString(mac)
	keystoreJson.Crypto.KDFParams.N = scryptN
	keystoreJson.Crypto.KDFParams.R = scryptR
	keystoreJson.Crypto.KDFParams.P = scryptP
	keystoreJson.Crypto.KDFParams.DkLen = scryptDKLen
	keystoreJson.Crypto.KDFParams.Salt = hex.EncodeToString(salt)

	buff, err := json.Marshal(keystoreJson)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, buff, 0644)
}

// LoadPrivateKeyFromPemFile loads a private key from a .pem file
func (w *wallet) LoadPrivateKeyFromPemFile(filename string) ([]byte, error) {
	buff, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return w.LoadPrivateKeyFromPemData(buff)
}

// LoadPrivateKeyFromPemData returns the private key from decoded pem data
func (w *wallet) LoadPrivateKeyFromPemData(buff []byte) ([]byte, error) {
	blk, _ := pem.Decode(buff)
	if blk == nil {
		return nil, ErrInvalidPemFile
	}

	privKey := blk.Bytes
	if len(privKey) > addressLen*2 {
		privKey = privKey[:addressLen*2]
	}

	return hex.DecodeString(string(privKey))
}

// SavePrivateKeyToPemFile saves the private key in a .pem file
func (w *wallet) SavePrivateKeyToPemFile(privateKey []byte, filename string) error {
	address, err := w.GetAddressFromPrivateKey(privateKey)
	if err != nil {
		return err
	}
	if len(privateKey) == addressLen {
		privateKey = append(privateKey, address.AddressBytes()...)
	}
	blk := pem.Block{
		Type:  "PRIVATE KEY for " + address.AddressAsBech32String(),
		Bytes: []byte(hex.EncodeToString(privateKey)),
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	return pem.Encode(file, &blk)
}
