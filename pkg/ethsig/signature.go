package ethsig

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeebo/errs"
)

// ErrSignatureParams indicates that an error occurred while creating a signature.
var ErrSignatureParams = errs.Class("create signature error")

// SignatureParams entity describes values for create signature.
type SignatureParams struct {
	Values     [][]byte
	PrivateKey *ecdsa.PrivateKey
}

// Signature defines signature type.
type Signature string

// EthereumSignedMessage defines message for signature.
const EthereumSignedMessage string = "\x19Ethereum Signed Message:\n"

// MetaTransactionFucnDescription defines name and parameters MetaTransaction function.
const MetaTransactionFucnDescription string = "MetaTransaction(uint256 nonce,address from,bytes functionSignature)"

// EthereumSignedMessageForApprove defines signature message for approve.
const EthereumSignedMessageForApprove string = "\x19\x01"

// SignHash is a function that calculates a hash for the given message.
func SignHash(data []byte) []byte {
	msg := fmt.Sprintf(EthereumSignedMessage+"%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

// SignWallet generates signature for user's wallet.
func SignWallet(addressWallet Address, addressContract Address, privateKey *ecdsa.PrivateKey) (Signature, error) {
	var values [][]byte
	if err := addressWallet.IsValidAddress(); err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}
	if err := addressContract.IsValidAddress(); err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	addressWalletByte, err := hex.DecodeString(string(addressWallet)[LengthHexPrefix:])
	if err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	addressContractByte, err := hex.DecodeString(string(addressContract)[LengthHexPrefix:])
	if err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	values = append(values, addressWalletByte, addressContractByte)
	createSignature := SignatureParams{
		Values:     values,
		PrivateKey: privateKey,
	}

	signatureByte, err := Sign(createSignature)
	if err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	signature, err := NormalizeSignature(signatureByte)

	return signature, ErrSignatureParams.Wrap(err)
}

// SignWalletWithValue generates signature for user's wallet with value.
func SignWalletWithValue(addressWallet Address, addressContract Address, value int64, privateKey *ecdsa.PrivateKey) (Signature, error) {
	var values [][]byte
	if err := addressWallet.IsValidAddress(); err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}
	if err := addressContract.IsValidAddress(); err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	addressWalletByte, err := hex.DecodeString(string(addressWallet)[LengthHexPrefix:])
	if err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	addressContractByte, err := hex.DecodeString(string(addressContract)[LengthHexPrefix:])
	if err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	valueStringWithZeros := CreateHexStringFixedLength(fmt.Sprintf("%x", value))
	valueByte, err := hex.DecodeString(string(valueStringWithZeros))
	if err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	values = append(values, addressWalletByte, addressContractByte, valueByte)
	createSignature := SignatureParams{
		Values:     values,
		PrivateKey: privateKey,
	}

	signatureByte, err := Sign(createSignature)
	if err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	signature, err := NormalizeSignature(signatureByte)

	return signature, ErrSignatureParams.Wrap(err)
}

// SignWalletWithNonce generates signature for user's wallet with value and nonce.
func SignWalletWithNonce(addressWallet Address, addressContract Address, value *big.Int, nonce int64, privateKey *ecdsa.PrivateKey) (Signature, error) {
	var values [][]byte
	if err := addressWallet.IsValidAddress(); err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}
	if err := addressContract.IsValidAddress(); err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	addressWalletByte, err := hex.DecodeString(string(addressWallet)[LengthHexPrefix:])
	if err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	addressContractByte, err := hex.DecodeString(string(addressContract)[LengthHexPrefix:])
	if err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	valueStringWithZeros := CreateHexStringFixedLength(fmt.Sprintf("%x", value))
	valueByte, err := hex.DecodeString(string(valueStringWithZeros))
	if err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	nonceStringWithZeros := CreateHexStringFixedLength(fmt.Sprintf("%x", nonce))
	nonceByte, err := hex.DecodeString(string(nonceStringWithZeros))
	if err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	values = append(values, addressWalletByte, addressContractByte, valueByte, nonceByte)
	createSignature := SignatureParams{
		Values:     values,
		PrivateKey: privateKey,
	}

	signatureByte, err := Sign(createSignature)
	if err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	signature, err := NormalizeSignature(signatureByte)

	return signature, ErrSignatureParams.Wrap(err)
}

// SignNFTTransaction generates signature for user's wallet with tokenID and value.
func SignNFTTransaction(addressWallet Address, addressSaleContract Address, addressNFTContract Address, tokenID int64, value *big.Int, privateKey *ecdsa.PrivateKey) (Signature, error) {
	var values [][]byte
	if err := addressWallet.IsValidAddress(); err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}
	if err := addressSaleContract.IsValidAddress(); err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}
	if err := addressNFTContract.IsValidAddress(); err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	addressWalletByte, err := hex.DecodeString(string(addressWallet)[LengthHexPrefix:])
	if err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	addressSaleContractByte, err := hex.DecodeString(string(addressSaleContract)[LengthHexPrefix:])
	if err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	addressNFTContractByte, err := hex.DecodeString(string(addressNFTContract)[LengthHexPrefix:])
	if err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	tokenIDStringWithZeros := CreateHexStringFixedLength(fmt.Sprintf("%x", tokenID))
	tokenIDByte, err := hex.DecodeString(string(tokenIDStringWithZeros))
	if err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	valueStringWithZeros := CreateHexStringFixedLength(fmt.Sprintf("%x", value))
	valueByte, err := hex.DecodeString(string(valueStringWithZeros))
	if err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	values = append(values, addressWalletByte, addressSaleContractByte, addressNFTContractByte, tokenIDByte, valueByte)
	createSignature := SignatureParams{
		Values:     values,
		PrivateKey: privateKey,
	}

	signatureByte, err := Sign(createSignature)
	if err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	signature, err := NormalizeSignature(signatureByte)

	return signature, ErrSignatureParams.Wrap(err)
}

// SignERC20Approval generates signature for user's wallet for approve ERC20.
func SignERC20Approval(contractMethodAddress Hex, to Address, value *big.Int, nonce int64, from Address, domainSeperator Hex, privateKey *ecdsa.PrivateKey) (Signature, error) {
	if err := to.IsValidAddress(); err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}
	if err := from.IsValidAddress(); err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	toStringWithZeros := CreateHexStringFixedLength(string(to[LengthHexPrefix:]))
	valueMoneyStringWithZeros := CreateHexStringFixedLength(fmt.Sprintf("%x", value))
	data := contractMethodAddress + toStringWithZeros + valueMoneyStringWithZeros
	dataByte, err := hex.DecodeString(string(data[LengthHexPrefix:]))
	if err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	functionSignatureBytes := crypto.Keccak256Hash(dataByte)

	var metaTransaction []byte
	metaTransactionFucnDescriptionBytes := crypto.Keccak256([]byte(MetaTransactionFucnDescription))
	metaTransaction = append(metaTransaction, metaTransactionFucnDescriptionBytes...)

	nonceStringWithZeros := CreateHexStringFixedLength(fmt.Sprintf("%x", nonce))
	nonceByte, err := hex.DecodeString(string(nonceStringWithZeros))
	if err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}
	metaTransaction = append(metaTransaction, nonceByte...)

	fromStringWithZeros := CreateHexStringFixedLength(string(from[LengthHexPrefix:]))
	fromByte, err := hex.DecodeString(string(fromStringWithZeros))
	if err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}
	metaTransaction = append(metaTransaction, fromByte...)
	metaTransaction = append(metaTransaction, functionSignatureBytes.Bytes()...)
	metaTransactionBytes := crypto.Keccak256Hash(metaTransaction)

	domainSeperatorStringWithZeros := CreateHexStringFixedLength(string(domainSeperator[LengthHexPrefix:]))
	domainSeperatorByte, err := hex.DecodeString(string(domainSeperatorStringWithZeros))
	if err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	dataSignature := crypto.Keccak256Hash([]byte(EthereumSignedMessageForApprove), domainSeperatorByte, metaTransactionBytes.Bytes())

	signatureByte, err := crypto.Sign(dataSignature.Bytes(), privateKey)
	if err != nil {
		return "", ErrSignatureParams.Wrap(err)
	}

	signature, err := NormalizeSignature(signatureByte)

	return signature, ErrSignatureParams.Wrap(err)
}

// Sign makes signature.
func Sign(createSignature SignatureParams) ([]byte, error) {
	var allValues []byte
	for _, value := range createSignature.Values {
		allValues = append(allValues, value...)
	}
	dataSignature := SignHash(crypto.Keccak256Hash(allValues).Bytes())
	signature, err := crypto.Sign(dataSignature, createSignature.PrivateKey)
	return signature, ErrSignatureParams.Wrap(err)
}

// NormalizeSignature reforms last two byte of signature from 00, 01 to 1b, 1c.
func NormalizeSignature(signatureByte []byte) (Signature, error) {
	signatureWithoutEnd := string(signatureByte)[:len(signatureByte)-1]
	signatureString := hex.EncodeToString(signatureByte)
	signatureLastSymbol := signatureString[len(signatureString)-1:]

	if signatureLastSymbol == fmt.Sprintf("%d", PrivateKeyVZero) {
		return Signature(hex.EncodeToString(append([]byte(signatureWithoutEnd), []byte{byte(PrivateKeyVTwentySeven)}...))), nil
	}

	if signatureLastSymbol == fmt.Sprintf("%d", PrivateKeyVOne) {
		return Signature(hex.EncodeToString(append([]byte(signatureWithoutEnd), []byte{byte(PrivateKeyVTwentyEight)}...))), nil
	}

	return "", ErrSignatureParams.New("error private key format")
}
