package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
)

const (
	RSAPublicKey = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArpUnWG0bprc6+6gfK52L+jX9dlXa1RTY4hZLUelSVYlk/szlFsBUVGkt+vFQq0HmoHDgg/TTadJ4odxkStyn2+q1AngAzWPRspCmstDhC96lwaH1h9bhdxEqAiR3nvxHu/pviuBfn9SKn4fBHqgubYS2IMQwDCDpjAbql7f24QTw5MtwQ6nCqLzz81BMH9x/y3c4ylfxvze48UJqKS9VLeIhIK6wBB7NjpWwTqzVwVveg/+RAvcp3Gb/tXqeNgx8Ofv+Nr9BvLKY1+E1LT2l0sAgt1AUFP3x6AhTaPp0MDWn0Uq+xNsOGSeHMu9BNS90TdlUmecUkseA7yB2M6m5MwIDAQAB"
)

func RSAVerify(parameterStr string, publicKey *rsa.PublicKey, sign string) (bool, error) {
	// Decode the base64 encoded signature
	signature, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %v", err)
	}

	// Hash the parameter string
	hashed := sha256.Sum256([]byte(parameterStr))

	// Verify the signature using the public key and hashed data
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return false, fmt.Errorf("verification failed: %v", err)
	}

	// Return true if verification succeeds
	return true, nil
}

func GetPublicKey(publicKey string) (*rsa.PublicKey, error) {
	// Decode the Base64 public key
	decodedKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, errors.New("failed to decode Base64 public key: " + err.Error())
	}

	// Parse the decoded public key to an x509 format
	parsedKey, err := x509.ParsePKIXPublicKey(decodedKey)
	if err != nil {
		return nil, errors.New("failed to parse public key: " + err.Error())
	}

	// Assert the key is an RSA public key
	rsaPubKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not a valid RSA public key")
	}

	return rsaPubKey, nil
}

// getPrivateKey decodes a Base64-encoded private key string and parses it as an RSA private key.
func getPrivateKey(privateKeyStr string) (*rsa.PrivateKey, error) {
	// Decode the Base64-encoded private key
	decodedKey, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		return nil, errors.New("failed to decode Base64 private key: " + err.Error())
	}

	// Parse the decoded key as PKCS#8 private key
	parsedKey, err := x509.ParsePKCS1PrivateKey(decodedKey)
	if err != nil {
		return nil, errors.New("failed to parse private key: " + err.Error())
	}

	// // Assert that the parsed key is an RSA private key
	// rsaKey, ok := parsedKey.(*rsa.PrivateKey)
	// if !ok {
	// 	return nil, errors.New("parsed key is not an RSA private key")
	// }

	return parsedKey, nil
}

func generateSignature(privateKey *rsa.PrivateKey, data string) (string, error) {
	// Hash the data
	hashed := sha256.Sum256([]byte(data))

	// Sign the hashed data using the RSA private key
	signature, err := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %v", err)
	}

	// Return the signature as a Base64-encoded string
	return base64.StdEncoding.EncodeToString(signature), nil
}

func generateRSAKeyPair(bits int) (privateKeyBase64 string, publicKeyBase64 string, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate RSA key pair: %v", err)
	}

	publicKey := &privateKey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBase64 = base64.StdEncoding.EncodeToString(privateKeyBytes)

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal public key: %v", err)
	}
	publicKeyBase64 = base64.StdEncoding.EncodeToString(publicKeyBytes)

	return privateKeyBase64, publicKeyBase64, nil
}
