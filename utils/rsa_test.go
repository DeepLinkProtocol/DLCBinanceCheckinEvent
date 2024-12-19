package utils

import (
	"testing"
)

const (
	publicKeyStr = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCbWoXkbbwfcZnLW43Vsh1YMu1W5a4reIHvcMYqFjWJl4huA7JK" +
		"ZdC/O3pmEqxdSGZPkerDoN70yfFUPJwKHF+Zc30CWSHTgN+ivR1W4EwyQd48b7WfdU6NVNu2p0p9B2dvcytsdIZ+FKjDwjXplw21//9zX7x" +
		"Lr2rF+YeP1mp20QIDAQAB"

	signatureStr = "VI6k2ILEFuB2ltAIYHrEeFjlxq4ZMHdoPTMLxFyHrg1ylnMpFJo2J/YStRKRdEh0Pv+beVWje" +
		"0Nz+rZ6z3RzPFFwFkgEGK4XT3PGnpYnZXWvvCBHhQg0OmypNftzktUxcekbazWvF4BSTxoFlIDYBdAt5L69lUnwY7GZ9pOXGoU="
	parameterStr = "a=b&c=[\"1\",\"2\",\"3\"]&recvWindow=5000&timestamp=1499827319559"
)

func TestGetPublicKey(t *testing.T) {
	_, err := GetPublicKey(publicKeyStr)
	if err != nil {
		t.Fatalf("failed to parse public key: %v", err)
	}
}

func TestVerifySignature(t *testing.T) {
	publicKey, err := GetPublicKey(publicKeyStr)
	if err != nil {
		t.Fatalf("failed to get public key: %v", err)
	}

	valid, err := RSAVerify(parameterStr, publicKey, signatureStr)
	if err != nil {
		t.Fatalf("verification failed: %v", err)
	}

	if !valid {
		t.Fatalf("expected signature to be valid, but it is not")
	}
}

func TestSignAndVerifySignature(t *testing.T) {
	privateKeyBase64, publicKeyBase64, err := generateRSAKeyPair(2048)
	t.Logf("privateKeyBase64 : %s \n", privateKeyBase64)
	t.Log("--------- ")
	t.Logf("publicKeyBase64 : %s \n", publicKeyBase64)
	if err != nil {
		t.Fatalf("Error generating RSA key pair: %v", err)
	}

	privateKey, err := getPrivateKey(privateKeyBase64)
	if err != nil {
		t.Fatalf("Failed to get private key: %v", err)
	}

	publicKey, err := GetPublicKey(publicKeyBase64)
	if err != nil {
		t.Fatalf("Failed to get public key: %v", err)
	}

	data := "walletAddress=0x5de8477A8A47e7F2c5cE05ad4532861a0AaAc909"
	signature, err := generateSignature(privateKey, data)
	if err != nil {
		t.Fatalf("Failed to generate signature: %v", err)
	}

	isValid, err := RSAVerify(data, publicKey, signature)
	if err != nil {
		t.Fatalf("Failed to verify signature: %v", err)
	}

	if !isValid {
		t.Errorf("Signature is invalid")
	}
}

func TestSignAndVerifySignature1(t *testing.T) {
	privateKeyBase64 := `MIIEpAIBAAKCAQEArpUnWG0bprc6+6gfK52L+jX9dlXa1RTY4hZLUelSVYlk/szlFsBUVGkt+vFQq0HmoHDgg/TTadJ4odxkStyn2+q1AngAzWPRspCmstDhC96lwaH1h9bhdxEqAiR3nvxHu/pviuBfn9SKn4fBHqgubYS2IMQwDCDpjAbql7f24QTw5MtwQ6nCqLzz81BMH9x/y3c4ylfxvze48UJqKS9VLeIhIK6wBB7NjpWwTqzVwVveg/+RAvcp3Gb/tXqeNgx8Ofv+Nr9BvLKY1+E1LT2l0sAgt1AUFP3x6AhTaPp0MDWn0Uq+xNsOGSeHMu9BNS90TdlUmecUkseA7yB2M6m5MwIDAQABAoIBAQCq9QQCY1Wge+0WYhuj6jMYYaZCKvCPmEqJmqtHGuO1P0XW/W2YSd5KinSsN4J3MaFVmzNABI7CIYhfeCH97PWzNLLsJ/chKY4+/cc/c6vso6pNYvu8eX9vyS8JygwZc5ILUcHIjM8XBFp6vcUu8CIGvN5cV8F1HwWUSXGQAGWREla7upPv6pPA57gPItAK7ZfVS8O16LItGWkzPaYcqjfySE6b92Mb9c9R1dd3iB5mz5jtVpDidUTJWaUs2e0EUtntFWS8rgAEf6KArAUBks6e+GRy7KpohVst85EMJbXpOINZwdVwIWXlTFgW7LYb0PH9shKzQM17Z0h2oUofkZSBAoGBAOD0YxgV9zR3sEEqyUXE4h4S7O5axO2XTMDo94P7ctPxQJwBGznj4AuZOUCGznwgWT/vm5DbsQL3ShASU++fr5DD0AQIt9AGkghjtGHtJ9SU7YovPHDK5RjYKVSlNh9nFIG926qZILgsphhlrHrm5UrTeSrz9rNeU5AzZh/uxvAXAoGBAMatIaDVLcku9iHwLtKaTSar89IOYkkw3vHa0n0lcBap9hV6LtBcRLZP1/a3VvRoV5IlZuE58zcMTrmdqnCvQUpuxVcZm
CXAFdBXRLyaLJoXDZytaNqrHHicKYvRZf2KtyZiz5WpsKySHwagkIZR1fJGeFnAKCjMl2lOW438A/VFAoGBAMyCwbg6+AQWMAH+4P7x8FRBBm/ny5Lo81mKMsQljI5MjV2Gz/bASYah5V/Zbs2AJ0OuFTML15CHuyiDURXPijBFJM8WEe4omwjPhEVm8sgcIRx4ty0f+Emu81xF2r+P2h/duGAPWKS3ysTxYm5Vje2J9mVraERHpBa788NiNDA/AoGAYMTZOU48Q96U0gj7tWakp25fjvOkmcOtA54yofQHOXLFQsbFYIVgjnArX6cDOn5MEQoYpyEjvq9G9Q+/ga1LHub/RaJYwiJiPZ0UBM0PZmpOHf80sDVh47kkX06535meBZthQqNpQ1TUudShMFtR2vTKD+URanXkVc1tuKWEhAUCgYAM1/sF8FmnTV455iywVSw3GgMs4+pGnBlbwn46cJzRcua2fJGYKxt8QineOmPmhqMgYcKxHbT68vczXdBKV/lyqk2ZpJ/MjOtl0M3r+Cqx5gkPGD0g6ysA+f1lybrSseE5enguAUIuAeWbCeJYW1HkWPtQmJx2EgAJ4qcejOu5rQ==`
	publicKeyBase64 := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArpUnWG0bprc6+6gfK52L+jX9dlXa1RTY4hZLUelSVYlk/szlFsBUVGkt+vFQq0HmoHDgg/TTadJ4odxkStyn2+q1AngAzWPRspCmstDhC96lwaH1h9bhdxEqAiR3nvxHu/pviuBfn9SKn4fBHqgubYS2IMQwDCDpjAbql7f24QTw5MtwQ6nCqLzz81BMH9x/y3c4ylfxvze48UJqKS9VLeIhIK6wBB7NjpWwTqzVwVveg/+RAvcp3Gb/tXqeNgx8Ofv+Nr9BvLKY1+E1LT2l0sAgt1AUFP3x6AhTaPp0MDWn0Uq+xNsOGSeHMu9BNS90TdlUmecUkseA7yB2M6m5MwIDAQAB"
	t.Logf("privateKeyBase64 : %s \n", privateKeyBase64)
	t.Log("--------- ")
	t.Logf("publicKeyBase64 : %s \n", publicKeyBase64)

	// privateKey, err := getPrivateKey(privateKeyBase64)
	// if err != nil {
	// 	t.Fatalf("Failed to get private key: %v", err)
	// }

	publicKey, err := GetPublicKey(publicKeyBase64)
	if err != nil {
		t.Fatalf("Failed to get public key: %v", err)
	}

	data := "walletAddress=0x5de8477A8A47e7F2c5cE05ad4532861a0AaAc909"
	signature := "iUsbAbfM0FT1Kk6frovXMuOHuLwOf8fipuUNHouVOTxg+qRJB2fmsOwjnMdxRn+Dwsm5L03joOyrb+o0mHHRPFQNow2I0xfyhkhw6NFc94Vg+9Q4bYYK3subTdBb1d/ybDFKmxsDeHS7SkykG6F5kwvvJggNABJDxIRJoHf9Z+yaFVncYydfKpUtNcL7vH83LLFd9nXRrFEo9V6u6MvD7dhdvpbm7+kDGsb0fhmn0q3oMKu3MS1ZZj7JgSajsr38CgoKkbfDnr4gUKR5A0X7UVFxUI3wm8beEJA6ndDg6IviVgUzvOFjiHBsiCbiiqLnvPHdcVtGQN2tC9xDClTV1g=="
	t.Logf("signature : %s \n", signature)

	isValid, err := RSAVerify(data, publicKey, signature)
	if err != nil {
		t.Fatalf("Failed to verify signature: %v", err)
	}

	if !isValid {
		t.Errorf("Signature is invalid")
	}
}
