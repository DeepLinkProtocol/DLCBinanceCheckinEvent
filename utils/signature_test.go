package utils

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestSignMessage(t *testing.T) {
	privateKey, _ := crypto.GenerateKey()
	publicKey := privateKey.PublicKey
	address := crypto.PubkeyToAddress(publicKey)
	privateKeyBytes := privateKey.D.Bytes()
	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	message := "hello world"
	signature, err := SignMessage(message, privateKeyHex)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(signature)

	ok1, _ := VerifyEthereumSignature(message, signature, address.Hex())
	if !ok1 {
		t.Fatal("verify1 signature failed")
	}
	ok2, _ := VerifyEthereumSignature(message, "0x"+signature, address.Hex())
	if !ok2 {
		t.Fatal("verify2 signature failed")
	}
}

func TestSignMessage1(t *testing.T) {
	privateKeyHex := "3674391010b4526c30a71d6174966badc7f76aed93507420f473179e7da9d70b"
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		t.Fail()
	}
	// privateKeyBytes := privateKey.D.Bytes()
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	t.Log("address: ", address.Hex())
	message := "hello world"
	signature, err := SignMessage(message, privateKeyHex)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("signature: ", signature)

	// ok := VerifySignature(message, signature, address.Hex())

	// if !ok {
	// 	t.Fatal("verify signature failed")
	// }
}
