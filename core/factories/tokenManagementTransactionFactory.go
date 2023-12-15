package factories

import (
	"encoding/hex"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-sdk-go/core"
	"math/big"
	"strconv"
	"strings"
)

type TokenType string

const (
	NFT  TokenType = "NFT"
	SFT  TokenType = "SFT"
	META TokenType = "META"
	FNG  TokenType = "FNG"
)

type TokenManagementTransactionFactory interface {
	CreateTransactionForIssuingFungible(
		sender core.Address,
		tokenName string,
		tokenTicker string,
		initialSupply int,
		numDecimals int,
		canFreeze bool,
		canWipe bool,
		canPause bool,
		canChangeOwner bool,
		canUpgrade bool,
		canAddSpecialRoles bool,
	) (*core.Transaction, error)

	CreateTransactionForIssuingSemiFungible(
		sender core.Address,
		tokenName string,
		tokenTicker string,
		canFreeze bool,
		canWipe bool,
		canPause bool,
		canTransferNFTCreateRole bool,
		canChangeOwner bool,
		canUpgrade bool,
		canAddSpecialRoles bool,
	) (*core.Transaction, error)

	CreateTransactionForIssuingNonFungible(
		sender core.Address,
		tokenName string,
		tokenTicker string,
		canFreeze bool,
		canWipe bool,
		canPause bool,
		canTransferNFTCreateRole bool,
		canChangeOwner bool,
		canUpgrade bool,
		canAddSpecialRoles bool,
	) (*core.Transaction, error)

	CreateTransactionForRegisteringMetaESDT(
		sender core.Address,
		tokenName string,
		tokenTicker string,
		numDecimals int,
		canFreeze bool,
		canWipe bool,
		canPause bool,
		canTransferNFTCreateRole bool,
		canChangeOwner bool,
		canUpgrade bool,
		canAddSpecialRoles bool,
	) (*core.Transaction, error)

	CreateTransactionForRegisteringAndSettingRoles(
		sender core.Address,
		tokenName string,
		tokenTicker string,
		tokenType TokenType,
		numDecimals int,
	) (*core.Transaction, error)

	CreateTransactionForSettingBurnRoleGlobally(
		sender core.Address,
		tokenIdentifier string,
	) (*core.Transaction, error)

	CreateTransactionForSettingSpecialRoleOnFungibleToken(
		sender core.Address,
		user core.Address,
		tokenIdentifier string,
		addRoleLocalMint bool,
		addRoleLocalBurn bool,
	) (*core.Transaction, error)

	CreateTransactionForSettingSpecialRoleOnSemiFungibleToken(
		sender core.Address,
		user core.Address,
		tokenIdentifier string,
		addRoleNFTCreate bool,
		addRoleNFTBurn bool,
		addRoleNFTAddQuantity bool,
		addRoleESDTTransferRole bool,
	) (*core.Transaction, error)

	CreateTransactionForSettingSpecialRoleOnNonFungibleToken(
		sender core.Address,
		user core.Address,
		tokenIdentifier string,
		addRoleNFTCreate bool,
		addRoleNFTBurn bool,
		addRoleNFTUpdateAttributes bool,
		addRoleNFTAddURI bool,
		addRoleESDTTransferRole bool,
	) (*core.Transaction, error)

	CreateTransactionForCreatingNFT(
		sender core.Address,
		tokenIdentifier string,
		initialQuantity int,
		name string,
		royalties int,
		hash string,
		attributes []byte,
		URIs []string,
	) (*core.Transaction, error)

	CreateTransactionForPausing(
		sender core.Address,
		tokenIdentifier string,
	) (*core.Transaction, error)

	CreateTransactionForUnpausing(
		sender core.Address,
		tokenIdentifier string,
	) (*core.Transaction, error)

	CreateTransactionForFreezing(
		sender core.Address,
		user core.Address,
		tokenIdentifier string,
	) (*core.Transaction, error)

	CreateTransactionForUnfreezing(
		sender core.Address,
		user core.Address,
		tokenIdentifier string,
	) (*core.Transaction, error)

	CreateTransactionForLocalMinting(
		sender core.Address,
		tokenIdentifier string,
		supplyToMint int,
	) (*core.Transaction, error)

	CreateTransactionForLocalBurning(
		sender core.Address,
		tokenIdentifier string,
		supplyToBurn int,
	) (*core.Transaction, error)

	CreateTransactionForUpdatingAttributes(
		sender core.Address,
		tokenIdentifier string,
		tokenNonce int,
		attributes []byte,
	) (*core.Transaction, error)

	CreateTransactionForAddingQuantity(
		sender core.Address,
		tokenIdentifier string,
		tokenNonce int,
		quantityToAdd int,
	) (*core.Transaction, error)

	CreateTransactionForBurningQuantity(
		sender core.Address,
		tokenIdentifier string,
		tokenNonce int,
		quantityToAdd int,
	) (*core.Transaction, error)
}

type tokenManagementTransactionFactory struct {
	Config *core.Config
}

func NewTokenManagementTransactionFactory(config *core.Config) TokenManagementTransactionFactory {
	return &tokenManagementTransactionFactory{Config: config}
}

func (t *tokenManagementTransactionFactory) CreateTransactionForIssuingFungible(
	sender core.Address,
	tokenName string,
	tokenTicker string,
	initialSupply int,
	numDecimals int,
	canFreeze bool,
	canWipe bool,
	canPause bool,
	canChangeOwner bool,
	canUpgrade bool,
	canAddSpecialRoles bool,
) (*core.Transaction, error) {
	notifyAboutUnsettingBurnRoleGlobally()

	parts := []string{
		"issue",
		hex.EncodeToString([]byte(tokenName)),
		hex.EncodeToString([]byte(tokenTicker)),
		hex.EncodeToString(core.EncodeUnsignedNumber(uint64(initialSupply))),
		hex.EncodeToString(core.EncodeUnsignedNumber(uint64(numDecimals))),
	}

	if canFreeze {
		parts = append(parts, hex.EncodeToString([]byte("canFreeze")), hex.EncodeToString([]byte("true")))
	}
	if canWipe {
		parts = append(parts, hex.EncodeToString([]byte("canWipe")), hex.EncodeToString([]byte("true")))
	}
	if canPause {
		parts = append(parts, hex.EncodeToString([]byte("canPause")), hex.EncodeToString([]byte("true")))
	}
	if canChangeOwner {
		parts = append(parts, hex.EncodeToString([]byte("canChangeOwner")), hex.EncodeToString([]byte("true")))
	}
	if canUpgrade {
		parts = append(parts, hex.EncodeToString([]byte("canUpgrade")),
			hex.EncodeToString([]byte(strconv.FormatBool(canUpgrade))))
	}
	if canAddSpecialRoles {
		parts = append(parts, hex.EncodeToString([]byte("canAddSpecialRoles")),
			hex.EncodeToString([]byte(strconv.FormatBool(canUpgrade))))
	}

	return core.NewTransactionBuilder().
		WithConfig(t.Config).
		WithSender(sender).
		WithReceiver(t.Config.ESDTContractAddress).
		WithAmount(big.NewInt(int64(t.Config.IssueCost))).
		WithProvidedGasLimit(uint64(t.Config.GasLimitIssue)).
		WithAddDataMovementGas(true).
		WithDataParts(parts).
		Build()
}

func (t *tokenManagementTransactionFactory) CreateTransactionForIssuingSemiFungible(
	sender core.Address,
	tokenName string,
	tokenTicker string,
	canFreeze bool,
	canWipe bool,
	canPause bool,
	canTransferNFTCreateRole bool,
	canChangeOwner bool,
	canUpgrade bool,
	canAddSpecialRoles bool,
) (*core.Transaction, error) {
	notifyAboutUnsettingBurnRoleGlobally()

	parts := []string{
		"issueSemiFungible",
		hex.EncodeToString([]byte(tokenName)),
		hex.EncodeToString([]byte(tokenTicker)),
	}

	if canFreeze {
		parts = append(parts, hex.EncodeToString([]byte("canFreeze")), hex.EncodeToString([]byte("true")))
	}
	if canWipe {
		parts = append(parts, hex.EncodeToString([]byte("canWipe")), hex.EncodeToString([]byte("true")))
	}
	if canPause {
		parts = append(parts, hex.EncodeToString([]byte("canPause")), hex.EncodeToString([]byte("true")))
	}
	if canTransferNFTCreateRole {
		parts = append(parts, hex.EncodeToString([]byte("canTransferNFTCreateRole")), hex.EncodeToString([]byte("true")))
	}
	if canChangeOwner {
		parts = append(parts, hex.EncodeToString([]byte("canChangeOwner")), hex.EncodeToString([]byte("true")))
	}
	if canUpgrade {
		parts = append(parts, hex.EncodeToString([]byte("canUpgrade")),
			hex.EncodeToString([]byte(strconv.FormatBool(canUpgrade))))
	}
	if canAddSpecialRoles {
		parts = append(parts, hex.EncodeToString([]byte("canAddSpecialRoles")),
			hex.EncodeToString([]byte(strconv.FormatBool(canUpgrade))))
	}

	return core.NewTransactionBuilder().
		WithConfig(t.Config).
		WithSender(sender).
		WithReceiver(t.Config.ESDTContractAddress).
		WithAmount(big.NewInt(int64(t.Config.IssueCost))).
		WithProvidedGasLimit(uint64(t.Config.GasLimitIssue)).
		WithAddDataMovementGas(true).
		WithDataParts(parts).
		Build()
}

func (t *tokenManagementTransactionFactory) CreateTransactionForIssuingNonFungible(
	sender core.Address,
	tokenName string,
	tokenTicker string,
	canFreeze bool,
	canWipe bool,
	canPause bool,
	canTransferNFTCreateRole bool,
	canChangeOwner bool,
	canUpgrade bool,
	canAddSpecialRoles bool,
) (*core.Transaction, error) {
	notifyAboutUnsettingBurnRoleGlobally()

	parts := []string{
		"issueNonFungible",
		hex.EncodeToString([]byte(tokenName)),
		hex.EncodeToString([]byte(tokenTicker)),
	}

	if canFreeze {
		parts = append(parts, hex.EncodeToString([]byte("canFreeze")), hex.EncodeToString([]byte("true")))
	}
	if canWipe {
		parts = append(parts, hex.EncodeToString([]byte("canWipe")), hex.EncodeToString([]byte("true")))
	}
	if canPause {
		parts = append(parts, hex.EncodeToString([]byte("canPause")), hex.EncodeToString([]byte("true")))
	}
	if canTransferNFTCreateRole {
		parts = append(parts, hex.EncodeToString([]byte("canTransferNFTCreateRole")), hex.EncodeToString([]byte("true")))
	}
	if canChangeOwner {
		parts = append(parts, hex.EncodeToString([]byte("canChangeOwner")), hex.EncodeToString([]byte("true")))
	}
	if canUpgrade {
		parts = append(parts, hex.EncodeToString([]byte("canUpgrade")),
			hex.EncodeToString([]byte(strconv.FormatBool(canUpgrade))))
	}
	if canAddSpecialRoles {
		parts = append(parts, hex.EncodeToString([]byte("canAddSpecialRoles")),
			hex.EncodeToString([]byte(strconv.FormatBool(canUpgrade))))
	}

	return core.NewTransactionBuilder().
		WithConfig(t.Config).
		WithSender(sender).
		WithReceiver(t.Config.ESDTContractAddress).
		WithAmount(big.NewInt(int64(t.Config.IssueCost))).
		WithProvidedGasLimit(uint64(t.Config.GasLimitIssue)).
		WithAddDataMovementGas(true).
		WithDataParts(parts).
		Build()
}

func (t *tokenManagementTransactionFactory) CreateTransactionForRegisteringMetaESDT(
	sender core.Address,
	tokenName string,
	tokenTicker string,
	numDecimals int,
	canFreeze bool,
	canWipe bool,
	canPause bool,
	canTransferNFTCreateRole bool,
	canChangeOwner bool,
	canUpgrade bool,
	canAddSpecialRoles bool,
) (*core.Transaction, error) {
	notifyAboutUnsettingBurnRoleGlobally()

	parts := []string{
		"registerMetaESDT",
		hex.EncodeToString([]byte(tokenName)),
		hex.EncodeToString([]byte(tokenTicker)),
		hex.EncodeToString(core.EncodeUnsignedNumber(uint64(numDecimals))),
	}

	if canFreeze {
		parts = append(parts, hex.EncodeToString([]byte("canFreeze")), hex.EncodeToString([]byte("true")))
	}
	if canWipe {
		parts = append(parts, hex.EncodeToString([]byte("canWipe")), hex.EncodeToString([]byte("true")))
	}
	if canPause {
		parts = append(parts, hex.EncodeToString([]byte("canPause")), hex.EncodeToString([]byte("true")))
	}
	if canTransferNFTCreateRole {
		parts = append(parts, hex.EncodeToString([]byte("canTransferNFTCreateRole")), hex.EncodeToString([]byte("true")))
	}
	if canChangeOwner {
		parts = append(parts, hex.EncodeToString([]byte("canChangeOwner")), hex.EncodeToString([]byte("true")))
	}
	if canUpgrade {
		parts = append(parts, hex.EncodeToString([]byte("canUpgrade")),
			hex.EncodeToString([]byte(strconv.FormatBool(canUpgrade))))
	}
	if canAddSpecialRoles {
		parts = append(parts, hex.EncodeToString([]byte("canAddSpecialRoles")),
			hex.EncodeToString([]byte(strconv.FormatBool(canUpgrade))))
	}

	return core.NewTransactionBuilder().
		WithConfig(t.Config).
		WithSender(sender).
		WithReceiver(t.Config.ESDTContractAddress).
		WithAmount(big.NewInt(int64(t.Config.IssueCost))).
		WithProvidedGasLimit(uint64(t.Config.GasLimitIssue)).
		WithAddDataMovementGas(true).
		WithDataParts(parts).
		Build()
}

func (t *tokenManagementTransactionFactory) CreateTransactionForRegisteringAndSettingRoles(
	sender core.Address,
	tokenName string,
	tokenTicker string,
	tokenType TokenType,
	numDecimals int,
) (*core.Transaction, error) {
	notifyAboutUnsettingBurnRoleGlobally()

	parts := []string{
		"registerAndSetAllRoles",
		hex.EncodeToString([]byte(tokenName)),
		hex.EncodeToString([]byte(tokenTicker)),
		hex.EncodeToString([]byte(tokenType)),
		hex.EncodeToString(core.EncodeUnsignedNumber(uint64(numDecimals))),
	}

	return core.NewTransactionBuilder().
		WithConfig(t.Config).
		WithSender(sender).
		WithReceiver(t.Config.ESDTContractAddress).
		WithAmount(big.NewInt(int64(t.Config.IssueCost))).
		WithProvidedGasLimit(uint64(t.Config.GasLimitIssue)).
		WithAddDataMovementGas(true).
		WithDataParts(parts).
		Build()
}

func (t *tokenManagementTransactionFactory) CreateTransactionForSettingBurnRoleGlobally(
	sender core.Address,
	tokenIdentifier string,
) (*core.Transaction, error) {
	parts := []string{
		"setBurnRoleGlobally",
		hex.EncodeToString([]byte(tokenIdentifier)),
	}

	return core.NewTransactionBuilder().
		WithConfig(t.Config).
		WithSender(sender).
		WithReceiver(t.Config.ESDTContractAddress).
		WithAmount(big.NewInt(0)).
		WithProvidedGasLimit(uint64(t.Config.GasLimitToggleBurnRoleGlobally)).
		WithAddDataMovementGas(true).
		WithDataParts(parts).
		Build()
}

func (t *tokenManagementTransactionFactory) CreateTransactionForSettingSpecialRoleOnFungibleToken(
	sender core.Address,
	user core.Address,
	tokenIdentifier string,
	addRoleLocalMint bool,
	addRoleLocalBurn bool,
) (*core.Transaction, error) {
	parts := []string{
		"setSpecialRole",
		hex.EncodeToString([]byte(tokenIdentifier)),
		user.ToHex(),
	}

	if addRoleLocalMint {
		parts = append(parts, hex.EncodeToString([]byte("ESDTRoleLocalMint")))
	}
	if addRoleLocalBurn {
		parts = append(parts, hex.EncodeToString([]byte("ESDTRoleLocalBurn")))
	}

	return core.NewTransactionBuilder().
		WithConfig(t.Config).
		WithSender(sender).
		WithReceiver(t.Config.ESDTContractAddress).
		WithAmount(big.NewInt(0)).
		WithProvidedGasLimit(uint64(t.Config.GasLimitSetSpecialRole)).
		WithAddDataMovementGas(true).
		WithDataParts(parts).
		Build()
}

func (t *tokenManagementTransactionFactory) CreateTransactionForSettingSpecialRoleOnSemiFungibleToken(
	sender core.Address,
	user core.Address,
	tokenIdentifier string,
	addRoleNFTCreate bool,
	addRoleNFTBurn bool,
	addRoleNFTAddQuantity bool,
	addRoleESDTTransferRole bool,
) (*core.Transaction, error) {
	parts := []string{
		"setSpecialRole",
		hex.EncodeToString([]byte(tokenIdentifier)),
		user.ToHex(),
	}

	if addRoleNFTCreate {
		parts = append(parts, hex.EncodeToString([]byte("ESDTRoleNFTCreate")))
	}
	if addRoleNFTBurn {
		parts = append(parts, hex.EncodeToString([]byte("ESDTRoleNFTBurn")))
	}
	if addRoleNFTAddQuantity {
		parts = append(parts, hex.EncodeToString([]byte("ESDTRoleNFTAddQuantity")))
	}
	if addRoleESDTTransferRole {
		parts = append(parts, hex.EncodeToString([]byte("ESDTTransferRole")))
	}

	return core.NewTransactionBuilder().
		WithConfig(t.Config).
		WithSender(sender).
		WithReceiver(t.Config.ESDTContractAddress).
		WithAmount(big.NewInt(0)).
		WithProvidedGasLimit(uint64(t.Config.GasLimitSetSpecialRole)).
		WithAddDataMovementGas(true).
		WithDataParts(parts).
		Build()
}

func (t *tokenManagementTransactionFactory) CreateTransactionForSettingSpecialRoleOnNonFungibleToken(
	sender core.Address,
	user core.Address,
	tokenIdentifier string,
	addRoleNFTCreate bool,
	addRoleNFTBurn bool,
	addRoleNFTUpdateAttributes bool,
	addRoleNFTAddURI bool,
	addRoleESDTTransferRole bool,
) (*core.Transaction, error) {
	parts := []string{
		"setSpecialRole",
		hex.EncodeToString([]byte(tokenIdentifier)),
		user.ToHex(),
	}

	if addRoleNFTCreate {
		parts = append(parts, hex.EncodeToString([]byte("ESDTRoleNFTCreate")))
	}
	if addRoleNFTBurn {
		parts = append(parts, hex.EncodeToString([]byte("ESDTRoleNFTBurn")))
	}
	if addRoleNFTUpdateAttributes {
		parts = append(parts, hex.EncodeToString([]byte("ESDTRoleNFTUpdateAttributes")))
	}
	if addRoleNFTAddURI {
		parts = append(parts, hex.EncodeToString([]byte("ESDTRoleNFTAddURI")))
	}
	if addRoleESDTTransferRole {
		parts = append(parts, hex.EncodeToString([]byte("ESDTTransferRole")))
	}

	return core.NewTransactionBuilder().
		WithConfig(t.Config).
		WithSender(sender).
		WithReceiver(t.Config.ESDTContractAddress).
		WithAmount(big.NewInt(0)).
		WithProvidedGasLimit(uint64(t.Config.GasLimitSetSpecialRole)).
		WithAddDataMovementGas(true).
		WithDataParts(parts).
		Build()
}

func (t *tokenManagementTransactionFactory) CreateTransactionForCreatingNFT(
	sender core.Address,
	tokenIdentifier string,
	initialQuantity int,
	name string,
	royalties int,
	hash string,
	attributes []byte,
	URIs []string,
) (*core.Transaction, error) {
	parts := []string{
		"ESDTNFTCreate",
		hex.EncodeToString([]byte(tokenIdentifier)),
		hex.EncodeToString(core.EncodeUnsignedNumber(uint64(initialQuantity))),
		hex.EncodeToString([]byte(name)),
		hex.EncodeToString(core.EncodeUnsignedNumber(uint64(royalties))),
		hex.EncodeToString([]byte(hash)),
		hex.EncodeToString(attributes),
	}

	for _, uri := range URIs {
		parts = append(parts, hex.EncodeToString([]byte(uri)))
	}

	nftData := name + hash + hex.EncodeToString(attributes) + strings.Join(URIs, "")
	storageGasLimit := len(nftData) * t.Config.GasLimitStorePerByte

	return core.NewTransactionBuilder().
		WithConfig(t.Config).
		WithSender(sender).
		WithReceiver(sender).
		WithAmount(big.NewInt(0)).
		WithProvidedGasLimit(uint64(t.Config.GasLimitESDTNFTCreate + storageGasLimit)).
		WithAddDataMovementGas(true).
		WithDataParts(parts).
		Build()
}

func (t *tokenManagementTransactionFactory) CreateTransactionForPausing(
	sender core.Address,
	tokenIdentifier string,
) (*core.Transaction, error) {
	parts := []string{
		"pause",
		hex.EncodeToString([]byte(tokenIdentifier)),
	}

	return core.NewTransactionBuilder().
		WithConfig(t.Config).
		WithSender(sender).
		WithReceiver(t.Config.ESDTContractAddress).
		WithAmount(big.NewInt(0)).
		WithProvidedGasLimit(uint64(t.Config.GasLimitPausing)).
		WithAddDataMovementGas(true).
		WithDataParts(parts).
		Build()
}

func (t *tokenManagementTransactionFactory) CreateTransactionForUnpausing(
	sender core.Address,
	tokenIdentifier string,
) (*core.Transaction, error) {
	parts := []string{
		"unPause",
		hex.EncodeToString([]byte(tokenIdentifier)),
	}

	return core.NewTransactionBuilder().
		WithConfig(t.Config).
		WithSender(sender).
		WithReceiver(t.Config.ESDTContractAddress).
		WithAmount(big.NewInt(0)).
		WithProvidedGasLimit(uint64(t.Config.GasLimitPausing)).
		WithAddDataMovementGas(true).
		WithDataParts(parts).
		Build()
}

func (t *tokenManagementTransactionFactory) CreateTransactionForFreezing(
	sender core.Address,
	user core.Address,
	tokenIdentifier string,
) (*core.Transaction, error) {
	parts := []string{
		"freeze",
		hex.EncodeToString([]byte(tokenIdentifier)),
		user.ToHex(),
	}

	return core.NewTransactionBuilder().
		WithConfig(t.Config).
		WithSender(sender).
		WithReceiver(t.Config.ESDTContractAddress).
		WithAmount(big.NewInt(0)).
		WithProvidedGasLimit(uint64(t.Config.GasLimitFreezing)).
		WithAddDataMovementGas(true).
		WithDataParts(parts).
		Build()
}

func (t *tokenManagementTransactionFactory) CreateTransactionForUnfreezing(
	sender core.Address,
	user core.Address,
	tokenIdentifier string,
) (*core.Transaction, error) {
	parts := []string{
		"unFreeze",
		hex.EncodeToString([]byte(tokenIdentifier)),
		user.ToHex(),
	}

	return core.NewTransactionBuilder().
		WithConfig(t.Config).
		WithSender(sender).
		WithReceiver(t.Config.ESDTContractAddress).
		WithAmount(big.NewInt(0)).
		WithProvidedGasLimit(uint64(t.Config.GasLimitFreezing)).
		WithAddDataMovementGas(true).
		WithDataParts(parts).
		Build()
}

func (t *tokenManagementTransactionFactory) CreateTransactionForLocalMinting(
	sender core.Address,
	tokenIdentifier string,
	supplyToMint int,
) (*core.Transaction, error) {
	parts := []string{
		"ESDTLocalMint",
		hex.EncodeToString([]byte(tokenIdentifier)),
		hex.EncodeToString(core.EncodeUnsignedNumber(uint64(supplyToMint))),
	}

	return core.NewTransactionBuilder().
		WithConfig(t.Config).
		WithSender(sender).
		WithReceiver(t.Config.ESDTContractAddress).
		WithAmount(big.NewInt(0)).
		WithProvidedGasLimit(uint64(t.Config.GasLimitESDTLocalMint)).
		WithAddDataMovementGas(true).
		WithDataParts(parts).
		Build()
}

func (t *tokenManagementTransactionFactory) CreateTransactionForLocalBurning(
	sender core.Address,
	tokenIdentifier string,
	supplyToBurn int,
) (*core.Transaction, error) {
	parts := []string{
		"ESDTLocalBurn",
		hex.EncodeToString([]byte(tokenIdentifier)),
		hex.EncodeToString(core.EncodeUnsignedNumber(uint64(supplyToBurn))),
	}

	return core.NewTransactionBuilder().
		WithConfig(t.Config).
		WithSender(sender).
		WithReceiver(t.Config.ESDTContractAddress).
		WithAmount(big.NewInt(0)).
		WithProvidedGasLimit(uint64(t.Config.GasLimitESDTLocalBurn)).
		WithAddDataMovementGas(true).
		WithDataParts(parts).
		Build()
}

func (t *tokenManagementTransactionFactory) CreateTransactionForUpdatingAttributes(
	sender core.Address,
	tokenIdentifier string,
	tokenNonce int,
	attributes []byte,
) (*core.Transaction, error) {
	parts := []string{
		"ESDTNFTUpdateAttributes",
		hex.EncodeToString([]byte(tokenIdentifier)),
		hex.EncodeToString(core.EncodeUnsignedNumber(uint64(tokenNonce))),
		hex.EncodeToString(attributes),
	}

	return core.NewTransactionBuilder().
		WithConfig(t.Config).
		WithSender(sender).
		WithReceiver(sender).
		WithAmount(big.NewInt(0)).
		WithProvidedGasLimit(uint64(t.Config.GasLimitESDTNFTUpdateAttributes)).
		WithAddDataMovementGas(true).
		WithDataParts(parts).
		Build()
}

func (t *tokenManagementTransactionFactory) CreateTransactionForAddingQuantity(
	sender core.Address,
	tokenIdentifier string,
	tokenNonce int,
	quantityToAdd int,
) (*core.Transaction, error) {
	parts := []string{
		"ESDTNFTAddQuantity",
		hex.EncodeToString([]byte(tokenIdentifier)),
		hex.EncodeToString(core.EncodeUnsignedNumber(uint64(tokenNonce))),
		hex.EncodeToString(core.EncodeUnsignedNumber(uint64(quantityToAdd))),
	}

	return core.NewTransactionBuilder().
		WithConfig(t.Config).
		WithSender(sender).
		WithReceiver(sender).
		WithAmount(big.NewInt(0)).
		WithProvidedGasLimit(uint64(t.Config.GasLimitESDTNFTAddQuantity)).
		WithAddDataMovementGas(true).
		WithDataParts(parts).
		Build()
}

func (t *tokenManagementTransactionFactory) CreateTransactionForBurningQuantity(
	sender core.Address,
	tokenIdentifier string,
	tokenNonce int,
	quantityToAdd int,
) (*core.Transaction, error) {
	parts := []string{
		"ESDTNFTBurn",
		hex.EncodeToString([]byte(tokenIdentifier)),
		hex.EncodeToString(core.EncodeUnsignedNumber(uint64(tokenNonce))),
		hex.EncodeToString(core.EncodeUnsignedNumber(uint64(quantityToAdd))),
	}

	return core.NewTransactionBuilder().
		WithConfig(t.Config).
		WithSender(sender).
		WithReceiver(sender).
		WithAmount(big.NewInt(0)).
		WithProvidedGasLimit(uint64(t.Config.GasLimitESDTNFTAddQuantity)).
		WithAddDataMovementGas(true).
		WithDataParts(parts).
		Build()
}

func notifyAboutUnsettingBurnRoleGlobally() {
	log := logger.GetOrCreate("tokenManagementTransactionFactory")
	log.Info(`
==========
IMPORTANT!
==========
You are about to issue (register) a new token. This will set the role "ESDTRoleBurnForAll" (globally).
Once the token is registered, you can unset this role by calling "unsetBurnRoleGlobally" (in a separate transaction).
`)
}
