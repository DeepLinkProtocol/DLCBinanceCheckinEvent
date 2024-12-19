package utils

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	ValidWalletAddress = "0x81725c71a98EF3B2b08Cd1fc6Ba052E6568e3A48"
)

func SignMessage(message string, privateKeyStr string) (string, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return "", err
	}

	prefix := "\x19Ethereum Signed Message:\n"
	prefixLen := strconv.Itoa(len(message))
	prefixedMessage := prefix + prefixLen + message

	data := []byte(prefixedMessage)
	hash := crypto.Keccak256Hash(data)

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	return common.Bytes2Hex(signature), nil
}

func VerifySignature(message, signatureStr, expectedWalletAddress string) (bool, error) {
	signatureStr = strings.TrimPrefix(signatureStr, "0x")

	// Decode the hex-encoded signature
	signature := common.Hex2Bytes(signatureStr)

	// Remove the recovery ID (v) from the signature for ECDSA
	if len(signature) != 65 {
		return false, fmt.Errorf("invalid signature length: %d", len(signature))
	}

	// Extract `v` value
	v := signature[64]
	if v < 27 {
		v += 27 // Ensure compatibility with Ethereum's recovery ID
	}
	signature[64] = v - 27 // Standardize v to 0 or 1 for `crypto.SigToPub`

	// Compute the Ethereum-specific hash
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(prefixedMessage))

	// Recover the public key from the signature
	pubKey, err := crypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		return false, fmt.Errorf("failed to recover public key: %v", err)
	}

	// Derive the wallet address from the recovered public key
	recoveredAddress := crypto.PubkeyToAddress(*pubKey)
	expectedAddress := common.HexToAddress(expectedWalletAddress)

	return recoveredAddress == expectedAddress, nil
}
